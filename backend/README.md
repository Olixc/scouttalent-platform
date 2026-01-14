# ScoutTalent Platform

An AI-powered football talent marketplace connecting players, scouts, and clubs through video showcases and intelligent matching.

## 🏗️ Monorepo Structure

```
scouttalent-platform/
├── backend/              # Go microservices backend
│   ├── services/         # Active microservices
│   │   ├── auth-service/           (Port 8080)
│   │   ├── profile-service/        (Port 8081)
│   │   ├── media-service/          (Port 8082)
│   │   ├── ai-moderation-worker/   (Background)
│   │   └── discovery-service/      (Port 8083)
│   ├── services-archive/ # Archived/WIP services
│   ├── deployment/       # Kubernetes, Helm, Terraform
│   ├── pkg/             # Shared Go packages
│   ├── docs/            # Architecture documentation
│   └── scripts/         # Testing and deployment scripts
├── frontend/            # Nuxt 3 + Tailwind CSS web app
└── .github/workflows/   # CI/CD pipelines
```

## 🚀 Quick Start

### Prerequisites

- Go 1.23+
- Node.js 20+
- PostgreSQL 16+
- Redis 7+
- NATS 2.10+

### Backend Setup

```bash
cd backend

# Start infrastructure (Docker)
docker compose up -d postgres redis nats

# Run database migrations
make migrate-up

# Start services (in separate terminals)
cd services/auth-service && go run cmd/main.go              # :8080
cd services/profile-service && go run cmd/server/main.go   # :8081
cd services/media-service && go run cmd/server/main.go     # :8082
cd services/ai-moderation-worker && go run cmd/worker/main.go
cd services/discovery-service && go run cmd/server/main.go # :8083
```

### Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev  # http://localhost:3000
```

### Run Tests

```bash
cd backend
./scripts/test-all-services.sh
```

## 📊 Active Services

| Service | Port | Status | Purpose |
|---------|------|--------|---------|
| **Auth Service** | 8080 | ✅ Active | User authentication & JWT |
| **Profile Service** | 8081 | ✅ Active | User profiles & player details |
| **Media Service** | 8082 | ✅ Active | Video upload & management |
| **AI Moderation Worker** | N/A | ✅ Active | Automated content moderation |
| **Discovery Service** | 8083 | ✅ Active | Search & recommendations |
| **Frontend** | 3000 | ✅ Active | Nuxt 3 web application |

## 🎯 Key Features

### Backend
- **Microservices Architecture** - Independent, scalable services
- **JWT Authentication** - Secure token-based auth
- **PostgreSQL per Service** - Database isolation
- **Redis Caching** - Performance optimization
- **NATS Messaging** - Event-driven communication
- **Azure Blob Storage** - Video file storage (test mode supported)
- **AI Content Moderation** - OpenAI-powered video analysis
- **Search & Discovery** - Advanced filtering and recommendations

### Frontend
- **Nuxt 3** - Modern Vue 3 framework with SSR
- **Tailwind CSS** - Utility-first styling
- **TypeScript** - Type-safe development
- **Pinia** - State management
- **Responsive Design** - Mobile-first approach
- **Authentication Flow** - Login, register, profile management
- **Video Upload** - Drag-and-drop interface
- **Search & Discovery** - Find players and videos

## 📚 Documentation

- **Backend Docs**: [backend/docs/](./backend/docs/)
- **API Documentation**: [backend/docs/API.md](./backend/docs/API.md)
- **Frontend Docs**: [frontend/README.md](./frontend/README.md)
- **Testing Guide**: [backend/TESTING_WITHOUT_DOCKER.md](./backend/TESTING_WITHOUT_DOCKER.md)
- **Completion Summary**: [PROJECT_COMPLETION_SUMMARY.md](./PROJECT_COMPLETION_SUMMARY.md)
- **Cleanup Plan**: [CLEANUP_PLAN.md](./CLEANUP_PLAN.md)

## 🔧 Configuration

### Backend Services

Each service has its own `.env` file. Example for Media Service:

```bash
SERVER_ADDRESS=:8082
DATABASE_URL=postgres://scout:scoutpass@localhost:5432/media_db?sslmode=disable
NATS_URL=nats://localhost:4222
JWT_SECRET=your-secret-key
AZURE_STORAGE_ACCOUNT=  # Leave empty for test mode
AZURE_STORAGE_KEY=
AZURE_CONTAINER_NAME=videos
```

### Frontend

Create `frontend/.env`:

```bash
NUXT_PUBLIC_API_BASE=http://localhost:8080
NUXT_PUBLIC_AUTH_SERVICE_URL=http://localhost:8080
NUXT_PUBLIC_PROFILE_SERVICE_URL=http://localhost:8081
NUXT_PUBLIC_MEDIA_SERVICE_URL=http://localhost:8082
NUXT_PUBLIC_DISCOVERY_SERVICE_URL=http://localhost:8083
```

## 🧪 Testing

### End-to-End Tests

```bash
cd backend
./scripts/test-all-services.sh
```

This script tests:
- ✅ Infrastructure (PostgreSQL, Redis, NATS)
- ✅ Service health checks
- ✅ User registration and login
- ✅ Profile creation and updates
- ✅ Video upload and management

### Frontend Tests

```bash
cd frontend
npm run test
npm run lint
npm run type-check
```

## 🚢 Deployment

### Docker

Each service has a Dockerfile:

```bash
# Build and run a service
cd backend/services/auth-service
docker build -t scouttalent/auth-service .
docker run -p 8080:8080 scouttalent/auth-service
```

### Kubernetes

Deployment manifests are in `backend/deployment/`:

```bash
# Deploy to Kubernetes
kubectl apply -f backend/deployment/k8s/

