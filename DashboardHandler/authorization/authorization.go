package authorization

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckAuthorization(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		expected := fmt.Sprintf("Bearer %s", apiKey)

		if authHeader != expected {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}
