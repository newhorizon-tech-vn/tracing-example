package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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

	"github.com/Depado/ginprom"
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

	log.InitLogger(&log.Configuration{
		JSONFormat:      true,
		LogLevel:        viper.GetString("log.console_level"),
		StacktraceLevel: viper.GetString("log.stacktrace_level"),
		File: &log.FileConfiguration{
			Filename:   viper.GetString("log.file_name"),
			MaxSize:    viper.GetInt("log.file_max_size"),
			MaxAge:     viper.GetInt("log.file_max_age"),
			MaxBackups: viper.GetInt("log.file_max_backups"),
		},
		Console: &log.ConsoleConfiguration{},
	})

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

	p := ginprom.New(
		ginprom.Engine(router),
		ginprom.Path("/metrics"),
		ginprom.Namespace("tracing_example_namespace"),
		ginprom.Subsystem("tracing_example_subsystem"),
	)
	router.Use(p.Instrument())

	router.GET("/v1/class/:classId", authorize.Auth(), h.CheckClass)
	router.GET("/v1/user/:userId", h.GetUser)
	router.POST("/v1/user", h.CreateUser)
	router.PUT("/v1/user/:userId", h.EditUser)

	go StartMetrics()
	go StartMetricsVector()

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

func StartMetrics() {
	totalRequest := promauto.NewCounter(prometheus.CounterOpts{
		Name: "total_request",
		Help: "Total request",
	})

	requestProcessed := promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "request_processed",
			Help:    "Request time",
			Buckets: prometheus.DefBuckets,
		},
	)

	for {
		totalRequest.Inc()
		start := time.Now()
		t := (time.Now().Unix() % 10) * 100
		time.Sleep(time.Duration(t) * time.Millisecond)
		requestProcessed.Observe(time.Since(start).Seconds())
	}
}

func StartMetricsVector() {
	totalRequest := promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "total_http_request",
		Help: "HTTP request",
	},
		[]string{"api", "method", "code"},
	)

	requestProcessed := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_request_processed",
		Help: "HTTP time",
	},
		[]string{"api", "method", "code"},
	)

	for {
		t := (time.Now().Unix() % 10) * 100
		start := time.Now()

		api, method, code := simulateMetrics(t)
		totalRequest.WithLabelValues(api, method, strconv.FormatInt(code, 10)).Inc()

		time.Sleep(time.Duration(t) * time.Millisecond)
		requestProcessed.WithLabelValues(api, method, strconv.FormatInt(code, 10)).Observe(time.Since(start).Seconds())
	}
}

func simulateMetrics(t int64) (string, string, int64) {
	method := "GET"
	if t%2 == 0 {
		method = "POST"
	}

	api := "v1/users"
	if t%4 == 0 {
		api = "v1/classes"
	}

	return api, method, t % 500
}
