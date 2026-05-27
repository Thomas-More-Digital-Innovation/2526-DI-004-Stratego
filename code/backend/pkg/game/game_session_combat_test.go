package game

import (
	"digital-innovation/gostrategy/internal/models"
	"testing"
	"time"
)

func TestSessionSubmitMove_NoPiece(t *testing.T) {
	player1 := NewPlayer(0, "Human", "red")
	player2 := NewPlayer(1, "AI", "blue")
	controller1 := NewHumanPlayerController(&player1)
	controller2 := NewHumanPlayerController(&player2)
	session := NewSession("move-no-piece", controller1, controller2)

	if err := session.StartGameFromSetup(false); err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}
	defer session.Stop()
	time.Sleep(50 * time.Millisecond)

	// Clear a spot
	pos := NewPosition(5, 5)
	session.GetBoard().SetPieceAt(pos, nil)

	move := NewMove(pos, NewPosition(5, 4), &player1)
	err := session.SubmitMove(0, move)
	if err == nil || err.Error() != "no piece at source position" {
		t.Errorf("Expected 'no piece at source position' error, got: %v", err)
	}
}

func TestSessionSubmitMove_WrongOwner(t *testing.T) {
	player1 := NewPlayer(0, "Human", "red")
	player2 := NewPlayer(1, "AI", "blue")
	controller1 := NewHumanPlayerController(&player1)
	controller2 := NewHumanPlayerController(&player2)
	session := NewSession("move-wrong-owner", controller1, controller2)

	if err := session.StartGameFromSetup(false); err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}
	defer session.Stop()
	time.Sleep(50 * time.Millisecond)

	// Find an opponent piece (Player 2 is at top rows 0-3)
	pos := NewPosition(0, 0)
	move := NewMove(pos, NewPosition(0, 1), &player1)
	err := session.SubmitMove(0, move)
	if err == nil || err.Error() != "piece at source position does not belong to current player" {
		t.Errorf("Expected ownership error, got: %v", err)
	}
}

func TestSessionSubmitMove_UnmovablePiece(t *testing.T) {
	player1 := NewPlayer(0, "Human", "red")
	player2 := NewPlayer(1, "AI", "blue")
	controller1 := NewHumanPlayerController(&player1)
	controller2 := NewHumanPlayerController(&player2)
	session := NewSession("move-unmovable", controller1, controller2)

	if err := session.StartGameFromSetup(false); err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}
	defer session.Stop()
	time.Sleep(50 * time.Millisecond)

	// Place a flag for player 1
	pos := NewPosition(0, 9)
	flag := NewPiece(models.Flag, &player1)
	session.GetBoard().SetPieceAt(pos, flag)

	move := NewMove(pos, NewPosition(0, 8), &player1)
	err := session.SubmitMove(0, move)
	if err == nil || err.Error() != "no movable piece at the given position" {
		t.Errorf("Expected 'no movable piece' error, got: %v", err)
	}
}

func TestSessionSubmitMove_IllegalMove(t *testing.T) {
	player1 := NewPlayer(0, "Human", "red")
	player2 := NewPlayer(1, "AI", "blue")
	controller1 := NewHumanPlayerController(&player1)
	controller2 := NewHumanPlayerController(&player2)
	session := NewSession("move-illegal", controller1, controller2)

	if err := session.StartGameFromSetup(false); err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}
	defer session.Stop()
	time.Sleep(50 * time.Millisecond)

	// Place a non-scout for player 1
	pos := NewPosition(1, 9)
	marshal := NewPiece(models.Marshal, &player1)
	session.GetBoard().SetPieceAt(pos, marshal)

	move := NewMove(pos, NewPosition(1, 7), &player1)
	err := session.SubmitMove(0, move)
	if err == nil || err.Error() != "illegal move for this piece" {
		t.Errorf("Expected 'illegal move' error, got: %v", err)
	}
}

func TestSessionAnimationSignaling(t *testing.T) {
	player1 := NewPlayer(0, "Player1", "red")
	player2 := NewPlayer(1, "Player2", "blue")

	controller1 := NewHumanPlayerController(&player1)
	controller2 := NewHumanPlayerController(&player2)

	session := NewSession("anim-test", controller1, controller2)

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

func TestSessionAnimationTimeout(t *testing.T) {
	player1 := NewPlayer(0, "Player1", "red")
	player2 := NewPlayer(1, "Player2", "blue")

	controller1 := NewHumanPlayerController(&player1)
	controller2 := NewHumanPlayerController(&player2)

	session := NewSession("anim-timeout-test", controller1, controller2)

	start := time.Now()
	session.WaitForAnimationComplete(50 * time.Millisecond)
	elapsed := time.Since(start)

	if elapsed < 50*time.Millisecond {
		t.Errorf("Expected timeout to be at least 50ms, got: %v", elapsed)
	}
}

func TestSessionGetBoard(t *testing.T) {
	player1 := NewPlayer(0, "Player1", "red")
	player2 := NewPlayer(1, "Player2", "blue")

	controller1 := NewHumanPlayerController(&player1)
	controller2 := NewHumanPlayerController(&player2)

	session := NewSession("board-test", controller1, controller2)

	board := session.GetBoard()
	if board == nil {
		t.Error("Expected GetBoard to return a board, got nil")
	}
}

func TestSessionGetSetupPieces(t *testing.T) {
	player1 := NewPlayer(0, "Player1", "red")
	player2 := NewPlayer(1, "Player2", "blue")

	controller1 := NewHumanPlayerController(&player1)
	controller2 := NewHumanPlayerController(&player2)

	session := NewSession("pieces-test", controller1, controller2)

	pieces1 := session.GetSetupPieces(0)
	pieces2 := session.GetSetupPieces(1)

	if len(pieces1) != 40 {
		t.Errorf("Expected 40 pieces for player 0, got: %d", len(pieces1))
	}

	if len(pieces2) != 40 {
		t.Errorf("Expected 40 pieces for player 1, got: %d", len(pieces2))
	}
}

func TestSessionStop(t *testing.T) {
	player1 := NewPlayer(0, "Player1", "red")
	player2 := NewPlayer(1, "Player2", "blue")

	controller1 := NewHumanPlayerController(&player1)
	controller2 := NewHumanPlayerController(&player2)

	session := NewSession("stop-test", controller1, controller2)

	err := session.StartGameFromSetup(false)
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
