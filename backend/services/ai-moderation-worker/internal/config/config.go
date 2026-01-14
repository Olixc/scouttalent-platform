package config

import (
	"os"

	"github.com/scouttalent/pkg/database"
	"github.com/scouttalent/pkg/messaging"
)

type Config struct {
	Database database.Config
	NATS     messaging.NATSConfig
	OpenAI   OpenAIConfig
}

type OpenAIConfig struct {
	APIKey string
	Model  string
}

func Load() (*Config, error) {
	return &Config{
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
		OpenAI: OpenAIConfig{
			APIKey: getEnv("OPENAI_API_KEY", ""),
			Model:  getEnv("OPENAI_MODEL", "gpt-4o-mini"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}