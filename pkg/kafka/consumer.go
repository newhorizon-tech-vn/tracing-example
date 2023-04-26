package kafka

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Shopify/sarama"
)

/* =================== Note: please check Kafka version ===================
To check kafka version, run command

jps -m

output:
3854 Kafka /zserver/kafka_2.12-1.0.0/config/server.properties
Here 2.12 is scala version and 1.0.0 is kafka version
=========================================================================== */

type ConsumerConfig struct {
	Version      string
	Brokers      []string
	GroupID      string
	Topics       []string
	NumberWorker int
	Handler      MessageHandleFunc
}

type Consumer struct {
	Brokers   []string
	GroupID   string
	Topics    []string
	Config    *sarama.Config
	Ready     chan bool
	Version   sarama.KafkaVersion
	Processor *ConsumerProcessor
}

// NewConsumer init kafka consumer
func NewConsumer(cfg *ConsumerConfig) (*Consumer, error) {

	if len(cfg.Brokers) < 1 {
		return nil, fmt.Errorf("Brokers invalid")
	}

	if len(cfg.Topics) < 1 {
		return nil, fmt.Errorf("Topics invalid")
	}

	if cfg.Version == "" {
		return nil, fmt.Errorf("Version invalid")
	}

	if cfg.GroupID == "" {
		return nil, fmt.Errorf("GroupID invalid")
	}

	if cfg.NumberWorker < 1 {
		return nil, fmt.Errorf("Worker number invalid")
	}

	config := sarama.NewConfig()
	//config.Version = sarama.V2_1_0_0
	config.Consumer.Return.Errors = true //required
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	// config.Consumer.Offsets.Initial = sarama.OffsetOldest

	var err error
	config.Version, err = sarama.ParseKafkaVersion(cfg.Version)
	if err != nil {
		fmt.Printf("ERROR: kafka version [%v] is invalid: %s \n", cfg.Version, err.Error())
		return nil, err
	}

	return &Consumer{
		Brokers: cfg.Brokers,
		GroupID: cfg.GroupID,
		Topics:  cfg.Topics,
		Config:  config,
		Processor: &ConsumerProcessor{
			NumberWorker: cfg.NumberWorker,
			Handler:      cfg.Handler,
		},
	}, nil
}

// Start start kafka consumer
func (c *Consumer) Start() error {
	c.Ready = make(chan bool)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := sarama.NewConsumerGroup(c.Brokers, c.GroupID, c.Config)
	if err != nil {
		return err
	}

	// Start worker
	go c.Processor.StartWorkerDispatcher()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			if err := client.Consume(ctx, c.Topics, c); err != nil {
				fmt.Printf("ERROR: error from consumer: %v", err)
			}

			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			c.Ready = make(chan bool)
		}
	}()

	<-c.Ready // Await till the consumer has been set up
	fmt.Printf("INFO: Sarama consumer up and running! ...... \n")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		fmt.Println("terminating: context cancelled")
	case <-sigterm:
		fmt.Println("terminating: via signal")
	}

	cancel()
	wg.Wait()
	if err := client.Close(); err != nil {
		fmt.Printf("ERROR: Error closing client: %s \n", err.Error())
	}

	return nil
}

// Stop stop
func (c *Consumer) Stop() {
	// close(c.Quit)
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// do something
	close(c.Ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29

	for msg := range claim.Messages() {
		c.processMessage(msg)
		session.MarkMessage(msg, "")
	}
	return nil
}

// ProcessMsg process message callback refund
func (c *Consumer) processMessage(msg *sarama.ConsumerMessage) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("ERROR: process-message recovering occur: %v \n", r)
		}
	}()

	ticket := &WorkTicket{
		EnqueueTime: msg.Timestamp,
		Topic:       msg.Topic,
		Data:        msg.Value,
	}

	c.Processor.SendTicket(*ticket)
}
