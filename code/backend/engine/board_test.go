package engine_test

import (
	"digital-innovation/stratego/engine"
	"digital-innovation/stratego/models"
	"testing"
)

func TestNewBoard(t *testing.T) {
	board := engine.NewBoard()

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

	expectedLakes := []engine.Position{
		engine.NewPosition(2, 4), engine.NewPosition(3, 4), engine.NewPosition(2, 5), engine.NewPosition(3, 5),
		engine.NewPosition(6, 4), engine.NewPosition(7, 4), engine.NewPosition(6, 5), engine.NewPosition(7, 5),
	}

	for _, lakePos := range expectedLakes {
		if !board.IsLake(lakePos) {
			t.Errorf("Expected position %v to be a lake, but it is not", lakePos)
		}
	}
}

func TestSetAndGetPieceAt(t *testing.T) {
	board := engine.NewBoard()
	player := engine.NewPlayer(1, "Alice", "red")
	piece := engine.NewPiece(models.Marshal, &player)
	position := engine.NewPosition(0, 0)

	board.SetPieceAt(position, piece)
	retrievedPiece := board.GetPieceAt(position)

	if retrievedPiece != piece {
		t.Errorf("Expected to retrieve the same piece that was set, but got a different piece")
	}
}

func TestGetPieceAt(t *testing.T) {
	board := engine.NewBoard()
	player := engine.NewPlayer(1, "Alice", "red")
	piece := engine.NewPiece(models.Marshal, &player)
	position := engine.NewPosition(1, 1)

	board.SetPieceAt(position, piece)
	retrievedPiece := board.GetPieceAt(position)

	if retrievedPiece != piece {
		t.Errorf("Expected to retrieve the same piece that was set, but got a different piece")
	}
}

func TestIsLake(t *testing.T) {
	board := engine.NewBoard()
	position := engine.NewPosition(2, 4)

	isLake := board.IsLake(position)

	if !isLake {
		t.Errorf("Expected position %v to be a lake, but it is not", position)
	}
}

func TestIsValidMove(t *testing.T) {
	board := engine.NewBoard()
	player := engine.NewPlayer(1, "Alice", "red")
	piece := engine.NewPiece(models.Marshal, &player)
	fromPos := engine.NewPosition(0, 0)
	toPos := engine.NewPosition(0, 1)

	board.SetPieceAt(fromPos, piece)
	move := engine.NewMove(fromPos, toPos, &player)

	if !board.IsValidMove(&move) {
		t.Errorf("Expected move from %v to %v to be valid, but it is not", fromPos, toPos)
	}

	lakePos := engine.NewPosition(2, 4)
	moveToLake := engine.NewMove(fromPos, lakePos, &player)

	if board.IsValidMove(&moveToLake) {
		t.Errorf("Expected move to lake position %v to be invalid, but it is valid", lakePos)
	}
}

func TestIsInvalidMoveOutsideField(t *testing.T) {
	board := engine.NewBoard()
	player := engine.NewPlayer(1, "Alice", "red")
	piece := engine.NewPiece(models.Marshal, &player)
	fromPos := engine.NewPosition(0, 0)
	toPos := engine.NewPosition(10, 10) // Outside the board

	board.SetPieceAt(fromPos, piece)
	move := engine.NewMove(fromPos, toPos, &player)

	if board.IsValidMove(&move) {
		t.Errorf("Expected move from %v to %v to be invalid (outside field), but it is valid", fromPos, toPos)
	}
}

func TestIsInvalidMoveIntoLake(t *testing.T) {
	board := engine.NewBoard()
	player := engine.NewPlayer(1, "Alice", "red")
	piece := engine.NewPiece(models.Marshal, &player)
	fromPos := engine.NewPosition(1, 4)
	toPos := engine.NewPosition(2, 4) // Lake position

	board.SetPieceAt(fromPos, piece)
	move := engine.NewMove(fromPos, toPos, &player)

	if board.IsValidMove(&move) {
		t.Errorf("Expected move from %v to lake position %v to be invalid, but it is valid", fromPos, toPos)
	}
}

func TestIsInvalidMoveToTeamPiece(t *testing.T) {
	board := engine.NewBoard()
	player := engine.NewPlayer(1, "Alice", "red")
	piece1 := engine.NewPiece(models.Marshal, &player)
	piece2 := engine.NewPiece(models.General, &player)
	fromPos := engine.NewPosition(0, 0)
	toPos := engine.NewPosition(0, 1)

	board.SetPieceAt(fromPos, piece1)
	board.SetPieceAt(toPos, piece2)
	move := engine.NewMove(fromPos, toPos, &player)

	if board.IsValidMove(&move) {
		t.Errorf("Expected move from %v to %v to be invalid (to own piece), but it is valid", fromPos, toPos)
	}
}

func TestRandomizeSetup(t *testing.T) {
	// unimplemented
}

