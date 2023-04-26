package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"
	"go.uber.org/zap"

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
		log.Fatal("get config failed", zap.Error(err))
		return
	}

	if err := models.InitMySQL(); err != nil {
		log.Fatal("connect to mysql failed", zap.Error(err))
		return
	}

	if err := cache.InitRedis(); err != nil {
		log.Fatal("connect to redis failed", zap.Error(err))
		return
	}

	log.InitLogger(viper.GetString("log.console_level"), viper.GetString("log.stacktrace_level"))

	/*
		// create kafka consumer
		c, err := kafka.NewConsumer(consumers.GetConsumerConfig())
		if err != nil {
			log.Fatal("start kafka consumer failed", zap.Error(err))
			return
		}
		go c.Start()

		// create kafka producer
		producer, err := kafka.NewProducer(producers.GetProducerConfig())
		if err != nil {
			log.Fatal("create kafka producer failed", zap.Error(err))
			return
		}
		producers.SetProducer(producer)

		err = producers.ProduceMessage(context.Background(), viper.GetString("kafka.topic"), "key-1", "abcdef")
		if err != nil {
			log.Fatal("produce message failed", zap.Error(err))
			return
		}
	*/

	serviceName := viper.GetString("jaeger.name")
	// jaegerEndPoint := "http://127.0.0.1:14268/v1/trace"
	// jaegerEndPoint := "http://127.0.0.1:14268/api/traces"
	/*
		if _, err := tracing.StartOpenTelemetry(serviceName, viper.GetString("jaeger.endpoint")); err != nil {
			log.Fatal("connect to jaeger failed", zap.Error(err))
			return
		}
	*/

	if _, err := tracing.StartOpenTelemetryByUDP(serviceName, viper.GetString("jaeger.udp_host"), viper.GetString("jaeger.udp_port")); err != nil {
		log.Fatal("connect to jaeger udp failed", zap.Error(err))
		return
	}

	h := controllers.MakeHandler("v1")

	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(otelgin.Middleware(serviceName))
	router.GET("/v1/class/:classId", authorize.Auth(), h.CheckClass)
	router.GET("/v1/user/:userId", h.GetUser)

	router.Run(fmt.Sprintf("localhost:%d", viper.GetInt("setting.port")))
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server ...")
}

func startChildService() {
	serviceName := "child-service"
	if _, err := tracing.StartOpenTelemetry(serviceName, viper.GetString("jaeger.endpoint")); err != nil {
		log.Fatal("connect to jaeger failed", zap.Error(err))
		return
	}

	h := controllers.MakeHandler("v1")

	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(otelgin.Middleware(serviceName))
	router.GET("/v1/user/:userId", h.GetUser)

	router.Run(fmt.Sprintf("localhost:%d", viper.GetInt("simulator.port")))
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server ...")
}
