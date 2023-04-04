package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/newhorizon-tech-vn/tracing-example/models/entities"
)

const KeyUserInfo = "nhz:dms:api:user:%d:info"

func SetUser(ctx context.Context, user *entities.User, expireTime time.Duration) error {
	key := fmt.Sprintf(KeyUserInfo, user.UserId)
	value := user.ToJSON()
	return redisClient.Set(ctx, key, value, expireTime).Err()
}

func GetUser(ctx context.Context, userId int) (*entities.User, error) {
	key := fmt.Sprintf(KeyUserInfo, userId)
	value := redisClient.Get(ctx, key)
	if value.Err() != nil {
		return nil, value.Err()
	}

	user := &entities.User{}
	err := user.FromJSON(value.Val())
	return user, err
}

func DeleteUser(ctx context.Context, userId int) error {
	return redisClient.Del(ctx, fmt.Sprintf(KeyUserInfo, userId)).Err()
}
