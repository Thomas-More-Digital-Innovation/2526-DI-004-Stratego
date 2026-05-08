package api

import (
	"digital-innovation/gostrategy/db"
	"digital-innovation/gostrategy/models"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func TestIntegration_WebSocket(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db.SetupTestDB(t)

	server := NewGameServer()
	// Add routes
	server.router.GET("/game/:gameID", server.HandleWebSocketConnection)

	ts := httptest.NewServer(server.router)
	defer ts.Close()

	// Convert http URL to ws URL
	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http")

	// 1. Create a game session manually in server
	gameID := "ws-test-game"
	handler, err := server.CreateGame(gameID, models.HumanVsHuman, "Player 1", "Player 2")
	if err != nil {
		t.Fatalf("Failed to create game: %v", err)
	}
	defer server.RemoveSession(gameID)

	t.Run("WebSocket Handshake and Basic Message", func(t *testing.T) {
		// Connect as Player 1 (Red)
		dialer := websocket.Dialer{}
		conn1, _, err := dialer.Dial(wsURL+"/game/"+gameID+"?player=0", nil)
		if err != nil {
			t.Fatalf("Failed to connect: %v", err)
		}
		defer func() { _ = conn1.Close() }()

		// Should receive initial game state and board state
		for i := range 2 {
			_, message, err := conn1.ReadMessage()
			if err != nil {
				t.Fatalf("Failed to read initial message %d: %v", i, err)
			}
			var msg WSMessage
			if err := json.Unmarshal(message, &msg); err != nil {
				t.Fatalf("Failed to unmarshal initial message %d: %v", i, err)
			}
		}

		// Connect as Player 2 (Blue)
		conn2, _, err := dialer.Dial(wsURL+"/game/"+gameID+"?player=1", nil)
		if err != nil {
			t.Fatalf("Failed to connect P2: %v", err)
		}
		defer func() { _ = conn2.Close() }()

		// Player 1 sends a chat message (if supported) or just wait
		// Let's simulate a move if possible, but move needs pieces.
		// For simplicity, let's just check if they are both registered in hub
		time.Sleep(50 * time.Millisecond)

		handler.Hub.mutex.RLock()
		clientCount := len(handler.Hub.clients)
		handler.Hub.mutex.RUnlock()

		if clientCount != 2 {
			t.Errorf("Expected 2 clients in hub, got %d", clientCount)
		}

		// Test Ping/Pong
		pingMsg := WSMessage{Type: MsgTypePing}
		pingBytes, _ := json.Marshal(pingMsg)
		if err := conn1.WriteMessage(websocket.TextMessage, pingBytes); err != nil {
			t.Fatalf("Failed to send ping: %v", err)
		}

		// Wait for pong
		timeout := time.After(200 * time.Millisecond)
		foundPong := false
		for !foundPong {
			select {
			case <-timeout:
				t.Fatal("Timed out waiting for pong")
			default:
				_, message, err := conn1.ReadMessage()
				if err != nil {
					t.Fatalf("Failed to read message for pong: %v", err)
				}
				var msg WSMessage
				if err := json.Unmarshal(message, &msg); err != nil {
					t.Logf("Failed to unmarshal message: %v", err)
				}
				if msg.Type == MsgTypePong {
					foundPong = true
				}
			}
		}
	})
}
