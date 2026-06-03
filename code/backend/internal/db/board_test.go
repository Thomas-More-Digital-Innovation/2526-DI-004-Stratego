package db

import (
	"context"
	"digital-innovation/gostrategy/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestBoardSetupLogic(t *testing.T) {
	SetupDBTest(t)
	ctx := context.Background()
	user, _ := CreateUser(ctx, "setupuser", "Pass1234!", "")
	ctx = WithUserID(ctx, user.ID)

	t.Run("CRUD BoardSetup", func(t *testing.T) {
		setup, err := CreateBoardSetup(ctx, user.ID, "My Setup", "Desc", "DATA", true)
		assert.NoError(t, err)

		// Read
		retrieved, err := GetBoardSetup(ctx, setup.ID, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, "My Setup", retrieved.Name)

		// Update
		err = UpdateBoardSetup(ctx, setup.ID, user.ID, "Updated Name", "", "", false)
		assert.NoError(t, err)

		// Verify update
		DB.First(&setup, setup.ID)
		assert.Equal(t, "Updated Name", setup.Name)

		// Delete
		err = DeleteBoardSetup(ctx, setup.ID, user.ID)
		assert.NoError(t, err)

		// Verify soft delete
		err = DB.First(&models.BoardSetup{}, setup.ID).Error
		assert.Equal(t, gorm.ErrRecordNotFound, err)

		// Verify it still exists in DB (soft deleted)
		var deleted models.BoardSetup
		err = DB.Unscoped().First(&deleted, setup.ID).Error
		assert.NoError(t, err)
	})
}
