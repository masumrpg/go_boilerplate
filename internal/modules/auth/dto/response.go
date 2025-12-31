package dto

import (
	"time"

	"go_boilerplate/internal/modules/user/dto"
	"github.com/google/uuid"
)

// AuthResponse represents an authentication response
type AuthResponse struct {
	AccessToken  string                    `json:"access_token"`
	RefreshToken string                    `json:"refresh_token"`
	ExpiresIn    int64                     `json:"expires_in"`
	User         dto.UserRoleResponse      `json:"user"`
}

// MessageResponse represents a simple message response
type MessageResponse struct {
	Message string `json:"message"`
}

// TokenInfo represents token information
type TokenInfo struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// RefreshToken represents a refresh token in the database
type RefreshToken struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Token     string    `json:"token" gorm:"type:varchar(500);uniqueIndex;not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName specifies the table name for RefreshToken
func (RefreshToken) TableName() string {
	return "t_refresh_tokens"
}
