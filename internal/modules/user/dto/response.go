package dto

import (
	"time"

	"github.com/google/uuid"
)

// UserResponse represents a user response (without password and role)
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	IsVerified bool     `json:"is_verified"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRoleResponse represents a user response with role information
type UserRoleResponse struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Role      *RoleInfo  `json:"role"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// RoleInfo represents simplified role information
type RoleInfo struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Permissions []string  `json:"permissions"`
}

// UsersResponse represents a paginated list of users
type UsersResponse struct {
	Users []UserResponse `json:"users"`
	Meta  PaginationMeta  `json:"meta"`
}

// PaginationMeta contains pagination metadata
type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}


