package mcts

import (
	"digital-innovation/gostrategy/internal/ai"
	"digital-innovation/gostrategy/internal/game"
	"digital-innovation/gostrategy/internal/game/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMctsAI_NoMoves(t *testing.T) {
	player1 := game.NewPlayer(0, "us", "red")
	aiObj := NewAI(&player1, false)
	board := game.NewBoard()

	move := aiObj.MakeMove(board)
	assert.Equal(t, game.Move{}, move)
}

func TestMctsAI_ConstructorAndParams(t *testing.T) {
	player1 := game.NewPlayer(0, "us", "red")
	params := &ai.Parameters{
		Weights: map[string]float64{
			"Marshal": 100.0,
		},
		Aggression: 0.5,
		Config: map[string]interface{}{
			"iterations":           float64(5),
			"exploration_constant": 1.414,
		},
	}

	aiObj := NewAIWithParams(&player1, true, params)
	assert.NotNil(t, aiObj)
	assert.Equal(t, params, aiObj.params)
}

func TestMctsAI_RolloutAndDecisions(t *testing.T) {
	player1 := game.NewPlayer(0, "us", "red")
	player2 := game.NewPlayer(1, "them", "blue")
	params := &ai.Parameters{
		Weights: map[string]float64{
			"Marshal": 100.0,
			"Flag":    10000.0,
		},
		Aggression: 0.5,
		Config: map[string]interface{}{
			"iterations": float64(2),
		},
	}

	t.Run("Flag capture rollout", func(t *testing.T) {
		aiObj := NewAIWithParams(&player1, false, params)
		board := game.NewBoard()

		p1Flag := game.NewPiece(models.Flag, &player1)
		p2Flag := game.NewPiece(models.Flag, &player2)
		p1Marshal := game.NewPiece(models.Marshal, &player1)

		board.SetPieceAt(game.NewPosition(0, 0), p1Flag)
		board.SetPieceAt(game.NewPosition(9, 9), p2Flag)
		board.SetPieceAt(game.NewPosition(0, 1), p1Marshal)

		player1.AddPiece(p1Flag, game.NewPosition(0, 0))
		player2.AddPiece(p2Flag, game.NewPosition(9, 9))
		player1.AddPiece(p1Marshal, game.NewPosition(0, 1))

		move := aiObj.MakeMove(board)
		assert.NotEqual(t, game.Move{}, move)
	})

	t.Run("rollout methods directly", func(t *testing.T) {
		aiObj := NewAIWithParams(&player1, false, params)
		board := game.NewBoard()

		p1Flag := game.NewPiece(models.Flag, &player1)
		p2Flag := game.NewPiece(models.Flag, &player2)
		board.SetPieceAt(game.NewPosition(0, 0), p1Flag)
		board.SetPieceAt(game.NewPosition(0, 1), p2Flag)

		player1.AddPiece(p1Flag, game.NewPosition(0, 0))
		player2.AddPiece(p2Flag, game.NewPosition(0, 1))

		assert.Equal(t, game.NewPosition(-1, -1), aiObj.findFlagPosition(board, nil))
		assert.Equal(t, game.NewPosition(0, 0), aiObj.findFlagPosition(board, &player1))

		boardNoFlag := game.NewBoard()
		assert.Equal(t, game.NewPosition(-1, -1), aiObj.findFlagPosition(boardNoFlag, &player1))

		assert.True(t, aiObj.isFlagCapturedAt(board, game.NewPosition(-1, -1), &player1))
		assert.False(t, aiObj.isFlagCapturedAt(board, game.NewPosition(0, 0), &player1))
		assert.True(t, aiObj.isFlagCapturedAt(board, game.NewPosition(0, 1), &player1))
		assert.True(t, aiObj.isFlagCapturedAt(boardNoFlag, game.NewPosition(0, 0), &player1))
	})
}

