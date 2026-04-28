package models

// HistoricalMove represents a move in the game history with its outcome
type HistoricalMove struct {
	MoveIndex int            `json:"moveIndex"`
	PlayerID  int            `json:"playerId"`
	FromX     int            `json:"fromX"`
	FromY     int            `json:"fromY"`
	ToX       int            `json:"toX"`
	ToY       int            `json:"toY"`
	Attacker  *PieceData     `json:"attacker,omitempty"`
	Defender  *PieceData     `json:"defender,omitempty"`
	Result    MoveResultType `json:"result"`
}

type PieceData struct {
	Type    string `json:"type"`
	Rank    string `json:"rank"`
	OwnerID int    `json:"ownerId"`
}

type MoveResultType string

const (
	ResultMove    MoveResultType = "move"    // Normal move to empty cell
	ResultWin     MoveResultType = "win"     // Attacker won combat
	ResultLoss    MoveResultType = "loss"    // Attacker lost combat
	ResultTie     MoveResultType = "tie"     // Both pieces died
	ResultCapture MoveResultType = "capture" // Flag captured (game over)
)

// GameHistory represents the full history of a game
type GameHistory struct {
	GameID       string           `json:"gameId"`
	InitialState interface{}      `json:"initialState"` // Will be JSONB/Array
	Moves        []HistoricalMove `json:"moves"`
	WinnerID     *int             `json:"winnerId"`
}
