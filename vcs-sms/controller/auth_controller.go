package controller

import (
	"fmt"
	"net/http"
	"vcs-sms/config/logger"
	"vcs-sms/model/dto"
	"vcs-sms/model/entity"
	"vcs-sms/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	service    *service.UserService
	jwtService *service.JWTService
}

func NewAuthController(service *service.UserService, jwtService *service.JWTService) *AuthController {
	return &AuthController{service: service, jwtService: jwtService}
}

func (controller *AuthController) Register(c *gin.Context) {
	log := logger.NewLogger()
	authRequest := &dto.AuthRequest{}
	if err := c.ShouldBindJSON(&authRequest); err != nil {
		log.Error(fmt.Sprintf("Failed to bind JSON: %s", err.Error()), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	existUser := controller.service.FindByUsername(authRequest.Username)
	if existUser.Username != "" {
		log.Error(fmt.Sprintf("Username already exist: %s", authRequest.Username), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exist"})
		return
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(authRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to hash password: %s", err.Error()), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = controller.service.CreateNewUser(authRequest.Username, string(hashed), []entity.Scope{})
	if err != nil {
		log.Error(fmt.Sprintf("Failed to create user: %s", err.Error()), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
		return
	}
	log.Info(fmt.Sprintf("User %s created", authRequest.Username), zap.String("client", c.ClientIP()))
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
	return
}

func (controller *AuthController) Login(c *gin.Context) {
	log := logger.NewLogger()
	authRequest := &dto.AuthRequest{}
	if err := c.ShouldBindJSON(&authRequest); err != nil {
		log.Error(fmt.Sprintf("Failed to bind JSON: %s", err.Error()), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := controller.service.FindByUsername(authRequest.Username)
	if user.Username == "" {
		log.Error(fmt.Sprintf("Auth failed: %s", "Wrong username or password"), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong username or password"})
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(authRequest.Password))
	if err != nil {
		log.Error(fmt.Sprintf("Auth failed: %s", "Wrong username or password"), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong username or password"})
		return
	}
	accessToken, err := controller.jwtService.GenerateToken(&user)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to generate token: %s", err.Error()), zap.String("client", c.ClientIP()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	log.Info(fmt.Sprintf("User logged in"), zap.String("client", c.ClientIP()))
	c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
	return
}
