package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/scouttalent/profile-service/internal/model"
)

var (
	ErrProfileNotFound      = errors.New("profile not found")
	ErrProfileAlreadyExists = errors.New("profile already exists")
)

type ProfileRepository struct {
	pool *pgxpool.Pool
}

func NewProfileRepository(pool *pgxpool.Pool) *ProfileRepository {
	return &ProfileRepository{pool: pool}
}

func (r *ProfileRepository) Create(ctx context.Context, profile *model.Profile) error {
	query := `
		INSERT INTO profiles (id, user_id, type, display_name, bio, avatar_url, 
		                     location_country, location_city, trust_level, 
		                     profile_completion_score, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := r.pool.Exec(ctx, query,
		profile.ID,
		profile.UserID,
		profile.Type,
		profile.DisplayName,
		profile.Bio,
		profile.AvatarURL,
		profile.LocationCountry,
		profile.LocationCity,
		profile.TrustLevel,
		profile.ProfileCompletionScore,
		profile.CreatedAt,
		profile.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create profile: %w", err)
	}

	return nil
}

func (r *ProfileRepository) GetByID(ctx context.Context, id string) (*model.Profile, error) {
	query := `
		SELECT id, user_id, type, display_name, bio, avatar_url, 
		       location_country, location_city, trust_level, 
		       profile_completion_score, created_at, updated_at
		FROM profiles
		WHERE id = $1
	`

	var profile model.Profile
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.Type,
		&profile.DisplayName,
		&profile.Bio,
		&profile.AvatarURL,
		&profile.LocationCountry,
		&profile.LocationCity,
		&profile.TrustLevel,
		&profile.ProfileCompletionScore,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProfileNotFound
		}
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	return &profile, nil
}

func (r *ProfileRepository) GetByUserID(ctx context.Context, userID string) (*model.Profile, error) {
	query := `
		SELECT id, user_id, type, display_name, bio, avatar_url, 
		       location_country, location_city, trust_level, 
		       profile_completion_score, created_at, updated_at
		FROM profiles
		WHERE user_id = $1
	`

	var profile model.Profile
	err := r.pool.QueryRow(ctx, query, userID).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.Type,
		&profile.DisplayName,
		&profile.Bio,
		&profile.AvatarURL,
		&profile.LocationCountry,
		&profile.LocationCity,
		&profile.TrustLevel,
		&profile.ProfileCompletionScore,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProfileNotFound
		}
		return nil, fmt.Errorf("failed to get profile by user ID: %w", err)
	}

	return &profile, nil
}

func (r *ProfileRepository) Update(ctx context.Context, profile *model.Profile) error {
	query := `
		UPDATE profiles
		SET display_name = $1, bio = $2, avatar_url = $3, 
		    location_country = $4, location_city = $5, 
		    profile_completion_score = $6, updated_at = $7
		WHERE id = $8
	`

	_, err := r.pool.Exec(ctx, query,
		profile.DisplayName,
		profile.Bio,
		profile.AvatarURL,
		profile.LocationCountry,
		profile.LocationCity,
		profile.ProfileCompletionScore,
		profile.UpdatedAt,
		profile.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	return nil
}

func (r *ProfileRepository) CreatePlayerDetails(ctx context.Context, details *model.PlayerDetails) error {
	query := `
		INSERT INTO player_details (profile_id, position, date_of_birth, height_cm, 
		                           weight_kg, preferred_foot, current_team)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.pool.Exec(ctx, query,
		details.ProfileID,
		details.Position,
		details.DateOfBirth,
		details.HeightCM,
		details.WeightKG,
		details.PreferredFoot,
		details.CurrentTeam,
	)

	if err != nil {
		return fmt.Errorf("failed to create player details: %w", err)
	}

	return nil
}

func (r *ProfileRepository) GetPlayerProfile(ctx context.Context, profileID string) (*model.PlayerProfile, error) {
	query := `
		SELECT 
			p.id, p.user_id, p.type, p.display_name, p.bio, p.avatar_url,
			p.location_country, p.location_city, p.trust_level,
			p.profile_completion_score, p.created_at, p.updated_at,
			pd.profile_id, pd.position, pd.date_of_birth, pd.height_cm,
			pd.weight_kg, pd.preferred_foot, pd.current_team, pd.skill_scores,
			pd.overall_score, pd.last_scored_at
		FROM profiles p
		JOIN player_details pd ON pd.profile_id = p.id
		WHERE p.id = $1 AND p.type = 'player'
	`

	var player model.PlayerProfile
	err := r.pool.QueryRow(ctx, query, profileID).Scan(
		&player.ID,
		&player.UserID,
		&player.Type,
		&player.DisplayName,
		&player.Bio,
		&player.AvatarURL,
		&player.LocationCountry,
		&player.LocationCity,
		&player.TrustLevel,
		&player.ProfileCompletionScore,
		&player.CreatedAt,
		&player.UpdatedAt,
		&player.PlayerDetails.ProfileID,
		&player.PlayerDetails.Position,
		&player.PlayerDetails.DateOfBirth,
		&player.PlayerDetails.HeightCM,
		&player.PlayerDetails.WeightKG,
		&player.PlayerDetails.PreferredFoot,
		&player.PlayerDetails.CurrentTeam,
		&player.PlayerDetails.SkillScores,
		&player.PlayerDetails.OverallScore,
		&player.PlayerDetails.LastScoredAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProfileNotFound
		}
		return nil, fmt.Errorf("failed to get player profile: %w", err)
	}

	return &player, nil
}

func (r *ProfileRepository) CalculateCompletionScore(ctx context.Context, profile *model.Profile) int {
	score := 0

	// Basic fields (60 points)
	if profile.DisplayName != "" {
		score += 10
	}
	if profile.Bio != nil && *profile.Bio != "" {
		score += 15
	}
	if profile.AvatarURL != nil && *profile.AvatarURL != "" {
		score += 15
	}
	if profile.LocationCountry != nil && *profile.LocationCountry != "" {
		score += 10
	}
	if profile.LocationCity != nil && *profile.LocationCity != "" {
		score += 10
	}

	// Type-specific fields (40 points)
	if profile.Type == model.UserTypePlayer {
		// Check player details
		var hasDetails bool
		query := `SELECT EXISTS(SELECT 1 FROM player_details WHERE profile_id = $1)`
		r.pool.QueryRow(ctx, query, profile.ID).Scan(&hasDetails)
		if hasDetails {
			score += 40
		}
	}

	return score
}