// Package models defines the data structures used in the game
package models

// AiTournamentData contains statistics for an AI player in a tournament
type AiTournamentData struct {
	Name                 string
	Wins                 int
	WinCauseFlagCaptured int
	WinCauseNoMovesWin   int
	WinCauseMaxTurns     int
}

// AiGameSummary contains the overall results of an AI vs AI match set
type AiGameSummary struct {
	Player1data          AiTournamentData
	Player2data          AiTournamentData
	Draws                int
	TotalRounds          int
	AverageRounds        float64
	LeastRounds          int
	Matches              int
	WinCauseFlagCaptured int
	WinCauseNoMovesWins  int
	WinCauseMaxTurns     int
}
