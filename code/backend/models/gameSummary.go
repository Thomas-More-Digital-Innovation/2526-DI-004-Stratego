package models

type GameSummary struct {
	Player1Name          string
	Player2Name          string
	Player1Wins          int
	Player2Wins          int
	Draws                int
	TotalRounds          int
	AverageRounds        float64
	Matches              int
	WinCauseFlagCaptured float64
	WinCauseNoMovesWins  float64
	WinCauseMaxTurns     float64
}
