package model

import (
	"time"
)

type UserType string
type TrustLevel string

const (
	UserTypePlayer  UserType = "player"
	UserTypeScout   UserType = "scout"
	UserTypeAcademy UserType = "academy"

	TrustLevelNewcomer    TrustLevel = "newcomer"
	TrustLevelEstablished TrustLevel = "established"
	TrustLevelVerified    TrustLevel = "verified"
	TrustLevelPro         TrustLevel = "pro"
)

type Profile struct {
	ID                     string     `json:"id" db:"id"`
	UserID                 string     `json:"user_id" db:"user_id"`
	Type                   UserType   `json:"type" db:"type"`
	DisplayName            string     `json:"display_name" db:"display_name"`
	Bio                    *string    `json:"bio,omitempty" db:"bio"`
	AvatarURL              *string    `json:"avatar_url,omitempty" db:"avatar_url"`
	LocationCountry        *string    `json:"location_country,omitempty" db:"location_country"`
	LocationCity           *string    `json:"location_city,omitempty" db:"location_city"`
	TrustLevel             TrustLevel `json:"trust_level" db:"trust_level"`
	ProfileCompletionScore int        `json:"profile_completion_score" db:"profile_completion_score"`
	CreatedAt              time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at" db:"updated_at"`
}

type PlayerDetails struct {
	ProfileID      string     `json:"profile_id" db:"profile_id"`
	Position       string     `json:"position" db:"position"`
	DateOfBirth    *time.Time `json:"date_of_birth,omitempty" db:"date_of_birth"`
	HeightCM       *int       `json:"height_cm,omitempty" db:"height_cm"`
	WeightKG       *int       `json:"weight_kg,omitempty" db:"weight_kg"`
	PreferredFoot  *string    `json:"preferred_foot,omitempty" db:"preferred_foot"`
	CurrentTeam    *string    `json:"current_team,omitempty" db:"current_team"`
	SkillScores    *string    `json:"skill_scores,omitempty" db:"skill_scores"` // JSONB
	OverallScore   *float64   `json:"overall_score,omitempty" db:"overall_score"`
	LastScoredAt   *time.Time `json:"last_scored_at,omitempty" db:"last_scored_at"`
}

type ScoutDetails struct {
	ProfileID            string     `json:"profile_id" db:"profile_id"`
	Organization         *string    `json:"organization,omitempty" db:"organization"`
	OrganizationType     *string    `json:"organization_type,omitempty" db:"organization_type"`
	RegionsOfInterest    []string   `json:"regions_of_interest,omitempty" db:"regions_of_interest"`
	PositionsOfInterest  []string   `json:"positions_of_interest,omitempty" db:"positions_of_interest"`
	VerifiedAt           *time.Time `json:"verified_at,omitempty" db:"verified_at"`
	VerifiedBy           *string    `json:"verified_by,omitempty" db:"verified_by"`
	VerificationDocuments *string   `json:"verification_documents,omitempty" db:"verification_documents"` // JSONB
}

type CreateProfileRequest struct {
	DisplayName     string  `json:"display_name" binding:"required,min=2,max=100"`
	Bio             *string `json:"bio" binding:"omitempty,max=500"`
	LocationCountry *string `json:"location_country" binding:"omitempty,max=100"`
	LocationCity    *string `json:"location_city" binding:"omitempty,max=100"`
}

type UpdateProfileRequest struct {
	DisplayName     *string `json:"display_name" binding:"omitempty,min=2,max=100"`
	Bio             *string `json:"bio" binding:"omitempty,max=500"`
	AvatarURL       *string `json:"avatar_url" binding:"omitempty,url"`
	LocationCountry *string `json:"location_country" binding:"omitempty,max=100"`
	LocationCity    *string `json:"location_city" binding:"omitempty,max=100"`
}

type CreatePlayerDetailsRequest struct {
	Position      string     `json:"position" binding:"required,oneof=goalkeeper defender midfielder forward"`
	DateOfBirth   *time.Time `json:"date_of_birth" binding:"omitempty"`
	HeightCM      *int       `json:"height_cm" binding:"omitempty,min=100,max=250"`
	WeightKG      *int       `json:"weight_kg" binding:"omitempty,min=30,max=150"`
	PreferredFoot *string    `json:"preferred_foot" binding:"omitempty,oneof=left right both"`
	CurrentTeam   *string    `json:"current_team" binding:"omitempty,max=100"`
}

type PlayerProfile struct {
	Profile
	PlayerDetails PlayerDetails `json:"player_details"`
}

type ScoutProfile struct {
	Profile
	ScoutDetails ScoutDetails `json:"scout_details"`
}