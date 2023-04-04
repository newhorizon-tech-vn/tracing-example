package services

import (
	"context"
	"fmt"

	"github.com/newhorizon-tech-vn/tracing-example/cache"
	"github.com/newhorizon-tech-vn/tracing-example/models"
	"github.com/newhorizon-tech-vn/tracing-example/models/entities"
	"github.com/newhorizon-tech-vn/tracing-example/pkg/util"
	"golang.org/x/exp/maps"
)

type ClassService struct {
	ClassId int
}

func (*ClassService) GetClassIds(ctx context.Context, userId, userRoleId int) ([]int, error) {

	// try to get from cache
	var (
		result []int
	)

	// try to get data from storage
	if userRoleId == entities.ROLE_STUDENT {
		rs, err := models.GetClassessByStudentId(ctx, userId)
		if err != nil {
			return nil, err
		}

		result = util.MapFunc(rs, func(c *entities.ClassStudent) int {
			return c.ClassId
		})
	} else if userRoleId == entities.ROLE_TEACHER {
		rs, err := models.GetClassessByTeacherId(ctx, userId)
		if err != nil {
			return nil, err
		}

		result = util.MapFunc(rs, func(c *entities.ClassTeacher) int {
			return c.ClassId
		})

	} else if userRoleId == entities.ROLE_ORGANIZATION_ADMIN {
		organization, err := models.GetOrganizationByUserId(ctx, userId)
		if err != nil {
			return nil, err
		}

		rs, err := models.GetOrganizationWithClasses(ctx, organization.GetId())
		if err != nil {
			return nil, err
		}

		classIds := make(map[int]bool)
		for _, school := range rs.GetSchools() {
			for _, class := range school.GetClasses() {
				classIds[class.ClassId] = true
			}
		}

		result = maps.Keys(classIds)

	} else if userRoleId == entities.ROLE_SCHOOL_ADMIN {
		rs, err := models.GetUserOrgSchoolsWithClassesByUserId(ctx, userId)
		if err != nil {
			return nil, err
		}

		classIds := make(map[int]bool)
		for _, val := range rs {
			for _, class := range val.GetSchool().GetClasses() {
				classIds[class.ClassId] = true
			}
		}

		result = maps.Keys(classIds)

	} else {
		return nil, fmt.Errorf("role invalid")
	}

	// return result
	return result, nil
}

func (*ClassService) GetClassIdsV2(ctx context.Context, userId, userRoleId int) ([]int, error) {

	// try to get from cache
	result, err := cache.GetClassIdsOfUser(ctx, userId)
	if err != nil {
		return result, nil
	}

	// try to get data from storage
	if userRoleId == entities.ROLE_STUDENT {
		rs, err := models.GetClassessByStudentId(ctx, userId)
		if err != nil {
			return nil, err
		}

		result = util.MapFunc(rs, func(c *entities.ClassStudent) int {
			return c.ClassId
		})
	} else if userRoleId == entities.ROLE_TEACHER {
		rs, err := models.GetClassessByTeacherId(ctx, userId)
		if err != nil {
			return nil, err
		}

		result = util.MapFunc(rs, func(c *entities.ClassTeacher) int {
			return c.ClassId
		})

	} else if userRoleId == entities.ROLE_ORGANIZATION_ADMIN {
		organization, err := models.GetOrganizationByUserId(ctx, userId)
		if err != nil {
			return nil, err
		}

		rs, err := models.GetOrganizationWithClasses(ctx, organization.GetId())
		if err != nil {
			return nil, err
		}

		classIds := make(map[int]bool)
		for _, school := range rs.GetSchools() {
			for _, class := range school.GetClasses() {
				classIds[class.ClassId] = true
			}
		}

		result = maps.Keys(classIds)

	} else if userRoleId == entities.ROLE_SCHOOL_ADMIN {
		rs, err := models.GetUserOrgSchoolsWithClassesByUserId(ctx, userId)
		if err != nil {
			return nil, err
		}

		classIds := make(map[int]bool)
		for _, val := range rs {
			for _, class := range val.GetSchool().GetClasses() {
				classIds[class.ClassId] = true
			}
		}

		result = maps.Keys(classIds)

	} else {
		return nil, fmt.Errorf("role invalid")
	}

	// update cache
	cache.SetClassIdsOfUser(ctx, userId, result)

	// return result
	return result, nil
}
