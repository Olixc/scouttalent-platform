.PHONY: help dev stop build test lint migrate-up migrate-down clean

help:
	@echo "ScoutTalent Development Commands"
	@echo ""
	@echo "  make dev              - Start local development environment"
	@echo "  make stop             - Stop all services"
	@echo "  make build            - Build all services"
	@echo "  make test             - Run tests"
	@echo "  make lint             - Run linters"
	@echo "  make migrate-up       - Run all database migrations"
	@echo "  make migrate-down     - Rollback all database migrations"
	@echo "  make migrate-auth     - Run auth service migrations"
	@echo "  make migrate-profile  - Run profile service migrations"
	@echo "  make clean            - Clean build artifacts"

dev:
	@echo "Starting local development environment..."
	docker-compose up -d
	@echo "Services started:"
	@echo "  - PostgreSQL: localhost:5432"
	@echo "  - Redis: localhost:6379"
	@echo "  - NATS: localhost:4222"
	@echo "  - Auth Service: localhost:8080"
	@echo "  - Profile Service: localhost:8081"

stop:
	@echo "Stopping all services..."
	docker-compose down

build:
	@echo "Building services..."
	@cd services/auth-service && go build -o bin/auth-service ./cmd/server
	@cd services/profile-service && go build -o bin/profile-service ./cmd/server
	@echo "Build complete!"

test:
	@echo "Running tests..."
	@cd services/auth-service && go test -v -race ./...
	@cd services/profile-service && go test -v -race ./...

lint:
	@echo "Running linters..."
	@cd services/auth-service && golangci-lint run ./...
	@cd services/profile-service && golangci-lint run ./...

migrate-up: migrate-auth migrate-profile
	@echo "All migrations complete!"

migrate-down:
	@echo "Rolling back migrations..."
	@migrate -path services/profile-service/migrations \
		-database "postgres://scout:scoutpass@localhost:5432/profile_db?sslmode=disable" down
	@migrate -path services/auth-service/migrations \
		-database "postgres://scout:scoutpass@localhost:5432/auth_db?sslmode=disable" down

migrate-auth:
	@echo "Running auth service migrations..."
	@migrate -path services/auth-service/migrations \
		-database "postgres://scout:scoutpass@localhost:5432/auth_db?sslmode=disable" up

migrate-profile:
	@echo "Running profile service migrations..."
	@migrate -path services/profile-service/migrations \
		-database "postgres://scout:scoutpass@localhost:5432/profile_db?sslmode=disable" up

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf services/*/bin
	@rm -rf services/*/vendor
	@echo "Clean complete!"

.DEFAULT_GOAL := help