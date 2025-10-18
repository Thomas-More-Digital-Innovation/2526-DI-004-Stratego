package game_test

import (
	"digital-innovation/stratego/engine"
	"digital-innovation/stratego/game"
	"testing"
)

func TestPieceListLength(t *testing.T) {
	player := engine.NewPlayer(1, "player1", "avatar1")
	pieceList := game.GetPieceList(&player)

	if len(pieceList) != 40 {
		t.Errorf("Expected 40 pieces, got %d", len(pieceList))
	}
}

func TestGetPieceListStrategicValue(t *testing.T) {
	// setup
	player := engine.NewPlayer(1, "player1", "avatar1")
	pieceList := game.GetPieceList(&player)
	expectedValue := 219

	// test
	strategicValue := game.GetPieceListStrategicValue(pieceList)
	if strategicValue != expectedValue {
		t.Errorf("Expected strategic value to be %d, got %d", expectedValue, strategicValue)
	}
}
