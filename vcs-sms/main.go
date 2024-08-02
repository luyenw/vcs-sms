package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"vcs-sms/config"
	"vcs-sms/config/logger"
	"vcs-sms/config/mq"
	"vcs-sms/config/rpc"
	"vcs-sms/config/storage"
	"vcs-sms/route"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	config.InitConfig()

	cli := storage.GetGCPStorage()
	fmt.Printf("Client: %+v\n", cli)
	p := mq.GetProducer()
	fmt.Printf("Producer: %+v\n", p)
	rpcClient := rpc.GetRpcClient()
	fmt.Println(rpcClient)

	fmt.Println("Hello, World!")
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
	router.Static("static", "./api-spec")
	config := cors.DefaultConfig()
	config.AddAllowHeaders("Authorization")
	config.AddAllowHeaders("Content-Type")
	config.AllowCredentials = true
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	router.Use(cors.New(config))

	router.InitServerRoute()
	router.InitAuthRoute()
	router.InitReportRoute()
	router.InitUserRoute()
	router.Run(":8081")
}
