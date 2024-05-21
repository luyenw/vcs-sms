package route

import (
	"time"
	"vcs-sms/config/cache"
	"vcs-sms/config/elasticsearch"
	"vcs-sms/config/sql"
	"vcs-sms/controller"
	SCOPE "vcs-sms/enum"
	"vcs-sms/middleware"
	"vcs-sms/repo"
	"vcs-sms/service"
)

func (r *Router) InitReportRoute() {
	reportController := controller.NewReportController(
		service.NewReportService(service.NewESService(&repo.ESClient{Client: elasticsearch.GetESClient()}),
			service.NewRegisteredMailService(),
			service.NewServerService(sql.GetPostgres()),
			service.NewCacheService(cache.GetRedis())),
	)
	reportController.PeriodicReport(24 * time.Hour)
	reportRouter := r.Group("/report")

	reportRouter.POST("/", middleware.TokenAuthorization(), middleware.CheckScope(SCOPE.API_REPORT_READ), reportController.SendReport)
}
