package auth

import (
	"context"
	"errors"
	"time"

	"go_boilerplate/internal/modules/auth/dto"
	"go_boilerplate/internal/modules/email"
	"go_boilerplate/internal/modules/user"
	userdto "go_boilerplate/internal/modules/user/dto"
	"go_boilerplate/internal/shared/config"
	"go_boilerplate/internal/shared/utils"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// AuthService defines the interface for authentication business logic
type AuthService interface {
	Register(req *dto.RegisterRequest, metadata dto.SessionMetadata) (*dto.AuthResponse, error)
	Login(req *dto.LoginRequest, metadata dto.SessionMetadata) (*dto.AuthResponse, error)
	RefreshToken(refreshToken string, metadata dto.SessionMetadata) (*dto.AuthResponse, error)
	Logout(refreshToken string) error
	VerifyEmail(req *dto.VerifyEmailRequest) error
	Verify2FA(req *dto.Verify2FARequest, metadata dto.SessionMetadata) (*dto.AuthResponse, error)
	ResendVerification(email string) error
	Resend2FA(email string) error
	GetSessions(userID uuid.UUID) ([]dto.Session, error)
	DeleteSession(userID uuid.UUID, sessionID uuid.UUID) error
	BlockSession(userID uuid.UUID, sessionID uuid.UUID) error
}

// authService implements AuthService interface
type authService struct {
	userService  user.UserService
	jwtManager   *utils.JWTManager
	db           *gorm.DB
	cfg          *config.Config
	emailService email.EmailService
	redis        *redis.Client
}

// NewAuthService creates a new auth service
func NewAuthService(
	userService user.UserService,
	db *gorm.DB,
	cfg *config.Config,
	emailService email.EmailService,
	redis *redis.Client,
) AuthService {
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
		redis:        redis,
	}
}

// Register registers a new user
func (s *authService) Register(req *dto.RegisterRequest, metadata dto.SessionMetadata) (*dto.AuthResponse, error) {
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

	// Check if email verification is enabled
	if s.cfg.Security.EmailVerificationEnabled {
		// Generate and send verification code
		code := utils.RandomIntString(6)
		// Save 6-digit code to Redis with 10m expiry
		key := "activation:" + req.Email
		if err := s.redis.Set(context.Background(), key, code, 10*time.Minute).Err(); err != nil {
			return nil, errors.New("failed to save verification code")
		}

		// Send email asynchronously
		go func() {
			if s.emailService != nil {
				s.emailService.SendVerificationEmail(req.Email, code)
			}
		}()

		return &dto.AuthResponse{
			Message: "Registration successful. Please check your email to activate your account.",
		}, nil
	}

	// If verification disabled, set verified = true immediately (if not already default)
	if !s.cfg.Security.EmailVerificationEnabled {
		s.db.Model(&user.User{}).Where("id = ?", createdUser.ID).Update("is_verified", true)
	}

	return s.generateAuthResponse(createdUser.ID, metadata)
}

// Login authenticates a user
func (s *authService) Login(req *dto.LoginRequest, metadata dto.SessionMetadata) (*dto.AuthResponse, error) {
	// Validate password
	authenticatedUser, err := s.userService.ValidatePassword(req.Email, req.Password)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Check verification status (skip for SuperAdmin)
	// Get full profile to check role
	userWithRole, err := s.userService.GetProfileWithRole(authenticatedUser.ID)
	if err != nil {
		return nil, errors.New("failed to load user profile")
	}

	isSuperAdmin := userWithRole.Role != nil && userWithRole.Role.Slug == "super_admin"

	// If verification enabled and user not verified, deny login (unless SuperAdmin)
	if s.cfg.Security.EmailVerificationEnabled && !userWithRole.IsVerified && !isSuperAdmin {
		return nil, errors.New("account not verified. please verify your email")
	}

	// Check Two-Factor Authentication
	// Applicable if enabled and NOT SuperAdmin
	if s.cfg.Security.TwoFactorEnabled && !isSuperAdmin {
		// Generate 2FA code
		code := utils.RandomIntString(6)
		key := "2fa:" + req.Email
		if err := s.redis.Set(context.Background(), key, code, 5*time.Minute).Err(); err != nil {
			return nil, errors.New("failed to generate 2fa code")
		}

		// Send Email
		go func() {
			if s.emailService != nil {
				s.emailService.SendTwoFactorEmail(req.Email, code)
			}
		}()

		return &dto.AuthResponse{
			User:        userWithRole,
			Message:     "2FA Required",
			Requires2FA: true,
		}, nil
	}

	// Normal Login
	return s.generateAuthResponse(authenticatedUser.ID, metadata)
}

// VerifyEmail verifies user email
func (s *authService) VerifyEmail(req *dto.VerifyEmailRequest) error {
	key := "activation:" + req.Email
	storedCode, err := s.redis.Get(context.Background(), key).Result()
	if err != nil || storedCode != req.Code {
		return errors.New("invalid or expired activation code")
	}

	// Update user status
	if err := s.db.Model(&user.User{}).Where("email = ?", req.Email).Update("is_verified", true).Error; err != nil {
		return errors.New("failed to verify user")
	}

	// Delete code
	s.redis.Del(context.Background(), key)
	return nil
}

