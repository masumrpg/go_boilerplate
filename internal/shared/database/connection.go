package database

import (
	"fmt"
	"time"

	"go_boilerplate/internal/shared/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB initializes the database connection
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	// Configure GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(getLogLevel(cfg)),
		// Disable foreign key constraints during development if needed
		// DisableForeignKeyConstraintWhenMigrating: true,
	}

	// Open connection
	db, err := gorm.Open(postgres.Open(cfg.Database.GetDSN()), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL DB instance to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
	setConnectionPoolSettings(sqlDB)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// getLogLevel returns the appropriate GORM log level based on server mode
func getLogLevel(cfg *config.Config) logger.LogLevel {
	if cfg.Server.IsDevelopment() {
		return logger.Info
	}
	return logger.Silent
}

// setConnectionPoolSettings configures the database connection pool
func setConnectionPoolSettings(sqlDB interface{}) {
	// Type assertion to access *sql.DB methods
	type db interface {
		SetMaxIdleConns(n int)
		SetMaxOpenConns(n int)
		SetConnMaxLifetime(d time.Duration)
		SetConnMaxIdleTime(d time.Duration)
	}

	if db, ok := sqlDB.(db); ok {
		// SetMaxIdleConns sets the maximum number of connections in the idle connection pool
		db.SetMaxIdleConns(10)

		// SetMaxOpenConns sets the maximum number of open connections to the database
		db.SetMaxOpenConns(100)

		// SetConnMaxLifetime sets the maximum amount of time a connection may be reused
		db.SetConnMaxLifetime(1 * time.Hour)

		// SetConnMaxIdleTime sets the maximum amount of time a connection may be idle
		db.SetConnMaxIdleTime(10 * time.Minute)
	}
}

// CloseDB closes the database connection
func CloseDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	return nil
}
