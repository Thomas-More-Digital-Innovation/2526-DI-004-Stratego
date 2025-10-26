package game_test

import (
	"digital-innovation/stratego/engine"
	"digital-innovation/stratego/game"
	"testing"
)

func TestSetupGame(t *testing.T) {
	player1 := engine.NewPlayer(0, "Player 1", "avatar1")
	player2 := engine.NewPlayer(1, "Player 2", "avatar2")

	controller1 := engine.NewHumanPlayerController(&player1)
	controller2 := engine.NewHumanPlayerController(&player2)

	g := game.NewGame(controller1, controller2)

	player1Pieces := game.GetPieceList(&player1)
	player2Pieces := game.GetPieceList(&player2)

	err := game.SetupGame(g, player1Pieces, player2Pieces)
	if err != nil {
		t.Fatalf("SetupGame failed: %v", err)
	}

	// Verify piece counts
	if len(player1.GetAlivePieces()) != 40 {
		t.Errorf("Player 1 should have 40 pieces, got %d", len(player1.GetAlivePieces()))
	}

	if len(player2.GetAlivePieces()) != 40 {
		t.Errorf("Player 2 should have 40 pieces, got %d", len(player2.GetAlivePieces()))
	}

	// Verify pieces are placed in correct rows
	field := g.Board.GetField()

	// Player 1 should be in rows 6-9
	for y := 6; y <= 9; y++ {
		for x := 0; x < 10; x++ {
			piece := field[y][x]
			if piece == nil {
				t.Errorf("Expected piece at (%d, %d) for player 1", x, y)
			} else if piece.GetOwner().GetID() != 0 {
				t.Errorf("Expected player 1 piece at (%d, %d), got player %d", x, y, piece.GetOwner().GetID())
			}
		}
	}

	// Player 2 should be in rows 0-3
	for y := 0; y <= 3; y++ {
		for x := 0; x < 10; x++ {
			piece := field[y][x]
			if piece == nil {
				t.Errorf("Expected piece at (%d, %d) for player 2", x, y)
			} else if piece.GetOwner().GetID() != 1 {
				t.Errorf("Expected player 2 piece at (%d, %d), got player %d", x, y, piece.GetOwner().GetID())
			}
		}
	}

	// Middle rows (4-5) should be empty (except lakes)
	for y := 4; y <= 5; y++ {
		for x := 0; x < 10; x++ {
			piece := field[y][x]
			if piece != nil {
				t.Errorf("Expected no piece at (%d, %d), got piece", x, y)
			}
		}
	}

	// Verify piece scores are initialized
	if player1.GetPieceScore() == 0 {
		t.Error("Player 1 piece score should be initialized")
	}

	if player2.GetPieceScore() == 0 {
		t.Error("Player 2 piece score should be initialized")
	}
}

