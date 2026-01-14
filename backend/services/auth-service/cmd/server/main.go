package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/scouttalent/auth-service/internal/config"
	"github.com/scouttalent/auth-service/internal/handler"
	"github.com/scouttalent/auth-service/internal/repository"
	"github.com/scouttalent/auth-service/internal/service"
	"github.com/scouttalent/pkg/database"
	"github.com/scouttalent/pkg/logging"
	"github.com/scouttalent/pkg/middleware"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Initialize logger
	logger, err := logging.NewLogger("auth-service")
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer logger.Sync()

	// Initialize database
	pool, err := database.NewPool(ctx, cfg.Database)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	logger.Info("connected to database")

	// Initialize layers
	repo := repository.NewUserRepository(pool)
	svc := service.NewAuthService(repo, cfg)
	h := handler.NewAuthHandler(svc, logger.Logger)

	// Setup router
	router := gin.Default()

	// Public routes
	router.GET("/health", h.Health)
	router.POST("/api/v1/auth/register", h.Register)
	router.POST("/api/v1/auth/login", h.Login)

	// Protected routes
	protected := router.Group("/api/v1/auth")
	protected.Use(middleware.AuthMiddleware(cfg.JWT))
	{
		protected.GET("/me", h.GetMe)
	}

	// Start server
	logger.Info("starting server", zap.String("address", cfg.ServerAddress))

	go func() {
		if err := router.Run(cfg.ServerAddress); err != nil {
			logger.Fatal("server failed", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")
}