package routes

import (
	"github.com/ashoksahoo/whatsapp-business-platform/internal/api/handlers"
	"github.com/ashoksahoo/whatsapp-business-platform/internal/api/middleware"
	"github.com/ashoksahoo/whatsapp-business-platform/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SetupRoutes registers all routes
func SetupRoutes(
	router *gin.Engine,
	messageHandler *handlers.MessageHandler,
	contactHandler *handlers.ContactHandler,
	templateHandler *handlers.TemplateHandler,
	webhookHandler *handlers.WebhookHandler,
	healthHandler *handlers.HealthHandler,
	authService *services.AuthService,
	logger *zap.Logger,
) {
	// Global middleware
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.LoggingMiddleware(logger))
	router.Use(middleware.RecoveryMiddleware(logger))
	router.Use(middleware.CORSMiddleware())

	// Public routes (no auth required)
	router.GET("/health", healthHandler.HealthCheck)

	// Webhook routes (no auth, signature verification inside handler)
	webhooks := router.Group("/webhooks")
	{
		webhooks.GET("/whatsapp", webhookHandler.VerifyWebhook)
		webhooks.POST("/whatsapp", webhookHandler.ReceiveWebhook)
	}

	// API v1 routes (auth required)
	v1 := router.Group("/api/v1")
	v1.Use(middleware.AuthMiddleware(authService))
	v1.Use(middleware.RateLimitMiddleware(1000)) // 1000 requests per minute
	{
		// Messages
		messages := v1.Group("/messages")
		{
			messages.POST("", messageHandler.SendMessage)
			messages.GET("", messageHandler.ListMessages)
			messages.GET("/search", messageHandler.SearchMessages)
			messages.GET("/:id", messageHandler.GetMessage)
		}

		// Contacts
		contacts := v1.Group("/contacts")
		{
			contacts.GET("", contactHandler.ListContacts)
			contacts.GET("/search", contactHandler.SearchContacts)
			contacts.GET("/:id", contactHandler.GetContact)
			contacts.PATCH("/:id", contactHandler.UpdateContact)
		}

		// Templates
		templates := v1.Group("/templates")
		{
			templates.GET("", templateHandler.ListTemplates)
			templates.POST("", templateHandler.CreateTemplate)
			templates.GET("/:id", templateHandler.GetTemplate)
			templates.PATCH("/:id", templateHandler.UpdateTemplate)
			templates.DELETE("/:id", templateHandler.DeleteTemplate)
		}
	}
}
