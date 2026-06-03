package api_test

import (
	"digital-innovation/gostrategy/internal/api"
	"digital-innovation/gostrategy/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGameServer(t *testing.T) {
	server := api.NewGameServer()
	assert.NotNil(t, server)
}

func TestCreateGameHumanVsAI(t *testing.T) {
	server := api.NewGameServer()
	handler, err := server.CreateGame("test-game-1", models.HumanVsAi, "Human Player", models.Fafo)
	require.NoError(t, err)
	require.NotNil(t, handler)
	require.NotNil(t, handler.Session)
	require.NotNil(t, handler.Hub)

	assert.Equal(t, models.HumanVsAi, handler.GameType)
	assert.True(t, handler.Session.IsSetupPhase())
}

func TestCreateGameAIVsAI(t *testing.T) {
	server := api.NewGameServer()
	handler, err := server.CreateGame("test-game-2", models.AiVsAi, models.Fafo, models.Fafo)
	require.NoError(t, err)
	require.NotNil(t, handler)

	assert.True(t, handler.Session.IsSetupPhase())
	assert.False(t, handler.Session.IsRunning())

	// Clean up
	handler.Session.Stop()
}

func TestCreateGameHumanVsHuman(t *testing.T) {
	server := api.NewGameServer()
	handler, err := server.CreateGame("test-game-3", models.HumanVsHuman, "Player 1", "Player 2")
	require.NoError(t, err)
	require.NotNil(t, handler)

	assert.Equal(t, models.HumanVsHuman, handler.GameType)
	assert.True(t, handler.Session.IsSetupPhase())
}

func TestCreateGameDuplicateID(t *testing.T) {
	server := api.NewGameServer()
	gameID := "duplicate-test"

	_, err := server.CreateGame(gameID, models.HumanVsAi, "Human", models.Fafo)
	require.NoError(t, err)

	_, err = server.CreateGame(gameID, models.HumanVsAi, "Human", models.Fafo)
	assert.Error(t, err)
}

func TestCreateGameInvalidType(t *testing.T) {
	server := api.NewGameServer()
	_, err := server.CreateGame("invalid-type-game", "InvalidGameType", "P1", "P2")
	assert.Error(t, err)
}

func TestGetSession(t *testing.T) {
	server := api.NewGameServer()
	gameID := "get-session-test"

	handler, err := server.CreateGame(gameID, models.HumanVsAi, "Human", models.Fafo)
	require.NoError(t, err)

	retrieved, exists := server.GetSession(gameID)
	assert.True(t, exists)
	assert.Equal(t, handler, retrieved)
}

func TestGetSessionNonExistent(t *testing.T) {
	server := api.NewGameServer()
	_, exists := server.GetSession("non-existent-game")
	assert.False(t, exists)
}
