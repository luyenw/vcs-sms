package main

import (
	"fmt"
	"healthcheck-server/config"
	"healthcheck-server/config/logger"
	"healthcheck-server/config/mq"
	"healthcheck-server/config/rpc"
	"healthcheck-server/controller"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	config.InitConfig()
	go func() {
		rpc.ConfigGRPC()
	}()
	log := logger.NewLogger()
	fmt.Println("Hello, World!, running on port 5000")
	log.Info("Server started")
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-signalChan
		log.Fatal("Received signal, shutting down", zap.String("signal", sig.String()))
		os.Exit(0)
	}()

	p := mq.GetProducer()
	fmt.Println(p)

	go func() {
		controller.UpdateServersStatus()
	}()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.POST("/healthcheck", controller.Healthcheck)
	router.GET("/healthcheck", controller.GetHealthcheck)
	router.Run(":5000")
}
