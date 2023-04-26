package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"go.opentelemetry.io/otel/trace"
)

const (
	DefaultKafkaProduceTimeout        int = 5000
	DefaultKafkaProduceFlushFrequency int = 500
)

type ProducerConfig struct {
	Brokers        []string
	Timeout        int // In mili second
	FlushFrequency int // In mili second
	Version        string
}

type Producer struct {
	Producer sarama.AsyncProducer
	Tracer   trace.Tracer
}

// AsyncProduceWithTopic msg to ZPI Kafka
func (p *Producer) AsyncProduce(key string, data []byte, topic string) {
	if p == nil {
		return
	}
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(data),
	}
	p.Producer.Input() <- msg
}

// Close connection
func (p *Producer) Close() error {
	if p == nil {
		return nil
	}

	if err := p.Producer.Close(); err != nil {
		return err
	}
	return nil
}

// produce to async DB writer
func NewProducer(cfg *ProducerConfig) (*Producer, error) {

	// config
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true // bắt buộc
	// config.Producer.Return.Errors = true

	var err error
	config.Version, err = sarama.ParseKafkaVersion(cfg.Version)
	if err != nil {
		fmt.Printf("ERROR: kafka version [%v] is invalid: %s \n", cfg.Version, err.Error())
		return nil, err
	}

	if cfg.Timeout > 0 {
		config.Producer.Timeout = time.Duration(cfg.Timeout) * time.Millisecond
	} else {
		fmt.Printf("WARN: kafka produce not set timeout \n")
	}

	if cfg.FlushFrequency > 0 {
		config.Producer.Flush.Frequency = time.Duration(cfg.FlushFrequency) * time.Millisecond
	} else {
		fmt.Printf("WARN: kafka produce not set flush frequency \n")
	}

	// new async producer
	asynProducer, err := sarama.NewAsyncProducer(cfg.Brokers, config)
	if err != nil {
		fmt.Printf("ERROR: init kafka producer failed: brokers %v error %s\n", cfg.Brokers, err.Error())
		return nil, err
	}
	// monitor notify error
	go func(producer sarama.AsyncProducer) {
		errors := producer.Errors()
		success := producer.Successes()
		for {
			select {
			case err := <-errors:
				if err != nil {
					fmt.Printf("ERROR: produce has error failed: brokers %v error %s\n", cfg.Brokers, err.Error())
				}
			case rs := <-success:
				fmt.Printf("DEBUG: produce has error failed: topic %s partition %d offset %d\n", rs.Topic, rs.Partition, rs.Offset)
			}
		}
	}(asynProducer)

	return &Producer{
		Producer: asynProducer,
	}, nil
}

// ProduceKafka produce data to kafka write DB
func (p *Producer) Produce(ctx context.Context, topic string, key string, data interface{}) error {
	if p == nil {
		return fmt.Errorf("nil pointer")
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("ERROR: produce recovering from panic recover %v\n", r)
		}
	}()

	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Key:       sarama.ByteEncoder(key),
		Value:     sarama.ByteEncoder(bytes),
		Timestamp: time.Now(),
	}

	p.Producer.Input() <- msg

	return nil
}
