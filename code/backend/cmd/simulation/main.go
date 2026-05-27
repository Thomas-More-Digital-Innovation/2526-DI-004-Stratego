// Package main is the entry point for autonomous simulation
package main

import (
	"digital-innovation/gostrategy/internal/models"
	aivsai "digital-innovation/gostrategy/pkg/ai/AIvsAI"
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
