// Package db provides database access and persistence for game data
package db

import (
	"context"
	"digital-innovation/stratego/models"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

// SaveGame persists the game metadata and initial state
func SaveGame(ctx context.Context, gameID string, p1ID, p2ID *int, gameType string, initialState interface{}, winnerID *int) error {
	stateJSON, err := json.Marshal(initialState)
	if err != nil {
		return fmt.Errorf("failed to marshal initial state: %w", err)
	}

	game := models.Game{
		ID:            gameID,
		Player1UserID: p1ID,
		Player2UserID: p2ID,
		WinnerID:      winnerID,
		GameType:      gameType,
		InitialState:  string(stateJSON),
	}

	err = DB.WithContext(ctx).Create(&game).Error
	if err != nil {
		return fmt.Errorf("failed to save game: %w", err)
	}
	return nil
}

// SaveMove persists a single move in a game's history
func SaveMove(ctx context.Context, gameID string, move models.HistoricalMove) error {
	var attackerData, defenderData string

	if move.Attacker != nil {
		b, _ := json.Marshal(move.Attacker)
		attackerData = string(b)
	}
	if move.Defender != nil {
		b, _ := json.Marshal(move.Defender)
		defenderData = string(b)
	}

	gameMove := models.GameMove{
		GameID:       gameID,
		MoveIndex:    move.MoveIndex,
		PlayerID:     move.PlayerID,
		FromX:        move.FromX,
		FromY:        move.FromY,
		ToX:          move.ToX,
		ToY:          move.ToY,
		AttackerData: attackerData,
		DefenderData: defenderData,
		Result:       move.Result,
	}

	err := DB.WithContext(ctx).Create(&gameMove).Error
	if err != nil {
		return fmt.Errorf("failed to save move: %w", err)
	}
	return nil
}

// SaveMoves persists multiple moves in a single transaction
func SaveMoves(ctx context.Context, gameID string, moves []models.HistoricalMove) error {
	if len(moves) == 0 {
		return nil
	}

	var gameMoves []models.GameMove
	for _, move := range moves {
		var attackerData, defenderData string

		if move.Attacker != nil {
			b, _ := json.Marshal(move.Attacker)
			attackerData = string(b)
		}
		if move.Defender != nil {
			b, _ := json.Marshal(move.Defender)
			defenderData = string(b)
		}

		gameMoves = append(gameMoves, models.GameMove{
			GameID:       gameID,
			MoveIndex:    move.MoveIndex,
			PlayerID:     move.PlayerID,
			FromX:        move.FromX,
			FromY:        move.FromY,
			ToX:          move.ToX,
			ToY:          move.ToY,
			AttackerData: attackerData,
			DefenderData: defenderData,
			Result:       move.Result,
		})
	}

	err := DB.WithContext(ctx).Create(&gameMoves).Error
	if err != nil {
		return fmt.Errorf("failed to save moves: %w", err)
	}
	return nil
}

// GetGameHistory retrieves the full move history of a completed game
func GetGameHistory(ctx context.Context, gameID string) (*models.GameHistory, error) {
	var game models.Game
	err := DB.WithContext(ctx).Preload("Moves", func(db *gorm.DB) *gorm.DB {
		return db.Order("move_index ASC")
	}).Where("id = ?", gameID).First(&game).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get game history: %w", err)
	}

	var history models.GameHistory
	history.GameID = game.ID
	history.WinnerID = game.WinnerID

	if err := json.Unmarshal([]byte(game.InitialState), &history.InitialState); err != nil {
		return nil, fmt.Errorf("failed to unmarshal initial state: %w", err)
	}

	for _, gm := range game.Moves {
		m := models.HistoricalMove{
			MoveIndex: gm.MoveIndex,
			PlayerID:  gm.PlayerID,
			FromX:     gm.FromX,
			FromY:     gm.FromY,
			ToX:       gm.ToX,
			ToY:       gm.ToY,
			Result:    gm.Result,
		}

		if gm.AttackerData != "" {
			var attacker models.PieceData
			if err := json.Unmarshal([]byte(gm.AttackerData), &attacker); err == nil {
				m.Attacker = &attacker
			}
		}
		if gm.DefenderData != "" {
			var defender models.PieceData
			if err := json.Unmarshal([]byte(gm.DefenderData), &defender); err == nil {
				m.Defender = &defender
			}
		}

		history.Moves = append(history.Moves, m)
	}

	return &history, nil
}
