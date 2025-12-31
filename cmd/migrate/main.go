package main

import (
	"flag"
	"log"

	"go_boilerplate/internal/shared/config"
	"go_boilerplate/internal/shared/database"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// Define flags
	up := flag.Bool("up", false, "Run up migrations")
	down := flag.Bool("down", false, "Run down migrations")
	steps := flag.Int("steps", 0, "Number of steps to migrate (0 for all)")
	version := flag.Bool("version", false, "Print current migration version")
	force := flag.Int("force", -1, "Force set version (useful for dirty state)")

	flag.Parse()

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database using existing helper or construct DSN manually
	// We'll use the DSN from config directly

	// Initialize database connection for driver
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB: %v", err)
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{
		MigrationsTable: "schema_migrations", // Default table name
	})
	if err != nil {
		log.Fatalf("Failed to create postgres driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}

	// Handle force version
	if *force >= 0 {
		if err := m.Force(*force); err != nil {
			log.Fatalf("Failed to force version: %v", err)
		}
		log.Printf("Forced version to %d", *force)
		return
	}

	// Handle version check
	if *version {
		v, dirty, err := m.Version()
		if err != nil && err != migrate.ErrNilVersion {
			log.Fatalf("Failed to get version: %v", err)
		}
		if err == migrate.ErrNilVersion {
			log.Println("No migrations applied")
		} else {
			log.Printf("Version: %d, Dirty: %v\n", v, dirty)
		}
		return
	}

	// Handle Up migration
	if *up {
		if *steps > 0 {
			if err := m.Steps(*steps); err != nil {
				if err == migrate.ErrNoChange {
					log.Println("No changes to apply")
				} else {
					log.Fatalf("Failed to migrate up %d steps: %v", *steps, err)
				}
			} else {
				log.Printf("Migrated up %d steps successfully", *steps)
			}
		} else {
			if err := m.Up(); err != nil {
				if err == migrate.ErrNoChange {
					log.Println("No changes to apply")
				} else {
					log.Fatalf("Failed to run up migrations: %v", err)
				}
			} else {
				log.Println("Migrated up successfully")
			}
		}
		return
	}

	// Handle Down migration
	if *down {
		if *steps > 0 {
			if err := m.Steps(-(*steps)); err != nil {
				if err == migrate.ErrNoChange {
					log.Println("No changes to revert")
				} else {
					log.Fatalf("Failed to migrate down %d steps: %v", *steps, err)
				}
			} else {
				log.Printf("Migrated down %d steps successfully", *steps)
			}
		} else {
			if err := m.Down(); err != nil {
				if err == migrate.ErrNoChange {
					log.Println("No changes to revert")
				} else {
					log.Fatalf("Failed to run down migrations: %v", err)
				}
			} else {
				log.Println("Migrated down successfully")
			}
		}
		return
	}

	// If no flags are set, print usage
	flag.Usage()
}
