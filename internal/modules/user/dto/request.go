package dto

import "github.com/google/uuid"

// CreateUserRequest represents a request to create a new user
type CreateUserRequest struct {
	Name     string    `json:"name" validate:"required,min=3,max=100"`
	Email    string    `json:"email" validate:"required,email"`
	Password string    `json:"password" validate:"required,min=6,max=50"`
	RoleID   *uuid.UUID `json:"role_id" validate:"omitempty"` // Optional: if not provided, defaults to user role
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UpdateUserRequest represents a request to update a user
type UpdateUserRequest struct {
	Name   string    `json:"name" validate:"omitempty,min=3,max=100"`
	Email  string    `json:"email" validate:"omitempty,email"`
	RoleID *uuid.UUID `json:"role_id" validate:"omitempty"` // Optional: can update role to user or admin only
}

// ChangePasswordRequest represents a request to change password
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6,max=50"`
}

// AssignRoleRequest represents a request to assign a role to a user
type AssignRoleRequest struct {
	RoleID uuid.UUID `json:"role_id" validate:"required"`
}

