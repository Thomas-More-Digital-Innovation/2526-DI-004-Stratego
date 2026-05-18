package game

import (
	"digital-innovation/gostrategy/engine"
	"digital-innovation/gostrategy/logging"
	"digital-innovation/gostrategy/models"
	"digital-innovation/gostrategy/utils"
	"errors"
	"fmt"
	"sync"
	"time"
)

// Session manages a game that can be controlled via API
// Supports async gameplay for human players
type Session struct {
	ID                    string
	game                  *Game
	runner                *Runner
	mutex                 sync.RWMutex
	running               bool
	isSetupPhase          bool
	headless              bool
	player1Pieces         []*engine.Piece
	player2Pieces         []*engine.Piece
	doneChan              chan *engine.Player // Signals when game is complete
	stopChan              chan bool           // Signals to stop the game
	waitingForAnimation   bool
	animationCompleteChan chan bool
	moveNotifyChan        chan bool // Signals when a move has been executed
	moveAckChan           chan bool // Signals that move has been processed (for synchronization)
	aborted               bool      // Signals that game was manually stopped
	// User info for players (nil/empty if guest/AI)
	Player1UserID     *int
	Player1Username   string
	Player2UserID     *int
	Player2Username   string
	StartTime         time.Time
	setupCompleteChan chan bool
	setupTimer        *time.Timer
	setupWarningTimer *time.Timer
	setupExpiresAt    time.Time
	setupTimeout      time.Duration
	setupWarning      time.Duration
}

// SessionOption is a function that configures a Session
type SessionOption func(*Session)

// WithSetupTimeout set a custom setup phase timeout
func WithSetupTimeout(d time.Duration) SessionOption {
	return func(gs *Session) { gs.setupTimeout = d }
}

// WithSetupWarning sets a custom delay for the setup phase warning
// d is the duration from the start of the session until the warning fires
func WithSetupWarning(d time.Duration) SessionOption {
	return func(gs *Session) { gs.setupWarning = d }
}

// NewSession creates a new Session instance
func NewSession(id string, controller1, controller2 engine.PlayerController, opts ...SessionOption) *Session {
	g := NewGame(controller1, controller2)

	// Generate initial piece setups for both players
	player1Pieces := RandomSetup(g.Players[0])
	player2Pieces := RandomSetup(g.Players[1])

	session := &Session{
		ID:                    id,
		game:                  g,
		runner:                NewRunner(g, 1*time.Nanosecond, 1000),
		isSetupPhase:          true,
		player1Pieces:         player1Pieces,
		player2Pieces:         player2Pieces,
		doneChan:              make(chan *engine.Player, 1),
		stopChan:              make(chan bool, 1),
		animationCompleteChan: make(chan bool, 1),
		moveNotifyChan:        make(chan bool, 100),
		moveAckChan:           make(chan bool, 1),
		StartTime:             time.Now(),
		setupCompleteChan:     make(chan bool, 1),
		setupTimeout:          5 * time.Minute, // Default
		Player1Username:       g.Players[0].GetName(),
		Player2Username:       g.Players[1].GetName(),
	}

	// Apply options
	for _, opt := range opts {
		opt(session)
	}

	session.runner.stopChan = session.stopChan
	session.runner.locker = &session.mutex

	session.runner.SetMoveCallback(func() {
		session.NotifyMoveExecuted()
		// If we are headless, don't wait for UI acknowledgment
		if !session.headless {
			session.WaitForMoveAck(10 * time.Second)
		}
	})

	// Calculate logical warning delay if not explicitly set
	warningDelay := session.setupWarning
	if warningDelay == 0 {
		warningDelay = session.setupTimeout - (1 * time.Minute)
		if warningDelay <= 0 {
			warningDelay = session.setupTimeout / 2
		}
	}

	session.setupWarningTimer = time.AfterFunc(warningDelay, func() {
		session.mutex.Lock()
		if !session.isSetupPhase {
			session.mutex.Unlock()
			return
		}
		session.mutex.Unlock()
		session.NotifySetupUpdate()
	})

	// Add setup timeout
	session.setupExpiresAt = time.Now().Add(session.setupTimeout)
	session.setupTimer = time.AfterFunc(session.setupTimeout, func() {
		session.mutex.Lock()
		if !session.isSetupPhase {
			session.mutex.Unlock()
			return
		}
		session.mutex.Unlock()

		logging.Debug(logging.TagGame, "Session %s: Setup timeout reached. Starting game with current/random setups.", session.ID)
		_ = session.StartGameFromSetup(false)
	})

	return session
}

// Start begins the game loop in a goroutine
// Returns immediately, game runs asynchronously
func (gs *Session) Start() error {
	gs.mutex.Lock()
	if gs.running {
		gs.mutex.Unlock()
		return errors.New("game already running")
	}
	gs.running = true
	gs.mutex.Unlock()

	go func() {
		winner := gs.runner.RunToCompletion()

		winnerLabel := "Draw"
		loserLabel := "Draw"

		p1 := logging.FormatUser(gs.Player1Username, utils.GetIntSafe(gs.Player1UserID))
		p2 := logging.FormatUser(gs.Player2Username, utils.GetIntSafe(gs.Player2UserID))

		if winner != nil {
			if winner.GetID() == 0 {
				winnerLabel = p1
				loserLabel = p2
			} else {
				winnerLabel = p2
				loserLabel = p1
			}
		}

		logging.GameFinished(gs.ID, winnerLabel, loserLabel, gs.game.GetRound())
		gs.doneChan <- winner
		gs.mutex.Lock()
		gs.running = false
		gs.mutex.Unlock()
	}()

	return nil
}

