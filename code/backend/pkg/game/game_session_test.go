package game

import (
	"digital-innovation/gostrategy/pkg/game/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewSession(t *testing.T) {
	session, _, _ := setupTestSession("test-session")

	assert.NotNil(t, session)
	assert.Equal(t, "test-session", session.ID)
	assert.True(t, session.IsSetupPhase())
	assert.False(t, session.IsRunning())
}

func TestSessionSwapSetupPieces(t *testing.T) {
	session, _, _ := setupTestSession("swap-test")

	pos1 := NewPosition(0, 6)
	pos2 := NewPosition(1, 6)

	err := session.SwapSetupPieces(0, pos1, pos2)
	assert.NoError(t, err)
}

func TestSessionSwapSetupPiecesNotInSetup(t *testing.T) {
	session, _, _ := setupTestSession("swap-test-2")
	session.SetSetupPhaseComplete()

	pos1 := NewPosition(0, 6)
	pos2 := NewPosition(1, 6)

	err := session.SwapSetupPieces(0, pos1, pos2)
	assert.Error(t, err)
}

func TestSessionSwapSetupPiecesInvalidPlayer(t *testing.T) {
	session, _, _ := setupTestSession("swap-test-3")

	pos1 := NewPosition(0, 6)
	pos2 := NewPosition(1, 6)

	err := session.SwapSetupPieces(99, pos1, pos2)
	assert.Error(t, err)
}

func TestSessionRandomizeSetup(t *testing.T) {
	session, _, _ := setupTestSession("randomize-test")

	err := session.RandomizeSetup(0)
	assert.NoError(t, err)

	pieces := session.GetSetupPieces(0)
	assert.Len(t, pieces, 40)
}

func TestSessionRandomizeSetupNotInSetupPhase(t *testing.T) {
	session, _, _ := setupTestSession("randomize-test-2")
	session.SetSetupPhaseComplete()

	err := session.RandomizeSetup(0)
	assert.Error(t, err)
}

func TestSessionStartGameFromSetup(t *testing.T) {
	session, _, _ := setupTestSession("start-test")

	err := session.StartGameFromSetup(false)
	assert.NoError(t, err)
	assert.False(t, session.IsSetupPhase())
	assert.True(t, session.IsRunning())

	session.Stop()
	time.Sleep(50 * time.Millisecond)
}

func TestSessionGetGameState(t *testing.T) {
	session, _, _ := setupTestSession("state-test")

	state := session.GetGameState()

	assert.Equal(t, 1, state.Round)
	assert.Equal(t, 0, state.CurrentPlayerID)
	assert.True(t, state.IsSetupPhase)
	assert.False(t, state.IsGameOver)
}

func TestSessionSubmitMove(t *testing.T) {
	session, player1, _ := setupTestSession("move-test")

	err := session.StartGameFromSetup(false)
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond) // Wait for game to be ready

	scout := NewPiece(models.Scout, player1)
	session.GetBoard().SetPieceAt(NewPosition(0, 6), scout)

	move := NewMove(NewPosition(0, 6), NewPosition(0, 5), player1)
	err = session.SubmitMove(0, move)
	assert.NoError(t, err)

	session.Stop()
	time.Sleep(50 * time.Millisecond)
}

func TestSessionSubmitMoveNotRunning(t *testing.T) {
	session, player1, _ := setupTestSession("move-test-2")

	move := NewMove(NewPosition(0, 6), NewPosition(0, 5), player1)
	err := session.SubmitMove(0, move)
	assert.Error(t, err)
}
