package game

import (
	"digital-innovation/gostrategy/pkg/game/models"
	"testing"
)

func TestAttackBomb(t *testing.T) {
	// setup
	player1 := NewPlayer(1, "player1", "avatar1")
	colonel := NewPiece(models.Colonel, &player1)
	player1.InitializePieceScore(colonel.GetStrategicValue())

	player2 := NewPlayer(2, "player2", "avatar2")
	bomb := NewPiece(models.Bomb, &player2)
	player2.InitializePieceScore(bomb.GetStrategicValue())

	// execute
	result := colonel.Attack(bomb)
	attacker, target := result[0], result[1]

	// verify
	if !target.IsAlive() {
		t.Errorf("Expected bomb to be alive, got eliminated")
	}
	if attacker.IsAlive() {
		t.Errorf("Expected colonel to be eliminated, got alive")
	}
	expectedScore1 := 0
	if player1.GetPieceScore() != expectedScore1 {
		t.Errorf("Expected player1 piece score to be %d, got %d", expectedScore1, player1.GetPieceScore())
	}

	expectedScore2 := models.Bomb.GetStrategicValue()
	if player2.GetPieceScore() != expectedScore2 {
		t.Errorf("Expected player2 piece score to be %d, got %d", expectedScore2, player2.GetPieceScore())
	}
}

func TestStandardAttack(t *testing.T) {
	// setup
	player1 := NewPlayer(1, "player1", "avatar1")
	major := NewPiece(models.Major, &player1)
	player1.InitializePieceScore(major.GetStrategicValue())

	player2 := NewPlayer(2, "player2", "avatar2")
	captain := NewPiece(models.Captain, &player2)
	player2.InitializePieceScore(captain.GetStrategicValue())

	// execute
	result := major.Attack(captain)
	attacker, target := result[0], result[1]

	// verify
	if target.IsAlive() {
		t.Errorf("Expected captain to be eliminated, got alive")
	}
	if !attacker.IsAlive() {
		t.Errorf("Expected major to be alive, got eliminated")
	}

	expectedScore1 := models.Major.GetStrategicValue()
	if player1.GetPieceScore() != expectedScore1 {
		t.Errorf("Expected player1 piece (major) score to be %d, got %d", expectedScore1, player1.GetPieceScore())
	}

	expectedScore2 := 0
	if player2.GetPieceScore() != expectedScore2 {
		t.Errorf("Expected player2 piece score to be %d, got %d", expectedScore2, player2.GetPieceScore())
	}
}

func TestStandardAttackEqualRank(t *testing.T) {
	// setup
	player1 := NewPlayer(1, "player1", "avatar1")
	sergeant1 := NewPiece(models.Sergeant, &player1)
	player1.InitializePieceScore(sergeant1.GetStrategicValue())

	player2 := NewPlayer(2, "player2", "avatar2")
	sergeant2 := NewPiece(models.Sergeant, &player2)
	player2.InitializePieceScore(sergeant2.GetStrategicValue())

	// execute
	result := sergeant1.Attack(sergeant2)
	attacker, target := result[0], result[1]

	// verify
	if attacker.IsAlive() {
		t.Errorf("Expected attacker sergeant to be eliminated, got alive")
	}
	if target.IsAlive() {
		t.Errorf("Expected target sergeant to be eliminated, got alive")
	}

	expectedScore1 := 0
	if player1.GetPieceScore() != expectedScore1 {
		t.Errorf("Expected player1 piece score to be %d, got %d", expectedScore1, player1.GetPieceScore())
	}

	expectedScore2 := 0
	if player2.GetPieceScore() != expectedScore2 {
		t.Errorf("Expected player2 piece score to be %d, got %d", expectedScore2, player2.GetPieceScore())
	}
}

func TestPieceHide(t *testing.T) {
	player := NewPlayer(1, "P1", "")
	piece := NewPiece(models.Marshal, &player)
	piece.Reveal()
	if !piece.IsRevealed() {
		t.Error("Piece should be revealed")
	}
	piece.Hide()
	if piece.IsRevealed() {
		t.Error("Piece should be hidden")
	}
}

func TestPieceClone(t *testing.T) {
	player := NewPlayer(1, "P1", "")
	piece := NewPiece(models.Marshal, &player)
	cloned := piece.Clone()
	if cloned.GetType().GetName() != piece.GetType().GetName() {
		t.Error("Clone type mismatch")
	}
	if cloned.GetOwner() != piece.GetOwner() {
		t.Error("Clone owner mismatch")
	}
}

func TestStandardAttackLoss(t *testing.T) {
	player1 := NewPlayer(1, "P1", "")
	captain := NewPiece(models.Captain, &player1)
	player2 := NewPlayer(2, "P2", "")
	major := NewPiece(models.Major, &player2)

	result := captain.Attack(major)
	attacker, target := result[0], result[1]

	if attacker.IsAlive() {
		t.Error("Attacker should be eliminated")
	}
	if !target.IsAlive() {
		t.Error("Target should be alive")
	}
}
