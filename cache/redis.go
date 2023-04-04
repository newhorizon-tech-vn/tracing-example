package cache

import (
	"context"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var (
	// RedisClient client process Redis commands
	redisClient *redis.Client
)

func InitRedis() error {

	redisClient = redis.NewClient(&redis.Options{
		Addr:       viper.GetString("redis.address"),
		MaxRetries: viper.GetInt("redis.max_retry"),
		Password:   viper.GetString("redis.password"),
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()

	return redisClient.Ping(ctx).Err()
}
