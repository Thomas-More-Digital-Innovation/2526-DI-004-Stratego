package minimax

import (
	"digital-innovation/gostrategy/internal/game"
	"digital-innovation/gostrategy/internal/game/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMinimaxAI(t *testing.T) {
	player1 := game.NewPlayer(0, "Piet", "red")
	aiObj := NewAI(&player1, false)
	board := game.NewBoard()

	p1 := game.NewPiece(models.Marshal, &player1)
	board.SetPieceAt(game.NewPosition(4, 4), p1)
	player1.AddPiece(p1, game.NewPosition(4, 4))

	move := aiObj.MakeMove(board)
	assert.NotEqual(t, game.Move{}, move)
	assert.Equal(t, game.NewPosition(4, 4), move.GetFrom())
}
