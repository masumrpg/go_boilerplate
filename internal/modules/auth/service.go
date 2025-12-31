package auth

import (
	"errors"
	"time"

	"go_boilerplate/internal/shared/config"
	"go_boilerplate/internal/shared/utils"
	"go_boilerplate/internal/modules/auth/dto"
	"go_boilerplate/internal/modules/email"
	"go_boilerplate/internal/modules/user"
	userdto "go_boilerplate/internal/modules/user/dto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuthService defines the interface for authentication business logic
type AuthService interface {
	Register(req *dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(req *dto.LoginRequest) (*dto.AuthResponse, error)
	RefreshToken(refreshToken string) (*dto.AuthResponse, error)
	Logout(refreshToken string) error
}

// authService implements AuthService interface
type authService struct {
	userService  user.UserService
	jwtManager   *utils.JWTManager
	db           *gorm.DB
	cfg          *config.Config
	emailService email.EmailService
}

// NewAuthService creates a new auth service
func NewAuthService(userService user.UserService, db *gorm.DB, cfg *config.Config, emailService email.EmailService) AuthService {
	jwtManager := utils.NewJWTManager(
		cfg.JWT.Secret,
		cfg.JWT.AccessExpiry,
		cfg.JWT.RefreshExpiry,
		cfg.JWT.Issuer,
	)

	return &authService{
		userService:  userService,
		jwtManager:   jwtManager,
		db:           db,
		cfg:          cfg,
		emailService: emailService,
	}
}

// Register registers a new user
func (s *authService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Create user request
	createUserReq := &userdto.CreateUserRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	// Create user (with default role assigned)
	createdUser, err := s.userService.CreateUser(createUserReq)
	if err != nil {
		return nil, err
	}

	// Load user with role information
	userWithRole, err := s.userService.GetProfileWithRole(createdUser.ID)
	if err != nil {
		return nil, errors.New("failed to load user role")
	}

	// Generate tokens with role information
	roleSlug := ""
	permissions := []string{}
	if userWithRole.Role != nil {
		roleSlug = userWithRole.Role.Slug
		permissions = userWithRole.Role.Permissions
	}

	accessToken, refreshToken, err := s.jwtManager.GenerateTokenPair(
		createdUser.ID,
		createdUser.Email,
		roleSlug,
		permissions,
	)
	if err != nil {
		return nil, errors.New("failed to generate tokens")
	}

	// Send welcome email if enabled
	if s.emailService != nil && s.cfg.Email.Enabled {
		// Send welcome email asynchronously (don't block the response)
		go func() {
			if err := s.emailService.SendWelcomeEmail(req.Email, req.Name); err != nil {
				// Log error but don't fail the registration
				println("Failed to send welcome email:", err.Error())
			}
		}()
	}

	// Save refresh token to database
	if err := s.saveRefreshToken(createdUser.ID, refreshToken); err != nil {
		return nil, err
	}

	// Calculate expires in (in seconds)
	expiresIn := int64(s.cfg.JWT.AccessExpiry.Seconds())

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		User:         *userWithRole,
	}, nil
}

// Login authenticates a user
func (s *authService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	// Validate password
	authenticatedUser, err := s.userService.ValidatePassword(req.Email, req.Password)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Load user with role information
	userWithRole, err := s.userService.GetProfileWithRole(authenticatedUser.ID)
	if err != nil {
		return nil, errors.New("failed to load user role")
	}

	// Generate tokens with role information
	roleSlug := ""
	permissions := []string{}
	if userWithRole.Role != nil {
		roleSlug = userWithRole.Role.Slug
		permissions = userWithRole.Role.Permissions
	}

	accessToken, refreshToken, err := s.jwtManager.GenerateTokenPair(
		authenticatedUser.ID,
		authenticatedUser.Email,
		roleSlug,
		permissions,
	)
	if err != nil {
		return nil, errors.New("failed to generate tokens")
	}

	// Save refresh token to database
	if err := s.saveRefreshToken(authenticatedUser.ID, refreshToken); err != nil {
		return nil, err
	}

	// Calculate expires in
	expiresIn := int64(s.cfg.JWT.AccessExpiry.Seconds())

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		User:         *userWithRole,
	}, nil
}

// RefreshToken refreshes an access token using a refresh token
func (s *authService) RefreshToken(refreshToken string) (*dto.AuthResponse, error) {
	// Validate refresh token
	claims, err := s.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	// Check if refresh token exists in database
	var storedToken dto.RefreshToken
	if err := s.db.Where("token = ? AND expires_at > ?", refreshToken, time.Now()).First(&storedToken).Error; err != nil {
		return nil, errors.New("refresh token not found or expired")
	}

	// Get user profile with role
	userProfile, err := s.userService.GetProfileWithRole(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Generate new tokens with role information
	roleSlug := ""
	permissions := []string{}
	if userProfile.Role != nil {
		roleSlug = userProfile.Role.Slug
		permissions = userProfile.Role.Permissions
	}

	newAccessToken, newRefreshToken, err := s.jwtManager.GenerateTokenPair(
		claims.UserID,
		claims.Email,
		roleSlug,
		permissions,
	)
	if err != nil {
		return nil, errors.New("failed to generate new tokens")
	}

	// Delete old refresh token
	s.db.Delete(&storedToken)

	// Save new refresh token
	if err := s.saveRefreshToken(claims.UserID, newRefreshToken); err != nil {
		return nil, err
	}

	// Calculate expires in
	expiresIn := int64(s.cfg.JWT.AccessExpiry.Seconds())

	return &dto.AuthResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    expiresIn,
		User:         *userProfile,
	}, nil
}

// Logout logs out a user by deleting their refresh token
func (s *authService) Logout(refreshToken string) error {
	// Delete refresh token from database
	if err := s.db.Where("token = ?", refreshToken).Delete(&dto.RefreshToken{}).Error; err != nil {
		return err
	}

	return nil
}

// saveRefreshToken saves a refresh token to the database
func (s *authService) saveRefreshToken(userID uuid.UUID, token string) error {
	expiresAt := time.Now().Add(s.cfg.JWT.RefreshExpiry)

	refreshToken := &dto.RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
	}

	if err := s.db.Create(refreshToken).Error; err != nil {
		return err
	}

	return nil
}
