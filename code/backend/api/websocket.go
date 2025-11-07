package api

import (
	"digital-innovation/stratego/game"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for development
		// TODO: Restrict in production
		return true
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
