package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGlobalStats(t *testing.T) {
	SetupDBTest(t)
	ctx := context.Background()

	t.Run("Count Users", func(t *testing.T) {
		// Reset cache for test predictability
		cache.mu.Lock()
		cache.lastUpdate = time.Time{}
		cache.mu.Unlock()

		_, _ = CreateUser(ctx, "user1", "Pass1!", "")
		_, _ = CreateUser(ctx, "user2", "Pass1!", "")

		count, err := GetTotalUserCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 2, count)
	})

	t.Run("Count Games Played", func(t *testing.T) {
		// Reset cache
		cache.mu.Lock()
		cache.lastUpdate = time.Time{}
		cache.mu.Unlock()

		user, _ := CreateUser(ctx, "statuser", "Pass1!", "")

		_ = UpdateUserStats(ctx, user.ID, true, 5, 60.0)
		_ = UpdateUserStats(ctx, user.ID, false, 8, 90.0)

		count, err := GetTotalGamesPlayedCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 2, count)
	})

	t.Run("Cache respects TTL", func(t *testing.T) {
		// Set cache to something old
		cache.mu.Lock()
		cache.userCount = 999
		cache.lastUpdate = time.Now()
		cache.mu.Unlock()

		// This should return the cached value 999 because TTL (1 min) hasn't passed
		count, _ := GetTotalUserCount(ctx)
		assert.Equal(t, 999, count)

		// Force expiry
		cache.mu.Lock()
		cache.lastUpdate = time.Now().Add(-2 * time.Minute)
		cache.mu.Unlock()

		// This should trigger a refresh
		count, _ = GetTotalUserCount(ctx)
		assert.NotEqual(t, 999, count)
	})

	t.Run("Stress concurrency", func(_ *testing.T) {
		// Reset cache
		cache.mu.Lock()
		cache.lastUpdate = time.Time{}
		cache.mu.Unlock()

		const goroutines = 50
		done := make(chan bool)

		for range goroutines {
			go func() {
				_, _ = GetTotalUserCount(ctx)
				done <- true
			}()
		}

		for range goroutines {
			<-done
		}
		// If we reached here without deadlocking, it's a pass
	})
}
