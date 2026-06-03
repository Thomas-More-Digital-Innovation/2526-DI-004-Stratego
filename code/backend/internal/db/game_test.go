package db

import (
	"context"
	"digital-innovation/gostrategy/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGameLogic(t *testing.T) {
	SetupDBTest(t)
	ctx := context.Background()
	user1, _ := CreateUser(ctx, "p1", "Pass1234!", "")
	user2, _ := CreateUser(ctx, "p2", "Pass1234!", "")
	ctx = WithUserID(ctx, user1.ID)

	t.Run("Save Game and Moves", func(t *testing.T) {
		gameID := "game-123"
		initialState := map[string]any{"board": "initial"}

		err := SaveGame(ctx, gameID, &user1.ID, &user2.ID, "ranked", initialState, nil, time.Now(), time.Now())
		assert.NoError(t, err)

		moves := []models.HistoricalMove{
			{MoveIndex: 1, PlayerID: 1, FromX: 0, FromY: 0, ToX: 0, ToY: 1, Result: "moved"},
			{MoveIndex: 2, PlayerID: 2, FromX: 9, FromY: 9, ToX: 9, ToY: 8, Result: "moved"},
		}

		err = SaveMoves(ctx, gameID, moves)
		assert.NoError(t, err)

		// Verify history
		history, err := GetGameHistory(ctx, gameID)
		assert.NoError(t, err)
		assert.Len(t, history.Moves, 2)
		assert.Equal(t, 1, history.Moves[0].MoveIndex)
	})
}
