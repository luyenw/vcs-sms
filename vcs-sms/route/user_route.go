package route

import (
	"vcs-sms/config/sql"
	"vcs-sms/controller"
	SCOPE "vcs-sms/enum"
	"vcs-sms/middleware"
	"vcs-sms/service"
)

func (r *Router) InitUserRoute() {
	userController := controller.NewUserController(service.NewUserService(sql.GetPostgres()))
	userRouter := r.Group("/users")
	userRouter.POST("/", middleware.TokenAuthorization(), middleware.CheckScope(SCOPE.API_USER_READ_WRITE), userController.CreateUser)
	userRouter.PUT("/:id/role", middleware.TokenAuthorization(), middleware.CheckScope(SCOPE.API_USER_READ_WRITE), userController.UpdateUserRole)
}
