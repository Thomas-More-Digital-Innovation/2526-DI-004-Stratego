package db

import (
	"context"
	"database/sql"
	"digital-innovation/stratego/utils"
	"fmt"
	"log"
	"time"

	"sync"

	_ "github.com/lib/pq"
)

var DB *sql.DB

type statsCache struct {
	userCount  int
	gameCount  int
	lastUpdate time.Time
	mu         sync.RWMutex
}

var cache = &statsCache{}

// InitDB initializes the database connection
func InitDB() error {
	dbHost := utils.GetEnv("DB_HOST", "localhost")
	dbPort := utils.GetEnv("DB_PORT", "5432")
	dbUser := utils.GetEnv("DB_USER", "stratego")
	dbPassword := utils.GetEnv("DB_PASSWORD", "pass")
	dbName := utils.GetEnv("DB_NAME", "stratego")
	sslMode := utils.GetEnv("DB_SSLMODE", "disable")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, sslMode,
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool limits
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(time.Hour)

	// Test the connection with retries
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		err = DB.Ping()
		if err == nil {
			log.Println("Database connection established")
			return nil
		}
		log.Printf("Failed to connect to database (attempt %d/%d): %v. Retrying in 2 seconds...", i+1, maxRetries, err)
		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("failed to ping database after %d attempts: %w", maxRetries, err)
}

// CloseDB closes the database connection
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

func updateStatsCache(ctx context.Context) error {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if time.Since(cache.lastUpdate) < time.Minute && !cache.lastUpdate.IsZero() {
		return nil
	}

	var userCount int
	err := DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&userCount)
	if err != nil {
		return err
	}

	var gameCount int
	// We aggregate the total games from user_stats as an indicator of platform activity
	err = DB.QueryRowContext(ctx, "SELECT COALESCE(SUM(total_games), 0) FROM user_stats").Scan(&gameCount)
	if err != nil {
		return err
	}

	cache.userCount = userCount
	cache.gameCount = gameCount
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
