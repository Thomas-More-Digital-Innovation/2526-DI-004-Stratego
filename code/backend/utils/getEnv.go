package utils

import (
	"errors"
	"os"
)

func GetEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func GetEnvOrError(key string) (string, error) {
	if value := os.Getenv(key); value == "" {
		return "", errors.New("environment variable " + key + " not set")
	}
	return os.Getenv(key), nil
}

var isProd bool

func init() {
	isProd = os.Getenv("APP_ENV") == "production"
}
func IsProduction() bool {
	return isProd
}
