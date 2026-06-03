package game

import (
	"digital-innovation/gostrategy/internal/game/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBoard(t *testing.T) {
	board := NewBoard()
	require.NotNil(t, board, "Expected a board to be created, but got nil")

	assert.Len(t, board.GetField(), 10, "Expected board field to have 10 rows")

	for _, row := range board.GetField() {
		assert.Len(t, row, 10, "Expected each row to have 10 columns")
	}

	expectedLakes := []Position{
		NewPosition(2, 4), NewPosition(3, 4), NewPosition(2, 5), NewPosition(3, 5),
		NewPosition(6, 4), NewPosition(7, 4), NewPosition(6, 5), NewPosition(7, 5),
	}

	for _, lakePos := range expectedLakes {
		assert.True(t, board.IsLake(lakePos), "Expected position %v to be a lake", lakePos)
	}
}

func TestSetAndGetPieceAt(t *testing.T) {
	board := NewBoard()
	player := NewPlayer(1, "Alice", "red")
	piece := NewPiece(models.Marshal, &player)
	position := NewPosition(0, 0)

	board.SetPieceAt(position, piece)
	retrievedPiece := board.GetPieceAt(position)

	assert.Equal(t, piece, retrievedPiece, "Expected to retrieve the same piece that was set")
}

func TestGetPieceAt(t *testing.T) {
	board := NewBoard()
	player := NewPlayer(1, "Alice", "red")
	piece := NewPiece(models.Marshal, &player)
	position := NewPosition(1, 1)

	board.SetPieceAt(position, piece)
	retrievedPiece := board.GetPieceAt(position)

	assert.Equal(t, piece, retrievedPiece, "Expected to retrieve the same piece that was set")
}

func TestIsLake(t *testing.T) {
	board := NewBoard()
	position := NewPosition(2, 4)

	isLake := board.IsLake(position)
	assert.True(t, isLake, "Expected position %v to be a lake", position)
}

func TestIsValidMove(t *testing.T) {
	board := NewBoard()
	player := NewPlayer(1, "Alice", "red")
	piece := NewPiece(models.Marshal, &player)
	fromPos := NewPosition(0, 0)
	toPos := NewPosition(0, 1)

	board.SetPieceAt(fromPos, piece)
	move := NewMove(fromPos, toPos, &player)

	assert.True(t, board.IsValidMove(&move), "Expected move from %v to %v to be valid", fromPos, toPos)

	lakePos := NewPosition(2, 4)
	moveToLake := NewMove(fromPos, lakePos, &player)

	assert.False(t, board.IsValidMove(&moveToLake), "Expected move to lake position %v to be invalid", lakePos)
}

func TestIsInvalidMoveOutsideField(t *testing.T) {
	board := NewBoard()
	player := NewPlayer(1, "Alice", "red")
	piece := NewPiece(models.Marshal, &player)
	fromPos := NewPosition(0, 0)
	toPos := NewPosition(10, 10) // Outside the board

	board.SetPieceAt(fromPos, piece)
	move := NewMove(fromPos, toPos, &player)

	assert.False(t, board.IsValidMove(&move), "Expected move from %v to %v to be invalid (outside field)", fromPos, toPos)
}

func TestIsInvalidMoveIntoLake(t *testing.T) {
	board := NewBoard()
	player := NewPlayer(1, "Alice", "red")
	piece := NewPiece(models.Marshal, &player)
	fromPos := NewPosition(1, 4)
	toPos := NewPosition(2, 4) // Lake position

	board.SetPieceAt(fromPos, piece)
	move := NewMove(fromPos, toPos, &player)

	assert.False(t, board.IsValidMove(&move), "Expected move from %v to lake position %v to be invalid", fromPos, toPos)
}

func TestIsInvalidMoveToTeamPiece(t *testing.T) {
	board := NewBoard()
	player := NewPlayer(1, "Alice", "red")
	piece1 := NewPiece(models.Marshal, &player)
	piece2 := NewPiece(models.General, &player)
	fromPos := NewPosition(0, 0)
	toPos := NewPosition(0, 1)

	board.SetPieceAt(fromPos, piece1)
	board.SetPieceAt(toPos, piece2)
	move := NewMove(fromPos, toPos, &player)

	assert.False(t, board.IsValidMove(&move), "Expected move from %v to %v to be invalid (to own piece)", fromPos, toPos)
}

func TestRandomizeSetup(_ *testing.T) {
	// unimplemented
}

func TestSwapPieces(t *testing.T) {
	board := NewBoard()
	player := NewPlayer(1, "Alice", "red")
	piece1 := NewPiece(models.Marshal, &player)
	piece2 := NewPiece(models.General, &player)
	pos1 := NewPosition(0, 0)
	pos2 := NewPosition(0, 1)

	board.SetPieceAt(pos1, piece1)
	board.SetPieceAt(pos2, piece2)

	err := board.SwapPieces(pos1, pos2)
	assert.NoError(t, err)

	assert.Equal(t, piece2, board.GetPieceAt(pos1))
	assert.Equal(t, piece1, board.GetPieceAt(pos2))
}

func TestSwapPiecesInvalidPosition(t *testing.T) {
	board := NewBoard()
	player := NewPlayer(1, "Alice", "red")
	piece1 := NewPiece(models.Marshal, &player)
	pos1 := NewPosition(0, 0)
	pos2 := NewPosition(0, 1)

	board.SetPieceAt(pos1, piece1)

	err := board.SwapPieces(pos1, pos2)
	assert.Error(t, err)
}

func TestRemovePieceAt(t *testing.T) {
	board := NewBoard()
	player := NewPlayer(1, "Alice", "red")
	piece := NewPiece(models.Marshal, &player)
	position := NewPosition(0, 0)

	board.SetPieceAt(position, piece)
	piece.Eliminate()
	err := board.RemovePieceAt(position)
	assert.NoError(t, err)

	assert.Nil(t, board.GetPieceAt(position))
}

func TestRemovePieceThatDoesntExist(t *testing.T) {
	board := NewBoard()
	position := NewPosition(0, 0)

	err := board.RemovePieceAt(position)
	assert.Error(t, err)
}

func TestRemoveAlivePiece(t *testing.T) {
	board := NewBoard()
	player := NewPlayer(1, "Alice", "red")
	piece := NewPiece(models.Marshal, &player)
	position := NewPosition(0, 0)

	board.SetPieceAt(position, piece)
	// piece is still alive
	err := board.RemovePieceAt(position)
	assert.Error(t, err)
}
