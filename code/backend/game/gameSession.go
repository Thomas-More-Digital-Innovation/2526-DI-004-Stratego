package game

import (
	"digital-innovation/stratego/engine"
	"errors"
	"log"
	"sync"
)

// GameSession manages a game that can be controlled via API
// Supports async gameplay for human players
type GameSession struct {
	ID       string
	game     *Game
	runner   *GameRunner
	mutex    sync.RWMutex
	running  bool
	doneChan chan *engine.Player // Signals when game is complete
}

func NewGameSession(id string, controller1, controller2 engine.PlayerController) *GameSession {
	g := NewGame(controller1, controller2)

	return &GameSession{
		ID:       id,
		game:     g,
		runner:   NewGameRunner(g, 0, 1000),
		doneChan: make(chan *engine.Player, 1),
	}
}

// Start begins the game loop in a goroutine
// Returns immediately, game runs asynchronously
func (gs *GameSession) Start() error {
	gs.mutex.Lock()
	if gs.running {
		gs.mutex.Unlock()
		return errors.New("game already running")
	}
	gs.running = true
	gs.mutex.Unlock()

	log.Printf("GameSession %s: Starting game loop", gs.ID)

	go func() {
		winner := gs.runner.RunToCompletion()
		log.Printf("GameSession %s: Game finished, winner=%v", gs.ID, winner)
		gs.doneChan <- winner
		gs.mutex.Lock()
		gs.running = false
		gs.mutex.Unlock()
	}()

	return nil
}

// SubmitMove submits a move for a human player
// Returns error if move is invalid or not the player's turn
func (gs *GameSession) SubmitMove(playerID int, move engine.Move) error {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()

	log.Printf("SubmitMove called: gameID=%s, playerID=%d, running=%v, currentPlayerID=%d, isGameOver=%v",
		gs.ID, playerID, gs.running, gs.game.CurrentPlayer.GetID(), gs.game.IsGameOver())

	if !gs.running {
		return errors.New("game not running")
	}

	// Verify it's the correct player's turn
	if gs.game.CurrentPlayer.GetID() != playerID {
		return errors.New("not your turn")
	}

	// Verify current controller is human
	controller := gs.game.GetCurrentController()
	if controller.GetControllerType() != engine.HumanController {
		return errors.New("current player is not human-controlled")
	}

	// Submit move to the controller
	humanController, ok := controller.(*engine.HumanPlayerController)
	if !ok {
		return errors.New("failed to cast to human controller")
	}

	humanController.SetPendingMove(move)
	return nil
}

// GetGameState returns current game state (for API responses)
func (gs *GameSession) GetGameState() GameState {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()

	return GameState{
		Round:              gs.game.GetRound(),
		CurrentPlayerID:    gs.game.CurrentPlayer.GetID(),
		CurrentPlayerName:  gs.game.CurrentPlayer.GetName(),
		IsGameOver:         gs.game.IsGameOver(),
		WinnerID:           getPlayerIDOrNil(gs.game.GetWinner()),
		Player1Score:       gs.game.Players[0].GetPieceScore(),
		Player2Score:       gs.game.Players[1].GetPieceScore(),
		WaitingForInput:    gs.runner.IsWaitingForInput(),
		MoveCount:          len(gs.game.MoveHistory),
		Player1AlivePieces: len(gs.game.Players[0].GetAlivePieces()),
		Player2AlivePieces: len(gs.game.Players[1].GetAlivePieces()),
	}
}

// GetBoard returns the current board state
func (gs *GameSession) GetBoard() *engine.Board {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	return gs.game.Board
}

// GetAvailableMoves returns valid moves for a piece at the given position
func (gs *GameSession) GetAvailableMoves(pos engine.Position) ([]engine.Move, error) {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()

	return gs.game.Board.ListMoves(pos)
}

// WaitForCompletion blocks until the game is complete and returns the winner
func (gs *GameSession) WaitForCompletion() *engine.Player {
	return <-gs.doneChan
}

// IsRunning returns whether the game is currently running
func (gs *GameSession) IsRunning() bool {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	return gs.running
}

// GetWinner returns the winner of the game
func (gs *GameSession) GetWinner() *engine.Player {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	return gs.game.GetWinner()
}

// GetWinCause returns the cause of the win
func (gs *GameSession) GetWinCause() WinCause {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	return gs.game.GetWinCause()
}

// GetLastCombat returns the last combat result if any
func (gs *GameSession) GetLastCombat() *CombatResult {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	return gs.game.GetLastCombat()
}

// HideCombatPieces hides the pieces from the last combat
func (gs *GameSession) HideCombatPieces() {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	gs.game.HideCombatPieces()
}

// GetGame returns the game instance (for advanced access)
func (gs *GameSession) GetGame() *Game {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	return gs.game
}

// GameState represents the current state of a game (for API responses)
type GameState struct {
	Round              int    `json:"round"`
	CurrentPlayerID    int    `json:"currentPlayerId"`
	CurrentPlayerName  string `json:"currentPlayerName"`
	IsGameOver         bool   `json:"isGameOver"`
	WinnerID           *int   `json:"winnerId,omitempty"`
	Player1Score       int    `json:"player1Score"`
	Player2Score       int    `json:"player2Score"`
	WaitingForInput    bool   `json:"waitingForInput"`
	MoveCount          int    `json:"moveCount"`
	Player1AlivePieces int    `json:"player1AlivePieces"`
	Player2AlivePieces int    `json:"player2AlivePieces"`
}

func getPlayerIDOrNil(player *engine.Player) *int {
	if player == nil {
		return nil
	}
	id := player.GetID()
	return &id
}
