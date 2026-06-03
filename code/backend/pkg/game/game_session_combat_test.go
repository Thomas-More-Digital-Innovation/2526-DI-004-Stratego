package game

import (
	"digital-innovation/gostrategy/pkg/game/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSessionSubmitMove_NoPiece(t *testing.T) {
	session, player1, _ := setupTestSession("move-no-piece")

	err := session.StartGameFromSetup(false)
	assert.NoError(t, err)
	defer session.Stop()
	time.Sleep(50 * time.Millisecond)

	// Clear a spot
	pos := NewPosition(5, 5)
	session.GetBoard().SetPieceAt(pos, nil)

	move := NewMove(pos, NewPosition(5, 4), player1)
	err = session.SubmitMove(0, move)
	assert.Error(t, err)
	assert.Equal(t, "no piece at source position", err.Error())
}

func TestSessionSubmitMove_WrongOwner(t *testing.T) {
	session, player1, _ := setupTestSession("move-wrong-owner")

	err := session.StartGameFromSetup(false)
	assert.NoError(t, err)
	defer session.Stop()
	time.Sleep(50 * time.Millisecond)

	// Find an opponent piece (Player 2 is at top rows 0-3)
	pos := NewPosition(0, 0)
	move := NewMove(pos, NewPosition(0, 1), player1)
	err = session.SubmitMove(0, move)
	assert.Error(t, err)
	assert.Equal(t, "piece at source position does not belong to current player", err.Error())
}

func TestSessionSubmitMove_UnmovablePiece(t *testing.T) {
	session, player1, _ := setupTestSession("move-unmovable")

	err := session.StartGameFromSetup(false)
	assert.NoError(t, err)
	defer session.Stop()
	time.Sleep(50 * time.Millisecond)

	// Place a flag for player 1
	pos := NewPosition(0, 9)
	flag := NewPiece(models.Flag, player1)
	session.GetBoard().SetPieceAt(pos, flag)

	move := NewMove(pos, NewPosition(0, 8), player1)
	err = session.SubmitMove(0, move)
	assert.Error(t, err)
	assert.Equal(t, "no movable piece at the given position", err.Error())
}

func TestSessionSubmitMove_IllegalMove(t *testing.T) {
	session, player1, _ := setupTestSession("move-illegal")

	err := session.StartGameFromSetup(false)
	assert.NoError(t, err)
	defer session.Stop()
	time.Sleep(50 * time.Millisecond)

	// Place a non-scout for player 1
	pos := NewPosition(1, 9)
	marshal := NewPiece(models.Marshal, player1)
	session.GetBoard().SetPieceAt(pos, marshal)

	move := NewMove(pos, NewPosition(1, 7), player1)
	err = session.SubmitMove(0, move)
	assert.Error(t, err)
	assert.Equal(t, "illegal move for this piece", err.Error())
}

func TestSessionAnimationSignaling(t *testing.T) {
	session, _, _ := setupTestSession("anim-test")

	assert.False(t, session.IsWaitingForAnimation())

	go func() {
		time.Sleep(50 * time.Millisecond)
		session.SignalAnimationComplete()
	}()

	session.WaitForAnimationComplete(200 * time.Millisecond)
	assert.False(t, session.IsWaitingForAnimation())
}

func TestSessionAnimationTimeout(t *testing.T) {
	session, _, _ := setupTestSession("anim-timeout-test")

	start := time.Now()
	session.WaitForAnimationComplete(50 * time.Millisecond)
	elapsed := time.Since(start)

	assert.GreaterOrEqual(t, elapsed, 50*time.Millisecond)
}

func TestSessionGetBoard(t *testing.T) {
	session, _, _ := setupTestSession("board-test")

	board := session.GetBoard()
	assert.NotNil(t, board)
}

func TestSessionGetSetupPieces(t *testing.T) {
	session, _, _ := setupTestSession("pieces-test")

	pieces1 := session.GetSetupPieces(0)
	pieces2 := session.GetSetupPieces(1)

	assert.Len(t, pieces1, 40)
	assert.Len(t, pieces2, 40)
}

func TestSessionStop(t *testing.T) {
	setupSession := func() *Session {
		session, _, _ := setupTestSession("stop-test")
		_ = session.StartGameFromSetup(false)
		return session
	}

	t.Run("stop-no-args", func(t *testing.T) {
		session := setupSession()
		assert.False(t, session.IsAborted())

		session.Stop()
		assert.True(t, session.IsAborted())
	})

	t.Run("stop-with-reason", func(t *testing.T) {
		session := setupSession()
		session.Stop("custom abort reason")
		assert.True(t, session.IsAborted())
	})

	t.Run("stop-with-reason-and-operator", func(t *testing.T) {
		session := setupSession()
		session.Stop("custom abort reason", "OperatorName")
		assert.True(t, session.IsAborted())
	})

	t.Run("stop-with-all-args", func(t *testing.T) {
		session := setupSession()
		session.Stop("custom abort reason", "OperatorName", 999)
		assert.True(t, session.IsAborted())
	})

	t.Run("stop-with-invalid-types-does-not-panic", func(t *testing.T) {
		session := setupSession()
		assert.NotPanics(t, func() {
			session.Stop(12345, true, "not-an-id", 3.14)
		})
		assert.True(t, session.IsAborted())
	})

	t.Run("stop-multiple-calls-does-not-panic", func(t *testing.T) {
		session := setupSession()
		session.Stop("first call")
		assert.True(t, session.IsAborted())

		assert.NotPanics(t, func() {
			session.Stop("second call")
		})
	})
}
