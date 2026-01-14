package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/scouttalent/discovery-service/internal/model"
)

type VideoRepository struct {
	db *pgxpool.Pool
}

func NewVideoRepository(db *pgxpool.Pool) *VideoRepository {
	return &VideoRepository{db: db}
}

func (r *VideoRepository) SearchVideos(ctx context.Context, query string, limit, offset int) ([]model.Video, int, error) {
	sqlQuery := `
		SELECT id, profile_id, title, description, file_name, blob_url, status, view_count, created_at
		FROM videos
		WHERE status = 'ready'
	`

	args := []interface{}{}
	argCount := 1

	if query != "" {
		sqlQuery += fmt.Sprintf(" AND (title ILIKE $%d OR description ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+query+"%")
		argCount++
	}

	// Get total count
	countQuery := "SELECT COUNT(*) FROM (" + sqlQuery + ") AS count_query"
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Add pagination
	sqlQuery += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, sqlQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	videos := []model.Video{}
	for rows.Next() {
		var v model.Video
		err := rows.Scan(&v.ID, &v.ProfileID, &v.Title, &v.Description, &v.FileName, &v.BlobURL, &v.Status, &v.ViewCount, &v.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		videos = append(videos, v)
	}

	return videos, total, nil
}

func (r *VideoRepository) GetFeed(ctx context.Context, limit, offset int) ([]model.Video, int, error) {
	query := `
		SELECT id, profile_id, title, description, file_name, blob_url, status, view_count, created_at
		FROM videos
		WHERE status = 'ready'
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	videos := []model.Video{}
	for rows.Next() {
		var v model.Video
		err := rows.Scan(&v.ID, &v.ProfileID, &v.Title, &v.Description, &v.FileName, &v.BlobURL, &v.Status, &v.ViewCount, &v.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		videos = append(videos, v)
	}

	// Get total count
	var total int
	err = r.db.QueryRow(ctx, "SELECT COUNT(*) FROM videos WHERE status = 'ready'").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return videos, total, nil
}

func (r *VideoRepository) GetTrendingVideos(ctx context.Context, limit int) ([]model.Video, error) {
	query := `
		SELECT id, profile_id, title, description, file_name, blob_url, status, view_count, created_at
		FROM videos
		WHERE status = 'ready'
		ORDER BY view_count DESC, created_at DESC
		LIMIT $1
	`

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	videos := []model.Video{}
	for rows.Next() {
		var v model.Video
		err := rows.Scan(&v.ID, &v.ProfileID, &v.Title, &v.Description, &v.FileName, &v.BlobURL, &v.Status, &v.ViewCount, &v.CreatedAt)
		if err != nil {
			return nil, err
		}
		videos = append(videos, v)
	}

	return videos, nil
}

func (r *VideoRepository) GetRecommendedVideos(ctx context.Context, profileID string, limit int) ([]model.Video, error) {
	// Get videos from similar profiles
	query := `
		SELECT v.id, v.profile_id, v.title, v.description, v.file_name, v.blob_url, v.status, v.view_count, v.created_at
		FROM videos v
		JOIN profiles p ON v.profile_id = p.id
		WHERE v.status = 'ready'
		  AND v.profile_id != $1
		  AND p.profile_type = (SELECT profile_type FROM profiles WHERE id = $1)
		ORDER BY v.view_count DESC, v.created_at DESC
		LIMIT $2
	`

	rows, err := r.db.Query(ctx, query, profileID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	videos := []model.Video{}
	for rows.Next() {
		var v model.Video
		err := rows.Scan(&v.ID, &v.ProfileID, &v.Title, &v.Description, &v.FileName, &v.BlobURL, &v.Status, &v.ViewCount, &v.CreatedAt)
		if err != nil {
			return nil, err
		}
		videos = append(videos, v)
	}

	return videos, nil
}