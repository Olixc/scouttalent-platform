package model

import "time"

type VideoStatus string
type UploadStatus string

const (
	VideoStatusPending    VideoStatus = "pending"
	VideoStatusProcessing VideoStatus = "processing"
	VideoStatusReady      VideoStatus = "ready"
	VideoStatusFailed     VideoStatus = "failed"
	VideoStatusModerated  VideoStatus = "moderated"
	VideoStatusRejected   VideoStatus = "rejected"

	UploadStatusInitiated  UploadStatus = "initiated"
	UploadStatusInProgress UploadStatus = "in_progress"
	UploadStatusCompleted  UploadStatus = "completed"
	UploadStatusFailed     UploadStatus = "failed"
)

type Video struct {
	ID           string      `json:"id" db:"id"`
	ProfileID    string      `json:"profile_id" db:"profile_id"`
	Title        string      `json:"title" db:"title"`
	Description  string      `json:"description" db:"description"`
	FileName     string      `json:"file_name" db:"file_name"`
	FileSize     int64       `json:"file_size" db:"file_size"`
	MimeType     string      `json:"mime_type" db:"mime_type"`
	Duration     *int        `json:"duration,omitempty" db:"duration"`
	ThumbnailURL *string     `json:"thumbnail_url,omitempty" db:"thumbnail_url"`
	BlobURL      string      `json:"blob_url" db:"blob_url"`
	StreamURL    string      `json:"stream_url,omitempty" db:"-"`
	Status       VideoStatus `json:"status" db:"status"`
	Visibility   string      `json:"visibility" db:"visibility"`
	ViewCount    int         `json:"view_count" db:"view_count"`
	CreatedAt    time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at" db:"updated_at"`
}

type VideoUpload struct {
	ID          string       `json:"id" db:"id"`
	VideoID     string       `json:"video_id" db:"video_id"`
	Status      UploadStatus `json:"status" db:"status"`
	Progress    int          `json:"progress" db:"progress"`
	ErrorMsg    *string      `json:"error_msg,omitempty" db:"error_msg"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
	CompletedAt *time.Time   `json:"completed_at,omitempty" db:"completed_at"`
}

type VideoUploadRequest struct {
	ProfileID   string `json:"profile_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	FileName    string `json:"file_name" binding:"required"`
	FileSize    int64  `json:"file_size" binding:"required"`
	MimeType    string `json:"mime_type" binding:"required"`
}

type VideoUploadResponse struct {
	VideoID   string    `json:"video_id"`
	UploadID  string    `json:"upload_id"`
	UploadURL string    `json:"upload_url"`
	ExpiresAt time.Time `json:"expires_at"`
	TestMode  bool      `json:"test_mode"`
}

type VideoUpdateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Visibility  string `json:"visibility"`
}

type VideoListResponse struct {
	Videos []*Video `json:"videos"`
	Total  int      `json:"total"`
	Limit  int      `json:"limit"`
	Offset int      `json:"offset"`
}