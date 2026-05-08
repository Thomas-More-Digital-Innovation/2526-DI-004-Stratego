package game

import (
	"digital-innovation/gostrategy/engine"
	"digital-innovation/gostrategy/logging"
	"fmt"
	"math/rand"
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
	stopChan             chan bool
	locker               sync.Locker // Mutex from GameSession to prevent race conditions
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
	}
}

// SetMoveCallback sets the callback to be called when a move is executed
func (gr *Runner) SetMoveCallback(callback func()) {
	gr.onMoveExecuted = callback
}

// RunToCompletion runs the game until it's over (for AI vs AI)
// Winner can be nil when max turns are reached and both AIs have a similar piece count
func (gr *Runner) RunToCompletion() *engine.Player {
	turnCount := 0
	logging.Debug(logging.TagGame, "Runner: Starting RunToCompletion loop")

	for {
		var isGameOver bool
		if gr.locker != nil {
			gr.locker.Lock()
			isGameOver = gr.game.IsGameOver()
			gr.locker.Unlock()
		} else {
			isGameOver = gr.game.IsGameOver()
		}

		if isGameOver || turnCount >= gr.maxTurns {
			break
		}
		// Check for stop signal
		select {
		case <-gr.stopChan:

			logging.Debug(logging.TagGame, "Runner: Stop signal received, ending game")

			return nil
		default:
			// No stop signal, continue
		}

		if gr.IsPaused() {
			select {
			case <-gr.stopChan:
				logging.Debug(logging.TagGame, "Runner: Stop signal received during pause")
				return nil
			case <-time.After(100 * time.Millisecond):
				continue
			}
		}

		executed := gr.ExecuteTurn()

		if executed {
			turnCount++
		} else {
			if gr.game.IsGameOver() {
				logging.Debug(logging.TagGame, "Runner: Game ended during ExecuteTurn")
				break
			}
			select {
			case <-gr.stopChan:
				logging.Debug(logging.TagGame, "Runner: Stop signal received while waiting for turn")
				return nil
			case <-time.After(100 * time.Millisecond):
				continue
			}
		}

		// Optional delay for visualization is handled in ExecuteTurn
	}

	if turnCount >= gr.maxTurns {
		logging.Debug(logging.TagGame, "Runner: Game ended: Maximum turns reached")
		return gr.calculateWinnerOnMaxTurnsExceeded()
	}

	if gr.locker != nil {
		gr.locker.Lock()
		defer gr.locker.Unlock()
		return gr.game.GetWinner()
	}

	return gr.game.GetWinner()
}

func (gr *Runner) calculateWinnerOnMaxTurnsExceeded() *engine.Player {
	if gr.locker != nil {
		gr.locker.Lock()
		defer gr.locker.Unlock()
	}
	if float64(gr.game.Players[0].GetPieceScore())/float64(gr.game.Players[1].GetPieceScore()) > 1.15 {
		gr.game.SetWinner(gr.game.Players[0], WinCauseMaxTurns)
	} else if float64(gr.game.Players[1].GetPieceScore())/float64(gr.game.Players[0].GetPieceScore()) > 1.15 {
		gr.game.SetWinner(gr.game.Players[1], WinCauseMaxTurns)
	}
	return gr.game.GetWinner()
}

// ExecuteTurn executes a single turn. Returns false if waiting for human input.
func (gr *Runner) ExecuteTurn() bool {
	return gr.executeTurn(false)
}

func (gr *Runner) executeTurn(ignorePause bool) bool {
	if gr.locker != nil {
		gr.locker.Lock()
	}

	if gr.game.IsGameOver() {
		if gr.locker != nil {
			gr.locker.Unlock()
		}
		return false
	}

	if !ignorePause && gr.paused {
		if gr.locker != nil {
			gr.locker.Unlock()
		}
		return false
	}

	controller := gr.game.GetCurrentController()

	// Human controller - wait for input or handle move
	if controller.GetControllerType() == engine.HumanController {
		humanController, ok := controller.(*engine.HumanPlayerController)
		if !ok || !humanController.HasPendingMove() {
			if !gr.waitingForHumanInput {
				gr.waitingForHumanInput = true
			}
			if gr.locker != nil {
				gr.locker.Unlock()
			}
			return false
		}

		move := humanController.GetPendingMove()
		if move == nil {
			if gr.locker != nil {
				gr.locker.Unlock()
			}
			return false
		}

		piece := gr.game.Board.GetPieceAt(move.GetFrom())
		if piece == nil {
			if gr.locker != nil {
				gr.locker.Unlock()
			}
			return false
		}

		gr.game.MakeMove(move, piece)
		gr.waitingForHumanInput = false

		if gr.locker != nil {
			gr.locker.Unlock()
		}

		if gr.onMoveExecuted != nil {
			gr.onMoveExecuted()
		}
		return true
	}

	// AI controller - make move
	// Clone board for thread-safe calculation
	boardCopy := gr.game.Board.Clone()
	currentPlayer := gr.game.CurrentPlayer

	// Release lock while AI thinks -> avoid blocking other goroutines
	if gr.locker != nil {
		gr.locker.Unlock()
	}

	start := time.Now()
	move := controller.MakeMove(boardCopy)
	elapsed := time.Since(start)

	// Pacing delay (outside the lock)
	if !ignorePause && gr.turnDelay > 0 {
		delay := gr.turnDelay
		if gr.turnDelay >= 100*time.Millisecond {
			// #nosec G404 - weak random is sufficient for turn delay pacing
			delay = time.Duration(float64(gr.turnDelay)*0.8 + float64(rand.Intn(int(gr.turnDelay)))*0.4)
		}
		sleepTime := delay - elapsed
		if sleepTime > 0 {
			select {
			case <-gr.stopChan:
				return false
			case <-time.After(sleepTime):
				// Continue
			}
		}

		// Re-check pause after delay
		if !ignorePause && gr.paused {
			return false
		}
	}

	// Re-acquire lock to apply the move
	if gr.locker != nil {
		gr.locker.Lock()
	}

	// Re-verify that the game is still on and it's still this AI's turn
	if gr.game.IsGameOver() || gr.game.CurrentPlayer != currentPlayer {
		if gr.locker != nil {
			gr.locker.Unlock()
		}
		return false
	}

	// Check if game was stopped while AI was thinking
	select {
	case <-gr.stopChan:
		if gr.locker != nil {
			gr.locker.Unlock()
		}
		return false
	default:
	}

	piece := gr.game.Board.GetPieceAt(move.GetFrom())
	if piece == nil || piece.GetOwner() != gr.game.CurrentPlayer {
		opponent := gr.getOpponent(gr.game.CurrentPlayer)
		gr.game.SetWinner(opponent, WinCauseNoMovablePieces)
		if gr.locker != nil {
			gr.locker.Unlock()
		}
		return false
	}

	if !gr.game.Board.IsValidMove(&move) {
		opponent := gr.getOpponent(gr.game.CurrentPlayer)
		gr.game.SetWinner(opponent, WinCauseNoMovablePieces)
		if gr.locker != nil {
			gr.locker.Unlock()
		}
		return false
	}

	gr.game.MakeMove(&move, piece)

	if gr.locker != nil {
		gr.locker.Unlock()
	}

	if gr.onMoveExecuted != nil {
		gr.onMoveExecuted()
	}
	return true
}

