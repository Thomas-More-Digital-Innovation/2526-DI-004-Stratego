package db

import (
	"context"
	"digital-innovation/stratego/models"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthenticateUser checks username and password, returns user if valid
func AuthenticateUser(ctx context.Context, username, password string) (*models.User, error) {
	var user models.User

	err := DB.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invalid username or password")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	return &user, nil
}

// SaveRefreshToken stores a new refresh token for a user
func SaveRefreshToken(ctx context.Context, userID int, token string, expiresAt time.Time) error {
	rt := models.RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
	}
	err := DB.WithContext(ctx).Create(&rt).Error
	if err != nil {
		return fmt.Errorf("failed to save refresh token: %w", err)
	}
	return nil
}

// GetUserIDByRefreshToken validates a refresh token and returns the owner's ID
func GetUserIDByRefreshToken(ctx context.Context, token string) (int, error) {
	var rt models.RefreshToken
	err := DB.WithContext(ctx).Where("token = ? AND expires_at > ?", token, time.Now()).First(&rt).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("invalid or expired refresh token")
		}
		return 0, fmt.Errorf("database error: %w", err)
	}
	return rt.UserID, nil
}

// DeleteRefreshToken removes a refresh token (e.g., on logout)
func DeleteRefreshToken(ctx context.Context, token string) error {
	return DB.WithContext(ctx).Where("token = ?", token).Delete(&models.RefreshToken{}).Error
}

// DeleteAllUserRefreshTokens revokes all sessions for a user
func DeleteAllUserRefreshTokens(ctx context.Context, userID int) error {
	return DB.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.RefreshToken{}).Error
}
