package engine

import (
	"digital-innovation/stratego/models"
	"testing"
)

func TestEncodeBoard(t *testing.T) {
	board := NewBoard()
	player1 := NewPlayer(1, "Player1", "red")
	player2 := NewPlayer(2, "Player2", "blue")

	board.SetPieceAt(NewPosition(0, 0), NewPiece(models.Flag, &player1))
	board.SetPieceAt(NewPosition(1, 0), NewPiece(models.Marshal, &player2))

	moved := map[Position]bool{NewPosition(1, 0): true}
	data := EncodeBoard(board, moved)

	if len(data) != 100 {
		t.Errorf("Expected 100 bytes, got %d", len(data))
	}

	// Check Flag encoding
	cell0 := data[0]
	if cell0&BitOccupied == 0 {
		t.Error("Position (0,0) should be occupied")
	}
	if (cell0&MaskPieceType)>>ShiftPieceType != PieceIDFlag {
		t.Errorf("Position (0,0) should be Flag, but got %v| %v", (cell0&MaskPieceType)>>ShiftPieceType, cell0)
	}
	if cell0&BitColor != 0 {
		t.Error("Position (0,0) should be Player 1, but got Player 2")
	}
	if cell0&BitMoved != 0 {
		t.Error("Position (0,0) should not have moved, but got moved")
	}

	// Check Marshal encoding
	cell1 := data[1]
	if cell1&BitOccupied == 0 {
		t.Error("Position (1,0) should be occupied")
	}
	if (cell1&MaskPieceType)>>ShiftPieceType != PieceIDMarshal {
		t.Errorf("Position (1,0) should be Marshal, but got %v", (cell1&MaskPieceType)>>ShiftPieceType)
	}
	if cell1&BitColor == 0 {
		t.Error("Position (1,0) should be Player 2, but got Player 1")
	}
	if cell1&BitMoved == 0 {
		t.Error("Position (1,0) should have moved, but got not moved")
	}
}

func TestDecodeBoard(t *testing.T) {
	player1 := NewPlayer(1, "Player1", "red")
	player2 := NewPlayer(2, "Player2", "blue")

	data := make([]byte, 100)
	data[0] = BitOccupied | (PieceIDFlag << ShiftPieceType)
	data[1] = BitOccupied | (PieceIDMarshal << ShiftPieceType) | BitColor | BitMoved

	board, moved := DecodeBoard(data, &player1, &player2)

	piece0 := board.GetPieceAt(NewPosition(0, 0))
	if piece0 == nil {
		t.Fatal("Position (0,0) should have a piece")
	}
	if piece0.GetType().GetRank() != '0' {
		t.Error("Position (0,0) should be Flag")
	}
	if piece0.GetOwner().GetID() != 1 {
		t.Error("Position (0,0) should belong to Player 1")
	}
	if moved[NewPosition(0, 0)] {
		t.Error("Position (0,0) should not have moved")
	}

	piece1 := board.GetPieceAt(NewPosition(1, 0))
	if piece1 == nil {
		t.Fatal("Position (1,0) should have a piece")
	}
	if piece1.GetType().GetRank() != 'M' {
		t.Error("Position (1,0) should be Marshal")
	}
	if piece1.GetOwner().GetID() != 2 {
		t.Error("Position (1,0) should belong to Player 2")
	}
	if !moved[NewPosition(1, 0)] {
		t.Error("Position (1,0) should have moved")
	}
}

func TestBoardBase64RoundTrip(t *testing.T) {
	board := NewBoard()
	player1 := NewPlayer(1, "Player1", "red")
	player2 := NewPlayer(2, "Player2", "blue")

	board.SetPieceAt(NewPosition(0, 0), NewPiece(models.Flag, &player1))
	board.SetPieceAt(NewPosition(9, 9), NewPiece(models.Marshal, &player2))

	encoded := EncodeBoardToBase64(board, nil)
	decoded, _, err := DecodeBoardFromBase64(encoded, &player1, &player2)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if decoded.GetPieceAt(NewPosition(0, 0)).GetType().GetRank() != '0' {
		t.Error("Flag not preserved")
	}
	if decoded.GetPieceAt(NewPosition(9, 9)).GetType().GetRank() != 'M' {
		t.Error("Marshal not preserved")
	}
}

func TestEncodeSetup(t *testing.T) {
	rows := []string{"BBBBBB2222", "4444555566", "6677788399", "M311110333"}

	data, err := EncodeSetup(rows, 1)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}
	if len(data) != 40 {
		t.Errorf("Expected 40 bytes, got %d", len(data))
	}

	// Check first cell (Bomb)
	cell0 := data[0]
	if cell0&BitOccupied == 0 {
		t.Error("First cell should be occupied")
	}
	if (cell0&MaskPieceType)>>ShiftPieceType != PieceIDBomb {
		t.Error("First cell should be Bomb")
	}
	if cell0&BitColor != 0 {
		t.Error("First cell should be Player 1")
	}
}

func TestDecodeSetup(t *testing.T) {
	originalRows := []string{"BBBBBB2222", "4444555566", "6677788399", "M311110333"}

	data, _ := EncodeSetup(originalRows, 1)
	rows, playerID, err := DecodeSetup(data)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if playerID != 1 {
		t.Errorf("Expected player ID 1, got %d", playerID)
	}

	for i, row := range rows {
		if row != originalRows[i] {
			t.Errorf("Row %d mismatch: expected %s, got %s", i, originalRows[i], row)
		}
	}
}

func TestValidateSetup(t *testing.T) {
	// Valid: 1 Flag, 6 Bombs, 1 Spy, 8 Scouts, 5 Miners, 4 Sergeants, 4 Lieutenants, 4 Captains, 3 Majors, 2 Colonels, 1 General, 1 Marshal
	valid := []string{
		"0BBBBBB122", // 6 Bombs + 4 Scouts = 10
		"2222223333", // 4 Scouts + 1 Spy + 4 Sergeants + 1 Lieutenant = 10
		"3444455556", // 3 Lieutenants + 4 Captains + 3 Majors = 10
		"666777889M", // 2 Colonels + 1 General + 1 Marshal + 5 Miners + 1 Flag = 10
	}
	if err := ValidateSetup(valid); err != nil {
		t.Errorf("Valid setup rejected: %v", err)
	}

	// Invalid: 2 Flags instead of 1
	invalid := []string{"BBBBBB2222", "2222224444", "5555666677", "00899M3333"}
	if err := ValidateSetup(invalid); err == nil {
		t.Error("Invalid setup accepted")
	}
}
