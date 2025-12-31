package oauth

import (
	"go_boilerplate/internal/shared/utils"

	"github.com/gofiber/fiber/v2"
)

// OAuthHandler defines the interface for OAuth HTTP handlers
type OAuthHandler interface {
	GoogleLogin(c *fiber.Ctx) error
	GoogleCallback(c *fiber.Ctx) error
	GitHubLogin(c *fiber.Ctx) error
	GitHubCallback(c *fiber.Ctx) error
}

// oauthHandler implements OAuthHandler interface
type oauthHandler struct {
	service OAuthService
}

// NewOAuthHandler creates a new OAuth handler
func NewOAuthHandler(service OAuthService) OAuthHandler {
	return &oauthHandler{service: service}
}

// GoogleLogin initiates Google OAuth login
func (h *oauthHandler) GoogleLogin(c *fiber.Ctx) error {
	url := h.service.GetGoogleAuthURL()

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"auth_url": url,
		},
	})
}

// GoogleCallback handles Google OAuth callback
func (h *oauthHandler) GoogleCallback(c *fiber.Ctx) error {
	// Get authorization code
	code := c.Query("code")
	if code == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Authorization code is required", nil)
	}

	// Handle OAuth callback
	response, err := h.service.HandleGoogleCallback(code)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "OAuth authentication failed", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, response, "OAuth authentication successful")
}

// GitHubLogin initiates GitHub OAuth login
func (h *oauthHandler) GitHubLogin(c *fiber.Ctx) error {
	url := h.service.GetGitHubAuthURL()

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"auth_url": url,
		},
	})
}

// GitHubCallback handles GitHub OAuth callback
func (h *oauthHandler) GitHubCallback(c *fiber.Ctx) error {
	// Get authorization code
	code := c.Query("code")
	if code == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Authorization code is required", nil)
	}

	// Handle OAuth callback
	response, err := h.service.HandleGitHubCallback(code)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "OAuth authentication failed", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, response, "OAuth authentication successful")
}
