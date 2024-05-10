package middleware

import (
	"log"
	"net/http"
	"time"
	"vcs-sms/config/sql"
	"vcs-sms/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func TokenAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")
		if authorizationHeader == "" {
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
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		if int64(claims["exp"].(float64)) < time.Now().Unix() {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		username := claims["username"]
		userService := service.NewUserService(sql.GetPostgres())
		user := userService.FindByUsername(username.(string))
		log.Println(user.Scopes)
		if user.Username == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
		}
		c.Set("scopes", user.Scopes)
		c.Next()
	}
}
