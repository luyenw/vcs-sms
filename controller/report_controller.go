package controller

import (
	"net/http"
	"time"
	"vcs-sms/model/dto"
	"vcs-sms/service"

	"github.com/gin-gonic/gin"
)

type ReportController struct {
	service *service.ReportService
}

func NewReportController(reportService *service.ReportService) *ReportController {
	return &ReportController{
		service: reportService,
	}
}

func (controller *ReportController) SendReport(c *gin.Context) {
	requestBody := &dto.ReportRequest{}
	if err := c.ShouldBindJSON(requestBody); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	startDate, err := time.Parse("02-01-2006", requestBody.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format. Use dd-mm-yyyy"})
		return
	}
	endDate, err := time.Parse("02-01-2006", requestBody.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format. Use dd-mm-yyyy"})
		return
	}
	startMils := startDate.UnixMilli()
	endMils := endDate.UnixMilli()

	err = controller.service.SendReport(startMils, endMils, []string{requestBody.Email})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Report sent"})
	return
}

func (controller *ReportController) PeriodicReport(interval time.Duration) {
	controller.service.PeriodicReport(interval)
}
