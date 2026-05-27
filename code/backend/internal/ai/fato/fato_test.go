package fato_test

import (
	"digital-innovation/gostrategy/internal/ai/fato"
	"digital-innovation/gostrategy/internal/models"
	"digital-innovation/gostrategy/pkg/game"
	"testing"
)

func TestNewAI(t *testing.T) {
	player := game.NewPlayer(0, "player", "red")
	ai := fato.NewAI(&player, true)

	if ai.GetPlayer() == nil {
		t.Errorf("Expected player to be set in FatoAI")
	}

	if ai.GetMemory() == nil {
		t.Errorf("Expected memory to be initialized in FatoAI")
	}
}

func TestIsPieceMemorized(t *testing.T) {
	player := game.NewPlayer(0, "player", "red")
	ai := fato.NewAI(&player, true)

	position := game.NewPosition(1, 1)

	// Should not be memorized initially
	if ai.GetMemory().Recall(position) != nil {
		t.Errorf("Expected piece not to be memorized initially")
	}

	// Remember a piece
	piece := game.NewPiece(models.Scout, &player)
	ai.GetMemory().Remember(position, piece, 1.0, 1)

	// Should now be memorized
	if ai.GetMemory().Recall(position) == nil {
		t.Errorf("Expected piece to be memorized after Remember()")
	}
}

func TestMakeMove(t *testing.T) {
	p1 := game.NewPlayer(0, "ai", "red")
	p2 := game.NewPlayer(1, "human", "blue")
	ai := fato.NewAI(&p1, true)
	board := game.NewBoard()

	// Place AI piece
	pos1 := game.NewPosition(5, 5)
	piece1 := game.NewPiece(models.Marshal, &p1)
	board.SetPieceAt(pos1, piece1)
	p1.AddPiece(piece1, pos1)

	// Test case 1: No targets, should find an exploration or random move
	move := ai.MakeMove(board)
	if move.IsEmpty() {
		t.Error("Expected AI to find a move, but got empty")
	}

	// Test case 2: Visible target, should attack it
	targetPos := game.NewPosition(5, 6)
	targetPiece := game.NewPiece(models.General, &p2)
	targetPiece.Reveal()
	board.SetPieceAt(targetPos, targetPiece)
	p2.AddPiece(targetPiece, targetPos)

	move = ai.MakeMove(board)
	if move.GetTo() != targetPos {
		t.Errorf("Expected AI to attack visible target at %v, but moved to %v", targetPos, move.GetTo())
	}
}

func TestAnalyzeMove(t *testing.T) {
	aiPlayer := game.NewPlayer(0, "ai", "red")
	humanPlayer := game.NewPlayer(1, "human", "blue")
	ai := fato.NewAI(&aiPlayer, true)

	// Normal move (1 square) - should NOT be remembered as scout
	normalMove := game.NewMove(game.NewPosition(1, 1), game.NewPosition(1, 2), &humanPlayer)
	ai.AnalyzeMove(normalMove, &humanPlayer, 1)

	if ai.GetMemory().Recall(normalMove.GetTo()) != nil {
		t.Errorf("Expected piece to not be remembered after normal move")
	}

	// Scout move (2+ squares) - should be remembered as scout
	scoutMove := game.NewMove(game.NewPosition(1, 1), game.NewPosition(1, 3), &humanPlayer)
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
	aiPlayer := game.NewPlayer(0, "ai", "red")
	humanPlayer := game.NewPlayer(1, "human", "blue")
	ai := fato.NewAI(&aiPlayer, true)

	// Remember a piece at position (2, 2)
	pos1 := game.NewPosition(2, 2)
	piece := game.NewPiece(models.Captain, &humanPlayer)
	ai.GetMemory().Remember(pos1, piece, 0.9, 1)

	// Verify it's there
	if ai.GetMemory().Recall(pos1) == nil {
		t.Fatal("Expected piece to be remembered at pos1")
	}

	// Simulate opponent moving that piece from (2,2) to (2,3)
	pos2 := game.NewPosition(2, 3)
	move := game.NewMove(pos1, pos2, &humanPlayer)
	ai.AnalyzeMove(move, &humanPlayer, 2)

	// Memory should have moved from pos1 to pos2
	if ai.GetMemory().Recall(pos1) != nil {
		t.Errorf("Expected memory to be cleared at original position")
	}

	if ai.GetMemory().Recall(pos2) == nil {
		t.Errorf("Expected memory to be moved to new position")
	}
}
