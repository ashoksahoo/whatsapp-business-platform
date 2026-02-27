package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ashoksahoo/whatsapp-business-platform/internal/api"
	"github.com/ashoksahoo/whatsapp-business-platform/internal/config"
	"github.com/ashoksahoo/whatsapp-business-platform/internal/database"
	"github.com/ashoksahoo/whatsapp-business-platform/pkg/logger"
	"go.uber.org/zap"
	gormlogger "gorm.io/gorm/logger"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log, err := logger.InitLogger(logger.Config{
		Level:      cfg.Logging.Level,
		Format:     cfg.Logging.Format,
		OutputPath: cfg.Logging.OutputPath,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	log.Info("Starting Vibecoded WA Client",
		zap.String("environment", cfg.Server.Environment),
		zap.Int("port", cfg.Server.Port),
	)

	// Initialize database connection
	dbLogLevel := gormlogger.Silent
	if cfg.IsDevelopment() {
		dbLogLevel = gormlogger.Info
	}

	db, err := database.NewConnection(cfg.GetDatabaseDriver(), cfg.GetDatabaseDSN(), dbLogLevel)
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}
	log.Info("Database connection established",
		zap.String("driver", cfg.GetDatabaseDriver()),
	)

	// Run auto migrations
	if err := database.AutoMigrate(db); err != nil {
		log.Fatal("Failed to run database migrations", zap.Error(err))
	}
	log.Info("Database migrations completed")

	// Create indexes and triggers
	if err := database.CreateIndexes(db); err != nil {
		log.Warn("Failed to create some indexes", zap.Error(err))
	}
	if err := database.CreateTriggers(db); err != nil {
		log.Warn("Failed to create triggers", zap.Error(err))
	}

	// Health check
	health := database.HealthCheck(db)
	if !health.Healthy {
		log.Fatal("Database health check failed", zap.String("error", health.Error))
	}
	log.Info("Database health check passed",
		zap.Duration("response_time", health.ResponseTime),
		zap.Int("open_connections", health.OpenConns),
	)

	// Initialize API server
	server, err := api.NewServer(cfg, db, log)
	if err != nil {
		log.Fatal("Failed to create API server", zap.Error(err))
	}

	// Start server in a goroutine
	serverErrors := make(chan error, 1)
	go func() {
		log.Info("API server starting",
			zap.String("host", cfg.Server.Host),
			zap.Int("port", cfg.Server.Port),
		)
		serverErrors <- server.Start()
	}()

	// Wait for interrupt signal or server error
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatal("Server error", zap.Error(err))
	case sig := <-quit:
		log.Info("Shutdown signal received", zap.String("signal", sig.String()))
	}

	log.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	// Shutdown API server
	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server shutdown error", zap.Error(err))
	}

	// Close database connection
	if err := database.CloseConnection(db); err != nil {
		log.Error("Error closing database connection", zap.Error(err))
	}

	// Sync logger
	logger.Sync()

	log.Info("Server shutdown complete")

	select {
	case <-ctx.Done():
		log.Warn("Shutdown timeout exceeded")
		os.Exit(1)
	default:
		log.Info("Shutdown gracefully completed")
	}
}
