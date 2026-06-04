package ai

import (
	"digital-innovation/gostrategy/internal/game"
	"digital-innovation/gostrategy/internal/game/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimulateMove(t *testing.T) {
	player1 := game.NewPlayer(0, "us", "red")
	player2 := game.NewPlayer(1, "them", "blue")

	t.Run("attacker is nil", func(t *testing.T) {
		board := game.NewBoard()
		move := game.NewMove(game.NewPosition(0, 0), game.NewPosition(0, 1), &player1)
		res := SimulateMove(board, move)
		assert.Nil(t, res.GetPieceAt(game.NewPosition(0, 0)))
		assert.Nil(t, res.GetPieceAt(game.NewPosition(0, 1)))
	})

	t.Run("attacker captures flag", func(t *testing.T) {
		board := game.NewBoard()
		p1 := game.NewPiece(models.General, &player1)
		p2 := game.NewPiece(models.Flag, &player2)
		board.SetPieceAt(game.NewPosition(0, 0), p1)
		board.SetPieceAt(game.NewPosition(0, 1), p2)

		move := game.NewMove(game.NewPosition(0, 0), game.NewPosition(0, 1), &player1)
		res := SimulateMove(board, move)
		assert.Nil(t, res.GetPieceAt(game.NewPosition(0, 0)))
		assert.Equal(t, p1, res.GetPieceAt(game.NewPosition(0, 1)))
	})

	t.Run("spy attacks marshal", func(t *testing.T) {
		board := game.NewBoard()
		p1 := game.NewPiece(models.Spy, &player1)
		p2 := game.NewPiece(models.Marshal, &player2)
		board.SetPieceAt(game.NewPosition(0, 0), p1)
		board.SetPieceAt(game.NewPosition(0, 1), p2)

		move := game.NewMove(game.NewPosition(0, 0), game.NewPosition(0, 1), &player1)
		res := SimulateMove(board, move)
		assert.Nil(t, res.GetPieceAt(game.NewPosition(0, 0)))
		assert.Equal(t, p1, res.GetPieceAt(game.NewPosition(0, 1)))
	})

	t.Run("miner attacks bomb", func(t *testing.T) {
		board := game.NewBoard()
		p1 := game.NewPiece(models.Miner, &player1)
		p2 := game.NewPiece(models.Bomb, &player2)
		board.SetPieceAt(game.NewPosition(0, 0), p1)
		board.SetPieceAt(game.NewPosition(0, 1), p2)

		move := game.NewMove(game.NewPosition(0, 0), game.NewPosition(0, 1), &player1)
		res := SimulateMove(board, move)
		assert.Nil(t, res.GetPieceAt(game.NewPosition(0, 0)))
		assert.Equal(t, p1, res.GetPieceAt(game.NewPosition(0, 1)))
	})

	t.Run("non-miner attacks bomb", func(t *testing.T) {
		board := game.NewBoard()
		p1 := game.NewPiece(models.General, &player1)
		p2 := game.NewPiece(models.Bomb, &player2)
		board.SetPieceAt(game.NewPosition(0, 0), p1)
		board.SetPieceAt(game.NewPosition(0, 1), p2)

		move := game.NewMove(game.NewPosition(0, 0), game.NewPosition(0, 1), &player1)
		res := SimulateMove(board, move)
		assert.Nil(t, res.GetPieceAt(game.NewPosition(0, 0)))
		assert.Equal(t, p2, res.GetPieceAt(game.NewPosition(0, 1)))
	})

	t.Run("higher rank wins", func(t *testing.T) {
		board := game.NewBoard()
		p1 := game.NewPiece(models.General, &player1)
		p2 := game.NewPiece(models.Miner, &player2)
		board.SetPieceAt(game.NewPosition(0, 0), p1)
		board.SetPieceAt(game.NewPosition(0, 1), p2)

		move := game.NewMove(game.NewPosition(0, 0), game.NewPosition(0, 1), &player1)
		res := SimulateMove(board, move)
		assert.Nil(t, res.GetPieceAt(game.NewPosition(0, 0)))
		assert.Equal(t, p1, res.GetPieceAt(game.NewPosition(0, 1)))
	})

	t.Run("lower rank loses", func(t *testing.T) {
		board := game.NewBoard()
		p1 := game.NewPiece(models.Miner, &player1)
		p2 := game.NewPiece(models.General, &player2)
		board.SetPieceAt(game.NewPosition(0, 0), p1)
		board.SetPieceAt(game.NewPosition(0, 1), p2)

		move := game.NewMove(game.NewPosition(0, 0), game.NewPosition(0, 1), &player1)
		res := SimulateMove(board, move)
		assert.Nil(t, res.GetPieceAt(game.NewPosition(0, 0)))
		assert.Equal(t, p2, res.GetPieceAt(game.NewPosition(0, 1)))
	})

	t.Run("same rank draws", func(t *testing.T) {
		board := game.NewBoard()
		p1 := game.NewPiece(models.Miner, &player1)
		p2 := game.NewPiece(models.Miner, &player2)
		board.SetPieceAt(game.NewPosition(0, 0), p1)
		board.SetPieceAt(game.NewPosition(0, 1), p2)

		move := game.NewMove(game.NewPosition(0, 0), game.NewPosition(0, 1), &player1)
		res := SimulateMove(board, move)
		assert.Nil(t, res.GetPieceAt(game.NewPosition(0, 0)))
		assert.Nil(t, res.GetPieceAt(game.NewPosition(0, 1)))
	})
}

