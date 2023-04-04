package models

import (
	"context"

	"github.com/newhorizon-tech-vn/tracing-example/models/entities"
)

func GetOrganizationWithClasses(ctx context.Context, organizationId int) (*entities.Organization, error) {
	result := &entities.Organization{}
	err := DBConnection.WithContext(ctx).Preload("Schools.Classes").First(&result, organizationId).Error
	return result, err
}
