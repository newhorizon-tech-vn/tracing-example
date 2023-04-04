package entities

import (
	"time"

	"gorm.io/gorm"
)

const (
	ROLE_UNKNOWN = iota // 0
	ROLE_ADMIN
	ROLE_ORGANIZATION_ADMIN
	ROLE_SCHOOL_ADMIN
	ROLE_TEACHER
	ROLE_PARENT
	ROLE_STUDENT
)

type Base struct {
	CreatedAt time.Time `gorm:"column:CreateDate"`
	UpdatedAt time.Time `gorm:"column:UpdateDate"`
	EndDate   time.Time `gorm:"column:EndDate;default:9999-12-31 23:59:59"`
}

type BaseV2 struct {
	UUID      string         `gorm:"column:uuid;type:char;size:36"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;column:deleted_at"`
}
