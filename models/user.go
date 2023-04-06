package models

import (
	"context"

	"github.com/newhorizon-tech-vn/tracing-example/models/entities"
)

func GetUser(ctx context.Context, userId int) (*entities.User, error) {
	user := &entities.User{}
	err := DBConnection.WithContext(ctx).Model(&entities.User{}).First(user, userId).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}
