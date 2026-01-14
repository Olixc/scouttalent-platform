# ScoutTalent Platform - Local Testing Guide

This guide provides step-by-step instructions for setting up and testing the ScoutTalent microservices platform locally.

## Prerequisites

- Docker and Docker Compose installed
- Git installed
- `curl` or Postman for API testing
- `golang-migrate` CLI tool (for running migrations)

## Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/Olixc/scouttalent-platform.git
cd scouttalent-platform
```

### 2. Start Infrastructure Services

Start PostgreSQL, Redis, and NATS:

```bash
docker-compose up -d postgres redis nats
```

Wait for services to be healthy (about 30 seconds):

```bash
docker-compose ps
```

### 3. Run Database Migrations

Install golang-migrate if you haven't:

```bash
# macOS
brew install golang-migrate

# Linux
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/

# Windows
choco install golang-migrate
```

Run migrations for all services:

```bash
make migrate-up
```

Or run individually:

```bash
make migrate-auth
make migrate-profile
make migrate-media
```

### 4. Start Application Services

Build and start all services:

```bash
docker-compose up --build auth-service profile-service media-service
```

Or start individually:

```bash
docker-compose up --build auth-service
docker-compose up --build profile-service
docker-compose up --build media-service
```

### 5. Verify Services are Running

Check service health:

```bash
# Auth Service (Port 8080)
curl http://localhost:8080/health

# Profile Service (Port 8081)
curl http://localhost:8081/health

# Media Service (Port 8082)
curl http://localhost:8082/health
```

Expected response: `{"status":"healthy","service":"<service-name>"}`

## Testing the Services

### Option 1: Using Test Scripts

We provide automated test scripts for each service:

```bash
# Test Auth Service
./scripts/test-api.sh

# Test Media Service
./scripts/test-media-api.sh
```

### Option 2: Manual Testing with cURL

#### Auth Service Tests

**1. Register a new user:**

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "player@test.com",
    "password": "Test123!@#",
    "full_name": "John Doe"
  }'
```

**2. Login:**

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "player@test.com",
    "password": "Test123!@#"
  }'
```

Save the `access_token` from the response.

**3. Get current user info:**

```bash
curl -X GET http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

#### Profile Service Tests

**1. Create a profile:**

```bash
curl -X POST http://localhost:8081/api/v1/profiles \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "bio": "Professional football player",
    "location": "London, UK",
    "profile_type": "player"
  }'
```

**2. Get my profile:**

```bash
curl -X GET http://localhost:8081/api/v1/profiles/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**3. Update profile:**

```bash
curl -X PUT http://localhost:8081/api/v1/profiles/PROFILE_ID \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "bio": "Updated bio",
    "location": "Manchester, UK"
  }'
```

**4. Create player details:**

```bash
curl -X POST http://localhost:8081/api/v1/profiles/PROFILE_ID/player-details \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "position": "Forward",
    "preferred_foot": "Right",
    "height": 180,
    "weight": 75,
    "date_of_birth": "2000-01-15"
  }'
```

#### Media Service Tests

**1. Initiate video upload:**

```bash
curl -X POST http://localhost:8082/api/v1/videos/upload \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My Skills Video",
    "description": "Showcasing my football skills",
    "file_name": "skills.mp4",
    "file_size": 52428800,
    "mime_type": "video/mp4"
  }'
```

Save the `video_id` and `upload_id` from the response.

**2. Update upload progress:**

```bash
curl -X PATCH http://localhost:8082/api/v1/videos/upload/UPLOAD_ID?progress=50 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**3. Complete upload:**

```bash
curl -X POST http://localhost:8082/api/v1/videos/VIDEO_ID/complete \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**4. Get video details:**

```bash
curl -X GET http://localhost:8082/api/v1/videos/VIDEO_ID \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**5. List profile videos:**

