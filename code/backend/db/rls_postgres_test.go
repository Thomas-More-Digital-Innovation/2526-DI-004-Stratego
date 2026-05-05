package db

import (
	"context"
	"digital-innovation/stratego/models"
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupPostgresDB(t *testing.T) *gorm.DB {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		t.Skip("Skipping Postgres RLS test: TEST_DB_DSN not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to postgres: %v", err)
	}

	// Register our RLS plugin for the test DB instance
	rlsPlugin(db)

	// Set global DB for RunMigrations to work
	oldDB := DB
	DB = db
	if err := RunMigrations(context.Background()); err != nil {
		DB = oldDB
		t.Fatalf("failed to run migrations: %v", err)
	}
	DB = oldDB

	return db
}

func TestRowLevelSecurity(t *testing.T) {
	testDB := setupPostgresDB(t)
	// Swap global DB
	oldDB := DB
	DB = testDB
	defer func() { DB = oldDB }()

	ctx := context.Background()

	// 1. Create two users
	userA, _ := CreateUser(ctx, "user_a", "Pass1234!", "")
	userB, _ := CreateUser(ctx, "user_b", "Pass1234!", "")

	// 2. User A creates a board setup
	setupA, err := CreateBoardSetup(WithUserID(ctx, userA.ID), userA.ID, "A's Setup", "", "DATA_A", false)
	if err != nil {
		t.Fatalf("failed to create setup for A: %v", err)
	}

	t.Run("User B cannot see User A's setup via RLS", func(t *testing.T) {
		// Attempt to fetch setupA using User B's context
		var retrieved models.BoardSetup
		err := DB.WithContext(WithUserID(ctx, userB.ID)).First(&retrieved, setupA.ID).Error
		
		if err == nil {
			t.Errorf("Security Breach: User B successfully retrieved User A's setup!")
		} else if err != gorm.ErrRecordNotFound {
			t.Errorf("Expected ErrRecordNotFound due to RLS, got: %v", err)
		}
	})

	t.Run("User A can see their own setup", func(t *testing.T) {
		var retrieved models.BoardSetup
		err := DB.WithContext(WithUserID(ctx, userA.ID)).First(&retrieved, setupA.ID).Error
		
		if err != nil {
			t.Errorf("User A failed to retrieve their own setup: %v", err)
		}
		if retrieved.ID != setupA.ID {
			t.Errorf("expected ID %d, got %d", setupA.ID, retrieved.ID)
		}
	})

	t.Run("User B cannot delete User A's setup", func(t *testing.T) {
		result := DB.WithContext(WithUserID(ctx, userB.ID)).Delete(&models.BoardSetup{}, setupA.ID)
		if result.Error != nil {
			// Some DBs might return error on RLS violation for Delete, 
			// but GORM usually reports 0 rows affected if RLS hides the row.
		}
		
		if result.RowsAffected > 0 {
			t.Errorf("Security Breach: User B successfully deleted User A's setup!")
		}
	})

	// Cleanup (unscoped to actually remove from DB if needed, or just let it be)
	DB.Unscoped().Delete(&models.BoardSetup{}, setupA.ID)
	DB.Unscoped().Delete(&models.User{}, userA.ID)
	DB.Unscoped().Delete(&models.User{}, userB.ID)
	DB.Unscoped().Exec("DELETE FROM user_stats WHERE user_id IN (?, ?)", userA.ID, userB.ID)
}
