package main

import (
	"vcs-sms/route"

	"github.com/gin-gonic/gin"
)

func main() {
	router := route.NewRouter(gin.Default())
	router.InitServerRoute()
	router.InitAuthRoute()
	router.InitReportRoute()
	// router.InitHealthCheckRoute()
	router.InitUserRoute()

	router.Run(":8081")
}
