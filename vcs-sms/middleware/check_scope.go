package middleware

import (
	"vcs-sms/config/logger"
	"vcs-sms/model/entity"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CheckScope(scope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.NewLogger()
		scopes := c.MustGet("scopes")
		for _, value := range scopes.([]entity.Scope) {
			if value.Name == scope || value.Name == "all" {
				c.Next()
				return
			}
		}
		log.Error("Forbidden", zap.String("client", c.ClientIP()))
		c.JSON(403, gin.H{"error": "Forbidden"})
		c.Abort()
		return

	}
}
