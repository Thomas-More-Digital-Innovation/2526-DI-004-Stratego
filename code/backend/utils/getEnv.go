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

func GetEnvOrErrorInProduction(key string, devFallback string) (string, error) {
	if !IsProduction() {
		return GetEnv(key, devFallback), nil
	}
	return GetEnvOrError(key)
}
