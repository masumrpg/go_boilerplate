package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	authModule "go_boilerplate/internal/modules/auth"
	"go_boilerplate/internal/modules/auth/dto"
	oauthModule "go_boilerplate/internal/modules/oauth"
	oauthdto "go_boilerplate/internal/modules/oauth/dto"
	roleModule "go_boilerplate/internal/modules/role"
	userModule "go_boilerplate/internal/modules/user"
	"go_boilerplate/internal/shared/config"
	"go_boilerplate/internal/shared/database"
	"go_boilerplate/internal/shared/middleware"
	"go_boilerplate/internal/shared/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sirupsen/logrus"
)

func main() {
	// 1. Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// 2. Initialize logger
	logger := utils.InitLogger(cfg)
	logger.Info("Starting Go Boilerplate API...")

	// 3. Initialize database
	db, err := database.InitDB(cfg)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	logger.Info("Database connected successfully")

	// 4. Initialize Redis
	redisClient, err := database.InitRedis(cfg, logger)
	if err != nil {
		logger.Warnf("Failed to connect to Redis: %v", err)
	} else {
		defer redisClient.Close()
	}

	// 5. Run database migrations

	// Step 1: Rename tables (drop old tables) - ONLY IN DEVELOPMENT
	if cfg.Server.IsDevelopment() {
		logger.Info("Running in development mode - dropping old tables...")
		if err := database.RenameTables(db, logger); err != nil {
			logger.Warnf("Failed to rename tables: %v", err)
		}
	}

	// Step 2: AutoMigrate models with new table names
	// This should only run in development. In production, use manual migrations (golang-migrate).
	if cfg.Server.IsDevelopment() {
		migrationModels := []any{
			&roleModule.Role{},
			&userModule.User{},
			&dto.RefreshToken{},
			&oauthdto.OAuthAccount{},
		}

		if err := database.AutoMigrate(db, migrationModels, logger); err != nil {
			logger.Fatalf("Failed to run migrations: %v", err)
		}
	} else {
		logger.Info("Running in production mode - skipping AutoMigrate")
	}

	// Step 3: Seed initial roles
	roleRepo := roleModule.NewRoleRepository(db)
	roleService := roleModule.NewRoleService(roleRepo)
	if err := roleService.SeedInitialRoles(); err != nil {
		logger.Warnf("Failed to seed initial roles: %v", err)
	} else {
		logger.Info("✓ Initial roles seeded successfully")
	}

	// Step 4: Seed SuperAdmin user
	if err := database.SeedSuperAdmin(db, cfg, logger); err != nil {
		logger.Warnf("Failed to seed SuperAdmin user: %v", err)
	}

	// 5. Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:               "Go Boilerplate API",
		DisableStartupMessage: false,
		EnablePrintRoutes:     cfg.Server.IsDevelopment(),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			// Log error
			logger.WithFields(logrus.Fields{
				"path":    c.Path(),
				"method":  c.Method(),
				"status":  code,
				"error":   err.Error(),
			}).Error("Request error")

			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		},
	})

	// 6. Register global middleware
	app.Use(middleware.HTTPLogger(logger))
	app.Use(middleware.CORS(cfg))
	app.Use(recover.New())

	// 7. Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"status":  "ok",
			"message": "API is running",
		})
	})

	// 8. Register module routes
	logger.Info("Registering module routes...")

	// Auth routes (register, login, refresh, logout)
	authModule.RegisterRoutes(app, db, cfg, logger, redisClient)
	logger.Info("✓ Auth routes registered")

	// User routes (CRUD operations)
	userModule.RegisterRoutes(app, db, cfg, logger)
	logger.Info("✓ User routes registered")

	// Role routes (manage roles - SuperAdmin only)
	roleModule.RegisterRoutes(app, db, cfg, logger)
	logger.Info("✓ Role routes registered")

	// OAuth routes (Google, GitHub)
	oauthModule.RegisterRoutes(app, db, cfg, logger)
	logger.Info("✓ OAuth routes registered")

	// 9. Graceful shutdown
	// Handle shutdown signals
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		logger.Info("Shutting down server...")

		if err := app.Shutdown(); err != nil {
			logger.Errorf("Error during server shutdown: %v", err)
		}

		// Close database connection
		if err := database.CloseDB(db); err != nil {
			logger.Errorf("Error closing database: %v", err)
		}

		logger.Info("Server shut down gracefully")
	}()

	// 10. Start server
	addr := ":" + cfg.Server.Port
	logger.Infof("Server starting on %s", addr)
	logger.Infof("Environment: %s", cfg.Server.Mode)
	logger.Infof("API Documentation: http://localhost%s/swagger", addr)

	if err := app.Listen(addr); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
