package db

import (
	"context"
	"digital-innovation/gostrategy/models"
	"fmt"

	"gorm.io/gorm"
)

// CreateBoardSetup saves a new board setup
func CreateBoardSetup(ctx context.Context, userID int, name, description, setupData string, isDefault bool) (*models.BoardSetup, error) {
	setup := models.BoardSetup{
		UserID:      userID,
		Name:        name,
		Description: description,
		SetupData:   setupData,
		IsDefault:   isDefault,
	}
	err := WithRLS(ctx, func(tx *gorm.DB) error {
		return tx.Create(&setup).Error
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create board setup: %w", err)
	}
	return &setup, nil
}

// GetBoardSetup retrieves a board setup by ID and verifying ownership
func GetBoardSetup(ctx context.Context, setupID, userID int) (*models.BoardSetup, error) {
	var setup models.BoardSetup
	err := WithRLS(ctx, func(tx *gorm.DB) error {
		return tx.Where("id = ? AND user_id = ?", setupID, userID).First(&setup).Error
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get board setup: %w", err)
	}
	return &setup, nil
}

// GetUserBoardSetups retrieves all board setups for a user
func GetUserBoardSetups(ctx context.Context, userID int) ([]models.BoardSetup, error) {
	var setups []models.BoardSetup
	err := WithRLS(ctx, func(tx *gorm.DB) error {
		return tx.Where("user_id = ?", userID).Order("is_default desc, created_at desc").Find(&setups).Error
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query board setups: %w", err)
	}
	return setups, nil
}

// UpdateBoardSetup updates an existing board setup and verifying ownership
func UpdateBoardSetup(ctx context.Context, setupID, userID int, name, description, setupData string, isDefault bool) error {
	updates := map[string]any{
		"is_default": isDefault,
	}
	if name != "" {
		updates["name"] = name
	}
	if description != "" {
		updates["description"] = description
	}
	if setupData != "" {
		updates["setup_data"] = setupData
	}

	return WithRLS(ctx, func(tx *gorm.DB) error {
		result := tx.Model(&models.BoardSetup{}).Where("id = ? AND user_id = ?", setupID, userID).Updates(updates)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("board setup not found or not owned by user")
		}
		return nil
	})
}

// DeleteBoardSetup deletes a board setup
func DeleteBoardSetup(ctx context.Context, setupID, userID int) error {
	return WithRLS(ctx, func(tx *gorm.DB) error {
		result := tx.Where("id = ? AND user_id = ?", setupID, userID).Delete(&models.BoardSetup{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("board setup not found or not owned by user")
		}
		return nil
	})
}
