package consumers

import (
	"context"
	"time"

	"github.com/newhorizon-tech-vn/tracing-example/pkg/kafka"
	"github.com/newhorizon-tech-vn/tracing-example/pkg/log"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func GetConsumerConfig() *kafka.ConsumerConfig {
	viper.SetDefault("kafka.number_worker", 128)
	rs := &kafka.ConsumerConfig{
		Version:      viper.GetString("kafka.version"),
		Brokers:      viper.GetStringSlice("kafka.brokers"),
		GroupID:      viper.GetString("kafka.group_id"),
		NumberWorker: viper.GetInt("kafka.number_worker"),
		Handler:      processMessage,
	}

	rs.Topics = append(rs.Topics, viper.GetString("kafka.topic"))
	return rs
}

func processMessage(ctx context.Context, topic string, data []byte) error {
	log.Debug("process message", zap.String("topic", topic), zap.String("data", string(data)))
	time.Sleep(time.Second)
	return nil
}
