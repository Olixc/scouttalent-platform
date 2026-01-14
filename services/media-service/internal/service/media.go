package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/scouttalent/services/media-service/internal/model"
	"github.com/scouttalent/services/media-service/internal/repository"
	"github.com/scouttalent/services/media-service/internal/storage"
)

type MediaService struct {
	repo    *repository.MediaRepository
	storage *storage.BlobStorage
}

func NewMediaService(repo *repository.MediaRepository, storage *storage.BlobStorage) *MediaService {
	return &MediaService{
		repo:    repo,
		storage: storage,
	}
}

// InitiateUpload creates a new video record and returns upload URL
func (s *MediaService) InitiateUpload(ctx context.Context, req *model.VideoUploadRequest) (*model.VideoUploadResponse, error) {
	// Create video record
	video := &model.Video{
		ID:          uuid.New().String(),
		ProfileID:   req.ProfileID,
		Title:       req.Title,
		Description: req.Description,
		FileName:    req.FileName,
		FileSize:    req.FileSize,
		MimeType:    req.MimeType,
		Status:      model.VideoStatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.CreateVideo(ctx, video); err != nil {
		return nil, fmt.Errorf("failed to create video: %w", err)
	}

	// Create upload record
	upload := &model.VideoUpload{
		ID:        uuid.New().String(),
		VideoID:   video.ID,
		Status:    model.UploadStatusInitiated,
		Progress:  0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.CreateUpload(ctx, upload); err != nil {
		return nil, fmt.Errorf("failed to create upload: %w", err)
	}

	// Generate upload URL
	uploadURL, err := s.storage.GenerateUploadURL(ctx, video.ID, req.FileName)
	if err != nil {
		return nil, fmt.Errorf("failed to generate upload URL: %w", err)
	}

	// Check if we're in test mode
	testMode := s.storage.IsTestMode()

	return &model.VideoUploadResponse{
		VideoID:   video.ID,
		UploadID:  upload.ID,
		UploadURL: uploadURL,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		TestMode:  testMode,
	}, nil
}

// UpdateUploadProgress updates the upload progress
func (s *MediaService) UpdateUploadProgress(ctx context.Context, uploadID string, progress int) error {
	upload, err := s.repo.GetUploadByID(ctx, uploadID)
	if err != nil {
		return fmt.Errorf("failed to get upload: %w", err)
	}

	upload.Progress = progress
	upload.UpdatedAt = time.Now()

	if progress >= 100 {
		upload.Status = model.UploadStatusCompleted
		upload.CompletedAt = &upload.UpdatedAt
	} else {
		upload.Status = model.UploadStatusInProgress
	}

	return s.repo.UpdateUpload(ctx, upload)
}

// CompleteUpload marks the upload as complete and updates video status
func (s *MediaService) CompleteUpload(ctx context.Context, videoID string) error {
	video, err := s.repo.GetVideoByID(ctx, videoID)
	if err != nil {
		return fmt.Errorf("failed to get video: %w", err)
	}

	// Generate blob URL
	video.BlobURL = s.storage.GetBlobURL(videoID, video.FileName)
	video.Status = model.VideoStatusProcessing
	video.UpdatedAt = time.Now()

	if err := s.repo.UpdateVideo(ctx, video); err != nil {
		return fmt.Errorf("failed to update video: %w", err)
	}

	// In test mode, immediately mark as ready since we're not actually processing
	if s.storage.IsTestMode() {
		video.Status = model.VideoStatusReady
		video.UpdatedAt = time.Now()
		if err := s.repo.UpdateVideo(ctx, video); err != nil {
			return fmt.Errorf("failed to update video status: %w", err)
		}
	}

	// TODO: Publish event to NATS for AI moderation
	// For now, we'll skip this in test mode

	return nil
}

// GetVideo retrieves a video by ID
func (s *MediaService) GetVideo(ctx context.Context, videoID string) (*model.Video, error) {
	video, err := s.repo.GetVideoByID(ctx, videoID)
	if err != nil {
		return nil, fmt.Errorf("failed to get video: %w", err)
	}

	// Generate download URL if video is ready
	if video.Status == model.VideoStatusReady && video.BlobURL != "" {
		downloadURL, err := s.storage.GenerateDownloadURL(ctx, video.ID, video.FileName)
		if err == nil {
			video.StreamURL = downloadURL
		}
	}

	return video, nil
}

// ListProfileVideos retrieves all videos for a profile
func (s *MediaService) ListProfileVideos(ctx context.Context, profileID string, limit, offset int) ([]*model.Video, int, error) {
	videos, total, err := s.repo.ListVideosByProfile(ctx, profileID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list videos: %w", err)
	}

	// Generate download URLs for ready videos
	for _, video := range videos {
		if video.Status == model.VideoStatusReady && video.BlobURL != "" {
			downloadURL, err := s.storage.GenerateDownloadURL(ctx, video.ID, video.FileName)
			if err == nil {
				video.StreamURL = downloadURL
			}
		}
	}

	return videos, total, nil
}

// UpdateVideo updates video metadata
func (s *MediaService) UpdateVideo(ctx context.Context, videoID string, req *model.VideoUpdateRequest) error {
	video, err := s.repo.GetVideoByID(ctx, videoID)
	if err != nil {
		return fmt.Errorf("failed to get video: %w", err)
	}

	if req.Title != "" {
		video.Title = req.Title
	}
	if req.Description != "" {
		video.Description = req.Description
	}
	if req.Visibility != "" {
		video.Visibility = req.Visibility
	}

	video.UpdatedAt = time.Now()

	return s.repo.UpdateVideo(ctx, video)
}

// DeleteVideo deletes a video and its blob
func (s *MediaService) DeleteVideo(ctx context.Context, videoID string) error {
	video, err := s.repo.GetVideoByID(ctx, videoID)
	if err != nil {
		return fmt.Errorf("failed to get video: %w", err)
	}

	// Delete from blob storage (will be no-op in test mode)
	if err := s.storage.DeleteVideo(ctx, video.ID, video.FileName); err != nil {
		return fmt.Errorf("failed to delete blob: %w", err)
	}

	// Delete from database
	if err := s.repo.DeleteVideo(ctx, videoID); err != nil {
		return fmt.Errorf("failed to delete video: %w", err)
	}

	return nil
}