package game

import (
	"digital-innovation/gostrategy/internal/models"
	"testing"
)

func TestMovePiece(t *testing.T) {
	board := NewBoard()
	player := NewPlayer(1, "Alice", "red")
	piece := NewPiece(models.Marshal, &player)
	fromPos := NewPosition(0, 0)
	toPos := NewPosition(0, 1)

	board.SetPieceAt(fromPos, piece)
	move := NewMove(fromPos, toPos, &player)

	board.MovePiece(&move, piece)

	if board.GetPieceAt(toPos) != piece {
		t.Errorf("Expected piece to be at %v after move, but got %v", toPos, board.GetPieceAt(toPos))
	}

	if board.GetPieceAt(fromPos) != nil {
		t.Errorf("Expected original position %v to be empty after move, but got %v", fromPos, board.GetPieceAt(fromPos))
	}
}

func TestListMovesStandardPiece(t *testing.T) {
	board := NewBoard()
	player := NewPlayer(1, "Alice", "red")
	piece := NewPiece(models.Marshal, &player)
	position := NewPosition(1, 1)

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
	board := NewBoard()
	player := NewPlayer(1, "Alice", "red")
	piece := NewPiece(models.Marshal, &player)
	position := NewPosition(1, 1)
	blockedPositions := []Position{
		NewPosition(1, 0), // Up
		NewPosition(1, 2), // Down
		NewPosition(0, 1), // Left
		NewPosition(2, 1), // Right
	}

	for _, pos := range blockedPositions {
		blockingPiece := NewPiece(models.Scout, &player)
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
	board := NewBoard()
	player := NewPlayer(1, "Alice", "red")
	piece := NewPiece(models.Scout, &player)
	position := NewPosition(0, 0)

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
	board := NewBoard()
	position := NewPosition(0, 0)

	moves, err := board.ListMoves(position)
	if err == nil {
		t.Errorf("Expected error when listing moves for empty position, but got none")
	}

	if len(moves) != 0 {
		t.Errorf("Expected no moves for empty position, but got %d", len(moves))
	}
}

func TestClone(t *testing.T) {
	board := NewBoard()
	player := NewPlayer(1, "Alice", "red")
	piece := NewPiece(models.Marshal, &player)
	pos := NewPosition(0, 0)
	board.SetPieceAt(pos, piece)

	cloned := board.Clone()
	clonedPiece := cloned.GetPieceAt(pos)
	if clonedPiece == nil {
		t.Fatal("Cloned board should have a piece at the same position")
	}
	if clonedPiece == piece {
		t.Error("Cloned board should have a NEW piece pointer (deep copy)")
	}
	if clonedPiece.GetType().GetRank() != piece.GetType().GetRank() {
		t.Error("Cloned piece should have the same rank")
	}
}

func TestBoardString(t *testing.T) {
	board := NewBoard()
	player := NewPlayer(1, "Alice", "red")
	piece := NewPiece(models.Marshal, &player)
	board.SetPieceAt(NewPosition(0, 0), piece)
	s := board.String()
	if s == "" {
		t.Error("Board.String() returned empty string")
	}
	// Check if it contains pieces, lakes and empty squares
	if !contains(s, piece.GetType().GetIcon()) {
		t.Error("Board string should contain piece icon")
	}
	if !contains(s, "~~") {
		t.Error("Board string should contain lake icon")
	}
	if !contains(s, "..") {
		t.Error("Board string should contain empty square icon")
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestIsValidMoveBoundaries(t *testing.T) {
	board := NewBoard()
	player := NewPlayer(1, "P1", "")
	piece := NewPiece(models.Marshal, &player)
	from := NewPosition(0, 0)
	board.SetPieceAt(from, piece)

	tests := []struct {
		x, y int
		want bool
	}{
		{-1, 0, false},
		{0, -1, false},
		{10, 0, false},
		{0, 10, false},
		{0, 0, false}, // same as from (blocked by team piece)
	}

	for _, tt := range tests {
		move := NewMove(from, NewPosition(tt.x, tt.y), &player)
		if board.IsValidMove(&move) != tt.want {
			t.Errorf("IsValidMove for (%d, %d) should be %v", tt.x, tt.y, tt.want)
		}
	}

	// No piece at from
	moveNoPiece := NewMove(NewPosition(5, 5), NewPosition(5, 6), &player)
	if board.IsValidMove(&moveNoPiece) {
		t.Error("IsValidMove should be false if no piece at from")
	}
}

func TestScoutMovesOutOfBounds(t *testing.T) {
	board := NewBoard()
	player := NewPlayer(1, "P1", "")
	scout := NewPiece(models.Scout, &player)
	from := NewPosition(0, 0)
	board.SetPieceAt(from, scout)

	// This will trigger handleScoutMoves and it will hit boundaries (IsValidMove will return false)
	moves, _ := board.ListMoves(from)
	if len(moves) == 0 {
		t.Error("Scout should have moves")
	}
}

func TestEncodeSetupEmptyCells(t *testing.T) {
	// Test EncodeSetup with '.' and ' '
	rows := []string{
		validSetupPrefix,
		"2222223333",
		"3444455556",
		"66677788. ",
	}
	data, err := EncodeSetup(rows, 1)
	if err != nil {
		t.Fatal(err)
	}
	// Cell 38 and 39 should not be occupied
	if data[38]&BitOccupied != 0 || data[39]&BitOccupied != 0 {
		t.Error("Empty cells should not be occupied in binary format")
	}

	decoded, _, _ := DecodeSetup(data)
	if decoded[3][8] != '.' || decoded[3][9] != '.' {
		t.Error("Empty cells should decode back to '.'")
	}
}

func TestIsValidMoveScoutLongDistance(t *testing.T) {
	board := NewBoard()
	player := NewPlayer(1, "Alice", "red")
	scout := NewPiece(models.Scout, &player)
	from := NewPosition(0, 0)
	to := NewPosition(0, 5)
	board.SetPieceAt(from, scout)

	move := NewMove(from, to, &player)
	if !board.IsValidMove(&move) {
		t.Error("Scout should be able to move long distance according to IsValidMove")
	}
}

func TestListMovesScoutBlocked(t *testing.T) {
	board := NewBoard()
	player := NewPlayer(1, "Alice", "red")
	scout := NewPiece(models.Scout, &player)
	from := NewPosition(0, 0)
	board.SetPieceAt(from, scout)

	// Blocked by lake
	fromLakeSide := NewPosition(2, 3)
	board.SetPieceAt(fromLakeSide, scout)

	moves, _ := board.ListMoves(fromLakeSide)
	for _, m := range moves {
		if m.GetTo().Y > 3 && m.GetTo().X == 2 {
			t.Errorf("Scout should be blocked by lake, but got move to %v", m.GetTo())
		}
	}

	// Blocked by piece
	board = NewBoard()
	board.SetPieceAt(from, scout)
	blocker := NewPiece(models.Miner, &player)
	board.SetPieceAt(NewPosition(0, 3), blocker)

	moves, _ = board.ListMoves(from)
	for _, m := range moves {
		if m.GetTo().Y > 3 && m.GetTo().X == 0 {
			t.Errorf("Scout should be blocked by piece, but got move to %v", m.GetTo())
		}
	}
}
