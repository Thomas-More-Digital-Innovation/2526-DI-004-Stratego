package main

import (
	aivsai "digital-innovation/stratego/ai/AIvsAI"
	"digital-innovation/stratego/api"
	"digital-innovation/stratego/auth"
	"digital-innovation/stratego/db"
	"digital-innovation/stratego/models"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"
)

func main() {
	serverMode := flag.Bool("server", false, "Run in WebSocket server mode")
	addr := flag.String("addr", ":8080", "Server address")
	aiTypes := flag.String("ai", "fafo:fafo", "Run AI vs AI matches instead of server")
	matches := flag.Int("matches", 100, "Number of AI vs AI matches to run")
	format := flag.String("format", "none", "The format used to print the results of an AI vs AI competition, either none or md")
	logging := flag.Bool("logging", true, "Show logs in stdout")

	flag.Parse()

	fmt.Println("=== Stratego Backend Running ===")

	if *serverMode {
		if err := db.InitDB(); err != nil {
			log.Fatalf("Failed to initialize database: %v", err)
		}
		defer db.CloseDB()

		auth.Store.StartCleanupRoutine()

		runServer(*addr) // websocket server
	} else {
		var ai1, ai2 string
		if aiTypes == nil {
			ai1, ai2 = models.Fato, models.Fato // TODO: choose the best AI by default
		} else {
			aiTypeSplit := strings.Split(*aiTypes, ":")
			ai1, ai2 = aiTypeSplit[0], aiTypeSplit[1]
		}
		start := time.Now()
		aivsai.RunAIvsAI(ai1, ai2, *matches, *format, *logging)
		elapsed := time.Since(start)
		fmt.Printf("\nAI vs AI matches completed in %.2f seconds\n", elapsed.Seconds())
	}
}

// runServer starts the WebSocket server
func runServer(addr string) {
	fmt.Printf("Starting Stratego WebSocket Server on %s\n", addr)
	fmt.Println("Endpoints:")
	fmt.Println("  POST /api/users/register - Register a new user")
	fmt.Println("  POST /api/users/login - Login")
	fmt.Println("  GET  /api/users - Get user by ID")
	fmt.Println("  GET  /api/users/stats - Get user statistics")
	fmt.Println("  GET  /api/board-setups - Get board setups for user")
	fmt.Println("  POST /api/board-setups - Create board setup")
	fmt.Println("  PUT  /api/board-setups - Update board setup")
	fmt.Println("  DELETE /api/board-setups - Delete board setup")
	fmt.Println("  POST /api/games - Create a new game")
	fmt.Println("  GET  /api/games - List all games")
	fmt.Println("  WS   /ws/game/{gameID}?player={0|1|spectator} - Connect to game")

	server := api.NewGameServer()
	if err := server.StartServer(addr); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
