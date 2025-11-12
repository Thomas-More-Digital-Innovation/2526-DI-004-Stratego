package fato_test

import (
	"digital-innovation/stratego/ai/fato"
	"digital-innovation/stratego/engine"
	"digital-innovation/stratego/models"
	"testing"
)

func TestNewFatoAI(t *testing.T) {
	player := engine.NewPlayer(0, "player", "red")
	ai := fato.NewFatoAI(&player, true)

	if ai.GetPlayer() == nil {
		t.Errorf("Expected player to be set in FatoAI")
	}

	if ai.GetMemory() == nil {
		t.Errorf("Expected memory to be initialized in FatoAI")
	}
}

func TestIsPieceMemorized(t *testing.T) {
	player := engine.NewPlayer(0, "player", "red")
	ai := fato.NewFatoAI(&player, true)

	position := engine.NewPosition(1, 1)

	// Should not be memorized initially
	if ai.GetMemory().Recall(position) != nil {
		t.Errorf("Expected piece not to be memorized initially")
	}

	// Remember a piece
	piece := engine.NewPiece(models.Scout, &player)
	ai.GetMemory().Remember(position, piece, 1.0, 1)

	// Should now be memorized
	if ai.GetMemory().Recall(position) == nil {
		t.Errorf("Expected piece to be memorized after Remember()")
	}
}

func TestMakeMove(t *testing.T) {
	// TODO
}

func TestAnalyzeMove(t *testing.T) {
	aiPlayer := engine.NewPlayer(0, "ai", "red")
	humanPlayer := engine.NewPlayer(1, "human", "blue")
	ai := fato.NewFatoAI(&aiPlayer, true)

	// Normal move (1 square) - should NOT be remembered as scout
	normalMove := engine.NewMove(engine.NewPosition(1, 1), engine.NewPosition(1, 2), &humanPlayer)
	ai.AnalyzeMove(normalMove, &humanPlayer, 1)

	if ai.GetMemory().Recall(normalMove.GetTo()) != nil {
		t.Errorf("Expected piece to not be remembered after normal move")
	}

	// Scout move (2+ squares) - should be remembered as scout
	scoutMove := engine.NewMove(engine.NewPosition(1, 1), engine.NewPosition(1, 3), &humanPlayer)
	ai.AnalyzeMove(scoutMove, &humanPlayer, 1)

	remembered := ai.GetMemory().Recall(scoutMove.GetTo())
	if remembered == nil {
		t.Errorf("Expected piece to be remembered after scout move")
	}

	if remembered != nil && remembered.Piece.GetType().GetName() != "Scout" {
		t.Errorf("Expected remembered piece to be Scout, got: %s", remembered.Piece.GetType().GetName())
	}

	if remembered != nil && remembered.Confidence != 1.0 {
		t.Errorf("Expected confidence 1.0 for scout guess, got: %.2f", remembered.Confidence)
	}
}

func TestMemoryUpdatesOnMove(t *testing.T) {
	aiPlayer := engine.NewPlayer(0, "ai", "red")
	humanPlayer := engine.NewPlayer(1, "human", "blue")
	ai := fato.NewFatoAI(&aiPlayer, true)

	// Remember a piece at position (2, 2)
	pos1 := engine.NewPosition(2, 2)
	piece := engine.NewPiece(models.Captain, &humanPlayer)
	ai.GetMemory().Remember(pos1, piece, 0.9, 1)

	// Verify it's there
	if ai.GetMemory().Recall(pos1) == nil {
		t.Fatal("Expected piece to be remembered at pos1")
	}

	// Simulate opponent moving that piece from (2,2) to (2,3)
	pos2 := engine.NewPosition(2, 3)
	move := engine.NewMove(pos1, pos2, &humanPlayer)
	ai.AnalyzeMove(move, &humanPlayer, 2)

	// Memory should have moved from pos1 to pos2
	if ai.GetMemory().Recall(pos1) != nil {
		t.Errorf("Expected memory to be cleared at original position")
	}

	if ai.GetMemory().Recall(pos2) == nil {
		t.Errorf("Expected memory to be moved to new position")
	}
}
