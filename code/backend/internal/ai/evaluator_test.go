package ai

import (
	"digital-innovation/gostrategy/internal/game"
	"digital-innovation/gostrategy/internal/game/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvaluateBoardAndSimulateMove(t *testing.T) {
	board := game.NewBoard()
	player1 := game.NewPlayer(0, "Piet", "red")
	player2 := game.NewPlayer(1, "Bob", "blue")

	p1 := game.NewPiece(models.Marshal, &player1)
	p2 := game.NewPiece(models.General, &player2)

	board.SetPieceAt(game.NewPosition(0, 0), p1)
	board.SetPieceAt(game.NewPosition(0, 1), p2)

	weights := map[string]float64{
		"Marshal": 100.0,
		"General": 80.0,
	}

	score := EvaluateBoard(board, &player1, nil, weights, 0.5)
	assert.Greater(t, score, 0.0)

	move := game.NewMove(game.NewPosition(0, 0), game.NewPosition(0, 1), &player1)
	simulated := SimulateMove(board, move)

	pieceAtDest := simulated.GetPieceAt(game.NewPosition(0, 1))
	assert.NotNil(t, pieceAtDest)
	assert.Equal(t, "Marshal", pieceAtDest.GetType().GetName())
	assert.Nil(t, simulated.GetPieceAt(game.NewPosition(0, 0)))
}
