package core

import (
	"digital-innovation/gostrategy/internal/api/ws"
	"digital-innovation/gostrategy/internal/db"
	"digital-innovation/gostrategy/internal/models"
	"digital-innovation/gostrategy/pkg/game"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testGameID = "test-game"
const testUserID = 123

func TestIsUserInActiveGame(t *testing.T) {
	s := NewGameServer()
	defer s.Stop()

	// Properly initialize session
	p1 := game.NewPlayer(0, "Red", "red")
	p2 := game.NewPlayer(1, "Blue", "blue")
	session := game.NewSession(testGameID, game.NewHumanPlayerController(&p1), game.NewHumanPlayerController(&p2))
	session.SetPlayer1Associate(testUserID, "Red")

	handler := &SessionHandler{
		Session: session,
	}
	s.Sessions[testGameID] = handler

	// Test user in game
	foundHandler, found := s.IsUserInActiveGame(testUserID)
	assert.True(t, found)
	assert.Equal(t, handler, foundHandler)

	// Test user not in game
	_, found = s.IsUserInActiveGame(456)
	assert.False(t, found)
}

func TestGetUserActiveGameSeat(t *testing.T) {
	s := NewGameServer()
	defer s.Stop()
	userID1 := testUserID
	userID2 := testUserID + 1

	p1 := game.NewPlayer(0, "Red", "red")
	p2 := game.NewPlayer(1, "Blue", "blue")
	session := game.NewSession(testGameID, game.NewHumanPlayerController(&p1), game.NewHumanPlayerController(&p2))
	session.SetPlayer1Associate(userID1, "Red")
	session.SetPlayer2Associate(userID2, "Blue")

	handler := &SessionHandler{
		Session: session,
	}
	s.Sessions[testGameID] = handler

	foundHandler, seatIndex, found := s.GetUserActiveGameSeat(userID1)
	assert.True(t, found)
	assert.Equal(t, handler, foundHandler)
	assert.Equal(t, 0, seatIndex)

	foundHandler, seatIndex, found = s.GetUserActiveGameSeat(userID2)
	assert.True(t, found)
	assert.Equal(t, handler, foundHandler)
	assert.Equal(t, 1, seatIndex)

	_, seatIndex, found = s.GetUserActiveGameSeat(999)
	assert.False(t, found)
	assert.Equal(t, -1, seatIndex)
}

func TestIsWaitingForCleanup(t *testing.T) {
	p1 := game.NewPlayer(0, "Red", "red")
	p2 := game.NewPlayer(1, "Blue", "blue")
	session := game.NewSession(testGameID, game.NewHumanPlayerController(&p1), game.NewHumanPlayerController(&p2))
	hub := ws.NewHub(session, models.HumanVsAi)
	defer hub.Stop()

	sh := &SessionHandler{
		Session: session,
		Hub:     hub,
	}

	// 1. Initially, user is not connected, so it should be true (waiting for cleanup/reconnect)
	assert.True(t, sh.IsWaitingForCleanup(testUserID))

	// 2. User connects
	client := ws.NewTestClient()
	client.UserID = testUserID

	// Start hub if not already running (Hub.Run() is normally started by GameServer, but here we created it manually)
	go hub.Run()

	hub.Register() <- client
	// Give it a bit of time to process registration
	time.Sleep(10 * time.Millisecond)

	// Now user is connected, so it should NOT be waiting for cleanup
	assert.False(t, sh.IsWaitingForCleanup(testUserID))

	// 3. Game is over
	sh.Session.SetWinner(&p1, game.WinCauseFlagCaptured)

	// Even if user is connected, if game is over, it should be waiting for cleanup
	assert.True(t, sh.IsWaitingForCleanup(testUserID))
}
func TestCreateGame(t *testing.T) {
	s := NewGameServer()
	defer s.Stop()

	// Test HumanVsAi
	handler, err := s.CreateGame("game-hva", models.HumanVsAi, "Alice", models.Fafo)
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	assert.Equal(t, models.HumanVsAi, handler.GameType)

	// Test AiVsAi
	handler, err = s.CreateGame("game-ava", models.AiVsAi, models.Fafo, models.Fafo)
	assert.NoError(t, err)
	assert.NotNil(t, handler)

	// Test HumanVsHuman
	handler, err = s.CreateGame("game-hvh", models.HumanVsHuman, "Alice", "Bob")
	assert.NoError(t, err)
	assert.NotNil(t, handler)

	// Test Unknown Type
	_, err = s.CreateGame("game-err", "unknown", "A", "B")
	assert.Error(t, err)

	// Test Duplicate ID
	_, err = s.CreateGame("game-hva", models.HumanVsAi, "Alice", models.Fafo)
	assert.Error(t, err)
}

func TestRemoveSession(t *testing.T) {
	s := NewGameServer()
	defer s.Stop()
	_, err := s.CreateGame("rem-test", models.HumanVsAi, "A", models.Fafo)
	assert.NoError(t, err)

	s.RemoveSession("rem-test")
	_, exists := s.GetSession("rem-test")
	assert.False(t, exists)

	// Test non-existent
	s.RemoveSession("non-existent")
}

func TestStop(_ *testing.T) {
	s := NewGameServer()
	s.Stop()
	// Should not panic or error
}

func TestPrintRoutes(_ *testing.T) {
	s := NewGameServer()
	defer s.Stop()
	s.PrintRoutes()
	// Should not panic
}

func TestGameServer_Completion(t *testing.T) {
	db.SetupTestDB(t)
	s := NewGameServer()
	defer s.Stop()

	handler, _ := s.CreateGame("test-comp", models.HumanVsHuman, "P1", "P2")

	// Simulate game over directly calling handleGameOver for coverage
	s.handleGameOver(handler.Session, handler.Hub)
}

func TestGameServer_Limits(t *testing.T) {
	s := NewGameServer()
	defer s.Stop()
	// Simulate many sessions
	for i := range 500 {
		s.Sessions[fmt.Sprintf("g%d", i)] = &SessionHandler{}
	}

	_, err := s.CreateGame("too-many", models.HumanVsAi, "A", models.Fafo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "server busy")
}

func TestGameServer_Abort(t *testing.T) {
	s := NewGameServer()
	defer s.Stop()

	handler, _ := s.CreateGame("test-abort", models.HumanVsHuman, "P1", "P2")

	// Give monitorGame time to start and sleep
	time.Sleep(150 * time.Millisecond)

	// Abort the session
	handler.Session.Stop()

	// Give monitorGame time to process the abort
	time.Sleep(300 * time.Millisecond)

	_, exists := s.GetSession("test-abort")
	assert.False(t, exists)
}
