package main

import (
	"digital-innovation/stratego/ai/fafo"
	"digital-innovation/stratego/api"
	"digital-innovation/stratego/engine"
	"digital-innovation/stratego/game"
	"flag"
	"fmt"
	"log"
)

func main() {
	// Command line flags
	serverMode := flag.Bool("server", false, "Run in WebSocket server mode")
	addr := flag.String("addr", ":8080", "Server address")
	flag.Parse()

	fmt.Println("=== Stratego Backend Running ===")

	if *serverMode {
		// Run WebSocket server
		runServer(*addr)
	} else {
		// Run AI vs AI matches (original behavior)
		runAIvsAI()
	}
}

// runServer starts the WebSocket server
func runServer(addr string) {
	fmt.Printf("Starting Stratego WebSocket Server on %s\n", addr)
	fmt.Println("Endpoints:")
	fmt.Println("  POST /api/games - Create a new game")
	fmt.Println("  GET  /api/games - List all games")
	fmt.Println("  WS   /ws/game/{gameID}?player={0|1|spectator} - Connect to game")

	server := api.NewGameServer()
	if err := server.StartServer(addr); err != nil {
		log.Fatalf("Server error: %v", err)
	}
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
