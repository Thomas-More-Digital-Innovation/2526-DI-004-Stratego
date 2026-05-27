package game

import (
	"digital-innovation/gostrategy/internal/logging"
	"errors"
	"fmt"
	"time"
)

// SwapSetupPieces swaps two pieces in the setup
func (gs *Session) SwapSetupPieces(playerID int, pos1, pos2 Position) error {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	if !gs.isSetupPhase {
		return errors.New("not in setup phase")
	}

	var pieces []*Piece
	switch playerID {
	case 0:
		pieces = gs.player1Pieces
	case 1:
		pieces = gs.player2Pieces
	default:
		return errors.New("invalid player ID")
	}

	// Calculate indices from positions (setup area is 4x10 = 40 pieces)
	idx1 := gs.positionToIndex(pos1, playerID)
	idx2 := gs.positionToIndex(pos2, playerID)

	if idx1 < 0 || idx1 >= len(pieces) || idx2 < 0 || idx2 >= len(pieces) {
		return fmt.Errorf("invalid positions for swap: %v, %v", pos1, pos2)
	}

	// Swap the pieces
	pieces[idx1], pieces[idx2] = pieces[idx2], pieces[idx1]

	logging.Debug(logging.TagGame, "Swapped pieces for player %d: index %d <-> %d", playerID, idx1, idx2)
	return nil
}

// positionToIndex converts a board position to piece array index
func (gs *Session) positionToIndex(pos Position, playerID int) int {
	var startRow, endRow int
	if playerID == 0 {
		startRow = 6
		endRow = 9
	} else {
		startRow = 0
		endRow = 3
	}

	// Check if position is in valid range
	if pos.Y < startRow || pos.Y > endRow || pos.X < 0 || pos.X >= 10 {
		return -1
	}

	// Calculate index
	rowOffset := pos.Y - startRow
	return rowOffset*10 + pos.X
}

// LoadSetup loads a predefined setup from binary data (40 bytes)
func (gs *Session) LoadSetup(playerID int, data []byte) error {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	if !gs.isSetupPhase {
		return errors.New("not in setup phase")
	}

	var player *Player
	switch playerID {
	case 0:
		player = gs.game.Players[0]
	case 1:
		player = gs.game.Players[1]
	default:
		return errors.New("invalid player ID")
	}

	pieces, err := ParseSetup(player, data)
	if err != nil {
		return fmt.Errorf("invalid setup data: %v", err)
	}

	switch playerID {
	case 0:
		gs.player1Pieces = pieces
	case 1:
		gs.player2Pieces = pieces
	}

	logging.Debug(logging.TagGame, "Loaded custom setup for player %d in game %s", playerID, gs.ID)
	return nil
}

// RandomizeSetup randomizes the setup for a player
func (gs *Session) RandomizeSetup(playerID int) error {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	if !gs.isSetupPhase {
		return errors.New("not in setup phase")
	}

	var player *Player
	switch playerID {
	case 0:
		player = gs.game.Players[0]
		gs.player1Pieces = RandomSetup(player)
	case 1:
		player = gs.game.Players[1]
		gs.player2Pieces = RandomSetup(player)
	default:
		return errors.New("invalid player ID")
	}

	logging.Debug(logging.TagGame, "Randomized setup for player %d in game %s", playerID, gs.ID)
	return nil
}

// StartGameFromSetup starts the game from setup phase
func (gs *Session) StartGameFromSetup(headless bool) error {
	gs.mutex.Lock()

	if !gs.isSetupPhase {
		gs.mutex.Unlock()
		return errors.New("not in setup phase")
	}

	// Set initial speed based on headless mode
	if headless {
		gs.runner.SetTurnDelay(0)
	} else {
		// Default pacing of 1 second for visible games
		gs.runner.SetTurnDelay(1 * time.Second)
	}

	// Place pieces on the board
	if err := SetupGame(gs.game, gs.player1Pieces, gs.player2Pieces); err != nil {
		gs.mutex.Unlock()
		return fmt.Errorf("failed to setup game: %v", err)
	}

	if gs.setupTimer != nil {
		gs.setupTimer.Stop()
	}
	if gs.setupWarningTimer != nil {
		gs.setupWarningTimer.Stop()
	}

	gs.game.InitialState = gs.game.GetInitialBoardState()

	// Exit setup phase BEFORE starting the game
	gs.isSetupPhase = false
	gs.headless = headless

	// Signal setup complete to any monitors
	select {
	case gs.setupCompleteChan <- true:
	default:
	}

	gs.mutex.Unlock()

	// Start the game (now outside the lock)
	if err := gs.Start(); err != nil {
		logging.Error("Error starting game "+gs.ID, err)
		return err
	}

	return nil
}

// GetSetupCompleteChan returns the channel that signals setup completion
func (gs *Session) GetSetupCompleteChan() <-chan bool {
	return gs.setupCompleteChan
}

// SetSetupPhaseComplete marks the setup phase as complete
func (gs *Session) SetSetupPhaseComplete() {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	if gs.setupTimer != nil {
		gs.setupTimer.Stop()
	}
	if gs.setupWarningTimer != nil {
		gs.setupWarningTimer.Stop()
	}

	gs.isSetupPhase = false
	select {
	case gs.setupCompleteChan <- true:
	default:
		// Already signaled or no one listening
	}
}

// IsSetupPhase returns whether the game is in setup phase
func (gs *Session) IsSetupPhase() bool {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	return gs.isSetupPhase
}

// NotifySetupUpdate signals that the setup phase state has changed (e.g. timer)
func (gs *Session) NotifySetupUpdate() {
	gs.NotifyMoveExecuted() // Currently reuse this to trigger broadcast
}
