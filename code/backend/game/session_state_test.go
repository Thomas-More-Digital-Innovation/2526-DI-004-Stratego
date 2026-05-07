package game

import (
	"digital-innovation/stratego/engine"
	"testing"
	"time"
)

func TestSession_SetupTimeout(t *testing.T) {
	player1 := engine.NewPlayer(0, "P1", "red")
	player2 := engine.NewPlayer(1, "P2", "blue")
	c1 := engine.NewHumanPlayerController(&player1)
	c2 := engine.NewHumanPlayerController(&player2)

	// Use functional option to set a short timeout for testing
	session := NewSession("test-timeout", c1, c2, WithSetupTimeout(100*time.Millisecond))

	if !session.IsSetupPhase() {
		t.Error("Expected game to start in setup phase")
	}

	// Wait for timeout to expire
	time.Sleep(200 * time.Millisecond)

	if session.IsSetupPhase() {
		t.Error("Expected setup phase to be complete after timeout")
	}
	if !session.IsRunning() {
		t.Error("Expected game to be running after timeout")
	}
}

func TestSession_Abort(t *testing.T) {
	player1 := engine.NewPlayer(0, "P1", "red")
	player2 := engine.NewPlayer(1, "P2", "blue")
	c1 := engine.NewHumanPlayerController(&player1)
	c2 := engine.NewHumanPlayerController(&player2)
	session := NewSession("test-abort", c1, c2)

	if err := session.StartGameFromSetup(true); err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}

	if !session.IsRunning() {
		t.Fatal("Game should be running")
	}

	session.Stop()

	if !session.IsAborted() {
		t.Error("Expected session to be marked as aborted")
	}

	// Wait a bit for the runner to stop
	time.Sleep(50 * time.Millisecond)

	if session.IsRunning() {
		t.Error("Expected game to stop running after Stop()")
	}
}

func TestSession_SwapSetupPieces(t *testing.T) {
	player1 := engine.NewPlayer(0, "P1", "red")
	player2 := engine.NewPlayer(1, "P2", "blue")
	c1 := engine.NewHumanPlayerController(&player1)
	c2 := engine.NewHumanPlayerController(&player2)
	session := NewSession("test-swap", c1, c2)

	p1 := engine.NewPosition(0, 6)
	p2 := engine.NewPosition(1, 6)

	piecesBefore := make([]*engine.Piece, len(session.player1Pieces))
	copy(piecesBefore, session.player1Pieces)

	err := session.SwapSetupPieces(0, p1, p2)
	if err != nil {
		t.Fatalf("Swap failed: %v", err)
	}

	if session.player1Pieces[0] != piecesBefore[1] || session.player1Pieces[1] != piecesBefore[0] {
		t.Error("Pieces were not swapped correctly at indices 0 and 1")
	}

	// Test invalid player
	err = session.SwapSetupPieces(2, p1, p2)
	if err == nil {
		t.Error("Expected error for invalid player ID")
	}

	// Test invalid position (outside setup area)
	err = session.SwapSetupPieces(0, p1, engine.NewPosition(0, 0))
	if err == nil {
		t.Error("Expected error for position outside setup area")
	}
}
func TestSession_SetupWarning(t *testing.T) {
	player1 := engine.NewPlayer(0, "P1", "red")
	player2 := engine.NewPlayer(1, "P2", "blue")
	c1 := engine.NewHumanPlayerController(&player1)
	c2 := engine.NewHumanPlayerController(&player2)

	// Set warning at 50ms, timeout at 200ms
	session := NewSession("test-warning", c1, c2,
		WithSetupTimeout(200*time.Millisecond),
		WithSetupWarning(50*time.Millisecond),
	)

	// We expect a notification to be sent to moveNotifyChan when the warning fires
	// We'll use WaitForMoveNotification to check for it
	received := session.WaitForMoveNotification(150 * time.Millisecond)
	if !received {
		t.Error("Expected to receive a warning notification, but timed out")
	}

	if !session.IsSetupPhase() {
		t.Error("Game should still be in setup phase after warning but before timeout")
	}
}
