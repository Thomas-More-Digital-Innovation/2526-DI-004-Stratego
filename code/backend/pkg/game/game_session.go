package game

import (
	"digital-innovation/gostrategy/internal/logging"
	"digital-innovation/gostrategy/internal/utils"
	"errors"
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
	player1Pieces         []*Piece
	player2Pieces         []*Piece
	doneChan              chan *Player // Signals when game is complete
	stopChan              chan bool    // Signals to stop the game
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
func NewSession(id string, controller1, controller2 PlayerController, opts ...SessionOption) *Session {
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
		doneChan:              make(chan *Player, 1),
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

// WaitForCompletion blocks until the game is complete and returns the winner
func (gs *Session) WaitForCompletion() *Player {
	return <-gs.doneChan
}

// IsRunning returns whether the game is currently running
func (gs *Session) IsRunning() bool {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	return gs.running
}

// GetWinner returns the winner of the game
func (gs *Session) GetWinner() *Player {
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
func (gs *Session) SetWinner(player *Player, cause WinCause) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	gs.game.SetWinner(player, cause)
}

// GetGame returns the game instance (for advanced access)
func (gs *Session) GetGame() *Game {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	return gs.game
}

// RLock locks the session mutex for reading
func (gs *Session) RLock() {
	gs.mutex.RLock()
}

// RUnlock unlocks the session mutex for reading
func (gs *Session) RUnlock() {
	gs.mutex.RUnlock()
}