func TestSwapPieces(t *testing.T) {
	board := engine.NewBoard()
	player := engine.NewPlayer(1, "Alice", "red")
	piece1 := engine.NewPiece(models.Marshal, &player)
	piece2 := engine.NewPiece(models.General, &player)
	pos1 := engine.NewPosition(0, 0)
	pos2 := engine.NewPosition(0, 1)

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
	board := engine.NewBoard()
	player := engine.NewPlayer(1, "Alice", "red")
	piece1 := engine.NewPiece(models.Marshal, &player)
	pos1 := engine.NewPosition(0, 0)
	pos2 := engine.NewPosition(0, 1)

	board.SetPieceAt(pos1, piece1)

	err := board.SwapPieces(pos1, pos2)

	if err == nil {
		t.Errorf("Expected swap to fail due to invalid position, but it succeeded")
	}
}

func TestRemovePieceAt(t *testing.T) {
	board := engine.NewBoard()
	player := engine.NewPlayer(1, "Alice", "red")
	piece := engine.NewPiece(models.Marshal, &player)
	position := engine.NewPosition(0, 0)

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
	board := engine.NewBoard()
	position := engine.NewPosition(0, 0)

	err := board.RemovePieceAt(position)

	if err == nil {
		t.Errorf("Expected removal to fail for non-existent piece, but it succeeded")
	}
}

func TestRemoveAlivePiece(t *testing.T) {
	board := engine.NewBoard()
	player := engine.NewPlayer(1, "Alice", "red")
	piece := engine.NewPiece(models.Marshal, &player)
	position := engine.NewPosition(0, 0)

	board.SetPieceAt(position, piece)
	// piece is still alive
	err := board.RemovePieceAt(position)

	if err == nil {
		t.Errorf("Expected removal to fail for alive piece, but it succeeded")
	}
}

func TestMovePiece(t *testing.T) {
	board := engine.NewBoard()
	player := engine.NewPlayer(1, "Alice", "red")
	piece := engine.NewPiece(models.Marshal, &player)
	fromPos := engine.NewPosition(0, 0)
	toPos := engine.NewPosition(0, 1)

	board.SetPieceAt(fromPos, piece)
	move := engine.NewMove(fromPos, toPos, &player)

	board.MovePiece(&move, piece)

	if board.GetPieceAt(toPos) != piece {
		t.Errorf("Expected piece to be at %v after move, but got %v", toPos, board.GetPieceAt(toPos))
	}

	if board.GetPieceAt(fromPos) != nil {
		t.Errorf("Expected original position %v to be empty after move, but got %v", fromPos, board.GetPieceAt(fromPos))
	}
}

func TestListMovesStandardPiece(t *testing.T) {
	board := engine.NewBoard()
	player := engine.NewPlayer(1, "Alice", "red")
	piece := engine.NewPiece(models.Marshal, &player)
	position := engine.NewPosition(1, 1)

	board.SetPieceAt(position, piece)

	moves, err := board.ListMoves(position)
	if err != nil {
		t.Errorf("Expected to list moves successfully, but got error: %v", err)
	}

	expectedMoveCount := 4 // Up, Down, Left, Right
	if len(moves) != expectedMoveCount {
		t.Errorf("Expected %d moves for standard piece, but got %d", expectedMoveCount, len(moves))
	}
}

func TestListMovesStandardPieceNoMovesAvailable(t *testing.T) {
	board := engine.NewBoard()
	player := engine.NewPlayer(1, "Alice", "red")
	piece := engine.NewPiece(models.Marshal, &player)
	position := engine.NewPosition(1, 1)
	blockedPositions := []engine.Position{
		engine.NewPosition(1, 0), // Up
		engine.NewPosition(1, 2), // Down
		engine.NewPosition(0, 1), // Left
		engine.NewPosition(2, 1), // Right
	}

	for _, pos := range blockedPositions {
		blockingPiece := engine.NewPiece(models.Scout, &player)
		board.SetPieceAt(pos, blockingPiece)
	}

	board.SetPieceAt(position, piece)

	moves, err := board.ListMoves(position)
	if err != nil {
		t.Errorf("Expected to list moves successfully, but got error: %v", err)
	}

	if len(moves) != 0 {
		t.Errorf("Expected no moves for standard piece, but got %d", len(moves))
	}
}

func TestListMovesScout(t *testing.T) {
	board := engine.NewBoard()
	player := engine.NewPlayer(1, "Alice", "red")
	piece := engine.NewPiece(models.Scout, &player)
	position := engine.NewPosition(0, 0)

	board.SetPieceAt(position, piece)

	moves, err := board.ListMoves(position)
	if err != nil {
		t.Errorf("Expected to list moves successfully, but got error: %v", err)
	}

	expectedMoveCount := 18 // 9 Down, 9 Right
	if len(moves) != expectedMoveCount {
		t.Errorf("Expected %d moves for scout piece, but got %d", expectedMoveCount, len(moves))
	}
}

func TestListMovesNoPiece(t *testing.T) {
	board := engine.NewBoard()
	position := engine.NewPosition(0, 0)

	moves, err := board.ListMoves(position)
	if err == nil {
		t.Errorf("Expected error when listing moves for empty position, but got none")
	}

	if len(moves) != 0 {
		t.Errorf("Expected no moves for empty position, but got %d", len(moves))
	}
}
