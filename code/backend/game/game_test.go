package game_test

import (
	"digital-innovation/stratego/engine"
	"digital-innovation/stratego/game"
	"digital-innovation/stratego/models"
	"testing"
)

func TestNewGame(t *testing.T) {
	player1 := engine.NewPlayer(1, "Alice", "red")
	player2 := engine.NewPlayer(2, "Bob", "blue")
	game := game.NewGame(&player1, &player2)

	if game == nil {
		t.Errorf("Expected a game to be created, but got nil")
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
	player1 := engine.NewPlayer(1, "Alice", "red")
	player2 := engine.NewPlayer(2, "Bob", "blue")
	game := game.NewGame(&player1, &player2)

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

func TestMakeMoveToEmptyCell(t *testing.T) {
	player1 := engine.NewPlayer(1, "Alice", "red")
	player2 := engine.NewPlayer(2, "Bob", "blue")
	game := game.NewGame(&player1, &player2)

	piece := engine.NewPiece(models.Major, &player1)
	move := engine.NewMove(engine.NewPosition(0, 0), engine.NewPosition(0, 1), &player1)

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
	player1 := engine.NewPlayer(1, "Alice", "red")
	player2 := engine.NewPlayer(2, "Bob", "blue")
	game := game.NewGame(&player1, &player2)

	attacker := engine.NewPiece(models.Captain, &player1)
	defender := engine.NewPiece(models.Scout, &player2)

	game.Board.SetPieceAt(engine.NewPosition(0, 0), attacker)
	game.Board.SetPieceAt(engine.NewPosition(0, 1), defender)
	move := engine.NewMove(engine.NewPosition(0, 0), engine.NewPosition(0, 1), &player1)

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
	player1 := engine.NewPlayer(1, "Alice", "red")
	player2 := engine.NewPlayer(2, "Bob", "blue")
	game := game.NewGame(&player1, &player2)

	attacker := engine.NewPiece(models.Scout, &player1)
	defender := engine.NewPiece(models.Captain, &player2)

	game.Board.SetPieceAt(engine.NewPosition(0, 0), attacker)
	game.Board.SetPieceAt(engine.NewPosition(0, 1), defender)
	move := engine.NewMove(engine.NewPosition(0, 0), engine.NewPosition(0, 1), &player1)

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

func TestMakeMoveWithMutualAttack(t *testing.T) {
	player1 := engine.NewPlayer(1, "Alice", "red")
	player2 := engine.NewPlayer(2, "Bob", "blue")
	game := game.NewGame(&player1, &player2)

	attacker := engine.NewPiece(models.Scout, &player1)
	defender := engine.NewPiece(models.Scout, &player2)

	game.Board.SetPieceAt(engine.NewPosition(0, 0), attacker)
	game.Board.SetPieceAt(engine.NewPosition(0, 1), defender)
	move := engine.NewMove(engine.NewPosition(0, 0), engine.NewPosition(0, 1), &player1)

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
