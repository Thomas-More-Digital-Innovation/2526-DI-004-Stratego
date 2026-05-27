package game

import (
	"digital-innovation/gostrategy/internal/models"
	"testing"
)

func TestNewGame(t *testing.T) {
	player1 := NewPlayer(1, "Alice", "red")
	controller1 := NewHumanPlayerController(&player1)
	player2 := NewPlayer(2, "Bob", "blue")
	controller2 := NewHumanPlayerController(&player2)
	game := NewGame(controller1, controller2)

	if game == nil {
		t.Fatal("Expected a game to be created, but got nil")
	}

	if len(game.Players) != 2 {
		t.Errorf("Expected 2 players, but got %d", len(game.Players))
	}

	if game.CurrentPlayer != &player1 {
		t.Errorf("Expected current player to be player1, but got %v", game.CurrentPlayer)
	}

	if game.Board == nil {
		t.Errorf("Expected a board to be created, but got nil")
	}
}

func TestNextTurn(t *testing.T) {
	player1 := NewPlayer(1, "Alice", "red")
	controller1 := NewHumanPlayerController(&player1)
	player2 := NewPlayer(2, "Bob", "blue")
	controller2 := NewHumanPlayerController(&player2)
	game := NewGame(controller1, controller2)

	if game.CurrentPlayer != &player1 {
		t.Errorf("Expected current player to be player1, but got %v", game.CurrentPlayer)
	}

	game.NextTurn()
	if game.CurrentPlayer != &player2 {
		t.Errorf("Expected current player to be player2 after next turn, but got %v", game.CurrentPlayer)
	}

	game.NextTurn()
	if game.CurrentPlayer != &player1 {
		t.Errorf("Expected current player to be player1 after next turn, but got %v", game.CurrentPlayer)
	}
}

func TestGetRound(t *testing.T) {
	player1 := NewPlayer(1, "Alice", "red")
	controller1 := NewHumanPlayerController(&player1)
	piece1 := NewPiece(models.Scout, &player1)

	player2 := NewPlayer(2, "Bob", "blue")
	controller2 := NewHumanPlayerController(&player2)
	piece2 := NewPiece(models.Major, &player2)

	game := NewGame(controller1, controller2)

	game.Board.SetPieceAt(NewPosition(0, 0), piece1)
	game.Board.SetPieceAt(NewPosition(1, 0), piece2)

	if game.GetRound() != 1 {
		t.Errorf("Expected round to be 1 at game start, but got %d", game.GetRound())
	}

	move1 := NewMove(NewPosition(0, 0), NewPosition(0, 4), &player1)
	game.MakeMove(&move1, piece1)

	if game.GetRound() != 1 {
		t.Errorf("Expected round to be 1 after one move, but got %d", game.GetRound())
	}

	move2 := NewMove(NewPosition(1, 0), NewPosition(1, 1), &player2)
	game.MakeMove(&move2, piece2)

	if game.GetRound() != 2 {
		t.Errorf("Expected round to be 2 after two moves, but got %d", game.GetRound())
	}
}

func TestMakeMoveToEmptyCell(t *testing.T) {
	player1 := NewPlayer(1, "Alice", "red")
	controller1 := NewHumanPlayerController(&player1)
	player2 := NewPlayer(2, "Bob", "blue")
	controller2 := NewHumanPlayerController(&player2)
	game := NewGame(controller1, controller2)

	piece := NewPiece(models.Major, &player1)
	move := NewMove(NewPosition(0, 0), NewPosition(0, 1), &player1)

	game.MakeMove(&move, piece)

	if game.Board.GetPieceAt(move.GetTo()) != piece {
		t.Errorf("Expected piece to be at the new position after move")
	}

	if game.Board.GetPieceAt(move.GetFrom()) != nil {
		t.Errorf("Expected original position to be empty after move")
	}

	if game.CurrentPlayer != &player2 {
		t.Errorf("Expected current player to be player2 after move, but got %v", game.CurrentPlayer)
	}
}

func TestMakeMoveWithWinningAttack(t *testing.T) {
	player1 := NewPlayer(1, "Alice", "red")
	controller1 := NewHumanPlayerController(&player1)
	player2 := NewPlayer(2, "Bob", "blue")
	controller2 := NewHumanPlayerController(&player2)
	game := NewGame(controller1, controller2)

	attacker := NewPiece(models.Captain, &player1)
	defender := NewPiece(models.Scout, &player2)

	game.Board.SetPieceAt(NewPosition(0, 0), attacker)
	game.Board.SetPieceAt(NewPosition(0, 1), defender)
	move := NewMove(NewPosition(0, 0), NewPosition(0, 1), &player1)

	game.MakeMove(&move, attacker)

	if game.Board.GetPieceAt(move.GetTo()) != attacker {
		t.Errorf("Expected attacker to be at the new position after winning attack")
	}

	if game.Board.GetPieceAt(move.GetFrom()) != nil {
		t.Errorf("Expected original position to be empty after move")
	}

	if !attacker.IsAlive() {
		t.Errorf("Expected attacker to be alive after winning attack")
	}

	if defender.IsAlive() {
		t.Errorf("Expected defender to be dead after losing attack")
	}

	if game.CurrentPlayer != &player2 {
		t.Errorf("Expected current player to be player2 after move, but got %v", game.CurrentPlayer)
	}
}

func TestMakeMoveWithLosingAttack(t *testing.T) {
	player1 := NewPlayer(1, "Alice", "red")
	controller1 := NewHumanPlayerController(&player1)
	player2 := NewPlayer(2, "Bob", "blue")
	controller2 := NewHumanPlayerController(&player2)
	game := NewGame(controller1, controller2)

	attacker := NewPiece(models.Scout, &player1)
	defender := NewPiece(models.Captain, &player2)

	game.Board.SetPieceAt(NewPosition(0, 0), attacker)
	game.Board.SetPieceAt(NewPosition(0, 1), defender)
	move := NewMove(NewPosition(0, 0), NewPosition(0, 1), &player1)

	game.MakeMove(&move, attacker)

	if game.Board.GetPieceAt(move.GetTo()) != defender {
		t.Errorf("Expected defender to remain at the position after winning attack")
	}

	if game.Board.GetPieceAt(move.GetFrom()) != nil {
		t.Errorf("Expected original position to be empty after move")
	}

	if attacker.IsAlive() {
		t.Errorf("Expected attacker to be dead after losing attack")
	}

	if !defender.IsAlive() {
		t.Errorf("Expected defender to be dead after losing attack")
	}

	if game.CurrentPlayer != &player2 {
		t.Errorf("Expected current player to be player2 after move, but got %v", game.CurrentPlayer)
	}
}
