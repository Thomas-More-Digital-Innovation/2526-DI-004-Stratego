package utils

import "os"

func GetEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

var isProd bool

func init() {
	isProd = os.Getenv("APP_ENV") == "production"
}
func IsProduction() bool {
	return isProd
}
