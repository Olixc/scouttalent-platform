# AI Moderation Worker

Background worker service that processes uploaded videos using AI for content moderation.

## Features

- **Event-Driven Processing**: Listens to NATS events for new video uploads
- **AI Content Analysis**: Uses OpenAI GPT-4 for intelligent content moderation
- **Automatic Approval/Rejection**: Makes moderation decisions based on AI analysis
- **Test Mode**: Works without OpenAI API key for development/testing
- **Confidence Scoring**: Provides confidence levels for moderation decisions
- **Tag Extraction**: Automatically suggests relevant tags for videos
- **Content Summarization**: Generates brief summaries of video content

## Architecture

```
NATS Event → Worker → AI Moderator → Database Update
                ↓
         OpenAI API (optional)
```

## Configuration

Environment variables:

```bash
DATABASE_URL=postgres://scout:scoutpass@localhost:5432/media_db?sslmode=disable
NATS_URL=nats://localhost:4222
OPENAI_API_KEY=sk-...  # Optional - runs in test mode if empty
OPENAI_MODEL=gpt-4o-mini
LOG_LEVEL=debug
```

## Test Mode

When `OPENAI_API_KEY` is not set, the worker operates in test mode:

- All videos are automatically approved
- Basic keyword-based moderation is applied
- No external API calls are made
- Useful for development and testing

## Running Locally

```bash
# Set environment variables
export DATABASE_URL="postgres://scout:scoutpass@localhost:5432/media_db?sslmode=disable"
export NATS_URL="nats://localhost:4222"
export OPENAI_API_KEY=""  # Leave empty for test mode

# Run migrations
migrate -path migrations \
  -database "$DATABASE_URL" \
  up

# Run worker
go run cmd/worker/main.go
```

## Event Format

The worker subscribes to `media.video.uploaded` events:

```json
{
  "event_type": "video.uploaded",
  "video_id": "uuid",
  "profile_id": "uuid",
  "title": "Video Title",
  "timestamp": 1234567890
}
```

## Moderation Result

```json
{
  "approved": true,
  "confidence": 0.95,
  "flags": [],
  "reason": "Content approved",
  "suggested_tags": ["football", "skills", "training"],
  "content_summary": "Football skills demonstration video"
}
```

## Database Schema

### moderation_results

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| video_id | UUID | Reference to videos table |
| approved | BOOLEAN | Moderation decision |
| confidence | DECIMAL | Confidence score (0-1) |
| flags | TEXT[] | Content flags if any |
| reason | TEXT | Explanation for decision |
| result_data | JSONB | Full moderation result |
| created_at | TIMESTAMP | When moderation occurred |

## Production Setup

1. **Get OpenAI API Key**: Sign up at https://platform.openai.com/
2. **Set Environment Variable**: `OPENAI_API_KEY=sk-your-key`
3. **Choose Model**: Default is `gpt-4o-mini` (cost-effective)
4. **Monitor Usage**: Track API usage in OpenAI dashboard
5. **Set Rate Limits**: Configure appropriate rate limiting

## Moderation Criteria

The AI analyzes videos for:

- **Inappropriate Content**: Violence, hate speech, explicit material
- **Spam**: Promotional or irrelevant content
- **Quality**: Video relevance to football/sports
- **Authenticity**: Genuine player showcases vs. fake content

## Monitoring

Key metrics to monitor:

- Processing time per video
- Approval/rejection rates
- Confidence score distribution
- Failed moderation attempts
- NATS queue depth

## Troubleshooting

### Worker not receiving events

```bash
# Check NATS connection
curl http://localhost:8222/connz

# Verify subscription
curl http://localhost:8222/subsz
```

### Database connection issues

```bash
# Test database connection
psql $DATABASE_URL -c "SELECT 1;"

# Check moderation_results table
psql $DATABASE_URL -c "\d moderation_results"
```

### OpenAI API errors

- Verify API key is valid
- Check rate limits and quotas
- Monitor OpenAI status page
- Review error logs for specific issues

## Future Enhancements

- [ ] Video frame analysis using computer vision
- [ ] Multi-language support
- [ ] Custom moderation rules per organization
- [ ] Human review queue for borderline cases
- [ ] Batch processing for efficiency
- [ ] Real-time moderation status updates via WebSocket