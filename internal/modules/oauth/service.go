package oauth

import (
	"context"
	"errors"

	authdto "go_boilerplate/internal/modules/auth/dto"
	"go_boilerplate/internal/modules/email"
	"go_boilerplate/internal/modules/oauth/dto"
	"go_boilerplate/internal/modules/user"
	userdto "go_boilerplate/internal/modules/user/dto"
	"go_boilerplate/internal/shared/config"
	"go_boilerplate/internal/shared/utils"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

// OAuthService defines the interface for OAuth operations
type OAuthService interface {
	GetGoogleAuthURL() string
	HandleGoogleCallback(code string) (*authdto.AuthResponse, error)
	GetGitHubAuthURL() string
	HandleGitHubCallback(code string) (*authdto.AuthResponse, error)
}

// oauthService implements OAuthService interface
type oauthService struct {
	db           *gorm.DB
	cfg          *config.Config
	userService  user.UserService
	emailService email.EmailService
	jwtManager   *utils.JWTManager
}

// NewOAuthService creates a new OAuth service
func NewOAuthService(db *gorm.DB, cfg *config.Config, userService user.UserService) OAuthService {
	jwtManager := utils.NewJWTManager(
		cfg.JWT.Secret,
		cfg.JWT.AccessExpiry,
		cfg.JWT.RefreshExpiry,
		cfg.JWT.Issuer,
	)

	// Initialize email service (optional, will check before sending)
	var emailService email.EmailService
	if cfg.Email.Enabled {
		// Import logger here - we'll get it from context or create a new one
		// For now, we'll initialize without logger
		emailService = email.NewEmailService(cfg, nil)
	}

	return &oauthService{
		db:           db,
		cfg:          cfg,
		userService:  userService,
		emailService: emailService,
		jwtManager:   jwtManager,
	}
}

// GetGoogleAuthURL returns the Google OAuth URL
func (s *oauthService) GetGoogleAuthURL() string {
	oauth2Config := &oauth2.Config{
		ClientID:     s.cfg.OAuth.Google.ClientID,
		ClientSecret: s.cfg.OAuth.Google.ClientSecret,
		RedirectURL:  s.cfg.OAuth.Google.RedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return oauth2Config.AuthCodeURL("state", oauth2.AccessTypeOffline)
}

// HandleGoogleCallback handles Google OAuth callback
func (s *oauthService) HandleGoogleCallback(code string) (*authdto.AuthResponse, error) {
	// Exchange code for token
	oauth2Config := &oauth2.Config{
		ClientID:     s.cfg.OAuth.Google.ClientID,
		ClientSecret: s.cfg.OAuth.Google.ClientSecret,
		RedirectURL:  s.cfg.OAuth.Google.RedirectURL,
		Endpoint:     google.Endpoint,
	}

	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		return nil, errors.New("failed to exchange token")
	}

	// Get user info from Google
	// Note: In production, you should make an HTTP request to get user info
	// For this boilerplate, we'll simulate it
	userInfo := &dto.OAuthUserInfo{
		ID:       "google_" + uuid.New().String(),
		Email:    "user@example.com", // In production, get from Google API
		Name:     "Google User",
		Provider: "google",
	}

	return s.handleOAuthUser(userInfo, token)
}

// GetGitHubAuthURL returns the GitHub OAuth URL
func (s *oauthService) GetGitHubAuthURL() string {
	oauth2Config := &oauth2.Config{
		ClientID:     s.cfg.OAuth.GitHub.ClientID,
		ClientSecret: s.cfg.OAuth.GitHub.ClientSecret,
		RedirectURL:  s.cfg.OAuth.GitHub.RedirectURL,
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}

	return oauth2Config.AuthCodeURL("state")
}

// HandleGitHubCallback handles GitHub OAuth callback
func (s *oauthService) HandleGitHubCallback(code string) (*authdto.AuthResponse, error) {
	// Exchange code for token
	oauth2Config := &oauth2.Config{
		ClientID:     s.cfg.OAuth.GitHub.ClientID,
		ClientSecret: s.cfg.OAuth.GitHub.ClientSecret,
		RedirectURL:  s.cfg.OAuth.GitHub.RedirectURL,
		Endpoint:     github.Endpoint,
	}

	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		return nil, errors.New("failed to exchange token")
	}

	// Get user info from GitHub
	// Note: In production, you should make an HTTP request to get user info
	// For this boilerplate, we'll simulate it
	userInfo := &dto.OAuthUserInfo{
		ID:       "github_" + uuid.New().String(),
		Email:    "user@example.com", // In production, get from GitHub API
		Name:     "GitHub User",
		Provider: "github",
	}

	return s.handleOAuthUser(userInfo, token)
}

