package api

import (
	"digital-innovation/stratego/engine"
	"digital-innovation/stratego/game"
	"digital-innovation/stratego/models"
	"testing"
	"time"
)

// MockConn is a simple mock for websocket.Conn enough for Hub testing
type MockConn struct{}

func TestWSHub_RegisterUnregister(t *testing.T) {
	player1 := engine.NewPlayer(0, "P1", "red")
	player2 := engine.NewPlayer(1, "P2", "blue")
	c1 := engine.NewHumanPlayerController(&player1)
	c2 := engine.NewHumanPlayerController(&player2)
	session := game.NewGameSession("test-reg", c1, c2)
	hub := NewWSHub(session, models.HumanVsHuman)

	go hub.Run()
	defer hub.Stop()

	client := &WSClient{
		hub:  hub,
		send: make(chan []byte, 10),
	}

	// Register
	hub.register <- client

	// Wait a bit for processing
	time.Sleep(10 * time.Millisecond)

	hub.mutex.RLock()
	count := len(hub.clients)
	hub.mutex.RUnlock()

	if count != 1 {
		t.Errorf("Expected 1 client, got %d", count)
	}

	// Unregister
	hub.unregister <- client
	time.Sleep(10 * time.Millisecond)

	hub.mutex.RLock()
	count = len(hub.clients)
	hub.mutex.RUnlock()

	if count != 0 {
		t.Errorf("Expected 0 clients after unregister, got %d", count)
	}
}

func TestWSHub_AICleanup(t *testing.T) {
	player1 := engine.NewPlayer(0, "AI1", "red")
	player2 := engine.NewPlayer(1, "AI2", "blue")
	c1 := engine.NewHumanPlayerController(&player1)
	c2 := engine.NewHumanPlayerController(&player2)
	session := game.NewGameSession("test-ai-cleanup", c1, c2)

	cleanupSignal := make(chan bool, 1)
	hub := NewWSHub(session, models.AiVsAi)
	hub.OnCleanup = func() {
		cleanupSignal <- true
	}

	go hub.Run()

	client := &WSClient{
		hub:  hub,
		send: make(chan []byte, 10),
	}

	// Register then unregister
	hub.register <- client
	time.Sleep(10 * time.Millisecond)
	hub.unregister <- client

	// Wait for cleanup with timeout
	select {
	case <-cleanupSignal:
		// Success
	case <-time.After(1 * time.Second):
		t.Error("Timed out waiting for OnCleanup to be called")
	}

	if !hub.IsStopped() {
		t.Error("Expected hub to be stopped")
	}
}

func TestWSHub_Broadcast(t *testing.T) {
	player1 := engine.NewPlayer(0, "P1", "red")
	player2 := engine.NewPlayer(1, "P2", "blue")
	c1 := engine.NewHumanPlayerController(&player1)
	c2 := engine.NewHumanPlayerController(&player2)
	session := game.NewGameSession("test-broadcast", c1, c2)
	hub := NewWSHub(session, models.HumanVsHuman)

	go hub.Run()
	defer hub.Stop()

	client := &WSClient{
		hub:  hub,
		send: make(chan []byte, 10),
	}
	hub.register <- client

	// Drain the initial messages (gameState and boardState) sent upon registration
	for range 2 {
		select {
		case <-client.send:
		case <-time.After(200 * time.Millisecond):
			t.Fatal("Timed out waiting for initial messages")
		}
	}

	message := []byte("test message")
	hub.broadcast <- message

	select {
	case received := <-client.send:
		if string(received) != string(message) {
			t.Errorf("Expected %s, got %s", message, received)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Timed out waiting for broadcast")
	}
}
