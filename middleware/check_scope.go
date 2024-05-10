package middleware

import (
	"vcs-sms/model/entity"

	"github.com/gin-gonic/gin"
)

func CheckScope(scope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		scopes := c.MustGet("scopes")
		for _, value := range scopes.([]entity.Scope) {
			if value.Name == scope || value.Name == "root" {
				c.Next()
				return
			}
		}
		c.JSON(403, gin.H{"error": "Forbidden"})
		c.Abort()
		return

	}
}
