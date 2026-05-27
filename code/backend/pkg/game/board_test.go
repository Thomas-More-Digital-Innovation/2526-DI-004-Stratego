package game

import (
	"digital-innovation/gostrategy/pkg/game/models"
	"testing"
)

func TestNewBoard(t *testing.T) {
	board := NewBoard()

	if board == nil {
		t.Errorf("Expected a board to be created, but got nil")
	}

	if len(board.GetField()) != 10 {
		t.Errorf("Expected board field to have 10 rows, but got %d", len(board.GetField()))
	}

	for _, row := range board.GetField() {
		if len(row) != 10 {
			t.Errorf("Expected each row to have 10 columns, but got %d", len(row))
		}
	}

	expectedLakes := []Position{
		NewPosition(2, 4), NewPosition(3, 4), NewPosition(2, 5), NewPosition(3, 5),
		NewPosition(6, 4), NewPosition(7, 4), NewPosition(6, 5), NewPosition(7, 5),
	}

	for _, lakePos := range expectedLakes {
		if !board.IsLake(lakePos) {
			t.Errorf("Expected position %v to be a lake, but it is not", lakePos)
		}
	}
}

func TestSetAndGetPieceAt(t *testing.T) {
	board := NewBoard()
	player := NewPlayer(1, "Alice", "red")
	piece := NewPiece(models.Marshal, &player)
	position := NewPosition(0, 0)

	board.SetPieceAt(position, piece)
	retrievedPiece := board.GetPieceAt(position)

	if retrievedPiece != piece {
		t.Errorf("Expected to retrieve the same piece that was set, but got a different piece")
	}
}

func TestGetPieceAt(t *testing.T) {
	board := NewBoard()
	player := NewPlayer(1, "Alice", "red")
	piece := NewPiece(models.Marshal, &player)
	position := NewPosition(1, 1)

	board.SetPieceAt(position, piece)
	retrievedPiece := board.GetPieceAt(position)

	if retrievedPiece != piece {
		t.Errorf("Expected to retrieve the same piece that was set, but got a different piece")
	}
}

func TestIsLake(t *testing.T) {
	board := NewBoard()
	position := NewPosition(2, 4)

	isLake := board.IsLake(position)

	if !isLake {
		t.Errorf("Expected position %v to be a lake, but it is not", position)
	}
}

func TestIsValidMove(t *testing.T) {
	board := NewBoard()
	player := NewPlayer(1, "Alice", "red")
	piece := NewPiece(models.Marshal, &player)
	fromPos := NewPosition(0, 0)
	toPos := NewPosition(0, 1)

	board.SetPieceAt(fromPos, piece)
	move := NewMove(fromPos, toPos, &player)

	if !board.IsValidMove(&move) {
		t.Errorf("Expected move from %v to %v to be valid, but it is not", fromPos, toPos)
	}

	lakePos := NewPosition(2, 4)
	moveToLake := NewMove(fromPos, lakePos, &player)

	if board.IsValidMove(&moveToLake) {
		t.Errorf("Expected move to lake position %v to be invalid, but it is valid", lakePos)
	}
}

func TestIsInvalidMoveOutsideField(t *testing.T) {
	board := NewBoard()
	player := NewPlayer(1, "Alice", "red")
	piece := NewPiece(models.Marshal, &player)
	fromPos := NewPosition(0, 0)
	toPos := NewPosition(10, 10) // Outside the board

	board.SetPieceAt(fromPos, piece)
	move := NewMove(fromPos, toPos, &player)

	if board.IsValidMove(&move) {
		t.Errorf("Expected move from %v to %v to be invalid (outside field), but it is valid", fromPos, toPos)
	}
}

