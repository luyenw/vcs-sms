package controller

import (
	"net/http"
	"strconv"
	"vcs-sms/model/dto"
	"vcs-sms/service"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	service *service.UserService
}

func NewUserController(service *service.UserService) *UserController {
	return &UserController{service: service}
}

func (controller *UserController) CreateUser(c *gin.Context) {
	createUserRequest := &dto.CreateUserRequest{}
	if err := c.ShouldBindJSON(createUserRequest); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	existUser := controller.service.FindByUsername(createUserRequest.Username)
	if existUser.Username != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exist"})
		return
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(createUserRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = controller.service.CreateNewUser(createUserRequest.Username, string(hashed), createUserRequest.Scopes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func (controller *UserController) UpdateUserScope(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	user := controller.service.FindUserByID(id)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	scopes := dto.UpdatePermissionRequest{}
	if err := c.ShouldBindJSON(&scopes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = controller.service.UpdateUserScope(user, scopes.Scopes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user scope: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User scope updated successfully"})
	return
}
