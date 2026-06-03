package game

import (
	"digital-innovation/gostrategy/internal/game/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMovePiece(t *testing.T) {
	board := NewBoard()
	player := NewPlayer(1, "Alice", "red")
	piece := NewPiece(models.Marshal, &player)
	fromPos := NewPosition(0, 0)
	toPos := NewPosition(0, 1)

	board.SetPieceAt(fromPos, piece)
	move := NewMove(fromPos, toPos, &player)

	board.MovePiece(&move, piece)

	assert.Equal(t, piece, board.GetPieceAt(toPos))
	assert.Nil(t, board.GetPieceAt(fromPos))
}

func TestListMovesStandardPiece(t *testing.T) {
	board := NewBoard()
	player := NewPlayer(1, "Alice", "red")
	piece := NewPiece(models.Marshal, &player)
	position := NewPosition(1, 1)

	board.SetPieceAt(position, piece)

	moves, err := board.ListMoves(position)
	require.NoError(t, err)

	assert.Len(t, moves, 4) // Up, Down, Left, Right
}

func TestListMovesStandardPieceNoMovesAvailable(t *testing.T) {
	board := NewBoard()
	player := NewPlayer(1, "Alice", "red")
	piece := NewPiece(models.Marshal, &player)
	position := NewPosition(1, 1)
	blockedPositions := []Position{
		NewPosition(1, 0), // Up
		NewPosition(1, 2), // Down
		NewPosition(0, 1), // Left
		NewPosition(2, 1), // Right
	}

	for _, pos := range blockedPositions {
		blockingPiece := NewPiece(models.Scout, &player)
		board.SetPieceAt(pos, blockingPiece)
	}

	board.SetPieceAt(position, piece)

	moves, err := board.ListMoves(position)
	require.NoError(t, err)
	assert.Empty(t, moves)
}

func TestListMovesScout(t *testing.T) {
	board := NewBoard()
	player := NewPlayer(1, "Alice", "red")
	piece := NewPiece(models.Scout, &player)
	position := NewPosition(0, 0)

	board.SetPieceAt(position, piece)

	moves, err := board.ListMoves(position)
	require.NoError(t, err)
	assert.Len(t, moves, 18) // 9 Down, 9 Right
}

func TestListMovesNoPiece(t *testing.T) {
	board := NewBoard()
	position := NewPosition(0, 0)

	moves, err := board.ListMoves(position)
	assert.Error(t, err)
	assert.Empty(t, moves)
}

func TestClone(t *testing.T) {
	board := NewBoard()
	player := NewPlayer(1, "Alice", "red")
	piece := NewPiece(models.Marshal, &player)
	pos := NewPosition(0, 0)
	board.SetPieceAt(pos, piece)

	cloned := board.Clone()
	clonedPiece := cloned.GetPieceAt(pos)
	require.NotNil(t, clonedPiece)
	assert.NotSame(t, piece, clonedPiece)
	assert.Equal(t, piece.GetType().GetRank(), clonedPiece.GetType().GetRank())
}

func TestBoardString(t *testing.T) {
	board := NewBoard()
	player := NewPlayer(1, "Alice", "red")
	piece := NewPiece(models.Marshal, &player)
	board.SetPieceAt(NewPosition(0, 0), piece)
	s := board.String()
	assert.NotEmpty(t, s)

	assert.Contains(t, s, piece.GetType().GetIcon())
	assert.Contains(t, s, "~~")
	assert.Contains(t, s, "..")
}

func TestIsValidMoveBoundaries(t *testing.T) {
	board := NewBoard()
	player := NewPlayer(1, "P1", "")
	piece := NewPiece(models.Marshal, &player)
	from := NewPosition(0, 0)
	board.SetPieceAt(from, piece)

	tests := []struct {
		x, y int
		want bool
	}{
		{-1, 0, false},
		{0, -1, false},
		{10, 0, false},
		{0, 10, false},
		{0, 0, false}, // same as from (blocked by team piece)
	}

	for _, tt := range tests {
		move := NewMove(from, NewPosition(tt.x, tt.y), &player)
		assert.Equal(t, tt.want, board.IsValidMove(&move))
	}

	// No piece at from
	moveNoPiece := NewMove(NewPosition(5, 5), NewPosition(5, 6), &player)
	assert.False(t, board.IsValidMove(&moveNoPiece))
}

func TestScoutMovesOutOfBounds(t *testing.T) {
	board := NewBoard()
	player := NewPlayer(1, "P1", "")
	scout := NewPiece(models.Scout, &player)
	from := NewPosition(0, 0)
	board.SetPieceAt(from, scout)

	// This will trigger handleScoutMoves and it will hit boundaries (IsValidMove will return false)
	moves, _ := board.ListMoves(from)
	assert.NotEmpty(t, moves)
}

func TestEncodeSetupEmptyCells(t *testing.T) {
	// Test EncodeSetup with '.' and ' '
	rows := []string{
		validSetupPrefix,
		"2222223333",
		"3444455556",
		"66677788. ",
	}
	data, err := EncodeSetup(rows, 1)
	require.NoError(t, err)

	// Cell 38 and 39 should not be occupied
	assert.Equal(t, uint8(0), data[38]&BitOccupied)
	assert.Equal(t, uint8(0), data[39]&BitOccupied)

	decoded, _, _ := DecodeSetup(data)
	assert.Equal(t, byte('.'), decoded[3][8])
	assert.Equal(t, byte('.'), decoded[3][9])
}

func TestIsValidMoveScoutLongDistance(t *testing.T) {
	board := NewBoard()
	player := NewPlayer(1, "Alice", "red")
	scout := NewPiece(models.Scout, &player)
	from := NewPosition(0, 0)
	to := NewPosition(0, 5)
	board.SetPieceAt(from, scout)

	move := NewMove(from, to, &player)
	assert.True(t, board.IsValidMove(&move))
}

func TestListMovesScoutBlocked(t *testing.T) {
	board := NewBoard()
	player := NewPlayer(1, "Alice", "red")
	scout := NewPiece(models.Scout, &player)
	from := NewPosition(0, 0)
	board.SetPieceAt(from, scout)

	// Blocked by lake
	fromLakeSide := NewPosition(2, 3)
	board.SetPieceAt(fromLakeSide, scout)

	moves, _ := board.ListMoves(fromLakeSide)
	for _, m := range moves {
		if m.GetTo().Y > 3 && m.GetTo().X == 2 {
			t.Errorf("Scout should be blocked by lake, but got move to %v", m.GetTo())
		}
	}

	// Blocked by piece
	board = NewBoard()
	board.SetPieceAt(from, scout)
	blocker := NewPiece(models.Miner, &player)
	board.SetPieceAt(NewPosition(0, 3), blocker)

	moves, _ = board.ListMoves(from)
	for _, m := range moves {
		if m.GetTo().Y > 3 && m.GetTo().X == 0 {
			t.Errorf("Scout should be blocked by piece, but got move to %v", m.GetTo())
		}
	}
}
