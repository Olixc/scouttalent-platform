package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/scouttalent/discovery-service/internal/model"
	"github.com/scouttalent/discovery-service/internal/service"
)

type SearchHandler struct {
	service *service.SearchService
}

func NewSearchHandler(service *service.SearchService) *SearchHandler {
	return &SearchHandler{service: service}
}

func (h *SearchHandler) SearchProfiles(c *gin.Context) {
	query := c.Query("q")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	var filters model.ProfileFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profiles, total, err := h.service.SearchProfiles(c.Request.Context(), query, filters, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search profiles"})
		return
	}

	c.JSON(http.StatusOK, model.SearchResponse{
		Results: profiles,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
	})
}

func (h *SearchHandler) SearchVideos(c *gin.Context) {
	query := c.Query("q")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	videos, total, err := h.service.SearchVideos(c.Request.Context(), query, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search videos"})
		return
	}

	c.JSON(http.StatusOK, model.SearchResponse{
		Results: videos,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
	})
}