// getOpponent returns the opponent of the given player
func (gr *Runner) getOpponent(player *engine.Player) *engine.Player {
	if gr.game.Players[0] == player {
		return gr.game.Players[1]
	}
	return gr.game.Players[0]
}

// IsWaitingForInput returns true if the game is waiting for human input
func (gr *Runner) IsWaitingForInput() bool {
	if gr.locker != nil {
		gr.locker.Lock()
		defer gr.locker.Unlock()
	}
	return gr.isWaitingForInput()
}

func (gr *Runner) isWaitingForInput() bool {
	return gr.waitingForHumanInput
}

// DebugSetWaitingForInput sets the waiting for human input flag to the given value.
// This is for debugging (& testing) purposes only and should not be used in production code.
func (gr *Runner) DebugSetWaitingForInput(value bool) {
	gr.waitingForHumanInput = value
}

// GetGame returns the underlying game
func (gr *Runner) GetGame() *Game {
	return gr.game
}

// SubmitHumanMove allows external code to submit a human player's move
func (gr *Runner) SubmitHumanMove(move engine.Move) error {
	if gr.locker != nil {
		gr.locker.Lock()
	}

	if !gr.waitingForHumanInput {
		if gr.locker != nil {
			gr.locker.Unlock()
		}
		return fmt.Errorf("not waiting for input")
	}

	controller := gr.game.GetCurrentController()
	if controller.GetControllerType() != engine.HumanController {
		if gr.locker != nil {
			gr.locker.Unlock()
		}
		return fmt.Errorf("current player is not human")
	}

	humanController, ok := controller.(*engine.HumanPlayerController)
	if !ok {
		if gr.locker != nil {
			gr.locker.Unlock()
		}
		return fmt.Errorf("invalid controller type")
	}

	// Check that the move's player matches the current player
	if move.GetPlayer() != gr.game.CurrentPlayer {
		if gr.locker != nil {
			gr.locker.Unlock()
		}
		return fmt.Errorf("move player does not match current player")
	}

	if !gr.game.Board.IsValidMove(&move) {
		if gr.locker != nil {
			gr.locker.Unlock()
		}
		return fmt.Errorf("invalid move")
	}

	humanController.SetPendingMove(move)

	if gr.locker != nil {
		gr.locker.Unlock()
	}

	gr.ExecuteTurn()

	return nil
}

// Pause pauses the game runner
func (gr *Runner) Pause() {
	if gr.locker != nil {
		gr.locker.Lock()
		defer gr.locker.Unlock()
	}
	gr.setPaused(true)
}

func (gr *Runner) setPaused(paused bool) {
	gr.paused = paused
}

// Unpause unpauses the game runner
func (gr *Runner) Unpause() {
	if gr.locker != nil {
		gr.locker.Lock()
		defer gr.locker.Unlock()
	}
	gr.setPaused(false)
}

// SetTurnDelay sets the delay between AI turns
func (gr *Runner) SetTurnDelay(delay time.Duration) {
	gr.turnDelay = delay
}

// Step executes a single turn even if the game is paused
func (gr *Runner) Step() bool {
	return gr.executeTurn(true)
}

// IsPaused returns whether the game runner is paused
func (gr *Runner) IsPaused() bool {
	if gr.locker != nil {
		gr.locker.Lock()
		defer gr.locker.Unlock()
	}
	return gr.isPaused()
}

func (gr *Runner) isPaused() bool {
	return gr.paused
}
