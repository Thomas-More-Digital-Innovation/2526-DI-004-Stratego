package engine_test

import (
	"digital-innovation/stratego/engine"
	"testing"
)

func TestGetFrom(t *testing.T) {
	player := engine.NewPlayer(1, "player1", "avatar1")
	move := engine.NewMove(engine.NewPosition(1, 2), engine.NewPosition(3, 4), &player)

	from := move.GetFrom()

	if from.X != 1 || from.Y != 2 {
		t.Errorf("Expected from coordinates to be (1,2), got (%d,%d)", from.X, from.Y)
	}
}

func TestGetTo(t *testing.T) {
	player := engine.NewPlayer(1, "player1", "avatar1")
	move := engine.NewMove(engine.NewPosition(1, 2), engine.NewPosition(3, 4), &player)

	to := move.GetTo()

	if to.X != 3 || to.Y != 4 {
		t.Errorf("Expected to coordinates to be (3,4), got (%d,%d)", to.X, to.Y)
	}
}

func TestGetPlayer(t *testing.T) {
	player := engine.NewPlayer(1, "player1", "avatar1")
	move := engine.NewMove(engine.NewPosition(1, 2), engine.NewPosition(3, 4), &player)

	movePlayer := move.GetPlayer()

	if movePlayer.GetID() != player.GetID() {
		t.Errorf("Expected move player ID to be %d, got %d", player.GetID(), movePlayer.GetID())
	}
}
