package models

import (
	"context"

	"github.com/newhorizon-tech-vn/tracing-example/models/entities"
)

func GetClassessByTeacherId(ctx context.Context, userId int) ([]*entities.ClassTeacher, error) {
	var result []*entities.ClassTeacher
	err := DBConnection.WithContext(ctx).Preload("Class").Where("user_UserId = ?", userId).Find(&result).Error
	return result, err
}

func GetClassessByStudentId(ctx context.Context, userId int) ([]*entities.ClassStudent, error) {
	var result []*entities.ClassStudent
	err := DBConnection.WithContext(ctx).Preload("Class").Where("user_UserId = ?", userId).Find(&result).Error
	return result, err
}
