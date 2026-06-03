package db

import (
	"context"
	"digital-innovation/gostrategy/internal/models"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGamePagination(t *testing.T) {
	SetupDBTest(t)
	ctx := context.Background()
	user, _ := CreateUser(ctx, "paginator", "Pass1234!", "")
	ctx = WithUserID(ctx, user.ID)

	// Create 15 games
	for i := range 15 {
		gameID := fmt.Sprintf("paged-game-%d", i)
		err := SaveGame(ctx, gameID, &user.ID, nil, "ranked", map[string]string{"test": "data"}, nil, time.Now(), time.Now())
		assert.NoError(t, err)
	}

	t.Run("Count Games", func(t *testing.T) {
		count, err := GetGamesCountForUser(ctx, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, int64(15), count)
	})

	t.Run("Paged Retrieval Page 1", func(t *testing.T) {
		games, err := GetGamesForUserPaged(ctx, user.ID, 10, 0)
		assert.NoError(t, err)
		assert.Len(t, games, 10)
	})

	t.Run("Paged Retrieval Page 2", func(t *testing.T) {
		games, err := GetGamesForUserPaged(ctx, user.ID, 10, 10)
		assert.NoError(t, err)
		assert.Len(t, games, 5)
	})
}

func TestComplexJSONSerialization(t *testing.T) {
	SetupDBTest(t)
	ctx := context.Background()
	ctx = WithUserID(ctx, 1) // Dummy ID for RLS
	gameID := "json-complex-test"

	t.Run("Special Characters in Piece Data", func(t *testing.T) {
		err := SaveGame(ctx, gameID, nil, nil, "funky", map[string]any{
			"metadata": "some \"quoted\" text and \n newlines",
			"unicode":  "🚀 Stratego!",
		}, nil, time.Now(), time.Now())
		assert.NoError(t, err)

		move := models.HistoricalMove{
			MoveIndex: 1,
			PlayerID:  0,
			FromX:     1, FromY: 1, ToX: 1, ToY: 2,
			Attacker: &models.PieceData{
				Type: "Special \"Char\" Unit",
				Rank: "Rank-10",
			},
			Result: models.ResultMove,
		}

		err = SaveMove(ctx, gameID, move)
		assert.NoError(t, err)

		history, err := GetGameHistory(ctx, gameID)
		assert.NoError(t, err)
		assert.Equal(t, "Special \"Char\" Unit", history.Moves[0].Attacker.Type)
	})

	t.Run("Empty and Nil Data", func(t *testing.T) {
		gameID2 := "json-nil-test"
		_ = SaveGame(ctx, gameID2, nil, nil, "test", map[string]any{}, nil, time.Now(), time.Now())

		move := models.HistoricalMove{
			MoveIndex: 1,
			Result:    models.ResultMove,
			Attacker:  nil,
			Defender:  nil,
		}

		err := SaveMove(ctx, gameID2, move)
		assert.NoError(t, err)

		history, err := GetGameHistory(ctx, gameID2)
		assert.NoError(t, err)
		assert.Nil(t, history.Moves[0].Attacker)
	})
}
