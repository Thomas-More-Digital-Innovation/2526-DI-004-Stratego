package ai_test

import (
	"digital-innovation/gostrategy/internal/ai"
	"digital-innovation/gostrategy/internal/game"
	"digital-innovation/gostrategy/internal/game/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeterminizeBoard_NoOpponent(t *testing.T) {
	board := game.NewBoard()
	ourPlayer := game.NewPlayer(0, "us", "red")
	ourPiece := game.NewPiece(models.Marshal, &ourPlayer)
	board.SetPieceAt(game.NewPosition(0, 0), ourPiece)

	detBoard := ai.DeterminizeBoard(board, &ourPlayer, nil)
	assert.NotNil(t, detBoard)
	assert.Equal(t, ourPiece, detBoard.GetPieceAt(game.NewPosition(0, 0)))
}

func TestDeterminizeBoard_NoAlivePieces(t *testing.T) {
	board := game.NewBoard()
	ourPlayer := game.NewPlayer(0, "us", "red")
	opponent := game.NewPlayer(1, "them", "blue")

	ourPiece := game.NewPiece(models.Marshal, &ourPlayer)
	board.SetPieceAt(game.NewPosition(0, 0), ourPiece)

	oppPiece := game.NewPiece(models.General, &opponent)
	board.SetPieceAt(game.NewPosition(0, 1), oppPiece)

	detBoard := ai.DeterminizeBoard(board, &ourPlayer, nil)
	assert.NotNil(t, detBoard)
	assert.Equal(t, oppPiece, detBoard.GetPieceAt(game.NewPosition(0, 1)))
}

func TestDeterminizeBoard_WithAndWithoutMemory(t *testing.T) {
	board := game.NewBoard()
	ourPlayer := game.NewPlayer(0, "us", "red")
	opponent := game.NewPlayer(1, "them", "blue")

	ourPiece := game.NewPiece(models.Marshal, &ourPlayer)
	board.SetPieceAt(game.NewPosition(0, 0), ourPiece)
	ourPlayer.AddPiece(ourPiece, game.NewPosition(0, 0))

	oppPiece1 := game.NewPiece(models.General, &opponent)
	pos1 := game.NewPosition(0, 1)
	board.SetPieceAt(pos1, oppPiece1)
	opponent.AddPiece(oppPiece1, pos1)

	oppPiece2 := game.NewPiece(models.Spy, &opponent)
	pos2 := game.NewPosition(0, 2)
	board.SetPieceAt(pos2, oppPiece2)
	opponent.AddPiece(oppPiece2, pos2)

	oppPiece1.Reveal()

	t.Run("without memory", func(t *testing.T) {
		detBoard := ai.DeterminizeBoard(board, &ourPlayer, nil)
		assert.NotNil(t, detBoard)

		assert.Equal(t, models.General.GetName(), detBoard.GetPieceAt(pos1).GetType().GetName())
		assert.Equal(t, models.Spy.GetName(), detBoard.GetPieceAt(pos2).GetType().GetName())
	})

	t.Run("with memory recalling the piece", func(t *testing.T) {
		mem := ai.NewMemory()
		mem.Remember(pos2, oppPiece2, 1.0, 1)

		detBoard := ai.DeterminizeBoard(board, &ourPlayer, mem)
		assert.NotNil(t, detBoard)
		assert.Equal(t, models.General.GetName(), detBoard.GetPieceAt(pos1).GetType().GetName())
		assert.Equal(t, models.Spy.GetName(), detBoard.GetPieceAt(pos2).GetType().GetName())
	})

	t.Run("with memory but low confidence", func(t *testing.T) {
		mem := ai.NewMemory()
		mem.Remember(pos2, oppPiece2, 0.5, 1)

		detBoard := ai.DeterminizeBoard(board, &ourPlayer, mem)
		assert.NotNil(t, detBoard)
		assert.Equal(t, models.General.GetName(), detBoard.GetPieceAt(pos1).GetType().GetName())
		assert.Equal(t, models.Spy.GetName(), detBoard.GetPieceAt(pos2).GetType().GetName())
	})
}
