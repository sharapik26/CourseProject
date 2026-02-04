package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID           int64          `json:"id"`
	Email        string         `json:"email"`
	PasswordHash string         `json:"-"`
	Username     string         `json:"username"`
	AvatarURL    sql.NullString `json:"avatar_url,omitempty"`
	BirthDate    sql.NullTime   `json:"birth_date,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type UserResponse struct {
	ID        int64   `json:"id"`
	Email     string  `json:"email"`
	Username  string  `json:"username"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	BirthDate *string `json:"birth_date,omitempty"`
	CreatedAt string  `json:"created_at"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Username string `json:"username" binding:"required,min=2"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateProfileRequest struct {
	Username  *string `json:"username,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	BirthDate *string `json:"birth_date,omitempty"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (u *User) ToResponse() UserResponse {
	resp := UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Username:  u.Username,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
	}

	if u.AvatarURL.Valid {
		resp.AvatarURL = &u.AvatarURL.String
	}

	if u.BirthDate.Valid {
		dateStr := u.BirthDate.Time.Format("2006-01-02")
		resp.BirthDate = &dateStr
	}

	return resp
}

