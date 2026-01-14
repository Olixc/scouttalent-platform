package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/scouttalent/auth-service/internal/config"
	"github.com/scouttalent/auth-service/internal/model"
	"github.com/scouttalent/auth-service/internal/repository"
	"github.com/scouttalent/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo   *repository.UserRepository
	config *config.Config
}

func NewAuthService(repo *repository.UserRepository, config *config.Config) *AuthService {
	return &AuthService{
		repo:   repo,
		config: config,
	}
}

func (s *AuthService) Register(ctx context.Context, req model.RegisterRequest) (*model.User, error) {
	// Check if user already exists
	existing, err := s.repo.GetByEmail(ctx, req.Email)
	if err == nil && existing != nil {
		return nil, repository.ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &model.User{
		ID:           uuid.New().String(),
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         req.Role,
		Status:       "active",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, req model.LoginRequest) (*auth.TokenPair, *model.User, error) {
	// Get user by email
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, nil, fmt.Errorf("invalid credentials")
	}

	// Check user status
	if user.Status != "active" {
		return nil, nil, fmt.Errorf("account is %s", user.Status)
	}

	// Generate tokens
	accessToken, err := auth.GenerateAccessToken(
		user.ID,
		user.Role,
		"newcomer", // Default trust level
		s.getPermissionsForRole(user.Role),
		s.config.JWT,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID, s.config.JWT)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Update last login
	if err := s.repo.UpdateLastLogin(ctx, user.ID); err != nil {
		// Log error but don't fail the login
		fmt.Printf("failed to update last login: %v\n", err)
	}

	return &auth.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, user, nil
}

func (s *AuthService) getPermissionsForRole(role string) []string {
	permissions := map[string][]string{
		"player": {
			"upload:video",
			"delete:video",
			"view:profiles",
			"edit:profile",
		},
		"scout": {
			"view:profiles",
			"view:all_videos",
			"contact:player",
			"view:analytics",
		},
		"academy": {
			"view:profiles",
			"view:all_videos",
			"contact:player",
			"verify:profile",
		},
		"admin": {
			"*", // All permissions
		},
	}

	return permissions[role]
}