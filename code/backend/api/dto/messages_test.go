package dto_test

import (
	"digital-innovation/gostrategy/api/dto"
	"digital-innovation/gostrategy/engine"
	"digital-innovation/gostrategy/models"
	"testing"
)

func TestPieceToDTO(t *testing.T) {
	player := engine.NewPlayer(0, "TestPlayer", "red")
	piece := engine.NewPiece(models.Captain, &player)

	// Test as owner
	d := dto.PieceToDTO(piece, 0, false)
	if d.Type != "Captain" {
		t.Errorf("Expected type Captain, got: %s", d.Type)
	}
	if d.Rank != "6" {
		t.Errorf("Expected rank '6' for Captain, got: %s", d.Rank)
	}
	if d.OwnerID != 0 {
		t.Errorf("Expected owner 0, got: %d", d.OwnerID)
	}

	// Test as opponent - piece should be hidden
	dtoHidden := dto.PieceToDTO(piece, 1, false)
	if dtoHidden.Type != "" {
		t.Errorf("Expected type to be empty for opponent, got: %s", dtoHidden.Type)
	}
	if dtoHidden.Rank != "" {
		t.Errorf("Expected rank to be empty for hidden piece, got: %s", dtoHidden.Rank)
	}
	if dtoHidden.OwnerID != 0 {
		t.Errorf("Expected owner ID to still be 0, got: %d", dtoHidden.OwnerID)
	}
}

func TestPieceToDTORevealed(t *testing.T) {
	player := engine.NewPlayer(0, "TestPlayer", "red")
	piece := engine.NewPiece(models.Scout, &player)
	piece.Reveal()

	// Even as opponent, revealed piece should show details
	d := dto.PieceToDTO(piece, 1, false)
	if d.Type != "Scout" {
		t.Errorf("Expected type Scout for revealed piece, got: %s", d.Type)
	}
	if !d.Revealed {
		t.Error("Expected Revealed to be true")
	}
}

func TestPieceToDTONil(t *testing.T) {
	d := dto.PieceToDTO(nil, 0, false)

	if d.Type != "" {
		t.Errorf("Expected empty type for nil piece, got: %s", d.Type)
	}
	if d.OwnerID != -1 {
		t.Errorf("Expected OwnerID -1 for nil piece, got: %d", d.OwnerID)
	}
}

func TestPositionToDTO(t *testing.T) {
	pos := engine.NewPosition(3, 7)
	d := dto.PositionToDTO(pos)

	if d.X != 3 {
		t.Errorf("Expected X=3, got: %d", d.X)
	}
	if d.Y != 7 {
		t.Errorf("Expected Y=7, got: %d", d.Y)
	}
}

func TestMoveToDTO(t *testing.T) {
	player := engine.NewPlayer(0, "TestPlayer", "red")
	from := engine.NewPosition(2, 6)
	to := engine.NewPosition(2, 5)
	move := engine.NewMove(from, to, &player)

	d := dto.MoveToDTO(move)

	if d.From.X != 2 || d.From.Y != 6 {
		t.Errorf("Expected from position (2,6), got: (%d,%d)", d.From.X, d.From.Y)
	}
	if d.To.X != 2 || d.To.Y != 5 {
		t.Errorf("Expected to position (2,5), got: (%d,%d)", d.To.X, d.To.Y)
	}
}

func TestPieceToDTOAllPieceTypes(t *testing.T) {
	player := engine.NewPlayer(0, "TestPlayer", "red")

	pieceTypes := []models.PieceType{
		models.Flag,
		models.Bomb,
		models.Spy,
		models.Scout,
		models.Miner,
		models.Sergeant,
		models.Lieutenant,
		models.Captain,
		models.Major,
		models.Colonel,
		models.General,
		models.Marshal,
	}

	for _, pieceType := range pieceTypes {
		piece := engine.NewPiece(pieceType, &player)
		d := dto.PieceToDTO(piece, 0, false)

		if d.Type != pieceType.GetName() {
			t.Errorf("Expected type %s, got: %s", pieceType.GetName(), d.Type)
		}
		if d.Icon != pieceType.GetIcon() {
			t.Errorf("Expected icon %s for %s, got: %s", pieceType.GetIcon(), pieceType.GetName(), d.Icon)
		}
	}
}
