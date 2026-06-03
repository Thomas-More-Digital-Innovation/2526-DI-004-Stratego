package dto_test

import (
	"digital-innovation/gostrategy/internal/api/dto"
	"digital-innovation/gostrategy/internal/game"
	gamemodels "digital-innovation/gostrategy/internal/game/models"
	"digital-innovation/gostrategy/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPieceToDTO(t *testing.T) {
	player := game.NewPlayer(0, "TestPlayer", "red")
	piece := game.NewPiece(models.Captain, &player)

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
	player := game.NewPlayer(0, "TestPlayer", "red")
	piece := game.NewPiece(models.Scout, &player)
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
	pos := game.NewPosition(3, 7)
	d := dto.PositionToDTO(pos)

	if d.X != 3 {
		t.Errorf("Expected X=3, got: %d", d.X)
	}
	if d.Y != 7 {
		t.Errorf("Expected Y=7, got: %d", d.Y)
	}
}

func TestMoveToDTO(t *testing.T) {
	player := game.NewPlayer(0, "TestPlayer", "red")
	from := game.NewPosition(2, 6)
	to := game.NewPosition(2, 5)
	move := game.NewMove(from, to, &player)

	d := dto.MoveToDTO(move)

	if d.From.X != 2 || d.From.Y != 6 {
		t.Errorf("Expected from position (2,6), got: (%d,%d)", d.From.X, d.From.Y)
	}
	if d.To.X != 2 || d.To.Y != 5 {
		t.Errorf("Expected to position (2,5), got: (%d,%d)", d.To.X, d.To.Y)
	}
}

func TestPieceToDTOAllPieceTypes(t *testing.T) {
	player := game.NewPlayer(0, "TestPlayer", "red")

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
		piece := game.NewPiece(pieceType, &player)
		d := dto.PieceToDTO(piece, 0, false)

		if d.Type != pieceType.GetName() {
			t.Errorf("Expected type %s, got: %s", pieceType.GetName(), d.Type)
		}
		if d.Icon != pieceType.GetIcon() {
			t.Errorf("Expected icon %s for %s, got: %s", pieceType.GetIcon(), pieceType.GetName(), d.Icon)
		}
	}
}

func TestBuildGameStateMessage(t *testing.T) {
	state := gamemodels.GameState{
		CurrentPlayerID: 1,
		IsGameOver:      true,
		WinnerID:        nil,
	}
	msg := dto.BuildGameStateMessage(state, "Alice", "Capture")
	assert.Equal(t, state, msg.GameState)
	assert.Equal(t, "Alice", msg.WinnerName)
	assert.Equal(t, "Capture", msg.WinCause)
}

func TestMapBoardToDTO(t *testing.T) {
	player := game.NewPlayer(0, "Piet", "red")
	var field [10][10]*game.Piece
	piece := game.NewPiece(models.Marshal, &player)
	field[3][5] = piece

	boardDTO := dto.MapBoardToDTO(field, 0, false)
	require.Len(t, boardDTO, 10)
	for y := 0; y < 10; y++ {
		require.Len(t, boardDTO[y], 10)
		for x := 0; x < 10; x++ {
			if x == 5 && y == 3 {
				assert.Equal(t, "Marshal", boardDTO[y][x].Type)
				assert.Equal(t, 0, boardDTO[y][x].OwnerID)
				assert.Equal(t, 5, boardDTO[y][x].Position.X)
				assert.Equal(t, 3, boardDTO[y][x].Position.Y)
			} else {
				assert.Equal(t, "", boardDTO[y][x].Type)
				assert.Equal(t, -1, boardDTO[y][x].OwnerID)
			}
		}
	}
}
