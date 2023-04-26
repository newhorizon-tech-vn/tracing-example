package cache

import (
	"context"
	"fmt"

	"github.com/newhorizon-tech-vn/tracing-example/pkg/log"
	"github.com/newhorizon-tech-vn/tracing-example/pkg/util"
	"go.uber.org/zap"
)

const KeyUserClassIds = "nhz:dms:api:user:%d:classids"

func SetClassIdsOfUser(ctx context.Context, userId int, classIds []int) error {
	key := fmt.Sprintf(KeyUserClassIds, userId)
	redisClient.Del(ctx, key)

	var value []interface{}
	for _, val := range classIds {
		value = append(value, val)
	}
	return redisClient.SAdd(ctx, key, value...).Err()
}

func GetClassIdsOfUser(ctx context.Context, userId int) ([]int, error) {
	key := fmt.Sprintf(KeyUserClassIds, userId)
	rs := redisClient.SMembers(ctx, key)
	if rs.Err() != nil {
		log.For(ctx).Error("[get-cache-cache-class-ids-of-user] query failed", zap.Int("user_id", userId), zap.String("key", key), zap.Error(rs.Err()))
		return nil, rs.Err()
	}

	var ids []int
	for _, val := range rs.Val() {
		id, err := util.StringToInt(val)
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}

	log.For(ctx).Debug("[get-cache-cache-class-ids-of-user] process success", zap.Int("user_id", userId), zap.Reflect("result", ids))
	return ids, nil
}

func DeleteClassIdsOfUser(ctx context.Context, userId int) error {
	return redisClient.Del(ctx, fmt.Sprintf(KeyUserClassIds, userId)).Err()
}
