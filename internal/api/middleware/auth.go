package middleware

import (
	"strings"

	"github.com/ashoksahoo/whatsapp-business-platform/internal/services"
	"github.com/ashoksahoo/whatsapp-business-platform/pkg/errors"
	"github.com/ashoksahoo/whatsapp-business-platform/pkg/utils"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates API keys
func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorJSON(c, errors.NewUnauthorized("Authorization header required"))
			c.Abort()
			return
		}

		// Check Bearer format
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorJSON(c, errors.NewUnauthorized("Invalid authorization format. Use: Bearer <api_key>"))
			c.Abort()
			return
		}

		apiKey := parts[1]
		if apiKey == "" {
			utils.ErrorJSON(c, errors.NewUnauthorized("API key is empty"))
			c.Abort()
			return
		}

		// Validate API key
		keyInfo, err := authService.ValidateAPIKey(apiKey)
		if err != nil {
			utils.ErrorJSON(c, err.(*errors.AppError))
			c.Abort()
			return
		}

		// Store API key info in context
		c.Set("api_key", keyInfo)
		c.Set("api_key_id", keyInfo.ID)

		c.Next()
	}
}
