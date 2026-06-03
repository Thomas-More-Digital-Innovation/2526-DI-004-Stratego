package core

import (
	"context"
	"digital-innovation/gostrategy/internal/db"
	"digital-innovation/gostrategy/internal/game"
	"digital-innovation/gostrategy/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResignGame(t *testing.T) {
	db.SetupTestDB(t)

	t.Run("non-existent game ID", func(t *testing.T) {
		s := NewGameServer()
		defer s.Stop()

		err := s.ResignGame("non-existent", 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("non-participant user", func(t *testing.T) {
		s := NewGameServer()
		defer s.Stop()

		handler, err := s.CreateGame("test-resign-non-part", models.HumanVsHuman, "P1", "P2")
		assert.NoError(t, err)

		handler.Session.SetPlayer1Associate(100, "P1")
		handler.Session.SetPlayer2Associate(200, "P2")

		err = s.ResignGame("test-resign-non-part", 999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not a participant")
	})

	t.Run("resignation during setup phase", func(t *testing.T) {
		s := NewGameServer()
		defer s.Stop()

		u1, err := db.CreateUser(context.Background(), "P1_setup", "pass", "")
		assert.NoError(t, err)
		u2, err := db.CreateUser(context.Background(), "P2_setup", "pass", "")
		assert.NoError(t, err)

		handler, err := s.CreateGame("test-resign-setup", models.HumanVsHuman, "P1_setup", "P2_setup")
		assert.NoError(t, err)

		handler.Session.SetPlayer1Associate(u1.ID, "P1_setup")
		handler.Session.SetPlayer2Associate(u2.ID, "P2_setup")

		// Resign during setup
		err = s.ResignGame("test-resign-setup", u1.ID)
		assert.NoError(t, err)

		// Session should be removed
		_, exists := s.GetSession("test-resign-setup")
		assert.False(t, exists)

		// Stats should NOT be updated (wins/losses should both be 0)
		stats1, err := db.GetUserStats(context.Background(), u1.ID)
		assert.NoError(t, err)
		assert.Equal(t, 0, stats1.TotalGames)

		stats2, err := db.GetUserStats(context.Background(), u2.ID)
		assert.NoError(t, err)
		assert.Equal(t, 0, stats2.TotalGames)
	})

	t.Run("resignation during active play", func(t *testing.T) {
		s := NewGameServer()
		defer s.Stop()

		u1, err := db.CreateUser(context.Background(), "P1_active", "pass", "")
		assert.NoError(t, err)
		u2, err := db.CreateUser(context.Background(), "P2_active", "pass", "")
		assert.NoError(t, err)

		handler, err := s.CreateGame("test-resign-active", models.HumanVsHuman, "P1_active", "P2_active")
		assert.NoError(t, err)

		handler.Session.SetPlayer1Associate(u1.ID, "P1_active")
		handler.Session.SetPlayer2Associate(u2.ID, "P2_active")

		// Move out of setup phase
		handler.Session.SetSetupPhaseComplete()

		// Resign
		err = s.ResignGame("test-resign-active", u1.ID)
		assert.NoError(t, err)

		// Stats should be updated: resigning player (u1) loses, opponent (u2) wins
		stats1, err := db.GetUserStats(context.Background(), u1.ID)
		assert.NoError(t, err)
		assert.Equal(t, 1, stats1.TotalGames)
		assert.Equal(t, 0, stats1.Wins)
		assert.Equal(t, 1, stats1.Losses)

		stats2, err := db.GetUserStats(context.Background(), u2.ID)
		assert.NoError(t, err)
		assert.Equal(t, 1, stats2.TotalGames)
		assert.Equal(t, 1, stats2.Wins)
		assert.Equal(t, 0, stats2.Losses)
	})
}

func TestWinGame(t *testing.T) {
	db.SetupTestDB(t)

	t.Run("non-existent game ID", func(t *testing.T) {
		s := NewGameServer()
		defer s.Stop()

		err := s.WinGame("non-existent", 1, "flag_captured")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("non-participant user", func(t *testing.T) {
		s := NewGameServer()
		defer s.Stop()

		handler, err := s.CreateGame("test-wingame-non-part", models.HumanVsHuman, "P1", "P2")
		assert.NoError(t, err)

		handler.Session.SetPlayer1Associate(100, "P1")
		handler.Session.SetPlayer2Associate(200, "P2")

		err = s.WinGame("test-wingame-non-part", 999, "flag_captured")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not a participant")
	})

	t.Run("valid win declaration", func(t *testing.T) {
		s := NewGameServer()
		defer s.Stop()

		u1, err := db.CreateUser(context.Background(), "P1_win", "pass", "")
		assert.NoError(t, err)
		u2, err := db.CreateUser(context.Background(), "P2_win", "pass", "")
		assert.NoError(t, err)

		handler, err := s.CreateGame("test-wingame-valid", models.HumanVsHuman, "P1_win", "P2_win")
		assert.NoError(t, err)

		handler.Session.SetPlayer1Associate(u1.ID, "P1_win")
		handler.Session.SetPlayer2Associate(u2.ID, "P2_win")

		// Move out of setup phase
		handler.Session.SetSetupPhaseComplete()

		// Declare u1 the winner
		err = s.WinGame("test-wingame-valid", u1.ID, "flag_captured")
		assert.NoError(t, err)

		// Verify game state winner set
		state := handler.Session.GetGameState()
		assert.True(t, state.IsGameOver)
		assert.NotNil(t, state.WinnerID)
		assert.Equal(t, 0, *state.WinnerID) // seat 0 represents Player 1 (u1)

		// Stats should be updated
		stats1, err := db.GetUserStats(context.Background(), u1.ID)
		assert.NoError(t, err)
		assert.Equal(t, 1, stats1.Wins)

		stats2, err := db.GetUserStats(context.Background(), u2.ID)
		assert.NoError(t, err)
		assert.Equal(t, 1, stats2.Losses)
	})

	t.Run("game already over", func(t *testing.T) {
		s := NewGameServer()
		defer s.Stop()

		u1, err := db.CreateUser(context.Background(), "P1_already", "pass", "")
		assert.NoError(t, err)
		u2, err := db.CreateUser(context.Background(), "P2_already", "pass", "")
		assert.NoError(t, err)

		handler, err := s.CreateGame("test-wingame-over", models.HumanVsHuman, "P1_already", "P2_already")
		assert.NoError(t, err)

		handler.Session.SetPlayer1Associate(u1.ID, "P1_already")
		handler.Session.SetPlayer2Associate(u2.ID, "P2_already")

		handler.Session.SetWinner(handler.Session.GetGame().Players[0], game.WinCauseFlagCaptured)

		err = s.WinGame("test-wingame-over", u1.ID, "flag_captured")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already over")
	})
}
