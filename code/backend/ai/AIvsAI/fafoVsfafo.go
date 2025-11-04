package aivsai

import (
	"digital-innovation/stratego/ai/fafo"
	"digital-innovation/stratego/engine"
	"digital-innovation/stratego/game"
	"digital-innovation/stratego/models"
	"fmt"
)

func runAIvsAI(matches int) models.GameSummary {
	aliceWins := 0
	bobWins := 0
	draws := 0

	// Track win causes
	flagCaptures := 0
	noMovesWins := 0
	maxTurnsWins := 0
	totalRounds := 0

	player1Name := ""
	player2Name := ""

	// Run 100 games - alternate who goes first for fairness
	for i := 0; i < matches; i++ {
		// Create FRESH players and controllers for EACH game
		playerAlice := engine.NewPlayer(0, "Alice AI", "red")
		playerBob := engine.NewPlayer(1, "Bob AI", "blue")

		player1Name = playerAlice.GetName()
		player2Name = playerBob.GetName()

		controllerAlice := fafo.NewFafoAI(&playerAlice)
		controllerBob := fafo.NewFafoAI(&playerBob)

		// Alternate who goes first (even games: Alice first, odd games: Bob first)
		var g *game.Game
		if i%2 == 0 {
			// Alice goes first
			g = game.QuickStart(controllerAlice, controllerBob)
			fmt.Printf("Game %3d (Alice 1st): ", i+1)
		} else {
			// Bob goes first (swap player order)
			g = game.QuickStart(controllerBob, controllerAlice)
			fmt.Printf("Game %3d (Bob 1st):   ", i+1)
		}

		runner := game.NewGameRunner(g, 0, 1000)
		winner := runner.RunToCompletion()

		// Get game statistics
		rounds := g.GetRound()
		winCause := g.GetWinCause()
		totalRounds += rounds

		// Track win causes
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

	// Print tournament summary
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
