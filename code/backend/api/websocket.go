package api

import (
	"digital-innovation/stratego/game"
	"digital-innovation/stratego/utils"
	"log"
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

		log.Printf("WebSocket: Rejected connection from unauthorized origin: %s", origin)
		return false
	},
}

// HandleWebSocket handles WebSocket connections
func HandleWebSocket(w http.ResponseWriter, r *http.Request, session *game.GameSession, hub *WSHub, seatIndex int) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &WSClient{
		conn:      conn,
		send:      make(chan []byte, 256),
		session:   session,
		seatIndex: seatIndex,
		hub:       hub,
	}

	hub.register <- client

	go client.writePump()
	go client.readPump()
}
