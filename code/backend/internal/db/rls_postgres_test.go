package db

import (
	"context"
	"digital-innovation/gostrategy/internal/models"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupPostgresDB(t *testing.T) *gorm.DB {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		t.Skip("Skipping Postgres RLS test: TEST_DB_DSN not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	oldDB := DB
	DB = db
	err = RunMigrations(context.Background())
	if err != nil {
		DB = oldDB
	}
	require.NoError(t, err)
	DB = oldDB

	return db
}

func TestRowLevelSecurity(t *testing.T) {
	testDB := setupPostgresDB(t)
	oldDB := DB
	DB = testDB
	defer func() { DB = oldDB }()

	ctx := context.Background()

	userA, _ := CreateUser(ctx, "user_a", "Pass1234!", "")
	userB, _ := CreateUser(ctx, "user_b", "Pass1234!", "")

	// User A creates a setup via WithRLS so the INSERT also respects the policy
	var setupA *models.BoardSetup
	err := WithRLS(WithUserID(ctx, userA.ID), func(tx *gorm.DB) error {
		setup := models.BoardSetup{
			UserID:    userA.ID,
			Name:      "A's Setup",
			SetupData: "DATA_A",
		}
		if err := tx.Create(&setup).Error; err != nil {
			return err
		}
		setupA = &setup
		return nil
	})
	require.NoError(t, err)

	t.Run("User B cannot see User A's setup via RLS", func(t *testing.T) {
		var retrieved models.BoardSetup
		err := WithRLS(WithUserID(ctx, userB.ID), func(tx *gorm.DB) error {
			return tx.First(&retrieved, setupA.ID).Error
		})

		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})

	t.Run("User A can see their own setup", func(t *testing.T) {
		var retrieved models.BoardSetup
		err := WithRLS(WithUserID(ctx, userA.ID), func(tx *gorm.DB) error {
			return tx.First(&retrieved, setupA.ID).Error
		})

		assert.NoError(t, err)
		assert.Equal(t, setupA.ID, retrieved.ID)
	})

	t.Run("User B cannot delete User A's setup", func(t *testing.T) {
		var rowsAffected int64
		_ = WithRLS(WithUserID(ctx, userB.ID), func(tx *gorm.DB) error {
			result := tx.Delete(&models.BoardSetup{}, setupA.ID)
			rowsAffected = result.RowsAffected
			return result.Error
		})

		assert.Zero(t, rowsAffected)
	})

	// Cleanup
	DB.Unscoped().Delete(&models.BoardSetup{}, setupA.ID)
	DB.Unscoped().Delete(&models.User{}, userA.ID)
	DB.Unscoped().Delete(&models.User{}, userB.ID)
	DB.Unscoped().Exec("DELETE FROM user_stats WHERE user_id IN (?, ?)", userA.ID, userB.ID)
}
