package game_test

import (
	AIhandler "digital-innovation/gostrategy/internal/ai/handler"
	"digital-innovation/gostrategy/internal/game"
	"digital-innovation/gostrategy/internal/game/models"
	"digital-innovation/gostrategy/internal/testutils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunToCompletion(t *testing.T) {
	player1Wins := 0
	player2Wins := 0
	draws := 0

	flagCaptures := 0
	noMovesWins := 0
	maxTurnsWins := 0
	totalRounds := 0

	numGames := 10

	const (
		p1 = "player 1"
		p2 = "player 2"
	)

	for i := 0; i < numGames; i++ {
		player1 := game.NewPlayer(0, p1, "red")
		player2 := game.NewPlayer(1, p2, "blue")

		controller1, err := AIhandler.CreateAI(models.Fafo, &player1)
		require.NoError(t, err)

		controller2, err := AIhandler.CreateAI(models.Fafo, &player2)
		require.NoError(t, err)

		var g *game.Game
		if i%2 == 0 {
			g = game.QuickStart(controller1, controller2)
		} else {
			g = game.QuickStart(controller2, controller1)
		}

		runner := game.NewRunner(g, 0, 1000)
		winner := runner.RunToCompletion()
		rounds := g.GetRound()

		winCause := g.GetWinCause()
		totalRounds += rounds

		switch winCause {
		case game.WinCauseFlagCaptured:
			flagCaptures++
		case game.WinCauseNoMovablePieces:
			noMovesWins++
		case game.WinCauseMaxTurns:
			maxTurnsWins++
		}

		switch {
		case winner == nil:
			draws++
		case winner.GetName() == p1:
			player1Wins++
		default:
			player2Wins++
		}
	}

	avgRounds := float64(totalRounds) / float64(numGames)
	assert.True(t, avgRounds <= 900, "Average rounds per game too high: %.2f", avgRounds)
	assert.True(t, draws <= numGames/2, "Too many draws: %d out of %d games", draws, numGames)
	assert.True(t, player1Wins <= 75 && player2Wins <= 75, "One player is winning too often: Player1=%d, Player2=%d", player1Wins, player2Wins)
	assert.True(t, avgRounds >= 10, "Average rounds per game too low: %.2f", avgRounds)
}

func TestSubmitHumanMove(t *testing.T) {
	player1 := game.NewPlayer(1, "fakeHuman", "red")
	player2 := game.NewPlayer(2, "realAI", "blue")

	controller1 := game.NewHumanPlayerController(&player1)
	controller2, err := AIhandler.CreateAI(models.Fafo, &player2)
	require.NoError(t, err)

	gameInstance := game.QuickStart(controller1, controller2)
	runner := game.NewRunner(gameInstance, 0, 1000)
	runner.DebugSetWaitingForInput(true)
	gameObj := runner.GetGame()

	time.Sleep(1 * time.Second)
	// Invalid move - bomb can't move
	move := game.NewMove(game.NewPosition(2, 6), game.NewPosition(2, 5), &player1)
	err = runner.SubmitHumanMove(move)
	assert.Error(t, err)

	// Valid move
	move = game.NewMove(game.NewPosition(0, 6), game.NewPosition(0, 5), &player1)
	err = runner.SubmitHumanMove(move)
	assert.NoError(t, err)

	assert.False(t, runner.IsWaitingForInput())
	assert.Equal(t, &player2, gameObj.CurrentPlayer)
}

func TestRunnerIsWaitingForInput(t *testing.T) {
	player1 := game.NewPlayer(0, "Human", "red")
	player2 := game.NewPlayer(1, "AI", "blue")

	controller1 := game.NewHumanPlayerController(&player1)
	controller2, err := AIhandler.CreateAI(models.Fafo, &player2)
	require.NoError(t, err)

	g := game.QuickStart(controller1, controller2)
	runner := game.NewRunner(g, 0, 1000)

	assert.False(t, runner.IsWaitingForInput(), "Expected not waiting for input before game starts")

	runner.DebugSetWaitingForInput(true)
	assert.True(t, runner.IsWaitingForInput(), "Expected waiting for input after DebugSetWaitingForInput(true)")

	runner.DebugSetWaitingForInput(false)
	assert.False(t, runner.IsWaitingForInput(), "Expected not waiting for input after DebugSetWaitingForInput(false)")
}

