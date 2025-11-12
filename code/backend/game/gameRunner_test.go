package game_test

import (
	AIhandler "digital-innovation/stratego/ai/handler"
	"digital-innovation/stratego/engine"
	"digital-innovation/stratego/game"
	"digital-innovation/stratego/models"
	"testing"
	"time"
)

func TestRunToCompletion(t *testing.T) {
	// test 100 games and see if there is predictable outcome
	player1Wins := 0
	player2Wins := 0
	draws := 0

	// Track win causes
	flagCaptures := 0
	noMovesWins := 0
	maxTurnsWins := 0
	totalRounds := 0

	numGames := 100

	const (
		p1 = "player 1"
		p2 = "player 2"
	)

	for i := 0; i < numGames; i++ {
		player1 := engine.NewPlayer(0, p1, "red")
		player2 := engine.NewPlayer(1, p2, "blue")

		controller1 := AIhandler.CreateAI(models.Fafo, &player1)
		controller2 := AIhandler.CreateAI(models.Fafo, &player2)

		var g *game.Game
		if i%2 == 0 {
			g = game.QuickStart(controller1, controller2)
		} else {
			g = game.QuickStart(controller2, controller1)
		}

		runner := game.NewGameRunner(g, 0, 1000)
		winner := runner.RunToCompletion(false) // we don't want cluttered logging in pipeline
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
	if avgRounds > 900 {
		t.Errorf("Average rounds per game too high: %.2f", avgRounds)
	}

	if draws > numGames/5 {
		t.Errorf("Too many draws: %d out of %d games", draws, numGames)

	}

	if player1Wins > 75 || player2Wins > 75 {
		t.Errorf("One player is winning too often: Player1=%d, Player2=%d", player1Wins, player2Wins)
	}

	if avgRounds < 10 {
		t.Errorf("Average rounds per game too low: %.2f", avgRounds)
	}

}

func TestSubmitHumanMove(t *testing.T) {
	player1 := engine.NewPlayer(1, "fakeHuman", "red")
	player2 := engine.NewPlayer(2, "realAI", "blue")

	controller1 := engine.NewHumanPlayerController(&player1)
	controller2 := AIhandler.CreateAI(models.Fafo, &player2)

	gameInstance := game.QuickStart(controller1, controller2)
	runner := game.NewGameRunner(gameInstance, 0, 1000)
	runner.DebugSetWaitingForInput(true)
	gameObj := runner.GetGame()

	time.Sleep(1 * time.Second)
	// Invalid move - bomb can't move
	move := engine.NewMove(engine.NewPosition(2, 6), engine.NewPosition(2, 5), &player1)
	err := runner.SubmitHumanMove(move)
	if err == nil {
		t.Errorf("Expected error when submitting invalid move")
	}

	// Valid move
	move = engine.NewMove(engine.NewPosition(0, 6), engine.NewPosition(0, 5), &player1)
	err = runner.SubmitHumanMove(move)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if runner.IsWaitingForInput() {
		t.Errorf("Expected not waiting for input")
	}

	if gameObj.CurrentPlayer != &player2 {
		t.Errorf("Expected current player to be %s, got: %s",
			player2.GetName(), gameObj.CurrentPlayer.GetName())
	}
}

func TestGameRunnerIsWaitingForInput(t *testing.T) {
	player1 := engine.NewPlayer(0, "Human", "red")
	player2 := engine.NewPlayer(1, "AI", "blue")

	controller1 := engine.NewHumanPlayerController(&player1)
	controller2 := AIhandler.CreateAI(models.Fafo, &player2)

	g := game.QuickStart(controller1, controller2)
	runner := game.NewGameRunner(g, 0, 1000)

	if runner.IsWaitingForInput() {
		t.Error("Expected not waiting for input before game starts")
	}

	runner.DebugSetWaitingForInput(true)
	if !runner.IsWaitingForInput() {
		t.Error("Expected waiting for input after DebugSetWaitingForInput(true)")
	}

	runner.DebugSetWaitingForInput(false)
	if runner.IsWaitingForInput() {
		t.Error("Expected not waiting for input after DebugSetWaitingForInput(false)")
	}
}

func TestGameRunnerGetGame(t *testing.T) {
	player1 := engine.NewPlayer(0, "P1", "red")
	player2 := engine.NewPlayer(1, "P2", "blue")

	controller1 := engine.NewHumanPlayerController(&player1)
	controller2 := engine.NewHumanPlayerController(&player2)

	g := game.NewGame(controller1, controller2)
	runner := game.NewGameRunner(g, 0, 1000)

	retrieved := runner.GetGame()
	if retrieved != g {
		t.Error("Expected GetGame to return the same game instance")
	}
}

func TestSubmitHumanMoveWrongPlayer(t *testing.T) {
	player1 := engine.NewPlayer(0, "Human1", "red")
	player2 := engine.NewPlayer(1, "Human2", "blue")

	controller1 := engine.NewHumanPlayerController(&player1)
	controller2 := engine.NewHumanPlayerController(&player2)

	g := game.QuickStart(controller1, controller2)
	runner := game.NewGameRunner(g, 0, 1000)
	runner.DebugSetWaitingForInput(true)

	// Try to submit move for player2 when it's player1's turn
	move := engine.NewMove(engine.NewPosition(0, 3), engine.NewPosition(0, 4), &player2)
	err := runner.SubmitHumanMove(move)

	if err == nil {
		t.Error("Expected error when submitting move for wrong player")
	}
}

func TestGameRunnerWithDelay(t *testing.T) {
	player1 := engine.NewPlayer(0, "AI1", "red")
	player2 := engine.NewPlayer(1, "AI2", "blue")

	controller1 := AIhandler.CreateAI(models.Fafo, &player1)
	controller2 := AIhandler.CreateAI(models.Fafo, &player2)

	g := game.QuickStart(controller1, controller2)

	// Create runner with short delay
	runner := game.NewGameRunner(g, 5*time.Millisecond, 10)

	start := time.Now()
	runner.RunToCompletion(false)
	elapsed := time.Since(start)

	// With 5ms delay per turn and up to 10 turns, should take at least 50ms
	// But we set max turns to 10, so it should be relatively quick
	if elapsed < 10*time.Millisecond {
		t.Errorf("Expected game to take at least 10ms with delays, took: %v", elapsed)
	}
}
