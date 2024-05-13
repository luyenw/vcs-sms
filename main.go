package main

import (
	"io"
	"os"
	"vcs-sms/route"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	router := route.NewRouter(gin.Default())
	router.InitServerRoute()
	router.InitAuthRoute()
	router.InitReportRoute()
	router.InitHealthCheckRoute()
	router.InitUserRoute()

	log.SetFormatter(&log.JSONFormatter{})
	if _, err := os.Stat("/var/log/vcs-sms"); os.IsNotExist(err) {
		os.Mkdir("/var/log/vcs-sms", 0755)
	}
	logFile, err := os.OpenFile("/var/log/vcs-sms/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(io.MultiWriter(logFile, os.Stdout))
	} else {
		log.Fatal("Failed to log to file, using default stderr")
	}
	router.Run(":8081")
}
