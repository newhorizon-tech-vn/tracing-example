package v1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/newhorizon-tech-vn/tracing-example/pkg/log"
)

func (h *Handler) CreateUser(c *gin.Context) {
	log.For(c.Request.Context()).Debug("[create-user] start process")
	time.Sleep(100 * time.Millisecond)
	c.JSON(http.StatusOK, nil)
}
