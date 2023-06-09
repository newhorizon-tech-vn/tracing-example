package v1

import (
	"net/http"
	"strconv"
	"time"

	"github.com/newhorizon-tech-vn/tracing-example/services"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/newhorizon-tech-vn/tracing-example/pkg/log"
)

func (h *Handler) GetUser(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	log.For(c.Request.Context()).Debug("[get-user] start process")
	t := (time.Now().Unix() % 10) * 100
	time.Sleep(time.Duration(t) * time.Millisecond)

	user, err := (&services.UserService{UserId: userId}).GetUser(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	log.Debug("start process", zap.Int("userId", userId))
	c.JSON(http.StatusOK, user)
}
