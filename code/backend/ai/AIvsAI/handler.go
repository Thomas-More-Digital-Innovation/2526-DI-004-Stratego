package aivsai

import (
	"digital-innovation/stratego/models"
	"fmt"
)

func RunAIvsAI(ai1, ai2 string, matches int, format string, logging bool) {
	summary := runAIvsAI(ai1, ai2, matches, logging)

	switch format {
	case "md":
		printMarkdownSummary(summary, matches)
	default:
		printDefaultSummary(summary, matches)
	}

}

func printMarkdownSummary(summary models.GameSummary, matches int) {
	// Top-level summary
	fmt.Printf("\n### AI vs AI Tournament Summary (%d games)\n\n", matches)
	fmt.Printf("**Total Matches:** %d  \n**Total Rounds:** %d  \n**Average Rounds (per game):** %.2f  \n**Shortest Game (rounds):** %d\n\n",
		matches, summary.TotalRounds, summary.AverageRounds, summary.LeastRounds)

	// Overall win causes (aggregate)
	totalFlag := summary.WinCauseFlagCaptured
	totalNoMoves := summary.WinCauseNoMovesWins
	totalMaxTurns := summary.WinCauseMaxTurns
	wonMatches := float64(matches - summary.Draws)

	fmt.Println("#### Overall Win Causes")
	fmt.Println()
	fmt.Println("| Cause | Count | % |")
	fmt.Println("|-------:|------:|---:|")
	fmt.Printf("| Flag captured | %d | %.1f%% |\n", totalFlag, float64(totalFlag)*100.0/wonMatches)
	fmt.Printf("| No movable pieces | %d | %.1f%% |\n", totalNoMoves, float64(totalNoMoves)*100.0/wonMatches)
	fmt.Printf("| Max turns | %d | %.1f%% |\n\n", totalMaxTurns, float64(totalMaxTurns)*100.0/wonMatches)

	// Per-player summary table
	p1 := summary.Player1data
	p2 := summary.Player2data

	fmt.Println("#### Player Results")
	fmt.Println()
	fmt.Println("| Player | Wins | Win % | Flag captures | No-move wins | Max-turn wins |")
	fmt.Println("|:-------|-----:|-----:|--------------:|-------------:|--------------:|")
	fmt.Printf("| %s | %d | %.1f%% | %d | %d | %d |\n",
		p1.Name, p1.Wins, float64(p1.Wins)*100.0/wonMatches, p1.WinCauseFlagCaptured, p1.WinCauseNoMovesWin, p1.WinCauseMaxTurns)
	fmt.Printf("| %s | %d | %.1f%% | %d | %d | %d |\n\n",
		p2.Name, p2.Wins, float64(p2.Wins)*100.0/wonMatches, p2.WinCauseFlagCaptured, p2.WinCauseNoMovesWin, p2.WinCauseMaxTurns)

	// Draws
	fmt.Printf("**Draws:** %d (%.1f%%)\n", summary.Draws, float64(summary.Draws)*100.0/wonMatches)
}

func printDefaultSummary(summary models.GameSummary, matches int) {
	// Human-readable plain text summary
	fmt.Println()
	fmt.Println("========================================")
	fmt.Printf("AI vs AI Tournament Summary (%d games)\n", matches)
	fmt.Println("========================================")
	fmt.Printf("Total Matches: %d\n", matches)
	fmt.Printf("Total Rounds: %d\n", summary.TotalRounds)
	fmt.Printf("Average Rounds (per game): %.2f\n", summary.AverageRounds)
	fmt.Printf("Shortest Game (rounds): %d\n", summary.LeastRounds)
	fmt.Println("----------------------------------------")

	// Overall win causes
	fmt.Println("Overall Win Causes:")
	totalFlag := summary.WinCauseFlagCaptured
	totalNoMoves := summary.WinCauseNoMovesWins
	totalMaxTurns := summary.WinCauseMaxTurns
	wonMatches := float64(matches - summary.Draws)

	fmt.Printf("  Flag captured:     %d (%.1f%%)\n", totalFlag, float64(totalFlag)*100.0/wonMatches)
	fmt.Printf("  No movable pieces: %d (%.1f%%)\n", totalNoMoves, float64(totalNoMoves)*100.0/wonMatches)
	fmt.Printf("  Max turns:         %d (%.1f%%)\n", totalMaxTurns, float64(totalMaxTurns)*100.0/wonMatches)
	fmt.Println("----------------------------------------")

	// Per-player breakdown
	p1 := summary.Player1data
	p2 := summary.Player2data
	fmt.Printf("Player: %s\n", p1.Name)
	fmt.Printf("  Wins: %d (%.1f%%)\n", p1.Wins, float64(p1.Wins)*100.0/wonMatches)
	fmt.Printf("  Win causes:\n")
	fmt.Printf("    Flag captured:     %d\n", p1.WinCauseFlagCaptured)
	fmt.Printf("    No movable pieces: %d\n", p1.WinCauseNoMovesWin)
	fmt.Printf("    Max turns:         %d\n", p1.WinCauseMaxTurns)
	fmt.Println("----------------------------------------")
	fmt.Printf("Player: %s\n", p2.Name)
	fmt.Printf("  Wins: %d (%.1f%%)\n", p2.Wins, float64(p2.Wins)*100.0/wonMatches)
	fmt.Printf("  Win causes:\n")
	fmt.Printf("    Flag captured:     %d\n", p2.WinCauseFlagCaptured)
	fmt.Printf("    No movable pieces: %d\n", p2.WinCauseNoMovesWin)
	fmt.Printf("    Max turns:         %d\n", p2.WinCauseMaxTurns)
	fmt.Println("----------------------------------------")

	fmt.Printf("Draws: %d (%.1f%%)\n", summary.Draws, float64(summary.Draws)*100.0/wonMatches)
}
