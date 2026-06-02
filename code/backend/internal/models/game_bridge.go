package models

import gamemodels "digital-innovation/gostrategy/pkg/game/models"

// PieceType is an alias for the game engine PieceType definition.
type PieceType = gamemodels.PieceType

// HistoricalMove is an alias for the game engine HistoricalMove representation.
type HistoricalMove = gamemodels.HistoricalMove

// PieceData is an alias for the game engine PieceData representation.
type PieceData = gamemodels.PieceData

// MoveResultType is an alias for the game engine MoveResultType outcome representation.
type MoveResultType = gamemodels.MoveResultType

// GameHistory is an alias for the game engine GameHistory log representation.
type GameHistory = gamemodels.GameHistory

// GameState is an alias for the game engine GameState representation.
type GameState = gamemodels.GameState

// AiTournamentData is an alias for the AI performance statistics.
type AiTournamentData = gamemodels.AiTournamentData

// AiGameSummary is an alias for the overall AI tournament outcomes.
type AiGameSummary = gamemodels.AiGameSummary

// Constant aliases
const (
	HumanVsAi    = gamemodels.HumanVsAi
	HumanVsHuman = gamemodels.HumanVsHuman
	AiVsAi       = gamemodels.AiVsAi

	ResultMove    = gamemodels.ResultMove
	ResultWin     = gamemodels.ResultWin
	ResultLoss    = gamemodels.ResultLoss
	ResultTie     = gamemodels.ResultTie
	ResultCapture = gamemodels.ResultCapture

	Fafo = gamemodels.Fafo
	Fato = gamemodels.Fato
)

// Variable aliases
var (
	Flag       = gamemodels.Flag
	Bomb       = gamemodels.Bomb
	Spy        = gamemodels.Spy
	Scout      = gamemodels.Scout
	Miner      = gamemodels.Miner
	Sergeant   = gamemodels.Sergeant
	Lieutenant = gamemodels.Lieutenant
	Captain    = gamemodels.Captain
	Major      = gamemodels.Major
	Colonel    = gamemodels.Colonel
	General    = gamemodels.General
	Marshal    = gamemodels.Marshal
)

// NewPieceType constructor alias
func NewPieceType(name string, rank byte, movable bool, description string, icon string, count int, strategicValue int) *PieceType {
	return gamemodels.NewPieceType(name, rank, movable, description, icon, count, strategicValue)
}
