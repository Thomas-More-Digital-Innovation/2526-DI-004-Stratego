package game

import (
	"digital-innovation/gostrategy/internal/game/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupGame(t *testing.T) {
	player1 := NewPlayer(0, "Player 1", "avatar1")
	player2 := NewPlayer(1, "Player 2", "avatar2")

	controller1 := NewHumanPlayerController(&player1)
	controller2 := NewHumanPlayerController(&player2)

	g := NewGame(controller1, controller2)

	player1Pieces := GetPieceList(&player1)
	player2Pieces := GetPieceList(&player2)

	err := SetupGame(g, player1Pieces, player2Pieces)
	require.NoError(t, err)

	assert.Len(t, player1.GetAlivePieces(), 40)
	assert.Len(t, player2.GetAlivePieces(), 40)

	field := g.Board.GetField()

	// Player 1 should be in rows 6-9
	for y := 6; y <= 9; y++ {
		for x := 0; x < 10; x++ {
			piece := field[y][x]
			require.NotNil(t, piece, "Expected piece at (%d, %d) for player 1", x, y)
			assert.Equal(t, 0, piece.GetOwner().GetID())
		}
	}

	// Player 2 should be in rows 0-3
	for y := 0; y <= 3; y++ {
		for x := 0; x < 10; x++ {
			piece := field[y][x]
			require.NotNil(t, piece, "Expected piece at (%d, %d) for player 2", x, y)
			assert.Equal(t, 1, piece.GetOwner().GetID())
		}
	}

	// Middle rows (4-5) should be empty (except lakes)
	for y := 4; y <= 5; y++ {
		for x := 0; x < 10; x++ {
			assert.Nil(t, field[y][x], "Expected no piece at (%d, %d)", x, y)
		}
	}

	assert.NotZero(t, player1.GetPieceScore())
	assert.NotZero(t, player2.GetPieceScore())
}

func TestSetupGame_InvalidPieceCount(t *testing.T) {
	player1 := NewPlayer(0, "Player 1", "avatar1")
	player2 := NewPlayer(1, "Player 2", "avatar2")

	controller1 := NewHumanPlayerController(&player1)
	controller2 := NewHumanPlayerController(&player2)

	g := NewGame(controller1, controller2)

	player1Pieces := GetPieceList(&player1)
	player2Pieces := GetPieceList(&player2)[:30] // Only 30 pieces

	err := SetupGame(g, player1Pieces, player2Pieces)
	assert.Error(t, err)
	assert.Equal(t, "each player must have exactly 40 pieces", err.Error())
}

func TestRandomSetup(t *testing.T) {
	player := NewPlayer(0, "Test Player", "avatar")

	pieces1 := RandomSetup(&player)
	pieces2 := RandomSetup(&player)

	assert.Len(t, pieces1, 40)

	different := false
	for i := 0; i < 40; i++ {
		if pieces1[i].GetType().GetName() != pieces2[i].GetType().GetName() {
			different = true
			break
		}
	}
	assert.True(t, different, "Two random setups should likely be different")
}

func TestQuickStart(t *testing.T) {
	player1 := NewPlayer(0, "Player 1", "avatar1")
	player2 := NewPlayer(1, "Player 2", "avatar2")

	controller1 := NewHumanPlayerController(&player1)
	controller2 := NewHumanPlayerController(&player2)

	g := QuickStart(controller1, controller2)
	require.NotNil(t, g)

	assert.Len(t, player1.GetAlivePieces(), 40)
	assert.Len(t, player2.GetAlivePieces(), 40)
	assert.NotNil(t, g.CurrentPlayer)
	assert.Equal(t, 1, g.GetRound())
	assert.Nil(t, g.GetWinner())
}

func TestPiecePositionTracking(t *testing.T) {
	player1 := NewPlayer(0, "Player 1", "avatar1")
	player2 := NewPlayer(1, "Player 2", "avatar2")

	controller1 := NewHumanPlayerController(&player1)
	controller2 := NewHumanPlayerController(&player2)

	QuickStart(controller1, controller2)

	for _, piece := range player1.GetAlivePieces() {
		pos, exists := player1.GetPiecePosition(piece)
		assert.True(t, exists)
		assert.True(t, pos.Y >= 6 && pos.Y <= 9)
	}

	for _, piece := range player2.GetAlivePieces() {
		pos, exists := player2.GetPiecePosition(piece)
		assert.True(t, exists)
		assert.True(t, pos.Y >= 0 && pos.Y <= 3)
	}
}

func TestPieceListConsistency(t *testing.T) {
	player := NewPlayer(0, "Test Player", "avatar")

	pieces := GetPieceList(&player)

	counts := make(map[string]int)
	for _, piece := range pieces {
		counts[piece.GetType().GetName()]++
	}

	expected := map[string]int{
		"Flag":       1,
		"Bomb":       6,
		"Spy":        1,
		"Scout":      8,
		"Miner":      5,
		"Sergeant":   4,
		"Lieutenant": 4,
		"Captain":    4,
		"Major":      3,
		"Colonel":    2,
		"General":    1,
		"Marshal":    1,
	}

	for name, expectedCount := range expected {
		assert.Equal(t, expectedCount, counts[name])
	}
}

func TestParseSetup(t *testing.T) {
	player := NewPlayer(1, "P1", "")

	// Rank string format
	rankStr := "0BBBBBB12222222233333444455556666777889M"
	pieces, err := ParseSetup(&player, []byte(rankStr))
	require.NoError(t, err)
	assert.Len(t, pieces, 40)

	// Bitpacked format
	bitpacked := make([]byte, 40)
	for i := range bitpacked {
		bitpacked[i] = BitOccupied | (PieceIDScout << ShiftPieceType)
	}
	// This will fail validation because counts are wrong (40 scouts)
	_, err = ParseSetup(&player, bitpacked)
	assert.Error(t, err)

	// Error: empty cell
	rankStrEmpty := "0BBBBBB12222222233333444455556666777889."
	_, err = ParseSetup(&player, []byte(rankStrEmpty))
	assert.Error(t, err)

	// Error: invalid rank
	rankStrInvalid := "0BBBBBB12222222233333444455556666777889X"
	_, err = ParseSetup(&player, []byte(rankStrInvalid))
	assert.Error(t, err)
}

func TestCustomSetup(t *testing.T) {
	player := NewPlayer(1, "P1", "")
	pieceMap := make(map[Position]*Piece)
	for i := 0; i < 40; i++ {
		pieceMap[NewPosition(i%10, i/10)] = NewPiece(models.Scout, &player)
	}

	pieces, err := CustomSetup(&player, pieceMap)
	require.NoError(t, err)
	assert.Len(t, pieces, 40)

	// Invalid count
	delete(pieceMap, NewPosition(0, 0))
	_, err = CustomSetup(&player, pieceMap)
	assert.Error(t, err)
}
