package database

import (
	"fmt"

	"go_boilerplate/internal/shared/config"

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
