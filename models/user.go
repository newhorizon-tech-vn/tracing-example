package models

import "github.com/newhorizon-tech-vn/tracing-example/models/entities"

func GetUser(userId int) (*entities.User, error) {
	user := &entities.User{}
	err := DBConnection.Model(&entities.User{}).First(user, userId).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}
