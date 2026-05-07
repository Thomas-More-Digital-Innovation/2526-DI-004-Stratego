package db

import (
	"context"
	"digital-innovation/gostrategy/models"
	"testing"
)

func TestUserLogic(t *testing.T) {
	testDB := setupSQLiteDB(t)
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

	t.Run("UpdateUserPassword", func(t *testing.T) {
		user, _ := CreateUser(ctx, "passuser", "OldPass1!", "")
		err := UpdateUserPassword(ctx, user.ID, "NewPass1!")
		if err != nil {
			t.Fatalf("UpdateUserPassword failed: %v", err)
		}

		// Verify authentication works with new password
		_, err = AuthenticateUser(ctx, "passuser", "NewPass1!")
		if err != nil {
			t.Errorf("failed to authenticate with new password: %v", err)
		}
	})
}

func TestUserSoftDelete(t *testing.T) {
	testDB := setupSQLiteDB(t)
	oldDB := DB
	DB = testDB
	defer func() { DB = oldDB }()

	ctx := context.Background()
	user, _ := CreateUser(ctx, "delete-me", "Pass1234!", "")

	t.Run("Soft Deleted User cannot authenticate", func(t *testing.T) {
		// Delete user
		err := DB.Delete(user).Error
		if err != nil {
			t.Fatalf("failed to delete user: %v", err)
		}

		// Try to authenticate
		_, err = AuthenticateUser(ctx, "delete-me", "Pass1234!")
		if err == nil {
			t.Error("expected authentication to fail for soft-deleted user, but it succeeded")
		}
	})
}

func TestStatsConcurrency(t *testing.T) {
	testDB := setupSQLiteDB(t)
	oldDB := DB
	DB = testDB
	defer func() { DB = oldDB }()

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
		if stats.TotalGames != expectedGames {
			t.Errorf("expected %d games, got %d", expectedGames, stats.TotalGames)
		}
		if stats.Wins != expectedGames {
			t.Errorf("expected %d wins, got %d", expectedGames, stats.Wins)
		}
	})
}
