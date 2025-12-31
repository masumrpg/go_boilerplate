package auth

import (
	"go_boilerplate/internal/shared/config"
	sharedmiddleware "go_boilerplate/internal/shared/middleware"
	"go_boilerplate/internal/modules/user"
	"go_boilerplate/internal/modules/auth/dto"
	"go_boilerplate/internal/modules/email"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// RegisterRoutes registers all auth-related routes
func RegisterRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config, logger *logrus.Logger) {
	// Initialize user service (auth service depends on it)
	userRepo := user.NewUserRepository(db)
	userService := user.NewUserService(userRepo)

	// Initialize email service (optional, will check before sending)
	var emailService email.EmailService
	if cfg.Email.Enabled {
		emailService = email.NewEmailService(cfg, logger)
	}

	// Initialize auth service
	authService := NewAuthService(userService, db, cfg, emailService)

	// Initialize auth handler
	authHandler := NewAuthHandler(authService)

	// Create API route group
	api := app.Group("/api/v1")

	// Public auth routes
	auth := api.Group("/auth")
	auth.Post("/register", sharedmiddleware.BodyValidator(&dto.RegisterRequest{}), authHandler.Register)
	auth.Post("/login", sharedmiddleware.BodyValidator(&dto.LoginRequest{}), authHandler.Login)
	auth.Post("/refresh", sharedmiddleware.BodyValidator(&dto.RefreshTokenRequest{}), authHandler.RefreshToken)
	auth.Post("/logout", sharedmiddleware.BodyValidator(&dto.RefreshTokenRequest{}), authHandler.Logout)
}
