# ScoutTalent Platform

An AI-powered football talent marketplace connecting players, scouts, and clubs through video showcases and intelligent matching.

## 🏗️ Project Structure

```
scouttalent-platform/
├── backend/           # Go microservices
│   ├── services/      # Individual microservices
│   │   ├── auth-service/
│   │   ├── profile-service/
│   │   ├── media-service/
│   │   ├── ai-moderation-worker/
│   │   └── discovery-service/
│   ├── pkg/           # Shared Go packages
│   ├── docs/          # Architecture documentation
│   └── scripts/       # Testing and deployment scripts
├── frontend/          # Nuxt 3 + Tailwind CSS web app
│   ├── components/
│   ├── pages/
│   ├── composables/
│   └── assets/
└── .github/           # CI/CD workflows
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

# Start infrastructure services (if using Docker)
docker compose up -d postgres redis nats

# Run database migrations
make migrate-up

# Start all services
make dev
```

### Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

Visit `http://localhost:3000` to access the application.

## 📚 Documentation

- [Backend Documentation](./backend/docs/)
- [API Documentation](./backend/docs/API.md)
- [Frontend Documentation](./frontend/README.md)
- [Testing Guide](./backend/TESTING_WITHOUT_DOCKER.md)

## 🧪 Testing

```bash
# Backend tests
cd backend
./scripts/test-all-services.sh

# Frontend tests
cd frontend
npm run test
```

## 🔧 Development

### Backend Services

- **Auth Service** (`:8080`) - User authentication and authorization
- **Profile Service** (`:8081`) - User profiles and player details
- **Media Service** (`:8082`) - Video upload and management
- **AI Moderation Worker** - Automated content moderation
- **Discovery Service** (`:8083`) - Search and recommendations

### Frontend

- **Nuxt 3** - Vue 3 framework with SSR
- **Tailwind CSS** - Utility-first styling
- **Pinia** - State management
- **VueUse** - Composition utilities

## 🚢 Deployment

See [CI/CD Documentation](./.github/workflows/README.md) for deployment instructions.

## 📝 License

MIT License - see [LICENSE](./LICENSE) for details.

## 🤝 Contributing

Contributions are welcome! Please read our [Contributing Guide](./CONTRIBUTING.md) first.