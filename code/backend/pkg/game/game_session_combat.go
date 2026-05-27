package game

import (
	"digital-innovation/gostrategy/internal/logging"
	"digital-innovation/gostrategy/internal/models"
	"time"
)

// GetLastCombat returns the last combat result if any
func (gs *Session) GetLastCombat() *CombatResult {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	return gs.game.GetLastCombat()
}

// GetLastHistoricalMove returns the last historical move if any
func (gs *Session) GetLastHistoricalMove() *models.HistoricalMove {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	if len(gs.game.HistoricalHistory) == 0 {
		return nil
	}
	move := gs.game.HistoricalHistory[len(gs.game.HistoricalHistory)-1]
	return &move
}

// ClearLastCombat clears the last combat result
func (gs *Session) ClearLastCombat() {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	gs.game.ClearLastCombat()
}

// HideCombatPieces hides the pieces from the last combat
func (gs *Session) HideCombatPieces() {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	gs.game.HideCombatPieces()
}

// WaitForAnimationComplete blocks until animation is complete or timeout
func (gs *Session) WaitForAnimationComplete(timeout time.Duration) {
	gs.mutex.Lock()
	gs.waitingForAnimation = true
	gs.mutex.Unlock()

	select {
	case <-gs.animationCompleteChan:
		logging.Debug(logging.TagGame, "Session %s: Animation complete signal received", gs.ID)
	case <-time.After(timeout):
		logging.Debug(logging.TagGame, "Session %s: Animation timeout", gs.ID)
	}

	gs.mutex.Lock()
	gs.waitingForAnimation = false
	gs.mutex.Unlock()
}

// SignalAnimationComplete signals that the animation has completed
func (gs *Session) SignalAnimationComplete() {
	gs.mutex.RLock()
	waiting := gs.waitingForAnimation
	gs.mutex.RUnlock()

	if waiting {
		select {
		case gs.animationCompleteChan <- true:
			logging.Debug(logging.TagGame, "Session %s: Animation complete signal sent", gs.ID)
		default:
			logging.Debug(logging.TagGame, "Session %s: Animation complete channel full", gs.ID)
		}
	}
}

// IsWaitingForAnimation returns whether the session is waiting for animation
func (gs *Session) IsWaitingForAnimation() bool {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	return gs.waitingForAnimation
}
