package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/scouttalent/media-service/internal/config"
	"github.com/scouttalent/media-service/internal/handler"
	"github.com/scouttalent/media-service/internal/repository"
	"github.com/scouttalent/media-service/internal/service"
	"github.com/scouttalent/media-service/internal/storage"
	"github.com/scouttalent/pkg/database"
	"github.com/scouttalent/pkg/logging"
	"github.com/scouttalent/pkg/messaging"
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
	logger, err := logging.NewLogger("media-service")
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

	// Initialize NATS
	nc, err := messaging.NewNATSClient(cfg.NATS)
	if err != nil {
		logger.Fatal("failed to connect to NATS", zap.Error(err))
	}
	defer nc.Close()

	logger.Info("connected to NATS")

	// Initialize Azure Blob Storage
	blobClient, err := storage.NewAzureBlobClient(cfg.Azure)
	if err != nil {
		logger.Fatal("failed to initialize Azure Blob Storage", zap.Error(err))
	}

	logger.Info("initialized Azure Blob Storage")

	// Initialize layers
	repo := repository.NewMediaRepository(pool)
	svc := service.NewMediaService(repo, blobClient, nc, logger.Logger)
	h := handler.NewMediaHandler(svc, logger.Logger)

	// Setup router
	router := gin.Default()

	// Public routes
	router.GET("/health", h.Health)

	// Protected routes
	api := router.Group("/api/v1/videos")
	api.Use(middleware.AuthMiddleware(cfg.JWT))
	{
		api.POST("/upload", h.InitiateUpload)
		api.PATCH("/upload/:id", h.ResumeUpload)
		api.POST("/:id/complete", h.CompleteUpload)
		api.GET("/:id", h.GetVideo)
		api.GET("/profile/:profile_id", h.ListProfileVideos)
		api.PUT("/:id", h.UpdateVideo)
		api.DELETE("/:id", h.DeleteVideo)
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