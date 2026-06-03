package game

import (
	"digital-innovation/gostrategy/internal/game/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWinCondition_FlagCapture(t *testing.T) {
	p1 := NewPlayer(1, "Player 1", "red")
	p2 := NewPlayer(2, "Player 2", "blue")

	attacker := NewPiece(models.Marshal, &p1)
	flag := NewPiece(models.Flag, &p2)

	assert.False(t, p1.HasWon())

	// Capture the flag
	attacker.Attack(flag)

	assert.True(t, p1.HasWon())
	assert.False(t, flag.IsAlive())
}

func TestWinCondition_TotalImmobilization(t *testing.T) {
	p1 := NewPlayer(1, "Player 1", "red")
	board := NewBoard()

	// Place only unmovable pieces for p1
	board.SetPieceAt(NewPosition(0, 0), NewPiece(models.Flag, &p1))
	board.SetPieceAt(NewPosition(0, 1), NewPiece(models.Bomb, &p1))

	hasMovable := false
	for _, piece := range p1.GetAlivePieces() {
		if piece.CanMove() {
			hasMovable = true
			break
		}
	}

	assert.False(t, hasMovable)
}
