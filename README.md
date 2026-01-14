# ScoutTalent Platform

An AI-powered football talent marketplace connecting players, scouts, and clubs through video showcases and intelligent matching.

## 🏗️ Monorepo Structure

```
scouttalent-platform/
├── backend/              # Go microservices backend
│   ├── services/         # Active microservices (5 services)
│   │   ├── auth-service/           (Port 8080)
│   │   ├── profile-service/        (Port 8081)
│   │   ├── media-service/          (Port 8082)
│   │   ├── ai-moderation-worker/   (Background)
│   │   └── discovery-service/      (Port 8083)
│   ├── services-archive/ # Archived/WIP services (5 services)
│   ├── deployment/       # Kubernetes, Helm, Terraform
│   ├── pkg/             # Shared Go packages
│   ├── docs/            # Architecture documentation
│   └── scripts/         # Testing and deployment scripts
├── frontend/            # Nuxt 3 + Tailwind CSS web app
└── .github/workflows/   # CI/CD pipelines
```

## 🚀 Quick Start - Run Locally

### Prerequisites

- **Go 1.23+** - [Download](https://go.dev/dl/)
- **Node.js 20+** - [Download](https://nodejs.org/)
- **PostgreSQL 16+** - [Download](https://www.postgresql.org/download/)
- **Redis 7+** - [Download](https://redis.io/download/)
- **NATS 2.10+** - [Download](https://nats.io/download/)

### Option 1: Using Docker (Recommended)

```bash
# 1. Start infrastructure services
cd backend
docker compose up -d postgres redis nats

# 2. Wait for services to be ready (30 seconds)
sleep 30

# 3. Initialize databases
psql -U scout -h localhost -c "CREATE DATABASE auth_db;"
psql -U scout -h localhost -c "CREATE DATABASE profile_db;"
psql -U scout -h localhost -c "CREATE DATABASE media_db;"
psql -U scout -h localhost -c "CREATE DATABASE discovery_db;"
psql -U scout -h localhost -c "CREATE DATABASE moderation_db;"

# 4. Run migrations for each service
cd services/auth-service && migrate -path migrations -database "postgres://scout:scoutpass@localhost:5432/auth_db?sslmode=disable" up
cd ../profile-service && migrate -path migrations -database "postgres://scout:scoutpass@localhost:5432/profile_db?sslmode=disable" up
cd ../media-service && migrate -path migrations -database "postgres://scout:scoutpass@localhost:5432/media_db?sslmode=disable" up
cd ../ai-moderation-worker && migrate -path migrations -database "postgres://scout:scoutpass@localhost:5432/moderation_db?sslmode=disable" up

# 5. Start backend services (open 5 terminal tabs)
# Terminal 1 - Auth Service
cd backend/services/auth-service
go run cmd/main.go

# Terminal 2 - Profile Service
cd backend/services/profile-service
go run cmd/server/main.go

# Terminal 3 - Media Service
cd backend/services/media-service
go run cmd/server/main.go

# Terminal 4 - AI Moderation Worker
cd backend/services/ai-moderation-worker
go run cmd/worker/main.go

# Terminal 5 - Discovery Service
cd backend/services/discovery-service
go run cmd/server/main.go

# 6. Start frontend (new terminal)
cd frontend
npm install
npm run dev
```

### Option 2: Without Docker

#### Step 1: Install and Start PostgreSQL

```bash
# macOS
brew install postgresql@16
brew services start postgresql@16

# Ubuntu/Debian
sudo apt install postgresql-16
sudo systemctl start postgresql

# Create databases
createdb -U postgres auth_db
createdb -U postgres profile_db
createdb -U postgres media_db
createdb -U postgres discovery_db
createdb -U postgres moderation_db
```

#### Step 2: Install and Start Redis

```bash
# macOS
brew install redis
brew services start redis

# Ubuntu/Debian
sudo apt install redis-server
sudo systemctl start redis
```

#### Step 3: Install and Start NATS

```bash
# macOS
brew install nats-server
nats-server &

# Ubuntu/Debian
wget https://github.com/nats-io/nats-server/releases/download/v2.10.7/nats-server-v2.10.7-linux-amd64.tar.gz
tar -xzf nats-server-v2.10.7-linux-amd64.tar.gz
./nats-server-v2.10.7-linux-amd64/nats-server &
```

#### Step 4: Configure Environment Variables

Create `.env` files for each service:

**backend/services/auth-service/.env**
```bash
SERVER_ADDRESS=:8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable
REDIS_URL=localhost:6379
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRY=24h
```

**backend/services/profile-service/.env**
```bash
SERVER_ADDRESS=:8081
DATABASE_URL=postgres://postgres:postgres@localhost:5432/profile_db?sslmode=disable
REDIS_URL=localhost:6379
JWT_SECRET=your-super-secret-jwt-key-change-in-production
```

**backend/services/media-service/.env**
```bash
SERVER_ADDRESS=:8082
DATABASE_URL=postgres://postgres:postgres@localhost:5432/media_db?sslmode=disable
NATS_URL=nats://localhost:4222
JWT_SECRET=your-super-secret-jwt-key-change-in-production
AZURE_STORAGE_ACCOUNT=
AZURE_STORAGE_KEY=
AZURE_CONTAINER_NAME=videos
```

**backend/services/ai-moderation-worker/.env**
```bash
DATABASE_URL=postgres://postgres:postgres@localhost:5432/moderation_db?sslmode=disable
NATS_URL=nats://localhost:4222
OPENAI_API_KEY=
TEST_MODE=true
```

**backend/services/discovery-service/.env**
```bash
SERVER_ADDRESS=:8083
DATABASE_URL=postgres://postgres:postgres@localhost:5432/discovery_db?sslmode=disable
REDIS_URL=localhost:6379
JWT_SECRET=your-super-secret-jwt-key-change-in-production
```

#### Step 5: Run Database Migrations

```bash
# Install golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations for each service
cd backend/services/auth-service
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable" up

cd ../profile-service
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/profile_db?sslmode=disable" up

cd ../media-service
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/media_db?sslmode=disable" up

cd ../ai-moderation-worker
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/moderation_db?sslmode=disable" up
```

#### Step 6: Start All Services

Open 5 terminal windows and run:

```bash
# Terminal 1 - Auth Service (Port 8080)
cd backend/services/auth-service
go run cmd/main.go

# Terminal 2 - Profile Service (Port 8081)
cd backend/services/profile-service
go run cmd/server/main.go

# Terminal 3 - Media Service (Port 8082)
cd backend/services/media-service
go run cmd/server/main.go

# Terminal 4 - AI Moderation Worker (Background)
cd backend/services/ai-moderation-worker
go run cmd/worker/main.go

# Terminal 5 - Discovery Service (Port 8083)
cd backend/services/discovery-service
go run cmd/server/main.go
```

#### Step 7: Start Frontend

```bash
cd frontend
npm install
npm run dev
```

Visit **http://localhost:3000** in your browser.

## 🧪 Testing the Platform

### Automated Testing

```bash
cd backend
./scripts/test-all-services.sh
```

This script will:
- ✅ Check infrastructure (PostgreSQL, Redis, NATS)
- ✅ Test all service health endpoints
- ✅ Run end-to-end user flow tests
- ✅ Verify video upload and moderation

### Manual Testing

#### 1. Register a User
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "player@example.com",
    "password": "SecurePass123!",
    "profile_type": "player"
  }'
```

#### 2. Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "player@example.com",
    "password": "SecurePass123!"
  }'
```

Save the returned JWT token.

#### 3. Create Profile
```bash
curl -X POST http://localhost:8081/api/profiles \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "John Doe",
    "profile_type": "player",
    "bio": "Professional footballer",
    "location": "London, UK"
  }'
```

#### 4. Upload Video
```bash
curl -X POST http://localhost:8082/api/videos/upload \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "file=@/path/to/video.mp4" \
  -F "title=My Skills Video" \
  -F "description=Showcasing my football skills"
```

#### 5. Search Players
```bash
curl "http://localhost:8083/api/search/profiles?query=footballer&profile_type=player" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 📊 Service Architecture

### Active Services

| Service | Port | Purpose | Database |
|---------|------|---------|----------|
| **Auth Service** | 8080 | User authentication & JWT | auth_db |
| **Profile Service** | 8081 | User profiles & player details | profile_db |
| **Media Service** | 8082 | Video upload & management | media_db |
| **AI Moderation Worker** | N/A | Automated content moderation | moderation_db |
| **Discovery Service** | 8083 | Search & recommendations | discovery_db |
| **Frontend** | 3000 | Nuxt 3 web application | N/A |

### Infrastructure

- **PostgreSQL** - 5 separate databases (one per service)
- **Redis** - Caching and session management
- **NATS** - Event-driven messaging between services
- **Azure Blob Storage** - Video file storage (optional, test mode available)

## 🎯 Key Features

### Backend
- ✅ Microservices architecture with service isolation
- ✅ JWT-based authentication
- ✅ Database per service pattern
- ✅ Event-driven communication via NATS
- ✅ Redis caching for performance
- ✅ AI-powered video moderation
- ✅ Advanced search and recommendations
- ✅ Test mode for development without external dependencies

### Frontend
- ✅ Nuxt 3 with SSR support
- ✅ Tailwind CSS for styling
- ✅ TypeScript for type safety
- ✅ Pinia state management
- ✅ Responsive design
- ✅ Authentication flow
- ✅ Video upload interface
- ✅ Search and discovery features

## 🔧 Configuration

### Test Mode (No External Dependencies)

The platform can run in test mode without Azure Blob Storage or OpenAI:

**Media Service** - Set in `.env`:
```bash
AZURE_STORAGE_ACCOUNT=  # Leave empty
AZURE_STORAGE_KEY=      # Leave empty
```
Videos will be stored locally in `uploads/` directory.

**AI Moderation Worker** - Set in `.env`:
```bash
OPENAI_API_KEY=         # Leave empty
TEST_MODE=true
```
Videos will be auto-approved without AI analysis.

### Production Mode

For production deployment:

1. **Azure Blob Storage**:
   - Create Azure Storage Account
   - Get account name and key
   - Update Media Service `.env`

2. **OpenAI API**:
   - Get API key from OpenAI
   - Update AI Moderation Worker `.env`
   - Set `TEST_MODE=false`

## 📚 Documentation

- **Architecture**: [backend/docs/](./backend/docs/)
- **Testing Guide**: [backend/TESTING_WITHOUT_DOCKER.md](./backend/TESTING_WITHOUT_DOCKER.md)
- **Quick Start**: [backend/QUICKSTART.md](./backend/QUICKSTART.md)
- **Frontend**: [frontend/README.md](./frontend/README.md)
- **Cleanup Plan**: [CLEANUP_PLAN.md](./CLEANUP_PLAN.md)
- **Project Summary**: [PROJECT_COMPLETION_SUMMARY.md](./PROJECT_COMPLETION_SUMMARY.md)

## 🐛 Troubleshooting

### Port Already in Use
```bash
# Find and kill process
lsof -ti:8080 | xargs kill -9
```

### Database Connection Failed
```bash
# Check PostgreSQL is running
pg_isready

# Test connection
psql -U postgres -d auth_db -c "SELECT 1;"
```

### Redis Connection Failed
```bash
# Check Redis is running
redis-cli ping
# Should return: PONG
```

### NATS Connection Failed
```bash
# Check NATS is running
ps aux | grep nats-server

# Test connection
nats-server --version
```

### Service Won't Start
```bash
# Check logs
cd backend/services/auth-service
go run cmd/main.go 2>&1 | tee service.log

# Verify dependencies
go mod download
go mod verify
```

### Frontend Build Errors
```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
npm run dev
```

## 🚢 Deployment

### Docker Compose (All Services)

```bash
cd backend
docker compose up -d
```

### Kubernetes

```bash
# Apply manifests
kubectl apply -f backend/deployment/k8s/

# Using Helm
helm install scouttalent backend/deployment/helm/charts/scouttalent
```

### Terraform (Infrastructure)

```bash
cd backend/deployment/terraform/environments/dev
terraform init
terraform plan
terraform apply
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

MIT License - see [LICENSE](./LICENSE) for details.

## 📞 Support

- **Documentation**: [docs/](./backend/docs/)
- **Issues**: [GitHub Issues](https://github.com/yourusername/scouttalent-platform/issues)

---

**Made with ❤️ for the football community**