package game

import (
	"digital-innovation/gostrategy/pkg/game/models"
	"testing"
	"time"
)

func TestNewSession(t *testing.T) {
	player1 := NewPlayer(0, "Human", "red")
	player2 := NewPlayer(1, "AI", "blue")

	controller1 := NewHumanPlayerController(&player1)
	controller2 := NewHumanPlayerController(&player2)

	session := NewSession("test-session", controller1, controller2)

	if session == nil {
		t.Fatal("Expected NewSession to return a session, but got nil")
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

func TestSessionSwapSetupPieces(t *testing.T) {
	player1 := NewPlayer(0, "Player1", "red")
	player2 := NewPlayer(1, "Player2", "blue")

	controller1 := NewHumanPlayerController(&player1)
	controller2 := NewHumanPlayerController(&player2)

	session := NewSession("swap-test", controller1, controller2)

	pos1 := NewPosition(0, 6)
	pos2 := NewPosition(1, 6)

	err := session.SwapSetupPieces(0, pos1, pos2)
	if err != nil {
		t.Errorf("Expected no error swapping pieces, got: %v", err)
	}
}

func TestSessionSwapSetupPiecesNotInSetup(t *testing.T) {
	player1 := NewPlayer(0, "Player1", "red")
	player2 := NewPlayer(1, "Player2", "blue")

	controller1 := NewHumanPlayerController(&player1)
	controller2 := NewHumanPlayerController(&player2)

	session := NewSession("swap-test-2", controller1, controller2)
	session.SetSetupPhaseComplete()

	pos1 := NewPosition(0, 6)
	pos2 := NewPosition(1, 6)

	err := session.SwapSetupPieces(0, pos1, pos2)
	if err == nil {
		t.Error("Expected error swapping pieces when not in setup phase")
	}
}

func TestSessionSwapSetupPiecesInvalidPlayer(t *testing.T) {
	player1 := NewPlayer(0, "Player1", "red")
	player2 := NewPlayer(1, "Player2", "blue")

	controller1 := NewHumanPlayerController(&player1)
	controller2 := NewHumanPlayerController(&player2)

	session := NewSession("swap-test-3", controller1, controller2)

	pos1 := NewPosition(0, 6)
	pos2 := NewPosition(1, 6)

	err := session.SwapSetupPieces(99, pos1, pos2)
	if err == nil {
		t.Error("Expected error swapping pieces for invalid player ID")
	}
}

func TestSessionRandomizeSetup(t *testing.T) {
	player1 := NewPlayer(0, "Player1", "red")
	player2 := NewPlayer(1, "Player2", "blue")

	controller1 := NewHumanPlayerController(&player1)
	controller2 := NewHumanPlayerController(&player2)

	session := NewSession("randomize-test", controller1, controller2)

	err := session.RandomizeSetup(0)
	if err != nil {
		t.Errorf("Expected no error randomizing setup, got: %v", err)
	}

	pieces := session.GetSetupPieces(0)
	if len(pieces) != 40 {
		t.Errorf("Expected 40 pieces after randomization, got: %d", len(pieces))
	}
}

func TestSessionRandomizeSetupNotInSetupPhase(t *testing.T) {
	player1 := NewPlayer(0, "Player1", "red")
	player2 := NewPlayer(1, "Player2", "blue")

	controller1 := NewHumanPlayerController(&player1)
	controller2 := NewHumanPlayerController(&player2)

	session := NewSession("randomize-test-2", controller1, controller2)
	session.SetSetupPhaseComplete()

	err := session.RandomizeSetup(0)
	if err == nil {
		t.Error("Expected error randomizing setup when not in setup phase")
	}
}

func TestSessionStartGameFromSetup(t *testing.T) {
	player1 := NewPlayer(0, "Player1", "red")
	player2 := NewPlayer(1, "Player2", "blue")

	controller1 := NewHumanPlayerController(&player1)
	controller2 := NewHumanPlayerController(&player2)

	session := NewSession("start-test", controller1, controller2)

	err := session.StartGameFromSetup(false)
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

func TestSessionGetGameState(t *testing.T) {
	player1 := NewPlayer(0, "Player1", "red")
	player2 := NewPlayer(1, "Player2", "blue")

	controller1 := NewHumanPlayerController(&player1)
	controller2 := NewHumanPlayerController(&player2)

	session := NewSession("state-test", controller1, controller2)

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

func TestSessionSubmitMove(t *testing.T) {
	player1 := NewPlayer(0, "Human", "red")
	player2 := NewPlayer(1, "AI", "blue")

	controller1 := NewHumanPlayerController(&player1)
	controller2 := NewHumanPlayerController(&player2)

	session := NewSession("move-test", controller1, controller2)

	// Start the game
	err := session.StartGameFromSetup(false)
	if err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}

	time.Sleep(100 * time.Millisecond) // Wait for game to be ready

	scout := NewPiece(models.Scout, &player1)
	session.GetBoard().SetPieceAt(NewPosition(0, 6), scout)

	// Submit a valid move
	move := NewMove(NewPosition(0, 6), NewPosition(0, 5), &player1)
	err = session.SubmitMove(0, move)
	if err != nil {
		t.Errorf("Expected no error submitting move, got: %v", err)
	}

	// Clean up
	session.Stop()
	time.Sleep(50 * time.Millisecond)
}

func TestSessionSubmitMoveNotRunning(t *testing.T) {
	player1 := NewPlayer(0, "Human", "red")
	player2 := NewPlayer(1, "AI", "blue")

	controller1 := NewHumanPlayerController(&player1)
	controller2 := NewHumanPlayerController(&player2)

	session := NewSession("move-test-2", controller1, controller2)

	// Try to submit move without starting the game
	move := NewMove(NewPosition(0, 6), NewPosition(0, 5), &player1)
	err := session.SubmitMove(0, move)
	if err == nil {
		t.Error("Expected error submitting move when game is not running")
	}
}
