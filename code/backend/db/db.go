package db

import (
	"context"
	"digital-innovation/stratego/logging"
	"digital-innovation/stratego/utils"
	"fmt"
	"time"

	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

type contextKey string

const UserIDContextKey contextKey = "user_id"

type statsCache struct {
	userCount  int
	gameCount  int
	lastUpdate time.Time
	mu         sync.RWMutex
}

var cache = &statsCache{}

// WithUserID returns a new context with the given user ID for DB operations
func WithUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, UserIDContextKey, userID)
}

func rlsPlugin(db *gorm.DB) {
	db.Callback().Query().Before("gorm:query").Register("rls:set_user_id", setUserIDCallback)
	db.Callback().Create().Before("gorm:create").Register("rls:set_user_id", setUserIDCallback)
	db.Callback().Update().Before("gorm:update").Register("rls:set_user_id", setUserIDCallback)
	db.Callback().Delete().Before("gorm:delete").Register("rls:set_user_id", setUserIDCallback)
}

func setUserIDCallback(db *gorm.DB) {
	if db.Dialector.Name() != "postgres" {
		return
	}
	userID := 0
	if db.Statement.Context != nil {
		if id, ok := db.Statement.Context.Value(UserIDContextKey).(int); ok {
			userID = id
		}
	}
	// Use SET instead of SET LOCAL to ensure it persists for the session if not in a transaction.
	// Since we set it for EVERY query (via the callback), it will be overwritten correctly
	// even if connections are reused from the pool.
	db.Exec(fmt.Sprintf("SET app.current_user_id = '%d'", userID))
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
			rlsPlugin(DB)
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

func updateStatsCache(ctx context.Context) error {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if time.Since(cache.lastUpdate) < time.Minute && !cache.lastUpdate.IsZero() {
		return nil
	}

	var userCount int64
	err := DB.WithContext(ctx).Table("users").Count(&userCount).Error
	if err != nil {
		return err
	}

	var gameCount int64
	err = DB.WithContext(ctx).Table("user_stats").Select("COALESCE(SUM(total_games), 0)").Scan(&gameCount).Error
	if err != nil {
		return err
	}

	cache.userCount = int(userCount)
	cache.gameCount = int(gameCount)
	cache.lastUpdate = time.Now()
	return nil
}

func GetTotalUserCount(ctx context.Context) (int, error) {
	if err := updateStatsCache(ctx); err != nil {
		return 0, err
	}
	cache.mu.RLock()
	defer cache.mu.RUnlock()
	return cache.userCount, nil
}

func GetTotalGamesPlayedCount(ctx context.Context) (int, error) {
	if err := updateStatsCache(ctx); err != nil {
		return 0, err
	}
	cache.mu.RLock()
	defer cache.mu.RUnlock()
	return cache.gameCount, nil
}
