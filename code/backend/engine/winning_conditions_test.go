package engine_test

import (
	"digital-innovation/gostrategy/engine"
	"digital-innovation/gostrategy/models"
	"testing"
)

func TestWinCondition_FlagCapture(t *testing.T) {
	p1 := engine.NewPlayer(1, "Player 1", "red")
	p2 := engine.NewPlayer(2, "Player 2", "blue")

	attacker := engine.NewPiece(models.Marshal, &p1)
	flag := engine.NewPiece(models.Flag, &p2)

	// Before capture
	if p1.HasWon() {
		t.Error("P1 should not have won before capturing the flag")
	}

	// Capture the flag
	attacker.Attack(flag)

	if !p1.HasWon() {
		t.Error("P1 should have won after capturing the flag")
	}
	if flag.IsAlive() {
		t.Error("Flag should be eliminated after capture")
	}
}

func TestWinCondition_TotalImmobilization(t *testing.T) {
	// This usually is detected in the game loop/session by checking available moves for the current player
	p1 := engine.NewPlayer(1, "Player 1", "red")
	board := engine.NewBoard()

	// Place only unmovable pieces for p1
	board.SetPieceAt(engine.NewPosition(0, 0), engine.NewPiece(models.Flag, &p1))
	board.SetPieceAt(engine.NewPosition(0, 1), engine.NewPiece(models.Bomb, &p1))

	// Check if p1 has any movable pieces
	hasMovable := false
	for _, piece := range p1.GetAlivePieces() {
		if piece.CanMove() {
			hasMovable = true
			break
		}
	}

	if hasMovable {
		t.Error("Player should have no movable pieces")
	}
}
