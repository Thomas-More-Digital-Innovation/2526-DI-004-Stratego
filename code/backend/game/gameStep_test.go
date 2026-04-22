package game_test

import (
	AIhandler "digital-innovation/stratego/ai/handler"
	"digital-innovation/stratego/engine"
	"digital-innovation/stratego/game"
	"digital-innovation/stratego/models"
	"testing"
	"time"
)

func TestStepWhilePaused(t *testing.T) {
	player1 := engine.NewPlayer(0, "AI1", "red")
	player2 := engine.NewPlayer(1, "AI2", "blue")

	controller1 := AIhandler.CreateAI(models.Fafo, &player1)
	controller2 := AIhandler.CreateAI(models.Fafo, &player2)

	g := game.QuickStart(controller1, controller2)

	// Create runner with a delay to trigger the second pause check
	runner := game.NewGameRunner(g, 10*time.Millisecond, 1000)

	// Pause the runner
	runner.Pause()

	if !runner.IsPaused() {
		t.Fatal("Runner should be paused")
	}

	initialMoveHistory := len(g.MoveHistory)

	// Attempt to step while paused
	success := runner.Step(true)

	if !success {
		t.Error("Step should have succeeded even when paused")
	}

	if len(g.MoveHistory) != initialMoveHistory+1 {
		t.Errorf("Expected move history length to increase by 1, got %d -> %d", initialMoveHistory, len(g.MoveHistory))
	}

	// Verify that the game is still paused
	if !runner.IsPaused() {
		t.Error("Runner should still be paused after step")
	}
}
