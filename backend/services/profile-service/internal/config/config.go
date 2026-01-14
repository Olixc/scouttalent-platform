package config

import (
	"fmt"
	"os"

	"github.com/scouttalent/pkg/auth"
	"github.com/scouttalent/pkg/database"
	"github.com/scouttalent/pkg/messaging"
)

type Config struct {
	ServerAddress string
	Database      database.Config
	NATS          messaging.NATSConfig
	JWT           auth.JWTConfig
}

func Load() (*Config, error) {
	cfg := &Config{
		ServerAddress: getEnv("SERVER_ADDRESS", ":8081"),
		Database: database.Config{
			URL:             getEnv("DATABASE_URL", ""),
			MaxConns:        25,
			MinConns:        5,
			MaxConnLifetime: "1h",
			MaxConnIdleTime: "30m",
		},
		NATS: messaging.NATSConfig{
			URL: getEnv("NATS_URL", "nats://localhost:4222"),
		},
		JWT: auth.JWTConfig{
			Secret: getEnv("JWT_SECRET", ""),
		},
	}

	if cfg.Database.URL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	if cfg.JWT.Secret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}