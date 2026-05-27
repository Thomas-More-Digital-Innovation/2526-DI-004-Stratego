package game

import (
	"digital-innovation/gostrategy/pkg/game/models"
	"testing"
)

func TestMakeMoveWithMutualAttack(t *testing.T) {
	player1 := NewPlayer(1, "Alice", "red")
	controller1 := NewHumanPlayerController(&player1)
	player2 := NewPlayer(2, "Bob", "blue")
	controller2 := NewHumanPlayerController(&player2)
	game := NewGame(controller1, controller2)

	attacker := NewPiece(models.Scout, &player1)
	defender := NewPiece(models.Scout, &player2)

	game.Board.SetPieceAt(NewPosition(0, 0), attacker)
	game.Board.SetPieceAt(NewPosition(0, 1), defender)
	move := NewMove(NewPosition(0, 0), NewPosition(0, 1), &player1)

	game.MakeMove(&move, attacker)

	if game.Board.GetPieceAt(move.GetTo()) != nil {
		t.Errorf("Expected both pieces to be removed after mutual annihilation")
	}

	if game.Board.GetPieceAt(move.GetFrom()) != nil {
		t.Errorf("Expected original position to be empty after move")
	}

	if attacker.IsAlive() {
		t.Errorf("Expected attacker to be dead after mutual annihilation")
	}

	if defender.IsAlive() {
		t.Errorf("Expected defender to be dead after mutual annihilation")
	}

	if game.CurrentPlayer != &player2 {
		t.Errorf("Expected current player to be player2 after move, but got %v", game.CurrentPlayer)
	}
}

func TestMakeMoveCapturingFlag(t *testing.T) {
	player1 := NewPlayer(1, "Alice", "red")
	controller1 := NewHumanPlayerController(&player1)
	player2 := NewPlayer(2, "Bob", "blue")
	controller2 := NewHumanPlayerController(&player2)
	game := NewGame(controller1, controller2)

	attacker := NewPiece(models.Major, &player1)
	flag := NewPiece(models.Flag, &player2)

	game.Board.SetPieceAt(NewPosition(0, 0), attacker)
	game.Board.SetPieceAt(NewPosition(0, 1), flag)
	move := NewMove(NewPosition(0, 0), NewPosition(0, 1), &player1)

	game.MakeMove(&move, attacker)

	if game.Board.GetPieceAt(move.GetTo()) != attacker {
		t.Errorf("Expected attacker to be at the new position after capturing flag")
	}

	if game.Board.GetPieceAt(move.GetFrom()) != nil {
		t.Errorf("Expected original position to be empty after move")
	}

	if !attacker.IsAlive() {
		t.Errorf("Expected attacker to be alive after capturing flag")
	}

	if flag.IsAlive() {
		t.Errorf("Expected flag to be dead after being captured")
	}

	if game.GetWinner() != &player1 {
		t.Errorf("Expected player1 to be the winner after capturing the flag")
	}
}

func TestGameClone(t *testing.T) {
	player1 := NewPlayer(1, "Alice", "red")
	controller1 := NewHumanPlayerController(&player1)
	player2 := NewPlayer(2, "Bob", "blue")
	controller2 := NewHumanPlayerController(&player2)
	g := NewGame(controller1, controller2)

	attacker := NewPiece(models.Major, &player1)
	g.Board.SetPieceAt(NewPosition(0, 0), attacker)

	cloned := g.Clone()
	if cloned.CurrentPlayer.GetID() != g.CurrentPlayer.GetID() {
		t.Error("Cloned game should have same current player ID")
	}
	if cloned.Board.GetPieceAt(NewPosition(0, 0)) == nil {
		t.Error("Cloned game board should have the piece")
	}
}

func TestGameInitializePieces(t *testing.T) {
	player1 := NewPlayer(1, "Alice", "red")
	controller1 := NewHumanPlayerController(&player1)
	player2 := NewPlayer(2, "Bob", "blue")
	controller2 := NewHumanPlayerController(&player2)
	g := NewGame(controller1, controller2)

	p1 := NewPiece(models.Flag, &player1)
	g.Board.SetPieceAt(NewPosition(0, 0), p1)

	g.InitializePieces()

	_, exists := player1.GetPiecePosition(p1)
	if !exists {
		t.Error("InitializePieces should update player piece tracking")
	}
}

func TestGameEnd(t *testing.T) {
	player1 := NewPlayer(1, "Alice", "red")
	controller1 := NewHumanPlayerController(&player1)
	player2 := NewPlayer(2, "Bob", "blue")
	controller2 := NewHumanPlayerController(&player2)
	g := NewGame(controller1, controller2)

	g.SetWinner(&player1, WinCauseFlagCaptured)
	if g.GetWinner() != &player1 {
		t.Error("Winner mismatch")
	}
	if !g.IsGameOver() {
		t.Error("Game should be over")
	}
	if g.GetWinCause() != WinCauseFlagCaptured {
		t.Error("Win cause mismatch")
	}
}
