package models

// GameState represents the current state of a game (for API responses)
type GameState struct {
	Round              int    `json:"round"`
	CurrentPlayerID    int    `json:"currentPlayerId"`
	CurrentPlayerName  string `json:"currentPlayerName"`
	IsGameOver         bool   `json:"isGameOver"`
	WinnerID           *int   `json:"winnerId,omitempty"`
	Player1Score       int    `json:"player1Score"`
	Player2Score       int    `json:"player2Score"`
	WaitingForInput    bool   `json:"waitingForInput"`
	Paused             bool   `json:"paused"`
	MoveCount          int    `json:"moveCount"`
	Player1AlivePieces int    `json:"player1AlivePieces"`
	Player2AlivePieces int    `json:"player2AlivePieces"`
	IsSetupPhase       bool   `json:"isSetupPhase"`
	Headless           bool   `json:"headless"`
	SetupRemainingSecs int    `json:"setupRemainingSecs,omitempty"`
}

// GameSummary represents a lightweight overview of a game
type GameSummary struct {
	GameID       string `json:"gameId"`
	Round        int    `json:"round"`
	IsRunning    bool   `json:"isRunning"`
	IsGameOver   bool   `json:"isGameOver"`
	IsSetupPhase bool   `json:"isSetupPhase"`
	GameType     string `json:"gameType"`
}
