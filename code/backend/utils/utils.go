package utils

import (
	"digital-innovation/gostrategy/models"
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

// TryGetUser returns username and ID of the user if not nil, otherwise "Unknown" and 0
func TryGetUser(user *models.User) (string, int) {
	if user == nil {
		return "Unknown", 0
	}
	return user.Username, user.ID
}

// TryGetUserOrError returns username and ID of the user if not nil, otherwise an error
func TryGetUserOrError(user *models.User) (string, int, error) {
	if user == nil {
		return "", 0, errors.New("couldn't find user")
	}
	return user.Username, user.ID, nil
}
