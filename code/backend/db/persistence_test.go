package db

import (
	"context"
	"digital-innovation/gostrategy/models"
	"fmt"
	"testing"
)

func TestGamePagination(t *testing.T) {
	testDB := SetupTestDB(t)
	oldDB := DB
	DB = testDB
	defer func() { DB = oldDB }()

	ctx := context.Background()
	user, _ := CreateUser(ctx, "paginator", "Pass1234!", "")

	// Create 15 games
	for i := range 15 {
		gameID := fmt.Sprintf("paged-game-%d", i)
		err := SaveGame(ctx, gameID, &user.ID, nil, "ranked", map[string]string{"test": "data"}, nil)
		if err != nil {
			t.Fatalf("failed to save game %d: %v", i, err)
		}
	}

	t.Run("Count Games", func(t *testing.T) {
		count, err := GetGamesCountForUser(ctx, user.ID)
		if err != nil {
			t.Errorf("GetGamesCountForUser failed: %v", err)
		}
		if count != 15 {
			t.Errorf("expected 15 games, got %d", count)
		}
	})

	t.Run("Paged Retrieval Page 1", func(t *testing.T) {
		games, err := GetGamesForUserPaged(ctx, user.ID, 10, 0)
		if err != nil {
			t.Errorf("GetGamesForUserPaged failed: %v", err)
		}
		if len(games) != 10 {
			t.Errorf("expected 10 games, got %d", len(games))
		}
	})

	t.Run("Paged Retrieval Page 2", func(t *testing.T) {
		games, err := GetGamesForUserPaged(ctx, user.ID, 10, 10)
		if err != nil {
			t.Errorf("GetGamesForUserPaged failed: %v", err)
		}
		if len(games) != 5 {
			t.Errorf("expected 5 games, got %d", len(games))
		}
	})
}

func TestComplexJSONSerialization(t *testing.T) {
	testDB := SetupTestDB(t)
	oldDB := DB
	DB = testDB
	defer func() { DB = oldDB }()

	ctx := context.Background()
	gameID := "json-complex-test"

	t.Run("Special Characters in Piece Data", func(t *testing.T) {
		err := SaveGame(ctx, gameID, nil, nil, "funky", map[string]any{
			"metadata": "some \"quoted\" text and \n newlines",
			"unicode":  "🚀 Stratego!",
		}, nil)
		if err != nil {
			t.Fatalf("SaveGame with complex JSON failed: %v", err)
		}

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

		if err := SaveMove(ctx, gameID, move); err != nil {
			t.Errorf("SaveMove with complex JSON failed: %v", err)
		}

		history, err := GetGameHistory(ctx, gameID)
		if err != nil {
			t.Fatalf("GetGameHistory failed: %v", err)
		}

		if history.Moves[0].Attacker.Type != "Special \"Char\" Unit" {
			t.Errorf("expected special char type, got %s", history.Moves[0].Attacker.Type)
		}
	})

	t.Run("Empty and Nil Data", func(t *testing.T) {
		gameID2 := "json-nil-test"
		_ = SaveGame(ctx, gameID2, nil, nil, "test", map[string]any{}, nil)

		move := models.HistoricalMove{
			MoveIndex: 1,
			Result:    models.ResultMove,
			Attacker:  nil,
			Defender:  nil,
		}

		if err := SaveMove(ctx, gameID2, move); err != nil {
			t.Errorf("SaveMove with nil data failed: %v", err)
		}

		history, err := GetGameHistory(ctx, gameID2)
		if err != nil {
			t.Fatalf("GetGameHistory failed: %v", err)
		}

		if history.Moves[0].Attacker != nil {
			t.Errorf("expected nil attacker, got %+v", history.Moves[0].Attacker)
		}
	})
}