// handleOAuthUser handles OAuth user login/registration
func (s *oauthService) handleOAuthUser(userInfo *dto.OAuthUserInfo, token *oauth2.Token) (*authdto.AuthResponse, error) {
	// Check if OAuth account exists
	var oauthAccount dto.OAuthAccount
	err := s.db.Where("provider = ? AND provider_id = ?", userInfo.Provider, userInfo.ID).First(&oauthAccount).Error

	var userID uuid.UUID
	isNewUser := false

	if err == nil {
		// OAuth account exists, use existing user
		userID = oauthAccount.UserID

		// Update token
		oauthAccount.AccessToken = token.AccessToken
		if token.RefreshToken != "" {
			oauthAccount.RefreshToken = token.RefreshToken
		}
		oauthAccount.ExpiresAt = token.Expiry
		s.db.Save(&oauthAccount)
	} else {
		// OAuth account doesn't exist, create new user
		isNewUser = true

		createUserReq := &userdto.CreateUserRequest{
			Name:     userInfo.Name,
			Email:    userInfo.Email,
			Password: uuid.New().String(), // Random password for OAuth users
		}

		createdUser, err := s.userService.CreateUser(createUserReq)
		if err != nil {
			// User might already exist with this email, link accounts
			// For simplicity, we'll return an error here
			return nil, errors.New("failed to create user")
		}

		userID = createdUser.ID

		// Create OAuth account
		oauthAccount = dto.OAuthAccount{
			UserID:       userID,
			Provider:     userInfo.Provider,
			ProviderID:   userInfo.ID,
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			ExpiresAt:    token.Expiry,
		}
		s.db.Create(&oauthAccount)
	}

	// Send welcome email if enabled and this is a new user
	if isNewUser && s.emailService != nil && s.cfg.Email.Enabled {
		// Check if welcome email is enabled for this provider
		sendWelcomeEmail := false
		if userInfo.Provider == "google" && s.cfg.OAuth.Google.SendWelcomeEmail {
			sendWelcomeEmail = true
		} else if userInfo.Provider == "github" && s.cfg.OAuth.GitHub.SendWelcomeEmail {
			sendWelcomeEmail = true
		}

		if sendWelcomeEmail {
			// Send welcome email asynchronously (don't block the response)
			go func() {
				if err := s.emailService.SendWelcomeEmail(userInfo.Email, userInfo.Name); err != nil {
					// Log error but don't fail the OAuth flow
					// In production, you might want to use proper logger
					println("Failed to send welcome email:", err.Error())
				}
			}()
		}
	}

	// Get user profile with role information
	userProfile, err := s.userService.GetProfileWithRole(userID)
	if err != nil {
		return nil, err
	}

	// Generate JWT tokens with role information
	roleSlug := ""
	permissions := []string{}
	if userProfile.Role != nil {
		roleSlug = userProfile.Role.Slug
		permissions = userProfile.Role.Permissions
	}

	accessToken, refreshToken, err := s.jwtManager.GenerateTokenPair(userID, userProfile.Email, roleSlug, permissions)
	if err != nil {
		return nil, errors.New("failed to generate tokens")
	}

	// Calculate expires in
	expiresIn := int64(s.cfg.JWT.AccessExpiry.Seconds())

	return &authdto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		User:         userProfile,
	}, nil
}
