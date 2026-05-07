package api

import (
	"digital-innovation/gostrategy/game"
	"digital-innovation/gostrategy/logging"
	"digital-innovation/gostrategy/models"
	"sync"
	"time"
)

// WSHub manages all WebSocket connections for a game
type WSHub struct {
	clients       map[*WSClient]bool
	broadcast     chan []byte
	register      chan *WSClient
	unregister    chan *WSClient
	session       *game.Session
	gameType      string
	mutex         sync.RWMutex
	cleanupTimer  *time.Timer
	timerMutex    sync.Mutex
	cleanupPeriod time.Duration
	OnCleanup     func() // Callback to remove session from GameServer
	stop          chan bool
	stopped       bool
}

// NewWSHub creates a new WSHub instance for a session
func NewWSHub(session *game.Session, gameType string) *WSHub {
	return &WSHub{
		clients:       make(map[*WSClient]bool),
		broadcast:     make(chan []byte, 256),
		register:      make(chan *WSClient, 128),
		unregister:    make(chan *WSClient, 128),
		session:       session,
		gameType:      gameType,
		cleanupPeriod: 1 * time.Minute, // 1 minute grace period for reconnection
		stop:          make(chan bool, 1),
	}
}

// Run starts the hub's main loop
func (h *WSHub) Run() {
	// Start an initial cleanup timer to prevent zombie sessions if no one ever connects
	// We give 2 minutes for the initial connection
	originalPeriod := h.cleanupPeriod
	h.cleanupPeriod = 2 * time.Minute
	h.startCleanupTimer()
	h.cleanupPeriod = originalPeriod

	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			clientCount := len(h.clients)
			h.mutex.Unlock()

			if clientCount > 0 {
				h.cancelCleanupTimer()
			}

			go h.sendGameState(client)

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Close()
			}
			clientCount := len(h.clients)
			h.mutex.Unlock()

			if clientCount == 0 {
				if h.session.GetGameState().IsGameOver {
					h.cleanupPeriod = 30 * time.Second
				} else {
					h.cleanupPeriod = 1 * time.Minute
				}

				switch h.gameType {
				case models.AiVsAi:
					// Stop AI vs AI games immediately - no point running without observers
					logging.Debug(logging.TagWeb, "All clients disconnected from AI vs AI game, stopping game immediately: %s", h.session.ID)
					h.session.Stop()
					h.Stop()
					if h.OnCleanup != nil {
						h.OnCleanup()
					}

				case models.HumanVsAi, models.HumanVsHuman:
					// Start cleanup timer with appropriate grace period
					logging.Debug(logging.TagWeb, "All clients disconnected from %s game, starting cleanup timer (%v): %s", h.gameType, h.cleanupPeriod, h.session.ID)
					h.startCleanupTimer()
				}
			}

		case message := <-h.broadcast:
			h.mutex.RLock()
			// Copy clients to a slice to avoid holding RLock during potential mutations/sends
			clients := make([]*WSClient, 0, len(h.clients))
			for client := range h.clients {
				clients = append(clients, client)
			}
			h.mutex.RUnlock()

			for _, client := range clients {
				if !client.Send(message, 100*time.Millisecond) {
					// Client is slow or disconnected - unregister them properly
					logging.Debug(logging.TagWeb, "Dropping slow client: %s", client.Username)
					h.mutex.Lock()
					if _, ok := h.clients[client]; ok {
						delete(h.clients, client)
						client.Close()
					}
					clientCount := len(h.clients)
					h.mutex.Unlock()

					// If this was the last client, we need to trigger the cleanup logic
					// To avoid duplicating logic, we send to the unregister case which handles timers
					if clientCount == 0 {
						select {
						case h.unregister <- client:
						default:
						}
					}
				}
			}
		case <-h.stop:
			logging.Debug(logging.TagWeb, "Stopping hub loop for game %s", h.session.ID)
			return
		}
	}
}

// Stop stops the hub loop and closes all connections
func (h *WSHub) Stop() {
	h.mutex.Lock()
	if h.stopped {
		h.mutex.Unlock()
		return
	}
	h.stopped = true

	// Close all client connections to trigger their pump exits
	for client := range h.clients {
		if client.conn != nil {
			_ = client.conn.Close()
		}
	}
	h.mutex.Unlock()

	select {
	case h.stop <- true:
	default:
		// Already stopped or no one listening
	}
}

