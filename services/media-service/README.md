# Media Service

Video upload and management service for the ScoutTalent platform.

## Features

- **Video Upload**: Resumable uploads using TUS protocol
- **Azure Blob Storage**: Secure video storage in Azure
- **Video Management**: CRUD operations for video metadata
- **Status Tracking**: Upload progress and processing status
- **Event Publishing**: Publishes video events to NATS for other services

## API Endpoints

### Health Check
```bash
GET /health
```

### Video Upload
```bash
# 1. Initiate upload
POST /api/v1/videos/upload
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "My Football Skills",
  "description": "Showcasing my dribbling skills",
  "file_name": "skills.mp4",
  "file_size": 52428800,
  "mime_type": "video/mp4"
}

# 2. Upload video chunks (TUS protocol)
PATCH /api/v1/videos/upload/:upload_id?progress=50
Authorization: Bearer <token>

# 3. Complete upload
POST /api/v1/videos/:video_id/complete
Authorization: Bearer <token>
```

### Video Management
```bash
# Get video details
GET /api/v1/videos/:id
Authorization: Bearer <token>

# List profile videos
GET /api/v1/videos/profile/:profile_id?limit=20&offset=0
Authorization: Bearer <token>

# Update video metadata
PUT /api/v1/videos/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Updated Title",
  "description": "Updated description"
}

# Delete video
DELETE /api/v1/videos/:id
Authorization: Bearer <token>
```

## Database Schema

### Videos Table
```sql
CREATE TABLE videos (
    id UUID PRIMARY KEY,
    profile_id UUID NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    blob_url TEXT NOT NULL,
    thumbnail_url TEXT,
    duration INTEGER DEFAULT 0,
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL,
    metadata JSONB,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
```

### Uploads Table
```sql
CREATE TABLE uploads (
    id UUID PRIMARY KEY,
    video_id UUID NOT NULL REFERENCES videos(id),
    upload_id VARCHAR(255) NOT NULL UNIQUE,
    status VARCHAR(20) NOT NULL,
    progress INTEGER DEFAULT 0,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
```

## Configuration

Environment variables:

```bash
SERVER_ADDRESS=:8082
DATABASE_URL=postgres://scout:scoutpass@localhost:5432/media_db?sslmode=disable
NATS_URL=nats://localhost:4222
JWT_SECRET=your-secret-key-change-in-production
AZURE_STORAGE_ACCOUNT=your-storage-account
AZURE_STORAGE_KEY=your-storage-key
AZURE_CONTAINER_NAME=videos
LOG_LEVEL=debug
```

## Azure Blob Storage Setup

1. Create an Azure Storage Account
2. Create a container named "videos"
3. Get the account name and access key
4. Set environment variables:
   - `AZURE_STORAGE_ACCOUNT`
   - `AZURE_STORAGE_KEY`

## Development

```bash
# Run locally
cd services/media-service
go run cmd/server/main.go

# Run with Docker
docker-compose up media-service

# Run migrations
make migrate-media
```

## Testing

```bash
# Run tests
cd services/media-service
go test ./...

# Test upload flow
./scripts/test-media-api.sh
```

## Integration with Other Services

### Profile Service
- Videos are linked to user profiles via `profile_id`
- Video count affects profile completion score

### AI Moderation Worker
- Publishes `media.video.uploaded` events to NATS
- AI worker subscribes to process new videos

### Discovery Service
- Videos are indexed for search
- Metadata used for recommendations

## Video Status Flow

```
uploading → processing → ready
                ↓
              failed
```

1. **uploading**: Video chunks being uploaded
2. **processing**: Upload complete, processing video (thumbnails, metadata extraction)
3. **ready**: Video ready for viewing
4. **failed**: Upload or processing failed

## Future Enhancements

- [ ] Automatic thumbnail generation
- [ ] Video transcoding (multiple resolutions)
- [ ] Streaming optimization (HLS/DASH)
- [ ] Video analytics (views, watch time)
- [ ] CDN integration for faster delivery
- [ ] Video compression before upload