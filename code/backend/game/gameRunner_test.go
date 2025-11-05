package game_test

import (
	"digital-innovation/stratego/ai/fafo"
	"digital-innovation/stratego/engine"
	"digital-innovation/stratego/game"
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

		controller1 := fafo.NewFafoAI(&player1)
		controller2 := fafo.NewFafoAI(&player2)

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
	controller2 := fafo.NewFafoAI(&player2)

	gameInstance := game.QuickStart(controller1, controller2)
	runner := game.NewGameRunner(gameInstance, 0, 1000)
	runner.DebugSetWaitingForInput(true)
	game := runner.GetGame()

	time.Sleep(1 * time.Second)
	move := engine.NewMove(engine.NewPosition(2, 6), engine.NewPosition(2, 5), &player1)
	err := runner.SubmitHumanMove(move)
	if err == nil {
		t.Errorf("Expected error when submitting invalid move")
	}

	move = engine.NewMove(engine.NewPosition(0, 6), engine.NewPosition(0, 5), &player1)
	err = runner.SubmitHumanMove(move)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if runner.IsWaitingForInput() {
		t.Errorf("Expected not waiting for input")
	}

	if game.CurrentPlayer != &player2 {
		t.Errorf("Expected current player to be %s, got: %s",
			player2.GetName(), game.CurrentPlayer.GetName())
	}
}
