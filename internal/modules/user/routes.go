package user

import (
	"go_boilerplate/internal/modules/role"
	"go_boilerplate/internal/modules/user/dto"
	"go_boilerplate/internal/shared/config"
	sharedmiddleware "go_boilerplate/internal/shared/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// RegisterRoutes registers all user-related routes
func RegisterRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config, logger *logrus.Logger) {
	// Initialize repositories
	userRepo := NewUserRepository(db)
	roleRepo := role.NewRoleRepository(db)

	// Initialize user service with role repository
	userService := NewUserServiceWithRole(userRepo, roleRepo)

	// Initialize handler
	userHandler := NewUserHandler(userService)

	// Create API route group
	api := app.Group("/api/v1")

	// Public routes (if any)
	// Currently, all user routes require authentication

	// Protected routes - All authenticated users
	protected := api.Group("/users")
	protected.Use(sharedmiddleware.JWTAuth(cfg))

	// Routes accessible by any authenticated user
	protected.Get("/me", userHandler.GetCurrentUser)                       // Get current user profile
	protected.Get("/:id", userHandler.GetUser)                             // Get user by ID
	protected.Put("/:id", sharedmiddleware.BodyValidator(&dto.UpdateUserRequest{}), userHandler.UpdateUser) // Update user (self-profile or with permission)

	// Routes accessible by Admin and SuperAdmin only
	adminOnly := protected.Group("/")
	adminOnly.Use(sharedmiddleware.RequireRole(cfg, "admin", "super_admin"))
	adminOnly.Get("/", userHandler.GetUsers)                               // Get all users (with pagination)
	adminOnly.Post("/", sharedmiddleware.BodyValidator(&dto.CreateUserRequest{}), userHandler.CreateUser) // Create user
	adminOnly.Delete("/:id", userHandler.DeleteUser)                       // Delete user

	// Routes accessible by SuperAdmin only
	superAdminOnly := protected.Group("/")
	superAdminOnly.Use(sharedmiddleware.RequireRole(cfg, "super_admin"))
	superAdminOnly.Patch("/:id/role", sharedmiddleware.BodyValidator(&dto.AssignRoleRequest{}), userHandler.AssignRole) // Assign role to user
}