// Stop forcefully stops the game session
func (gs *Session) Stop() {
	gs.mutex.Lock()
	if gs.aborted {
		gs.mutex.Unlock()
		return
	}
	gs.aborted = true
	gs.mutex.Unlock()

	// Always close stopChan to signal abortion to any listeners (like monitorGame)
	close(gs.stopChan)

	logging.GameAborted(gs.ID, "Manual stop requested", "", 0)

	if gs.setupTimer != nil {
		gs.setupTimer.Stop()
	}
	if gs.setupWarningTimer != nil {
		gs.setupWarningTimer.Stop()
	}
}

// IsAborted returns whether the game was aborted
func (gs *Session) IsAborted() bool {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	return gs.aborted
}

// Pause pauses the game session
func (gs *Session) Pause() {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	gs.runner.setPaused(true)
	logging.Debug(logging.TagGame, "Session %s: Paused", gs.ID)
}

// Unpause unpauses the game session
func (gs *Session) Unpause() {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	gs.runner.setPaused(false)
	logging.Debug(logging.TagGame, "Session %s: Unpaused", gs.ID)
}

// SetTurnDelay sets the delay between AI turns
func (gs *Session) SetTurnDelay(delay time.Duration) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	gs.runner.SetTurnDelay(delay)
}

// StepAI executes a single AI turn even if the game is paused
func (gs *Session) StepAI() bool {
	return gs.runner.Step()
}

// SubmitMove submits a move for a human player
// Returns error if move is invalid or not the player's turn
func (gs *Session) SubmitMove(playerID int, move engine.Move) error {
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
	if controller.GetControllerType() != engine.HumanController {
		return errors.New("current player is not human-controlled")
	}

	humanController, ok := controller.(*engine.HumanPlayerController)
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
func (gs *Session) GetBoard() *engine.Board {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	return gs.game.Board
}

// GetAvailableMoves returns valid moves for a piece at the given position
// It returns an error if the piece does not belong to the requesting player
func (gs *Session) GetAvailableMoves(playerID int, pos engine.Position) ([]engine.Move, error) {
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

// WaitForCompletion blocks until the game is complete and returns the winner
func (gs *Session) WaitForCompletion() *engine.Player {
	return <-gs.doneChan
}

// IsRunning returns whether the game is currently running
func (gs *Session) IsRunning() bool {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	return gs.running
}

// GetWinner returns the winner of the game
func (gs *Session) GetWinner() *engine.Player {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	return gs.game.GetWinner()
}

// GetWinCause returns the cause of the win
func (gs *Session) GetWinCause() WinCause {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	return gs.game.GetWinCause()
}

// SetWinner manually declares a winner and ends the game
func (gs *Session) SetWinner(player *engine.Player, cause WinCause) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	gs.game.SetWinner(player, cause)
}

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

// GetGame returns the game instance (for advanced access)
func (gs *Session) GetGame() *Game {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	return gs.game
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

func getPlayerIDOrNil(player *engine.Player) *int {
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

// NotifySetupUpdate signals that the setup phase state has changed (e.g. timer)
func (gs *Session) NotifySetupUpdate() {
	gs.NotifyMoveExecuted() // Currently reuse this to trigger broadcast
}

// IsSetupPhase returns whether the game is in setup phase
func (gs *Session) IsSetupPhase() bool {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	return gs.isSetupPhase
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

// IsHeadless returns whether the game is running in headless simulation mode
func (gs *Session) IsHeadless() bool {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	return gs.headless
}

// GetSetupPieces returns the setup pieces for a player
func (gs *Session) GetSetupPieces(playerID int) []*engine.Piece {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()

	if playerID == 0 {
		return gs.player1Pieces
	}
	return gs.player2Pieces
}

// SwapSetupPieces swaps two pieces in the setup
func (gs *Session) SwapSetupPieces(playerID int, pos1, pos2 engine.Position) error {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	if !gs.isSetupPhase {
		return errors.New("not in setup phase")
	}

	var pieces []*engine.Piece
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
func (gs *Session) positionToIndex(pos engine.Position, playerID int) int {
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

	var player *engine.Player
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

	var player *engine.Player
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

// IsAbortedChan returns the channel that signals if the game was aborted
func (gs *Session) IsAbortedChan() <-chan bool {
	return gs.stopChan
}

// GetMoveNotifyChan returns the channel that signals when a move has been executed
func (gs *Session) GetMoveNotifyChan() <-chan bool {
	return gs.moveNotifyChan
}

// RLock locks the session mutex for reading
func (gs *Session) RLock() {
	gs.mutex.RLock()
}

// RUnlock unlocks the session mutex for reading
func (gs *Session) RUnlock() {
	gs.mutex.RUnlock()
}
