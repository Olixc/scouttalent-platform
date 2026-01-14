package model

import (
	"time"
)

type User struct {
	ID              string    `json:"id" db:"id"`
	Email           string    `json:"email" db:"email"`
	EmailVerified   bool      `json:"email_verified" db:"email_verified"`
	Phone           *string   `json:"phone,omitempty" db:"phone"`
	PhoneVerified   bool      `json:"phone_verified" db:"phone_verified"`
	PasswordHash    string    `json:"-" db:"password_hash"`
	Role            string    `json:"role" db:"role"` // player, scout, academy, admin
	Status          string    `json:"status" db:"status"` // active, suspended, banned
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
	LastLoginAt     *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Role     string `json:"role" binding:"required,oneof=player scout academy"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type VerifyEmailRequest struct {
	Code string `json:"code" binding:"required,len=6"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Code        string `json:"code" binding:"required,len=6"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}