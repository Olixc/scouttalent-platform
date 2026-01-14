package config

import (
	"os"

	"github.com/scouttalent/pkg/auth"
	"github.com/scouttalent/pkg/database"
)

type Config struct {
	ServerAddress string
	Database      database.Config
	JWT           auth.JWTConfig
}

func Load() (*Config, error) {
	return &Config{
		ServerAddress: getEnv("SERVER_ADDRESS", ":8083"),
		Database: database.Config{
			URL:             getEnv("DATABASE_URL", "postgres://scout:scoutpass@localhost:5432/profile_db?sslmode=disable"),
			MaxConns:        20,
			MinConns:        5,
			MaxConnLifetime: "1h",
			MaxConnIdleTime: "30m",
		},
		JWT: auth.JWTConfig{
			Secret: getEnv("JWT_SECRET", "test-secret-key-for-development"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}