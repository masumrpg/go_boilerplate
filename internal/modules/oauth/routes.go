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

	// Public OAuth routes
	oauth := api.Group("/oauth")
	oauth.Get("/google", oauthHandler.GoogleLogin)
	oauth.Get("/google/callback", oauthHandler.GoogleCallback)
	oauth.Get("/github", oauthHandler.GitHubLogin)
	oauth.Get("/github/callback", oauthHandler.GitHubCallback)
}
