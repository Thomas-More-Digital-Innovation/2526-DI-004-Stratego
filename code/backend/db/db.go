package db

import (
	"context"
	"digital-innovation/stratego/logging"
	"digital-innovation/stratego/utils"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

type contextKey string

const UserIDContextKey contextKey = "user_id"

// WithUserID returns a new context with the given user ID for DB operations
func WithUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, UserIDContextKey, userID)
}

// WithRLS executes fn inside a transaction where app.current_user_id is set from ctx.
// SET LOCAL ensures the variable is scoped to the transaction, and the transaction
// guarantees both SET and the query run on the exact same connection from the pool.
func WithRLS(ctx context.Context, fn func(tx *gorm.DB) error) error {
	if !isPostgresDialect() {
		return DB.WithContext(ctx).Transaction(fn)
	}
	userID := 0
	if id, ok := ctx.Value(UserIDContextKey).(int); ok {
		userID = id
	}
	return DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(fmt.Sprintf("SET LOCAL app.current_user_id = '%d'", userID)).Error; err != nil {
			return fmt.Errorf("failed to set RLS context: %w", err)
		}
		return fn(tx)
	})
}

// InitDB initializes the database connection
func InitDB() error {
	dbHost := utils.GetEnv("DB_HOST", "localhost")
	dbPort := utils.GetEnv("DB_PORT", "5432")
	dbUser := utils.GetEnv("DB_USER", "stratego")
	dbPassword := utils.GetEnv("DB_PASSWORD", "pass")
	dbName := utils.GetEnv("DB_NAME", "stratego")
	sslMode := utils.GetEnv("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		dbHost, dbUser, dbPassword, dbName, dbPort, sslMode,
	)

	var err error

	// Test the connection with retries
	maxRetries := 10
	for i := range maxRetries {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

		if err == nil {
			sqlDB, err := DB.DB()
			if err == nil {
				// Set connection pool limits
				sqlDB.SetMaxOpenConns(25)
				sqlDB.SetMaxIdleConns(5)
				sqlDB.SetConnMaxLifetime(time.Hour)

				err = sqlDB.Ping()
				if err == nil {
					logging.Debug(logging.TagAuth, "Database connection established")

					// Run migrations automatically on startup
					if err := RunMigrations(context.Background()); err != nil {
						return fmt.Errorf("failed to run migrations: %w", err)
					}

					return nil
				}
			}
		}
		logging.Error(fmt.Sprintf("Failed to connect to database (attempt %d/%d)", i+1, maxRetries), err)
		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("failed to ping database after %d attempts: %w", maxRetries, err)
}

// CloseDB closes the database connection
func CloseDB() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
