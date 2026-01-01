package dto

import (
	"time"

	"go_boilerplate/internal/modules/user/dto"

	"github.com/google/uuid"
)

// AuthResponse represents an authentication response
type AuthResponse struct {
	AccessToken  string                     `json:"access_token,omitempty"`
	RefreshToken string                     `json:"refresh_token,omitempty"`
	ExpiresIn    int64                      `json:"expires_in,omitempty"`
	User         *dto.UserRoleResponse      `json:"user,omitempty"`
	Message      string                     `json:"message,omitempty"`
	Requires2FA  bool                       `json:"requires_2fa,omitempty"`
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

// Session represents a user session/refresh token in the database
type Session struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Token     string    `json:"token" gorm:"type:varchar(500);uniqueIndex;not null"`
	IPAddress string    `json:"ip_address" gorm:"type:varchar(45)"`
	UserAgent string    `json:"user_agent" gorm:"type:text"`
	DeviceID  string    `json:"device_id" gorm:"type:varchar(255)"`
	IsBlocked bool      `json:"is_blocked" gorm:"default:false"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	LastActive time.Time `json:"last_active"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName specifies the table name for Session
func (Session) TableName() string {
	return "t_sessions"
}
