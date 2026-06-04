// Package playground implements the local TUI and simulation runner.
package playground

import (
	"digital-innovation/gostrategy/internal/game/models"
	"encoding/json"
	"os"
)

// SimulationExport stores the entire simulation telemetry dataset
type SimulationExport struct {
	MatchesCount int                   `json:"matchesCount"`
	TotalRounds  int                   `json:"totalRounds"`
	AvgRounds    float64               `json:"avgRounds"`
	AliceAI      string                `json:"aliceAi"`
	BobAI        string                `json:"bobAi"`
	AliceWins    int                   `json:"aliceWins"`
	BobWins      int                   `json:"bobWins"`
	Draws        int                   `json:"draws"`
	Games        []GameTelemetryExport `json:"games"`
}

// GameTelemetryExport stores specific details of an executed match
type GameTelemetryExport struct {
	GameIndex    int                     `json:"gameIndex"`
	WinnerID     int                     `json:"winnerId"` // 0 = Alice, 1 = Bob, -1 = Draw
	WinCause     string                  `json:"winCause"`
	TotalTurns   int                     `json:"totalTurns"`
	InitialBoard [][]models.PieceData    `json:"initialBoard"`
	Moves        []models.HistoricalMove `json:"moves"`
}

// ExportSimulation writes the simulation data structure to a JSON file
func ExportSimulation(filepath string, data SimulationExport) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath, bytes, 0600)
}
