# Quick Start Guide

Get the ScoutTalent platform running locally in 5 minutes.

## Prerequisites

Install these tools:
```bash
# Go 1.22+
brew install go  # macOS
# or download from https://go.dev/dl/

# Docker Desktop
# Download from https://www.docker.com/products/docker-desktop

# golang-migrate (for database migrations)
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

## Step 1: Clone and Setup

```bash
git clone <your-repo-url>
cd scouting
```

## Step 2: Start Infrastructure

```bash
# Start PostgreSQL and Redis
docker-compose up -d postgres redis

# Wait for services to be healthy (about 10 seconds)
docker-compose ps
```

## Step 3: Run Database Migrations

```bash
# Run auth service migrations
migrate -path services/auth-service/migrations \
  -database "postgres://scout:scoutpass@localhost:5432/auth_db?sslmode=disable" up

# Verify tables were created
docker-compose exec postgres psql -U scout -d auth_db -c "\dt"
```

## Step 4: Start Auth Service

```bash
cd services/auth-service

# Set environment variables
export DATABASE_URL="postgres://scout:scoutpass@localhost:5432/auth_db?sslmode=disable"
export REDIS_URL="redis://localhost:6379"
export JWT_SECRET="your-secret-key-change-in-production"
export LOG_LEVEL="debug"

# Run the service
go run cmd/server/main.go
```

You should see:
```
{"level":"info","timestamp":"...","service":"auth-service","msg":"connected to database"}
{"level":"info","timestamp":"...","service":"auth-service","msg":"starting server","address":":8080"}
```

## Step 5: Test the API

Open a new terminal and test the endpoints:

### Register a Player
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "player@example.com",
    "password": "password123",
    "role": "player"
  }'
```

Expected response:
```json
{
  "user": {
    "id": "...",
    "email": "player@example.com",
    "role": "player",
    ...
  },
  "message": "Registration successful. Please verify your email."
}
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "player@example.com",
    "password": "password123"
  }'
```

Expected response:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "...",
    "email": "player@example.com",
    "role": "player"
  }
}
```

### Access Protected Endpoint
```bash
# Copy the access_token from login response
export TOKEN="your-access-token-here"

curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer $TOKEN"
```

Expected response:
```json
{
  "user_id": "...",
  "message": "Authenticated"
}
```

## Step 6: Check Health

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "healthy",
  "service": "auth-service"
}
```

## Troubleshooting

### Database connection failed
```bash
# Check if PostgreSQL is running
docker-compose ps postgres

# Check logs
docker-compose logs postgres

# Restart if needed
docker-compose restart postgres
```

### Port already in use
```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>
```

### Migrations failed
```bash
# Check current migration version
migrate -path services/auth-service/migrations \
  -database "postgres://scout:scoutpass@localhost:5432/auth_db?sslmode=disable" version

# Force to specific version if needed
migrate -path services/auth-service/migrations \
  -database "postgres://scout:scoutpass@localhost:5432/auth_db?sslmode=disable" force 1
```

## Next Steps

1. **Add more services**: Profile Service, Media Service
2. **Setup frontend**: Next.js web application
3. **Deploy to Kubernetes**: Use Kind for local K8s testing
4. **Add monitoring**: Prometheus, Grafana

See [README.md](README.md) for full documentation.