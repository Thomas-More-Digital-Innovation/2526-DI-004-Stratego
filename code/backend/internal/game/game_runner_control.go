package game

import "fmt"

// IsWaitingForInput returns true if the game is waiting for human input
func (gr *Runner) IsWaitingForInput() bool {
	if gr.locker != nil {
		gr.locker.Lock()
		defer gr.locker.Unlock()
	}
	return gr.isWaitingForInput()
}

func (gr *Runner) isWaitingForInput() bool {
	gr.stateMu.RLock()
	defer gr.stateMu.RUnlock()
	return gr.waitingForHumanInput
}

// DebugSetWaitingForInput sets the waiting for human input flag to the given value.
// This is for debugging (& testing) purposes only and should not be used in production code.
func (gr *Runner) DebugSetWaitingForInput(value bool) {
	gr.stateMu.Lock()
	defer gr.stateMu.Unlock()
	gr.waitingForHumanInput = value
}

// Stop stops the game runner
func (gr *Runner) Stop() {
	gr.stateMu.Lock()
	gr.stopped = true
	gr.stateMu.Unlock()

	select {
	case gr.stopChan <- true:
	default:
		// Channel already has a stop signal
	}
}

// SubmitHumanMove allows external code to submit a human player's move
func (gr *Runner) SubmitHumanMove(move Move) error {
	if gr.locker != nil {
		gr.locker.Lock()
	}

	if !gr.isWaitingForInput() {
		if gr.locker != nil {
			gr.locker.Unlock()
		}
		return fmt.Errorf("not waiting for input")
	}

	controller := gr.game.GetCurrentController()
	if controller.GetControllerType() != HumanController {
		if gr.locker != nil {
			gr.locker.Unlock()
		}
		return fmt.Errorf("current player is not human")
	}

	humanController, ok := controller.(*HumanPlayerController)
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
	gr.stateMu.Lock()
	defer gr.stateMu.Unlock()
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

// IsPaused returns whether the game runner is paused
func (gr *Runner) IsPaused() bool {
	if gr.locker != nil {
		gr.locker.Lock()
		defer gr.locker.Unlock()
	}
	return gr.isPaused()
}

func (gr *Runner) isPaused() bool {
	gr.stateMu.RLock()
	defer gr.stateMu.RUnlock()
	return gr.paused
}

// IsStopped returns whether the game runner has been stopped
func (gr *Runner) IsStopped() bool {
	gr.stateMu.RLock()
	defer gr.stateMu.RUnlock()
	return gr.stopped
}
