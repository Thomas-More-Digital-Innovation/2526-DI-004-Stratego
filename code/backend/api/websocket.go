package api

import (
	"digital-innovation/stratego/game"
	"digital-innovation/stratego/logging"
	"digital-innovation/stratego/utils"
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
			return true // Allow non-browser clients
		}

		// Get allowed origins from env
		allowedOrigins := utils.GetEnv("ALLOWED_ORIGINS", "")
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
func HandleWebSocket(w http.ResponseWriter, r *http.Request, session *game.GameSession, hub *WSHub, seatIndex int) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logging.ConnectionError(session.ID, "", 0, err)
		return
	}

	client := &WSClient{
		conn:      conn,
		send:      make(chan []byte, 256),
		session:   session,
		seatIndex: seatIndex,
		hub:       hub,
		Username:  "",
		UserID:    0,
	}

	// Try to get user info from the session/hub if available
	switch seatIndex {
	case 0:
		client.Username = session.Player1Username
		client.UserID = utils.GetIntSafe(session.Player1UserID)
	case 1:
		client.Username = session.Player2Username
		client.UserID = utils.GetIntSafe(session.Player2UserID)
	}

	hub.register <- client

	go client.writePump()
	go client.readPump()
}
