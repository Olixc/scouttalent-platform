package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/scouttalent/discovery-service/internal/model"
	"github.com/scouttalent/discovery-service/internal/service"
)

type FeedHandler struct {
	service *service.FeedService
}

func NewFeedHandler(service *service.FeedService) *FeedHandler {
	return &FeedHandler{service: service}
}

func (h *FeedHandler) GetFeed(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	videos, total, err := h.service.GetFeed(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get feed"})
		return
	}

	c.JSON(http.StatusOK, model.SearchResponse{
		Results: videos,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
	})
}

func (h *FeedHandler) GetTrendingVideos(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	videos, err := h.service.GetTrendingVideos(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get trending videos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"videos": videos})
}