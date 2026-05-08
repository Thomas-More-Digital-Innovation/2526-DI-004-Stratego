package api

import (
	"digital-innovation/gostrategy/db"
	"digital-innovation/gostrategy/engine"
	"digital-innovation/gostrategy/game"
	"testing"
	"time"
)

func TestWSHub_Cleanup(t *testing.T) {
	// Setup session
	p1 := engine.NewPlayer(0, "P1", "red")
	p2 := engine.NewPlayer(1, "P2", "blue")
	c1 := engine.NewHumanPlayerController(&p1)
	c2 := engine.NewHumanPlayerController(&p2)
	session := game.NewSession("cleanup-test", c1, c2)

	// Create hub
	hub := NewWSHub(session, "human_vs_ai")

	// Set a very short cleanup period for testing
	hub.cleanupPeriod = 100 * time.Millisecond

	cleanupChan := make(chan bool, 1)
	hub.OnCleanup = func() {
		cleanupChan <- true
	}

	// Start hub loop in background
	go hub.Run()
	defer hub.Stop()

	// Initially, Run() starts a cleanup timer (double cleanupPeriod in our case)
	// Wait for the initial cleanup timer to fire
	select {
	case <-cleanupChan:
		// OK
	case <-time.After(500 * time.Millisecond):
		t.Error("Expected cleanup to be called when no clients connect")
	}
}

func TestWSHub_CleanupCancelledOnConnect(t *testing.T) {
	// Setup session
	p1 := engine.NewPlayer(0, "P1", "red")
	p2 := engine.NewPlayer(1, "P2", "blue")
	c1 := engine.NewHumanPlayerController(&p1)
	c2 := engine.NewHumanPlayerController(&p2)
	session := game.NewSession("cleanup-cancel-test", c1, c2)

	// Create hub
	hub := NewWSHub(session, "human_vs_ai")
	hub.cleanupPeriod = 200 * time.Millisecond

	cleanupChan := make(chan bool, 1)
	hub.OnCleanup = func() {
		cleanupChan <- true
	}

	go hub.Run()
	defer hub.Stop()

	// Register a client
	client := &WSClient{
		hub:      hub,
		Username: db.TestUser,
		send:     make(chan []byte, 10),
	}
	// Background consumer for client messages to avoid blocking hub
	go func() {
		for range client.send {
			_ = 0 // Consume and discard messages to avoid blocking the hub
		}
	}()

	hub.register <- client

	// Wait for hub to process registration and cancel initial timer
	time.Sleep(100 * time.Millisecond)

	select {
	case <-cleanupChan:
		t.Error("Cleanup should have been cancelled when client connected")
	default:
		// OK
	}

	// Unregister client
	hub.unregister <- client

	// Now it should trigger cleanup
	select {
	case <-cleanupChan:
		// OK
	case <-time.After(1 * time.Second):
		t.Error("Cleanup should have been triggered after client disconnected")
	}
}
