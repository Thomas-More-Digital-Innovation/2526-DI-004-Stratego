package models

type AiTournamentData struct {
	Name                 string
	Wins                 int
	WinCauseFlagCaptured int
	WinCauseNoMovesWin   int
	WinCauseMaxTurns     int
}

type GameSummary struct {
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
