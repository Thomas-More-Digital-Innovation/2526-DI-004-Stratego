// Package main is the entry point for autonomous simulation
package main

import (
	aivsai "digital-innovation/gostrategy/internal/ai/AIvsAI"
	"digital-innovation/gostrategy/internal/models"
	"flag"
	"fmt"
	"strings"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func run() error {
	aiTypes := flag.String("ai", "fafo:fafo", "Run AI vs AI matches")
	matches := flag.Int("matches", 100, "Number of AI vs AI matches to run")
	format := flag.String("format", "none", "The format used to print the results of an AI vs AI competition, either none or md")
	loggingEnabled := flag.Bool("logging", true, "Show logs in stdout")

	flag.Parse()

	fmt.Println("=== GoStrategy Autonomous Simulation Running ===")

	var ai1, ai2 string
	switch {
	case aiTypes == nil:
		fmt.Print("No AIs specified, running default FATO vs FATO")
		ai1, ai2 = models.Fato, models.Fato
	case strings.Contains(*aiTypes, ":"):
		aiTypeSplit := strings.Split(*aiTypes, ":")
		ai1, ai2 = aiTypeSplit[0], aiTypeSplit[1]
	default:
		ai1, ai2 = *aiTypes, *aiTypes

	}

	start := time.Now()
	aivsai.RunAIvsAI(ai1, ai2, *matches, *format, *loggingEnabled)
	elapsed := time.Since(start)
	fmt.Printf("\nAI vs AI matches completed in %.2f seconds\n", elapsed.Seconds())
	return nil
}
