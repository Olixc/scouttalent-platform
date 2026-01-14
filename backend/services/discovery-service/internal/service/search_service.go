package service

import (
	"context"

	"github.com/scouttalent/discovery-service/internal/model"
	"github.com/scouttalent/discovery-service/internal/repository"
	"go.uber.org/zap"
)

type SearchService struct {
	profileRepo *repository.ProfileRepository
	videoRepo   *repository.VideoRepository
	logger      *zap.Logger
}

func NewSearchService(profileRepo *repository.ProfileRepository, videoRepo *repository.VideoRepository, logger *zap.Logger) *SearchService {
	return &SearchService{
		profileRepo: profileRepo,
		videoRepo:   videoRepo,
		logger:      logger,
	}
}

func (s *SearchService) SearchProfiles(ctx context.Context, query string, filters model.ProfileFilters, limit, offset int) ([]model.Profile, int, error) {
	return s.profileRepo.SearchProfiles(ctx, query, filters, limit, offset)
}

func (s *SearchService) SearchVideos(ctx context.Context, query string, limit, offset int) ([]model.Video, int, error) {
	return s.videoRepo.SearchVideos(ctx, query, limit, offset)
}