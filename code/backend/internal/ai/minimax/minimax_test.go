package minimax

import (
	"digital-innovation/gostrategy/internal/ai"
	ai_const "digital-innovation/gostrategy/internal/ai/const"
	"digital-innovation/gostrategy/internal/game"
	"digital-innovation/gostrategy/internal/game/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMinimaxAI_NoMoves(t *testing.T) {
	player1 := game.NewPlayer(0, "us", "red")
	aiObj := NewAI(&player1, false)
	board := game.NewBoard()

	move := aiObj.MakeMove(board)
	assert.Equal(t, game.Move{}, move)
}

func TestMinimaxAI_ConstructorAndParams(t *testing.T) {
	player1 := game.NewPlayer(0, "us", "red")
	params := &ai.Parameters{
		Weights: map[string]float64{
			ai_const.Marshal: 100.0,
		},
		Aggression: 0.5,
		Config: map[string]interface{}{
			"depth": float64(3),
		},
	}

	aiObj := NewAIWithParams(&player1, true, params)
	assert.NotNil(t, aiObj)
	assert.Equal(t, params, aiObj.params)
	assert.True(t, aiObj.GetMemory() != nil)
}

func TestMinimaxAI_MinimaxBranches(t *testing.T) {
	player1 := game.NewPlayer(0, "us", "red")
	player2 := game.NewPlayer(1, "them", "blue")
	params := &ai.Parameters{
		Weights: map[string]float64{
			ai_const.Marshal: 100.0,
			ai_const.General: 80.0,
		},
		Aggression: 0.0,
		Config: map[string]interface{}{
			"depth": float64(2),
		},
	}

	t.Run("basic search, max player moves, opponent moves, pruning", func(t *testing.T) {
		aiObj := NewAIWithParams(&player1, false, params)
		board := game.NewBoard()

		p1 := game.NewPiece(models.Marshal, &player1)
		p2 := game.NewPiece(models.General, &player2)
		board.SetPieceAt(game.NewPosition(0, 0), p1)
		board.SetPieceAt(game.NewPosition(0, 1), p2)

		player1.AddPiece(p1, game.NewPosition(0, 0))
		player2.AddPiece(p2, game.NewPosition(0, 1))

		move := aiObj.MakeMove(board)
		assert.NotEqual(t, game.Move{}, move)
		assert.Equal(t, game.NewPosition(0, 0), move.GetFrom())
		assert.Equal(t, game.NewPosition(0, 1), move.GetTo())
	})

	t.Run("no opponent - should return evaluation", func(t *testing.T) {
		aiObj := NewAIWithParams(&player1, false, params)
		board := game.NewBoard()
		p1 := game.NewPiece(models.Marshal, &player1)
		board.SetPieceAt(game.NewPosition(0, 0), p1)
		player1.AddPiece(p1, game.NewPosition(0, 0))

		score := aiObj.minimax(board, 1, -1e9, 1e9, false, nil)
		assert.Greater(t, score, 0.0)
	})

	t.Run("transposition table hits", func(t *testing.T) {
		aiObj := NewAIWithParams(&player1, false, params)
		board := game.NewBoard()
		p1 := game.NewPiece(models.Marshal, &player1)
		board.SetPieceAt(game.NewPosition(0, 0), p1)
		player1.AddPiece(p1, game.NewPosition(0, 0))

		aiObj.tt = make(map[string]ttEntry)
		key := getBoardStateKey(board, true)

		aiObj.tt[key] = ttEntry{score: 45.0, depth: 2, flag: ttExact}
		valExact := aiObj.minimax(board, 2, -1e9, 1e9, true, nil)
		assert.Equal(t, 45.0, valExact)

		aiObj.tt[key] = ttEntry{score: 10.0, depth: 2, flag: ttAlpha}
		aiObj.tt[key] = ttEntry{score: -100.0, depth: 2, flag: ttAlpha}
		valAlpha := aiObj.minimax(board, 2, -50.0, 50.0, true, nil)
		assert.Equal(t, -100.0, valAlpha)

		aiObj.tt[key] = ttEntry{score: 100.0, depth: 2, flag: ttBeta}
		valBeta := aiObj.minimax(board, 2, -50.0, 50.0, true, nil)
		assert.Equal(t, 100.0, valBeta)
	})

	t.Run("maximize player with no moves", func(t *testing.T) {
		aiObj := NewAIWithParams(&player1, false, params)
		board := game.NewBoard()
		aiObj.tt = make(map[string]ttEntry)
		score := aiObj.minimax(board, 1, -1e9, 1e9, true, nil)
		assert.Equal(t, -1e9, score)
	})

	t.Run("minimize player with no moves", func(t *testing.T) {
		aiObj := NewAIWithParams(&player1, false, params)
		board := game.NewBoard()
		aiObj.tt = make(map[string]ttEntry)
		score := aiObj.minimax(board, 1, -1e9, 1e9, false, &player2)
		assert.Equal(t, 1e9, score)
	})

	t.Run("alpha-beta pruning cutoff assignment", func(t *testing.T) {
		aiObj := NewAIWithParams(&player1, false, params)
		board := game.NewBoard()

		p1 := game.NewPiece(models.Marshal, &player1)
		board.SetPieceAt(game.NewPosition(0, 0), p1)
		player1.AddPiece(p1, game.NewPosition(0, 0))

		aiObj.tt = make(map[string]ttEntry)
		_ = aiObj.minimax(board, 1, -1e9, -500.0, true, nil)
		key := getBoardStateKey(board, true)
		entry, ok := aiObj.tt[key]
		assert.True(t, ok)
		assert.Equal(t, ttBeta, entry.flag)
	})
}
