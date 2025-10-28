package api

import "digital-innovation/stratego/engine"

// WebSocket message types
const (
	// Client -> Server
	MsgTypeMove         = "move"
	MsgTypeRequestState = "requestState"
	MsgTypePing         = "ping"

	// Server -> Client
	MsgTypeGameState  = "gameState"
	MsgTypeMoveResult = "moveResult"
	MsgTypeAIMove     = "aiMove"
	MsgTypeGameOver   = "gameOver"
	MsgTypeError      = "error"
	MsgTypePong       = "pong"
	MsgTypeBoardState = "boardState"
)

// Base message structure
type WSMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
}

// Client messages
type MoveMessage struct {
	From PositionDTO `json:"from"`
	To   PositionDTO `json:"to"`
}

// Server messages
type GameStateMessage struct {
	Round              int    `json:"round"`
	CurrentPlayerID    int    `json:"currentPlayerId"`
	CurrentPlayerName  string `json:"currentPlayerName"`
	IsGameOver         bool   `json:"isGameOver"`
	WinnerID           *int   `json:"winnerId,omitempty"`
	WinnerName         string `json:"winnerName,omitempty"`
	WinCause           string `json:"winCause,omitempty"`
	Player1Score       int    `json:"player1Score"`
	Player2Score       int    `json:"player2Score"`
	WaitingForInput    bool   `json:"waitingForInput"`
	MoveCount          int    `json:"moveCount"`
	Player1AlivePieces int    `json:"player1AlivePieces"`
	Player2AlivePieces int    `json:"player2AlivePieces"`
}

type MoveResultMessage struct {
	Success      bool       `json:"success"`
	Error        string     `json:"error,omitempty"`
	Move         MoveDTO    `json:"move,omitempty"`
	AttackerDied bool       `json:"attackerDied"`
	DefenderDied bool       `json:"defenderDied"`
	CombatResult *CombatDTO `json:"combatResult,omitempty"`
}

type AIMoveMessage struct {
	PlayerID     int        `json:"playerId"`
	PlayerName   string     `json:"playerName"`
	Move         MoveDTO    `json:"move"`
	AttackerDied bool       `json:"attackerDied"`
	DefenderDied bool       `json:"defenderDied"`
	CombatResult *CombatDTO `json:"combatResult,omitempty"`
}

type GameOverMessage struct {
	WinnerID   *int   `json:"winnerId,omitempty"`
	WinnerName string `json:"winnerName,omitempty"`
	WinCause   string `json:"winCause"`
	Round      int    `json:"round"`
}

type ErrorMessage struct {
	Error string `json:"error"`
}

type BoardStateMessage struct {
	Board  [][]PieceDTO `json:"board"`
	Width  int          `json:"width"`
	Height int          `json:"height"`
}

// DTOs for data transfer
type PositionDTO struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type MoveDTO struct {
	From PositionDTO `json:"from"`
	To   PositionDTO `json:"to"`
}

type PieceDTO struct {
	Type      string      `json:"type,omitempty"`
	Rank      string      `json:"rank,omitempty"`
	OwnerID   int         `json:"ownerId"`
	OwnerName string      `json:"ownerName,omitempty"`
	Revealed  bool        `json:"revealed"`
	Icon      string      `json:"icon,omitempty"`
	Position  PositionDTO `json:"position"`
}

type CombatDTO struct {
	AttackerRank     string `json:"attackerRank"`
	DefenderRank     string `json:"defenderRank"`
	AttackerRevealed bool   `json:"attackerRevealed"`
	DefenderRevealed bool   `json:"defenderRevealed"`
}

// Helper functions to convert engine types to DTOs
func PositionToDTO(pos engine.Position) PositionDTO {
	return PositionDTO{X: pos.X, Y: pos.Y}
}

func MoveToDTO(move engine.Move) MoveDTO {
	return MoveDTO{
		From: PositionToDTO(move.GetFrom()),
		To:   PositionToDTO(move.GetTo()),
	}
}

func PieceToDTO(piece *engine.Piece, viewerID int) PieceDTO {
	if piece == nil {
		return PieceDTO{}
	}

	ownerID := piece.GetOwner().GetID()
	canSee := piece.IsRevealed() || ownerID == viewerID

	dto := PieceDTO{
		OwnerID:   ownerID,
		OwnerName: piece.GetOwner().GetName(),
		Revealed:  piece.IsRevealed(),
		Position:  PositionDTO{}, // Will be set by caller
	}

	if canSee {
		pieceType := piece.GetType()
		dto.Type = pieceType.GetName()
		dto.Rank = string(pieceType.GetRank())
		dto.Icon = pieceType.GetIcon()
	}

	return dto
}