func TestRunnerGetGame(t *testing.T) {
	g, _, _ := testutils.SetupTestGame()
	runner := game.NewRunner(g, 0, 1000)

	retrieved := runner.GetGame()
	assert.Equal(t, g, retrieved)
}

func TestSubmitHumanMoveWrongPlayer(t *testing.T) {
	player1 := game.NewPlayer(0, "Human1", "red")
	player2 := game.NewPlayer(1, "Human2", "blue")

	controller1 := game.NewHumanPlayerController(&player1)
	controller2 := game.NewHumanPlayerController(&player2)

	g := game.QuickStart(controller1, controller2)
	runner := game.NewRunner(g, 0, 1000)
	runner.DebugSetWaitingForInput(true)

	// Try to submit move for player2 when it's player1's turn
	move := game.NewMove(game.NewPosition(0, 3), game.NewPosition(0, 4), &player2)
	err := runner.SubmitHumanMove(move)
	assert.Error(t, err)
}

func TestRunnerWithDelay(t *testing.T) {
	player1 := game.NewPlayer(0, "AI1", "red")
	player2 := game.NewPlayer(1, "AI2", "blue")

	controller1, err := AIhandler.CreateAI(models.Fafo, &player1)
	require.NoError(t, err)

	controller2, err := AIhandler.CreateAI(models.Fafo, &player2)
	require.NoError(t, err)

	g := game.QuickStart(controller1, controller2)

	// Create runner with short delay
	runner := game.NewRunner(g, 5*time.Millisecond, 10)

	start := time.Now()
	runner.RunToCompletion()
	elapsed := time.Since(start)

	assert.True(t, elapsed >= 10*time.Millisecond, "Expected game to take at least 10ms with delays, took: %v", elapsed)
}

func TestRunner_Immobilization(t *testing.T) {
	player1 := game.NewPlayer(0, "AI1", "red")
	player2 := game.NewPlayer(1, "AI2", "blue")

	controller1, err := AIhandler.CreateAI(models.Fafo, &player1)
	require.NoError(t, err)

	controller2, err := AIhandler.CreateAI(models.Fafo, &player2)
	require.NoError(t, err)

	g := game.QuickStart(controller1, controller2)
	runner := game.NewRunner(g, 0, 1000)

	// Eliminate all movable pieces for player 1
	alive := append([]*game.Piece{}, player1.GetAlivePieces()...)
	for _, piece := range alive {
		if piece.CanMove() {
			pos, _ := player1.GetPiecePosition(piece)
			piece.Eliminate()
			// We should also remove them from the board to be realistic
			g.Board.SetPieceAt(pos, nil)
		}
	}

	// Try to execute turn for player 1 (who has no movable pieces)
	runner.ExecuteTurn()

	require.True(t, g.IsGameOver())
	assert.Equal(t, game.WinCauseNoMovablePieces, g.GetWinCause())
	assert.Equal(t, &player2, g.GetWinner())
}

func TestRunnerAbort(t *testing.T) {
	player1 := game.NewPlayer(0, "AI1", "red")
	player2 := game.NewPlayer(1, "AI2", "blue")
	controller1, err := AIhandler.CreateAI(models.Fafo, &player1)
	require.NoError(t, err)

	controller2, err := AIhandler.CreateAI(models.Fafo, &player2)
	require.NoError(t, err)

	g := game.QuickStart(controller1, controller2)

	runner := game.NewRunner(g, 100*time.Millisecond, 1000)

	go func() {
		time.Sleep(50 * time.Millisecond)
		runner.Stop()
	}()

	winner := runner.RunToCompletion()
	assert.Nil(t, winner, "Winner should be nil when aborted")
}
