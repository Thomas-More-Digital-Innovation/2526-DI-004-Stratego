package db

import (
	"context"
	"digital-innovation/gostrategy/internal/models"
	"testing"
	"time"
)

func TestAuthLogic(t *testing.T) {
	testDB := SetupTestDB(t)
	oldDB := DB
	DB = testDB
	defer func() { DB = oldDB }()

	ctx := context.Background()

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

	t.Run("RefreshToken Lifecycle", func(t *testing.T) {
		user, _ := CreateUser(ctx, "rtuser", "Pass1234!", "")
		token := "some-random-token"
		expiry := time.Now().Add(time.Hour)

		ctx = WithUserID(ctx, user.ID)
		err := SaveRefreshToken(ctx, user.ID, token, expiry)
		if err != nil {
			t.Fatalf("SaveRefreshToken failed: %v", err)
		}

		// Validate
		id, err := GetUserIDByRefreshToken(ctx, token)
		if err != nil {
			t.Fatalf("GetUserIDByRefreshToken failed: %v", err)
		}
		if id != user.ID {
			t.Errorf("expected user ID %d, got %d", user.ID, id)
		}

		// Delete
		err = DeleteRefreshToken(ctx, token)
		if err != nil {
			t.Fatalf("DeleteRefreshToken failed: %v", err)
		}

		// Verify deleted
		_, err = GetUserIDByRefreshToken(ctx, token)
		if err == nil {
			t.Error("expected error for deleted token, got nil")
		}
	})

	t.Run("Token is Hashed in DB", func(t *testing.T) {
		user, _ := CreateUser(ctx, "hashuser", "Pass1234!", "")
		token := "secret-token-123"
		ctx = WithUserID(ctx, user.ID)
		_ = SaveRefreshToken(ctx, user.ID, token, time.Now().Add(time.Hour))

		// Try to find the raw token using direct GORM (bypass our helper)
		var rt models.RefreshToken
		err := DB.Where("token = ?", token).First(&rt).Error
		if err == nil {
			t.Error("Security Vulnerability: found raw token in database!")
		}

		// Try to find the hashed token
		err = DB.Where("token = ?", hashToken(token)).First(&rt).Error
		if err != nil {
			t.Errorf("expected to find hashed token in DB, got error: %v", err)
		}
	})
}