func TestEvaluateBoard(t *testing.T) {
	player1 := game.NewPlayer(0, "us", "red")
	player2 := game.NewPlayer(1, "them", "blue")

	t.Run("basic evaluation", func(t *testing.T) {
		board := game.NewBoard()
		p1 := game.NewPiece(models.General, &player1)
		p2 := game.NewPiece(models.Miner, &player2)
		board.SetPieceAt(game.NewPosition(0, 0), p1)
		board.SetPieceAt(game.NewPosition(0, 1), p2)

		weights := map[string]float64{
			"General": 100.0,
			"Miner":   50.0,
		}

		score := EvaluateBoard(board, &player1, nil, weights, 1.0)
		assert.Equal(t, 125.0, score)
	})

	t.Run("exploration weight customized", func(t *testing.T) {
		board := game.NewBoard()
		p1 := game.NewPiece(models.General, &player1)
		board.SetPieceAt(game.NewPosition(0, 0), p1)

		weights := map[string]float64{
			"General":           100.0,
			"explorationWeight": 10.0,
		}
		score := EvaluateBoard(board, &player1, nil, weights, 0.5)
		assert.Equal(t, 145.0, score)
	})

	t.Run("player 1 exploration bonus", func(t *testing.T) {
		board := game.NewBoard()
		p1 := game.NewPiece(models.General, &player2)
		board.SetPieceAt(game.NewPosition(0, 8), p1)

		weights := map[string]float64{
			"General":           100.0,
			"explorationWeight": 10.0,
		}
		score := EvaluateBoard(board, &player2, nil, weights, 0.5)
		assert.Equal(t, 140.0, score)
	})

	t.Run("opponent piece revealed or remembered", func(t *testing.T) {
		board := game.NewBoard()
		p1 := game.NewPiece(models.General, &player1)
		p2 := game.NewPiece(models.Miner, &player2)
		board.SetPieceAt(game.NewPosition(0, 0), p1)
		pos2 := game.NewPosition(0, 1)
		board.SetPieceAt(pos2, p2)

		weights := map[string]float64{
			"General": 100.0,
			"Miner":   50.0,
		}

		p2.Reveal()
		scoreA := EvaluateBoard(board, &player1, nil, weights, 0.0)
		assert.Equal(t, 50.0, scoreA)

		p2 = game.NewPiece(models.Miner, &player2)
		board.SetPieceAt(pos2, p2)
		mem := NewMemory()
		mem.Remember(pos2, p2, 0.8, 1)
		scoreB := EvaluateBoard(board, &player1, mem, weights, 0.0)
		assert.Equal(t, 60.0, scoreB)
	})
}

func TestGetOpponent(t *testing.T) {
	board := game.NewBoard()
	player1 := game.NewPlayer(0, "us", "red")
	player2 := game.NewPlayer(1, "them", "blue")

	assert.Nil(t, GetOpponent(board, 0))

	p1 := game.NewPiece(models.General, &player1)
	board.SetPieceAt(game.NewPosition(0, 0), p1)
	assert.Nil(t, GetOpponent(board, 0))

	p2 := game.NewPiece(models.Miner, &player2)
	board.SetPieceAt(game.NewPosition(0, 1), p2)
	opp := GetOpponent(board, 0)
	assert.NotNil(t, opp)
	assert.Equal(t, player2.GetID(), opp.GetID())
}

func TestGetMovesAndIndices(t *testing.T) {
	player1 := game.NewPlayer(0, "us", "red")
	board := game.NewBoard()

	moves := GetMoves(board, &player1)
	assert.Empty(t, moves)

	p1 := game.NewPiece(models.General, &player1)
	pos := game.NewPosition(0, 0)
	board.SetPieceAt(pos, p1)

	moves = GetMoves(board, &player1)
	assert.Len(t, moves, 2)

	pBomb := game.NewPiece(models.Bomb, &player1)
	board.SetPieceAt(game.NewPosition(0, 2), pBomb)
	moves = GetMoves(board, &player1)
	assert.Len(t, moves, 2)

	idx := BuildMobileIndex(board, &player1)
	assert.Len(t, idx, 1)
	assert.Equal(t, pos, idx[0])

	movesFromIdx := GetMovesFromIndex(board, idx, &player1)
	assert.Len(t, movesFromIdx, 2)
}