func TestIsInvalidMoveIntoLake(t *testing.T) {
	board := NewBoard()
	player := NewPlayer(1, "Alice", "red")
	piece := NewPiece(models.Marshal, &player)
	fromPos := NewPosition(1, 4)
	toPos := NewPosition(2, 4) // Lake position

	board.SetPieceAt(fromPos, piece)
	move := NewMove(fromPos, toPos, &player)

	if board.IsValidMove(&move) {
		t.Errorf("Expected move from %v to lake position %v to be invalid, but it is valid", fromPos, toPos)
	}
}

func TestIsInvalidMoveToTeamPiece(t *testing.T) {
	board := NewBoard()
	player := NewPlayer(1, "Alice", "red")
	piece1 := NewPiece(models.Marshal, &player)
	piece2 := NewPiece(models.General, &player)
	fromPos := NewPosition(0, 0)
	toPos := NewPosition(0, 1)

	board.SetPieceAt(fromPos, piece1)
	board.SetPieceAt(toPos, piece2)
	move := NewMove(fromPos, toPos, &player)

	if board.IsValidMove(&move) {
		t.Errorf("Expected move from %v to %v to be invalid (to own piece), but it is valid", fromPos, toPos)
	}
}

func TestRandomizeSetup(_ *testing.T) {
	// unimplemented
}

func TestSwapPieces(t *testing.T) {
	board := NewBoard()
	player := NewPlayer(1, "Alice", "red")
	piece1 := NewPiece(models.Marshal, &player)
	piece2 := NewPiece(models.General, &player)
	pos1 := NewPosition(0, 0)
	pos2 := NewPosition(0, 1)

	board.SetPieceAt(pos1, piece1)
	board.SetPieceAt(pos2, piece2)

	err := board.SwapPieces(pos1, pos2)
	if err != nil {
		t.Errorf("Expected swap to succeed, but got error: %v", err)
	}

	if board.GetPieceAt(pos1) != piece2 {
		t.Errorf("Expected piece at %v to be %v after swap, but got %v", pos1, piece2, board.GetPieceAt(pos1))
	}

	if board.GetPieceAt(pos2) != piece1 {
		t.Errorf("Expected piece at %v to be %v after swap, but got %v", pos2, piece1, board.GetPieceAt(pos2))
	}
}

func TestSwapPiecesInvalidPosition(t *testing.T) {
	board := NewBoard()
	player := NewPlayer(1, "Alice", "red")
	piece1 := NewPiece(models.Marshal, &player)
	pos1 := NewPosition(0, 0)
	pos2 := NewPosition(0, 1)

	board.SetPieceAt(pos1, piece1)

	err := board.SwapPieces(pos1, pos2)

	if err == nil {
		t.Errorf("Expected swap to fail due to invalid position, but it succeeded")
	}
}

func TestRemovePieceAt(t *testing.T) {
	board := NewBoard()
	player := NewPlayer(1, "Alice", "red")
	piece := NewPiece(models.Marshal, &player)
	position := NewPosition(0, 0)

	board.SetPieceAt(position, piece)
	piece.Eliminate()
	err := board.RemovePieceAt(position)

	if err != nil {
		t.Errorf("Expected removal to succeed, but got error: %v", err)
	}

	if board.GetPieceAt(position) != nil {
		t.Errorf("Expected piece at %v to be nil after removal, but got %v", position, board.GetPieceAt(position))
	}
}

func TestRemovePieceThatDoesntExist(t *testing.T) {
	board := NewBoard()
	position := NewPosition(0, 0)

	err := board.RemovePieceAt(position)

	if err == nil {
		t.Errorf("Expected removal to fail for non-existent piece, but it succeeded")
	}
}

func TestRemoveAlivePiece(t *testing.T) {
	board := NewBoard()
	player := NewPlayer(1, "Alice", "red")
	piece := NewPiece(models.Marshal, &player)
	position := NewPosition(0, 0)

	board.SetPieceAt(position, piece)
	// piece is still alive
	err := board.RemovePieceAt(position)

	if err == nil {
		t.Errorf("Expected removal to fail for alive piece, but it succeeded")
	}
}
