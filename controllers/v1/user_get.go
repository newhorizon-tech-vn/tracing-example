package v1

import (
	"net/http"
	"strconv"
	"time"

	"github.com/newhorizon-tech-vn/tracing-example/clients/user"

	"github.com/gin-gonic/gin"
	"github.com/newhorizon-tech-vn/tracing-example/pkg/log"
)

func (h *Handler) GetUser(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	time.Sleep(100 * time.Millisecond)

	// simulator response
	user := &user.User{
		ID:   1,
		Name: "ABC",
	}

	log.Debug("start process", "userId", userId)
	c.JSON(http.StatusOK, user)
}
