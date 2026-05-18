// Package ws provides WebSocket functionality for the game server.
package ws

import (
	"digital-innovation/gostrategy/game"
	"digital-innovation/gostrategy/logging"
	"digital-innovation/gostrategy/utils"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" {
			// In production, we strictly require an Origin header for browser safety
			if utils.IsProduction() {
				logging.SecurityWarningWithIP("WebSocket: Rejected connection with missing Origin header", "", r.RemoteAddr)
				return false
			}
			return true // Allow non-browser clients in dev
		}

		// Get allowed origins from env
		allowedOrigins := utils.GetEnv("ALLOWED_ORIGINS", "")
		if allowedOrigins == "" && !utils.IsProduction() {
			return true // Default allow for local dev if not set
		}

		for _, allowed := range strings.Split(allowedOrigins, ",") {
			if origin == allowed {
				return true
			}
		}

		logging.SecurityWarningWithIP("WebSocket: Rejected connection from unauthorized origin", "Origin: "+origin, r.RemoteAddr)
		return false
	},
}

// HandleWebSocket handles WebSocket connections
func HandleWebSocket(w http.ResponseWriter, r *http.Request, session *game.Session, hub *Hub, seatIndex int, username string, userID int) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logging.ConnectionError(session.ID, "", 0, err)
		return
	}

	client := &Client{
		conn:      conn,
		send:      make(chan []byte, 256),
		stopChan:  make(chan struct{}),
		Session:   session,
		SeatIndex: seatIndex,
		Hub:       hub,
		Username:  username,
		UserID:    userID,
	}

	select {
	case hub.register <- client:
	default:
		// Hub is stopped or not accepting new clients
		_ = conn.Close()
		return
	}

	go client.writePump()
	go client.readPump()
}
