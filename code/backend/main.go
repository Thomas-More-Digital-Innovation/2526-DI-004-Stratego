package main

import (
	aivsai "digital-innovation/stratego/ai/AIvsAI"
	"digital-innovation/stratego/api"
	"digital-innovation/stratego/db"
	"digital-innovation/stratego/logging"
	"digital-innovation/stratego/models"
	"digital-innovation/stratego/utils"
	"flag"
	"fmt"
	"strings"
	"time"
)

func main() {
	if err := run(); err != nil {
		logging.Fatalf("%v", err)
	}
}

func run() error {
	serverMode := flag.Bool("server", false, "Run in WebSocket server mode")
	defaultAddr := fmt.Sprintf(":%s", utils.GetEnv("PORT", "8080"))
	addr := flag.String("addr", defaultAddr, "Server address")
	aiTypes := flag.String("ai", "fafo:fafo", "Run AI vs AI matches instead of server")
	matches := flag.Int("matches", 100, "Number of AI vs AI matches to run")
	format := flag.String("format", "none", "The format used to print the results of an AI vs AI competition, either none or md")
	loggingEnabled := flag.Bool("logging", true, "Show logs in stdout")

	flag.Parse()

	fmt.Println("=== Stratego Backend Running ===")

	if *serverMode {
		if err := db.InitDB(); err != nil {
			return fmt.Errorf("failed to initialize database: %w", err)
		}
		defer func() {
			if err := db.CloseDB(); err != nil {
				logging.Error("Error closing database", err)
			}
		}()

		return startServer(*addr)
	}

	var ai1, ai2 string
	if aiTypes == nil {
		ai1, ai2 = models.Fato, models.Fato
	} else {
		aiTypeSplit := strings.Split(*aiTypes, ":")
		ai1, ai2 = aiTypeSplit[0], aiTypeSplit[1]
	}

	start := time.Now()
	aivsai.RunAIvsAI(ai1, ai2, *matches, *format, *loggingEnabled)
	elapsed := time.Since(start)
	fmt.Printf("\nAI vs AI matches completed in %.2f seconds\n", elapsed.Seconds())
	return nil
}

func startServer(addr string) error {
	fmt.Printf("Starting Stratego Game Server on %s\n", addr)

	server := api.NewGameServer()
	if err := server.StartServer(addr); err != nil {
		return fmt.Errorf("server error: %w", err)
	}
	return nil
}
