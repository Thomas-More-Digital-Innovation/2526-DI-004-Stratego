package db

import (
	"context"
	"testing"
	"time"
)

func TestGlobalStats(t *testing.T) {
	testDB := setupSQLiteDB(t)
	oldDB := DB
	DB = testDB
	defer func() { DB = oldDB }()

	ctx := context.Background()

	t.Run("Count Users", func(t *testing.T) {
		// Reset cache for test predictability
		cache.mu.Lock()
		cache.lastUpdate = time.Time{}
		cache.mu.Unlock()

		_, _ = CreateUser(ctx, "user1", "Pass1!", "")
		_, _ = CreateUser(ctx, "user2", "Pass1!", "")

		count, err := GetTotalUserCount(ctx)
		if err != nil {
			t.Fatalf("GetTotalUserCount failed: %v", err)
		}
		if count != 2 {
			t.Errorf("expected 2 users, got %d", count)
		}
	})

	t.Run("Count Games Played", func(t *testing.T) {
		// Reset cache
		cache.mu.Lock()
		cache.lastUpdate = time.Time{}
		cache.mu.Unlock()

		user, _ := CreateUser(ctx, "statuser", "Pass1!", "")

		// Add some games via stats update
		_ = UpdateUserStats(ctx, user.ID, true, 5, 60.0)
		_ = UpdateUserStats(ctx, user.ID, false, 8, 90.0)

		count, err := GetTotalGamesPlayedCount(ctx)
		if err != nil {
			t.Fatalf("GetTotalGamesPlayedCount failed: %v", err)
		}
		if count != 2 {
			t.Errorf("expected 2 games played, got %d", count)
		}
	})

	t.Run("Cache respects TTL", func(t *testing.T) {
		// Set cache to something old
		cache.mu.Lock()
		cache.userCount = 999
		cache.lastUpdate = time.Now()
		cache.mu.Unlock()

		// This should return the cached value 999 because TTL (1 min) hasn't passed
		count, _ := GetTotalUserCount(ctx)
		if count != 999 {
			t.Errorf("expected cached value 999, got %d", count)
		}

		// Force expiry
		cache.mu.Lock()
		cache.lastUpdate = time.Now().Add(-2 * time.Minute)
		cache.mu.Unlock()

		// This should trigger a refresh
		count, _ = GetTotalUserCount(ctx)
		if count == 999 {
			t.Error("expected fresh value after cache expiry, but got cached value")
		}
	})

	t.Run("Stress concurrency", func(_ *testing.T) {
		// Reset cache
		cache.mu.Lock()
		cache.lastUpdate = time.Time{}
		cache.mu.Unlock()

		const goroutines = 50
		done := make(chan bool)

		for i := 0; i < goroutines; i++ {
			go func() {
				_, _ = GetTotalUserCount(ctx)
				done <- true
			}()
		}

		for i := 0; i < goroutines; i++ {
			<-done
		}
		// If we reached here without deadlocking, it's a pass
	})
}
