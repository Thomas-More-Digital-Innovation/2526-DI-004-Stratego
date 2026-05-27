package game

import (
	"digital-innovation/gostrategy/internal/logging"
	"math/rand"
	"time"
)

// RunToCompletion runs the game until it's over (for AI vs AI)
// Winner can be nil when max turns are reached and both AIs have a similar piece count
func (gr *Runner) RunToCompletion() *Player {
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

		if isGameOver || turnCount >= gr.maxTurns || gr.IsStopped() {
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
			if gr.isGameOverSafe() {
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

func (gr *Runner) calculateWinnerOnMaxTurnsExceeded() *Player {
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

func (gr *Runner) isGameOverSafe() bool {
	if gr.locker != nil {
		gr.locker.Lock()
		defer gr.locker.Unlock()
	}
	return gr.game.IsGameOver()
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
	if controller.GetControllerType() == HumanController {
		humanController, ok := controller.(*HumanPlayerController)
		if !ok || !humanController.HasPendingMove() {
			gr.stateMu.Lock()
			if !gr.waitingForHumanInput {
				gr.waitingForHumanInput = true
			}
			gr.stateMu.Unlock()

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
		gr.stateMu.Lock()
		gr.waitingForHumanInput = false
		gr.stateMu.Unlock()

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

// Step executes a single turn even if the game is paused
func (gr *Runner) Step() bool {
	return gr.executeTurn(true)
}
