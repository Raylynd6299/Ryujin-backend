package persistence

import (
	"embed"
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"gorm.io/gorm"
)

// RunMigrations applies all pending up-migrations from the embedded SQL files.
// It is safe to call on every startup — already-applied migrations are skipped.
// Returns an error only if something unexpected happens (not for "no change").
func RunMigrations(db *gorm.DB, fs embed.FS, dir string) error {
	// Unwrap the underlying *sql.DB from GORM
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("migrations: failed to get sql.DB from gorm: %w", err)
	}

	// Source: read SQL files from the embedded FS
	src, err := iofs.New(fs, dir)
	if err != nil {
		return fmt.Errorf("migrations: failed to create iofs source: %w", err)
	}

	// Driver: postgres adapter for golang-migrate
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("migrations: failed to create postgres driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", src, "postgres", driver)
	if err != nil {
		return fmt.Errorf("migrations: failed to create migrator: %w", err)
	}

	if err := m.Up(); err != nil {
		// ErrNoChange means the DB is already up to date — not an error.
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("✓ Migrations: already up to date")
			return nil
		}
		return fmt.Errorf("migrations: failed to run: %w", err)
	}

	version, _, _ := m.Version()
	log.Printf("✓ Migrations: applied successfully (version %d)", version)
	return nil
}
