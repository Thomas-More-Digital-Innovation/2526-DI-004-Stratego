package api

import (
	"digital-innovation/stratego/game"
	"digital-innovation/stratego/models"
	"log"
	"sync"
	"time"
)

// WSHub manages all WebSocket connections for a game
type WSHub struct {
	clients       map[*WSClient]bool
	broadcast     chan []byte
	register      chan *WSClient
	unregister    chan *WSClient
	session       *game.GameSession
	gameType      string
	mutex         sync.RWMutex
	cleanupTimer  *time.Timer
	timerMutex    sync.Mutex
	cleanupPeriod time.Duration
}

func NewWSHub(session *game.GameSession, gameType string) *WSHub {
	return &WSHub{
		clients:       make(map[*WSClient]bool),
		broadcast:     make(chan []byte, 256),
		register:      make(chan *WSClient),
		unregister:    make(chan *WSClient),
		session:       session,
		gameType:      gameType,
		cleanupPeriod: 1 * time.Minute, // 1 minute grace period for reconnection
	}
}

// Run starts the hub's main loop
// If user disconnects, the hub will stop the game after 1 minute, if human is playing
func (h *WSHub) Run() {
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
				close(client.send)
			}
			clientCount := len(h.clients)
			h.mutex.Unlock()

			if clientCount == 0 {
				switch h.gameType {
				case models.AiVsAi:
					// Stop AI vs AI games immediately - no point running without observers
					log.Printf("WSHub: All clients disconnected from AI vs AI game, stopping game immediately")
					h.session.Stop()

				case models.HumanVsAi:
					// Start cleanup timer for Human vs AI - allow reconnection grace period
					log.Printf("WSHub: All clients disconnected from Human vs AI game, starting cleanup timer")
					h.startCleanupTimer()

				case models.HumanVsHuman:
					// For Human vs Human, start timer to allow reconnection if both players leave
					log.Printf("WSHub: All clients disconnected from Human vs Human game, starting cleanup timer")
					h.startCleanupTimer()
				}
			}

		case message := <-h.broadcast:
			h.mutex.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					// Client is slow or disconnected
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mutex.RUnlock()
		}
	}
}

// startCleanupTimer starts a timer to stop the game after the cleanup period
func (h *WSHub) startCleanupTimer() {
	h.timerMutex.Lock()
	defer h.timerMutex.Unlock()

	if h.cleanupTimer != nil {
		h.cleanupTimer.Stop()
	}

	log.Printf("WSHub: Starting cleanup timer for %s game (will stop in %v)", h.gameType, h.cleanupPeriod)

	h.cleanupTimer = time.AfterFunc(h.cleanupPeriod, func() {
		log.Printf("WSHub: Cleanup timer expired for %s game, stopping game", h.gameType)
		h.session.Stop()
	})
}

// cancelCleanupTimer cancels the cleanup timer if it's running
func (h *WSHub) cancelCleanupTimer() {
	h.timerMutex.Lock()
	defer h.timerMutex.Unlock()

	if h.cleanupTimer != nil {
		wasActive := h.cleanupTimer.Stop()
		if wasActive {
			log.Printf("WSHub: Cleanup timer cancelled for %s game (client reconnected)", h.gameType)
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
				dto := PieceToDTO(piece, -1) // Hide pieces
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
