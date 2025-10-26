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

	// Run 100 games
	for i := 0; i < 100; i++ {
		// Create FRESH players and controllers for EACH game
		player1 := engine.NewPlayer(0, "Alice AI", "red")
		player2 := engine.NewPlayer(1, "Bob AI", "blue")

		controller1 := fafo.NewFafoAI(&player1)
		controller2 := fafo.NewFafoAI(&player2)

		// Setup and start game
		g := game.QuickStart(controller1, controller2)
		runner := game.NewGameRunner(g, 1*time.Nanosecond, 1000)

		winner := runner.RunToCompletion()

		if winner == nil {
			fmt.Printf("Game %3d: Draw (max turns)\n", i+1)
			draws++
		} else if winner.GetID() == 0 {
			fmt.Printf("Game %3d: Alice wins\n", i+1)
			aliceWins++
		} else {
			fmt.Printf("Game %3d: Bob wins\n", i+1)
			bobWins++
		}
	}

	// Print tournament summary
	fmt.Println("\n========================================")
	fmt.Println("📊 Tournament Results (100 games)")
	fmt.Println("========================================")
	fmt.Printf("Alice wins: %3d (%.1f%%)\n", aliceWins, float64(aliceWins))
	fmt.Printf("Bob wins:   %3d (%.1f%%)\n", bobWins, float64(bobWins))
	fmt.Printf("Draws:      %3d (%.1f%%)\n", draws, float64(draws))
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
