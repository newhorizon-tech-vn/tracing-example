package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/newhorizon-tech-vn/tracing-example/middleware/authorize"
	"github.com/newhorizon-tech-vn/tracing-example/services"
)

type Handler struct{}

func (h *Handler) CheckClass(c *gin.Context) {
	_, err := strconv.Atoi(c.Param("classId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	userId, exists := c.Get(authorize.KeyUserId)
	if exists == false {
		c.JSON(http.StatusUnauthorized, nil)
		return
	}

	roleId, exists := c.Get(authorize.KeyRoleId)
	if exists == false {
		c.JSON(http.StatusNonAuthoritativeInfo, nil)
		return
	}

	_, err = (&services.ClassService{}).GetClassIds(c.Request.Context(), userId.(int), roleId.(int))
	/*
		index := slices.IndexFunc(classIds, func(id int) bool { return id == classId })
		if index < 0 {
			c.JSON(http.StatusInternalServerError, nil)
			return
		}
	*/
	fmt.Println(err)
	c.JSON(http.StatusOK, nil)
}
