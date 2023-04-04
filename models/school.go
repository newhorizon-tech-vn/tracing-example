package models

import (
	"context"
	"time"

	"github.com/newhorizon-tech-vn/tracing-example/models/entities"
)

func GetUserOrgSchoolsWithClassesByUserId(ctx context.Context, userId int) ([]*entities.UserOrgSchool, error) {
	var result []*entities.UserOrgSchool
	err := DBConnection.WithContext(ctx).Preload("School.Classes").Where("UserId = ?", userId).Find(&result).Error
	return result, err
}

func GetOrganizationByUserId(ctx context.Context, userId int) (*entities.Organization, error) {
	result := &entities.UserOrgSchool{}
	err := DBConnection.WithContext(ctx).Preload("Organization").Where("UserId = ? AND EndDate > ?", userId, time.Now().UTC()).First(&result).Error
	return result.GetOrganization(), err
}
