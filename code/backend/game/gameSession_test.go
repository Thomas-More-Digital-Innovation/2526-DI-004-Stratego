package game_test

import (
	"digital-innovation/stratego/engine"
	"digital-innovation/stratego/game"
	"testing"
	"time"
)

func TestNewGameSession(t *testing.T) {
	player1 := engine.NewPlayer(0, "Human", "red")
	player2 := engine.NewPlayer(1, "AI", "blue")

	controller1 := engine.NewHumanPlayerController(&player1)
	controller2 := engine.NewHumanPlayerController(&player2)

	session := game.NewGameSession("test-session", controller1, controller2)

	if session == nil {
		t.Fatal("Expected NewGameSession to return a session, but got nil")
	}

	if session.ID != "test-session" {
		t.Errorf("Expected session ID to be 'test-session', got: %s", session.ID)
	}

	if !session.IsSetupPhase() {
		t.Error("Expected new session to be in setup phase")
	}

	if session.IsRunning() {
		t.Error("Expected new session to not be running")
	}
}

func TestGameSessionSwapSetupPieces(t *testing.T) {
	player1 := engine.NewPlayer(0, "Player1", "red")
	player2 := engine.NewPlayer(1, "Player2", "blue")

	controller1 := engine.NewHumanPlayerController(&player1)
	controller2 := engine.NewHumanPlayerController(&player2)

	session := game.NewGameSession("swap-test", controller1, controller2)

	pos1 := engine.NewPosition(0, 6)
	pos2 := engine.NewPosition(1, 6)

	err := session.SwapSetupPieces(0, pos1, pos2)
	if err != nil {
		t.Errorf("Expected no error swapping pieces, got: %v", err)
	}
}

func TestGameSessionSwapSetupPiecesNotInSetup(t *testing.T) {
	player1 := engine.NewPlayer(0, "Player1", "red")
	player2 := engine.NewPlayer(1, "Player2", "blue")

	controller1 := engine.NewHumanPlayerController(&player1)
	controller2 := engine.NewHumanPlayerController(&player2)

	session := game.NewGameSession("swap-test-2", controller1, controller2)
	session.SetSetupPhaseComplete()

	pos1 := engine.NewPosition(0, 6)
	pos2 := engine.NewPosition(1, 6)

	err := session.SwapSetupPieces(0, pos1, pos2)
	if err == nil {
		t.Error("Expected error swapping pieces when not in setup phase")
	}
}

func TestGameSessionSwapSetupPiecesInvalidPlayer(t *testing.T) {
	player1 := engine.NewPlayer(0, "Player1", "red")
	player2 := engine.NewPlayer(1, "Player2", "blue")

	controller1 := engine.NewHumanPlayerController(&player1)
	controller2 := engine.NewHumanPlayerController(&player2)

	session := game.NewGameSession("swap-test-3", controller1, controller2)

	pos1 := engine.NewPosition(0, 6)
	pos2 := engine.NewPosition(1, 6)

	err := session.SwapSetupPieces(99, pos1, pos2)
	if err == nil {
		t.Error("Expected error swapping pieces for invalid player ID")
	}
}

func TestGameSessionRandomizeSetup(t *testing.T) {
	player1 := engine.NewPlayer(0, "Player1", "red")
	player2 := engine.NewPlayer(1, "Player2", "blue")

	controller1 := engine.NewHumanPlayerController(&player1)
	controller2 := engine.NewHumanPlayerController(&player2)

	session := game.NewGameSession("randomize-test", controller1, controller2)

	err := session.RandomizeSetup(0)
	if err != nil {
		t.Errorf("Expected no error randomizing setup, got: %v", err)
	}

	pieces := session.GetSetupPieces(0)
	if len(pieces) != 40 {
		t.Errorf("Expected 40 pieces after randomization, got: %d", len(pieces))
	}
}

func TestGameSessionRandomizeSetupNotInSetupPhase(t *testing.T) {
	player1 := engine.NewPlayer(0, "Player1", "red")
	player2 := engine.NewPlayer(1, "Player2", "blue")

	controller1 := engine.NewHumanPlayerController(&player1)
	controller2 := engine.NewHumanPlayerController(&player2)

	session := game.NewGameSession("randomize-test-2", controller1, controller2)
	session.SetSetupPhaseComplete()

	err := session.RandomizeSetup(0)
	if err == nil {
		t.Error("Expected error randomizing setup when not in setup phase")
	}
}

func TestGameSessionStartGameFromSetup(t *testing.T) {
	player1 := engine.NewPlayer(0, "Player1", "red")
	player2 := engine.NewPlayer(1, "Player2", "blue")

	controller1 := engine.NewHumanPlayerController(&player1)
	controller2 := engine.NewHumanPlayerController(&player2)

	session := game.NewGameSession("start-test", controller1, controller2)

	err := session.StartGameFromSetup()
	if err != nil {
		t.Errorf("Expected no error starting game from setup, got: %v", err)
	}

	if session.IsSetupPhase() {
		t.Error("Expected session to not be in setup phase after starting")
	}

	if !session.IsRunning() {
		t.Error("Expected session to be running after starting")
	}

	// Clean up
	session.Stop()
	time.Sleep(50 * time.Millisecond)
}

