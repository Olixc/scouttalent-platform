package model

import (
	"time"

	"github.com/google/uuid"
)

type VideoStatus string

const (
	VideoStatusUploading  VideoStatus = "uploading"
	VideoStatusProcessing VideoStatus = "processing"
	VideoStatusReady      VideoStatus = "ready"
	VideoStatusFailed     VideoStatus = "failed"
)

type Video struct {
	ID           uuid.UUID   `json:"id"`
	ProfileID    uuid.UUID   `json:"profile_id"`
	Title        string      `json:"title"`
	Description  string      `json:"description,omitempty"`
	BlobURL      string      `json:"blob_url"`
	ThumbnailURL string      `json:"thumbnail_url,omitempty"`
	Duration     int         `json:"duration"` // in seconds
	FileSize     int64       `json:"file_size"`
	MimeType     string      `json:"mime_type"`
	Status       VideoStatus `json:"status"`
	Metadata     VideoMetadata `json:"metadata,omitempty"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

type VideoMetadata struct {
	Width       int    `json:"width,omitempty"`
	Height      int    `json:"height,omitempty"`
	Codec       string `json:"codec,omitempty"`
	Bitrate     int    `json:"bitrate,omitempty"`
	FrameRate   float64 `json:"frame_rate,omitempty"`
}

type Upload struct {
	ID        uuid.UUID    `json:"id"`
	VideoID   uuid.UUID    `json:"video_id"`
	UploadID  string       `json:"upload_id"` // TUS upload ID
	Status    VideoStatus  `json:"status"`
	Progress  int          `json:"progress"` // 0-100
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type CreateVideoRequest struct {
	Title       string `json:"title" binding:"required,min=3,max=200"`
	Description string `json:"description" binding:"max=2000"`
	FileName    string `json:"file_name" binding:"required"`
	FileSize    int64  `json:"file_size" binding:"required,gt=0"`
	MimeType    string `json:"mime_type" binding:"required"`
}

type UpdateVideoRequest struct {
	Title       *string `json:"title" binding:"omitempty,min=3,max=200"`
	Description *string `json:"description" binding:"omitempty,max=2000"`
}

type VideoResponse struct {
	Video   *Video `json:"video"`
	Message string `json:"message,omitempty"`
}

type VideoListResponse struct {
	Videos []Video `json:"videos"`
	Total  int     `json:"total"`
}