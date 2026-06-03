package game_test

import (
	AIhandler "digital-innovation/gostrategy/pkg/ai/handler"
	"digital-innovation/gostrategy/pkg/game"
	"digital-innovation/gostrategy/pkg/game/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStepWhilePaused(t *testing.T) {
	player1 := game.NewPlayer(0, "AI1", "red")
	player2 := game.NewPlayer(1, "AI2", "blue")

	controller1, err := AIhandler.CreateAI(models.Fafo, &player1)
	require.NoError(t, err)

	controller2, err := AIhandler.CreateAI(models.Fafo, &player2)
	require.NoError(t, err)

	g := game.QuickStart(controller1, controller2)

	// Create runner with a delay to trigger the second pause check
	runner := game.NewRunner(g, 10*time.Millisecond, 1000)

	// Pause the runner
	runner.Pause()
	assert.True(t, runner.IsPaused(), "Runner should be paused")

	initialMoveHistory := len(g.MoveHistory)

	// Attempt to step while paused
	success := runner.Step()
	assert.True(t, success, "Step should have succeeded even when paused")
	assert.Len(t, g.MoveHistory, initialMoveHistory+1)

	// Verify that the game is still paused
	assert.True(t, runner.IsPaused(), "Runner should still be paused after step")
}
