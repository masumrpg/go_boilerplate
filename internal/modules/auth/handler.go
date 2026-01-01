package auth

import (
	"go_boilerplate/internal/modules/auth/dto"
	"go_boilerplate/internal/shared/middleware"
	"go_boilerplate/internal/shared/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// AuthHandler defines the interface for auth HTTP handlers
type AuthHandler interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	RefreshToken(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	VerifyEmail(c *fiber.Ctx) error
	Verify2FA(c *fiber.Ctx) error
	ResendVerification(c *fiber.Ctx) error
	Resend2FA(c *fiber.Ctx) error
	GetSessions(c *fiber.Ctx) error
	DeleteSession(c *fiber.Ctx) error
	BlockSession(c *fiber.Ctx) error
}

// authHandler implements AuthHandler interface
type authHandler struct {
	service AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(service AuthService) AuthHandler {
	return &authHandler{service: service}
}

// Register registers a new user
// @Summary Register a new user
// @Description Create a new user account with name, email and password.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration data"
// @Success 201 {object} utils.APIResponse{data=dto.AuthResponse} "Registration successful"
// @Failure 400 {object} utils.APIResponse "Invalid request data"
// @Router /auth/register [post]
func (h *authHandler) Register(c *fiber.Ctx) error {
	// Get validated body from context
	req := c.Locals("validatedBody").(*dto.RegisterRequest)

	// Register user
	response, err := h.service.Register(req, h.getMetadata(c))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Registration failed", err)
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, response, "Registration successful")
}

// Login logs in a user
// @Summary Login user
// @Description Authenticate user and return tokens.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} utils.APIResponse{data=dto.AuthResponse} "Login successful"
// @Failure 401 {object} utils.APIResponse "Invalid credentials"
// @Router /auth/login [post]
func (h *authHandler) Login(c *fiber.Ctx) error {
	// Get validated body from context
	req := c.Locals("validatedBody").(*dto.LoginRequest)

	// Login user
	response, err := h.service.Login(req, h.getMetadata(c))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Login failed", err)
	}

	message := "Login successful"
	if response.Requires2FA {
		message = "2FA Required"
	}

	return utils.SuccessResponse(c, fiber.StatusOK, response, message)
}

// RefreshToken refreshes an access token
// @Summary Refresh access token
// @Description Get a new access token using a refresh token.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} utils.APIResponse{data=dto.AuthResponse} "Token refreshed"
// @Failure 401 {object} utils.APIResponse "Invalid or expired refresh token"
// @Router /auth/refresh [post]
func (h *authHandler) RefreshToken(c *fiber.Ctx) error {
	// Get validated body from context
	req := c.Locals("validatedBody").(*dto.RefreshTokenRequest)

	// Refresh token
	response, err := h.service.RefreshToken(req.RefreshToken, h.getMetadata(c))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Token refresh failed", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, response, "Token refreshed successfully")
}

// Logout logs out a user
// @Summary Logout user
// @Description Invalidate the refresh token.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token to invalidate"
// @Success 200 {object} utils.APIResponse "Logout successful"
// @Failure 500 {object} utils.APIResponse "Logout failed"
// @Router /auth/logout [post]
func (h *authHandler) Logout(c *fiber.Ctx) error {
	// Get validated body from context
	req := c.Locals("validatedBody").(*dto.RefreshTokenRequest)

	// Logout user
	if err := h.service.Logout(req.RefreshToken); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Logout failed", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Logout successful")
}

// VerifyEmail verifies a user's email
// @Summary Verify email
// @Description Complete account activation using the code sent to email.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.VerifyEmailRequest true "Verification data"
// @Success 200 {object} utils.APIResponse "Email verified"
// @Failure 400 {object} utils.APIResponse "Invalid or expired code"
// @Router /auth/verify-email [post]
func (h *authHandler) VerifyEmail(c *fiber.Ctx) error {
	req := c.Locals("validatedBody").(*dto.VerifyEmailRequest)

	if err := h.service.VerifyEmail(req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Email verification failed", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Email verified successfully")
}

