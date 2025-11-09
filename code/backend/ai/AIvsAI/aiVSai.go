package aivsai

import (
	"digital-innovation/stratego/ai"
	"digital-innovation/stratego/ai/fafo"
	"digital-innovation/stratego/ai/fato"
	"digital-innovation/stratego/engine"
	"digital-innovation/stratego/game"
	"digital-innovation/stratego/models"
	"fmt"
)

func runAIvsAI(ai1, ai2 string, matches int, logging bool) models.GameSummary {
	draws := 0

	flagCaptures := 0
	noMovesWins := 0
	maxTurnsWins := 0
	totalRounds := 0
	leastRounds := 1000 // we start with our max rounds possible

	player1Name := "Alice AI - " + ai1
	player2Name := "Bob AI - " + ai2

	player1Data := models.AiTournamentData{Name: player1Name}
	player2Data := models.AiTournamentData{Name: player2Name}

	for i := range matches {
		// Create FRESH players and controllers for EACH game, use same ID & name
		playerAlice := engine.NewPlayer(0, player1Name, "red")
		playerBob := engine.NewPlayer(1, player2Name, "blue")

		controllerAlice := createAI(ai1, &playerAlice)
		controllerBob := createAI(ai2, &playerBob)

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
		totalRounds += rounds

		if rounds < leastRounds {
			leastRounds = rounds
		}

		if winner != nil {
			var winnerData *models.AiTournamentData
			winCause := g.GetWinCause()
			if winner.GetName() == player1Name {
				winnerData = &player1Data
				fmt.Printf("%v wins - %s (%d rounds)\n", player1Name, winCause, rounds)

			} else {
				winnerData = &player2Data
				fmt.Printf("%v wins - %s (%d rounds)\n", player2Name, winCause, rounds)
			}

			switch winCause {
			case game.WinCauseFlagCaptured:
				winnerData.WinCauseFlagCaptured++
				flagCaptures++
			case game.WinCauseNoMovablePieces:
				winnerData.WinCauseNoMovesWin++
				noMovesWins++
			default:
				winnerData.WinCauseMaxTurns++
				maxTurnsWins++
			}

			winnerData.Wins++

		} else {
			fmt.Printf("Draw after %d rounds\n", rounds)
			draws++
		}

	}

	avgRounds := float64(totalRounds) / float64(matches)

	gameSummary := models.GameSummary{
		Player1data:          player1Data,
		Player2data:          player2Data,
		Draws:                draws,
		TotalRounds:          totalRounds,
		AverageRounds:        avgRounds,
		LeastRounds:          leastRounds,
		Matches:              matches,
		WinCauseFlagCaptured: flagCaptures,
		WinCauseNoMovesWins:  noMovesWins,
		WinCauseMaxTurns:     maxTurnsWins,
	}

	return gameSummary
}

func createAI(ai string, player *engine.Player) ai.AI {
	switch ai {
	case models.Fafo:
		return fafo.NewFafoAI(player)
	case models.Fato:
		return fato.NewFatoAI(player)
	default:
		panic("I don't know that AI! " + ai)
	}

}
