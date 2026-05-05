package db

import (
	"context"
	"sync"
	"time"
)

type statsCache struct {
	userCount  int
	gameCount  int
	lastUpdate time.Time
	mu         sync.RWMutex
}

var cache = &statsCache{}

func updateStatsCache(ctx context.Context) error {
	// Fast path: check with read lock
	cache.mu.RLock()
	needsUpdate := time.Since(cache.lastUpdate) >= time.Minute || cache.lastUpdate.IsZero()
	cache.mu.RUnlock()

	if !needsUpdate {
		return nil
	}

	// Slow path: acquire write lock
	cache.mu.Lock()
	defer cache.mu.Unlock()

	// Check again in case someone else updated it while we were waiting for the lock
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

// GetTotalUserCount returns the cached or fresh count of total registered users
func GetTotalUserCount(ctx context.Context) (int, error) {
	if err := updateStatsCache(ctx); err != nil {
		return 0, err
	}
	cache.mu.RLock()
	defer cache.mu.RUnlock()
	return cache.userCount, nil
}

// GetTotalGamesPlayedCount returns the cached or fresh count of total games played across all users
func GetTotalGamesPlayedCount(ctx context.Context) (int, error) {
	if err := updateStatsCache(ctx); err != nil {
		return 0, err
	}
	cache.mu.RLock()
	defer cache.mu.RUnlock()
	return cache.gameCount, nil
}