func TestMctsAI_SimulateMoveInPlace(t *testing.T) {
	player1 := game.NewPlayer(0, "us", "red")
	player2 := game.NewPlayer(1, "them", "blue")
	aiObj := NewAI(&player1, false)

	t.Run("nil attacker", func(t *testing.T) {
		b := game.NewBoard()
		captured := aiObj.applySimulatedMoveInPlace(b, game.NewMove(game.NewPosition(0, 0), game.NewPosition(0, 1), &player1))
		assert.False(t, captured)
	})

	t.Run("nil target", func(t *testing.T) {
		b := game.NewBoard()
		p1 := game.NewPiece(models.General, &player1)
		b.SetPieceAt(game.NewPosition(0, 0), p1)
		captured := aiObj.applySimulatedMoveInPlace(b, game.NewMove(game.NewPosition(0, 0), game.NewPosition(0, 1), &player1))
		assert.False(t, captured)
		assert.Equal(t, p1, b.GetPieceAt(game.NewPosition(0, 1)))
		assert.Nil(t, b.GetPieceAt(game.NewPosition(0, 0)))
	})

	t.Run("target is flag", func(t *testing.T) {
		b := game.NewBoard()
		p1 := game.NewPiece(models.General, &player1)
		p2 := game.NewPiece(models.Flag, &player2)
		b.SetPieceAt(game.NewPosition(0, 0), p1)
		b.SetPieceAt(game.NewPosition(0, 1), p2)
		captured := aiObj.applySimulatedMoveInPlace(b, game.NewMove(game.NewPosition(0, 0), game.NewPosition(0, 1), &player1))
		assert.False(t, captured)
		assert.Equal(t, p1, b.GetPieceAt(game.NewPosition(0, 1)))
	})

	t.Run("spy vs marshal", func(t *testing.T) {
		b := game.NewBoard()
		p1 := game.NewPiece(models.Spy, &player1)
		p2 := game.NewPiece(models.Marshal, &player2)
		b.SetPieceAt(game.NewPosition(0, 0), p1)
		b.SetPieceAt(game.NewPosition(0, 1), p2)
		captured := aiObj.applySimulatedMoveInPlace(b, game.NewMove(game.NewPosition(0, 0), game.NewPosition(0, 1), &player1))
		assert.False(t, captured)
		assert.Equal(t, p1, b.GetPieceAt(game.NewPosition(0, 1)))
	})

	t.Run("miner vs bomb", func(t *testing.T) {
		b := game.NewBoard()
		p1 := game.NewPiece(models.Miner, &player1)
		p2 := game.NewPiece(models.Bomb, &player2)
		b.SetPieceAt(game.NewPosition(0, 0), p1)
		b.SetPieceAt(game.NewPosition(0, 1), p2)
		captured := aiObj.applySimulatedMoveInPlace(b, game.NewMove(game.NewPosition(0, 0), game.NewPosition(0, 1), &player1))
		assert.False(t, captured)
		assert.Equal(t, p1, b.GetPieceAt(game.NewPosition(0, 1)))
	})

	t.Run("non-miner vs bomb", func(t *testing.T) {
		b := game.NewBoard()
		p1 := game.NewPiece(models.General, &player1)
		p2 := game.NewPiece(models.Bomb, &player2)
		b.SetPieceAt(game.NewPosition(0, 0), p1)
		b.SetPieceAt(game.NewPosition(0, 1), p2)
		captured := aiObj.applySimulatedMoveInPlace(b, game.NewMove(game.NewPosition(0, 0), game.NewPosition(0, 1), &player1))
		assert.True(t, captured)
		assert.Nil(t, b.GetPieceAt(game.NewPosition(0, 0)))
		assert.Equal(t, p2, b.GetPieceAt(game.NewPosition(0, 1)))
	})

	t.Run("attacker vs defender ranks", func(t *testing.T) {
		b := game.NewBoard()
		p1 := game.NewPiece(models.General, &player1)
		p2 := game.NewPiece(models.Miner, &player2)
		b.SetPieceAt(game.NewPosition(0, 0), p1)
		b.SetPieceAt(game.NewPosition(0, 1), p2)
		captured := aiObj.applySimulatedMoveInPlace(b, game.NewMove(game.NewPosition(0, 0), game.NewPosition(0, 1), &player1))
		assert.False(t, captured)
		assert.Equal(t, p1, b.GetPieceAt(game.NewPosition(0, 1)))

		b = game.NewBoard()
		p1 = game.NewPiece(models.Miner, &player1)
		p2 = game.NewPiece(models.General, &player2)
		b.SetPieceAt(game.NewPosition(0, 0), p1)
		b.SetPieceAt(game.NewPosition(0, 1), p2)
		captured = aiObj.applySimulatedMoveInPlace(b, game.NewMove(game.NewPosition(0, 0), game.NewPosition(0, 1), &player1))
		assert.True(t, captured)
		assert.Nil(t, b.GetPieceAt(game.NewPosition(0, 0)))
		assert.Equal(t, p2, b.GetPieceAt(game.NewPosition(0, 1)))

		b = game.NewBoard()
		p1 = game.NewPiece(models.Miner, &player1)
		p2 = game.NewPiece(models.Miner, &player2)
		b.SetPieceAt(game.NewPosition(0, 0), p1)
		b.SetPieceAt(game.NewPosition(0, 1), p2)
		captured = aiObj.applySimulatedMoveInPlace(b, game.NewMove(game.NewPosition(0, 0), game.NewPosition(0, 1), &player1))
		assert.True(t, captured)
		assert.Nil(t, b.GetPieceAt(game.NewPosition(0, 0)))
		assert.Nil(t, b.GetPieceAt(game.NewPosition(0, 1)))
	})
}

func TestMctsAI_UpdateIndex(t *testing.T) {
	idx := []game.Position{
		game.NewPosition(0, 0),
		game.NewPosition(1, 1),
	}

	updateIndex(&idx, game.NewPosition(0, 0), game.NewPosition(0, 1), false)
	assert.Len(t, idx, 2)
	assert.Contains(t, idx, game.NewPosition(0, 1))

	updateIndex(&idx, game.NewPosition(0, 1), game.NewPosition(0, 2), true)
	assert.Len(t, idx, 1)
	assert.NotContains(t, idx, game.NewPosition(0, 1))
	assert.NotContains(t, idx, game.NewPosition(0, 2))

	updateIndex(&idx, game.NewPosition(5, 5), game.NewPosition(5, 6), false)
	assert.Len(t, idx, 1)
}

func BenchmarkMctsAI(b *testing.B) {
	player1 := game.NewPlayer(0, "Piet", "red")
	aiObj := NewAI(&player1, false)
	board := game.NewBoard()

	p1 := game.NewPiece(models.Marshal, &player1)
	board.SetPieceAt(game.NewPosition(4, 4), p1)
	player1.AddPiece(p1, game.NewPosition(4, 4))

	for b.Loop() {
		_ = aiObj.MakeMove(board)
	}
}
