package config

import (
	"os"

	"github.com/scouttalent/pkg/auth"
	"github.com/scouttalent/pkg/azure"
	"github.com/scouttalent/pkg/database"
	"github.com/scouttalent/pkg/messaging"
)

type Config struct {
	ServerAddress string
	Database      database.Config
	NATS          messaging.NATSConfig
	JWT           auth.JWTConfig
	Azure         azure.BlobConfig
}

func Load() (*Config, error) {
	return &Config{
		ServerAddress: getEnv("SERVER_ADDRESS", ":8082"),
		Database: database.Config{
			URL:             getEnv("DATABASE_URL", "postgres://scout:scoutpass@localhost:5432/media_db?sslmode=disable"),
			MaxConns:        10,
			MinConns:        2,
			MaxConnLifetime: "1h",
			MaxConnIdleTime: "30m",
		},
		NATS: messaging.NATSConfig{
			URL: getEnv("NATS_URL", "nats://localhost:4222"),
		},
		JWT: auth.JWTConfig{
			Secret: getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		},
		Azure: azure.BlobConfig{
			AccountName:   getEnv("AZURE_STORAGE_ACCOUNT", ""),
			AccountKey:    getEnv("AZURE_STORAGE_KEY", ""),
			ContainerName: getEnv("AZURE_CONTAINER_NAME", "videos"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}