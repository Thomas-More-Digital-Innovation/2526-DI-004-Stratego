package ws_test

import (
	"digital-innovation/gostrategy/internal/api/ws"
	"digital-innovation/gostrategy/internal/game"
	"digital-innovation/gostrategy/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHub_Broadcasters(t *testing.T) {
	player1 := game.NewPlayer(0, "P1", "red")
	player2 := game.NewPlayer(1, "P2", "blue")
	c1 := game.NewHumanPlayerController(&player1)
	c2 := game.NewHumanPlayerController(&player2)
	session := game.NewSession("test-broadcasters", c1, c2)
	hub := ws.NewHub(session, models.HumanVsHuman)

	go hub.Run()
	defer hub.Stop()

	client := ws.NewTestClient()
	client.Hub = hub
	client.Session = session
	client.SeatIndex = 0

	hub.Register() <- client
	time.Sleep(50 * time.Millisecond)

	// Drain initial messages
	drainMessages(client, 2)

	// 1. BroadcastMessage
	hub.BroadcastMessage("test-type", "test-data")
	msg := <-client.SendChan()
	assert.Contains(t, string(msg), "test-type")
	assert.Contains(t, string(msg), "test-data")

	// 2. BroadcastGameState
	hub.BroadcastGameState()
	msg = <-client.SendChan()
	assert.Contains(t, string(msg), "gameState")

	// 3. BroadcastSetupBoard
	hub.BroadcastSetupBoard()
	msg = <-client.SendChan()
	assert.Contains(t, string(msg), "boardState")

	// 4. BroadcastFullState
	hub.BroadcastFullState()
	drainMessages(client, 2) // gameState and boardState

	// 5. BroadcastCombat
	combat := &game.CombatResult{
		Occurred:         true,
		AttackerPiece:    game.NewPiece(models.Marshal, &player1),
		DefenderPiece:    game.NewPiece(models.General, &player2),
		AttackerPosition: game.NewPosition(0, 0),
		DefenderPosition: game.NewPosition(0, 1),
	}
	hub.BroadcastCombat(combat)
	msg = <-client.SendChan()
	assert.Contains(t, string(msg), "combat")

	// 6. BroadcastMoveHistory
	hub.BroadcastMoveHistory()
	msg = <-client.SendChan()
	assert.Contains(t, string(msg), "moveHistory")

	// 7. BroadcastGameTransition (starts game)
	hub.BroadcastGameTransition()
	drainMessages(client, 2)
}

func TestHub_RevealedBroadcasters(_ *testing.T) {
	player1 := game.NewPlayer(0, "AI1", "red")
	player2 := game.NewPlayer(1, "AI2", "blue")
	c1 := game.NewHumanPlayerController(&player1)
	c2 := game.NewHumanPlayerController(&player2)
	session := game.NewSession("test-revealed", c1, c2)
	hub := ws.NewHub(session, models.AiVsAi)

	go hub.Run()
	defer hub.Stop()

	client := ws.NewTestClient()
	client.Hub = hub
	client.Session = session
	client.SeatIndex = -1 // Spectator

	hub.Register() <- client
	time.Sleep(50 * time.Millisecond)
	drainMessages(client, 2)

	hub.BroadcastGameTransition()
	drainMessages(client, 2)
}

func TestHub_Perspectives(_ *testing.T) {
	player1 := game.NewPlayer(0, "P1", "red")
	player2 := game.NewPlayer(1, "P2", "blue")
	c1 := game.NewHumanPlayerController(&player1)
	c2 := game.NewHumanPlayerController(&player2)
	session := game.NewSession("test-perspectives", c1, c2)
	hub := ws.NewHub(session, models.HumanVsHuman)

	go hub.Run()
	defer hub.Stop()

	// Player 1
	client1 := ws.NewTestClient()
	client1.Hub = hub
	client1.Session = session
	client1.SeatIndex = 0
	hub.Register() <- client1
	time.Sleep(50 * time.Millisecond)
	drainMessages(client1, 2)

	// Player 2
	client2 := ws.NewTestClient()
	client2.Hub = hub
	client2.Session = session
	client2.SeatIndex = 1
	hub.Register() <- client2
	time.Sleep(50 * time.Millisecond)
	drainMessages(client2, 2)

	// Spectator
	client3 := ws.NewTestClient()
	client3.Hub = hub
	client3.Session = session
	client3.SeatIndex = -1
	hub.Register() <- client3
	time.Sleep(50 * time.Millisecond)
	drainMessages(client3, 2)

	// Start game and move to test history filtering
	_ = session.RandomizeSetup(0)
	_ = session.RandomizeSetup(1)
	_ = session.StartGameFromSetup(false)

	// Manually place a piece to ensure board state has content
	board := session.GetBoard()
	board.SetPieceAt(game.NewPosition(5, 5), game.NewPiece(models.Marshal, &player1))
	board.SetPieceAt(game.NewPosition(5, 6), game.NewPiece(models.General, &player2))

	hub.BroadcastGameState()
	drainMessages(client1, 1)
	drainMessages(client2, 1)
	drainMessages(client3, 1)
}

func TestHub_StatusAndCleanup(t *testing.T) {
	player1 := game.NewPlayer(0, "P1", "red")
	player2 := game.NewPlayer(1, "P2", "blue")
	c1 := game.NewHumanPlayerController(&player1)
	c2 := game.NewHumanPlayerController(&player2)
	session := game.NewSession("test-status", c1, c2)
	hub := ws.NewHub(session, models.HumanVsHuman)

	go hub.Run()
	defer hub.Stop()

	client := ws.NewTestClient()
	client.UserID = 123
	hub.Register() <- client
	time.Sleep(50 * time.Millisecond)

	assert.True(t, hub.IsUserConnected(123))
	assert.False(t, hub.IsUserConnected(456))

	// Test GameOver cleanup
	hub.StartGameOverCleanup(10 * time.Millisecond)
	time.Sleep(50 * time.Millisecond)
	// Hub should be stopped
	assert.True(t, hub.IsStopped())
}

func TestClient_Methods(t *testing.T) {
	player1 := game.NewPlayer(0, "P1", "red")
	player2 := game.NewPlayer(1, "P2", "blue")
	c1 := game.NewHumanPlayerController(&player1)
	c2 := game.NewHumanPlayerController(&player2)
	session := game.NewSession("test-client", c1, c2)
	hub := ws.NewHub(session, models.HumanVsHuman)

	client := &ws.Client{
		Hub:       hub,
		Session:   session,
		SeatIndex: 0,
		Username:  "alice",
		UserID:    123,
	}

	assert.Equal(t, 0, client.GetSeatIndex())
	assert.Equal(t, "alice", client.GetUsername())
	assert.Equal(t, 123, client.GetUserID())
	assert.True(t, client.IsAuthorized())
}

func drainMessages(client *ws.Client, count int) {
	for range count {
		select {
		case <-client.SendChan():
		case <-time.After(500 * time.Millisecond):
			return
		}
	}
}
