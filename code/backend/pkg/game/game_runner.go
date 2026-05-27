package game

import (
	"sync"
	"time"
)

// Runner handles the execution of a game turn-by-turn
type Runner struct {
	game                 *Game
	turnDelay            time.Duration // Optional delay between AI turns for visualization, can be 0 to remove the delay
	maxTurns             int
	waitingForHumanInput bool
	paused               bool // flag to indicate if game is paused
	onMoveExecuted       func()
	stopped              bool
	stopChan             chan bool
	locker               sync.Locker // Mutex from GameSession to prevent race conditions
	stateMu              sync.RWMutex
}

// NewRunner creates a new Runner instance
func NewRunner(game *Game, turnDelay time.Duration, maxTurns int) *Runner {
	if maxTurns <= 0 {
		maxTurns = 1000 // Default safety limit of 1000 turns, prevents infinite loops (especially for AI vs AI)
	}
	return &Runner{
		game:      game,
		turnDelay: turnDelay,
		maxTurns:  maxTurns,
		paused:    false,
		stopChan:  make(chan bool, 1),
	}
}

// SetMoveCallback sets the callback to be called when a move is executed
func (gr *Runner) SetMoveCallback(callback func()) {
	gr.onMoveExecuted = callback
}

// getOpponent returns the opponent of the given player
func (gr *Runner) getOpponent(player *Player) *Player {
	if gr.game.Players[0] == player {
		return gr.game.Players[1]
	}
	return gr.game.Players[0]
}

// GetGame returns the underlying game
func (gr *Runner) GetGame() *Game {
	return gr.game
}

// SetTurnDelay sets the delay between AI turns
func (gr *Runner) SetTurnDelay(delay time.Duration) {
	gr.turnDelay = delay
}
