package dto

import (
	"time"

	"go_boilerplate/internal/shared/utils"

	"github.com/google/uuid"
)

// RoleResponse represents a role response
type RoleResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Permissions []string  `json:"permissions"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// RolesResponse represents a paginated list of roles
type RolesResponse struct {
	Roles []RoleResponse     `json:"roles"`
	Meta  utils.PaginationMeta `json:"meta"`
}

// UserRoleResponse represents user with role information
type UserRoleResponse struct {
	ID        uuid.UUID       `json:"id"`
	Name      string          `json:"name"`
	Email     string          `json:"email"`
	Role      *RoleInfo       `json:"role"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// RoleInfo represents simplified role information
type RoleInfo struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Permissions []string  `json:"permissions"`
}
