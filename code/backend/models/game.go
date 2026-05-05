package models

import (
	"time"

	"gorm.io/gorm"
)

// Game represents a played or ongoing game session
type Game struct {
	ID            string         `json:"id" gorm:"primaryKey;type:varchar(100)"`
	Player1UserID *int           `json:"player1_user_id" gorm:"index"`
	Player2UserID *int           `json:"player2_user_id" gorm:"index"`
	WinnerID      *int           `json:"winner_id"`
	GameType      string         `json:"game_type" gorm:"not null;size:50"`
	InitialState  string         `json:"initial_state" gorm:"type:jsonb;not null"`
	CreatedAt     time.Time      `json:"created_at"`
	FinishedAt    *time.Time     `json:"finished_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
	Moves         []GameMove     `json:"moves,omitempty" gorm:"foreignKey:GameID;references:ID"`
}

// GameMove represents a single move in a game
type GameMove struct {
	ID           int            `json:"id" gorm:"primaryKey"`
	GameID       string         `json:"game_id" gorm:"index;type:varchar(100);not null"`
	MoveIndex    int            `json:"move_index" gorm:"not null"`
	PlayerID     int            `json:"player_id" gorm:"not null"`
	FromX        int            `json:"from_x" gorm:"not null"`
	FromY        int            `json:"from_y" gorm:"not null"`
	ToX          int            `json:"to_x" gorm:"not null"`
	ToY          int            `json:"to_y" gorm:"not null"`
	AttackerData string         `json:"attacker_data,omitempty" gorm:"type:jsonb"`
	DefenderData string         `json:"defender_data,omitempty" gorm:"type:jsonb"`
	Result       MoveResultType `json:"result" gorm:"not null;size:20"`
	CreatedAt    time.Time      `json:"created_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}
