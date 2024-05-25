package route

import (
	"vcs-sms/config/cache"
	"vcs-sms/config/sql"
	"vcs-sms/controller"
	SCOPE "vcs-sms/enum"
	"vcs-sms/middleware"
	"vcs-sms/service"
)

func (r *Router) InitServerRoute() {
	serverController := controller.NewServerController(
		service.NewServerService(sql.GetPostgres()),
		service.NewCacheService(cache.GetRedis()),
		service.NewXLSXService())

	serverRouter := r.Group("/servers")
	serverRouter.GET("", middleware.TokenAuthorization(), middleware.CheckScope(SCOPE.API_SERVER_READ), serverController.GetServer)
	serverRouter.GET("/export", middleware.TokenAuthorization(), middleware.CheckScope(SCOPE.API_SERVER_READ), serverController.ExportServers)
	serverRouter.POST("", middleware.TokenAuthorization(), middleware.CheckScope(SCOPE.API_SERVER_READ), serverController.CreateServer)
	serverRouter.POST("/import", middleware.TokenAuthorization(), middleware.CheckScope(SCOPE.API_SERVER_READ_WRITE), serverController.ImportServers)
	serverRouter.PATCH("/:id", middleware.TokenAuthorization(), middleware.CheckScope(SCOPE.API_SERVER_READ_WRITE), serverController.UpdateServer)
	serverRouter.DELETE("/:id", middleware.TokenAuthorization(), middleware.CheckScope(SCOPE.API_SERVER_READ_WRITE), serverController.DeleteServer)
}
