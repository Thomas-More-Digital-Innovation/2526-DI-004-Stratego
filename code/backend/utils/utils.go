package utils

import (
	"digital-innovation/stratego/models"
	"errors"
)

// GetIntSafe returns the value of a pointer to an int, or 0 if the pointer is nil
func GetIntSafe(ptr *int) int {
	if ptr == nil {
		return 0
	}
	return *ptr
}

var isProd bool

func init() {
	// Default to development if not set to avoid panics during tests
	appEnv := GetEnv("APP_ENV", "development")
	isProd = appEnv == "production"
}

// IsProduction checks if the application is running in production mode
func IsProduction() bool {
	return isProd
}

func TryGetUser(user *models.User) (string, int) {
	if user == nil {
		return "Unknown", 0
	}
	return user.Username, user.ID
}

func TryGetUserOrError(user *models.User) (string, int, error) {
	if user == nil {
		return "", 0, errors.New("Couldn't find user.")
	}
	return user.Username, user.ID, nil
}
