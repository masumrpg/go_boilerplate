package auth

import (
	"go_boilerplate/internal/modules/auth/dto"
	"go_boilerplate/internal/modules/email"
	"go_boilerplate/internal/modules/role"
	"go_boilerplate/internal/modules/user"
	"go_boilerplate/internal/shared/config"
	sharedmiddleware "go_boilerplate/internal/shared/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// RegisterRoutes registers all auth-related routes
func RegisterRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config, logger *logrus.Logger, redisClient *redis.Client) {
	// Initialize repositories
	userRepo := user.NewUserRepository(db)
	roleRepo := role.NewRoleRepository(db)

	// Initialize user service with role repository
	userService := user.NewUserServiceWithRole(userRepo, roleRepo)

	// Initialize email service (optional, will check before sending)
	var emailService email.EmailService
	if cfg.Email.Enabled {
		emailService = email.NewEmailService(cfg, logger)
	}

	// Initialize auth service
	authService := NewAuthService(userService, db, cfg, emailService, redisClient)

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

	// Add new verification endpoints
	auth.Post("/verify-email", sharedmiddleware.BodyValidator(&dto.VerifyEmailRequest{}), authHandler.VerifyEmail)
	auth.Post("/verify-2fa", sharedmiddleware.BodyValidator(&dto.Verify2FARequest{}), authHandler.Verify2FA)
	auth.Post("/resend-verification", sharedmiddleware.BodyValidator(&dto.ResendCodeRequest{}), authHandler.ResendVerification)
	auth.Post("/resend-2fa", sharedmiddleware.BodyValidator(&dto.ResendCodeRequest{}), authHandler.Resend2FA)

	// Protected session management routes
	sessions := auth.Group("/sessions", sharedmiddleware.JWTAuth(cfg))
	sessions.Get("/", authHandler.GetSessions)
	sessions.Delete("/:id", authHandler.DeleteSession)
	sessions.Patch("/:id/block", authHandler.BlockSession)
}
