package db

import (
	"context"
	"digital-innovation/gostrategy/internal/models"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// CreateUser creates a new user with hashed password
func CreateUser(ctx context.Context, username, password, profilePicture string) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := models.User{
		Username:       username,
		PasswordHash:   string(hashedPassword),
		ProfilePicture: profilePicture,
	}

	err = DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		if isPostgresDialect() {
			if err := tx.Exec(fmt.Sprintf("SET LOCAL app.current_user_id = '%d'", user.ID)).Error; err != nil {
				return err
			}
		}
		stats := models.UserStats{UserID: user.ID}
		return tx.Create(&stats).Error
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

// UpdateUserPassword updates a user's password with a new hash
func UpdateUserPassword(ctx context.Context, userID int, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	err = DB.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("password_hash", string(hashedPassword)).Error
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// GetUserByID retrieves a user by ID
func GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	var user models.User
	err := DB.WithContext(ctx).First(&user, userID).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// GetUserStats retrieves stats for a user
func GetUserStats(ctx context.Context, userID int) (*models.UserStats, error) {
	var stats models.UserStats
	err := DB.WithContext(ctx).Where("user_id = ?", userID).First(&stats).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}
	return &stats, nil
}

// UpdateUserStats updates game statistics for a user
func UpdateUserStats(ctx context.Context, userID int, won bool, moveCount int, durationSecs float64) error {
	winsIncrement, lossesIncrement := 0, 0
	if won {
		winsIncrement = 1
	} else {
		lossesIncrement = 1
	}

	err := DB.WithContext(ctx).Model(&models.UserStats{}).Where("user_id = ?", userID).
		Updates(map[string]any{
			"total_games":               gorm.Expr("total_games + 1"),
			"wins":                      gorm.Expr("wins + ?", winsIncrement),
			"losses":                    gorm.Expr("losses + ?", lossesIncrement),
			"total_moves":               gorm.Expr("total_moves + ?", moveCount),
			"avg_game_duration_seconds": gorm.Expr("((avg_game_duration_seconds * total_games) + ?) / (total_games + 1)", durationSecs),
		}).Error

	if err != nil {
		return fmt.Errorf("failed to update user stats: %w", err)
	}
	return nil
}