func TestGameSessionGetGameState(t *testing.T) {
	player1 := engine.NewPlayer(0, "Player1", "red")
	player2 := engine.NewPlayer(1, "Player2", "blue")

	controller1 := engine.NewHumanPlayerController(&player1)
	controller2 := engine.NewHumanPlayerController(&player2)

	session := game.NewGameSession("state-test", controller1, controller2)

	state := session.GetGameState()

	if state.Round != 1 {
		t.Errorf("Expected round 1, got: %d", state.Round)
	}

	if state.CurrentPlayerID != 0 {
		t.Errorf("Expected current player ID 0, got: %d", state.CurrentPlayerID)
	}

	if !state.IsSetupPhase {
		t.Error("Expected IsSetupPhase to be true")
	}

	if state.IsGameOver {
		t.Error("Expected IsGameOver to be false")
	}
}

func TestGameSessionSubmitMove(t *testing.T) {
	player1 := engine.NewPlayer(0, "Human", "red")
	player2 := engine.NewPlayer(1, "AI", "blue")

	controller1 := engine.NewHumanPlayerController(&player1)
	controller2 := engine.NewHumanPlayerController(&player2)

	session := game.NewGameSession("move-test", controller1, controller2)

	// Start the game
	err := session.StartGameFromSetup()
	if err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}

	time.Sleep(100 * time.Millisecond) // Wait for game to be ready

	// Submit a valid move
	move := engine.NewMove(engine.NewPosition(0, 6), engine.NewPosition(0, 5), &player1)
	err = session.SubmitMove(0, move)
	if err != nil {
		t.Errorf("Expected no error submitting move, got: %v", err)
	}

	// Clean up
	session.Stop()
	time.Sleep(50 * time.Millisecond)
}

func TestGameSessionSubmitMoveNotRunning(t *testing.T) {
	player1 := engine.NewPlayer(0, "Human", "red")
	player2 := engine.NewPlayer(1, "AI", "blue")

	controller1 := engine.NewHumanPlayerController(&player1)
	controller2 := engine.NewHumanPlayerController(&player2)

	session := game.NewGameSession("move-test-2", controller1, controller2)

	// Try to submit move without starting the game
	move := engine.NewMove(engine.NewPosition(0, 6), engine.NewPosition(0, 5), &player1)
	err := session.SubmitMove(0, move)
	if err == nil {
		t.Error("Expected error submitting move when game is not running")
	}
}

func TestGameSessionAnimationSignaling(t *testing.T) {
	player1 := engine.NewPlayer(0, "Player1", "red")
	player2 := engine.NewPlayer(1, "Player2", "blue")

	controller1 := engine.NewHumanPlayerController(&player1)
	controller2 := engine.NewHumanPlayerController(&player2)

	session := game.NewGameSession("anim-test", controller1, controller2)

	if session.IsWaitingForAnimation() {
		t.Error("Expected not waiting for animation initially")
	}

	go func() {
		time.Sleep(50 * time.Millisecond)
		session.SignalAnimationComplete()
	}()

	session.WaitForAnimationComplete(200 * time.Millisecond)

	if session.IsWaitingForAnimation() {
		t.Error("Expected not waiting for animation after signal")
	}
}

func TestGameSessionAnimationTimeout(t *testing.T) {
	player1 := engine.NewPlayer(0, "Player1", "red")
	player2 := engine.NewPlayer(1, "Player2", "blue")

	controller1 := engine.NewHumanPlayerController(&player1)
	controller2 := engine.NewHumanPlayerController(&player2)

	session := game.NewGameSession("anim-timeout-test", controller1, controller2)

	start := time.Now()
	session.WaitForAnimationComplete(50 * time.Millisecond)
	elapsed := time.Since(start)

	if elapsed < 50*time.Millisecond {
		t.Errorf("Expected timeout to be at least 50ms, got: %v", elapsed)
	}
}

func TestGameSessionGetBoard(t *testing.T) {
	player1 := engine.NewPlayer(0, "Player1", "red")
	player2 := engine.NewPlayer(1, "Player2", "blue")

	controller1 := engine.NewHumanPlayerController(&player1)
	controller2 := engine.NewHumanPlayerController(&player2)

	session := game.NewGameSession("board-test", controller1, controller2)

	board := session.GetBoard()
	if board == nil {
		t.Error("Expected GetBoard to return a board, got nil")
	}
}

func TestGameSessionGetSetupPieces(t *testing.T) {
	player1 := engine.NewPlayer(0, "Player1", "red")
	player2 := engine.NewPlayer(1, "Player2", "blue")

	controller1 := engine.NewHumanPlayerController(&player1)
	controller2 := engine.NewHumanPlayerController(&player2)

	session := game.NewGameSession("pieces-test", controller1, controller2)

	pieces1 := session.GetSetupPieces(0)
	pieces2 := session.GetSetupPieces(1)

	if len(pieces1) != 40 {
		t.Errorf("Expected 40 pieces for player 0, got: %d", len(pieces1))
	}

	if len(pieces2) != 40 {
		t.Errorf("Expected 40 pieces for player 1, got: %d", len(pieces2))
	}
}

func TestGameSessionStop(t *testing.T) {
	player1 := engine.NewPlayer(0, "Player1", "red")
	player2 := engine.NewPlayer(1, "Player2", "blue")

	controller1 := engine.NewHumanPlayerController(&player1)
	controller2 := engine.NewHumanPlayerController(&player2)

	session := game.NewGameSession("stop-test", controller1, controller2)

	err := session.StartGameFromSetup()
	if err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	session.Stop()
	time.Sleep(100 * time.Millisecond)

	if session.IsRunning() {
		t.Error("Expected session to stop running after Stop()")
	}
}
