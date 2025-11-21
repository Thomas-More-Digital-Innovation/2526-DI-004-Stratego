package models

import "time"

// User represents a user in the system
type User struct {
	ID             int       `json:"id"`
	Username       string    `json:"username"`
	PasswordHash   string    `json:"-"` // never send to client
	ProfilePicture string    `json:"profile_picture,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// UserStats represents game statistics for a user
type UserStats struct {
	ID                  int       `json:"id"`
	UserID              int       `json:"user_id"`
	TotalGames          int       `json:"total_games"`
	Wins                int       `json:"wins"`
	Losses              int       `json:"losses"`
	Draws               int       `json:"draws"`
	TotalMoves          int       `json:"total_moves"`
	AvgGameDurationSecs float64   `json:"avg_game_duration_seconds"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// BoardSetup represents a saved board configuration
type BoardSetup struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	SetupData   string    `json:"setup_data"` // JSON string of piece positions
	IsDefault   bool      `json:"is_default"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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
