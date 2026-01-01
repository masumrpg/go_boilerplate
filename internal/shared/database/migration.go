package database

import (
	"fmt"

	roleModule "go_boilerplate/internal/modules/role"
	userModule "go_boilerplate/internal/modules/user"
	"go_boilerplate/internal/shared/config"
	"go_boilerplate/internal/shared/utils"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// AutoMigrate runs auto migration for given models
func AutoMigrate(db *gorm.DB, models []any, logger *logrus.Logger) error {
	logger.Info("Running database migrations...")

	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("failed to run auto migration: %w", err)
	}

	logger.Info("Database migrations completed successfully")
	return nil
}

// DropAllTables drops all tables in the database
// WARNING: This will delete all data. Use only in testing/development
func DropAllTables(db *gorm.DB, logger *logrus.Logger) error {
	logger.Warn("Dropping all tables...")

	if err := db.Migrator().DropTable(
		// Add all table names here
		// "users",
		// "refresh_tokens",
	); err != nil {
		return fmt.Errorf("failed to drop tables: %w", err)
	}

	logger.Warn("All tables dropped successfully")
	return nil
}

// RenameTables drops old tables so AutoMigrate can create new ones with prefixes
// This is a destructive operation and should only be used in development
// Old tables to drop: users, oauth_accounts, refresh_tokens
// New tables will be created with prefixes: m_users, m_roles, t_oauth_accounts, t_refresh_tokens
func RenameTables(db *gorm.DB, logger *logrus.Logger) error {
	logger.Info("Starting table rename migration...")

	// List of old tables to drop
	oldTables := []string{
		"users",
		"oauth_accounts",
		"refresh_tokens",
		"t_refresh_tokens",
	}

	// Drop old tables if they exist
	for _, table := range oldTables {
		logger.Infof("Dropping old table: %s", table)

		// Check if table exists before dropping
		if db.Migrator().HasTable(table) {
			if err := db.Migrator().DropTable(table); err != nil {
				logger.Errorf("Failed to drop table %s: %v", table, err)
				return err
			}
			logger.Infof("Successfully dropped table: %s", table)
		} else {
			logger.Infof("Table %s does not exist, skipping", table)
		}
	}

	logger.Info("Table rename migration completed successfully")
	return nil
}

// CreateIndexes creates indexes for optimized queries
func CreateIndexes(db *gorm.DB, logger *logrus.Logger) error {
	logger.Info("Creating database indexes...")

	// Example: Create index on users.email
	// if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)").Error; err != nil {
	// 	return fmt.Errorf("failed to create index on users.email: %w", err)
	// }

	logger.Info("Database indexes created successfully")
	return nil
}

// SeedDatabase seeds the database with initial data
func SeedDatabase(db *gorm.DB, cfg *config.Config, logger *logrus.Logger) error {
	logger.Info("Seeding database...")

	// Add your seed data here
	// Example: Create default admin user
	// adminUser := &user.User{
	// 	Name:  "Admin",
	// 	Email: "admin@example.com",
	// }
	// if err := db.FirstOrCreate(adminUser, user.User{Email: adminUser.Email}).Error; err != nil {
	// 	return fmt.Errorf("failed to seed admin user: %w", err)
	// }

	logger.Info("Database seeded successfully")
	return nil
}

// SeedSuperAdmin creates a default SuperAdmin user if it doesn't exist
// This should be called AFTER roles are seeded
func SeedSuperAdmin(db *gorm.DB, cfg *config.Config, logger *logrus.Logger) error {
	logger.Info("Checking for SuperAdmin user...")

	// Get SuperAdmin role
	var superAdminRole roleModule.Role
	if err := db.Where("slug = ?", "super_admin").First(&superAdminRole).Error; err != nil {
		return fmt.Errorf("super_admin role not found: %w", err)
	}

	// Check if SuperAdmin user already exists by email
	var existingUser userModule.User
	if err := db.Where("email = ?", cfg.SuperAdmin.Email).First(&existingUser).Error; err == nil {
		// User exists, update password and role
		hashedPassword, err := utils.HashPassword(cfg.SuperAdmin.Password)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		existingUser.Password = hashedPassword
		existingUser.Name = cfg.SuperAdmin.Name
		existingUser.RoleID = superAdminRole.ID
		existingUser.IsVerified = true

		if err := db.Save(&existingUser).Error; err != nil {
			return fmt.Errorf("failed to update superadmin user: %w", err)
		}

		logger.Info("✓ SuperAdmin user updated successfully")
		logger.Infof("  Email: %s", cfg.SuperAdmin.Email)
		logger.Warn("  ⚠️  Please change the default password after first login!")
		return nil
	}

	// Hash password from config
	hashedPassword, err := utils.HashPassword(cfg.SuperAdmin.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Create SuperAdmin user
	superAdminUser := &userModule.User{
		ID:       uuid.New(),
		Name:     cfg.SuperAdmin.Name,
		Email:    cfg.SuperAdmin.Email,
		Password: hashedPassword,
		RoleID:   superAdminRole.ID,
		IsVerified: true,
	}

	if err := db.Create(superAdminUser).Error; err != nil {
		return fmt.Errorf("failed to create superadmin user: %w", err)
	}

	logger.Info("✓ SuperAdmin user created successfully")
	logger.Infof("  Email: %s", cfg.SuperAdmin.Email)
	logger.Warn("  ⚠️  Please change the default password after first login!")

	return nil
}
