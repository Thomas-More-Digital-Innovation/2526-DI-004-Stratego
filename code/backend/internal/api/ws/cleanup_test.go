package ws_test

import (
	"digital-innovation/gostrategy/internal/api/ws"
	"digital-innovation/gostrategy/internal/game"
	"testing"
	"time"
)

func TestWSHub_Cleanup(t *testing.T) {
	// Setup session
	p1 := game.NewPlayer(0, "P1", "red")
	p2 := game.NewPlayer(1, "P2", "blue")
	c1 := game.NewHumanPlayerController(&p1)
	c2 := game.NewHumanPlayerController(&p2)
	session := game.NewSession("cleanup-test", c1, c2)

	// Create hub
	hub := ws.NewHub(session, "human_vs_ai")

	// Set a very short cleanup period for testing
	hub.SetCleanupPeriod(100 * time.Millisecond)

	cleanupChan := make(chan bool, 1)
	hub.OnCleanup = func() {
		cleanupChan <- true
	}

	// Start hub loop in background
	go hub.Run()
	defer hub.Stop()

	// Initially, Run() starts a cleanup timer
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
	p1 := game.NewPlayer(0, "P1", "red")
	p2 := game.NewPlayer(1, "P2", "blue")
	c1 := game.NewHumanPlayerController(&p1)
	c2 := game.NewHumanPlayerController(&p2)
	session := game.NewSession("cleanup-cancel-test", c1, c2)

	// Create hub
	hub := ws.NewHub(session, "human_vs_ai")
	hub.SetCleanupPeriod(200 * time.Millisecond)

	cleanupChan := make(chan bool, 1)
	hub.OnCleanup = func() {
		cleanupChan <- true
	}

	go hub.Run()
	defer hub.Stop()

	// Register a client
	client := ws.NewTestClient()
	// Background consumer for client messages
	go func() {
		for range client.SendChan() {
			_ = 0
		}
	}()

	hub.Register() <- client

	// Wait for hub to process registration and cancel initial timer
	time.Sleep(100 * time.Millisecond)

	select {
	case <-cleanupChan:
		t.Error("Cleanup should have been cancelled when client connected")
	default:
		// OK
	}

	// Unregister client
	hub.Unregister() <- client

	// Now it should trigger cleanup
	select {
	case <-cleanupChan:
		// OK
	case <-time.After(1 * time.Second):
		t.Error("Cleanup should have been triggered after client disconnected")
	}
}
