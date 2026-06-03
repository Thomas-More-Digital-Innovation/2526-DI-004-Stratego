package models

import "time"

// AIParameter represents the stored GORM model for AI aggression and weights.
type AIParameter struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	AIType     string    `json:"ai_type" gorm:"uniqueIndex:idx_ai_type_name;size:50;not null"`
	Name       string    `json:"name" gorm:"uniqueIndex:idx_ai_type_name;size:50;not null"`
	Aggression float64   `json:"aggression" gorm:"not null;default:0.5"`
	Weights    string    `json:"weights" gorm:"type:text"`
	Config     string    `json:"config" gorm:"type:text"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
