package cache

import (
	"context"
	"encoding/json"
	"time"

	redisotel "github.com/go-redis/redis/extra/redisotel/v8"
	redis "github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

type Ints []int

func (s Ints) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s Ints) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &s)
}

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
	redisClient.AddHook(redisotel.NewTracingHook())

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()

	return redisClient.Ping(ctx).Err()
}
