package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/scouttalent/media-service/internal/model"
	"github.com/scouttalent/media-service/internal/service"
	"go.uber.org/zap"
)

type MediaHandler struct {
	service *service.MediaService
	logger  *zap.Logger
}

func NewMediaHandler(service *service.MediaService, logger *zap.Logger) *MediaHandler {
	return &MediaHandler{
		service: service,
		logger:  logger,
	}
}

func (h *MediaHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "media-service",
	})
}

func (h *MediaHandler) InitiateUpload(c *gin.Context) {
	var req model.CreateVideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get profile ID from JWT claims
	profileID, exists := c.Get("profile_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "profile_id not found in token"})
		return
	}

	profileUUID, err := uuid.Parse(profileID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid profile_id"})
		return
	}

	video, upload, err := h.service.InitiateUpload(c.Request.Context(), profileUUID, &req)
	if err != nil {
		h.logger.Error("failed to initiate upload", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to initiate upload"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"video":     video,
		"upload":    upload,
		"upload_url": video.BlobURL,
	})
}

func (h *MediaHandler) ResumeUpload(c *gin.Context) {
	uploadID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid upload_id"})
		return
	}

	// Get progress from request
	progressStr := c.Query("progress")
	progress, err := strconv.Atoi(progressStr)
	if err != nil || progress < 0 || progress > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid progress value"})
		return
	}

	if err := h.service.UpdateUploadProgress(c.Request.Context(), uploadID, progress); err != nil {
		h.logger.Error("failed to update upload progress", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update upload progress"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "upload progress updated"})
}

func (h *MediaHandler) CompleteUpload(c *gin.Context) {
	videoID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video_id"})
		return
	}

	if err := h.service.CompleteUpload(c.Request.Context(), videoID); err != nil {
		h.logger.Error("failed to complete upload", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to complete upload"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "upload completed successfully"})
}

func (h *MediaHandler) GetVideo(c *gin.Context) {
	videoID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video_id"})
		return
	}

	video, err := h.service.GetVideo(c.Request.Context(), videoID)
	if err != nil {
		h.logger.Error("failed to get video", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "video not found"})
		return
	}

	c.JSON(http.StatusOK, model.VideoResponse{Video: video})
}

func (h *MediaHandler) ListProfileVideos(c *gin.Context) {
	profileID, err := uuid.Parse(c.Param("profile_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid profile_id"})
		return
	}

	// Parse pagination parameters
	limit := 20
	offset := 0

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	videos, total, err := h.service.ListProfileVideos(c.Request.Context(), profileID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list videos", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list videos"})
		return
	}

	c.JSON(http.StatusOK, model.VideoListResponse{
		Videos: videos,
		Total:  total,
	})
}

func (h *MediaHandler) UpdateVideo(c *gin.Context) {
	videoID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video_id"})
		return
	}

	var req model.UpdateVideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	video, err := h.service.UpdateVideo(c.Request.Context(), videoID, &req)
	if err != nil {
		h.logger.Error("failed to update video", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update video"})
		return
	}

	c.JSON(http.StatusOK, model.VideoResponse{
		Video:   video,
		Message: "video updated successfully",
	})
}

func (h *MediaHandler) DeleteVideo(c *gin.Context) {
	videoID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video_id"})
		return
	}

	if err := h.service.DeleteVideo(c.Request.Context(), videoID); err != nil {
		h.logger.Error("failed to delete video", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete video"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "video deleted successfully"})
}