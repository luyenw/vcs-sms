package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"vcs-sms/config/logger"
	"vcs-sms/model/dto"
	"vcs-sms/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	service *service.UserService
}

func NewUserController(service *service.UserService) *UserController {
	return &UserController{service: service}
}

func (controller *UserController) CreateUser(c *gin.Context) {
	log := logger.NewLogger()
	createUserRequest := &dto.CreateUserRequest{}
	if err := c.ShouldBindJSON(createUserRequest); err != nil {
		log.Error(fmt.Sprintf("Failed to bind json: %s", err.Error()), zap.String("client", c.ClientIP()))
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	existUser := controller.service.FindByUsername(createUserRequest.Username)
	if existUser.Username != "" {
		log.Error(fmt.Sprintf("Username already exist: %s", createUserRequest.Username), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exist"})
		return
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(createUserRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to hash password: %s", err.Error()), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = controller.service.CreateNewUser(createUserRequest.Username, string(hashed), createUserRequest.Scopes)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to create user: %s", err.Error()), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
		return
	}
	log.Info(fmt.Sprintf("User created successfully: %s", createUserRequest.Username), zap.String("client", c.ClientIP()))
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func (controller *UserController) UpdateUserScope(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	log := logger.NewLogger()
	if err != nil {
		log.Error(fmt.Sprintf("Failed to parse id: %s", err.Error()), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	user := controller.service.FindUserByID(id)
	if user == nil {
		log.Error(fmt.Sprintf("User not found: %d", id), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	scopes := dto.UpdatePermissionRequest{}
	if err := c.ShouldBindJSON(&scopes); err != nil {
		log.Error(fmt.Sprintf("Failed to bind json: %s", err.Error()), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = controller.service.UpdateUserScope(user, scopes.Scopes)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to update user scope: %s", err.Error()), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user scope: " + err.Error()})
		return
	}
	log.Info(fmt.Sprintf("User scope updated successfully: %d", id), zap.String("client", c.ClientIP()))
	c.JSON(http.StatusOK, gin.H{"message": "User scope updated successfully"})
	return
}
