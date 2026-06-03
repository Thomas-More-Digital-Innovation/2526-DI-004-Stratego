package auth_test

import (
	"digital-innovation/gostrategy/internal/auth"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateAndVerifyToken(t *testing.T) {
	userID := 123
	username := "testuser"

	// Generate token
	token, err := auth.GenerateToken(userID, username)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify token
	user, err := auth.VerifyToken(token)
	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, username, user.Username)
}

func TestVerifyInvalidToken(t *testing.T) {
	// Test with malformed token
	user, err := auth.VerifyToken("invalid.token.string")
	assert.Error(t, err)
	assert.Nil(t, user)

	// Test with expired token
	claims := auth.CustomClaims{
		Username: "expired",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "456",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(auth.GetJWTSecret())

	user, err = auth.VerifyToken(tokenString)
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestGenerateRefreshToken(t *testing.T) {
	token1, err := auth.GenerateRefreshToken()
	require.NoError(t, err)
	assert.NotEmpty(t, token1)

	token2, err := auth.GenerateRefreshToken()
	require.NoError(t, err)
	assert.NotEmpty(t, token2)
	assert.NotEqual(t, token1, token2)
}
