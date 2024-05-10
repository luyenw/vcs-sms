package route

import (
	"vcs-sms/config/elasticsearch"
	"vcs-sms/config/sql"
	"vcs-sms/controller"
	"vcs-sms/service"
)

func (r *Router) InitHealthCheckRoute() {
	healthCheckController := controller.NewHealthCheckController(service.NewHealthCheckService(), service.NewESService(elasticsearch.GetESClient()), service.NewServerService(sql.GetPostgres()))
	healthCheckController.HealthCheck()
}
