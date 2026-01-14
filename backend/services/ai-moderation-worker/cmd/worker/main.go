package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/scouttalent/ai-moderation-worker/internal/config"
	"github.com/scouttalent/ai-moderation-worker/internal/moderator"
	"github.com/scouttalent/ai-moderation-worker/internal/worker"
	"github.com/scouttalent/pkg/database"
	"github.com/scouttalent/pkg/logging"
	"github.com/scouttalent/pkg/messaging"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logger, err := logging.NewLogger("ai-moderation-worker")
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

	// Connect to NATS
	nc, err := messaging.NewNATSClient(cfg.NATS)
	if err != nil {
		logger.Fatal("Failed to connect to NATS", err)
	}
	defer nc.Close()

	logger.Info("Connected to NATS")

	// Initialize AI moderator
	mod := moderator.NewAIModerator(cfg.OpenAI, logger.Logger)

	// Initialize worker
	w := worker.NewWorker(db, nc, mod, logger.Logger)

	// Start worker
	if err := w.Start(ctx); err != nil {
		logger.Fatal("Failed to start worker", err)
	}

	logger.Info("AI Moderation Worker started successfully")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutting down AI Moderation Worker...")

	// Graceful shutdown
	if err := w.Stop(); err != nil {
		logger.Error("Error during shutdown", err)
	}

	logger.Info("AI Moderation Worker stopped")
}