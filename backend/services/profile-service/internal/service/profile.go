package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/scouttalent/profile-service/internal/model"
	"github.com/scouttalent/profile-service/internal/repository"
)

type ProfileService struct {
	repo *repository.ProfileRepository
}

func NewProfileService(repo *repository.ProfileRepository) *ProfileService {
	return &ProfileService{repo: repo}
}

func (s *ProfileService) CreateProfile(ctx context.Context, userID string, userType model.UserType, req model.CreateProfileRequest) (*model.Profile, error) {
	// Check if profile already exists
	existing, err := s.repo.GetByUserID(ctx, userID)
	if err == nil && existing != nil {
		return nil, repository.ErrProfileAlreadyExists
	}

	profile := &model.Profile{
		ID:                     uuid.New().String(),
		UserID:                 userID,
		Type:                   userType,
		DisplayName:            req.DisplayName,
		Bio:                    req.Bio,
		LocationCountry:        req.LocationCountry,
		LocationCity:           req.LocationCity,
		TrustLevel:             model.TrustLevelNewcomer,
		ProfileCompletionScore: 0,
		CreatedAt:              time.Now(),
		UpdatedAt:              time.Now(),
	}

	if err := s.repo.Create(ctx, profile); err != nil {
		return nil, err
	}

	// Calculate initial completion score
	profile.ProfileCompletionScore = s.repo.CalculateCompletionScore(ctx, profile)
	if err := s.repo.Update(ctx, profile); err != nil {
		return nil, err
	}

	return profile, nil
}

func (s *ProfileService) GetProfile(ctx context.Context, profileID string) (*model.Profile, error) {
	return s.repo.GetByID(ctx, profileID)
}

func (s *ProfileService) GetProfileByUserID(ctx context.Context, userID string) (*model.Profile, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *ProfileService) UpdateProfile(ctx context.Context, profileID string, req model.UpdateProfileRequest) (*model.Profile, error) {
	profile, err := s.repo.GetByID(ctx, profileID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.DisplayName != nil {
		profile.DisplayName = *req.DisplayName
	}
	if req.Bio != nil {
		profile.Bio = req.Bio
	}
	if req.AvatarURL != nil {
		profile.AvatarURL = req.AvatarURL
	}
	if req.LocationCountry != nil {
		profile.LocationCountry = req.LocationCountry
	}
	if req.LocationCity != nil {
		profile.LocationCity = req.LocationCity
	}

	profile.UpdatedAt = time.Now()

	// Recalculate completion score
	profile.ProfileCompletionScore = s.repo.CalculateCompletionScore(ctx, profile)

	if err := s.repo.Update(ctx, profile); err != nil {
		return nil, err
	}

	return profile, nil
}

func (s *ProfileService) CreatePlayerDetails(ctx context.Context, profileID string, req model.CreatePlayerDetailsRequest) (*model.PlayerProfile, error) {
	// Verify profile exists and is a player
	profile, err := s.repo.GetByID(ctx, profileID)
	if err != nil {
		return nil, err
	}

	if profile.Type != model.UserTypePlayer {
		return nil, fmt.Errorf("profile is not a player")
	}

	details := &model.PlayerDetails{
		ProfileID:     profileID,
		Position:      req.Position,
		DateOfBirth:   req.DateOfBirth,
		HeightCM:      req.HeightCM,
		WeightKG:      req.WeightKG,
		PreferredFoot: req.PreferredFoot,
		CurrentTeam:   req.CurrentTeam,
	}

	if err := s.repo.CreatePlayerDetails(ctx, details); err != nil {
		return nil, err
	}

	// Recalculate completion score
	profile.ProfileCompletionScore = s.repo.CalculateCompletionScore(ctx, profile)
	profile.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, profile); err != nil {
		return nil, err
	}

	return s.repo.GetPlayerProfile(ctx, profileID)
}

func (s *ProfileService) GetPlayerProfile(ctx context.Context, profileID string) (*model.PlayerProfile, error) {
	return s.repo.GetPlayerProfile(ctx, profileID)
}