// Verify2FA verifies login OTP
func (s *authService) Verify2FA(req *dto.Verify2FARequest, metadata dto.SessionMetadata) (*dto.AuthResponse, error) {
	key := "2fa:" + req.Email
	storedCode, err := s.redis.Get(context.Background(), key).Result()
	if err != nil || storedCode != req.Code {
		return nil, errors.New("invalid or expired OTP")
	}

	// Get User
	foundUser, err := s.userService.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Delete code
	s.redis.Del(context.Background(), key)

	return s.generateAuthResponse(foundUser.ID, metadata)
}

// ResendVerification resends the activation code
func (s *authService) ResendVerification(email string) error {
	if !s.cfg.Security.EmailVerificationEnabled {
		return errors.New("email verification is not enabled")
	}

	// Check if user exists and is not verified
	user, err := s.userService.GetByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}

	if user.IsVerified {
		return errors.New("account already verified")
	}

	// Generate and send code
	code := utils.RandomIntString(6)
	key := "activation:" + email
	if err := s.redis.Set(context.Background(), key, code, 10*time.Minute).Err(); err != nil {
		return errors.New("failed to resend verification code")
	}

	go func() {
		if s.emailService != nil {
			s.emailService.SendVerificationEmail(email, code)
		}
	}()

	return nil
}

// Resend2FA resends the 2FA code
func (s *authService) Resend2FA(email string) error {
	if !s.cfg.Security.TwoFactorEnabled {
		return errors.New("2FA is not enabled")
	}

	// Check if user exists
	_, err := s.userService.GetByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}

	// Generate and send code
	code := utils.RandomIntString(6)
	key := "2fa:" + email
	if err := s.redis.Set(context.Background(), key, code, 5*time.Minute).Err(); err != nil {
		return errors.New("failed to resend 2FA code")
	}

	go func() {
		if s.emailService != nil {
			s.emailService.SendTwoFactorEmail(email, code)
		}
	}()

	return nil
}

// generateAuthResponse helps to dry up token generation logic
func (s *authService) generateAuthResponse(userID uuid.UUID, metadata dto.SessionMetadata) (*dto.AuthResponse, error) {
	// Load user with role information
	userWithRole, err := s.userService.GetProfileWithRole(userID)
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
		userID,
		userWithRole.Email,
		roleSlug,
		permissions,
	)
	if err != nil {
		return nil, errors.New("failed to generate tokens")
	}

	// Save session to database
	if err := s.saveSession(userID, refreshToken, metadata); err != nil {
		return nil, err
	}

	// Calculate expires in
	expiresIn := int64(s.cfg.JWT.AccessExpiry.Seconds())

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		User:         userWithRole,
	}, nil
}


// RefreshToken refreshes an access token using a refresh token
func (s *authService) RefreshToken(refreshToken string, metadata dto.SessionMetadata) (*dto.AuthResponse, error) {
	// Validate refresh token
	claims, err := s.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	// Check if session exists in database
	var storedSession dto.Session
	if err := s.db.Where("token = ? AND expires_at > ? AND is_blocked = ?", refreshToken, time.Now(), false).First(&storedSession).Error; err != nil {
		return nil, errors.New("session not found, expired, or blocked")
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

	// Delete old session
	s.db.Delete(&storedSession)

	// Save new session
	if err := s.saveSession(claims.UserID, newRefreshToken, metadata); err != nil {
		return nil, err
	}

	// Calculate expires in
	expiresIn := int64(s.cfg.JWT.AccessExpiry.Seconds())

	return &dto.AuthResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    expiresIn,
		User:         userProfile,
	}, nil
}

// Logout logs out a user by deleting their refresh token
func (s *authService) Logout(refreshToken string) error {
	// Delete session from database
	if err := s.db.Where("token = ?", refreshToken).Delete(&dto.Session{}).Error; err != nil {
		return err
	}

	return nil
}

// saveSession saves a session to the database
func (s *authService) saveSession(userID uuid.UUID, token string, metadata dto.SessionMetadata) error {
	expiresAt := time.Now().Add(s.cfg.JWT.RefreshExpiry)

	session := &dto.Session{
		UserID:    userID,
		Token:     token,
		IPAddress: metadata.IPAddress,
		UserAgent: metadata.UserAgent,
		DeviceID:  metadata.DeviceID,
		ExpiresAt: expiresAt,
		LastActive: time.Now(),
	}

	if err := s.db.Create(session).Error; err != nil {
		return err
	}

	return nil
}

// GetSessions returns all active sessions for a user
func (s *authService) GetSessions(userID uuid.UUID) ([]dto.Session, error) {
	var sessions []dto.Session
	if err := s.db.Where("user_id = ?", userID).Order("last_active desc").Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

// DeleteSession deletes a specific session
func (s *authService) DeleteSession(userID uuid.UUID, sessionID uuid.UUID) error {
	if err := s.db.Where("id = ? AND user_id = ?", sessionID, userID).Delete(&dto.Session{}).Error; err != nil {
		return err
	}
	return nil
}

// BlockSession blocks a specific session
func (s *authService) BlockSession(userID uuid.UUID, sessionID uuid.UUID) error {
	if err := s.db.Model(&dto.Session{}).Where("id = ? AND user_id = ?", sessionID, userID).Update("is_blocked", true).Error; err != nil {
		return err
	}
	return nil
}
