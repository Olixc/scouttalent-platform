package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scouttalent/auth-service/internal/model"
	"github.com/scouttalent/auth-service/internal/service"
	"go.uber.org/zap"
)

type AuthHandler struct {
	service *service.AuthService
	logger  *zap.Logger
}

func NewAuthHandler(service *service.AuthService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		service: service,
		logger:  logger,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("failed to register user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user": user,
		"message": "Registration successful. Please verify your email.",
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, user, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("failed to login", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"message": "Authenticated",
	})
}

func (h *AuthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"service": "auth-service",
	})
}