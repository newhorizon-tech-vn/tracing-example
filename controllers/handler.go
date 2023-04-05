package controllers

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/newhorizon-tech-vn/tracing-example/controllers/v1"
)

type IHandler interface {
	CheckClass(c *gin.Context)
	GetUser(c *gin.Context)
}

func MakeHandler(version string) IHandler {
	return &v1.Handler{}
}
