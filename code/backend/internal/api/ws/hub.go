// Package ws provides WebSocket functionality for the game server.
package ws

import (
	"digital-innovation/gostrategy/internal/game"
	"digital-innovation/gostrategy/internal/logging"
	"digital-innovation/gostrategy/internal/models"
	"sync"
	"time"
)

// Hub manages all WebSocket connections for a game
type Hub struct {
	clients       map[*Client]bool
	broadcast     chan []byte
	register      chan *Client
	unregister    chan *Client
	Session       *game.Session
	GameType      string
	mutex         sync.RWMutex
	cleanupTimer  *time.Timer
	timerMutex    sync.Mutex
	cleanupPeriod time.Duration
	OnCleanup     func() // Callback to remove session from GameServer
	stop          chan bool
	stopped       bool
}

// GetGameType returns the game type
func (h *Hub) GetGameType() string { return h.GameType }

// NewHub creates a new Hub instance for a session
func NewHub(session *game.Session, gameType string) *Hub {
	return &Hub{
		clients:       make(map[*Client]bool),
		broadcast:     make(chan []byte, 256),
		register:      make(chan *Client, 128),
		unregister:    make(chan *Client, 128),
		Session:       session,
		GameType:      gameType,
		cleanupPeriod: 1 * time.Minute, // 1 minute grace period for reconnection
		stop:          make(chan bool, 1),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	// Start an initial cleanup timer to prevent zombie sessions if no one ever connects
	// We give 2 minutes for the initial connection (or double the custom period)
	initialPeriod := 2 * time.Minute
	if h.cleanupPeriod != 1*time.Minute {
		initialPeriod = h.cleanupPeriod * 2
	}

	originalPeriod := h.cleanupPeriod
	h.cleanupPeriod = initialPeriod
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
				if h.cleanupPeriod == 1*time.Minute {
					if h.Session.GetGameState().IsGameOver {
						h.cleanupPeriod = 30 * time.Second
					} else {
						h.cleanupPeriod = 1 * time.Minute
					}
				}

				switch h.GameType {
				case models.AiVsAi:
					// AI vs AI game always gets a 30 second cleanup period unless shortened for tests
					if h.cleanupPeriod > 30*time.Second {
						h.cleanupPeriod = 30 * time.Second
					}
					logging.Debug(logging.TagWeb, "All clients disconnected from AI vs AI game, starting cleanup timer (%v): %s", h.cleanupPeriod, h.Session.ID)
					h.startCleanupTimer()

				case models.HumanVsAi, models.HumanVsHuman:
					// Start cleanup timer with appropriate grace period
					logging.Debug(logging.TagWeb, "All clients disconnected from %s game, starting cleanup timer (%v): %s", h.GameType, h.cleanupPeriod, h.Session.ID)
					h.startCleanupTimer()
				}
			}

		case message := <-h.broadcast:
			h.mutex.RLock()
			// Copy clients to a slice to avoid holding RLock during potential mutations/sends
			clients := make([]*Client, 0, len(h.clients))
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
			logging.Debug(logging.TagWeb, "Stopping hub loop for game %s", h.Session.ID)
			return
		}
	}
}

// Stop stops the hub loop and closes all connections
func (h *Hub) Stop() {
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
func (h *Hub) IsStopped() bool {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return h.stopped
}

// ClientCount returns the number of connected clients
func (h *Hub) ClientCount() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return len(h.clients)
}

// IsUserConnected returns whether a specific user is currently connected to the hub
func (h *Hub) IsUserConnected(userID int) bool {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	for client := range h.clients {
		if client.UserID == userID {
			return true
		}
	}
	return false
}

// startCleanupTimer starts a timer to stop the game after the cleanup period
func (h *Hub) startCleanupTimer() {
	h.timerMutex.Lock()
	defer h.timerMutex.Unlock()

	if h.cleanupTimer != nil {
		h.cleanupTimer.Stop()
	}

	logging.Debug(logging.TagWeb, "Starting cleanup timer for %s game (will stop in %v): %s", h.GameType, h.cleanupPeriod, h.Session.ID)

	h.cleanupTimer = time.AfterFunc(h.cleanupPeriod, func() {
		logging.Debug(logging.TagWeb, "Cleanup timer expired for %s game, stopping and cleaning up: %s", h.GameType, h.Session.ID)
		h.Session.Stop("Inactivity timeout")
		h.Stop()
		if h.OnCleanup != nil {
			h.OnCleanup()
		}
	})
}

// StartGameOverCleanup starts a mandatory timer to stop the game after it's over
// This timer is NOT cancelled by reconnections, ensuring we eventually clean up
func (h *Hub) StartGameOverCleanup(duration time.Duration) {
	h.timerMutex.Lock()
	defer h.timerMutex.Unlock()

	// If there's already a cleanup timer, stop it to prefer the game-over one
	if h.cleanupTimer != nil {
		h.cleanupTimer.Stop()
	}

	logging.Debug(logging.TagWeb, "Starting mandatory game-over cleanup for game (will stop in %v): %s", duration, h.Session.ID)

	h.cleanupTimer = time.AfterFunc(duration, func() {
		logging.Debug(logging.TagWeb, "Mandatory game-over cleanup triggered for game: %s", h.Session.ID)
		h.Session.Stop("Game completed")
		h.Stop()
		if h.OnCleanup != nil {
			h.OnCleanup()
		}
	})
}

// cancelCleanupTimer cancels the cleanup timer if it's running
func (h *Hub) cancelCleanupTimer() {
	h.timerMutex.Lock()
	defer h.timerMutex.Unlock()

	if h.cleanupTimer != nil {
		wasActive := h.cleanupTimer.Stop()
		if wasActive {
			logging.Debug(logging.TagWeb, "Cleanup timer cancelled for %s game (client reconnected): %s", h.GameType, h.Session.ID)
		}
		h.cleanupTimer = nil
	}
}
