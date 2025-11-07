package aivsai

import (
	"digital-innovation/stratego/ai/fafo"
	"digital-innovation/stratego/engine"
	"digital-innovation/stratego/game"
	"digital-innovation/stratego/models"
	"fmt"
)

func runAIvsAI(matches int, logging bool) models.GameSummary {
	aliceWins := 0
	bobWins := 0
	draws := 0

	flagCaptures := 0
	noMovesWins := 0
	maxTurnsWins := 0
	totalRounds := 0

	player1Name := ""
	player2Name := ""

	for i := 0; i < matches; i++ {
		// Create FRESH players and controllers for EACH game, use same ID & name
		playerAlice := engine.NewPlayer(0, "Alice AI", "red")
		playerBob := engine.NewPlayer(1, "Bob AI", "blue")

		player1Name = playerAlice.GetName()
		player2Name = playerBob.GetName()

		controllerAlice := fafo.NewFafoAI(&playerAlice)
		controllerBob := fafo.NewFafoAI(&playerBob)

		// Alternate who goes first
		// Without this, player 1 wins more often than the other
		var g *game.Game
		if i%2 == 0 {
			g = game.QuickStart(controllerAlice, controllerBob)
			fmt.Printf("Game %3d (Alice starts): ", i+1)
		} else {
			g = game.QuickStart(controllerBob, controllerAlice)
			fmt.Printf("Game %3d (Bob starts):   ", i+1)
		}

		runner := game.NewGameRunner(g, 0, 1000)
		winner := runner.RunToCompletion(logging)

		rounds := g.GetRound()
		winCause := g.GetWinCause()
		totalRounds += rounds

		switch winCause {
		case game.WinCauseFlagCaptured:
			flagCaptures++
		case game.WinCauseNoMovablePieces:
			noMovesWins++
		default:
			maxTurnsWins++
		}

		switch {
		case winner == nil:
			fmt.Printf("Draw after %d rounds\n", rounds)
			draws++
		case winner.GetName() == "Alice AI":
			fmt.Printf("Alice wins - %s (%d rounds)\n", winCause, rounds)
			aliceWins++
		default:
			fmt.Printf("Bob wins - %s (%d rounds)\n", winCause, rounds)
			bobWins++
		}

	}

	avgRounds := float64(totalRounds) / float64(matches)

	gameSummary := models.GameSummary{
		Player1Name:          player1Name,
		Player2Name:          player2Name,
		Player1Wins:          aliceWins,
		Player2Wins:          bobWins,
		Draws:                draws,
		TotalRounds:          totalRounds,
		AverageRounds:        avgRounds,
		Matches:              matches,
		WinCauseFlagCaptured: float64(flagCaptures),
		WinCauseNoMovesWins:  float64(noMovesWins),
		WinCauseMaxTurns:     float64(maxTurnsWins),
	}

	return gameSummary
}
