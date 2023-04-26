package producers

import (
	"context"

	"github.com/newhorizon-tech-vn/tracing-example/pkg/kafka"
	"github.com/spf13/viper"
)

const (
	DefaultProduceTimeout        int = 5000
	DefaultProduceFlushFrequency int = 500
)

var (
	producer *kafka.Producer
)

func SetProducer(p *kafka.Producer) {
	producer = p
}

func GetProducerConfig() *kafka.ProducerConfig {
	viper.SetDefault("kafka.produce_timeout", DefaultProduceTimeout)
	viper.SetDefault("kafka.produce_flush_frequency", DefaultProduceFlushFrequency)

	return &kafka.ProducerConfig{
		Brokers:        viper.GetStringSlice("kafka.brokers"),
		Timeout:        viper.GetInt("kafka.produce_timeout"),
		FlushFrequency: viper.GetInt("kafka.produce_flush_frequency"),
		Version:        viper.GetString("kafka.version"),
	}
}

func ProduceMessage(ctx context.Context, topic, key string, data interface{}) error {
	return producer.Produce(ctx, topic, key, data)
}
