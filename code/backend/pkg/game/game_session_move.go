package game

import (
	"digital-innovation/gostrategy/internal/logging"
	"errors"
	"time"
)

// SubmitMove submits a move for a human player
// Returns error if move is invalid or not the player's turn
func (gs *Session) SubmitMove(playerID int, move Move) error {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()

	logging.Debug(logging.TagGame, "SubmitMove called: gameID=%s, playerID=%d, running=%v, currentPlayerID=%d, isGameOver=%v",
		gs.ID, playerID, gs.running, gs.game.CurrentPlayer.GetID(), gs.game.IsGameOver())

	if !gs.running {
		return errors.New("game not running")
	}

	if gs.game.CurrentPlayer.GetID() != playerID {
		return errors.New("not your turn")
	}

	controller := gs.game.GetCurrentController()
	if controller.GetControllerType() != HumanController {
		return errors.New("current player is not human-controlled")
	}

	humanController, ok := controller.(*HumanPlayerController)
	if !ok {
		return errors.New("failed to cast to human controller")
	}

	// Validate piece ownership
	piece := gs.game.Board.GetPieceAt(move.GetFrom())
	if piece == nil {
		return errors.New("no piece at source position")
	}
	if piece.GetOwner() == nil || piece.GetOwner().GetID() != playerID {
		return errors.New("piece at source position does not belong to current player")
	}

	// Validate move legality
	validMoves, err := gs.game.Board.ListMoves(move.GetFrom())
	if err != nil {
		return err
	}

	isLegal := false
	for _, m := range validMoves {
		if m.GetTo() == move.GetTo() {
			isLegal = true
			break
		}
	}

	if !isLegal {
		return errors.New("illegal move for this piece")
	}

	humanController.SetPendingMove(move)
	return nil
}

// NotifyMoveExecuted signals that a move has been executed
func (gs *Session) NotifyMoveExecuted() {
	select {
	case gs.moveNotifyChan <- true:
	default:
		logging.Debug(logging.TagGame, "Session %s: Move notification channel full", gs.ID)
	}
}

// WaitForMoveNotification waits for a move to be executed
func (gs *Session) WaitForMoveNotification(timeout time.Duration) bool {
	select {
	case <-gs.moveNotifyChan:
		return true
	case <-time.After(timeout):
		return false
	}
}

// AckMoveProcessed signals that the move has been processed by the monitor
func (gs *Session) AckMoveProcessed() {
	select {
	case gs.moveAckChan <- true:
	default:
		// Channel full
	}
}

// WaitForMoveAck waits for move to be acknowledged as processed
func (gs *Session) WaitForMoveAck(timeout time.Duration) bool {
	select {
	case <-gs.moveAckChan:
		return true
	case <-time.After(timeout):
		return false
	}
}