func TestSetupGame_InvalidPieceCount(t *testing.T) {
	player1 := engine.NewPlayer(0, "Player 1", "avatar1")
	player2 := engine.NewPlayer(1, "Player 2", "avatar2")

	controller1 := engine.NewHumanPlayerController(&player1)
	controller2 := engine.NewHumanPlayerController(&player2)

	g := game.NewGame(controller1, controller2)

	player1Pieces := game.GetPieceList(&player1)
	player2Pieces := game.GetPieceList(&player2)[:30] // Only 30 pieces

	err := game.SetupGame(g, player1Pieces, player2Pieces)
	if err == nil {
		t.Fatal("Expected error for invalid piece count, got nil")
	}

	if err.Error() != "each player must have exactly 40 pieces" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestRandomSetup(t *testing.T) {
	player := engine.NewPlayer(0, "Test Player", "avatar")

	pieces1 := game.RandomSetup(&player)
	pieces2 := game.RandomSetup(&player)

	// Should have 40 pieces
	if len(pieces1) != 40 {
		t.Errorf("Expected 40 pieces, got %d", len(pieces1))
	}

	// Two random setups should likely be different (not guaranteed but highly probable)
	different := false
	for i := 0; i < 40; i++ {
		if pieces1[i].GetType().GetName() != pieces2[i].GetType().GetName() {
			different = true
			break
		}
	}

	if !different {
		t.Log("Warning: Two random setups are identical (low probability but possible)")
	}
}

func TestQuickStart(t *testing.T) {
	player1 := engine.NewPlayer(0, "Player 1", "avatar1")
	player2 := engine.NewPlayer(1, "Player 2", "avatar2")

	controller1 := engine.NewHumanPlayerController(&player1)
	controller2 := engine.NewHumanPlayerController(&player2)

	g := game.QuickStart(controller1, controller2)

	if g == nil {
		t.Fatal("QuickStart returned nil game")
	}

	// Verify game is properly initialized
	if len(player1.GetAlivePieces()) != 40 {
		t.Errorf("Player 1 should have 40 pieces, got %d", len(player1.GetAlivePieces()))
	}

	if len(player2.GetAlivePieces()) != 40 {
		t.Errorf("Player 2 should have 40 pieces, got %d", len(player2.GetAlivePieces()))
	}

	// Verify current player is set
	if g.CurrentPlayer == nil {
		t.Error("Current player should be set")
	}

	// Verify round is initialized
	if g.GetRound() != 1 {
		t.Errorf("Expected round 1, got %d", g.GetRound())
	}

	// Verify no winner yet
	if g.GetWinner() != nil {
		t.Error("Game should not have a winner at start")
	}
}

func TestPiecePositionTracking(t *testing.T) {
	player1 := engine.NewPlayer(0, "Player 1", "avatar1")
	player2 := engine.NewPlayer(1, "Player 2", "avatar2")

	controller1 := engine.NewHumanPlayerController(&player1)
	controller2 := engine.NewHumanPlayerController(&player2)

	game.QuickStart(controller1, controller2)

	// Verify all pieces have tracked positions
	for _, piece := range player1.GetAlivePieces() {
		pos, exists := player1.GetPiecePosition(piece)
		if !exists {
			t.Error("Piece position not tracked for player 1")
		}

		// Verify position is in valid range for player 1 (rows 6-9)
		if pos.Y < 6 || pos.Y > 9 {
			t.Errorf("Player 1 piece at invalid row %d", pos.Y)
		}
	}

	for _, piece := range player2.GetAlivePieces() {
		pos, exists := player2.GetPiecePosition(piece)
		if !exists {
			t.Error("Piece position not tracked for player 2")
		}

		// Verify position is in valid range for player 2 (rows 0-3)
		if pos.Y < 0 || pos.Y > 3 {
			t.Errorf("Player 2 piece at invalid row %d", pos.Y)
		}
	}
}

func TestPieceListConsistency(t *testing.T) {
	player := engine.NewPlayer(0, "Test Player", "avatar")

	pieces := game.GetPieceList(&player)

	// Verify piece list has correct composition
	// 1 Flag, 6 Bombs, 1 Spy, 8 Scouts, 5 Miners, 4 Sergeants,
	// 4 Lieutenants, 4 Captains, 3 Majors, 2 Colonels, 1 General, 1 Marshal

	counts := make(map[string]int)
	for _, piece := range pieces {
		counts[piece.GetType().GetName()]++
	}

	expected := map[string]int{

		"Flag":       1,
		"Bomb":       6,
		"Spy":        1,
		"Scout":      8,
		"Miner":      5,
		"Sergeant":   4,
		"Lieutenant": 4,
		"Captain":    4,
		"Major":      3,
		"Colonel":    2,
		"General":    1,
		"Marshal":    1,
	}

	for name, expectedCount := range expected {
		if counts[name] != expectedCount {
			t.Errorf("Expected %d %s, got %d", expectedCount, name, counts[name])
		}
	}
}
