package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ashoksahoo/whatsapp-business-platform/internal/api/handlers"
	"github.com/ashoksahoo/whatsapp-business-platform/internal/api/routes"
	"github.com/ashoksahoo/whatsapp-business-platform/internal/config"
	"github.com/ashoksahoo/whatsapp-business-platform/internal/repositories"
	"github.com/ashoksahoo/whatsapp-business-platform/internal/services"
	"github.com/ashoksahoo/whatsapp-business-platform/internal/whatsapp"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Server represents the API server
type Server struct {
	router     *gin.Engine
	httpServer *http.Server
	config     *config.Config
	logger     *zap.Logger
}

// NewServer creates a new API server
func NewServer(cfg *config.Config, db *gorm.DB, logger *zap.Logger) (*Server, error) {
	// Set Gin mode
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	router := gin.New()

	// Initialize WhatsApp client
	waClient, err := whatsapp.NewClient(whatsapp.Config{
		APIToken:      cfg.WhatsApp.APIToken,
		PhoneNumberID: cfg.WhatsApp.PhoneNumberID,
		APIBaseURL:    cfg.WhatsApp.APIBaseURL,
		APIVersion:    cfg.WhatsApp.APIVersion,
		Logger:        logger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create WhatsApp client: %w", err)
	}

	// Initialize repositories
	messageRepo := repositories.NewMessageRepository(db)
	contactRepo := repositories.NewContactRepository(db)
	templateRepo := repositories.NewTemplateRepository(db)
	apiKeyRepo := repositories.NewAPIKeyRepository(db)

	// Initialize services
	messageService := services.NewMessageService(messageRepo, contactRepo, waClient, logger)
	contactService := services.NewContactService(contactRepo)
	templateService := services.NewTemplateService(templateRepo)
	authService := services.NewAuthService(apiKeyRepo)

	// Initialize handlers
	messageHandler := handlers.NewMessageHandler(messageService)
	contactHandler := handlers.NewContactHandler(contactService)
	templateHandler := handlers.NewTemplateHandler(templateService)
	webhookHandler := handlers.NewWebhookHandler(
		messageService,
		cfg.WhatsApp.WebhookVerifyToken,
		cfg.WhatsApp.WebhookSecret,
		logger,
	)
	healthHandler := handlers.NewHealthHandler(db)

	// Setup routes
	routes.SetupRoutes(
		router,
		messageHandler,
		contactHandler,
		templateHandler,
		webhookHandler,
		healthHandler,
		authService,
		logger,
	)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	return &Server{
		router:     router,
		httpServer: httpServer,
		config:     cfg,
		logger:     logger,
	}, nil
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.logger.Info("Starting HTTP server",
		zap.String("address", s.httpServer.Addr),
		zap.String("environment", s.config.Server.Environment),
	)

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down HTTP server...")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	s.logger.Info("HTTP server stopped")
	return nil
}

// GetRouter returns the Gin router (useful for testing)
func (s *Server) GetRouter() *gin.Engine {
	return s.router
}
