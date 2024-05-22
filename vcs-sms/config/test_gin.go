package config

import (
	"github.com/gin-gonic/gin"
)

func GetTestGin() *gin.Engine {
	var testGin *gin.Engine
	gin.SetMode(gin.ReleaseMode)
	testGin = gin.Default()
	return testGin
}
