package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID             int            `json:"id" gorm:"primaryKey"`
	Username       string         `json:"username" gorm:"unique;not null;size:50"`
	PasswordHash   string         `json:"-" gorm:"not null;size:255"` // never send to client
	ProfilePicture string         `json:"profile_picture,omitempty" gorm:"size:255"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

// UserStats represents game statistics for a user
type UserStats struct {
	ID                  int            `json:"id" gorm:"primaryKey"`
	UserID              int            `json:"user_id" gorm:"unique;not null"`
	TotalGames          int            `json:"total_games" gorm:"default:0"`
	Wins                int            `json:"wins" gorm:"default:0"`
	Losses              int            `json:"losses" gorm:"default:0"`
	Draws               int            `json:"draws" gorm:"default:0"`
	TotalMoves          int            `json:"total_moves" gorm:"default:0"`
	AvgGameDurationSecs float64        `json:"avg_game_duration_seconds" gorm:"default:0"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `json:"-" gorm:"index"`
}

// BoardSetup represents a saved board configuration
type BoardSetup struct {
	ID          int            `json:"id" gorm:"primaryKey"`
	UserID      int            `json:"user_id" gorm:"index;not null"`
	Name        string         `json:"name" gorm:"not null;size:100"`
	Description string         `json:"description,omitempty" gorm:"type:text"`
	SetupData   string         `json:"setup_data" gorm:"not null;size:40"` // JSON string of piece positions
	IsDefault   bool           `json:"is_default" gorm:"default:false"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// CreateUserRequest for user registration
type CreateUserRequest struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	ProfilePicture string `json:"profile_picture,omitempty"`
}

// LoginRequest for user login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// CreateBoardSetupRequest for creating a board setup
type CreateBoardSetupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	SetupData   string `json:"setup_data"`
	IsDefault   bool   `json:"is_default"`
}

// UpdateBoardSetupRequest for updating a board setup
type UpdateBoardSetupRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	SetupData   string `json:"setup_data,omitempty"`
	IsDefault   bool   `json:"is_default"`
}

// ChangePasswordRequest for updating user password
type ChangePasswordRequest struct {
	OldPassword     string `json:"old_password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}
