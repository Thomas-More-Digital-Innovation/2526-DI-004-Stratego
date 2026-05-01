package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateAndVerifyToken(t *testing.T) {
	userID := 123
	username := "testuser"

	// Generate token
	token, err := GenerateToken(userID, username)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	if token == "" {
		t.Fatal("Generated token is empty")
	}

	// Verify token
	user, err := VerifyToken(token)
	if err != nil {
		t.Fatalf("Failed to verify token: %v", err)
	}
	if user == nil {
		t.Fatal("Verified user is nil")
	}
	if user.ID != userID {
		t.Errorf("Expected userID %d, got %d", userID, user.ID)
	}
	if user.Username != username {
		t.Errorf("Expected username %s, got %s", username, user.Username)
	}
}

func TestVerifyInvalidToken(t *testing.T) {
	// Test with malformed token
	user, err := VerifyToken("invalid.token.string")
	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	}
	if user != nil {
		t.Error("Expected nil user for invalid token")
	}

	// Test with expired token
	claims := CustomClaims{
		Username: "expired",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "456",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtSecret)

	user, err = VerifyToken(tokenString)
	if err == nil {
		t.Error("Expected error for expired token, got nil")
	}
	if user != nil {
		t.Error("Expected nil user for expired token")
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	token1, err := GenerateRefreshToken()
	if err != nil {
		t.Fatalf("Failed to generate refresh token 1: %v", err)
	}
	if token1 == "" {
		t.Fatal("Refresh token 1 is empty")
	}

	token2, err := GenerateRefreshToken()
	if err != nil {
		t.Fatalf("Failed to generate refresh token 2: %v", err)
	}
	if token2 == "" {
		t.Fatal("Refresh token 2 is empty")
	}
	if token1 == token2 {
		t.Fatal("Refresh tokens should be unique")
	}
}
