package main

import (
	"healthcheck-service/config/cache"
	"healthcheck-service/config/elasticsearch"
	"healthcheck-service/config/logger"
	"healthcheck-service/config/sql"
	"healthcheck-service/controller"
	"healthcheck-service/service"
	"os"
	"os/signal"
	"syscall"
	"vcs-sms/repo"

	"go.uber.org/zap"
)

func main() {
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
