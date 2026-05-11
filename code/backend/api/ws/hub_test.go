package ws_test

import (
	"digital-innovation/gostrategy/api/ws"
	"digital-innovation/gostrategy/engine"
	"digital-innovation/gostrategy/game"
	"digital-innovation/gostrategy/models"
	"testing"
	"time"
)

func TestWSHub_RegisterUnregister(t *testing.T) {
	player1 := engine.NewPlayer(0, "P1", "red")
	player2 := engine.NewPlayer(1, "P2", "blue")
	c1 := engine.NewHumanPlayerController(&player1)
	c2 := engine.NewHumanPlayerController(&player2)
	session := game.NewSession("test-reg", c1, c2)
	hub := ws.NewHub(session, models.HumanVsHuman)

	go hub.Run()
	defer hub.Stop()

	client := &ws.Client{}

	// Register
	hub.Register() <- client

	// Wait a bit for processing
	time.Sleep(10 * time.Millisecond)

	count := len(hub.Clients())

	if count != 1 {
		t.Errorf("Expected 1 client, got %d", count)
	}

	// Unregister
	hub.Unregister() <- client
	time.Sleep(10 * time.Millisecond)

	count = len(hub.Clients())

	if count != 0 {
		t.Errorf("Expected 0 clients after unregister, got %d", count)
	}
}

func TestWSHub_AICleanup(t *testing.T) {
	player1 := engine.NewPlayer(0, "AI1", "red")
	player2 := engine.NewPlayer(1, "AI2", "blue")
	c1 := engine.NewHumanPlayerController(&player1)
	c2 := engine.NewHumanPlayerController(&player2)
	session := game.NewSession("test-ai-cleanup", c1, c2)

	cleanupSignal := make(chan bool, 1)
	hub := ws.NewHub(session, models.AiVsAi)
	hub.OnCleanup = func() {
		cleanupSignal <- true
	}

	go hub.Run()

	client := &ws.Client{}

	// Register then unregister
	hub.Register() <- client
	time.Sleep(10 * time.Millisecond)
	hub.Unregister() <- client

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
	session := game.NewSession("test-broadcast", c1, c2)
	hub := ws.NewHub(session, models.HumanVsHuman)

	go hub.Run()
	defer hub.Stop()

	client := &ws.Client{}
	hub.Register() <- client

	// Drain the initial messages (gameState and boardState) sent upon registration
	for range 2 {
		select {
		case <-client.SendChan():
		case <-time.After(200 * time.Millisecond):
			t.Fatal("Timed out waiting for initial messages")
		}
	}

	message := []byte("test message")
	hub.Broadcast() <- message

	select {
	case received := <-client.SendChan():
		if string(received) != string(message) {
			t.Errorf("Expected %s, got %s", message, received)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Timed out waiting for broadcast")
	}
}

func TestNewHub(t *testing.T) {
	player1 := engine.NewPlayer(0, "Player1", "red")
	player2 := engine.NewPlayer(1, "Player2", "blue")
	controller1 := engine.NewHumanPlayerController(&player1)
	controller2 := engine.NewHumanPlayerController(&player2)
	session := game.NewSession("test-hub", controller1, controller2)

	hub := ws.NewHub(session, models.HumanVsAi)

	if hub == nil {
		t.Fatal("Expected NewHub to return a hub, but got nil")
	}
}

func TestWSHubWithDifferentGameTypes(t *testing.T) {
	player1 := engine.NewPlayer(0, "AI1", "red")
	player2 := engine.NewPlayer(1, "AI2", "blue")
	controller1 := engine.NewHumanPlayerController(&player1)
	controller2 := engine.NewHumanPlayerController(&player2)

	testCases := []struct {
		name     string
		gameType string
	}{
		{"HumanVsAI", models.HumanVsAi},
		{"AIVsAI", models.AiVsAi},
		{"HumanVsHuman", models.HumanVsHuman},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			session := game.NewSession("test-"+tc.name, controller1, controller2)
			hub := ws.NewHub(session, tc.gameType)

			if hub == nil {
				t.Errorf("Expected hub to be created for game type %s", tc.gameType)
			}
		})
	}
}