```bash
curl -X GET http://localhost:8082/api/v1/videos/profile/PROFILE_ID?limit=10&offset=0 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**6. Update video metadata:**

```bash
curl -X PUT http://localhost:8082/api/v1/videos/VIDEO_ID \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Title",
    "description": "Updated description"
  }'
```

**7. Delete video:**

```bash
curl -X DELETE http://localhost:8082/api/v1/videos/VIDEO_ID \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## Viewing Logs

View logs for specific services:

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f auth-service
docker-compose logs -f profile-service
docker-compose logs -f media-service

# Infrastructure
docker-compose logs -f postgres
docker-compose logs -f redis
docker-compose logs -f nats
```

## Database Access

Connect to PostgreSQL:

```bash
docker exec -it scouttalent-postgres psql -U scout -d auth_db
docker exec -it scouttalent-postgres psql -U scout -d profile_db
docker exec -it scouttalent-postgres psql -U scout -d media_db
```

Useful SQL queries:

```sql
-- List all users
SELECT id, email, full_name, created_at FROM users;

-- List all profiles
SELECT id, user_id, bio, profile_type, created_at FROM profiles;

-- List all videos
SELECT id, profile_id, title, status, created_at FROM videos;
```

## Redis Access

Connect to Redis:

```bash
docker exec -it scouttalent-redis redis-cli
```

Useful Redis commands:

```redis
# List all keys
KEYS *

# Get a specific key
GET key_name

# Monitor all commands
MONITOR
```

## NATS Access

Access NATS monitoring UI:

```
http://localhost:8222
```

## Troubleshooting

### Services Won't Start

1. Check if ports are already in use:
   ```bash
   lsof -i :8080  # Auth service
   lsof -i :8081  # Profile service
   lsof -i :8082  # Media service
   lsof -i :5432  # PostgreSQL
   lsof -i :6379  # Redis
   lsof -i :4222  # NATS
   ```

2. Check Docker logs:
   ```bash
   docker-compose logs
   ```

3. Rebuild services:
   ```bash
   docker-compose down
   docker-compose up --build
   ```

### Database Connection Issues

1. Ensure PostgreSQL is healthy:
   ```bash
   docker-compose ps postgres
   ```

2. Check database exists:
   ```bash
   docker exec -it scouttalent-postgres psql -U scout -c "\l"
   ```

3. Verify connection string in docker-compose.yml

### Migration Errors

1. Check migration status:
   ```bash
   migrate -path services/auth-service/migrations \
     -database "postgres://scout:scoutpass@localhost:5432/auth_db?sslmode=disable" \
     version
   ```

2. Force migration version (if stuck):
   ```bash
   migrate -path services/auth-service/migrations \
     -database "postgres://scout:scoutpass@localhost:5432/auth_db?sslmode=disable" \
     force VERSION_NUMBER
   ```

### JWT Token Issues

1. Ensure JWT_SECRET is set correctly in docker-compose.yml
2. Token expires after 15 minutes - get a new one by logging in again
3. Check token format: `Bearer <token>`

## Stopping Services

Stop all services:

```bash
docker-compose down
```

Stop and remove volumes (clean slate):

```bash
docker-compose down -v
```

## Next Steps

After successfully testing locally:

1. **Configure Azure Blob Storage** for Media Service:
   - Create Azure Storage Account
   - Set `AZURE_STORAGE_ACCOUNT` and `AZURE_STORAGE_KEY` environment variables

2. **Set up CI/CD Pipeline** (see `docs/04-CICD-PIPELINE.md`)

3. **Deploy to Kubernetes** (see `docs/03-KUBERNETES-ARCHITECTURE.md`)

4. **Implement remaining services**:
   - AI Moderation Worker
   - Discovery Service
   - Notification Service
   - Analytics Service

## Useful Commands

```bash
# Build all services
make build

# Run all tests
make test

# Run linters
make lint

# Clean up
make clean

# View all available commands
make help
```

## Support

For issues or questions:
- Check documentation in `docs/` directory
- Review service-specific README files
- Check GitHub Issues: https://github.com/Olixc/scouttalent-platform/issues