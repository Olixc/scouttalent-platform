package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/scouttalent/media-service/internal/model"
	"github.com/scouttalent/media-service/internal/repository"
	"github.com/scouttalent/media-service/internal/storage"
	"go.uber.org/zap"
)

type MediaService struct {
	repo       *repository.MediaRepository
	blobClient *storage.AzureBlobClient
	nats       *nats.Conn
	logger     *zap.Logger
}

func NewMediaService(repo *repository.MediaRepository, blobClient *storage.AzureBlobClient, nc *nats.Conn, logger *zap.Logger) *MediaService {
	return &MediaService{
		repo:       repo,
		blobClient: blobClient,
		nats:       nc,
		logger:     logger,
	}
}

func (s *MediaService) InitiateUpload(ctx context.Context, profileID uuid.UUID, req *model.CreateVideoRequest) (*model.Video, *model.Upload, error) {
	// Create video record
	video := &model.Video{
		ID:          uuid.New(),
		ProfileID:   profileID,
		Title:       req.Title,
		Description: req.Description,
		FileSize:    req.FileSize,
		MimeType:    req.MimeType,
		Status:      model.VideoStatusUploading,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Generate blob URL
	blobName := fmt.Sprintf("%s/%s", profileID.String(), video.ID.String())
	video.BlobURL = s.blobClient.GetBlobURL(blobName)

	// Save video to database
	if err := s.repo.CreateVideo(ctx, video); err != nil {
		return nil, nil, fmt.Errorf("failed to create video: %w", err)
	}

	// Create upload record
	upload := &model.Upload{
		ID:        uuid.New(),
		VideoID:   video.ID,
		UploadID:  uuid.New().String(), // TUS upload ID
		Status:    model.VideoStatusUploading,
		Progress:  0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.CreateUpload(ctx, upload); err != nil {
		return nil, nil, fmt.Errorf("failed to create upload: %w", err)
	}

	s.logger.Info("upload initiated",
		zap.String("video_id", video.ID.String()),
		zap.String("upload_id", upload.ID.String()),
	)

	return video, upload, nil
}

func (s *MediaService) UpdateUploadProgress(ctx context.Context, uploadID uuid.UUID, progress int) error {
	status := model.VideoStatusUploading
	if progress >= 100 {
		status = model.VideoStatusProcessing
	}

	if err := s.repo.UpdateUploadProgress(ctx, uploadID, progress, status); err != nil {
		return fmt.Errorf("failed to update upload progress: %w", err)
	}

	return nil
}

func (s *MediaService) CompleteUpload(ctx context.Context, videoID uuid.UUID) error {
	video, err := s.repo.GetVideoByID(ctx, videoID)
	if err != nil {
		return fmt.Errorf("failed to get video: %w", err)
	}

	// Update video status
	video.Status = model.VideoStatusReady
	video.UpdatedAt = time.Now()

	if err := s.repo.UpdateVideo(ctx, video); err != nil {
		return fmt.Errorf("failed to update video: %w", err)
	}

	// Publish video uploaded event to NATS
	if err := s.publishVideoUploadedEvent(video); err != nil {
		s.logger.Error("failed to publish video uploaded event", zap.Error(err))
	}

	s.logger.Info("upload completed", zap.String("video_id", videoID.String()))

	return nil
}

func (s *MediaService) GetVideo(ctx context.Context, videoID uuid.UUID) (*model.Video, error) {
	return s.repo.GetVideoByID(ctx, videoID)
}

func (s *MediaService) ListProfileVideos(ctx context.Context, profileID uuid.UUID, limit, offset int) ([]model.Video, int, error) {
	videos, err := s.repo.GetVideosByProfileID(ctx, profileID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get videos: %w", err)
	}

	total, err := s.repo.CountVideosByProfileID(ctx, profileID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count videos: %w", err)
	}

	return videos, total, nil
}

func (s *MediaService) UpdateVideo(ctx context.Context, videoID uuid.UUID, req *model.UpdateVideoRequest) (*model.Video, error) {
	video, err := s.repo.GetVideoByID(ctx, videoID)
	if err != nil {
		return nil, fmt.Errorf("failed to get video: %w", err)
	}

	// Update fields
	if req.Title != nil {
		video.Title = *req.Title
	}
	if req.Description != nil {
		video.Description = *req.Description
	}
	video.UpdatedAt = time.Now()

	if err := s.repo.UpdateVideo(ctx, video); err != nil {
		return nil, fmt.Errorf("failed to update video: %w", err)
	}

	return video, nil
}

func (s *MediaService) DeleteVideo(ctx context.Context, videoID uuid.UUID) error {
	video, err := s.repo.GetVideoByID(ctx, videoID)
	if err != nil {
		return fmt.Errorf("failed to get video: %w", err)
	}

	// Delete from blob storage
	blobName := fmt.Sprintf("%s/%s", video.ProfileID.String(), video.ID.String())
	if err := s.blobClient.DeleteBlob(ctx, blobName); err != nil {
		s.logger.Error("failed to delete blob", zap.Error(err))
	}

	// Delete from database
	if err := s.repo.DeleteVideo(ctx, videoID); err != nil {
		return fmt.Errorf("failed to delete video: %w", err)
	}

	s.logger.Info("video deleted", zap.String("video_id", videoID.String()))

	return nil
}

func (s *MediaService) publishVideoUploadedEvent(video *model.Video) error {
	event := map[string]interface{}{
		"event_type": "video.uploaded",
		"video_id":   video.ID.String(),
		"profile_id": video.ProfileID.String(),
		"title":      video.Title,
		"timestamp":  time.Now().Unix(),
	}

	// Publish to NATS subject
	subject := "media.video.uploaded"
	if err := s.nats.Publish(subject, []byte(fmt.Sprintf("%v", event))); err != nil {
		return err
	}

	return nil
}