package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/scouttalent/media-service/internal/model"
)

type MediaRepository struct {
	pool *pgxpool.Pool
}

func NewMediaRepository(pool *pgxpool.Pool) *MediaRepository {
	return &MediaRepository{pool: pool}
}

func (r *MediaRepository) CreateVideo(ctx context.Context, video *model.Video) error {
	query := `
		INSERT INTO videos (id, profile_id, title, description, blob_url, thumbnail_url, 
			duration, file_size, mime_type, status, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := r.pool.Exec(ctx, query,
		video.ID,
		video.ProfileID,
		video.Title,
		video.Description,
		video.BlobURL,
		video.ThumbnailURL,
		video.Duration,
		video.FileSize,
		video.MimeType,
		video.Status,
		video.Metadata,
		video.CreatedAt,
		video.UpdatedAt,
	)

	return err
}

func (r *MediaRepository) GetVideoByID(ctx context.Context, id uuid.UUID) (*model.Video, error) {
	query := `
		SELECT id, profile_id, title, description, blob_url, thumbnail_url, 
			duration, file_size, mime_type, status, metadata, created_at, updated_at
		FROM videos
		WHERE id = $1
	`

	var video model.Video
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&video.ID,
		&video.ProfileID,
		&video.Title,
		&video.Description,
		&video.BlobURL,
		&video.ThumbnailURL,
		&video.Duration,
		&video.FileSize,
		&video.MimeType,
		&video.Status,
		&video.Metadata,
		&video.CreatedAt,
		&video.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &video, nil
}

func (r *MediaRepository) GetVideosByProfileID(ctx context.Context, profileID uuid.UUID, limit, offset int) ([]model.Video, error) {
	query := `
		SELECT id, profile_id, title, description, blob_url, thumbnail_url, 
			duration, file_size, mime_type, status, metadata, created_at, updated_at
		FROM videos
		WHERE profile_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, profileID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []model.Video
	for rows.Next() {
		var video model.Video
		err := rows.Scan(
			&video.ID,
			&video.ProfileID,
			&video.Title,
			&video.Description,
			&video.BlobURL,
			&video.ThumbnailURL,
			&video.Duration,
			&video.FileSize,
			&video.MimeType,
			&video.Status,
			&video.Metadata,
			&video.CreatedAt,
			&video.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}

	return videos, rows.Err()
}

func (r *MediaRepository) CountVideosByProfileID(ctx context.Context, profileID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM videos WHERE profile_id = $1`

	var count int
	err := r.pool.QueryRow(ctx, query, profileID).Scan(&count)
	return count, err
}

func (r *MediaRepository) UpdateVideo(ctx context.Context, video *model.Video) error {
	query := `
		UPDATE videos
		SET title = $2, description = $3, blob_url = $4, thumbnail_url = $5,
			duration = $6, status = $7, metadata = $8, updated_at = $9
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query,
		video.ID,
		video.Title,
		video.Description,
		video.BlobURL,
		video.ThumbnailURL,
		video.Duration,
		video.Status,
		video.Metadata,
		video.UpdatedAt,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("video not found")
	}

	return nil
}

func (r *MediaRepository) DeleteVideo(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM videos WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("video not found")
	}

	return nil
}

func (r *MediaRepository) CreateUpload(ctx context.Context, upload *model.Upload) error {
	query := `
		INSERT INTO uploads (id, video_id, upload_id, status, progress, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.pool.Exec(ctx, query,
		upload.ID,
		upload.VideoID,
		upload.UploadID,
		upload.Status,
		upload.Progress,
		upload.CreatedAt,
		upload.UpdatedAt,
	)

	return err
}

func (r *MediaRepository) GetUploadByID(ctx context.Context, id uuid.UUID) (*model.Upload, error) {
	query := `
		SELECT id, video_id, upload_id, status, progress, created_at, updated_at
		FROM uploads
		WHERE id = $1
	`

	var upload model.Upload
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&upload.ID,
		&upload.VideoID,
		&upload.UploadID,
		&upload.Status,
		&upload.Progress,
		&upload.CreatedAt,
		&upload.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &upload, nil
}

func (r *MediaRepository) UpdateUploadProgress(ctx context.Context, id uuid.UUID, progress int, status model.VideoStatus) error {
	query := `
		UPDATE uploads
		SET progress = $2, status = $3, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query, id, progress, status)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("upload not found")
	}

	return nil
}