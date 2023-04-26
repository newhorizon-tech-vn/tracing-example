package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/newhorizon-tech-vn/tracing-example/pkg/http/client"
	"github.com/newhorizon-tech-vn/tracing-example/pkg/log"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var userClient *client.Client

func GetUser(ctx context.Context, userId int) (*User, error) {
	if userClient == nil {
		userClient = client.DefaultClient()
	}

	bytes, statusCode, err := userClient.Get(ctx, fmt.Sprintf("http://localhost:%d/v1/user/%d", viper.GetInt("simulator.port"), userId))
	if err != nil {
		log.For(ctx).Error("[get-user] send request failed", zap.Int("user_id", userId), zap.Int("status_code", statusCode), zap.Error(err))
		return nil, err
	}

	if statusCode != http.StatusOK {
		log.For(ctx).Error("[get-user] http code invalid", zap.Int("user_id", userId), zap.Int("status_code", statusCode))
		return nil, fmt.Errorf("HTTP Code %d", statusCode)
	}

	log.For(ctx).Debug("[get-user] data response", zap.Int("user_id", userId), zap.String("status_code", string(bytes)))
	user := &User{}
	if err = json.Unmarshal(bytes, user); err != nil {
		log.For(ctx).Error("[get-user] unmarshal json failed", zap.Int("user_id", userId), zap.Error(err))
		return nil, err
	}

	return user, nil
}
