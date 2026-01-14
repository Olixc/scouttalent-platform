package service

import (
	"context"

	"github.com/scouttalent/discovery-service/internal/model"
	"github.com/scouttalent/discovery-service/internal/repository"
	"go.uber.org/zap"
)

type RecommendationService struct {
	profileRepo *repository.ProfileRepository
	videoRepo   *repository.VideoRepository
	logger      *zap.Logger
}

func NewRecommendationService(profileRepo *repository.ProfileRepository, videoRepo *repository.VideoRepository, logger *zap.Logger) *RecommendationService {
	return &RecommendationService{
		profileRepo: profileRepo,
		videoRepo:   videoRepo,
		logger:      logger,
	}
}

func (s *RecommendationService) GetProfileRecommendations(ctx context.Context, profileID string, limit int) ([]model.Profile, error) {
	return s.profileRepo.GetSimilarProfiles(ctx, profileID, limit)
}

func (s *RecommendationService) GetVideoRecommendations(ctx context.Context, profileID string, limit int) ([]model.Video, error) {
	return s.videoRepo.GetRecommendedVideos(ctx, profileID, limit)
}