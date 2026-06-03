package playground

import (
	"digital-innovation/gostrategy/internal/game"
	tea "github.com/charmbracelet/bubbletea"
)

// ProgressMsg streams live progress updates from the simulation goroutine
type ProgressMsg struct {
	Current int
	Total   int
}

// ResultsMsg signals simulation completion with final telemetry
type ResultsMsg struct {
	Results SimulationExport
	Error   error
}

// ListenToSimChan listens to progress/results updates from the simulation channel
func ListenToSimChan(ch <-chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-ch
		if !ok {
			return nil
		}
		return msg
	}
}

// StartHeadless runs a simulation in headless mode without booting the TUI
func StartHeadless(aliceAI, bobAI string, matches int, aggression float64, setupStr string, exportPath string) error {
	playerAlice := game.NewPlayer(0, "Alice AI - "+aliceAI, "red")
	playerBob := game.NewPlayer(1, "Bob AI - "+bobAI, "blue")

	sim := SimulationExport{
		MatchesCount: matches,
		AliceAI:      aliceAI,
		BobAI:        bobAI,
		Games:        make([]GameTelemetryExport, matches),
	}

	for i := 0; i < matches; i++ {
		gameExport, err := RunSingleMatch(i, aliceAI, bobAI, aggression, setupStr, &playerAlice, &playerBob)
		if err != nil {
			return err
		}
		sim.Games[i] = gameExport
		sim.TotalRounds += gameExport.TotalTurns
		switch gameExport.WinnerID {
		case 0:
			sim.AliceWins++
		case 1:
			sim.BobWins++
		default:
			sim.Draws++
		}
	}
	sim.AvgRounds = float64(sim.TotalRounds) / float64(matches)

	if exportPath != "" {
		return ExportSimulation(exportPath, sim)
	}
	return nil
}
