package heuristic

import (
	"digital-innovation/gostrategy/internal/ai"
	"digital-innovation/gostrategy/internal/game"
	"digital-innovation/gostrategy/internal/game/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAI(t *testing.T) {
	player := game.NewPlayer(0, "Alice", "red")
	aiObj := NewAI(&player, false)
	assert.NotNil(t, aiObj)
	assert.Equal(t, &player, aiObj.GetPlayer())
}

func TestNewAIWithParams(t *testing.T) {
	player := game.NewPlayer(0, "Alice", "red")
	params := &ai.Parameters{
		Weights:    map[string]float64{"Marshal": 100.0},
		Aggression: 0.5,
	}
	aiObj := NewAIWithParams(&player, true, params)
	assert.NotNil(t, aiObj)
	assert.Equal(t, params, aiObj.params)
	assert.NotNil(t, aiObj.GetMemory())
}

func TestHeuristicAI_NoMoves(t *testing.T) {
	player1 := game.NewPlayer(0, "Alice", "red")
	aiObj := NewAI(&player1, false)
	board := game.NewBoard()

	move := aiObj.MakeMove(board)
	assert.Equal(t, game.Move{}, move)
}

func TestHeuristicAI_EdgeCases(t *testing.T) {
	player1 := game.NewPlayer(0, "Alice", "red")
	aiObj := NewAI(&player1, false)
	board := game.NewBoard()

	pBomb := game.NewPiece(models.Bomb, &player1)
	board.SetPieceAt(game.NewPosition(0, 0), pBomb)
	player1.AddPiece(pBomb, game.NewPosition(0, 0))

	move := aiObj.MakeMove(board)
	assert.Equal(t, game.Move{}, move)

	board.SetPieceAt(game.NewPosition(0, 0), nil)
	move = aiObj.MakeMove(board)
	assert.Equal(t, game.Move{}, move)

	player1.RemovePiece(pBomb)
	pMarshal := game.NewPiece(models.Marshal, &player1)
	player1.AddPiece(pMarshal, game.NewPosition(0, 0))
	board.SetPieceAt(game.NewPosition(0, 0), nil)

	move = aiObj.MakeMove(board)
	assert.Equal(t, game.Move{}, move)
}
