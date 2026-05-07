package db

import (
	"context"
	"digital-innovation/gostrategy/models"
	"testing"

	"gorm.io/gorm"
)

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
