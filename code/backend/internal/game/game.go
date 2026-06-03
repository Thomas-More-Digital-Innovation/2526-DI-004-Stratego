// Package game handles the game state and flow
package game

import (
	"digital-innovation/gostrategy/internal/game/models"
)

// WinCause defines why a game ended
type WinCause string

// Win causes
const (
	WinCauseFlagCaptured    WinCause = "flag_captured"
	WinCauseNoMovablePieces WinCause = "no_movable_pieces"
	WinCauseMaxTurns        WinCause = "max_turns"
)

// CombatResult stores details about a resolved battle
type CombatResult struct {
	Occurred         bool
	AttackerPiece    *Piece
	DefenderPiece    *Piece
	AttackerPosition Position
	DefenderPosition Position
}

// Game represents a single match instance
type Game struct {
	Players           []*Player
	PlayerControllers []PlayerController // AI or Human controllers
	Board             *Board
	CurrentPlayer     *Player
	CurrentController PlayerController
	MoveHistory       []Move
	HistoricalHistory []models.HistoricalMove
	InitialState      [][]models.PieceData
	LastCombat        *CombatResult // Track last combat for broadcasting
	round             int
	winner            *Player
	winCause          WinCause
	gameOver          bool
}

// NewGame creates a new Game instance with given controllers
func NewGame(controller1, controller2 PlayerController) *Game {
	board := NewBoard()
	player1 := controller1.GetPlayer()
	player2 := controller2.GetPlayer()

	return &Game{
		Players:           []*Player{player1, player2},
		PlayerControllers: []PlayerController{controller1, controller2},
		Board:             board,
		CurrentPlayer:     player1,
		CurrentController: controller1,
		MoveHistory:       []Move{},
		HistoricalHistory: []models.HistoricalMove{},
		round:             1,
		gameOver:          false,
	}
}

// Clone creates a deep copy of the game state (excluding controllers which are shared or recreated)
func (g *Game) Clone() *Game {
	clonedBoard := g.Board.Clone()

	// Create new game with same controllers
	cloned := NewGame(g.PlayerControllers[0], g.PlayerControllers[1])
	cloned.Board = clonedBoard
	cloned.round = g.round
	cloned.gameOver = g.gameOver
	cloned.winCause = g.winCause
	cloned.winner = g.winner

	// Set current player and controller correctly
	if g.CurrentPlayer == g.Players[0] {
		cloned.CurrentPlayer = cloned.Players[0]
		cloned.CurrentController = cloned.PlayerControllers[0]
	} else {
		cloned.CurrentPlayer = cloned.Players[1]
		cloned.CurrentController = cloned.PlayerControllers[1]
	}

	return cloned
}

// NextTurn advances the game state to the next player's turn
func (g *Game) NextTurn() {
	switch {
	case g.Players[0].HasWon():
		g.winner = g.Players[0]
		g.winCause = WinCauseFlagCaptured
		g.gameOver = true
	case g.Players[1].HasWon():
		g.winner = g.Players[1]
		g.winCause = WinCauseFlagCaptured
		g.gameOver = true
	case g.CurrentPlayer == g.Players[0]:
		g.CurrentPlayer = g.Players[1]
		g.CurrentController = g.PlayerControllers[1]
	default:
		g.CurrentPlayer = g.Players[0]
		g.CurrentController = g.PlayerControllers[0]
		g.round++
		// Hide all revealed pieces at the start of a new round
		g.HideAllRevealedPieces()
	}
}

// IsGameOver returns true if the game has ended
func (g *Game) IsGameOver() bool {
	return g.gameOver
}

// GetCurrentController returns the controller for the current player
func (g *Game) GetCurrentController() PlayerController {
	return g.CurrentController
}

// GetRound returns the current round number
func (g *Game) GetRound() int {
	return g.round
}

// GetWinner returns the player who won the game, or nil
func (g *Game) GetWinner() *Player {
	return g.winner
}

// GetWinCause returns the cause of victory
func (g *Game) GetWinCause() WinCause {
	return g.winCause
}

// SetWinner manually declares a winner and ends the game
func (g *Game) SetWinner(player *Player, cause WinCause) {
	g.winner = player
	g.winCause = cause
	g.gameOver = true
	player.SetWinner()
}
