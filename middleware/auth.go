package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// APIKeyAuth middleware validates the x-api-key header
func APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("x-api-key")
		expectedKey := os.Getenv("API_KEY")

		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "API key is required",
				"data":    nil,
			})
			c.Abort()
			return
		}

		if apiKey != expectedKey {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid API key",
				"data":    nil,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
