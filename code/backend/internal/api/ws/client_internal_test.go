package ws

import (
	"digital-innovation/gostrategy/internal/models"
	"digital-innovation/gostrategy/pkg/game"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClient_IsAuthorized(t *testing.T) {
	player1 := game.NewPlayer(0, "P1", "red")
	player2 := game.NewPlayer(1, "P2", "blue")
	c1 := game.NewHumanPlayerController(&player1)
	c2 := game.NewHumanPlayerController(&player2)
	session := game.NewSession("auth-test", c1, c2)

	// Set player IDs
	p1ID := 101
	session.SetPlayer1Associate(p1ID, "P1")

	hub := NewHub(session, models.HumanVsAi)

	t.Run("player is authorized", func(t *testing.T) {
		client := NewTestClient()
		client.Session = session
		client.Hub = hub
		client.SeatIndex = 0 // Player 1
		client.UserID = p1ID

		assert.True(t, client.IsAuthorized())
	})

	t.Run("spectator creator is authorized in AI game", func(t *testing.T) {
		client := NewTestClient()
		client.Session = session
		client.Hub = hub
		client.SeatIndex = -1 // Spectator
		client.UserID = p1ID  // Creator

		assert.True(t, client.IsAuthorized())
	})

	t.Run("spectator non-creator is NOT authorized", func(t *testing.T) {
		client := NewTestClient()
		client.Session = session
		client.Hub = hub
		client.SeatIndex = -1
		client.UserID = 999

		assert.False(t, client.IsAuthorized())
	})
}

func TestClient_Send(t *testing.T) {
	client := NewTestClient()

	t.Run("send success", func(t *testing.T) {
		data := []byte("test")
		success := client.Send(data, 10*time.Millisecond)
		assert.True(t, success)

		received := <-client.send
		assert.Equal(t, data, received)
	})

	t.Run("send timeout", func(t *testing.T) {
		// Fill buffer completely (256 slots)
		for i := 0; i < 256; i++ {
			success := client.Send([]byte("fill"), 10*time.Millisecond)
			assert.True(t, success, "Failed to fill buffer at index %d", i)
		}

		// This one should definitely fail if buffer is 256
		success := client.Send([]byte("timeout"), 1*time.Millisecond)
		assert.False(t, success)
	})

	t.Run("send closed", func(t *testing.T) {
		client.Close()
		success := client.Send([]byte("closed"), 0)
		assert.False(t, success)
	})
}

func TestClient_SendJSON(t *testing.T) {
	client := NewTestClient()

	client.SendJSON("test-type", "test-data")

	msgBytes := <-client.send
	var msg struct {
		Type string `json:"type"`
		Data string `json:"data"`
	}
	_ = json.Unmarshal(msgBytes, &msg)

	assert.Equal(t, "test-type", msg.Type)
	assert.Equal(t, "test-data", msg.Data)
}
