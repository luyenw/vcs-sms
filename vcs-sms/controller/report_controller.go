package controller

import (
	"fmt"
	"net/http"
	"time"
	"vcs-sms/config/logger"
	"vcs-sms/model/dto"
	"vcs-sms/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ReportController struct {
	service service.IReportService
}

func NewReportController(reportService service.IReportService) *ReportController {
	return &ReportController{
		service: reportService,
	}
}

func (controller *ReportController) SendReport(c *gin.Context) {
	log := logger.NewLogger()
	requestBody := &dto.ReportRequest{}
	if err := c.ShouldBindJSON(requestBody); err != nil {
		log.Error(fmt.Sprintf("Failed to bind json: %s", err.Error()), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format. Use dd-mm-yyyy"})
		return
	}
	startDate, err := time.Parse("02-01-2006", requestBody.StartDate)
	if err != nil {
		log.Error(fmt.Sprintf("Invalid start date format: %s", err.Error()), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format. Use dd-mm-yyyy"})
		return
	}
	endDate, err := time.Parse("02-01-2006", requestBody.EndDate)
	if err != nil {
		log.Error(fmt.Sprintf("Invalid end date format: %s", err.Error()), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format. Use dd-mm-yyyy"})
		return
	}
	startMils := startDate.UnixMilli()
	endMils := endDate.UnixMilli()

	err = controller.service.SendReport(startMils, endMils, []string{requestBody.Email})
	if err != nil {
		log.Error(fmt.Sprintf("Failed to send report: %s", err.Error()), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Info(fmt.Sprintf("Report sent to %s", requestBody.Email), zap.String("client", c.ClientIP()))
	c.JSON(200, gin.H{"message": "Report sent"})
	return
}

func (controller *ReportController) PeriodicReport(interval time.Duration) {
	controller.service.PeriodicReport(interval)
}
