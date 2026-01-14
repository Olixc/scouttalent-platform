# Testing ScoutTalent Platform Without Docker

This guide explains how to test the ScoutTalent platform services without Docker, which is useful when Docker is not available in your environment (like Atoms platform).

## Prerequisites

- Go 1.23+ installed
- PostgreSQL 16+ installed and running
- Redis installed and running
- NATS installed and running
- `golang-migrate` CLI tool

## Setup Infrastructure Services

### 1. PostgreSQL Setup

Start PostgreSQL and create databases:

```bash
# Start PostgreSQL (if not running)
# macOS: brew services start postgresql
# Linux: sudo systemctl start postgresql

# Create databases
psql -U postgres << EOF
CREATE DATABASE auth_db;
CREATE DATABASE profile_db;
CREATE DATABASE media_db;
CREATE USER scout WITH PASSWORD 'scoutpass';
GRANT ALL PRIVILEGES ON DATABASE auth_db TO scout;
GRANT ALL PRIVILEGES ON DATABASE profile_db TO scout;
GRANT ALL PRIVILEGES ON DATABASE media_db TO scout;
EOF
```

### 2. Redis Setup

```bash
# Start Redis (if not running)
# macOS: brew services start redis
# Linux: sudo systemctl start redis

# Test connection
redis-cli ping
# Should return: PONG
```

### 3. NATS Setup

```bash
# Install NATS server
# macOS: brew install nats-server
# Linux: Download from https://github.com/nats-io/nats-server/releases

# Start NATS with JetStream
nats-server -js -m 8222

# In another terminal, verify it's running
curl http://localhost:8222/healthz
# Should return: ok
```

## Run Database Migrations

```bash
# Install golang-migrate if not installed
# macOS: brew install golang-migrate
# Linux: See https://github.com/golang-migrate/migrate/releases

# Run migrations for each service
cd /workspace

# Auth Service
migrate -path services/auth-service/migrations \
  -database "postgres://scout:scoutpass@localhost:5432/auth_db?sslmode=disable" \
  up

# Profile Service
migrate -path services/profile-service/migrations \
  -database "postgres://scout:scoutpass@localhost:5432/profile_db?sslmode=disable" \
  up

# Media Service
migrate -path services/media-service/migrations \
  -database "postgres://scout:scoutpass@localhost:5432/media_db?sslmode=disable" \
  up
```

## Configure Services

Create `.env` files for each service (or export environment variables):

### Auth Service (.env)

```bash
export SERVER_ADDRESS=:8080
export DATABASE_URL="postgres://scout:scoutpass@localhost:5432/auth_db?sslmode=disable"
export REDIS_URL="redis://localhost:6379"
export JWT_SECRET="test-secret-key-for-development"
export LOG_LEVEL="debug"
```

### Profile Service (.env)

```bash
export SERVER_ADDRESS=:8081
export DATABASE_URL="postgres://scout:scoutpass@localhost:5432/profile_db?sslmode=disable"
export NATS_URL="nats://localhost:4222"
export JWT_SECRET="test-secret-key-for-development"
export LOG_LEVEL="debug"
```

### Media Service (.env)

```bash
export SERVER_ADDRESS=:8082
export DATABASE_URL="postgres://scout:scoutpass@localhost:5432/media_db?sslmode=disable"
export NATS_URL="nats://localhost:4222"
export JWT_SECRET="test-secret-key-for-development"
export AZURE_STORAGE_ACCOUNT=""
export AZURE_STORAGE_KEY=""
export AZURE_CONTAINER_NAME="videos"
export LOG_LEVEL="debug"
```

**Note:** Leaving Azure credentials empty enables **test mode** where the service generates mock URLs without requiring actual Azure Blob Storage.

## Build and Run Services

### Option 1: Run All Services (Recommended)

Open 3 terminal windows and run each service:

**Terminal 1 - Auth Service:**
```bash
cd /workspace/services/auth-service
go mod download
go run cmd/main.go
```

**Terminal 2 - Profile Service:**
```bash
cd /workspace/services/profile-service
go mod download
go run cmd/server/main.go
```

**Terminal 3 - Media Service:**
```bash
cd /workspace/services/media-service
go mod download
go run cmd/server/main.go
```

### Option 2: Build Binaries First

```bash
# Build all services
cd /workspace

# Auth Service
cd services/auth-service
go build -o ../../bin/auth-service cmd/main.go

# Profile Service
cd ../profile-service
go build -o ../../bin/profile-service cmd/server/main.go

# Media Service
cd ../media-service
go build -o ../../bin/media-service cmd/server/main.go

# Run the binaries
cd /workspace
./bin/auth-service &
./bin/profile-service &
./bin/media-service &
```

