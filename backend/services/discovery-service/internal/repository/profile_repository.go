package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/scouttalent/discovery-service/internal/model"
)

type ProfileRepository struct {
	db *pgxpool.Pool
}

func NewProfileRepository(db *pgxpool.Pool) *ProfileRepository {
	return &ProfileRepository{db: db}
}

func (r *ProfileRepository) SearchProfiles(ctx context.Context, query string, filters model.ProfileFilters, limit, offset int) ([]model.Profile, int, error) {
	// Build search query
	sqlQuery := `
		SELECT p.id, p.user_id, p.bio, p.location, p.profile_type, p.avatar_url, p.created_at,
		       pd.position, pd.preferred_foot, pd.height, pd.weight
		FROM profiles p
		LEFT JOIN player_details pd ON p.id = pd.profile_id
		WHERE 1=1
	`

	args := []interface{}{}
	argCount := 1

	// Add search conditions
	if query != "" {
		sqlQuery += fmt.Sprintf(" AND (p.bio ILIKE $%d OR p.location ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+query+"%")
		argCount++
	}

	if filters.ProfileType != "" {
		sqlQuery += fmt.Sprintf(" AND p.profile_type = $%d", argCount)
		args = append(args, filters.ProfileType)
		argCount++
	}

	if filters.Position != "" {
		sqlQuery += fmt.Sprintf(" AND pd.position = $%d", argCount)
		args = append(args, filters.Position)
		argCount++
	}

	if filters.Location != "" {
		sqlQuery += fmt.Sprintf(" AND p.location ILIKE $%d", argCount)
		args = append(args, "%"+filters.Location+"%")
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
	sqlQuery += fmt.Sprintf(" ORDER BY p.created_at DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	// Execute query
	rows, err := r.db.Query(ctx, sqlQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	profiles := []model.Profile{}
	for rows.Next() {
		var p model.Profile
		var position, preferredFoot *string
		var height, weight *int

		err := rows.Scan(
			&p.ID, &p.UserID, &p.Bio, &p.Location, &p.ProfileType, &p.AvatarURL, &p.CreatedAt,
			&position, &preferredFoot, &height, &weight,
		)
		if err != nil {
			return nil, 0, err
		}

		if position != nil {
			p.Position = *position
		}
		if preferredFoot != nil {
			p.PreferredFoot = *preferredFoot
		}
		if height != nil {
			p.Height = *height
		}
		if weight != nil {
			p.Weight = *weight
		}

		profiles = append(profiles, p)
	}

	return profiles, total, nil
}

func (r *ProfileRepository) GetSimilarProfiles(ctx context.Context, profileID string, limit int) ([]model.Profile, error) {
	// Get profiles with similar characteristics
	query := `
		SELECT p.id, p.user_id, p.bio, p.location, p.profile_type, p.avatar_url, p.created_at,
		       pd.position, pd.preferred_foot, pd.height, pd.weight
		FROM profiles p
		LEFT JOIN player_details pd ON p.id = pd.profile_id
		WHERE p.id != $1
		  AND p.profile_type = (SELECT profile_type FROM profiles WHERE id = $1)
		ORDER BY p.created_at DESC
		LIMIT $2
	`

	rows, err := r.db.Query(ctx, query, profileID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	profiles := []model.Profile{}
	for rows.Next() {
		var p model.Profile
		var position, preferredFoot *string
		var height, weight *int

		err := rows.Scan(
			&p.ID, &p.UserID, &p.Bio, &p.Location, &p.ProfileType, &p.AvatarURL, &p.CreatedAt,
			&position, &preferredFoot, &height, &weight,
		)
		if err != nil {
			return nil, err
		}

		if position != nil {
			p.Position = *position
		}
		if preferredFoot != nil {
			p.PreferredFoot = *preferredFoot
		}
		if height != nil {
			p.Height = *height
		}
		if weight != nil {
			p.Weight = *weight
		}

		profiles = append(profiles, p)
	}

	return profiles, nil
}