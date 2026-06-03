package game

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSession_SetupTimeout(t *testing.T) {
	session, _, _ := setupTestSession("test-timeout", WithSetupTimeout(100*time.Millisecond))

	assert.True(t, session.IsSetupPhase())

	// Wait for timeout to expire
	time.Sleep(200 * time.Millisecond)

	assert.False(t, session.IsSetupPhase())
	assert.True(t, session.IsRunning())
}

func TestSession_Abort(t *testing.T) {
	session, _, _ := setupTestSession("test-abort")

	err := session.StartGameFromSetup(true)
	assert.NoError(t, err)
	assert.True(t, session.IsRunning())

	session.Stop()
	assert.True(t, session.IsAborted())

	// Wait a bit for the runner to stop
	time.Sleep(50 * time.Millisecond)
	assert.False(t, session.IsRunning())
}

func TestSession_SwapSetupPieces(t *testing.T) {
	session, _, _ := setupTestSession("test-swap")

	p1 := NewPosition(0, 6)
	p2 := NewPosition(1, 6)

	piecesBefore := make([]*Piece, len(session.GetPlayer1Pieces()))
	copy(piecesBefore, session.GetPlayer1Pieces())

	err := session.SwapSetupPieces(0, p1, p2)
	assert.NoError(t, err)

	assert.Equal(t, piecesBefore[1], session.GetPlayer1Pieces()[0])
	assert.Equal(t, piecesBefore[0], session.GetPlayer1Pieces()[1])

	// Test invalid player
	err = session.SwapSetupPieces(2, p1, p2)
	assert.Error(t, err)

	// Test invalid position (outside setup area)
	err = session.SwapSetupPieces(0, p1, NewPosition(0, 0))
	assert.Error(t, err)
}

func TestSession_SetupWarning(t *testing.T) {
	session, _, _ := setupTestSession("test-warning",
		WithSetupTimeout(200*time.Millisecond),
		WithSetupWarning(50*time.Millisecond),
	)

	// We expect a notification to be sent to moveNotifyChan when the warning fires
	// We'll use WaitForMoveNotification to check for it
	received := session.WaitForMoveNotification(150 * time.Millisecond)
	assert.True(t, received)
	assert.True(t, session.IsSetupPhase())
}
