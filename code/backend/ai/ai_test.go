package ai

import (
	"digital-innovation/stratego/engine"
	"testing"
)

func TestNewBaseAInoMemory(t *testing.T) {
	player := engine.NewPlayer(0, "player", "red")

	ai := NewBaseAI(&player, false)

	if ai.player == nil {
		t.Errorf("Expected player to be set in BaseAI")
	}

	if ai.memory != nil {
		t.Errorf("Expected memory to be nil in BaseAI")
	}
}

func TestNewBaseAIWithMemory(t *testing.T) {
	player := engine.NewPlayer(0, "player", "red")
	ai := NewBaseAI(&player, true)

	if ai.player == nil {
		t.Errorf("Expected player to be set in BaseAI")
	}

	if ai.memory == nil {
		t.Errorf("Expected memory to be initialized in BaseAI")
	}
}

func TestGetPlayer(t *testing.T) {
	player := engine.NewPlayer(0, "player", "red")
	ai := NewBaseAI(&player, false)

	if ai.GetPlayer() != &player {
		t.Errorf("Expected player to be returned")
	}
}

func TestGetControllerType(t *testing.T) {
	player := engine.NewPlayer(0, "player", "red")
	ai := NewBaseAI(&player, false)

	controllerType := ai.GetControllerType()

	if controllerType != engine.AIController {
		t.Errorf("Expected controller type to be AI")
	}
}
