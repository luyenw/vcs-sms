package main

import (
	"os"
	"os/signal"
	"syscall"
	"vcs-sms/config/logger"
	"vcs-sms/route"

	"github.com/gin-gonic/gin"
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

	gin.SetMode(gin.ReleaseMode)
	router := route.NewRouter(gin.Default())
	router.InitServerRoute()
	router.InitAuthRoute()
	router.InitReportRoute()
	router.InitUserRoute()
	router.Run(":8081")
}
