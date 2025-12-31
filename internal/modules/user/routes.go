package user

import (
	"go_boilerplate/internal/shared/config"
	"go_boilerplate/internal/shared/middleware"
	"go_boilerplate/internal/modules/user/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// RegisterRoutes registers all user-related routes
func RegisterRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config, logger *logrus.Logger) {
	// Initialize repository
	userRepo := NewUserRepository(db)

	// Initialize service
	userService := NewUserService(userRepo)

	// Initialize handler
	userHandler := NewUserHandler(userService)

	// Create API route group
	api := app.Group("/api/v1")

	// Public routes (if any)
	// Currently, all user routes require authentication

	// Protected routes
	protected := api.Group("/users")
	protected.Use(middleware.JWTAuth(cfg))

	// User CRUD routes
	protected.Get("/", userHandler.GetUsers)                    // Get all users (with pagination)
	protected.Get("/me", userHandler.GetCurrentUser)            // Get current user profile
	protected.Get("/:id", userHandler.GetUser)                  // Get user by ID
	protected.Post("/", middleware.BodyValidator(&dto.CreateUserRequest{}), userHandler.CreateUser) // Create user (admin only in production)
	protected.Put("/:id", middleware.BodyValidator(&dto.UpdateUserRequest{}), userHandler.UpdateUser) // Update user
	protected.Delete("/:id", userHandler.DeleteUser)            // Delete user (admin only in production)
}
