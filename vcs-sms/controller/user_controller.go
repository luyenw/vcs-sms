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
	service service.IUserService
}

func NewUserController(service service.IUserService) *UserController {
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
	err = controller.service.CreateNewUser(createUserRequest.Username, string(hashed), createUserRequest.Role)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to create user: %s", err.Error()), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
		return
	}
	log.Info(fmt.Sprintf("User created successfully: %s", createUserRequest.Username), zap.String("client", c.ClientIP()))
	createdUser := controller.service.FindByUsername(createUserRequest.Username)
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "data": dto.UserEntity{&createdUser}.ToDTO()})
}
func (controller *UserController) UpdateUserRole(c *gin.Context) {
	log := logger.NewLogger()
	idParam := c.Param("id")
	parseId, err := strconv.ParseInt(idParam, 10, 64)
	id := int(parseId)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to parse id: %s", err.Error()), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse id"})
		return
	}
	user := controller.service.FindUserByID(id)
	if user == nil {
		log.Error(fmt.Sprintf("User not found: %d", id), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	roleDto := dto.UpdateRoleRequest{}
	if err := c.ShouldBindJSON(&roleDto); err != nil {
		log.Error(fmt.Sprintf("Failed to bind json: %s", err.Error()), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	role, err := controller.service.FindRoleByName(roleDto.Role.Name)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to find role: %s", err.Error()), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role not found"})
		return
	}
	err = controller.service.UpdateUserRole(user, role)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to update user role: %s", err.Error()), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user role: " + err.Error()})
		return
	}
	log.Info(fmt.Sprintf("User role updated successfully: %d", id), zap.String("client", c.ClientIP()))
	c.JSON(http.StatusOK, gin.H{
		"data":    dto.UserEntity{controller.service.FindUserByID(id)}.ToDTO(),
		"message": "User role updated successfully",
	})
}
