package db

import (
	"context"
	"digital-innovation/stratego/models"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupSQLiteDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	// AutoMigrate all models
	err = db.AutoMigrate(
		&models.User{},
		&models.UserStats{},
		&models.BoardSetup{},
		&models.RefreshToken{},
		&models.Game{},
		&models.GameMove{},
	)
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	return db
}

func TestUserLogic(t *testing.T) {
	testDB := setupSQLiteDB(t)
	// Swap global DB for testing
	oldDB := DB
	DB = testDB
	defer func() { DB = oldDB }()

	ctx := context.Background()

	t.Run("CreateUser and Stats", func(t *testing.T) {
		user, err := CreateUser(ctx, "testuser", "StrongPass1!", "pic.png")
		if err != nil {
			t.Fatalf("CreateUser failed: %v", err)
		}

		if user.Username != "testuser" {
			t.Errorf("expected username testuser, got %s", user.Username)
		}

		// Verify stats were created
		var stats models.UserStats
		err = DB.Where("user_id = ?", user.ID).First(&stats).Error
		if err != nil {
			t.Errorf("stats not created for user: %v", err)
		}
	})

	t.Run("AuthenticateUser", func(t *testing.T) {
		_, err := CreateUser(ctx, "authuser", "Pass1234!", "")
		if err != nil {
			t.Fatalf("failed to create user: %v", err)
		}

		user, err := AuthenticateUser(ctx, "authuser", "Pass1234!")
		if err != nil {
			t.Errorf("AuthenticateUser failed: %v", err)
		}
		if user == nil {
			t.Fatal("expected user, got nil")
		}

		_, err = AuthenticateUser(ctx, "authuser", "wrongpass")
		if err == nil {
			t.Error("expected error for wrong password, got nil")
		}
	})
}

func TestBoardSetupLogic(t *testing.T) {
	testDB := setupSQLiteDB(t)
	oldDB := DB
	DB = testDB
	defer func() { DB = oldDB }()

	ctx := context.Background()
	user, _ := CreateUser(ctx, "setupuser", "Pass1234!", "")

	t.Run("CRUD BoardSetup", func(t *testing.T) {
		setup, err := CreateBoardSetup(ctx, user.ID, "My Setup", "Desc", "DATA", true)
		if err != nil {
			t.Fatalf("CreateBoardSetup failed: %v", err)
		}

		// Read
		retrieved, err := GetBoardSetup(ctx, setup.ID, user.ID)
		if err != nil {
			t.Fatalf("GetBoardSetup failed: %v", err)
		}
		if retrieved.Name != "My Setup" {
			t.Errorf("expected name 'My Setup', got %s", retrieved.Name)
		}

		// Update
		err = UpdateBoardSetup(ctx, setup.ID, user.ID, "Updated Name", "", "", false)
		if err != nil {
			t.Fatalf("UpdateBoardSetup failed: %v", err)
		}

		// Verify update
		DB.First(&setup, setup.ID)
		if setup.Name != "Updated Name" {
			t.Errorf("expected updated name, got %s", setup.Name)
		}

		// Delete
		err = DeleteBoardSetup(ctx, setup.ID, user.ID)
		if err != nil {
			t.Fatalf("DeleteBoardSetup failed: %v", err)
		}

		// Verify soft delete
		err = DB.First(&models.BoardSetup{}, setup.ID).Error
		if err != gorm.ErrRecordNotFound {
			t.Error("expected record to be soft-deleted (not found by First)")
		}

		// Verify it still exists in DB (soft deleted)
		var deleted models.BoardSetup
		err = DB.Unscoped().First(&deleted, setup.ID).Error
		if err != nil {
			t.Errorf("expected to find soft-deleted record with Unscoped, got error: %v", err)
		}
	})
}
