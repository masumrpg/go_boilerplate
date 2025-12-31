package dto

import (
	"github.com/google/uuid"
)

// CreateRoleRequest represents a request to create a role
type CreateRoleRequest struct {
	Name        string   `json:"name" validate:"required,min=3,max=100"`
	Slug        string   `json:"slug" validate:"required,min=2,max=50,alphanum"`
	Permissions []string `json:"permissions" validate:"required,min=1"`
	Description string   `json:"description" validate:"omitempty,max=500"`
}

// UpdateRoleRequest represents a request to update a role
type UpdateRoleRequest struct {
	Name        string   `json:"name" validate:"omitempty,min=3,max=100"`
	Permissions []string `json:"permissions" validate:"omitempty,min=1"`
	Description string   `json:"description" validate:"omitempty,max=500"`
}

// AssignRoleRequest represents a request to assign a role to a user
type AssignRoleRequest struct {
	RoleID uuid.UUID `json:"role_id" validate:"required"`
}