## Verify Services are Running

```bash
# Check Auth Service
curl http://localhost:8080/health
# Expected: {"status":"healthy","service":"auth-service"}

# Check Profile Service
curl http://localhost:8081/health
# Expected: {"status":"healthy","service":"profile-service"}

# Check Media Service
curl http://localhost:8082/health
# Expected: {"status":"healthy","service":"media-service"}
```

## Test the APIs

### 1. Register a User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@player.com",
    "password": "Test123!@#",
    "full_name": "Test Player"
  }'
```

### 2. Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@player.com",
    "password": "Test123!@#"
  }'
```

Save the `access_token` from the response.

### 3. Create Profile

```bash
TOKEN="your_access_token_here"

curl -X POST http://localhost:8081/api/v1/profiles \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "bio": "Professional football player",
    "location": "London, UK",
    "profile_type": "player"
  }'
```

### 4. Upload Video (Test Mode)

```bash
curl -X POST http://localhost:8082/api/v1/videos/upload \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My Skills Video",
    "description": "Showcasing football skills",
    "file_name": "skills.mp4",
    "file_size": 52428800,
    "mime_type": "video/mp4"
  }'
```

**Note:** In test mode (no Azure credentials), the response will include `"test_mode": true` and a mock upload URL.

### 5. Complete Upload

```bash
VIDEO_ID="video_id_from_previous_response"

curl -X POST http://localhost:8082/api/v1/videos/$VIDEO_ID/complete \
  -H "Authorization: Bearer $TOKEN"
```

In test mode, the video status will immediately change to "ready" since no actual processing occurs.

## Test Mode Features

When running without Azure credentials, the Media Service operates in **test mode**:

✅ **What Works:**
- All API endpoints function normally
- Video records are created in the database
- Upload progress tracking works
- Video metadata management works
- Mock URLs are generated for uploads/downloads
- Videos are immediately marked as "ready" after completion

❌ **What Doesn't Work:**
- Actual file uploads to Azure Blob Storage
- Real video processing
- Thumbnail generation
- Video transcoding

## Troubleshooting

### Port Already in Use

If you get "address already in use" errors:

```bash
# Find process using the port
lsof -i :8080  # or :8081, :8082

# Kill the process
kill -9 <PID>
```

### Database Connection Errors

```bash
# Verify PostgreSQL is running
psql -U scout -d auth_db -c "SELECT 1;"

# Check connection string
echo $DATABASE_URL
```

### Redis Connection Errors

```bash
# Test Redis connection
redis-cli ping

# Check if Redis is listening
netstat -an | grep 6379
```

### NATS Connection Errors

```bash
# Check NATS is running
curl http://localhost:8222/healthz

# View NATS logs
# (depends on how you started NATS)
```

### Migration Errors

```bash
# Check migration status
migrate -path services/auth-service/migrations \
  -database "postgres://scout:scoutpass@localhost:5432/auth_db?sslmode=disable" \
  version

# Force to a specific version if stuck
migrate -path services/auth-service/migrations \
  -database "postgres://scout:scoutpass@localhost:5432/auth_db?sslmode=disable" \
  force 1
```

## Stopping Services

```bash
# If running in foreground, press Ctrl+C in each terminal

# If running in background
pkill -f auth-service
pkill -f profile-service
pkill -f media-service

# Stop infrastructure
# macOS:
brew services stop postgresql
brew services stop redis
pkill nats-server

# Linux:
sudo systemctl stop postgresql
sudo systemctl stop redis
pkill nats-server
```

## Next Steps

Once local testing is complete:

1. **Add Azure Credentials** for production deployment
2. **Set up CI/CD Pipeline** for automated testing
3. **Deploy to Kubernetes** for production environment
4. **Implement AI Moderation Worker** to process uploaded videos
5. **Build Discovery Service** for search and recommendations

## Useful Commands

```bash
# View service logs (if running in background)
tail -f /tmp/auth-service.log
tail -f /tmp/profile-service.log
tail -f /tmp/media-service.log

# Check database tables
psql -U scout -d media_db -c "\dt"

# View recent videos
psql -U scout -d media_db -c "SELECT id, title, status, created_at FROM videos ORDER BY created_at DESC LIMIT 5;"

# Monitor NATS
nats sub ">"  # Subscribe to all subjects
```