// IsStopped returns whether the hub is stopped
func (h *WSHub) IsStopped() bool {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return h.stopped
}

// startCleanupTimer starts a timer to stop the game after the cleanup period
func (h *WSHub) startCleanupTimer() {
	h.timerMutex.Lock()
	defer h.timerMutex.Unlock()

	if h.cleanupTimer != nil {
		h.cleanupTimer.Stop()
	}

	logging.Debug(logging.TagWeb, "Starting cleanup timer for %s game (will stop in %v): %s", h.gameType, h.cleanupPeriod, h.session.ID)

	h.cleanupTimer = time.AfterFunc(h.cleanupPeriod, func() {
		logging.Debug(logging.TagWeb, "Cleanup timer expired for %s game, stopping and cleaning up: %s", h.gameType, h.session.ID)
		h.session.Stop()
		h.Stop()
		if h.OnCleanup != nil {
			h.OnCleanup()
		}
	})
}

// StartGameOverCleanup starts a mandatory timer to stop the game after it's over
// This timer is NOT cancelled by reconnections, ensuring we eventually clean up
func (h *WSHub) StartGameOverCleanup(duration time.Duration) {
	h.timerMutex.Lock()
	defer h.timerMutex.Unlock()

	// If there's already a cleanup timer, stop it to prefer the game-over one
	if h.cleanupTimer != nil {
		h.cleanupTimer.Stop()
	}

	logging.Debug(logging.TagWeb, "Starting mandatory game-over cleanup for game (will stop in %v): %s", duration, h.session.ID)

	h.cleanupTimer = time.AfterFunc(duration, func() {
		logging.Debug(logging.TagWeb, "Mandatory game-over cleanup triggered for game: %s", h.session.ID)
		h.session.Stop()
		h.Stop()
		if h.OnCleanup != nil {
			h.OnCleanup()
		}
	})
}

// cancelCleanupTimer cancels the cleanup timer if it's running
func (h *WSHub) cancelCleanupTimer() {
	h.timerMutex.Lock()
	defer h.timerMutex.Unlock()

	if h.cleanupTimer != nil {
		wasActive := h.cleanupTimer.Stop()
		if wasActive {
			logging.Debug(logging.TagWeb, "Cleanup timer cancelled for %s game (client reconnected): %s", h.gameType, h.session.ID)
		}
		h.cleanupTimer = nil
	}
}

func (h *WSHub) setupBoard() BoardStateMessage {
	session := h.session

	boardDTO := make([][]PieceDTO, 10)
	for y := range 10 {
		boardDTO[y] = make([]PieceDTO, 10)
	}

	// Place player 1 pieces in setup area (rows 6-9)
	player1Pieces := session.GetSetupPieces(0)
	idx := 0
	for y := 6; y <= 9; y++ {
		for x := range 10 {
			if idx < len(player1Pieces) {
				piece := player1Pieces[idx]
				dto := PieceToDTO(piece, 0) // Player 0 can see their own pieces
				if h.gameType == models.AiVsAi {
					dto.Revealed = true // Force visibility for spectators during setup
				}
				dto.Position = PositionDTO{X: x, Y: y}
				boardDTO[y][x] = dto
				idx++
			}
		}
	}

	// Place player 2 pieces in setup area (rows 0-3)
	// Hide opponent pieces during setup
	player2Pieces := session.GetSetupPieces(1)
	idx = 0
	for y := 0; y <= 3; y++ {
		for x := range 10 {
			if idx < len(player2Pieces) {
				piece := player2Pieces[idx]
				viewerID := -1
				if h.gameType == models.AiVsAi {
					viewerID = 1 // Show all pieces in AI vs AI
				}
				dto := PieceToDTO(piece, viewerID) // Hide pieces
				if h.gameType == models.AiVsAi {
					dto.Revealed = true // Force visibility for spectators during setup
				}
				dto.Position = PositionDTO{X: x, Y: y}
				boardDTO[y][x] = dto
				idx++
			}
		}
	}

	return BoardStateMessage{
		Board:  boardDTO,
		Width:  10,
		Height: 10,
	}
}
