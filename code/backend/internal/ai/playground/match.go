package playground

import (
	AIhandler "digital-innovation/gostrategy/internal/ai/handler"
	"digital-innovation/gostrategy/internal/game"
	"fmt"
)

// RunSingleMatch simulates a single match of AI vs AI with the configured parameters
func RunSingleMatch(index int, aliceAI, bobAI string, aggression float64, setupStr string, pAlice, pBob *game.Player) (GameTelemetryExport, error) {
	// Create controllers
	optsAlice := map[string]any{"aggression": aggression}
	optsBob := map[string]any{"aggression": aggression}

	ctrlAlice, err := AIhandler.CreateAIWithOptions(aliceAI, pAlice, optsAlice)
	if err != nil {
		return GameTelemetryExport{}, fmt.Errorf("failed to create Alice AI: %v", err)
	}

	ctrlBob, err := AIhandler.CreateAIWithOptions(bobAI, pBob, optsBob)
	if err != nil {
		return GameTelemetryExport{}, fmt.Errorf("failed to create Bob AI: %v", err)
	}

	// Alternate who goes first based on match index
	var g *game.Game
	if index%2 == 0 {
		g = game.NewGame(ctrlAlice, ctrlBob)
	} else {
		g = game.NewGame(ctrlBob, ctrlAlice)
	}

	// Prepare setups
	var pAlicePieces, pBobPieces []*game.Piece
	if setupStr == "random" || setupStr == "" {
		pAlicePieces = game.RandomSetup(pAlice)
		pBobPieces = game.RandomSetup(pBob)
	} else {
		normalized, err := NormalizeSetup(setupStr)
		if err != nil {
			return GameTelemetryExport{}, fmt.Errorf("invalid setup string: %v", err)
		}
		pAlicePieces, err = game.ParseSetup(pAlice, []byte(normalized))
		if err != nil {
			return GameTelemetryExport{}, fmt.Errorf("invalid Alice setup layout: %v", err)
		}
		pBobPieces, err = game.ParseSetup(pBob, []byte(normalized))
		if err != nil {
			return GameTelemetryExport{}, fmt.Errorf("invalid Bob setup layout: %v", err)
		}
	}

	// Place pieces and prepare game
	if err := game.SetupGame(g, pAlicePieces, pBobPieces); err != nil {
		return GameTelemetryExport{}, fmt.Errorf("failed to place setup pieces: %v", err)
	}
	g.InitialState = g.GetInitialBoardState()

	// Execute headless simulation
	runner := game.NewRunner(g, 0, 1000)
	winner := runner.RunToCompletion()

	winnerID := -1
	if winner != nil {
		if winner.GetName() == pAlice.GetName() {
			winnerID = 0
		} else {
			winnerID = 1
		}
	}

	return GameTelemetryExport{
		GameIndex:    index,
		WinnerID:     winnerID,
		WinCause:     string(g.GetWinCause()),
		TotalTurns:   g.GetRound(),
		InitialBoard: g.InitialState,
		Moves:        g.HistoricalHistory,
	}, nil
}
