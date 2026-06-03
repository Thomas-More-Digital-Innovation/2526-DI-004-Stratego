package db

import (
	"context"
	"digital-innovation/gostrategy/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserLogic(t *testing.T) {
	SetupDBTest(t)
	ctx := context.Background()

	t.Run("CreateUser and Stats", func(t *testing.T) {
		user, err := CreateUser(ctx, "testuser", "StrongPass1!", "pic.png")
		assert.NoError(t, err)
		assert.Equal(t, "testuser", user.Username)

		// Verify stats were created
		var stats models.UserStats
		err = DB.Where("user_id = ?", user.ID).First(&stats).Error
		assert.NoError(t, err)
	})

	t.Run("UpdateUserPassword", func(t *testing.T) {
		user, _ := CreateUser(ctx, "passuser", "OldPass1!", "")
		err := UpdateUserPassword(ctx, user.ID, "NewPass1!")
		assert.NoError(t, err)

		// Verify authentication works with new password
		_, err = AuthenticateUser(ctx, "passuser", "NewPass1!")
		assert.NoError(t, err)
	})
}

func TestUserSoftDelete(t *testing.T) {
	SetupDBTest(t)
	ctx := context.Background()
	user, _ := CreateUser(ctx, "delete-me", "Pass1234!", "")

	t.Run("Soft Deleted User cannot authenticate", func(t *testing.T) {
		err := DB.Delete(user).Error
		assert.NoError(t, err)

		_, err = AuthenticateUser(ctx, "delete-me", "Pass1234!")
		assert.Error(t, err)
	})
}

func TestStatsConcurrency(t *testing.T) {
	SetupDBTest(t)
	ctx := context.Background()
	user, _ := CreateUser(ctx, "concuruser", "Pass1234!", "")

	t.Run("Concurrent stats updates are consistent", func(t *testing.T) {
		const goroutines = 10
		const updatesPerRoutine = 5
		done := make(chan bool)

		for i := 0; i < goroutines; i++ {
			go func() {
				for j := 0; j < updatesPerRoutine; j++ {
					_ = UpdateUserStats(ctx, user.ID, true, 10, 100.0)
				}
				done <- true
			}()
		}

		for i := 0; i < goroutines; i++ {
			<-done
		}

		// Verify results
		var stats models.UserStats
		DB.Where("user_id = ?", user.ID).First(&stats)

		expectedGames := goroutines * updatesPerRoutine
		assert.Equal(t, expectedGames, stats.TotalGames)
		assert.Equal(t, expectedGames, stats.Wins)
	})
}
