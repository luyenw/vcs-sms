package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"vcs-sms/config/logger"
	"vcs-sms/model/dto"
	"vcs-sms/model/entity"
	"vcs-sms/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type ServerController struct {
	service      service.IServerService
	cacheService service.ICacheService
	xlsxService  service.IXLSXService
}

func NewServerController(serverService service.IServerService, cacheService service.ICacheService, xlsxService service.IXLSXService) *ServerController {
	return &ServerController{service: serverService, cacheService: cacheService, xlsxService: xlsxService}
}

func (controller *ServerController) GetServer(c *gin.Context) {
	log := logger.NewLogger()
	queryParam := &dto.QueryParam{}
	if err := c.ShouldBindQuery(queryParam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	validate := validator.New()
	if err := validate.Struct(queryParam); err != nil {
		log.Error(fmt.Sprintf("Failed to validate query param: %s", err.Error()), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	servers := controller.service.GetServer(queryParam)
	c.JSON(http.StatusOK, gin.H{"total": len(servers), "data": servers})
	return
}

func (controller *ServerController) UpdateServer(c *gin.Context) {
	id_param := c.Param("id")
	parseId, err := strconv.ParseInt(id_param, 10, 64)
	id := int(parseId)
	log := logger.NewLogger()
	if err != nil {
		log.Error(fmt.Sprintf("Failed to convert id: %s", err.Error()), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found"})
		return
	}
	server := controller.service.FindServerById(id)
	if server == nil {
		log.Error(fmt.Sprintf("Failed to find server: %d", id), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}
	inputServer := &dto.InputServer{}
	if err := c.ShouldBindJSON(inputServer); err != nil {
		log.Error(fmt.Sprintf("Failed to bind input server: %s", err.Error()), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		log.Error(fmt.Sprintf("Failed to update server: %s", err.Error()), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := controller.cacheService.Set("server:"+strconv.Itoa(int(server.ID)), server); err != nil {
		log.Error(fmt.Sprintf("Failed to set cache: %s", err.Error()), zap.String("client", c.ClientIP()))
	} else {
		log.Info(fmt.Sprintf("Set cache successfully: %s", "server:"+strconv.Itoa(int(server.ID))), zap.String("client", c.ClientIP()))
	}
	controller.cacheService.Set("server:all", controller.service.GetAllServers())
	log.Info(fmt.Sprintf("Updated server: %+v", server), zap.String("client", c.ClientIP()))
	c.JSON(http.StatusOK, gin.H{"message": "Record updated successfully", "data": server})
	return
}
func (controller *ServerController) DeleteServer(c *gin.Context) {
	log := logger.NewLogger()
	id_param := c.Param("id")
	parseId, err := strconv.ParseInt(id_param, 10, 64)
	id := int(parseId)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to convert id: %s", err.Error()), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found"})
		return
	}
	err = controller.service.DeleteServerById(id)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to delete server: %s", err.Error()), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := controller.cacheService.Set("server:"+strconv.Itoa(id), nil); err != nil {
		log.Error(fmt.Sprintf("Failed to set cache: %s", err.Error()), zap.String("client", c.ClientIP()))
	} else {
		log.Info(fmt.Sprintf("Set cache successfully: %s", "server:"+strconv.Itoa(id)), zap.String("client", c.ClientIP()))
	}
	controller.cacheService.Set("server:all", nil)
	log.Info(fmt.Sprintf("Deleted server: %d", id), zap.String("client", c.ClientIP()))
	c.JSON(http.StatusOK, gin.H{"message": "Record deleted successfully"})
	return
}
func (controller *ServerController) CreateServer(c *gin.Context) {
	log := logger.NewLogger()
	validate := validator.New()
	inputServer := &dto.InputServer{}
	//
	if err := c.ShouldBindJSON(inputServer); err != nil {
		log.Error(fmt.Sprintf("Failed to bind input server: %s", err.Error()), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validate.Struct(inputServer); err != nil {
		log.Error(fmt.Sprintf("Failed to validate input server: %s", err.Error()), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
	if err := controller.service.CreateServer(server); err != nil {
		log.Error(fmt.Sprintf("Failed to create server: %s", err.Error()), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	allServersString, _ := json.Marshal(controller.service.GetAllServers())
	if err := controller.cacheService.Set("server:all", allServersString); err != nil {
		log.Error("Failed to set all servers")
	} else {
		log.Info(fmt.Sprintf("Set all servers %s", allServersString))
	}
	log.Info(fmt.Sprintf("Created server: %+v", server), zap.String("client", c.ClientIP()))
	c.JSON(http.StatusOK, gin.H{"message": "Record created successfully", "data": server})
	return
}

func (controller *ServerController) ExportServers(c *gin.Context) {
	log := logger.NewLogger()
	queryParam := &dto.QueryParam{}
	if err := c.ShouldBindQuery(queryParam); err != nil {
		log.Error(fmt.Sprintf("Failed to bind query param: %s", err.Error()), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	validate := validator.New()
	if err := validate.Struct(queryParam); err != nil {
		log.Error(fmt.Sprintf("Failed to validate query param: %s", err.Error()), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	servers := controller.service.GetServer(queryParam)
	exportURL, err := controller.xlsxService.ExportXLSX(servers)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to export: %s", err.Error()), zap.String("client", c.ClientIP()))
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Info(fmt.Sprintf("Exported successfully: %s", exportURL), zap.String("client", c.ClientIP()))
	c.JSON(http.StatusOK, gin.H{"message": "Exported successfully", "url": exportURL})
}
func (controller *ServerController) ImportServers(c *gin.Context) {
	log := logger.NewLogger()
	file, header, err := c.Request.FormFile("input")
	if err != nil {
		log.Error(fmt.Sprintf("Failed to get file from request: %s", err.Error()), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot get file from request"})
		return
	}
	parts := strings.Split(header.Filename, ".")
	if len(parts) < 2 || parts[len(parts)-1] != "xlsx" {
		log.Error(fmt.Sprintf("Invalid file format: %s", header.Filename), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file format"})
		return
	}
	currentTimeF := time.Now().Format("06-01-02_15-04-05")
	tmpFilePath := "./tmp/" + currentTimeF + "_" + header.Filename
	out, err := os.Create(tmpFilePath)
	if err != nil {
		fmt.Println(err.Error())
		log.Error(fmt.Sprintf("Failed to create file: %s", err.Error()), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Error(fmt.Sprintf(": %s", err.Error()), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to copy file"})
		return
	}

	rows, err := controller.xlsxService.ImportXLSX(tmpFilePath)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to import file: %s", err.Error()), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to import file"})
		return
	}
	successCount := 0
	failureCount := 0
	successNames := []string{}
	validate := validator.New()
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
		if err := validate.Struct(server); err != nil {
			failureCount++
		} else {
			if err1 := controller.service.CreateServer(server); err1 != nil {
				failureCount++
			} else {
				successCount++
				successNames = append(successNames, server.Name)
			}
		}

	}
	allServersString, _ := json.Marshal(controller.service.GetAllServers())
	if err := controller.cacheService.Set("server:all", allServersString); err != nil {
		log.Error(fmt.Sprintf("Failed to set cache: %s", err.Error()), zap.String("client", c.ClientIP()))
	} else {
		log.Info(fmt.Sprintf("Set cache successfully: %s", "server:all"), zap.String("client", c.ClientIP()))
	}

	log.Info(fmt.Sprintf("Imported successfully: %d success, %d failure", successCount, failureCount), zap.String("client", c.ClientIP()))
	c.JSON(http.StatusOK, gin.H{
		"message": "Imported successfully",
		"data": gin.H{
			"failure_count": failureCount,
			"success_count": successCount,
			"success_names": successNames,
		},
	})
}
