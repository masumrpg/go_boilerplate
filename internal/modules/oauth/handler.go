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
// @Summary Google Login
// @Description Get the URL to initiate Google OAuth2 login.
// @Tags OAuth
// @Produce json
// @Success 200 {object} utils.APIResponse "Auth URL retrieved"
// @Router /oauth/google [get]
func (h *oauthHandler) GoogleLogin(c *fiber.Ctx) error {
	url := h.service.GetGoogleAuthURL()

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"auth_url": url,
	}, "Auth URL retrieved successfully")
}

// GoogleCallback handles Google OAuth callback
// @Summary Google Callback
// @Description Handle the callback from Google OAuth2.
// @Tags OAuth
// @Produce json
// @Param code query string true "Authorization code from Google"
// @Success 200 {object} utils.APIResponse "Login successful"
// @Failure 400 {object} utils.APIResponse "Authentication failed"
// @Router /oauth/google/callback [get]
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
// @Summary GitHub Login
// @Description Get the URL to initiate GitHub OAuth2 login.
// @Tags OAuth
// @Produce json
// @Success 200 {object} utils.APIResponse "Auth URL retrieved"
// @Router /oauth/github [get]
func (h *oauthHandler) GitHubLogin(c *fiber.Ctx) error {
	url := h.service.GetGitHubAuthURL()

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"auth_url": url,
	}, "Auth URL retrieved successfully")
}

// GitHubCallback handles GitHub OAuth callback
// @Summary GitHub Callback
// @Description Handle the callback from GitHub OAuth2.
// @Tags OAuth
// @Produce json
// @Param code query string true "Authorization code from GitHub"
// @Success 200 {object} utils.APIResponse "Login successful"
// @Failure 400 {object} utils.APIResponse "Authentication failed"
// @Router /oauth/github/callback [get]
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
