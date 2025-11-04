package aivsai

import (
	"digital-innovation/stratego/models"
	"fmt"
)

func RunAIvsAI(ai1, ai2 string, matches int, format string) {
	// currently we are only running Fafo vs Fafo
	summary := runAIvsAI(matches)

	switch format {
	case "md":
		printMarkdownSummary(summary, matches)
	default:
		printDefaultSummary(summary, matches)
	}

}

func printMarkdownSummary(summary models.GameSummary, matches int) {
	fmt.Printf("\n### 📊 Tournament Results (%d games)\n", matches)
	fmt.Println("\n| Player | Wins | Win % |")
	fmt.Println("|--------|------|-------|")
	fmt.Printf("| %v | %3d | %.1f%% |\n", summary.Player1Name, summary.Player1Wins, float64(summary.Player1Wins*100)/float64(matches))
	fmt.Printf("| %v | %3d | %.1f%% |\n", summary.Player2Name, summary.Player2Wins, float64(summary.Player2Wins*100)/float64(matches))
	fmt.Printf("| Draws | %3d | %.1f%% |\n", summary.Draws, float64(summary.Draws*100)/float64(matches))
	fmt.Println("\n| Win Causes | Count | % |")
	fmt.Println("|-------------|-------|---|")
	fmt.Printf("| Flag captured | %3d | %.1f%% |\n", int(summary.WinCauseFlagCaptured), float64(summary.WinCauseFlagCaptured*100)/float64(matches))
	fmt.Printf("| No movable pieces | %3d | %.1f%% |\n", int(summary.WinCauseNoMovesWins), float64(summary.WinCauseNoMovesWins*100)/float64(matches))
	fmt.Printf("| Max turns | %3d | %.1f%% |\n\n", int(summary.WinCauseMaxTurns), float64(summary.WinCauseMaxTurns*100)/float64(matches))
	fmt.Println("Average game length:", fmt.Sprintf("%.1f rounds", summary.AverageRounds))
}

func printDefaultSummary(summary models.GameSummary, matches int) {
	fmt.Println("\n========================================")
	fmt.Printf("📊 Tournament Results (%d games)\n", matches)
	fmt.Println("========================================")
	fmt.Printf("%v: %3d (%.1f%%)\n", summary.Player1Name, summary.Player1Wins, float64(summary.Player1Wins*100)/float64(matches))
	fmt.Printf("%v:   %3d (%.1f%%)\n", summary.Player2Name, summary.Player2Wins, float64(summary.Player2Wins*100)/float64(matches))
	fmt.Printf("Draws:      %3d (%.1f%%)\n", summary.Draws, float64(summary.Draws*100)/float64(matches))
	fmt.Println("========================================")
	fmt.Println("Win Causes:")
	fmt.Printf("  Flag captured:     %3d (%.1f%%)\n", int(summary.WinCauseFlagCaptured), float64(summary.WinCauseFlagCaptured*100)/float64(matches))
	fmt.Printf("  No movable pieces: %3d (%.1f%%)\n", int(summary.WinCauseNoMovesWins), float64(summary.WinCauseNoMovesWins*100)/float64(matches))
	fmt.Printf("  Max turns:         %3d (%.1f%%)\n", int(summary.WinCauseMaxTurns), float64(summary.WinCauseMaxTurns*100)/float64(matches))
	fmt.Println("========================================")
	fmt.Printf("Average game length: %.1f rounds\n", summary.AverageRounds)
	fmt.Println("========================================")

}
