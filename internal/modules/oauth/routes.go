package oauth

import (
	"go_boilerplate/internal/shared/config"
	"go_boilerplate/internal/modules/user"
	"go_boilerplate/internal/modules/oauth/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// RegisterRoutes registers all OAuth-related routes
func RegisterRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config, logger *logrus.Logger) {
	// Auto migrate OAuth account model
	db.AutoMigrate(&dto.OAuthAccount{})

	// Initialize user service (OAuth service depends on it)
	userRepo := user.NewUserRepository(db)
	userService := user.NewUserService(userRepo)

	// Initialize OAuth service
	oauthService := NewOAuthService(db, cfg, userService)

	// Initialize OAuth handler
	oauthHandler := NewOAuthHandler(oauthService)

	// Create API route group
	api := app.Group("/api/v1")

	// Register Google OAuth routes if enabled
	if cfg.OAuth.Google.Enabled {
		logger.Info("✓ Google OAuth routes registered (enabled)")
		oauth := api.Group("/oauth")
		oauth.Get("/google", oauthHandler.GoogleLogin)
		oauth.Get("/google/callback", oauthHandler.GoogleCallback)
	} else {
		logger.Info("✗ Google OAuth routes skipped (disabled)")
	}

	// Register GitHub OAuth routes if enabled
	if cfg.OAuth.GitHub.Enabled {
		logger.Info("✓ GitHub OAuth routes registered (enabled)")
		oauth := api.Group("/oauth")
		oauth.Get("/github", oauthHandler.GitHubLogin)
		oauth.Get("/github/callback", oauthHandler.GitHubCallback)
	} else {
		logger.Info("✗ GitHub OAuth routes skipped (disabled)")
	}
}
