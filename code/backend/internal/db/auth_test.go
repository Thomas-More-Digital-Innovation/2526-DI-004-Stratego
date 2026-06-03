package db

import (
	"context"
	"digital-innovation/gostrategy/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAuthLogic(t *testing.T) {
	SetupDBTest(t)
	ctx := context.Background()

	t.Run("AuthenticateUser", func(t *testing.T) {
		_, err := CreateUser(ctx, "authuser", "Pass1234!", "")
		assert.NoError(t, err)

		user, err := AuthenticateUser(ctx, "authuser", "Pass1234!")
		assert.NoError(t, err)
		assert.NotNil(t, user)

		_, err = AuthenticateUser(ctx, "authuser", "wrongpass")
		assert.Error(t, err)
	})

	t.Run("RefreshToken Lifecycle", func(t *testing.T) {
		user, _ := CreateUser(ctx, "rtuser", "Pass1234!", "")
		token := "some-random-token"
		expiry := time.Now().Add(time.Hour)

		ctx = WithUserID(ctx, user.ID)
		err := SaveRefreshToken(ctx, user.ID, token, expiry)
		assert.NoError(t, err)

		id, err := GetUserIDByRefreshToken(ctx, token)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, id)

		err = DeleteRefreshToken(ctx, token)
		assert.NoError(t, err)

		_, err = GetUserIDByRefreshToken(ctx, token)
		assert.Error(t, err)
	})

	t.Run("Token is Hashed in DB", func(t *testing.T) {
		user, _ := CreateUser(ctx, "hashuser", "Pass1234!", "")
		token := "secret-token-123"
		ctx = WithUserID(ctx, user.ID)
		_ = SaveRefreshToken(ctx, user.ID, token, time.Now().Add(time.Hour))

		// G117: try to find the raw token using direct GORM (bypass our helper)
		var rt models.RefreshToken
		err := DB.Where("token = ?", token).First(&rt).Error
		assert.Error(t, err)

		err = DB.Where("token = ?", hashToken(token)).First(&rt).Error
		assert.NoError(t, err)
	})
}
