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
	Redis         RedisConfig
	JWT           auth.TokenConfig
}

type RedisConfig struct {
	URL string
}

func Load() (*Config, error) {
	cfg := &Config{
		ServerAddress: getEnv("SERVER_ADDRESS", ":8080"),
		Database: database.Config{
			DSN:             getEnv("DATABASE_URL", ""),
			MaxOpenConns:    25,
			MaxIdleConns:    5,
			ConnMaxLifetime: 1 * time.Hour,
			ConnMaxIdleTime: 30 * time.Minute,
		},
		Redis: RedisConfig{
			URL: getEnv("REDIS_URL", "redis://localhost:6379"),
		},
		JWT: auth.TokenConfig{
			AccessTokenDuration:  15 * time.Minute,
			RefreshTokenDuration: 7 * 24 * time.Hour,
			Issuer:               "scouttalent.com",
			Audience:             []string{"api.scouttalent.com"},
			SecretKey:            getEnv("JWT_SECRET", ""),
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