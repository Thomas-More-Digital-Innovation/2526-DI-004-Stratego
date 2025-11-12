package fafo_test

import (
	"digital-innovation/stratego/ai/fafo"
	"digital-innovation/stratego/engine"
	"digital-innovation/stratego/game"
	"testing"
)

func TestPickRandomPiece(t *testing.T) {
	player1 := engine.NewPlayer(1, "Alice", "red")
	controller1 := fafo.NewFafoAI(&player1, false)

	player2 := engine.NewPlayer(2, "Bob", "blue")
	controller2 := fafo.NewFafoAI(&player2, false)

	game.QuickStart(controller1, controller2)

	piece1 := controller1.PickRandomPiece()

	if piece1 == nil {
		t.Errorf("Expected to pick a piece, but got nil")
	}

	if piece1.GetOwner() != &player1 {
		t.Errorf("Expected piece owner to be Alice, but got %v", piece1.GetOwner().GetName())
	}

	piece2 := controller2.PickRandomPiece()

	if piece2 == nil {
		t.Errorf("Expected to pick a piece, but got nil")
	}

	if piece2.GetOwner() != &player2 {
		t.Errorf("Expected piece owner to be Bob, but got %v", piece2.GetOwner().GetName())
	}
}

func TestMakeMove(t *testing.T) {
	player1 := engine.NewPlayer(1, "Alice", "red")
	controller1 := fafo.NewFafoAI(&player1, false)

	player2 := engine.NewPlayer(2, "Bob", "blue")
	controller2 := fafo.NewFafoAI(&player2, false)

	g := game.QuickStart(controller1, controller2)

	move1 := controller1.MakeMove(g.Board)
	piece := g.Board.GetPieceAt(move1.GetFrom())

	if piece.GetOwner() != &player1 {
		t.Errorf("Expected piece owner to be Alice, but got %v", piece.GetOwner().GetName())
	}

	move2 := controller2.MakeMove(g.Board)
	piece2 := g.Board.GetPieceAt(move2.GetFrom())

	if piece2.GetOwner() != &player2 {
		t.Errorf("Expected piece owner to be Bob, but got %v", piece2.GetOwner().GetName())
	}

}

func TestNoMovesLeft(t *testing.T) {
	player1 := engine.NewPlayer(1, "Alice", "red")
	controller1 := fafo.NewFafoAI(&player1, false)

	player2 := engine.NewPlayer(2, "Bob", "blue")
	controller2 := fafo.NewFafoAI(&player2, false)

	g := game.QuickStart(controller1, controller2)

	board := g.Board

	for _, y := range board.GetField() {
		for _, piece := range y {

			if piece != nil && piece.GetOwner() == &player1 {
				piece.Eliminate()
			}

		}
	}

	if len(player1.GetAlivePieces()) != 0 {
		t.Errorf("Expected all pieces to be eliminated, but got %d", len(player1.GetAlivePieces()))
	}

	move1 := controller1.MakeMove(g.Board)

	if !move1.IsEmpty() {
		t.Errorf("Expected no move to be made, but got %v", move1)
	}

	move2 := controller2.MakeMove(g.Board)
	piece2 := g.Board.GetPieceAt(move2.GetFrom())

	if piece2 == nil {
		t.Fatalf("Expected a piece to be selected, but got nil")
	}

	if piece2.GetOwner() != &player2 {
		t.Errorf("Expected piece owner to be Bob, but got %v", piece2.GetOwner().GetName())
	}

}
