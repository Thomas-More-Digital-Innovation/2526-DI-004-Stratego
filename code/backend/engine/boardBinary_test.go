package engine_test

import (
	"digital-innovation/gostrategy/engine"
	"digital-innovation/gostrategy/models"
	"testing"
)

const (
	testSetupRow     = "BBBBBB2222"
	validSetupPrefix = "0BBBBBB122"
	rowStr           = "row"
)

func TestEncodeBoard(t *testing.T) {
	board := engine.NewBoard()
	player1 := engine.NewPlayer(1, "Player1", "red")
	player2 := engine.NewPlayer(2, "Player2", "blue")

	board.SetPieceAt(engine.NewPosition(0, 0), engine.NewPiece(models.Flag, &player1))
	board.SetPieceAt(engine.NewPosition(1, 0), engine.NewPiece(models.Marshal, &player2))

	moved := map[engine.Position]bool{engine.NewPosition(1, 0): true}
	data := engine.EncodeBoard(board, moved)

	if len(data) != 100 {
		t.Errorf("Expected 100 bytes, got %d", len(data))
	}

	// Check Flag encoding
	cell0 := data[0]
	if cell0&engine.BitOccupied == 0 {
		t.Error("Position (0,0) should be occupied")
	}
	if (cell0&engine.MaskPieceType)>>engine.ShiftPieceType != engine.PieceIDFlag {
		t.Errorf("Position (0,0) should be Flag, but got %v| %v", (cell0&engine.MaskPieceType)>>engine.ShiftPieceType, cell0)
	}
	if cell0&engine.BitColor != 0 {
		t.Error("Position (0,0) should be Player 1, but got Player 2")
	}
	if cell0&engine.BitMoved != 0 {
		t.Error("Position (0,0) should not have moved, but got moved")
	}

	// Check Marshal encoding
	cell1 := data[1]
	if cell1&engine.BitOccupied == 0 {
		t.Error("Position (1,0) should be occupied")
	}
	if (cell1&engine.MaskPieceType)>>engine.ShiftPieceType != engine.PieceIDMarshal {
		t.Errorf("Position (1,0) should be Marshal, but got %v", (cell1&engine.MaskPieceType)>>engine.ShiftPieceType)
	}
	if cell1&engine.BitColor == 0 {
		t.Error("Position (1,0) should be Player 2, but got Player 1")
	}
	if cell1&engine.BitMoved == 0 {
		t.Error("Position (1,0) should have moved, but got not moved")
	}
}

func TestDecodeBoard(t *testing.T) {
	player1 := engine.NewPlayer(1, "Player1", "red")
	player2 := engine.NewPlayer(2, "Player2", "blue")

	data := make([]byte, 100)
	data[0] = engine.BitOccupied | (engine.PieceIDFlag << engine.ShiftPieceType)
	data[1] = engine.BitOccupied | (engine.PieceIDMarshal << engine.ShiftPieceType) | engine.BitColor | engine.BitMoved

	board, moved := engine.DecodeBoard(data, &player1, &player2)

	piece0 := board.GetPieceAt(engine.NewPosition(0, 0))
	if piece0 == nil {
		t.Fatal("Position (0,0) should have a piece")
	}
	if piece0.GetType().GetRank() != '0' {
		t.Error("Position (0,0) should be Flag")
	}
	if piece0.GetOwner().GetID() != 1 {
		t.Error("Position (0,0) should belong to Player 1")
	}
	if moved[engine.NewPosition(0, 0)] {
		t.Error("Position (0,0) should not have moved")
	}

	piece1 := board.GetPieceAt(engine.NewPosition(1, 0))
	if piece1 == nil {
		t.Fatal("Position (1,0) should have a piece")
	}
	if piece1.GetType().GetRank() != 'M' {
		t.Error("Position (1,0) should be Marshal")
	}
	if piece1.GetOwner().GetID() != 2 {
		t.Error("Position (1,0) should belong to Player 2")
	}
	if !moved[engine.NewPosition(1, 0)] {
		t.Error("Position (1,0) should have moved")
	}
}

func TestBoardBase64RoundTrip(t *testing.T) {
	board := engine.NewBoard()
	player1 := engine.NewPlayer(1, "Player1", "red")
	player2 := engine.NewPlayer(2, "Player2", "blue")

	board.SetPieceAt(engine.NewPosition(0, 0), engine.NewPiece(models.Flag, &player1))
	board.SetPieceAt(engine.NewPosition(9, 9), engine.NewPiece(models.Marshal, &player2))

	encoded := engine.EncodeBoardToBase64(board, nil)
	decoded, _, err := engine.DecodeBoardFromBase64(encoded, &player1, &player2)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if decoded.GetPieceAt(engine.NewPosition(0, 0)).GetType().GetRank() != '0' {
		t.Error("Flag not preserved")
	}
	if decoded.GetPieceAt(engine.NewPosition(9, 9)).GetType().GetRank() != 'M' {
		t.Error("Marshal not preserved")
	}
}

