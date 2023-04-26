package services

import (
	"context"
	"fmt"

	"github.com/newhorizon-tech-vn/tracing-example/cache"
	"github.com/newhorizon-tech-vn/tracing-example/clients/user"
	"github.com/newhorizon-tech-vn/tracing-example/models"
	"github.com/newhorizon-tech-vn/tracing-example/models/entities"
	"github.com/newhorizon-tech-vn/tracing-example/pkg/log"
	"github.com/newhorizon-tech-vn/tracing-example/pkg/util"
	"go.uber.org/zap"
	"golang.org/x/exp/maps"
)

type ClassService struct {
	ClassId int
}

func (*ClassService) GetClassIds(ctx context.Context, userId, userRoleId int) ([]int, error) {
	log.For(ctx).Debug("[get-class-ids] start process", zap.Int("user_id", userId), zap.Int("user_role_id", userRoleId))
	// test code other service
	_, err := user.GetUser(ctx, userId)
	if err != nil {
		log.For(ctx).Error("[get-class-ids] get user failed", zap.Int("user_id", userId), zap.Int("user_role_id", userRoleId), zap.Error(err))
		return nil, err
	}

	// try to get from cache
	result, err := cache.GetClassIdsOfUser(ctx, userId)
	if err != nil {
		log.For(ctx).Error("[get-class-ids] get cache failed", zap.Int("user_id", userId), zap.Int("user_role_id", userRoleId), zap.Error(err))
		return result, nil
	}

	// try to get data from storage
	if userRoleId == entities.ROLE_STUDENT {
		log.For(ctx).Debug("[get-class-ids] trying to process with student role", zap.Int("user_id", userId), zap.Int("user_role_id", userRoleId))
		rs, err := models.GetClassessByStudentId(ctx, userId)
		if err != nil {
			log.For(ctx).Error("[get-class-ids] get classess by student id failed", zap.Int("user_id", userId), zap.Int("user_role_id", userRoleId), zap.Error(err))
			return nil, err
		}

		result = util.MapFunc(rs, func(c *entities.ClassStudent) int {
			return c.ClassId
		})
	} else if userRoleId == entities.ROLE_TEACHER {
		log.For(ctx).Debug("[get-class-ids] trying to process with teacher role", zap.Int("user_id", userId), zap.Int("user_role_id", userRoleId))
		rs, err := models.GetClassessByTeacherId(ctx, userId)
		if err != nil {
			log.For(ctx).Error("[get-class-ids] get classess by teacher id failed", zap.Int("user_id", userId), zap.Int("user_role_id", userRoleId), zap.Error(err))
			return nil, err
		}

		result = util.MapFunc(rs, func(c *entities.ClassTeacher) int {
			return c.ClassId
		})

	} else if userRoleId == entities.ROLE_ORGANIZATION_ADMIN {
		log.For(ctx).Debug("[get-class-ids] trying to process with organization admin role", zap.Int("user_id", userId), zap.Int("user_role_id", userRoleId))
		organization, err := models.GetOrganizationByUserId(ctx, userId)
		if err != nil {
			log.For(ctx).Error("[get-class-ids] get organization info failed", zap.Int("user_id", userId), zap.Int("user_role_id", userRoleId), zap.Error(err))
			return nil, err
		}

		rs, err := models.GetOrganizationWithClasses(ctx, organization.GetId())
		if err != nil {
			log.For(ctx).Error("[get-class-ids] get org with classes failed", zap.Int("user_id", userId), zap.Int("user_role_id", userRoleId), zap.Error(err))
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
		log.For(ctx).Debug("[get-class-ids] trying to process with school admin role", zap.Int("user_id", userId), zap.Int("user_role_id", userRoleId))
		rs, err := models.GetUserOrgSchoolsWithClassesByUserId(ctx, userId)
		if err != nil {
			log.For(ctx).Error("[get-class-ids] get user org schools with classes by userid failed", zap.Int("user_id", userId), zap.Int("user_role_id", userRoleId), zap.Error(err))
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
		log.For(ctx).Error("[get-class-ids] role invalid", zap.Int("user_id", userId), zap.Int("user_role_id", userRoleId))
		return nil, fmt.Errorf("role invalid")
	}

	// update cache
	cache.SetClassIdsOfUser(ctx, userId, result)

	// return result
	log.For(ctx).Error("[get-class-ids] process success", zap.Int("user_id", userId), zap.Int("user_role_id", userRoleId), zap.Reflect("result", result))
	return result, nil
}
