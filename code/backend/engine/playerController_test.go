package engine

import (
	"testing"
)

func TestHumanPlayerController(t *testing.T) {
	p := NewPlayer(1, "Player 1", "")
	player := &p
	controller := NewHumanPlayerController(player)

	if controller.GetPlayer() != player {
		t.Errorf("GetPlayer() returned wrong player")
	}

	if controller.GetControllerType() != HumanController {
		t.Errorf("GetControllerType() returned wrong type")
	}

	// MakeMove should return empty move
	move := controller.MakeMove(nil)
	if !move.IsEmpty() {
		t.Errorf("MakeMove() should return empty move for human")
	}

	// Test pending move
	if controller.HasPendingMove() {
		t.Errorf("New controller should not have pending move")
	}

	testMove := NewMove(NewPosition(0, 0), NewPosition(1, 0), player)
	controller.SetPendingMove(testMove)

	if !controller.HasPendingMove() {
		t.Errorf("Controller should have pending move after SetPendingMove")
	}

	gotMove := controller.GetPendingMove()
	if gotMove == nil || *gotMove != testMove {
		t.Errorf("GetPendingMove() returned wrong move")
	}

	if controller.HasPendingMove() {
		t.Errorf("Controller should not have pending move after GetPendingMove")
	}
}
