# ScoutTalent Frontend

Modern web application built with Nuxt 3 and Tailwind CSS for the ScoutTalent platform.

## 🚀 Features

- **Authentication**: User registration, login, and session management
- **Profile Management**: Create and edit player/scout profiles
- **Video Upload**: Upload and manage football showcase videos
- **Discovery**: Search and discover players and videos
- **Feed**: Browse latest videos and trending content
- **Recommendations**: Personalized profile and video recommendations
- **Responsive Design**: Mobile-first design with Tailwind CSS

## 🛠️ Tech Stack

- **Nuxt 3**: Vue 3 framework with SSR support
- **TypeScript**: Type-safe development
- **Tailwind CSS**: Utility-first CSS framework
- **Pinia**: State management
- **VueUse**: Composition API utilities
- **Nuxt Image**: Optimized image handling

## 📦 Setup

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

## 🔧 Configuration

Environment variables (create `.env` file):

```bash
NUXT_PUBLIC_API_BASE=http://localhost:8080
NUXT_PUBLIC_AUTH_SERVICE_URL=http://localhost:8080
NUXT_PUBLIC_PROFILE_SERVICE_URL=http://localhost:8081
NUXT_PUBLIC_MEDIA_SERVICE_URL=http://localhost:8082
NUXT_PUBLIC_DISCOVERY_SERVICE_URL=http://localhost:8083
```

## 📁 Project Structure

```
frontend/
├── assets/          # CSS, images, fonts
├── components/      # Reusable Vue components
├── composables/     # Composition API utilities
├── layouts/         # Page layouts
├── pages/           # Application pages (auto-routed)
├── plugins/         # Nuxt plugins
├── public/          # Static files
├── stores/          # Pinia stores
└── types/           # TypeScript type definitions
```

## 🎨 Design System

### Colors

- **Primary**: Blue (#0ea5e9) - Main brand color
- **Secondary**: Stone gray (#78716c) - Supporting elements
- **Success**: Green - Success states
- **Error**: Red - Error states
- **Warning**: Yellow - Warning states

### Typography

- **Font Family**: Inter (Google Fonts)
- **Headings**: Bold, various sizes
- **Body**: Regular, 16px base

### Components

- **Buttons**: `.btn`, `.btn-primary`, `.btn-secondary`, `.btn-outline`
- **Cards**: `.card` - Rounded with shadow
- **Inputs**: `.input` - Styled form inputs

## 🔐 Authentication Flow

1. User registers/logs in via Auth Service
2. JWT token stored in localStorage and Pinia store
3. Token included in all authenticated API requests
4. Middleware protects authenticated routes

## 📱 Pages

- `/` - Landing page
- `/login` - User login
- `/register` - User registration
- `/profile` - User profile management
- `/upload` - Video upload
- `/discover` - Discover players
- `/search` - Search functionality
- `/feed` - Video feed

## 🧪 Testing

```bash
# Run tests
npm run test

# Run tests with coverage
npm run test:coverage

# Run E2E tests
npm run test:e2e
```

## 🚀 Deployment

### Vercel (Recommended)

```bash
# Install Vercel CLI
npm i -g vercel

# Deploy
vercel
```

### Docker

```bash
# Build image
docker build -t scouttalent-frontend .

# Run container
docker run -p 3000:3000 scouttalent-frontend
```

### Static Hosting

```bash
# Generate static site
npm run generate

# Deploy dist/ folder to any static host
```

## 🔄 API Integration

The frontend communicates with backend microservices:

- **Auth Service** (`:8080`): Authentication and user management
- **Profile Service** (`:8081`): User profiles and player details
- **Media Service** (`:8082`): Video upload and management
- **Discovery Service** (`:8083`): Search and recommendations

## 📝 Development Guidelines

1. **Component Naming**: Use PascalCase for components
2. **Composables**: Use camelCase with `use` prefix
3. **Types**: Define interfaces in `types/` directory
4. **API Calls**: Use `$fetch` with proper error handling
5. **State Management**: Use Pinia stores for global state
6. **Styling**: Prefer Tailwind utility classes

## 🐛 Troubleshooting

### Port already in use

```bash
# Kill process on port 3000
lsof -ti:3000 | xargs kill -9
```

### Module not found

```bash
# Clear node_modules and reinstall
rm -rf node_modules package-lock.json
npm install
```

### Build errors

```bash
# Clear Nuxt cache
rm -rf .nuxt
npm run dev
```

## 📚 Resources

- [Nuxt 3 Documentation](https://nuxt.com/docs)
- [Tailwind CSS Documentation](https://tailwindcss.com/docs)
- [Pinia Documentation](https://pinia.vuejs.org/)
- [VueUse Documentation](https://vueuse.org/)

## 🤝 Contributing

1. Create a feature branch
2. Make your changes
3. Run tests and linting
4. Submit a pull request

## 📄 License

MIT License - see LICENSE file for details