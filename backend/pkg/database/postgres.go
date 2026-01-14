package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	URL             string
	DSN             string
	MaxConns        int
	MinConns        int
	MaxConnLifetime string
	MaxConnIdleTime string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// NewPool creates a new PostgreSQL connection pool
func NewPool(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	// Support both URL and DSN
	connString := cfg.URL
	if connString == "" {
		connString = cfg.DSN
	}

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Connection pool settings
	if cfg.MaxConns > 0 {
		config.MaxConns = int32(cfg.MaxConns)
	} else if cfg.MaxOpenConns > 0 {
		config.MaxConns = int32(cfg.MaxOpenConns)
	}

	if cfg.MinConns > 0 {
		config.MinConns = int32(cfg.MinConns)
	} else if cfg.MaxIdleConns > 0 {
		config.MinConns = int32(cfg.MaxIdleConns)
	}

	// Parse duration strings if provided
	if cfg.MaxConnLifetime != "" {
		if duration, err := time.ParseDuration(cfg.MaxConnLifetime); err == nil {
			config.MaxConnLifetime = duration
		}
	} else if cfg.ConnMaxLifetime > 0 {
		config.MaxConnLifetime = cfg.ConnMaxLifetime
	}

	if cfg.MaxConnIdleTime != "" {
		if duration, err := time.ParseDuration(cfg.MaxConnIdleTime); err == nil {
			config.MaxConnIdleTime = duration
		}
	} else if cfg.ConnMaxIdleTime > 0 {
		config.MaxConnIdleTime = cfg.ConnMaxIdleTime
	}

	// Health check configuration
	config.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}

// DefaultConfig returns default database configuration
func DefaultConfig() Config {
	return Config{
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 1 * time.Hour,
		ConnMaxIdleTime: 30 * time.Minute,
	}
}