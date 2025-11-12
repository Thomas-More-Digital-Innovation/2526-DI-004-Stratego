package api_test

import (
	"digital-innovation/stratego/api"
	"digital-innovation/stratego/models"
	"testing"
)

func TestNewGameServer(t *testing.T) {
	server := api.NewGameServer()

	if server == nil {
		t.Fatal("Expected NewGameServer to return a server, but got nil")
	}
}

func TestCreateGameHumanVsAI(t *testing.T) {
	server := api.NewGameServer()
	handler, err := server.CreateGame("test-game-1", models.HumanVsAi, models.Fafo, models.Fafo)

	if err != nil {
		t.Fatalf("Expected no error creating HumanVsAI game, got: %v", err)
	}

	if handler == nil {
		t.Fatal("Expected handler to be created, but got nil")
	}

	if handler.Session == nil {
		t.Fatal("Expected session to be created, but got nil")
	}

	if handler.Hub == nil {
		t.Fatal("Expected hub to be created, but got nil")
	}

	if handler.GameType != models.HumanVsAi {
		t.Errorf("Expected game type to be HumanVsAi, got: %s", handler.GameType)
	}

	if !handler.Session.IsSetupPhase() {
		t.Error("Expected HumanVsAI game to start in setup phase")
	}
}

func TestCreateGameAIVsAI(t *testing.T) {
	server := api.NewGameServer()
	handler, err := server.CreateGame("test-game-2", models.AiVsAi, models.Fafo, models.Fafo)

	if err != nil {
		t.Fatalf("Expected no error creating AiVsAi game, got: %v", err)
	}

	if handler == nil {
		t.Fatal("Expected handler to be created, but got nil")
	}

	if handler.Session.IsSetupPhase() {
		t.Error("Expected AiVsAi game to skip setup phase")
	}

	if !handler.Session.IsRunning() {
		t.Error("Expected AiVsAi game to start running immediately")
	}

	// Clean up
	handler.Session.Stop()
}

func TestCreateGameHumanVsHuman(t *testing.T) {
	server := api.NewGameServer()
	handler, err := server.CreateGame("test-game-3", models.HumanVsHuman, models.Fafo, models.Fafo)

	if err != nil {
		t.Fatalf("Expected no error creating HumanVsHuman game, got: %v", err)
	}

	if handler.GameType != models.HumanVsHuman {
		t.Errorf("Expected game type to be HumanVsHuman, got: %s", handler.GameType)
	}

	if !handler.Session.IsSetupPhase() {
		t.Error("Expected HumanVsHuman game to start in setup phase")
	}
}

func TestCreateGameDuplicateID(t *testing.T) {
	server := api.NewGameServer()
	gameID := "duplicate-test"

	_, err := server.CreateGame(gameID, models.HumanVsAi, models.Fafo, models.Fafo)
	if err != nil {
		t.Fatalf("Expected no error on first create, got: %v", err)
	}

	_, err = server.CreateGame(gameID, models.HumanVsAi, models.Fafo, models.Fafo)
	if err == nil {
		t.Error("Expected error when creating game with duplicate ID, got nil")
	}
}

func TestCreateGameInvalidType(t *testing.T) {
	server := api.NewGameServer()
	_, err := server.CreateGame("invalid-type-game", "InvalidGameType", models.Fafo, models.Fafo)

	if err == nil {
		t.Error("Expected error for invalid game type, got nil")
	}
}

func TestGetSession(t *testing.T) {
	server := api.NewGameServer()
	gameID := "get-session-test"

	handler, err := server.CreateGame(gameID, models.HumanVsAi, models.Fafo, models.Fafo)
	if err != nil {
		t.Fatalf("Failed to create game: %v", err)
	}

	retrieved, exists := server.GetSession(gameID)
	if !exists {
		t.Error("Expected session to exist after creation")
	}

	if retrieved != handler {
		t.Error("Expected retrieved handler to match created handler")
	}
}

func TestGetSessionNonExistent(t *testing.T) {
	server := api.NewGameServer()
	_, exists := server.GetSession("non-existent-game")

	if exists {
		t.Error("Expected session to not exist for non-existent game ID")
	}
}
