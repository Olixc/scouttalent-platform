package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/scouttalent/profile-service/internal/config"
	"github.com/scouttalent/profile-service/internal/handler"
	"github.com/scouttalent/profile-service/internal/repository"
	"github.com/scouttalent/profile-service/internal/service"
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
	logger, err := logging.NewLogger("profile-service")
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
	repo := repository.NewProfileRepository(pool)
	svc := service.NewProfileService(repo)
	h := handler.NewProfileHandler(svc, logger.Logger)

	// Setup router
	router := gin.Default()

	// Public routes
	router.GET("/health", h.Health)

	// Protected routes
	api := router.Group("/api/v1/profiles")
	api.Use(middleware.AuthMiddleware(cfg.JWT))
	{
		api.POST("", h.CreateProfile)
		api.GET("/me", h.GetMyProfile)
		api.GET("/:id", h.GetProfile)
		api.PUT("/:id", h.UpdateProfile)
		
		// Player-specific routes
		api.POST("/:id/player-details", h.CreatePlayerDetails)
		api.GET("/:id/player", h.GetPlayerProfile)
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