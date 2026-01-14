package service

import (
	"context"

	"github.com/scouttalent/discovery-service/internal/model"
	"github.com/scouttalent/discovery-service/internal/repository"
	"go.uber.org/zap"
)

type FeedService struct {
	videoRepo *repository.VideoRepository
	logger    *zap.Logger
}

func NewFeedService(videoRepo *repository.VideoRepository, logger *zap.Logger) *FeedService {
	return &FeedService{
		videoRepo: videoRepo,
		logger:    logger,
	}
}

func (s *FeedService) GetFeed(ctx context.Context, limit, offset int) ([]model.Video, int, error) {
	return s.videoRepo.GetFeed(ctx, limit, offset)
}

func (s *FeedService) GetTrendingVideos(ctx context.Context, limit int) ([]model.Video, error) {
	return s.videoRepo.GetTrendingVideos(ctx, limit)
}