package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/scouttalent/discovery-service/internal/service"
)

type RecommendationHandler struct {
	service *service.RecommendationService
}

func NewRecommendationHandler(service *service.RecommendationService) *RecommendationHandler {
	return &RecommendationHandler{service: service}
}

func (h *RecommendationHandler) GetProfileRecommendations(c *gin.Context) {
	profileID := c.GetString("profile_id")
	if profileID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Profile ID not found in token"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	profiles, err := h.service.GetProfileRecommendations(c.Request.Context(), profileID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get recommendations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"recommendations": profiles})
}

func (h *RecommendationHandler) GetVideoRecommendations(c *gin.Context) {
	profileID := c.GetString("profile_id")
	if profileID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Profile ID not found in token"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	videos, err := h.service.GetVideoRecommendations(c.Request.Context(), profileID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get recommendations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"recommendations": videos})
}