package dto

import (
	"time"

	"github.com/google/uuid"
)

// UserResponse represents a user response (without password)
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UsersResponse represents a paginated list of users
type UsersResponse struct {
	Users []UserResponse `json:"users"`
	Meta  PaginationMeta `json:"meta"`
}

// PaginationMeta contains pagination metadata
type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// AuthResponse represents authentication response with tokens
type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}
