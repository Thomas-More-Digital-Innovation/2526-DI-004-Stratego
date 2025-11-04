package main

import (
	aivsai "digital-innovation/stratego/ai/AIvsAI"
	"digital-innovation/stratego/api"
	"flag"
	"fmt"
	"log"
)

func main() {
	// Command line flags
	serverMode := flag.Bool("server", false, "Run in WebSocket server mode")
	addr := flag.String("addr", ":8080", "Server address")
	aiVsaiMode := flag.String("aivsai", "fafo:fafo", "Run AI vs AI matches instead of server")
	matches := flag.Int("matches", 100, "Number of AI vs AI matches to run")
	format := flag.String("format", "none", "The format used to print the results of an AI vs AI competition, either none or md")

	flag.Parse()

	fmt.Println("=== Stratego Backend Running ===")

	if *serverMode {
		// Run WebSocket server
		runServer(*addr)
	} else {
		var ai1, ai2 string
		if aiVsaiMode != nil {
			// Not implemented yet
		} else {
			ai1 = "fafo"
			ai2 = "fafo"
		}
		aivsai.RunAIvsAI(ai1, ai2, *matches, *format)
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
