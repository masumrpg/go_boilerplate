package dto

import (
	"time"

	"github.com/google/uuid"
)

// OAuthUserInfo represents user information from OAuth provider
type OAuthUserInfo struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Provider string `json:"provider"` // google, github
}

// OAuthAccount represents an OAuth account linked to a user
type OAuthAccount struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID       uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Provider     string    `json:"provider" gorm:"type:varchar(50);not null"`
	ProviderID   string    `json:"provider_id" gorm:"type:varchar(255);not null"`
	AccessToken  string    `json:"access_token" gorm:"type:text"`
	RefreshToken string    `json:"refresh_token" gorm:"type:text"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName specifies the table name for OAuthAccount model
func (OAuthAccount) TableName() string {
	return "t_oauth_accounts"
}
