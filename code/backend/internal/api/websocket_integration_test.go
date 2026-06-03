package api

import (
	"digital-innovation/gostrategy/internal/api/dto"
	"digital-innovation/gostrategy/internal/db"
	"digital-innovation/gostrategy/internal/models"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_WebSocket(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db.SetupTestDB(t)

	server := NewGameServer()
	server.SetupRoutes()

	ts := httptest.NewServer(server.Router)
	defer ts.Close()

	// Convert http URL to ws URL
	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http")

	// 1. Create a game session manually in server
	gameID := "ws-test-game"
	handler, err := server.CreateGame(gameID, models.HumanVsHuman, "Player 1", "Player 2")
	require.NoError(t, err)
	defer server.RemoveSession(gameID)

	t.Run("WebSocket Handshake and Basic Message", func(t *testing.T) {
		// Connect as Player 1 (Red)
		dialer := websocket.Dialer{}
		conn1, _, err := dialer.Dial(wsURL+"/game/"+gameID+"?player=0", nil)
		require.NoError(t, err)
		defer func() { _ = conn1.Close() }()

		// Should receive initial game state and board state
		for i := range 2 {
			_, message, err := conn1.ReadMessage()
			require.NoError(t, err, "Failed to read initial message %d", i)

			var msg dto.WSMessage
			err = json.Unmarshal(message, &msg)
			require.NoError(t, err, "Failed to unmarshal initial message %d", i)
		}

		// Connect as Player 2 (Blue)
		conn2, _, err := dialer.Dial(wsURL+"/game/"+gameID+"?player=1", nil)
		require.NoError(t, err, "Failed to connect P2")
		defer func() { _ = conn2.Close() }()

		// Player 1 sends a chat message (if supported) or just wait
		time.Sleep(50 * time.Millisecond)

		clientCount := handler.Hub.ClientCount()
		assert.Equal(t, 2, clientCount)

		// Test Ping/Pong
		pingMsg := dto.WSMessage{Type: dto.MsgTypePing}
		pingBytes, err := json.Marshal(pingMsg)
		require.NoError(t, err)
		err = conn1.WriteMessage(websocket.TextMessage, pingBytes)
		require.NoError(t, err, "Failed to send ping")

		// Wait for pong
		err = conn1.SetReadDeadline(time.Now().Add(2 * time.Second))
		require.NoError(t, err)

		foundPong := false
		for !foundPong {
			_, message, err := conn1.ReadMessage()
			require.NoError(t, err, "Timed out or failed waiting for pong")

			var msg dto.WSMessage
			if err := json.Unmarshal(message, &msg); err != nil {
				continue
			}
			if msg.Type == dto.MsgTypePong {
				foundPong = true
			}
		}
		assert.True(t, foundPong)
	})
}
