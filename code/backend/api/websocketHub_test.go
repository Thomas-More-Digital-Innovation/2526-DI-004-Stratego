package api_test

import (
	"digital-innovation/stratego/api"
	"digital-innovation/stratego/engine"
	"digital-innovation/stratego/game"
	"digital-innovation/stratego/models"
	"testing"
)

func TestNewWSHub(t *testing.T) {
	player1 := engine.NewPlayer(0, "Player1", "red")
	player2 := engine.NewPlayer(1, "Player2", "blue")
	controller1 := engine.NewHumanPlayerController(&player1)
	controller2 := engine.NewHumanPlayerController(&player2)
	session := game.NewGameSession("test-hub", controller1, controller2)

	hub := api.NewWSHub(session, models.HumanVsAi)

	if hub == nil {
		t.Fatal("Expected NewWSHub to return a hub, but got nil")
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
			session := game.NewGameSession("test-"+tc.name, controller1, controller2)
			hub := api.NewWSHub(session, tc.gameType)

			if hub == nil {
				t.Errorf("Expected hub to be created for game type %s", tc.gameType)
			}
		})
	}
}
