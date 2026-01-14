package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/scouttalent/ai-moderation-worker/internal/moderator"
	"go.uber.org/zap"
)

type VideoUploadedEvent struct {
	EventType string `json:"event_type"`
	VideoID   string `json:"video_id"`
	ProfileID string `json:"profile_id"`
	Title     string `json:"title"`
	Timestamp int64  `json:"timestamp"`
}

type Worker struct {
	db         *pgxpool.Pool
	nats       *nats.Conn
	moderator  *moderator.AIModerator
	logger     *zap.Logger
	sub        *nats.Subscription
}

func NewWorker(db *pgxpool.Pool, nc *nats.Conn, mod *moderator.AIModerator, logger *zap.Logger) *Worker {
	return &Worker{
		db:        db,
		nats:      nc,
		moderator: mod,
		logger:    logger,
	}
}

func (w *Worker) Start(ctx context.Context) error {
	// Subscribe to video upload events
	sub, err := w.nats.Subscribe("media.video.uploaded", w.handleVideoUpload)
	if err != nil {
		return fmt.Errorf("failed to subscribe to NATS: %w", err)
	}

	w.sub = sub
	w.logger.Info("Subscribed to media.video.uploaded events")

	return nil
}

func (w *Worker) Stop() error {
	if w.sub != nil {
		if err := w.sub.Unsubscribe(); err != nil {
			return err
		}
	}
	return nil
}

func (w *Worker) handleVideoUpload(msg *nats.Msg) {
	w.logger.Info("Received video upload event", zap.String("data", string(msg.Data)))

	var event VideoUploadedEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		w.logger.Error("Failed to parse event", zap.Error(err))
		return
	}

	ctx := context.Background()

	// Get video details from database
	video, err := w.getVideo(ctx, event.VideoID)
	if err != nil {
		w.logger.Error("Failed to get video", zap.Error(err))
		return
	}

	// Moderate video content
	result, err := w.moderator.ModerateVideo(ctx, event.VideoID, video.Title, video.Description)
	if err != nil {
		w.logger.Error("Failed to moderate video", zap.Error(err))
		w.updateVideoStatus(ctx, event.VideoID, "failed", "Moderation failed")
		return
	}

	// Update video status based on moderation result
	if result.Approved {
		w.updateVideoStatus(ctx, event.VideoID, "approved", result.Reason)
		w.logger.Info("Video approved",
			zap.String("video_id", event.VideoID),
			zap.Float64("confidence", result.Confidence),
		)
	} else {
		w.updateVideoStatus(ctx, event.VideoID, "rejected", result.Reason)
		w.logger.Warn("Video rejected",
			zap.String("video_id", event.VideoID),
			zap.Strings("flags", result.Flags),
		)
	}

	// Store moderation result
	if err := w.storeModerationResult(ctx, event.VideoID, result); err != nil {
		w.logger.Error("Failed to store moderation result", zap.Error(err))
	}
}

type Video struct {
	ID          string
	Title       string
	Description string
	Status      string
}

func (w *Worker) getVideo(ctx context.Context, videoID string) (*Video, error) {
	var video Video
	err := w.db.QueryRow(ctx,
		"SELECT id, title, description, status FROM videos WHERE id = $1",
		videoID,
	).Scan(&video.ID, &video.Title, &video.Description, &video.Status)

	if err != nil {
		return nil, err
	}

	return &video, nil
}

func (w *Worker) updateVideoStatus(ctx context.Context, videoID, status, reason string) error {
	_, err := w.db.Exec(ctx,
		"UPDATE videos SET status = $1, updated_at = $2 WHERE id = $3",
		status, time.Now(), videoID,
	)
	return err
}

func (w *Worker) storeModerationResult(ctx context.Context, videoID string, result *moderator.ModerationResult) error {
	resultJSON, err := result.ToJSON()
	if err != nil {
		return err
	}

	_, err = w.db.Exec(ctx,
		`INSERT INTO moderation_results (video_id, approved, confidence, flags, reason, result_data, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		videoID,
		result.Approved,
		result.Confidence,
		result.Flags,
		result.Reason,
		resultJSON,
		time.Now(),
	)

	return err
}