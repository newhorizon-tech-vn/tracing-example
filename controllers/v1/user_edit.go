package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) EditUser(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}
