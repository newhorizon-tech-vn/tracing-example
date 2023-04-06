package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/newhorizon-tech-vn/tracing-example/pkg/http/client"
	"github.com/spf13/viper"
)

var userClient *client.Client

func GetUser(ctx context.Context, userId int) (*User, error) {
	if userClient == nil {
		userClient = client.DefaultClient()
	}

	bytes, statusCode, err := userClient.Get(ctx, fmt.Sprintf("http://localhost:%d/v1/user/%d", viper.GetInt("simulator.port"), userId))
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP Code %d", statusCode)
	}

	user := &User{}
	if err = json.Unmarshal(bytes, user); err != nil {
		return nil, err
	}

	return user, nil
}
