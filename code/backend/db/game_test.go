package db

import (
	"context"
	"digital-innovation/stratego/models"
	"testing"
)

func TestGameLogic(t *testing.T) {
	testDB := setupSQLiteDB(t)
	oldDB := DB
	DB = testDB
	defer func() { DB = oldDB }()

	ctx := context.Background()
	user1, _ := CreateUser(ctx, "p1", "Pass1234!", "")
	user2, _ := CreateUser(ctx, "p2", "Pass1234!", "")

	t.Run("Save Game and Moves", func(t *testing.T) {
		gameID := "game-123"
		initialState := map[string]interface{}{"board": "initial"}

		err := SaveGame(ctx, gameID, &user1.ID, &user2.ID, "ranked", initialState, nil)
		if err != nil {
			t.Fatalf("SaveGame failed: %v", err)
		}

		moves := []models.HistoricalMove{
			{MoveIndex: 1, PlayerID: 1, FromX: 0, FromY: 0, ToX: 0, ToY: 1, Result: "moved"},
			{MoveIndex: 2, PlayerID: 2, FromX: 9, FromY: 9, ToX: 9, ToY: 8, Result: "moved"},
		}

		err = SaveMoves(ctx, gameID, moves)
		if err != nil {
			t.Fatalf("SaveMoves failed: %v", err)
		}

		// Verify history
		history, err := GetGameHistory(ctx, gameID)
		if err != nil {
			t.Fatalf("GetGameHistory failed: %v", err)
		}

		if len(history.Moves) != 2 {
			t.Errorf("expected 2 moves, got %d", len(history.Moves))
		}
		if history.Moves[0].MoveIndex != 1 {
			t.Errorf("expected first move index 1, got %d", history.Moves[0].MoveIndex)
		}
	})
}
