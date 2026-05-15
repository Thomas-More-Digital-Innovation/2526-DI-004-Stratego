package game_test

import (
	"digital-innovation/gostrategy/engine"
	"digital-innovation/gostrategy/game"
	"digital-innovation/gostrategy/models"
	"testing"
	"time"
)

func TestSessionComprehensive(t *testing.T) {
	player1 := engine.NewPlayer(0, "P1", "")
	player2 := engine.NewPlayer(1, "P2", "")
	c1 := engine.NewHumanPlayerController(&player1)
	c2 := engine.NewHumanPlayerController(&player2)

	// Test WithSetupTimeout and WithSetupWarning
	session := game.NewSession("comp-test", c1, c2,
		game.WithSetupTimeout(10*time.Second),
		game.WithSetupWarning(5*time.Second),
	)

	// Test Player association
	session.SetPlayer1Associate(123, "Alice")
	session.SetPlayer2Associate(456, "Bob")
	p1, p2 := session.GetPlayerIDs()
	if *p1 != 123 || *p2 != 456 {
		t.Errorf("Player IDs mismatch: %d, %d", *p1, *p2)
	}

	// Test LoadSetup
	setupData := []byte("0BBBBBB12222222233333444455556666777889M")
	err := session.LoadSetup(0, setupData)
	if err != nil {
		t.Errorf("LoadSetup failed: %v", err)
	}
	err = session.LoadSetup(99, setupData)
	if err == nil {
		t.Error("LoadSetup should fail for invalid player ID")
	}

	// Test NotifySetupUpdate
	session.NotifySetupUpdate()

	// Test GetSetupCompleteChan
	ch := session.GetSetupCompleteChan()
	if ch == nil {
		t.Error("GetSetupCompleteChan returned nil")
	}

	// Test IsAbortedChan
	abCh := session.IsAbortedChan()
	if abCh == nil {
		t.Error("IsAbortedChan returned nil")
	}

	// Test StartGameFromSetup with headless=true
	err = session.StartGameFromSetup(true)
	if err != nil {
		t.Errorf("StartGameFromSetup failed: %v", err)
	}

	if !session.IsHeadless() {
		t.Error("Expected headless session")
	}

	// Test Move notifications
	session.NotifyMoveExecuted()
	if !session.WaitForMoveNotification(100 * time.Millisecond) {
		t.Error("WaitForMoveNotification timed out")
	}

	// Test Move Ack
	session.AckMoveProcessed()
	if !session.WaitForMoveAck(100 * time.Millisecond) {
		t.Error("WaitForMoveAck timed out")
	}

	// Additional Session methods
	session.Pause()
	session.Unpause()
	session.SetTurnDelay(100 * time.Millisecond)
	session.StepAI()

	session.GetGameSummary("classic")
	_, _ = session.GetAvailableMoves(0, engine.NewPosition(0, 6))
	session.GetWinner()
	session.GetWinCause()
	session.SetWinner(&player1, game.WinCauseFlagCaptured)
	session.GetLastCombat()
	session.GetLastHistoricalMove()
	session.ClearLastCombat()
	session.HideCombatPieces()
	session.GetGame()

	// Test SubmitMove error paths
	marshal := engine.NewPiece(models.Marshal, &player1)
	session.GetBoard().SetPieceAt(engine.NewPosition(0, 6), marshal)

	// Verify available moves for the piece we just placed
	moves, err := session.GetAvailableMoves(0, engine.NewPosition(0, 6))
	if err != nil || len(moves) == 0 {
		t.Errorf("Expected available moves for Marshal, got error: %v, count: %d", err, len(moves))
	}

	move := engine.NewMove(engine.NewPosition(0, 6), engine.NewPosition(0, 5), &player1)

	// Wrong player
	err = session.SubmitMove(1, move)
	if err == nil || err.Error() != "not your turn" {
		t.Errorf("Expected 'not your turn' error, got %v", err)
	}

	// No piece at source
	moveEmpty := engine.NewMove(engine.NewPosition(5, 5), engine.NewPosition(5, 4), &player1)
	err = session.SubmitMove(0, moveEmpty)
	if err == nil || err.Error() != "no piece at source position" {
		t.Errorf("Expected 'no piece at source position' error, got %v", err)
	}

	// Piece belongs to opponent
	oppPiece := engine.NewPiece(models.Scout, &player2)
	session.GetBoard().SetPieceAt(engine.NewPosition(5, 5), oppPiece)
	moveOpp := engine.NewMove(engine.NewPosition(5, 5), engine.NewPosition(5, 4), &player1)
	err = session.SubmitMove(0, moveOpp)
	if err == nil || err.Error() != "piece at source position does not belong to current player" {
		t.Errorf("Expected ownership error, got %v", err)
	}

	// Illegal move
	moveIllegal := engine.NewMove(engine.NewPosition(0, 6), engine.NewPosition(9, 9), &player1)
	err = session.SubmitMove(0, moveIllegal)
	if err == nil {
		t.Error("Expected error for illegal move")
	}
	// Test WaitForCompletion (async)
	go func() {
		time.Sleep(50 * time.Millisecond)
		session.SetWinner(&player1, game.WinCauseFlagCaptured)
	}()
	session.WaitForCompletion()

	// Test GetLastHistoricalMove when empty
	emptySession := game.NewSession("empty", c1, c2)
	if emptySession.GetLastHistoricalMove() != nil {
		t.Error("Expected nil historical move for new session")
	}

	session.Stop()
}

func TestRunnerComprehensive(t *testing.T) {
	player1 := engine.NewPlayer(0, "P1", "")
	player2 := engine.NewPlayer(1, "P2", "")
	c1 := engine.NewHumanPlayerController(&player1)
	c2 := engine.NewHumanPlayerController(&player2)
	g := game.NewGame(c1, c2)
	runner := game.NewRunner(g, 0, 100)

	// Test SetMoveCallback
	called := false
	runner.SetMoveCallback(func() {
		called = true
	})

	// Test Pause/Unpause
	runner.Pause()
	if !runner.IsPaused() {
		t.Error("Runner should be paused")
	}
	runner.Unpause()
	if runner.IsPaused() {
		t.Error("Runner should be unpaused")
	}

	// Test Step
	runner.Pause()
	// Place a piece to move
	scout := engine.NewPiece(models.Scout, &player1)
	g.Board.SetPieceAt(engine.NewPosition(0, 6), scout)
	c1.SetPendingMove(engine.NewMove(engine.NewPosition(0, 6), engine.NewPosition(0, 5), &player1))

	runner.DebugSetWaitingForInput(true)
	if !runner.Step() {
		t.Error("Step should execute turn even if paused")
	}
	if !called {
		t.Error("Move callback should have been called")
	}
}
