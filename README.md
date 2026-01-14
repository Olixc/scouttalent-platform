# ScoutTalent Platform

An AI-powered football talent marketplace connecting African youth players with global scouts.

## Architecture Overview

This is a microservices-based platform built with:
- **Backend**: Go 1.22+ with Gin framework
- **Frontend**: Next.js 15 with React 19
- **Databases**: PostgreSQL 18.1 (database-per-service)
- **Cache**: Redis 8.4
- **Message Queue**: NATS JetStream
- **Container Orchestration**: Kubernetes (AKS)
- **AI Services**: Azure Video Indexer, Face API, Content Moderator

## Project Structure

```
scouting/
├── services/          # Microservices (Go)
│   ├── auth-service/          ✅ COMPLETE
│   ├── profile-service/       ✅ COMPLETE
│   ├── media-service/         🚧 In Progress
│   └── ...
├── pkg/              # Shared Go packages
├── web/              # Next.js frontend
├── helm/             # Helm charts
├── k8s/              # Kubernetes manifests
├── terraform/        # Infrastructure as Code
└── docs/             # Architecture documentation
```

## Services

### Auth Service (Port 8080) ✅
- User registration with email/password
- Login with JWT tokens (15min access, 7day refresh)
- Role-based access control (player, scout, academy, admin)
- Password reset flow
- OAuth integration ready

### Profile Service (Port 8081) ✅
- Profile creation and management
- Player-specific details (position, stats, physical attributes)
- Scout-specific details (organization, regions of interest)
- Profile completion scoring (0-100)
- Trust level progression (newcomer → established → verified → pro)

## Getting Started

### Prerequisites

- Go 1.22+
- Docker & Docker Compose
- Node.js 20+
- golang-migrate: `go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`

### Quick Start

```bash
# 1. Start all services
make dev

# 2. Run database migrations
make migrate-up

# 3. Test the API
./scripts/test-api.sh
```

### Manual Testing

#### 1. Register a Player
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "player@example.com",
    "password": "password123",
    "role": "player"
  }'
```

#### 2. Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "player@example.com",
    "password": "password123"
  }'
```

Save the `access_token` from the response.

#### 3. Create Profile
```bash
export TOKEN="your-access-token-here"

curl -X POST http://localhost:8081/api/v1/profiles \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "display_name": "John Doe",
    "bio": "Aspiring professional footballer",
    "location_country": "Nigeria",
    "location_city": "Lagos"
  }'
```

#### 4. Add Player Details
```bash
curl -X POST http://localhost:8081/api/v1/profiles/{profile_id}/player-details \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "position": "forward",
    "height_cm": 180,
    "weight_kg": 75,
    "preferred_foot": "right",
    "current_team": "Lagos FC"
  }'
```

#### 5. Get Player Profile
```bash
curl http://localhost:8081/api/v1/profiles/{profile_id}/player \
  -H "Authorization: Bearer $TOKEN"
```

## Development Phases

### Phase 1: Foundation (Current) 🚧
- [x] Project structure
- [x] Shared packages (database, logging, auth, middleware)
- [x] Auth Service (registration, login, JWT)
- [x] Profile Service (profiles, player details, trust levels)
- [x] Docker Compose for local development
- [ ] Media Service (TUS upload)
- [ ] Basic AI Moderation

### Phase 2: AI Core (Next)
- [ ] AI Scoring Worker
- [ ] Highlight Generator
- [ ] Discovery Service
- [ ] Trust & Verification System

### Phase 3: Marketplace
- [ ] Scout Experience
- [ ] Engagement Service
- [ ] Notification Service
- [ ] Payment Service

### Phase 4: Scale & Optimize
- [ ] Custom AI Training
- [ ] Performance Optimization
- [ ] Production Deployment

## Key Features

### Profile Completion Scoring
Automatic calculation based on:
- Basic info (60 points): name, bio, avatar, location
- Type-specific (40 points): player details, scout verification

### Trust Level Progression
- **Newcomer** (default): Email verified, basic profile
- **Established**: 3+ quality videos, 30+ days, 50%+ completion
- **Verified**: Document verification, phone verified
- **Pro**: Academy endorsed, high engagement

### Database Design
- **Database-per-service**: Each service has its own PostgreSQL database
- **Event-driven sync**: Services communicate via NATS JetStream
- **Proper isolation**: No cross-database queries

## Available Commands

```bash
make dev              # Start all services
make stop             # Stop all services
make build            # Build all services
make test             # Run tests
make lint             # Run linters
make migrate-up       # Run all migrations
make migrate-down     # Rollback all migrations
make migrate-auth     # Run auth migrations only
make migrate-profile  # Run profile migrations only
make clean            # Clean build artifacts
```

## API Endpoints

### Auth Service (8080)
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login
- `GET /api/v1/auth/me` - Get current user (protected)
- `GET /health` - Health check

### Profile Service (8081)
- `POST /api/v1/profiles` - Create profile (protected)
- `GET /api/v1/profiles/me` - Get my profile (protected)
- `GET /api/v1/profiles/:id` - Get profile by ID (protected)
- `PUT /api/v1/profiles/:id` - Update profile (protected)
- `POST /api/v1/profiles/:id/player-details` - Add player details (protected)
- `GET /api/v1/profiles/:id/player` - Get player profile (protected)
- `GET /health` - Health check

## Documentation

See `docs/` directory for comprehensive architecture documentation:
- [00-PROJECT-OVERVIEW.md](docs/00-PROJECT-OVERVIEW.md)
- [01-MICROSERVICES-BREAKDOWN.md](docs/01-MICROSERVICES-BREAKDOWN.md)
- [02-TECHNOLOGY-STACK.md](docs/02-TECHNOLOGY-STACK.md)
- And more...

## Contributing

This is currently a solo project. Contributions welcome after MVP launch.

## License

Proprietary - All rights reserved