// Verify2FA verifies 2FA code
// @Summary Verify 2FA code
// @Description Complete login using the 2FA code sent to email.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.Verify2FARequest true "2FA data"
// @Success 200 {object} utils.APIResponse{data=dto.AuthResponse} "2FA verified"
// @Failure 401 {object} utils.APIResponse "Invalid or expired OTP"
// @Router /auth/verify-2fa [post]
func (h *authHandler) Verify2FA(c *fiber.Ctx) error {
	req := c.Locals("validatedBody").(*dto.Verify2FARequest)

	response, err := h.service.Verify2FA(req, h.getMetadata(c))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "2FA verification failed", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, response, "2FA verified successfully")
}

// ResendVerification resends the activation code
// @Summary Resend verification email
// @Description Request a new activation code to be sent to email.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.ResendCodeRequest true "Email address"
// @Success 200 {object} utils.APIResponse "Code resent"
// @Failure 400 {object} utils.APIResponse "Invalid request"
// @Router /auth/resend-verification [post]
func (h *authHandler) ResendVerification(c *fiber.Ctx) error {
	req := c.Locals("validatedBody").(*dto.ResendCodeRequest)

	if err := h.service.ResendVerification(req.Email); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to resend activation code", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Activation code resent successfully")
}

// Resend2FA resends the 2FA code
// @Summary Resend 2FA email
// @Description Request a new 2FA code to be sent to email.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.ResendCodeRequest true "Email address"
// @Success 200 {object} utils.APIResponse "Code resent"
// @Failure 400 {object} utils.APIResponse "Invalid request"
// @Router /auth/resend-2fa [post]
func (h *authHandler) Resend2FA(c *fiber.Ctx) error {
	req := c.Locals("validatedBody").(*dto.ResendCodeRequest)

	if err := h.service.Resend2FA(req.Email); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to resend 2FA code", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, nil, "2FA code resent successfully")
}

// GetSessions returns all active sessions for a user
// @Summary Get active sessions
// @Description Retrieve all active login sessions/devices for the current user.
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.APIResponse{data=[]dto.Session} "Sessions retrieved"
// @Failure 401 {object} utils.APIResponse "Unauthorized"
// @Router /auth/sessions [get]
func (h *authHandler) GetSessions(c *fiber.Ctx) error {
	userIDStr, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	userID, _ := uuid.Parse(userIDStr)
	sessions, err := h.service.GetSessions(userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get sessions", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, sessions, "Sessions retrieved successfully")
}

// DeleteSession deletes a specific session
// @Summary Logout from a specific device
// @Description Terminate a specific session by its ID.
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Param id path string true "Session ID (UUID)"
// @Success 200 {object} utils.APIResponse "Session deleted"
// @Failure 400 {object} utils.APIResponse "Invalid session ID"
// @Failure 401 {object} utils.APIResponse "Unauthorized"
// @Router /auth/sessions/{id} [delete]
func (h *authHandler) DeleteSession(c *fiber.Ctx) error {
	userIDStr, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	userID, _ := uuid.Parse(userIDStr)
	sessionID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid session ID", err)
	}

	if err := h.service.DeleteSession(userID, sessionID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete session", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Session deleted successfully")
}

// BlockSession blocks a specific session
// @Summary Block a specific device/session
// @Description Block a session so it cannot be used anymore.
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Param id path string true "Session ID (UUID)"
// @Success 200 {object} utils.APIResponse "Session blocked"
// @Failure 400 {object} utils.APIResponse "Invalid session ID"
// @Failure 401 {object} utils.APIResponse "Unauthorized"
// @Router /auth/sessions/{id}/block [patch]
func (h *authHandler) BlockSession(c *fiber.Ctx) error {
	userIDStr, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	userID, _ := uuid.Parse(userIDStr)
	sessionID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid session ID", err)
	}

	if err := h.service.BlockSession(userID, sessionID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to block session", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Session blocked successfully")
}

// getMetadata extracts session metadata from fiber.Ctx
func (h *authHandler) getMetadata(c *fiber.Ctx) dto.SessionMetadata {
	return dto.SessionMetadata{
		IPAddress: c.IP(),
		UserAgent: string(c.Request().Header.UserAgent()),
		DeviceID:  c.Get("X-Device-ID"),
	}
}
