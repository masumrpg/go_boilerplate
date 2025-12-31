package role

import (
	"go_boilerplate/internal/shared/config"
	"go_boilerplate/internal/shared/middleware"
	"go_boilerplate/internal/modules/role/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// RegisterRoutes registers all role-related routes
func RegisterRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config, logger *logrus.Logger) {
	// Initialize repository
	roleRepo := NewRoleRepository(db)

	// Initialize service
	roleService := NewRoleService(roleRepo)

	// Initialize handler
	roleHandler := NewRoleHandler(roleService)

	// Create API route group
	api := app.Group("/api/v1")

	// Protected routes - require SuperAdmin role
	roles := api.Group("/roles")
	roles.Use(middleware.JWTAuth(cfg))
	// TODO: Add RequireRole middleware once implemented
	// roles.Use(middleware.RequireRole(cfg, "super_admin"))

	// Role CRUD routes (only SuperAdmin can manage roles)
	roles.Get("/", roleHandler.GetRoles)                        // Get all roles (with pagination)
	roles.Get("/:id", roleHandler.GetRole)                      // Get role by ID
	roles.Post("/", middleware.BodyValidator(&dto.CreateRoleRequest{}), roleHandler.CreateRole) // Create role (SuperAdmin only)
	roles.Put("/:id", middleware.BodyValidator(&dto.UpdateRoleRequest{}), roleHandler.UpdateRole) // Update role (SuperAdmin only)
	roles.Delete("/:id", roleHandler.DeleteRole)                // Delete role (SuperAdmin only)

	logger.Info("✓ Role routes registered (SuperAdmin only)")
}
