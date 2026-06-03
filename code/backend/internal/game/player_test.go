package game

import (
	"digital-innovation/gostrategy/internal/game/models"
	"testing"
)

func TestGetId(t *testing.T) {
	player := NewPlayer(1, "player1", "avatar1")

	if player.GetID() == 0 {
		t.Errorf("Expected player ID to be non-zero, got %d", player.GetID())
	}
}

func TestGetName(t *testing.T) {
	player := NewPlayer(1, "player1", "avatar1")

	if player.GetName() != "player1" {
		t.Errorf("Expected player name to be 'player1', got %s", player.GetName())
	}
}

func TestSetName(t *testing.T) {
	player := NewPlayer(1, "player1", "avatar1")
	player.SetName("newName")

	if player.GetName() != "newName" {
		t.Errorf("Expected player name to be 'newName', got %s", player.GetName())
	}
}

func TestGetAvatar(t *testing.T) {
	player := NewPlayer(1, "player1", "avatar1")

	if player.GetAvatar() != "avatar1" {
		t.Errorf("Expected player avatar to be 'avatar1', got %s", player.GetAvatar())
	}
}

func TestSetAvatar(t *testing.T) {
	player := NewPlayer(1, "player1", "avatar1")
	player.SetAvatar("newAvatar")

	if player.GetAvatar() != "newAvatar" {
		t.Errorf("Expected player avatar to be 'newAvatar', got %s", player.GetAvatar())
	}
}

func TestInitializePieceScore(t *testing.T) {
	player := NewPlayer(1, "player1", "avatar1")
	player.InitializePieceScore(10)

	if player.GetPieceScore() != 10 {
		t.Errorf("Expected player piece score to be 10, got %d", player.GetPieceScore())
	}
}
func TestResetPieceScore(t *testing.T) {
	player := NewPlayer(1, "player1", "avatar1")
	player.InitializePieceScore(10)
	player.ResetPieceScore()

	if player.GetPieceScore() != 0 {
		t.Errorf("Expected player piece score to be 0 after reset, got %d", player.GetPieceScore())
	}
}

func TestUpdatePieceScore(t *testing.T) {
	initPieceValue := 100
	player := NewPlayer(1, "player1", "avatar1")
	player.InitializePieceScore(initPieceValue)
	piece := NewPiece(models.Sergeant, &player)

	player.UpdatePieceScore(piece)

	if player.GetPieceScore() != initPieceValue-models.Sergeant.GetStrategicValue() {
		t.Errorf("Expected player piece score to be %d after update, got %d", initPieceValue-models.Sergeant.GetStrategicValue(), player.GetPieceScore())
	}
}

func TestPlayerWinning(t *testing.T) {
	player := NewPlayer(1, "player1", "avatar1")
	if player.HasWon() {
		t.Error("New player should not have won")
	}
	player.SetWinner()
	if !player.HasWon() {
		t.Error("Player should have won after SetWinner")
	}
}

func TestPlayerPieceTracking(t *testing.T) {
	player := NewPlayer(1, "player1", "avatar1")
	piece := NewPiece(models.Marshal, &player)
	pos := NewPosition(1, 2)

	player.AddPiece(piece, pos)

	if len(player.GetAlivePieces()) != 1 {
		t.Errorf("Expected 1 alive piece, got %d", len(player.GetAlivePieces()))
	}

	gotPos, exists := player.GetPiecePosition(piece)
	if !exists || !gotPos.Equals(pos) {
		t.Error("GetPiecePosition returned wrong position or doesn't exist")
	}

	newPos := NewPosition(3, 4)
	player.UpdatePiecePosition(piece, newPos)
	gotPos, _ = player.GetPiecePosition(piece)
	if !gotPos.Equals(newPos) {
		t.Error("UpdatePiecePosition failed")
	}

	player.RemovePiece(piece)
	if len(player.GetAlivePieces()) != 0 {
		t.Error("Alive pieces should be empty after removal")
	}
	_, exists = player.GetPiecePosition(piece)
	if exists {
		t.Error("Piece position should not exist after removal")
	}
}
