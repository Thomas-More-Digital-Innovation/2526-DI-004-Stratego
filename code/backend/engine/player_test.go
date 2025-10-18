package engine_test

import (
	"digital-innovation/stratego/engine"
	"digital-innovation/stratego/models"
	"testing"
)

func TestGetId(t *testing.T) {
	player := engine.NewPlayer(1, "player1", "avatar1")

	if player.GetID() == 0 {
		t.Errorf("Expected player ID to be non-zero, got %d", player.GetID())
	}
}

func TestGetName(t *testing.T) {
	player := engine.NewPlayer(1, "player1", "avatar1")

	if player.GetName() != "player1" {
		t.Errorf("Expected player name to be 'player1', got %s", player.GetName())
	}
}

func TestSetName(t *testing.T) {
	player := engine.NewPlayer(1, "player1", "avatar1")
	player.SetName("newName")

	if player.GetName() != "newName" {
		t.Errorf("Expected player name to be 'newName', got %s", player.GetName())
	}
}

func TestGetAvatar(t *testing.T) {
	player := engine.NewPlayer(1, "player1", "avatar1")

	if player.GetAvatar() != "avatar1" {
		t.Errorf("Expected player avatar to be 'avatar1', got %s", player.GetAvatar())
	}
}

func TestSetAvatar(t *testing.T) {
	player := engine.NewPlayer(1, "player1", "avatar1")
	player.SetAvatar("newAvatar")

	if player.GetAvatar() != "newAvatar" {
		t.Errorf("Expected player avatar to be 'newAvatar', got %s", player.GetAvatar())
	}
}

func TestInitializePieceScore(t *testing.T) {
	player := engine.NewPlayer(1, "player1", "avatar1")
	player.InitializePieceScore(10)

	if player.GetPieceScore() != 10 {
		t.Errorf("Expected player piece score to be 10, got %d", player.GetPieceScore())
	}
}
func TestResetPieceScore(t *testing.T) {
	player := engine.NewPlayer(1, "player1", "avatar1")
	player.InitializePieceScore(10)
	player.ResetPieceScore()

	if player.GetPieceScore() != 0 {
		t.Errorf("Expected player piece score to be 0 after reset, got %d", player.GetPieceScore())
	}
}

func TestUpdatePieceScore(t *testing.T) {
	initPieceValue := 100
	player := engine.NewPlayer(1, "player1", "avatar1")
	player.InitializePieceScore(initPieceValue)
	piece := engine.NewPiece(models.Sergeant, &player)

	player.UpdatePieceScore(piece)

	if player.GetPieceScore() != initPieceValue-models.Sergeant.GetStrategicValue() {
		t.Errorf("Expected player piece score to be 60 after update, got %d", player.GetPieceScore())
	}
}
