package core

import (
	"digital-innovation/gostrategy/api/ws"
	"digital-innovation/gostrategy/engine"
	"digital-innovation/gostrategy/game"
	"digital-innovation/gostrategy/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsUserInActiveGame(t *testing.T) {
	s := NewGameServer()
	userID := 123
	gameID := "test-game"

	// Properly initialize session
	p1 := engine.NewPlayer(0, "Red", "red")
	p2 := engine.NewPlayer(1, "Blue", "blue")
	session := game.NewSession(gameID, engine.NewHumanPlayerController(&p1), engine.NewHumanPlayerController(&p2))
	session.Player1UserID = &userID

	handler := &SessionHandler{
		Session: session,
	}
	s.Sessions[gameID] = handler

	// Test user in game
	foundHandler, found := s.IsUserInActiveGame(userID)
	assert.True(t, found)
	assert.Equal(t, handler, foundHandler)

	// Test user not in game
	_, found = s.IsUserInActiveGame(456)
	assert.False(t, found)
}

func TestIsWaitingForCleanup(t *testing.T) {
	userID := 123
	gameID := "test-game"

	p1 := engine.NewPlayer(0, "Red", "red")
	p2 := engine.NewPlayer(1, "Blue", "blue")
	session := game.NewSession(gameID, engine.NewHumanPlayerController(&p1), engine.NewHumanPlayerController(&p2))
	hub := ws.NewHub(session, models.HumanVsAi)

	sh := &SessionHandler{
		Session: session,
		Hub:     hub,
	}

	// 1. Initially, user is not connected, so it should be true (waiting for cleanup/reconnect)
	assert.True(t, sh.IsWaitingForCleanup(userID))

	// 2. User connects
	client := ws.NewTestClient()
	client.UserID = userID

	// Start hub if not already running (Hub.Run() is normally started by GameServer, but here we created it manually)
	go hub.Run()

	hub.Register() <- client
	// Give it a bit of time to process registration
	time.Sleep(10 * time.Millisecond)

	// Now user is connected, so it should NOT be waiting for cleanup
	assert.False(t, sh.IsWaitingForCleanup(userID))

	// 3. Game is over
	sh.Session.GetGame().SetWinner(&p1, game.WinCauseFlagCaptured)

	// Even if user is connected, if game is over, it should be waiting for cleanup
	assert.True(t, sh.IsWaitingForCleanup(userID))
}
