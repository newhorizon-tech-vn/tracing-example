package kafka

import (
	"context"
	"fmt"
	"time"
)

type MessageHandleFunc func(context.Context, string, []byte) error

// WorkTicket ...
type WorkTicket struct {
	EnqueueTime time.Time
	Topic       string
	Data        []byte
}

// Worker ...
type Worker struct {
	ID int

	//chan chan for send and receiver workticket
	WorkItem    chan WorkTicket
	WorkerQueue chan chan WorkTicket

	//chan for quit worker
	QuitChan chan bool

	// handler
	Handler MessageHandleFunc
}

// NewWorker ...
func NewWorker(id int, workerQueue chan chan WorkTicket, handler MessageHandleFunc) Worker {
	return Worker{
		ID:          id,
		WorkerQueue: workerQueue,
		WorkItem:    make(chan WorkTicket),
		QuitChan:    make(chan bool),
		Handler:     handler,
	}
}

// Start ...
func (w *Worker) Start() {
	go func() {
		for {
			//put to chan chan for the first time
			w.WorkerQueue <- w.WorkItem

			select {
			case workItem := <-w.WorkItem:
				go w.processWorkItem(workItem)

			case <-w.QuitChan:
				fmt.Printf("INFO: worker is shutdowning\n")
				return
			}
		}
	}()
}

func (w *Worker) processWorkItem(item WorkTicket) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("ERROR: worker is recovering occur %v \n", r)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	w.Handler(ctx, item.Topic, item.Data)

}

// Stop ...
func (w *Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}

type ConsumerProcessor struct {
	NumberWorker int

	//WorkerQueueVar ...
	WorkerQueueVar chan chan WorkTicket

	//WorkItemVar ...
	WorkItemVar chan WorkTicket

	//Workers ...
	Workers []Worker

	// Handler ...
	Handler MessageHandleFunc
}

// SendTicket ...
func (c *ConsumerProcessor) SendTicket(item WorkTicket) {
	c.WorkItemVar <- item
}

// StopWorkerDispatcher ...
func (c *ConsumerProcessor) StopWorkerDispatcher() {
	for i := 0; i < len(c.Workers); i++ {
		fmt.Printf("INFO: stop worker dispatche id=%d \n", i+1)

		//stop worker
		c.Workers[i].Stop()
	}
}

// StartWorkerDispatcher ...
func (c *ConsumerProcessor) StartWorkerDispatcher() {
	//initial global var
	c.WorkerQueueVar = make(chan chan WorkTicket, c.NumberWorker)
	c.WorkItemVar = make(chan WorkTicket, c.NumberWorker)

	//instantialize workers
	c.Workers = make([]Worker, c.NumberWorker)

	for i := 0; i < c.NumberWorker; i++ {
		// log.Bg().Debug("start worker", zap.Int("id", i+1))

		//start worker
		worker := NewWorker(i, c.WorkerQueueVar, c.Handler)
		worker.Start()

		//keep worker for shutdown all workers
		c.Workers[i] = worker
	}

	//start receive workItem
	go func() {
		for {
			select {
			case workItemRecv := <-c.WorkItemVar:
				// log.Bg().Debug("received work item", zap.String("data", workItemRecv.Data))
				go func() {
					worker := <-c.WorkerQueueVar
					// log.Bg().Info("dispatching work item", zap.String("data", workItemRecv.Data))
					worker <- workItemRecv
				}()
			}
		}
	}()
}
