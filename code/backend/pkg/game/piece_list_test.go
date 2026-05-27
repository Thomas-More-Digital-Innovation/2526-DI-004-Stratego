package game

import (
	"testing"
)

func TestPieceListLength(t *testing.T) {
	player := NewPlayer(1, "player1", "avatar1")
	pieceList := GetPieceList(&player)

	if len(pieceList) != 40 {
		t.Errorf("Expected 40 pieces, got %d", len(pieceList))
	}
}

func TestGetPieceListStrategicValue(t *testing.T) {
	// setup
	player := NewPlayer(1, "player1", "avatar1")
	pieceList := GetPieceList(&player)
	expectedValue := 219

	// test
	strategicValue := GetPieceListStrategicValue(pieceList)
	if strategicValue != expectedValue {
		t.Errorf("Expected strategic value to be %d, got %d", expectedValue, strategicValue)
	}
}
