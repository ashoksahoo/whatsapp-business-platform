package handlers

import (
	"time"

	"github.com/ashoksahoo/whatsapp-business-platform/internal/database"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	db        *gorm.DB
	startTime time.Time
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{
		db:        db,
		startTime: time.Now(),
	}
}

// HealthCheck handles GET /health
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	// Check database health
	dbHealth := database.HealthCheck(h.db)

	uptime := time.Since(h.startTime)

	response := gin.H{
		"status":  "healthy",
		"version": "0.1.0",
		"uptime":  uptime.String(),
		"checks": gin.H{
			"database": map[string]interface{}{
				"status":        "connected",
				"healthy":       dbHealth.Healthy,
				"response_time": dbHealth.ResponseTime.String(),
				"connections":   dbHealth.OpenConns,
			},
		},
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	// If database is unhealthy, return 503
	if !dbHealth.Healthy {
		response["status"] = "unhealthy"
		response["checks"].(gin.H)["database"].(map[string]interface{})["status"] = "error"
		response["checks"].(gin.H)["database"].(map[string]interface{})["error"] = dbHealth.Error
		c.JSON(503, response)
		return
	}

	c.JSON(200, response)
}
