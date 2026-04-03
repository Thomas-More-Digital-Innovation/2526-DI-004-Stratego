package utils

import "os"

func GetEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func IsProduction() bool {
	return GetEnv("APP_ENV", "development") == "production"
}
