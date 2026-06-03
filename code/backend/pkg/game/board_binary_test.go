package game

import (
	"digital-innovation/gostrategy/pkg/game/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testSetupRow     = "BBBBBB2222"
	validSetupPrefix = "0BBBBBB122"
	rowStr           = "row"
)

func TestEncodeBoard(t *testing.T) {
	board := NewBoard()
	player1 := NewPlayer(1, "Player1", "red")
	player2 := NewPlayer(2, "Player2", "blue")

	board.SetPieceAt(NewPosition(0, 0), NewPiece(models.Flag, &player1))
	board.SetPieceAt(NewPosition(1, 0), NewPiece(models.Marshal, &player2))

	moved := map[Position]bool{NewPosition(1, 0): true}
	data := EncodeBoard(board, moved)

	assert.Len(t, data, 100)

	// Check Flag encoding
	cell0 := data[0]
	assert.NotEqual(t, uint8(0), cell0&BitOccupied, "Position (0,0) should be occupied")
	assert.Equal(t, uint8(PieceIDFlag), (cell0&MaskPieceType)>>ShiftPieceType, "Position (0,0) should be Flag")
	assert.Equal(t, uint8(0), cell0&BitColor, "Position (0,0) should be Player 1")
	assert.Equal(t, uint8(0), cell0&BitMoved, "Position (0,0) should not have moved")

	// Check Marshal encoding
	cell1 := data[1]
	assert.NotEqual(t, uint8(0), cell1&BitOccupied, "Position (1,0) should be occupied")
	assert.Equal(t, uint8(PieceIDMarshal), (cell1&MaskPieceType)>>ShiftPieceType, "Position (1,0) should be Marshal")
	assert.NotEqual(t, uint8(0), cell1&BitColor, "Position (1,0) should be Player 2")
	assert.NotEqual(t, uint8(0), cell1&BitMoved, "Position (1,0) should have moved")
}

func TestDecodeBoard(t *testing.T) {
	player1 := NewPlayer(1, "Player1", "red")
	player2 := NewPlayer(2, "Player2", "blue")

	data := make([]byte, 100)
	data[0] = BitOccupied | (PieceIDFlag << ShiftPieceType)
	data[1] = BitOccupied | (PieceIDMarshal << ShiftPieceType) | BitColor | BitMoved

	board, moved := DecodeBoard(data, &player1, &player2)

	piece0 := board.GetPieceAt(NewPosition(0, 0))
	require.NotNil(t, piece0, "Position (0,0) should have a piece")
	assert.Equal(t, uint8('0'), piece0.GetType().GetRank())
	assert.Equal(t, 1, piece0.GetOwner().GetID())
	assert.False(t, moved[NewPosition(0, 0)])

	piece1 := board.GetPieceAt(NewPosition(1, 0))
	require.NotNil(t, piece1, "Position (1,0) should have a piece")
	assert.Equal(t, uint8('M'), piece1.GetType().GetRank())
	assert.Equal(t, 2, piece1.GetOwner().GetID())
	assert.True(t, moved[NewPosition(1, 0)])
}

func TestBoardBase64RoundTrip(t *testing.T) {
	board := NewBoard()
	player1 := NewPlayer(1, "Player1", "red")
	player2 := NewPlayer(2, "Player2", "blue")

	board.SetPieceAt(NewPosition(0, 0), NewPiece(models.Flag, &player1))
	board.SetPieceAt(NewPosition(9, 9), NewPiece(models.Marshal, &player2))

	encoded := EncodeBoardToBase64(board, nil)
	decoded, _, err := DecodeBoardFromBase64(encoded, &player1, &player2)
	require.NoError(t, err)

	assert.Equal(t, uint8('0'), decoded.GetPieceAt(NewPosition(0, 0)).GetType().GetRank())
	assert.Equal(t, uint8('M'), decoded.GetPieceAt(NewPosition(9, 9)).GetType().GetRank())
}

func TestEncodeSetup(t *testing.T) {
	rows := []string{testSetupRow, "4444555566", "6677788399", "M311110333"}

	data, err := EncodeSetup(rows, 1)
	require.NoError(t, err)
	assert.Len(t, data, 40)

	// Check first cell (Bomb)
	cell0 := data[0]
	assert.NotEqual(t, uint8(0), cell0&BitOccupied)
	assert.Equal(t, uint8(PieceIDBomb), (cell0&MaskPieceType)>>ShiftPieceType)
	assert.Equal(t, uint8(0), cell0&BitColor)
}

func TestDecodeSetup(t *testing.T) {
	originalRows := []string{testSetupRow, "4444555566", "6677788399", "M311110333"}

	data, _ := EncodeSetup(originalRows, 1)
	rows, playerID, err := DecodeSetup(data)
	require.NoError(t, err)
	assert.Equal(t, 1, playerID)
	assert.Equal(t, originalRows, rows)
}

func TestValidateSetup(t *testing.T) {
	// Valid: 1 Flag, 6 Bombs, 1 Spy, 8 Scouts, 5 Miners, 4 Sergeants, 4 Lieutenants, 4 Captains, 3 Majors, 2 Colonels, 1 General, 1 Marshal
	valid := []string{
		validSetupPrefix,
		"2222223333",
		"3444455556",
		"666777889M",
	}
	assert.NoError(t, ValidateSetup(valid))

	// Invalid: 2 Flags instead of 1
	invalid := []string{testSetupRow, "2222224444", "5555666677", "00899M3333"}
	assert.Error(t, ValidateSetup(invalid))

	// Wrong row count
	assert.Error(t, ValidateSetup([]string{rowStr}))

	// Wrong row length
	assert.Error(t, ValidateSetup([]string{rowStr, rowStr, rowStr, rowStr}))
}

func TestBinaryHelpers(t *testing.T) {
	id, ok := GetPieceIDFromRank('M')
	assert.True(t, ok)
	assert.Equal(t, uint8(PieceIDMarshal), id)

	_, ok = GetPieceIDFromRank('X')
	assert.False(t, ok)

	pt := GetPieceTypeFromID(PieceIDMarshal)
	assert.Equal(t, "Marshal", pt.GetName())

	cell := BitOccupied | (PieceIDMarshal << ShiftPieceType)
	pt2 := GetPieceTypeFromCell(cell)
	assert.Equal(t, "Marshal", pt2.GetName())
}

func TestBoardBinaryErrors(t *testing.T) {
	player1 := NewPlayer(1, "P1", "")
	player2 := NewPlayer(2, "P2", "")

	_, _, err := DecodeBoardFromBase64("!!!", &player1, &player2)
	assert.Error(t, err)

	_, _, err = DecodeBoardFromBase64("YWJjZA==", &player1, &player2) // "abcd" in base64
	assert.Error(t, err)

	_, err = EncodeSetup([]string{rowStr}, 1)
	assert.Error(t, err)

	_, err = EncodeSetup([]string{"123", "123", "123", "123"}, 1)
	assert.Error(t, err)

	_, _, err = DecodeSetup(make([]byte, 10))
	assert.Error(t, err)
}

func TestEncodeSetupPlayer2(t *testing.T) {
	rows := []string{"0BBBBBB122", "2222223333", "3444455556", "666777889M"}
	data, err := EncodeSetup(rows, 2)
	require.NoError(t, err)
	assert.NotEqual(t, uint8(0), data[0]&BitColor)

	decodedRows, playerID, err := DecodeSetup(data)
	require.NoError(t, err)
	assert.Equal(t, 2, playerID)
	assert.Equal(t, rows, decodedRows)
}
