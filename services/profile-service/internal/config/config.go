package config

import (
	"fmt"
	"os"
	"time"

	"github.com/scouttalent/pkg/auth"
	"github.com/scouttalent/pkg/database"
)

type Config struct {
	ServerAddress string
	Database      database.Config
	NATS          NATSConfig
	JWT           auth.TokenConfig
}

type NATSConfig struct {
	URL string
}

func Load() (*Config, error) {
	cfg := &Config{
		ServerAddress: getEnv("SERVER_ADDRESS", ":8081"),
		Database: database.Config{
			DSN:             getEnv("DATABASE_URL", ""),
			MaxOpenConns:    25,
			MaxIdleConns:    5,
			ConnMaxLifetime: 1 * time.Hour,
			ConnMaxIdleTime: 30 * time.Minute,
		},
		NATS: NATSConfig{
			URL: getEnv("NATS_URL", "nats://localhost:4222"),
		},
		JWT: auth.TokenConfig{
			SecretKey: getEnv("JWT_SECRET", ""),
			Issuer:    "scouttalent.com",
			Audience:  []string{"api.scouttalent.com"},
		},
	}

	if cfg.Database.DSN == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	if cfg.JWT.SecretKey == "" {
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