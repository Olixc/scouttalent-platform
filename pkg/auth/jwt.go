package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID      string   `json:"sub"`
	ProfileID   string   `json:"profile_id,omitempty"`
	Role        string   `json:"role"`
	TrustLevel  string   `json:"trust_level"`
	Permissions []string `json:"permissions"`
}

type JWTConfig struct {
	Secret string
}

type TokenConfig struct {
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	Issuer               string
	Audience             []string
	SecretKey            string
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

// GenerateAccessToken creates a new access token
func GenerateAccessToken(userID, role, trustLevel string, permissions []string, config TokenConfig) (string, error) {
	now := time.Now()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.Issuer,
			Subject:   userID,
			Audience:  config.Audience,
			ExpiresAt: jwt.NewNumericDate(now.Add(config.AccessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
		UserID:      userID,
		Role:        role,
		TrustLevel:  trustLevel,
		Permissions: permissions,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.SecretKey))
}

// GenerateRefreshToken creates a new refresh token
func GenerateRefreshToken(userID string, config TokenConfig) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Issuer:    config.Issuer,
		Subject:   userID,
		Audience:  config.Audience,
		ExpiresAt: jwt.NewNumericDate(now.Add(config.RefreshTokenDuration)),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.SecretKey))
}

// ValidateToken validates and parses a JWT token
func ValidateToken(tokenString string, config TokenConfig) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.SecretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}