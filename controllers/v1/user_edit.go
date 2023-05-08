package v1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/newhorizon-tech-vn/tracing-example/pkg/log"
)

func (h *Handler) EditUser(c *gin.Context) {
	log.For(c.Request.Context()).Debug("[edit-user] start process")
	time.Sleep(300 * time.Millisecond)
	c.JSON(http.StatusOK, nil)
}
