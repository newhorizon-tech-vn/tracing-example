package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
	"github.com/newhorizon-tech-vn/tracing-example/cache"
	"github.com/newhorizon-tech-vn/tracing-example/controllers"
	"github.com/newhorizon-tech-vn/tracing-example/middleware/authorize"
	"github.com/newhorizon-tech-vn/tracing-example/models"
	"github.com/newhorizon-tech-vn/tracing-example/pkg/log"
	"github.com/newhorizon-tech-vn/tracing-example/pkg/tracing"
	"github.com/newhorizon-tech-vn/tracing-example/setting"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {

	if err := setting.InitSetting(); err != nil {
		log.Fatal("get config failed", "error", err)
		return
	}

	if err := models.InitMySQL(); err != nil {
		log.Fatal("connect to mysql failed", "error", err)
		return
	}

	if err := cache.InitRedis(); err != nil {
		log.Fatal("connect to redis failed", "error", err)
		return
	}

	serviceName := "tracing-example"
	// jaegerEndPoint := "http://127.0.0.1:14268/v1/trace"
	jaegerEndPoint := "http://127.0.0.1:14268/api/traces"
	if _, err := tracing.StartOpenTelemetryHTTP(serviceName, jaegerEndPoint); err != nil {
		log.Fatal("connect to jaeger failed", "error", err)
		return
	}

	h := controllers.MakeHandler("v1")

	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(otelgin.Middleware(serviceName))
	router.GET("/v1/class/:classId", authorize.Auth(), h.CheckClass)

	router.Run(fmt.Sprintf("localhost:%d", viper.GetInt("setting.port")))
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server ...")
}

func GetAlbums(c *gin.Context) {

}

func PostAlbums(c *gin.Context) {

}
