package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/scouttalent/discovery-service/internal/config"
	"github.com/scouttalent/discovery-service/internal/handler"
	"github.com/scouttalent/discovery-service/internal/repository"
	"github.com/scouttalent/discovery-service/internal/service"
	"github.com/scouttalent/pkg/database"
	"github.com/scouttalent/pkg/logging"
	"github.com/scouttalent/pkg/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logger, err := logging.NewLogger("discovery-service")
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	ctx := context.Background()

	// Connect to database
	db, err := database.NewPool(ctx, cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", err)
	}
	defer db.Close()

	logger.Info("Connected to database")

	// Initialize repositories
	profileRepo := repository.NewProfileRepository(db)
	videoRepo := repository.NewVideoRepository(db)

	// Initialize services
	searchService := service.NewSearchService(profileRepo, videoRepo, logger.Logger)
	recommendationService := service.NewRecommendationService(profileRepo, videoRepo, logger.Logger)
	feedService := service.NewFeedService(videoRepo, logger.Logger)

	// Initialize handlers
	searchHandler := handler.NewSearchHandler(searchService)
	recommendationHandler := handler.NewRecommendationHandler(recommendationService)
	feedHandler := handler.NewFeedHandler(feedService)

	// Setup router
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "discovery-service",
		})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Search endpoints
		search := api.Group("/search")
		{
			search.GET("/profiles", searchHandler.SearchProfiles)
			search.GET("/videos", searchHandler.SearchVideos)
		}

		// Recommendation endpoints (require authentication)
		recommendations := api.Group("/recommendations")
		recommendations.Use(middleware.AuthMiddleware(cfg.JWT))
		{
			recommendations.GET("/profiles", recommendationHandler.GetProfileRecommendations)
			recommendations.GET("/videos", recommendationHandler.GetVideoRecommendations)
		}

		// Feed endpoints (require authentication)
		feed := api.Group("/feed")
		feed.Use(middleware.AuthMiddleware(cfg.JWT))
		{
			feed.GET("", feedHandler.GetFeed)
			feed.GET("/trending", feedHandler.GetTrendingVideos)
		}
	}

	// Start server
	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: router,
	}

	go func() {
		logger.Info(fmt.Sprintf("Discovery Service starting on %s", cfg.ServerAddress))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", err)
	}

	logger.Info("Server exited")
}