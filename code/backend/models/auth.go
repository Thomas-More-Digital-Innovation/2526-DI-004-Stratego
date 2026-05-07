// Package models defines the data structures used in the game
package models

import (
	"time"

	"gorm.io/gorm"
)

// RefreshToken represents a long-lived session token in the database
type RefreshToken struct {
	ID        int            `json:"id" gorm:"primaryKey"`
	UserID    int            `json:"user_id" gorm:"index;not null"`
	Token     string         `json:"token" gorm:"unique;not null;size:255"`
	ExpiresAt time.Time      `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
