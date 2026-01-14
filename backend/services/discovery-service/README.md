# Discovery Service

Service for searching, discovering, and getting recommendations for players and videos.

## Features

- **Profile Search**: Search for players, scouts, and clubs
- **Video Search**: Find videos by title, description, or tags
- **Recommendations**: Get personalized profile and video recommendations
- **Feed**: Browse latest videos from the platform
- **Trending**: Discover most popular videos

## API Endpoints

### Search

```bash
# Search profiles
GET /api/v1/search/profiles?q=forward&profile_type=player&position=Forward&location=London&limit=20&offset=0

# Search videos
GET /api/v1/search/videos?q=skills&limit=20&offset=0
```

### Recommendations (Authenticated)

```bash
# Get profile recommendations
GET /api/v1/recommendations/profiles?limit=10
Authorization: Bearer <token>

# Get video recommendations
GET /api/v1/recommendations/videos?limit=10
Authorization: Bearer <token>
```

### Feed (Authenticated)

```bash
# Get video feed
GET /api/v1/feed?limit=20&offset=0
Authorization: Bearer <token>

# Get trending videos
GET /api/v1/feed/trending?limit=10
Authorization: Bearer <token>
```

## Configuration

```bash
SERVER_ADDRESS=:8083
DATABASE_URL=postgres://scout:scoutpass@localhost:5432/profile_db?sslmode=disable
JWT_SECRET=your-secret-key
LOG_LEVEL=debug
```

## Running Locally

```bash
# Set environment variables
export DATABASE_URL="postgres://scout:scoutpass@localhost:5432/profile_db?sslmode=disable"
export JWT_SECRET="test-secret-key"

# Run service
go run cmd/server/main.go
```

## Search Filters

### Profile Search
- `q`: Search query (searches bio and location)
- `profile_type`: Filter by type (player, scout, club)
- `position`: Filter by player position
- `location`: Filter by location
- `limit`: Results per page (default: 20)
- `offset`: Pagination offset (default: 0)

### Video Search
- `q`: Search query (searches title and description)
- `limit`: Results per page (default: 20)
- `offset`: Pagination offset (default: 0)

## Recommendation Algorithm

The service uses a simple collaborative filtering approach:

1. **Profile Recommendations**: Finds profiles with similar characteristics (type, position, location)
2. **Video Recommendations**: Suggests videos from similar profiles based on user's profile type

Future enhancements:
- Machine learning-based recommendations
- User interaction tracking (views, likes, saves)
- Content-based filtering using video tags
- Hybrid recommendation system

## Future Features

- [ ] Advanced search with multiple filters
- [ ] Fuzzy search for better matching
- [ ] Elasticsearch integration for better search performance
- [ ] User preference learning
- [ ] A/B testing for recommendation algorithms
- [ ] Search analytics and insights