package auth

import (
	"crypto/rand"
	"digital-innovation/gostrategy/models"
	"digital-innovation/gostrategy/utils"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

func init() {
	secret, err := utils.GetEnvOrErrorInProduction("JWT_SECRET", "dev_fallback_secret")
	if err != nil {
		panic(fmt.Errorf("critical: JWT_SECRET is not set in production: %w", err))
	}
	jwtSecret = []byte(secret)
}

// CustomClaims defines the structure of our JWT payload
type CustomClaims struct {
	Username string `json:"name"`
	jwt.RegisteredClaims
}

// GenerateToken creates a new JWT for a user
func GenerateToken(userID int, username string) (string, error) {
	claims := CustomClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", userID),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(MaxCookieAge) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// VerifyToken validates the JWT and returns the user information
func VerifyToken(tokenString string) (*models.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		var userID int
		_, err = fmt.Sscanf(claims.Subject, "%d", &userID)
		if err != nil {
			return nil, fmt.Errorf("invalid token subject: %w", err)
		}

		return &models.User{
			ID:       userID,
			Username: claims.Username,
		}, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// GenerateRefreshToken creates a random string to be used as a refresh token
func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
