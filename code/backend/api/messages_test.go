package api_test

import (
	"digital-innovation/stratego/api"
	"digital-innovation/stratego/engine"
	"digital-innovation/stratego/models"
	"testing"
)

func TestPieceToDTO(t *testing.T) {
	player := engine.NewPlayer(0, "TestPlayer", "red")
	piece := engine.NewPiece(models.Captain, &player)

	// Test as owner
	dto := api.PieceToDTO(piece, 0)
	if dto.Type != "Captain" {
		t.Errorf("Expected type Captain, got: %s", dto.Type)
	}
	if dto.Rank != "6" {
		t.Errorf("Expected rank '6' for Captain, got: %s", dto.Rank)
	}
	if dto.OwnerID != 0 {
		t.Errorf("Expected owner 0, got: %d", dto.OwnerID)
	}

	// Test as opponent - piece should be hidden
	dtoHidden := api.PieceToDTO(piece, 1)
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
	dto := api.PieceToDTO(piece, 1)
	if dto.Type != "Scout" {
		t.Errorf("Expected type Scout for revealed piece, got: %s", dto.Type)
	}
	if !dto.Revealed {
		t.Error("Expected Revealed to be true")
	}
}

func TestPieceToDTONil(t *testing.T) {
	dto := api.PieceToDTO(nil, 0)

	if dto.Type != "" {
		t.Errorf("Expected empty type for nil piece, got: %s", dto.Type)
	}
	if dto.OwnerID != 0 {
		t.Errorf("Expected OwnerID 0 for nil piece, got: %d", dto.OwnerID)
	}
}

func TestPositionToDTO(t *testing.T) {
	pos := engine.NewPosition(3, 7)
	dto := api.PositionToDTO(pos)

	if dto.X != 3 {
		t.Errorf("Expected X=3, got: %d", dto.X)
	}
	if dto.Y != 7 {
		t.Errorf("Expected Y=7, got: %d", dto.Y)
	}
}

func TestMoveToDTO(t *testing.T) {
	player := engine.NewPlayer(0, "TestPlayer", "red")
	from := engine.NewPosition(2, 6)
	to := engine.NewPosition(2, 5)
	move := engine.NewMove(from, to, &player)

	dto := api.MoveToDTO(move)

	if dto.From.X != 2 || dto.From.Y != 6 {
		t.Errorf("Expected from position (2,6), got: (%d,%d)", dto.From.X, dto.From.Y)
	}
	if dto.To.X != 2 || dto.To.Y != 5 {
		t.Errorf("Expected to position (2,5), got: (%d,%d)", dto.To.X, dto.To.Y)
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
		dto := api.PieceToDTO(piece, 0)

		if dto.Type != pieceType.GetName() {
			t.Errorf("Expected type %s, got: %s", pieceType.GetName(), dto.Type)
		}
		if dto.Icon != pieceType.GetIcon() {
			t.Errorf("Expected icon %s for %s, got: %s", pieceType.GetIcon(), pieceType.GetName(), dto.Icon)
		}
	}
}
