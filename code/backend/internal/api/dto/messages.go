// Package dto provides data transfer objects for WebSocket messages.
package dto

import (
	"digital-innovation/gostrategy/internal/models"
	"digital-innovation/gostrategy/pkg/game"
)

// WebSocket message types
const (
	// Client -> Server
	MsgTypeMove              = "move"
	MsgTypeGetValidMoves     = "getValidMoves"
	MsgTypePing              = "ping"
	MsgTypeAnimationComplete = "animationComplete"
	MsgTypeSwapPieces        = "swapPieces"
	MsgTypeRandomizeSetup    = "randomizeSetup"
	MsgTypeStartGame         = "startGame"
	MsgTypeLoadSetup         = "loadSetup"
	MsgTypePause             = "pause"
	MsgTypeUnpause           = "unpause"
	MsgTypeSetSpeed          = "setSpeed"
	MsgTypeStep              = "step"

	// Server -> Client
	MsgTypeGameState   = "gameState"
	MsgTypeMoveResult  = "moveResult"
	MsgTypeGameOver    = "gameOver"
	MsgTypeError       = "error"
	MsgTypePong        = "pong"
	MsgTypeBoardState  = "boardState"
	MsgTypeCombat      = "combat"
	MsgTypeValidMoves  = "validMoves"
	MsgTypeSetupPhase  = "setupPhase"
	MsgTypeMoveHistory = "moveHistory"
	MsgTypeMessage     = "message"
)

// WSMessage is the base structure for all WebSocket messages
type WSMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
}

// MoveMessage is sent by the client to request a move
type MoveMessage struct {
	From PositionDTO `json:"from"`
	To   PositionDTO `json:"to"`
}

// GetValidMovesMessage is sent by the client to request valid moves for a position
type GetValidMovesMessage struct {
	Position PositionDTO `json:"position"`
}

// SwapPiecesMessage is sent by the client to swap two pieces during setup
type SwapPiecesMessage struct {
	Pos1 PositionDTO `json:"pos1"`
	Pos2 PositionDTO `json:"pos2"`
}

// StartGameMessage is sent by the client to start the game
type StartGameMessage struct {
	Headless bool `json:"headless"`
}

// LoadSetupMessage is sent by the client to load a saved piece setup
type LoadSetupMessage struct {
	PlayerID  *int   `json:"playerId,omitempty"`
	SetupData string `json:"setupData"` // Base64 encoded 40 bytes
}

// RandomizeSetupMessage is sent by the client to request a random piece setup
type RandomizeSetupMessage struct {
	PlayerID *int `json:"playerId,omitempty"`
}

// SetSpeedMessage is sent by the client to change game speed
type SetSpeedMessage struct {
	SpeedMs int `json:"speedMs"`
}

// GameStateMessage contains the current state of the game
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
	Paused             bool   `json:"paused"`
	MoveCount          int    `json:"moveCount"`
	Player1AlivePieces int    `json:"player1AlivePieces"`
	Player2AlivePieces int    `json:"player2AlivePieces"`
	IsSetupPhase       bool   `json:"isSetupPhase"`
	Headless           bool   `json:"headless"`
	SetupRemainingSecs int    `json:"setupRemainingSecs,omitempty"`
	Player1Username    string `json:"player1Username"`
	Player2Username    string `json:"player2Username"`
}

// MoveResultMessage is sent by the server to report move success/failure
type MoveResultMessage struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// ValidMovesMessage is sent by the server with a list of valid moves
type ValidMovesMessage struct {
	Position   PositionDTO   `json:"position"`
	ValidMoves []PositionDTO `json:"validMoves"`
}

// GameOverMessage is sent by the server when the game ends
type GameOverMessage struct {
	WinnerID   *int   `json:"winnerId,omitempty"`
	WinnerName string `json:"winnerName,omitempty"`
	WinCause   string `json:"winCause"`
	Round      int    `json:"round"`
}

// ErrorMessage is sent by the server when an error occurs
type ErrorMessage struct {
	Error string `json:"error"`
}

// BoardStateMessage contains the current pieces on the board
type BoardStateMessage struct {
	Board    [][]PieceDTO           `json:"board"`
	Width    int                    `json:"width"`
	Height   int                    `json:"height"`
	LastMove *models.HistoricalMove `json:"lastMove,omitempty"`
}

// CombatMessage contains results of a combat between two pieces
type CombatMessage struct {
	Attacker     PieceDTO `json:"attacker"`
	Defender     PieceDTO `json:"defender"`
	AttackerWon  bool     `json:"attackerWon"`
	DefenderWon  bool     `json:"defenderWon"`
	AttackerDied bool     `json:"attackerDied"`
	DefenderDied bool     `json:"defenderDied"`
}

// MoveHistoryMessage contains the history of moves in a game
type MoveHistoryMessage struct {
	Moves        []MoveDTO               `json:"moves"`
	FullHistory  []models.HistoricalMove `json:"fullHistory"`
	InitialState [][]models.PieceData    `json:"initialState"`
}

// PositionDTO is a data transfer object for coordinates
type PositionDTO struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// MoveDTO is a data transfer object for moves
type MoveDTO struct {
	From PositionDTO `json:"from"`
	To   PositionDTO `json:"to"`
}

// PieceDTO is a data transfer object for pieces
type PieceDTO struct {
	Type      string      `json:"type,omitempty"`
	Rank      string      `json:"rank,omitempty"`
	OwnerID   int         `json:"ownerId"`
	OwnerName string      `json:"ownerName,omitempty"`
	Revealed  bool        `json:"revealed"`
	Icon      string      `json:"icon,omitempty"`
	Position  PositionDTO `json:"position"`
}

// CombatDTO is a data transfer object for combat results
type CombatDTO struct {
	AttackerRank     string `json:"attackerRank"`
	DefenderRank     string `json:"defenderRank"`
	AttackerRevealed bool   `json:"attackerRevealed"`
	DefenderRevealed bool   `json:"defenderRevealed"`
}

// PositionToDTO converts an game.Position to a PositionDTO
func PositionToDTO(pos game.Position) PositionDTO {
	return PositionDTO{X: pos.X, Y: pos.Y}
}

// MoveToDTO converts an game.Move to a MoveDTO
func MoveToDTO(move game.Move) MoveDTO {
	return MoveDTO{
		From: PositionToDTO(move.GetFrom()),
		To:   PositionToDTO(move.GetTo()),
	}
}

// PieceToDTO converts an game.Piece to a PieceDTO, hiding information if necessary
func PieceToDTO(piece *game.Piece, viewerID int, forceReveal bool) PieceDTO {
	if piece == nil {
		return PieceDTO{OwnerID: -1}
	}

	ownerID := piece.GetOwner().GetID()
	canSee := forceReveal || piece.IsRevealed() || ownerID == viewerID

	dto := PieceDTO{
		OwnerID:   ownerID,
		OwnerName: piece.GetOwner().GetName(),
		Revealed:  piece.IsRevealed() || forceReveal,
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
