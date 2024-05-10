package controller

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"vcs-sms/model/dto"
	"vcs-sms/model/entity"
	"vcs-sms/service"

	"github.com/gin-gonic/gin"
)

type ServerController struct {
	service *service.ServerService
}

func NewServerController(service *service.ServerService) *ServerController {
	return &ServerController{service: service}
}

func (controller *ServerController) GetServer(c *gin.Context) {
	queryParam := &dto.QueryParam{}
	if err := c.ShouldBindQuery(queryParam); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	servers := controller.service.GetServer(queryParam)
	c.JSON(200, servers)
	return
}

func (controller *ServerController) UpdateServer(c *gin.Context) {
	id_param := c.Param("id")
	id, err := strconv.Atoi(id_param)
	if err != nil {
		c.JSON(400, gin.H{"error": "Record not found"})
		return
	}
	server := controller.service.FindServerById(id)
	if server == nil {
		c.JSON(404, gin.H{"error": "Record not found"})
		return
	}
	inputServer := &dto.InputServer{}
	if err := c.ShouldBindJSON(inputServer); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if inputServer.Name != "" {
		server.Name = inputServer.Name
	}
	if inputServer.IPv4 != "" {
		server.IPv4 = inputServer.IPv4
	}
	if inputServer.Status != 0 {
		server.Status = inputServer.Status
	}
	server.LastUpdated = time.Now()
	if err := controller.service.UpdateServer(server); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, server)
	return
}
func (controller *ServerController) DeleteServer(c *gin.Context) {
	id_param := c.Param("id")
	id, err := strconv.Atoi(id_param)
	if err != nil {
		c.JSON(400, gin.H{"error": "Record not found"})
		return
	}
	err = controller.service.DeleteServerById(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Record deleted successfully"})
	return
}
func (controller *ServerController) CreateServer(c *gin.Context) {
	inputServer := &dto.InputServer{}
	if err := c.ShouldBindJSON(inputServer); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	now := time.Now()
	server := &entity.Server{
		Name:        inputServer.Name,
		IPv4:        inputServer.IPv4,
		Status:      inputServer.Status,
		CreatedTime: now,
		LastUpdated: now,
	}
	err := controller.service.CreateServer(server)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, server)
	return
}

func (controller *ServerController) ExportServers(c *gin.Context) {
	queryParam := &dto.QueryParam{}
	if err := c.ShouldBindQuery(queryParam); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	servers := controller.service.GetServer(queryParam)
	service := service.NewXLSXService()
	exportURL, err := service.ExportXLSX(servers)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Exported successfully", "url": exportURL})
}
func (controller *ServerController) ImportServers(c *gin.Context) {
	file, header, err := c.Request.FormFile("input")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	currentTimeF := time.Now().Format("06-01-02_15-04-05")
	tmpFilePath := "./tmp/" + currentTimeF + "_" + header.Filename
	out, err := os.Create(tmpFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := service.NewXLSXService().ImportXLSX(tmpFilePath)

	successCount := 0
	failureCount := 0
	var successNames []string
	for idx, row := range rows {
		currentTime := time.Now()
		if idx == 0 {
			continue
		}
		name, ipv4 := row[0], row[1]
		server := &entity.Server{
			Name:        name,
			IPv4:        ipv4,
			Status:      0,
			CreatedTime: currentTime,
		}
		err := controller.service.CreateServer(server)
		if err != nil {
			failureCount++
		} else {
			successCount++
			successNames = append(successNames, server.Name)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"failure_count": failureCount,
		"success_count": successCount,
		"success_names": successNames,
	})
}
