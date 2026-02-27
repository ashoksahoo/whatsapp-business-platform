package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/ashoksahoo/whatsapp-business-platform/pkg/errors"
	"github.com/ashoksahoo/whatsapp-business-platform/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RecoveryMiddleware recovers from panics and logs them
func RecoveryMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Get stack trace
				stack := string(debug.Stack())

				// Log the panic
				logger.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("stack", stack),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)

				// Return error response
				appErr := errors.NewInternalError(fmt.Errorf("%v", err))
				utils.ErrorJSON(c, appErr)
				c.Abort()
			}
		}()

		c.Next()
	}
}
