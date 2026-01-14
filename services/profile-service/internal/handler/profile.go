package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scouttalent/profile-service/internal/model"
	"github.com/scouttalent/profile-service/internal/service"
	"go.uber.org/zap"
)

type ProfileHandler struct {
	service *service.ProfileService
	logger  *zap.Logger
}

func NewProfileHandler(service *service.ProfileService, logger *zap.Logger) *ProfileHandler {
	return &ProfileHandler{
		service: service,
		logger:  logger,
	}
}

func (h *ProfileHandler) CreateProfile(c *gin.Context) {
	var req model.CreateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var userType model.UserType
	switch role.(string) {
	case "player":
		userType = model.UserTypePlayer
	case "scout":
		userType = model.UserTypeScout
	case "academy":
		userType = model.UserTypeAcademy
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user role"})
		return
	}

	profile, err := h.service.CreateProfile(c.Request.Context(), userID.(string), userType, req)
	if err != nil {
		h.logger.Error("failed to create profile", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create profile"})
		return
	}

	c.JSON(http.StatusCreated, profile)
}

func (h *ProfileHandler) GetMyProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	profile, err := h.service.GetProfileByUserID(c.Request.Context(), userID.(string))
	if err != nil {
		h.logger.Error("failed to get profile", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	profileID := c.Param("id")

	profile, err := h.service.GetProfile(c.Request.Context(), profileID)
	if err != nil {
		h.logger.Error("failed to get profile", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	profileID := c.Param("id")

	var req model.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile, err := h.service.UpdateProfile(c.Request.Context(), profileID, req)
	if err != nil {
		h.logger.Error("failed to update profile", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (h *ProfileHandler) CreatePlayerDetails(c *gin.Context) {
	profileID := c.Param("id")

	var req model.CreatePlayerDetailsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	player, err := h.service.CreatePlayerDetails(c.Request.Context(), profileID, req)
	if err != nil {
		h.logger.Error("failed to create player details", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create player details"})
		return
	}

	c.JSON(http.StatusCreated, player)
}

func (h *ProfileHandler) GetPlayerProfile(c *gin.Context) {
	profileID := c.Param("id")

	player, err := h.service.GetPlayerProfile(c.Request.Context(), profileID)
	if err != nil {
		h.logger.Error("failed to get player profile", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "player profile not found"})
		return
	}

	c.JSON(http.StatusOK, player)
}

func (h *ProfileHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "profile-service",
	})
}