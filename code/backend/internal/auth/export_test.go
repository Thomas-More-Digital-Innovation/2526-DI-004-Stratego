package auth

// Export internal secret for testing
func GetJWTSecret() []byte {
	return jwtSecret
}

func SetJWTSecret(secret []byte) {
	jwtSecret = secret
}
