package db

import (
	"context"
	"digital-innovation/stratego/logging"
	"digital-innovation/stratego/models"
	"embed"
	"fmt"
	"sort"
	"strings"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// RunMigrations applies all pending SQL migrations to the database
func RunMigrations(ctx context.Context) error {
	// 0. Pre-migration: Drop old constraints that might conflict with GORM
	// The manual schema.sql used 'UNIQUE', which Postgres names 'users_username_key'
	// GORM sometimes gets confused if these exist under different names or if it fails a previous migration
	_ = DB.WithContext(ctx).Exec("ALTER TABLE users DROP CONSTRAINT IF EXISTS users_username_key").Error
	_ = DB.WithContext(ctx).Exec("ALTER TABLE users DROP CONSTRAINT IF EXISTS uni_users_username").Error
	_ = DB.WithContext(ctx).Exec("ALTER TABLE user_stats DROP CONSTRAINT IF EXISTS user_stats_user_id_key").Error
	_ = DB.WithContext(ctx).Exec("ALTER TABLE user_stats DROP CONSTRAINT IF EXISTS uni_user_stats_user_id").Error

	// 1. AutoMigrate GORM models
	err := DB.WithContext(ctx).AutoMigrate(
		&models.User{},
		&models.UserStats{},
		&models.BoardSetup{},
		&models.RefreshToken{},
		&models.Game{},
		&models.GameMove{},
	)
	if err != nil {
		return fmt.Errorf("failed to auto-migrate models: %w", err)
	}

	// 2. Run manual SQL migrations (for RLS, etc.)
	err = DB.WithContext(ctx).Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMPTZ DEFAULT now()
		)
	`).Error
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			files = append(files, entry.Name())
		}
	}
	sort.Strings(files)

	for _, file := range files {
		var exists bool
		err := DB.WithContext(ctx).Raw("SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = ?)", file).Scan(&exists).Error
		if err != nil {
			return fmt.Errorf("failed to check migration status for %s: %w", file, err)
		}

		if exists {
			continue
		}

		logging.Debug(logging.TagAuth, "Applying migration: %s", file)
		content, err := migrationsFS.ReadFile("migrations/" + file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		// Execute the migration SQL
		err = DB.WithContext(ctx).Exec(string(content)).Error
		if err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", file, err)
		}

		// Record that this migration has been applied
		err = DB.WithContext(ctx).Exec("INSERT INTO schema_migrations (version) VALUES (?)", file).Error
		if err != nil {
			return fmt.Errorf("failed to record migration %s: %w", file, err)
		}
	}

	return nil
}
