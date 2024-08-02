package middleware

import (
	"fmt"
	"net/http"
	"time"
	"vcs-sms/config/logger"
	"vcs-sms/config/sql"
	"vcs-sms/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

func TokenAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.NewLogger()
		authorizationHeader := c.GetHeader("Authorization")
		if authorizationHeader == "" {
			log.Error(fmt.Sprintf("Unauthorized: %s", "No token provided"), zap.String("client", c.ClientIP()))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(authorizationHeader, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		claims = token.Claims.(jwt.MapClaims)
		if err != nil {
			log.Error(fmt.Sprintf("Unauthorized: %s", err.Error()), zap.String("client", c.ClientIP()))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		if int64(claims["exp"].(float64)) < time.Now().Unix() {
			log.Error(fmt.Sprintf("Unauthorized: %s", "Token expired"), zap.String("client", c.ClientIP()))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		username := claims["username"]
		userService := service.NewUserService(sql.GetPostgres())
		user := userService.FindByUsername(username.(string))
		if user.Username == "" {
			log.Error(fmt.Sprintf("Unauthorized: %s", "User not found"), zap.String("client", c.ClientIP()))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
		}
		c.Set("scopes", user.Role.Scopes)
		c.Next()
	}
}
