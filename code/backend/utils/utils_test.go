package utils

import (
	"digital-innovation/gostrategy/models"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testUserStr = "testuser"

func TestGetIntSafe(t *testing.T) {
	val := 42
	tests := []struct {
		name string
		ptr  *int
		want int
	}{
		{"nil pointer", nil, 0},
		{"valid pointer", &val, 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetIntSafe(tt.ptr); got != tt.want {
				t.Errorf("GetIntSafe() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsProduction(t *testing.T) {
	// Clear any existing APP_ENV to ensure clean state
	originalEnv := os.Getenv("APP_ENV")
	defer func() { _ = os.Setenv("APP_ENV", originalEnv) }()

	_ = os.Setenv("APP_ENV", "development")
	assert.False(t, IsProduction())

	_ = os.Setenv("APP_ENV", "production")
	assert.True(t, IsProduction())
}

func TestTryGetUser(t *testing.T) {
	user := &models.User{ID: 1, Username: testUserStr}
	tests := []struct {
		name     string
		user     *models.User
		wantUser string
		wantID   int
	}{
		{"nil user", nil, "Unknown", 0},
		{"valid user", user, testUserStr, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUser, gotID := TryGetUser(tt.user)
			if gotUser != tt.wantUser || gotID != tt.wantID {
				t.Errorf("TryGetUser() = (%v, %v), want (%v, %v)", gotUser, gotID, tt.wantUser, tt.wantID)
			}
		})
	}
}

func TestTryGetUserOrError(t *testing.T) {
	user := &models.User{ID: 1, Username: "testuser"}
	tests := []struct {
		name     string
		user     *models.User
		wantUser string
		wantID   int
		wantErr  bool
	}{
		{"nil user", nil, "", 0, true},
		{"valid user", user, "testuser", 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUser, gotID, err := TryGetUserOrError(tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("TryGetUserOrError() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUser != tt.wantUser || gotID != tt.wantID {
				t.Errorf("TryGetUserOrError() = (%v, %v), want (%v, %v)", gotUser, gotID, tt.wantUser, tt.wantID)
			}
		})
	}
}
