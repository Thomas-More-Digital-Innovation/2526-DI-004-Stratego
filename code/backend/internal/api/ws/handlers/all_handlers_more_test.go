package handlers

import (
	"digital-innovation/gostrategy/internal/api/dto"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleMove_UnmarshalError(t *testing.T) {
	_, _, client := setupTest(t)
	err := handleMove(client, []byte("invalid-json"))
	assert.Error(t, err)
}

func TestHandleMove_SubmitMoveError(t *testing.T) {
	_, session, client := setupTest(t)
	_ = session.StartGameFromSetup(false)

	// Try to move an empty space (which should trigger a game rule error)
	moveMsg := dto.MoveMessage{
		From: dto.PositionDTO{X: 0, Y: 5},
		To:   dto.PositionDTO{X: 0, Y: 4},
	}
	data, _ := json.Marshal(moveMsg)
	err := handleMove(client, data)
	assert.NoError(t, err) // Should return nil, but send a MoveResult with Success: false
}

func TestHandleUnpause_Unauthorized(t *testing.T) {
	_, _, client := setupTest(t)
	client.SeatIndex = -1
	client.UserID = 999
	err := handleUnpause(client, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")
}

func TestHandleSetSpeed_Errors(t *testing.T) {
	_, _, client := setupTest(t)

	// Test unauthorized
	client.SeatIndex = -1
	client.UserID = 999
	err := handleSetSpeed(client, nil)
	assert.Error(t, err)

	// Test unmarshal error
	client.SeatIndex = 0
	err = handleSetSpeed(client, []byte("invalid-json"))
	assert.Error(t, err)

	// Test valid speeds (below 500 and above 5000)
	msgBelow := dto.SetSpeedMessage{SpeedMs: 100}
	dataBelow, _ := json.Marshal(msgBelow)
	err = handleSetSpeed(client, dataBelow)
	assert.NoError(t, err)

	msgAbove := dto.SetSpeedMessage{SpeedMs: 9999}
	dataAbove, _ := json.Marshal(msgAbove)
	err = handleSetSpeed(client, dataAbove)
	assert.NoError(t, err)
}

func TestHandleSwapPieces_Errors(t *testing.T) {
	_, _, client := setupTest(t)

	// Test unmarshal error
	err := handleSwapPieces(client, []byte("invalid-json"))
	assert.Error(t, err)

	// Test spectator swap error
	client.SeatIndex = -1
	client.UserID = 999
	msg := dto.SwapPiecesMessage{
		Pos1: dto.PositionDTO{X: 0, Y: 6},
		Pos2: dto.PositionDTO{X: 1, Y: 6},
	}
	data, _ := json.Marshal(msg)
	err = handleSwapPieces(client, data)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "spectators cannot swap pieces")

	// Test player out of setup bounds error (Player 0 setup is Y 6-9)
	client.SeatIndex = 0
	msgBadBounds := dto.SwapPiecesMessage{
		Pos1: dto.PositionDTO{X: 0, Y: 4}, // out of setup bounds
		Pos2: dto.PositionDTO{X: 1, Y: 6},
	}
	dataBadBounds, _ := json.Marshal(msgBadBounds)
	err = handleSwapPieces(client, dataBadBounds)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "positions outside setup area")

	// Test player 1 out of setup bounds error (Player 1 setup is Y 0-3)
	client.SeatIndex = 1
	msgBadBounds1 := dto.SwapPiecesMessage{
		Pos1: dto.PositionDTO{X: 0, Y: 6}, // out of setup bounds
		Pos2: dto.PositionDTO{X: 1, Y: 2},
	}
	dataBadBounds1, _ := json.Marshal(msgBadBounds1)
	err = handleSwapPieces(client, dataBadBounds1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "positions outside setup area")
}

func TestHandleRandomizeSetup_Errors(t *testing.T) {
	_, _, client := setupTest(t)

	// Test unauthorized
	client.SeatIndex = -1
	client.UserID = 999
	err := handleRandomizeSetup(client, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")

	// Test target player < 0
	client.SeatIndex = 0
	playerID := -1
	msg := dto.RandomizeSetupMessage{PlayerID: &playerID}
	data, _ := json.Marshal(msg)
	err = handleRandomizeSetup(client, data)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "target player required")
}

func TestHandleLoadSetup_Errors(t *testing.T) {
	_, _, client := setupTest(t)

	// Test unauthorized
	client.SeatIndex = -1
	client.UserID = 999
	err := handleLoadSetup(client, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")

	// Test unmarshal error
	client.SeatIndex = 0
	err = handleLoadSetup(client, []byte("invalid-json"))
	assert.Error(t, err)

	// Test target player < 0
	playerID := -1
	msg := dto.LoadSetupMessage{SetupData: "0BBBBBB12222222233333444455556666777889M", PlayerID: &playerID}
	data, _ := json.Marshal(msg)
	err = handleLoadSetup(client, data)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "target player required")
}
