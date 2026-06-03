package game

import (
	"digital-innovation/gostrategy/internal/game/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGame(t *testing.T) {
	g, player1, _ := setupTestGame()

	assert.NotNil(t, g)
	assert.Len(t, g.Players, 2)
	assert.Equal(t, player1, g.CurrentPlayer)
	assert.NotNil(t, g.Board)
}

func TestNextTurn(t *testing.T) {
	g, player1, player2 := setupTestGame()

	assert.Equal(t, player1, g.CurrentPlayer)

	g.NextTurn()
	assert.Equal(t, player2, g.CurrentPlayer)

	g.NextTurn()
	assert.Equal(t, player1, g.CurrentPlayer)
}

func TestGetRound(t *testing.T) {
	g, player1, player2 := setupTestGame()

	piece1 := NewPiece(models.Scout, player1)
	piece2 := NewPiece(models.Major, player2)

	g.Board.SetPieceAt(NewPosition(0, 0), piece1)
	g.Board.SetPieceAt(NewPosition(1, 0), piece2)

	assert.Equal(t, 1, g.GetRound())

	move1 := NewMove(NewPosition(0, 0), NewPosition(0, 4), player1)
	g.MakeMove(&move1, piece1)

	assert.Equal(t, 1, g.GetRound())

	move2 := NewMove(NewPosition(1, 0), NewPosition(1, 1), player2)
	g.MakeMove(&move2, piece2)

	assert.Equal(t, 2, g.GetRound())
}

func TestMakeMoveToEmptyCell(t *testing.T) {
	g, player1, player2 := setupTestGame()

	piece := NewPiece(models.Major, player1)
	move := NewMove(NewPosition(0, 0), NewPosition(0, 1), player1)

	g.MakeMove(&move, piece)

	assert.Equal(t, piece, g.Board.GetPieceAt(move.GetTo()))
	assert.Nil(t, g.Board.GetPieceAt(move.GetFrom()))
	assert.Equal(t, player2, g.CurrentPlayer)
}

func TestMakeMoveWithWinningAttack(t *testing.T) {
	g, player1, player2 := setupTestGame()

	attacker := NewPiece(models.Captain, player1)
	defender := NewPiece(models.Scout, player2)

	g.Board.SetPieceAt(NewPosition(0, 0), attacker)
	g.Board.SetPieceAt(NewPosition(0, 1), defender)
	move := NewMove(NewPosition(0, 0), NewPosition(0, 1), player1)

	g.MakeMove(&move, attacker)

	assert.Equal(t, attacker, g.Board.GetPieceAt(move.GetTo()))
	assert.Nil(t, g.Board.GetPieceAt(move.GetFrom()))
	assert.True(t, attacker.IsAlive())
	assert.False(t, defender.IsAlive())
	assert.Equal(t, player2, g.CurrentPlayer)
}

func TestMakeMoveWithLosingAttack(t *testing.T) {
	g, player1, player2 := setupTestGame()

	attacker := NewPiece(models.Scout, player1)
	defender := NewPiece(models.Captain, player2)

	g.Board.SetPieceAt(NewPosition(0, 0), attacker)
	g.Board.SetPieceAt(NewPosition(0, 1), defender)
	move := NewMove(NewPosition(0, 0), NewPosition(0, 1), player1)

	g.MakeMove(&move, attacker)

	assert.Equal(t, defender, g.Board.GetPieceAt(move.GetTo()))
	assert.Nil(t, g.Board.GetPieceAt(move.GetFrom()))
	assert.False(t, attacker.IsAlive())
	assert.True(t, defender.IsAlive())
	assert.Equal(t, player2, g.CurrentPlayer)
}
