package main

import (
	"fmt"
	"healthcheck-worker/config"
	"healthcheck-worker/config/cache"
	"healthcheck-worker/config/elasticsearch"
	"healthcheck-worker/config/logger"
	"healthcheck-worker/config/mq"
	"healthcheck-worker/config/sql"
	"healthcheck-worker/controller"
	"healthcheck-worker/repo"
	"healthcheck-worker/service"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	config.InitConfig()
	fmt.Println("Hello, World!")

	c := mq.GetConsumer()
	c.SubscribeTopics([]string{"healthcheck-topic"}, nil)
	fmt.Println(c)

	log := logger.NewLogger()
	log.Info("Server started")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-signalChan
		log.Fatal("Received signal, shutting down", zap.String("signal", sig.String()))
		os.Exit(0)
	}()

	healthCheckController := controller.NewHealthCheckController(
		service.NewHealthCheckService(),
		service.NewESService(&repo.ESClient{Client: elasticsearch.GetESClient()}),
		service.NewServerService(sql.GetPostgres()),
		service.NewCacheService(cache.GetRedis()),
	)

	done := make(chan interface{})
	healthCheckController.HealthCheck()
	<-done
}
