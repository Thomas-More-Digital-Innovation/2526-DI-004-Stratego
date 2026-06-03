package game

import (
	"digital-innovation/gostrategy/pkg/game/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionComprehensive(t *testing.T) {
	player1 := NewPlayer(0, "P1", "")
	player2 := NewPlayer(1, "P2", "")
	c1 := NewHumanPlayerController(&player1)
	c2 := NewHumanPlayerController(&player2)

	// Test WithSetupTimeout and WithSetupWarning
	session := NewSession("comp-test", c1, c2,
		WithSetupTimeout(10*time.Second),
		WithSetupWarning(5*time.Second),
	)

	// Test Player association
	session.SetPlayer1Associate(123, "Alice")
	session.SetPlayer2Associate(456, "Bob")
	p1, p2 := session.GetPlayerIDs()
	require.NotNil(t, p1)
	require.NotNil(t, p2)
	assert.Equal(t, 123, *p1)
	assert.Equal(t, 456, *p2)

	// Test LoadSetup
	setupData := []byte("0BBBBBB12222222233333444455556666777889M")
	err := session.LoadSetup(0, setupData)
	assert.NoError(t, err)

	err = session.LoadSetup(99, setupData)
	assert.Error(t, err)

	// Test NotifySetupUpdate
	session.NotifySetupUpdate()

	// Test GetSetupCompleteChan
	ch := session.GetSetupCompleteChan()
	assert.NotNil(t, ch)

	// Test IsAbortedChan
	abCh := session.IsAbortedChan()
	assert.NotNil(t, abCh)

	// Test StartGameFromSetup with headless=true
	err = session.StartGameFromSetup(true)
	assert.NoError(t, err)
	assert.True(t, session.IsHeadless())

	// Test Move notifications
	session.NotifyMoveExecuted()
	assert.True(t, session.WaitForMoveNotification(100*time.Millisecond))

	// Test Move Ack
	session.AckMoveProcessed()
	assert.True(t, session.WaitForMoveAck(100*time.Millisecond))

	// Additional Session methods
	session.Pause()
	session.Unpause()
	session.SetTurnDelay(100 * time.Millisecond)
	session.StepAI()

	session.GetGameSummary("classic")
	_, _ = session.GetAvailableMoves(0, NewPosition(0, 6))
	session.GetLastCombat()
	session.GetLastHistoricalMove()
	session.ClearLastCombat()
	session.HideCombatPieces()
	session.GetGame()

	// Test SubmitMove error paths
	marshal := NewPiece(models.Marshal, &player1)
	session.GetBoard().SetPieceAt(NewPosition(0, 6), marshal)

	// Verify available moves for the piece we just placed
	moves, err := session.GetAvailableMoves(0, NewPosition(0, 6))
	require.NoError(t, err)
	assert.NotEmpty(t, moves)

	move := NewMove(NewPosition(0, 6), NewPosition(0, 5), &player1)

	// Wrong player
	err = session.SubmitMove(1, move)
	assert.Error(t, err)
	assert.Equal(t, "not your turn", err.Error())

	// No piece at source
	moveEmpty := NewMove(NewPosition(5, 5), NewPosition(5, 4), &player1)
	err = session.SubmitMove(0, moveEmpty)
	assert.Error(t, err)
	assert.Equal(t, "no piece at source position", err.Error())

	// Piece belongs to opponent
	oppPiece := NewPiece(models.Scout, &player2)
	session.GetBoard().SetPieceAt(NewPosition(5, 5), oppPiece)
	moveOpp := NewMove(NewPosition(5, 5), NewPosition(5, 4), &player1)
	err = session.SubmitMove(0, moveOpp)
	assert.Error(t, err)
	assert.Equal(t, "piece at source position does not belong to current player", err.Error())

	// Illegal move
	moveIllegal := NewMove(NewPosition(0, 6), NewPosition(9, 9), &player1)
	err = session.SubmitMove(0, moveIllegal)
	assert.Error(t, err)

	// Smoke tests for game end state methods
	session.GetWinner()
	session.GetWinCause()
	session.SetWinner(&player1, WinCauseFlagCaptured)
	// Test WaitForCompletion (async)
	go func() {
		time.Sleep(50 * time.Millisecond)
		session.SetWinner(&player1, WinCauseFlagCaptured)
	}()
	session.WaitForCompletion()

	// Test GetLastHistoricalMove when empty
	emptySession := NewSession("empty", c1, c2)
	assert.Nil(t, emptySession.GetLastHistoricalMove())

	session.Stop()
}

func TestRunnerComprehensive(t *testing.T) {
	player1 := NewPlayer(0, "P1", "")
	player2 := NewPlayer(1, "P2", "")
	c1 := NewHumanPlayerController(&player1)
	c2 := NewHumanPlayerController(&player2)
	g := NewGame(c1, c2)
	runner := NewRunner(g, 0, 100)

	// Test SetMoveCallback
	called := false
	runner.SetMoveCallback(func() {
		called = true
	})

	// Test Pause/Unpause
	runner.Pause()
	assert.True(t, runner.IsPaused())
	runner.Unpause()
	assert.False(t, runner.IsPaused())

	// Test Step
	runner.Pause()
	// Place a piece to move
	scout := NewPiece(models.Scout, &player1)
	g.Board.SetPieceAt(NewPosition(0, 6), scout)
	c1.SetPendingMove(NewMove(NewPosition(0, 6), NewPosition(0, 5), &player1))

	runner.DebugSetWaitingForInput(true)
	assert.True(t, runner.Step())
	assert.True(t, called)
}
