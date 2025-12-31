package auth

import (
	"go_boilerplate/internal/shared/utils"
	"go_boilerplate/internal/modules/auth/dto"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler defines the interface for auth HTTP handlers
type AuthHandler interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	RefreshToken(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
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

	return utils.SuccessResponse(c, fiber.StatusOK, response, "Login successful")
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