func TestEncodeSetup(t *testing.T) {
	rows := []string{testSetupRow, "4444555566", "6677788399", "M311110333"}

	data, err := engine.EncodeSetup(rows, 1)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}
	if len(data) != 40 {
		t.Errorf("Expected 40 bytes, got %d", len(data))
	}

	// Check first cell (Bomb)
	cell0 := data[0]
	if cell0&engine.BitOccupied == 0 {
		t.Error("First cell should be occupied")
	}
	if (cell0&engine.MaskPieceType)>>engine.ShiftPieceType != engine.PieceIDBomb {
		t.Error("First cell should be Bomb")
	}
	if cell0&engine.BitColor != 0 {
		t.Error("First cell should be Player 1")
	}
}

func TestDecodeSetup(t *testing.T) {
	originalRows := []string{testSetupRow, "4444555566", "6677788399", "M311110333"}

	data, _ := engine.EncodeSetup(originalRows, 1)
	rows, playerID, err := engine.DecodeSetup(data)
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
		validSetupPrefix,
		"2222223333",
		"3444455556",
		"666777889M",
	}
	if err := engine.ValidateSetup(valid); err != nil {
		t.Errorf("Valid setup rejected: %v", err)
	}

	// Invalid: 2 Flags instead of 1
	invalid := []string{testSetupRow, "2222224444", "5555666677", "00899M3333"}
	if err := engine.ValidateSetup(invalid); err == nil {
		t.Error("Invalid setup accepted")
	}

	// Wrong row count
	if err := engine.ValidateSetup([]string{rowStr}); err == nil {
		t.Error("Wrong row count should be rejected")
	}

	// Wrong row length
	if err := engine.ValidateSetup([]string{rowStr, rowStr, rowStr, rowStr}); err == nil {
		t.Error("Wrong row length should be rejected")
	}
}

func TestBinaryHelpers(t *testing.T) {
	id, ok := engine.GetPieceIDFromRank('M')
	if !ok || id != engine.PieceIDMarshal {
		t.Error("GetPieceIDFromRank failed for Marshal")
	}

	_, ok = engine.GetPieceIDFromRank('X')
	if ok {
		t.Error("GetPieceIDFromRank should fail for 'X'")
	}

	pt := engine.GetPieceTypeFromID(engine.PieceIDMarshal)
	if pt.GetName() != "Marshal" {
		t.Error("GetPieceTypeFromID failed")
	}

	cell := engine.BitOccupied | (engine.PieceIDMarshal << engine.ShiftPieceType)
	pt2 := engine.GetPieceTypeFromCell(cell)
	if pt2.GetName() != "Marshal" {
		t.Error("GetPieceTypeFromCell failed")
	}
}

func TestBoardBinaryErrors(t *testing.T) {
	player1 := engine.NewPlayer(1, "P1", "")
	player2 := engine.NewPlayer(2, "P2", "")

	_, _, err := engine.DecodeBoardFromBase64("!!!", &player1, &player2)
	if err == nil {
		t.Error("DecodeBoardFromBase64 should fail for invalid base64")
	}

	_, _, err = engine.DecodeBoardFromBase64("YWJjZA==", &player1, &player2) // "abcd" in base64
	if err == nil {
		t.Error("DecodeBoardFromBase64 should fail for wrong data length")
	}

	_, err = engine.EncodeSetup([]string{rowStr}, 1)
	if err == nil {
		t.Error("EncodeSetup should fail for wrong row count")
	}

	_, err = engine.EncodeSetup([]string{"123", "123", "123", "123"}, 1)
	if err == nil {
		t.Error("EncodeSetup should fail for wrong row length")
	}

	_, _, err = engine.DecodeSetup(make([]byte, 10))
	if err == nil {
		t.Error("DecodeSetup should fail for wrong data length")
	}
}

func TestEncodeSetupPlayer2(t *testing.T) {
	rows := []string{"0BBBBBB122", "2222223333", "3444455556", "666777889M"}
	data, err := engine.EncodeSetup(rows, 2)
	if err != nil {
		t.Fatal(err)
	}
	if data[0]&engine.BitColor == 0 {
		t.Error("Player 2 piece should have BitColor set")
	}

	decodedRows, playerID, _ := engine.DecodeSetup(data)
	if playerID != 2 {
		t.Errorf("Expected playerID 2, got %d", playerID)
	}
	if decodedRows[0] != rows[0] {
		t.Error("Setup rows not preserved for Player 2")
	}
}
