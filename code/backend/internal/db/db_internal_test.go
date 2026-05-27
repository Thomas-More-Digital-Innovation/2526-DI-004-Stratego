package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestDB_WithUserID(t *testing.T) {
	ctx := context.Background()
	userID := 123

	newCtx := WithUserID(ctx, userID)

	val := newCtx.Value(UserIDContextKey)
	assert.Equal(t, userID, val)
}

func TestDB_WithRLS_Fallback(t *testing.T) {
	SetupTestDB(t) // Uses SQLite

	ctx := context.Background()
	ctx = WithUserID(ctx, 123)

	err := WithRLS(ctx, func(tx *gorm.DB) error {
		// Just verify we can execute something
		return tx.Exec("SELECT 1").Error
	})

	assert.NoError(t, err)
}

func TestDB_WithRLS_Error(t *testing.T) {
	SetupTestDB(t)

	// Missing user ID in context
	ctx := context.Background()

	err := WithRLS(ctx, func(_ *gorm.DB) error {
		return nil
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "RLS context missing")
}

func TestDB_CloseDB(t *testing.T) {
	SetupTestDB(t)

	err := CloseDB()
	assert.NoError(t, err)

	// Subsequent calls should be fine
	err = CloseDB()
	assert.NoError(t, err)
}
