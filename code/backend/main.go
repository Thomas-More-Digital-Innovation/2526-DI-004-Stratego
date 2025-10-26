package main

import (
	"digital-innovation/stratego/ai/fafo"
	"digital-innovation/stratego/engine"
	"digital-innovation/stratego/game"
	"fmt"
	"time"
)

func main() {
	fmt.Println("=== Stratego Backend Running ===\n")

	// Example 1: AI vs AI match
	runAIvsAI()
}

// runAIvsAI demonstrates an AI vs AI game that runs to completion
func runAIvsAI() {
	aliceWins := 0
	bobWins := 0
	draws := 0
	
	// Track win causes
	flagCaptures := 0
	noMovesWins := 0
	maxTurnsWins := 0
	totalRounds := 0

	// Run 100 games - alternate who goes first for fairness
	numGames := 100
	for i := 0; i < numGames; i++ {
		// Create FRESH players and controllers for EACH game
		playerAlice := engine.NewPlayer(0, "Alice AI", "red")
		playerBob := engine.NewPlayer(1, "Bob AI", "blue")

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

		runner := game.NewGameRunner(g, 1*time.Nanosecond, 1000)
		winner := runner.RunToCompletion()

		// Get game statistics
		rounds := g.GetRound()
		winCause := g.GetWinCause()
		totalRounds += rounds
		
		// Track win causes
		if winCause == game.WinCauseFlagCaptured {
			flagCaptures++
		} else if winCause == game.WinCauseNoMovablePieces {
			noMovesWins++
		} else if winCause == game.WinCauseMaxTurns {
			maxTurnsWins++
		}
		
		if winner == nil {
			fmt.Printf("Draw after %d rounds\n", rounds)
			draws++
		} else if winner.GetName() == "Alice AI" {
			fmt.Printf("Alice wins - %s (%d rounds)\n", winCause, rounds)
			aliceWins++
		} else {
			fmt.Printf("Bob wins - %s (%d rounds)\n", winCause, rounds)
			bobWins++
		}
	}

	// Print tournament summary
	avgRounds := float64(totalRounds) / float64(numGames)
	fmt.Println("\n========================================")
	fmt.Printf("📊 Tournament Results (%d games)\n", numGames)
	fmt.Println("========================================")
	fmt.Printf("Alice wins: %3d (%.1f%%)\n", aliceWins, float64(aliceWins*100)/float64(numGames))
	fmt.Printf("Bob wins:   %3d (%.1f%%)\n", bobWins, float64(bobWins*100)/float64(numGames))
	fmt.Printf("Draws:      %3d (%.1f%%)\n", draws, float64(draws*100)/float64(numGames))
	fmt.Println("========================================")
	fmt.Println("Win Causes:")
	fmt.Printf("  Flag captured:     %3d (%.1f%%)\n", flagCaptures, float64(flagCaptures*100)/float64(numGames))
	fmt.Printf("  No movable pieces: %3d (%.1f%%)\n", noMovesWins, float64(noMovesWins*100)/float64(numGames))
	fmt.Printf("  Max turns:         %3d (%.1f%%)\n", maxTurnsWins, float64(maxTurnsWins*100)/float64(numGames))
	fmt.Println("========================================")
	fmt.Printf("Average game length: %.1f rounds\n", avgRounds)
	fmt.Println("========================================")
}

// runAIvsHuman demonstrates how an AI vs Human game would work
func runAIvsHuman() {
	// Create players
	player1 := engine.NewPlayer(0, "Alice AI", "red")
	player2 := engine.NewPlayer(1, "Human Player", "blue")

	// Player 1 is AI, Player 2 is Human
	controller1 := fafo.NewFafoAI(&player1)
	controller2 := engine.NewHumanPlayerController(&player2)

	// Setup game
	g := game.QuickStart(controller1, controller2)
	runner := game.NewGameRunner(g, 0, 1000)

	fmt.Printf("Game started: %s (AI) vs %s (Human)\n", player1.GetName(), player2.GetName())

	// Simulate a few turns
	for turn := 0; turn < 5 && !g.IsGameOver(); turn++ {
		currentPlayer := g.CurrentPlayer
		fmt.Printf("\n--- Turn %d: %s's turn ---\n", turn+1, currentPlayer.GetName())

		if runner.IsWaitingForInput() {
			// In a real scenario, this would wait for HTTP request with move
			fmt.Println("⏳ Waiting for human input...")

			// Simulate human move (in real app, this comes from frontend)
			pieces := currentPlayer.GetAlivePieces()
			if len(pieces) > 0 {
				// Find a movable piece
				for _, piece := range pieces {
					if !piece.CanMove() {
						continue
					}
					pos, exists := currentPlayer.GetPiecePosition(piece)
					if !exists {
						continue
					}

					moves, err := g.Board.ListMoves(pos)
					if err == nil && len(moves) > 0 {
						// Submit the first valid move
						move := moves[0]
						fmt.Printf("📥 Human submits move: %v -> %v\n", move.GetFrom(), move.GetTo())
						err := runner.SubmitHumanMove(move)
						if err != nil {
							fmt.Printf("❌ Error: %v\n", err)
						}
						break
					}
				}
			}
		} else {
			// AI turn - executes automatically
			if runner.ExecuteTurn() {
				fmt.Println("🤖 AI made its move")
			}
		}
	}

	fmt.Println("\n✅ Demo complete (game may still be ongoing)")
	fmt.Printf("Current Round: %d\n", g.GetRound())
	fmt.Println("\nCurrent Board:\n" + g.Board.String())
}
