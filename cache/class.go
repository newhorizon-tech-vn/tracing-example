package cache

import (
	"context"
	"fmt"

	"github.com/newhorizon-tech-vn/tracing-example/pkg/util"
)

const KeyUserClassIds = "nhz:dms:api:user:%d:classids"

func SetClassIdsOfUser(ctx context.Context, userId int, classIds []int) error {
	key := fmt.Sprintf(KeyUserClassIds, userId)
	redisClient.Del(ctx, key)

	return redisClient.SAdd(ctx, key, classIds).Err()
}

func GetClassIdsOfUser(ctx context.Context, userId int) ([]int, error) {
	key := fmt.Sprintf(KeyUserClassIds, userId)
	rs := redisClient.SMembers(ctx, key)
	if rs.Err() != nil {
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
	return ids, nil
}

func DeleteClassIdsOfUser(ctx context.Context, userId int) error {
	return redisClient.Del(ctx, fmt.Sprintf(KeyUserClassIds, userId)).Err()
}
