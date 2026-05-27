package game

import (
	"digital-innovation/gostrategy/pkg/game/models"
	"testing"
)

func TestWinCondition_FlagCapture(t *testing.T) {
	p1 := NewPlayer(1, "Player 1", "red")
	p2 := NewPlayer(2, "Player 2", "blue")

	attacker := NewPiece(models.Marshal, &p1)
	flag := NewPiece(models.Flag, &p2)

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
	p1 := NewPlayer(1, "Player 1", "red")
	board := NewBoard()

	// Place only unmovable pieces for p1
	board.SetPieceAt(NewPosition(0, 0), NewPiece(models.Flag, &p1))
	board.SetPieceAt(NewPosition(0, 1), NewPiece(models.Bomb, &p1))

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
