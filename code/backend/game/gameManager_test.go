package game_test

import (
	"digital-innovation/stratego/engine"
	"digital-innovation/stratego/game"
	"testing"
)

func TestNewGameManager(t *testing.T) {
	gm := game.NewGameManager()

	if gm == nil {
		t.Error("Expected GameManager instance, got nil")
	}

	if len(gm.Games) != 0 {
		t.Errorf("Expected 0 games, got %d", len(gm.Games))
	}
}

func TestRunMatch(t *testing.T) {
	gm := game.NewGameManager()
	player1 := engine.NewPlayer(1, "Alice", "red")
	controller1 := engine.NewHumanPlayerController(&player1)
	player2 := engine.NewPlayer(2, "Bob", "blue")
	controller2 := engine.NewHumanPlayerController(&player2)

	gameInstance := gm.RunMatch(controller1, controller2)

	if gameInstance == nil {
		t.Error("Expected Game instance, got nil")
	}

	if len(gm.Games) != 1 {
		t.Errorf("Expected 1 game, got %d", len(gm.Games))
	}

	if gm.ActiveGame != 0 {
		t.Errorf("Expected ActiveGame index to be 0, got %d", gm.ActiveGame)
	}

	if gameInstance.CurrentPlayer != &player1 {
		t.Errorf("Expected current player to be player1, but got %v", gameInstance.CurrentPlayer)
	}

	if gameInstance.Board == nil {
		t.Errorf("Expected a board to be created, but got nil")
	}

	if gameInstance.GetRound() != 1 {
		t.Errorf("Expected round to be 1 at game start, but got %d", gameInstance.GetRound())
	}

	if len(gameInstance.MoveHistory) != 0 {
		t.Errorf("Expected move history to be empty, but got %d moves", len(gameInstance.MoveHistory))
	}
}

func TestRunTournament(t *testing.T) {
	gm := game.NewGameManager()
	player1 := engine.NewPlayer(1, "Alice", "red")
	controller1 := engine.NewHumanPlayerController(&player1)
	player2 := engine.NewPlayer(2, "Bob", "blue")
	controller2 := engine.NewHumanPlayerController(&player2)

	rounds := 3

	games := gm.RunTournament(rounds, controller1, controller2)

	if len(games) != rounds {
		t.Errorf("Expected %d games in tournament, got %d", rounds, len(games))
	}

	if gm.GetGameCount() != rounds {
		t.Errorf("Expected GameManager to have %d games, got %d", rounds, gm.GetGameCount())
	}

	if gm.ActiveGame != rounds-1 {
		t.Errorf("Expected ActiveGame index to be %d, got %d", rounds-1, gm.ActiveGame)
	}
}

func TestGetActiveGame(t *testing.T) {
	gm := game.NewGameManager()
	player1 := engine.NewPlayer(1, "Alice", "red")
	controller1 := engine.NewHumanPlayerController(&player1)
	player2 := engine.NewPlayer(2, "Bob", "blue")
	controller2 := engine.NewHumanPlayerController(&player2)

	if gm.GetActiveGame() != nil {
		t.Error("Expected nil active game when no games exist")
	}

	gm.RunMatch(controller1, controller2)
	activeGame := gm.GetActiveGame()

	if activeGame == nil {
		t.Error("Expected active game instance, got nil")
	}

	if activeGame.CurrentPlayer != &player1 {
		t.Errorf("Expected current player to be player1, but got %v", activeGame.CurrentPlayer)
	}
}

func TestNextGame(t *testing.T) {
	gm := game.NewGameManager()
	player1 := engine.NewPlayer(1, "Alice", "red")
	controller1 := engine.NewHumanPlayerController(&player1)
	player2 := engine.NewPlayer(2, "Bob", "blue")
	controller2 := engine.NewHumanPlayerController(&player2)

	gm.RunMatch(controller1, controller2)
	gm.RunMatch(controller1, controller2)
	gm.RunMatch(controller1, controller2)

	gm.NextGame(1)
	if gm.ActiveGame != 1 {
		t.Errorf("Expected ActiveGame index to be 1, got %d", gm.ActiveGame)
	}

	gm.NextGame(5) // Invalid index
	if gm.ActiveGame != 1 {
		t.Errorf("Expected ActiveGame index to remain 1 after invalid index, got %d", gm.ActiveGame)
	}

	gm.NextGame(-1) // Invalid index
	if gm.ActiveGame != 1 {
		t.Errorf("Expected ActiveGame index to remain 1 after invalid index, got %d", gm.ActiveGame)
	}
}

func TestGetGameCount(t *testing.T) {
	gm := game.NewGameManager()
	player1 := engine.NewPlayer(1, "Alice", "red")
	controller1 := engine.NewHumanPlayerController(&player1)
	player2 := engine.NewPlayer(2, "Bob", "blue")
	controller2 := engine.NewHumanPlayerController(&player2)

	if gm.GetGameCount() != 0 {
		t.Errorf("Expected game count to be 0, got %d", gm.GetGameCount())
	}

	gm.RunMatch(controller1, controller2)
	gm.RunMatch(controller1, controller2)

	if gm.GetGameCount() != 2 {
		t.Errorf("Expected game count to be 2, got %d", gm.GetGameCount())
	}
}
