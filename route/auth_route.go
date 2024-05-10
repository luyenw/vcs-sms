package route

import (
	"vcs-sms/config/sql"
	"vcs-sms/controller"
	"vcs-sms/service"
)

func (r *Router) InitAuthRoute() {
	authController := controller.NewAuthController(service.NewUserService(sql.GetPostgres()), service.NewJWTService())
	authRouter := r.Group("/auth")
	authRouter.POST("/register", authController.Register)
	authRouter.POST("/login", authController.Login)
}
