// Package main is the entry point for autonomous simulation
package main

import (
	aivsai "digital-innovation/gostrategy/internal/ai/AIvsAI"
	"digital-innovation/gostrategy/internal/ai/playground"
	"digital-innovation/gostrategy/internal/ai/playground/tui"
	"digital-innovation/gostrategy/internal/models"
	"flag"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func run() error {
	aiTypes := flag.String("ai", "", "Run AI vs AI matches (e.g. 'fafo:fato')")
	matches := flag.Int("matches", 100, "Number of AI vs AI matches to run")
	format := flag.String("format", "none", "The format used to print the results of an AI vs AI competition, either none or md")
	loggingEnabled := flag.Bool("logging", true, "Show logs in stdout")
	tuiMode := flag.Bool("tui", false, "Run in interactive TUI mode")
	aggression := flag.Float64("aggression", 0.5, "FATO aggression value (0.0 to 1.0)")
	setup := flag.String("setup", "Honey Pot", "Board setup configuration ('Honey Pot', 'random', or 40-char string)")
	exportPath := flag.String("export", "", "Filepath to export simulation dataset to (JSON)")

	flag.Parse()

	// Launch interactive TUI if requested or if no arguments are passed
	if *tuiMode || *aiTypes == "" {
		p := tea.NewProgram(tui.NewModel(), tea.WithAltScreen())
		_, err := p.Run()
		return err
	}

	fmt.Println("=== GoStrategy Autonomous Simulation Running ===")

	var ai1, ai2 string
	aiTypeSplit := strings.Split(*aiTypes, ":")
	if len(aiTypeSplit) == 2 {
		ai1, ai2 = aiTypeSplit[0], aiTypeSplit[1]
	} else {
		ai1, ai2 = models.Fato, models.Fato
	}

	start := time.Now()
	if *exportPath != "" {
		err := playground.StartHeadless(ai1, ai2, *matches, *aggression, *setup, *exportPath)
		if err != nil {
			return err
		}
	} else {
		aivsai.RunAIvsAI(ai1, ai2, *matches, *format, *loggingEnabled)
	}
	elapsed := time.Since(start)
	fmt.Printf("\nAI vs AI matches completed in %.2f seconds\n", elapsed.Seconds())
	return nil
}
