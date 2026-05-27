package game

import (
	"digital-innovation/gostrategy/internal/models"
	"errors"
	"time"
)

func getPlayerIDOrNil(player *Player) *int {
	if player == nil {
		return nil
	}
	id := player.GetID()
	return &id
}

// SetPlayer1Associate associates a user with Player 1 slot
func (gs *Session) SetPlayer1Associate(userID int, username string) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	gs.Player1UserID = &userID
	gs.Player1Username = username
	gs.game.Players[0].SetName(username)
}

// SetPlayer2Associate associates a user with Player 2 slot
func (gs *Session) SetPlayer2Associate(userID int, username string) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	gs.Player2UserID = &userID
	gs.Player2Username = username
	gs.game.Players[1].SetName(username)
}

// GetPlayerIDs returns the user IDs associated with both players
func (gs *Session) GetPlayerIDs() (p1 *int, p2 *int) {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	return gs.Player1UserID, gs.Player2UserID
}

// GetGameState returns current game state (for API responses)
func (gs *Session) GetGameState() models.GameState {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()

	return models.GameState{
		Round:              gs.game.GetRound(),
		CurrentPlayerID:    gs.game.CurrentPlayer.GetID(),
		CurrentPlayerName:  gs.game.CurrentPlayer.GetName(),
		IsGameOver:         gs.game.IsGameOver(),
		WinnerID:           getPlayerIDOrNil(gs.game.GetWinner()),
		Player1Score:       gs.game.Players[0].GetPieceScore(),
		Player2Score:       gs.game.Players[1].GetPieceScore(),
		WaitingForInput:    gs.runner.isWaitingForInput(),
		Paused:             gs.runner.isPaused(),
		MoveCount:          len(gs.game.MoveHistory),
		Player1AlivePieces: len(gs.game.Players[0].GetAlivePieces()),
		Player2AlivePieces: len(gs.game.Players[1].GetAlivePieces()),
		IsSetupPhase:       gs.isSetupPhase,
		Headless:           gs.headless,
		SetupRemainingSecs: gs.getSetupRemainingSecs(),
		Player1Username:    gs.Player1Username,
		Player2Username:    gs.Player2Username,
	}
}

func (gs *Session) getSetupRemainingSecs() int {
	if !gs.isSetupPhase {
		return 0
	}
	remaining := time.Until(gs.setupExpiresAt)
	if remaining < 0 {
		return 0
	}
	return int(remaining.Seconds())
}

// GetGameSummary returns a lightweight summary of the game state
func (gs *Session) GetGameSummary(gameType string) models.GameSummary {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()

	return models.GameSummary{
		GameID:          gs.ID,
		Round:           gs.game.GetRound(),
		IsRunning:       gs.running,
		IsGameOver:      gs.game.IsGameOver(),
		IsSetupPhase:    gs.isSetupPhase,
		GameType:        gameType,
		Player1Username: gs.Player1Username,
		Player2Username: gs.Player2Username,
	}
}

// GetBoard returns the current board state
func (gs *Session) GetBoard() *Board {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	return gs.game.Board
}

// GetAvailableMoves returns valid moves for a piece at the given position
// It returns an error if the piece does not belong to the requesting player
func (gs *Session) GetAvailableMoves(playerID int, pos Position) ([]Move, error) {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()

	piece := gs.game.Board.GetPieceAt(pos)
	if piece == nil {
		return nil, errors.New("no piece at the given position")
	}

	if piece.GetOwner().GetID() != playerID {
		return nil, errors.New("you can only request moves for your own pieces")
	}

	return gs.game.Board.ListMoves(pos)
}

// IsHeadless returns whether the game is running in headless simulation mode
func (gs *Session) IsHeadless() bool {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	return gs.headless
}

// GetSetupPieces returns the setup pieces for a player
func (gs *Session) GetSetupPieces(playerID int) []*Piece {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()

	if playerID == 0 {
		return gs.player1Pieces
	}
	return gs.player2Pieces
}

// IsAbortedChan returns the channel that signals if the game was aborted
func (gs *Session) IsAbortedChan() <-chan bool {
	return gs.stopChan
}

// GetMoveNotifyChan returns the channel that signals when a move has been executed
func (gs *Session) GetMoveNotifyChan() <-chan bool {
	return gs.moveNotifyChan
}
