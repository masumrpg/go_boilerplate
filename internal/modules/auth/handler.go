package auth

import (
	"go_boilerplate/internal/modules/auth/dto"
	"go_boilerplate/internal/shared/utils"

	"github.com/gofiber/fiber/v2"
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
func (h *authHandler) Register(c *fiber.Ctx) error {
	// Get validated body from context
	req := c.Locals("validatedBody").(*dto.RegisterRequest)

	// Register user
	response, err := h.service.Register(req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Registration failed", err)
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, response, "Registration successful")
}

// Login logs in a user
func (h *authHandler) Login(c *fiber.Ctx) error {
	// Get validated body from context
	req := c.Locals("validatedBody").(*dto.LoginRequest)

	// Login user
	response, err := h.service.Login(req)
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
func (h *authHandler) RefreshToken(c *fiber.Ctx) error {
	// Get validated body from context
	req := c.Locals("validatedBody").(*dto.RefreshTokenRequest)

	// Refresh token
	response, err := h.service.RefreshToken(req.RefreshToken)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Token refresh failed", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, response, "Token refreshed successfully")
}

// Logout logs out a user
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
func (h *authHandler) VerifyEmail(c *fiber.Ctx) error {
	req := c.Locals("validatedBody").(*dto.VerifyEmailRequest)

	if err := h.service.VerifyEmail(req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Email verification failed", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Email verified successfully")
}

// Verify2FA verifies 2FA code
func (h *authHandler) Verify2FA(c *fiber.Ctx) error {
	req := c.Locals("validatedBody").(*dto.Verify2FARequest)

	response, err := h.service.Verify2FA(req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "2FA verification failed", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, response, "2FA verified successfully")
}

// ResendVerification resends the activation code
func (h *authHandler) ResendVerification(c *fiber.Ctx) error {
	req := c.Locals("validatedBody").(*dto.ResendCodeRequest)

	if err := h.service.ResendVerification(req.Email); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to resend activation code", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Activation code resent successfully")
}

// Resend2FA resends the 2FA code
func (h *authHandler) Resend2FA(c *fiber.Ctx) error {
	req := c.Locals("validatedBody").(*dto.ResendCodeRequest)

	if err := h.service.Resend2FA(req.Email); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to resend 2FA code", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, nil, "2FA code resent successfully")
}