# Using Helm
helm install scouttalent backend/deployment/helm/charts/scouttalent
```

### CI/CD

GitHub Actions workflows automatically:
- Run tests on push/PR
- Build Docker images
- Deploy to staging (on develop branch)
- Deploy to production (on main branch)

See `.github/workflows/` for details.

## 🔐 Security Features

- ✅ JWT-based authentication
- ✅ Password hashing (bcrypt)
- ✅ CORS configuration
- ✅ Input validation
- ✅ SQL injection prevention
- ✅ XSS protection
- ✅ Rate limiting ready

## 📈 Architecture

### Backend Microservices

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│   Client    │────▶│ Auth Service │────▶│  PostgreSQL │
└─────────────┘     └──────────────┘     └─────────────┘
                           │
                           ▼
                    ┌──────────────┐
                    │    Redis     │
                    └──────────────┘

┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│   Client    │────▶│Profile Service│───▶│  PostgreSQL │
└─────────────┘     └──────────────┘     └─────────────┘
                           │
                           ▼
                    ┌──────────────┐
                    │     NATS     │
                    └──────────────┘

┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│   Client    │────▶│Media Service │────▶│  PostgreSQL │
└─────────────┘     └──────────────┘     └─────────────┘
                           │
                           ▼
                    ┌──────────────┐
                    │ Azure Blob   │
                    └──────────────┘
                           │
                           ▼
                    ┌──────────────┐
                    │AI Moderation │
                    │   Worker     │
                    └──────────────┘
```

### Event Flow

```
Video Upload → Media Service → NATS Event
                                    ↓
                          AI Moderation Worker
                                    ↓
                          Approve/Reject Video
                                    ↓
                          Update Video Status
```

## 🛠️ Development

### Adding a New Service

1. Create service directory in `backend/services/`
2. Follow existing service structure
3. Add to `backend/docker-compose.yml`
4. Add to CI/CD workflows
5. Update this README

### Code Style

- **Go**: Follow [Effective Go](https://golang.org/doc/effective_go.html)
- **TypeScript/Vue**: Follow [Vue Style Guide](https://vuejs.org/style-guide/)
- **Commits**: Use [Conventional Commits](https://www.conventionalcommits.org/)

## 🐛 Troubleshooting

### Port Already in Use

```bash
# Find and kill process
lsof -ti:8080 | xargs kill -9
```

### Database Connection Issues

```bash
# Test PostgreSQL connection
psql -U scout -d auth_db -c "SELECT 1;"

# Reset database
cd backend
make db-reset
```

### Service Not Starting

```bash
# Check logs
cd backend/services/auth-service
go run cmd/main.go 2>&1 | tee service.log

# Verify dependencies
go mod download
go mod verify
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
- **Issues**: [GitHub Issues](https://github.com/Olixc/scouttalent-platform/issues)
- **Discussions**: [GitHub Discussions](https://github.com/Olixc/scouttalent-platform/discussions)

## 🎉 Acknowledgments

Built with modern technologies:
- [Go](https://golang.org/)
- [Nuxt 3](https://nuxt.com/)
- [PostgreSQL](https://www.postgresql.org/)
- [Redis](https://redis.io/)
- [NATS](https://nats.io/)
- [Tailwind CSS](https://tailwindcss.com/)

---

**Made with ❤️ for the football community**