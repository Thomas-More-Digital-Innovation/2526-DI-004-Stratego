package handlers

import (
	"digital-innovation/gostrategy/api/dto"
	"digital-innovation/gostrategy/api/ws"
	"digital-innovation/gostrategy/engine"
	"digital-innovation/gostrategy/game"
	"digital-innovation/gostrategy/models"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTest(_ *testing.T) (*ws.Hub, *game.Session, *ws.Client) {
	player1 := engine.NewPlayer(0, "P1", "red")
	player2 := engine.NewPlayer(1, "P2", "blue")
	c1 := engine.NewHumanPlayerController(&player1)
	c2 := engine.NewHumanPlayerController(&player2)
	session := game.NewSession("test-session", c1, c2)
	hub := ws.NewHub(session, models.HumanVsHuman)
	client := ws.NewTestClient()
	client.Session = session
	client.Hub = hub
	client.SeatIndex = 0
	client.UserID = 1
	client.Username = "P1"

	_ = session.RandomizeSetup(0)
	_ = session.RandomizeSetup(1)

	return hub, session, client
}

func TestHandleMove(t *testing.T) {
	_, session, client := setupTest(t)
	_ = session.StartGameFromSetup(false)

	// Ensure there is a piece at 0, 6 (Marshal for example)
	marshal := engine.NewPiece(models.Marshal, session.GetGame().Players[0])
	session.GetGame().Board.SetPieceAt(engine.NewPosition(0, 6), marshal)

	moveMsg := dto.MoveMessage{
		From: dto.PositionDTO{X: 0, Y: 6},
		To:   dto.PositionDTO{X: 0, Y: 5},
	}
	data, _ := json.Marshal(moveMsg)

	err := handleMove(client, data)
	assert.NoError(t, err)
}

func TestHandleGetValidMoves(t *testing.T) {
	_, session, client := setupTest(t)
	_ = session.StartGameFromSetup(false)

	// Ensure there is a piece at 0, 6
	marshal := engine.NewPiece(models.Marshal, session.GetGame().Players[0])
	session.GetGame().Board.SetPieceAt(engine.NewPosition(0, 6), marshal)

	reqMsg := dto.GetValidMovesMessage{
		Position: dto.PositionDTO{X: 0, Y: 6},
	}
	data, _ := json.Marshal(reqMsg)

	err := handleGetValidMoves(client, data)
	assert.NoError(t, err)
}

func TestHandlePauseUnpause(t *testing.T) {
	_, _, client := setupTest(t)

	err := handlePause(client, nil)
	assert.NoError(t, err)

	err = handleUnpause(client, nil)
	assert.NoError(t, err)
}

func TestHandleSetSpeed(t *testing.T) {
	_, _, client := setupTest(t)

	msg := dto.SetSpeedMessage{SpeedMs: 1000}
	data, _ := json.Marshal(msg)

	err := handleSetSpeed(client, data)
	assert.NoError(t, err)
}

func TestHandleSwapPieces(t *testing.T) {
	_, _, client := setupTest(t)

	msg := dto.SwapPiecesMessage{
		Pos1: dto.PositionDTO{X: 0, Y: 6},
		Pos2: dto.PositionDTO{X: 1, Y: 6},
	}
	data, _ := json.Marshal(msg)

	err := handleSwapPieces(client, data)
	assert.NoError(t, err)
}

func TestHandleRandomizeSetup(t *testing.T) {
	_, _, client := setupTest(t)

	playerID := 0
	msg := dto.RandomizeSetupMessage{PlayerID: &playerID}
	data, _ := json.Marshal(msg)

	err := handleRandomizeSetup(client, data)
	assert.NoError(t, err)
}

func TestHandleStartGame(t *testing.T) {
	_, _, client := setupTest(t)

	msg := dto.StartGameMessage{Headless: false}
	data, _ := json.Marshal(msg)

	err := handleStartGame(client, data)
	assert.NoError(t, err)
}

func TestHandleLoadSetup(t *testing.T) {
	_, _, client := setupTest(t)

	// Need 40 bytes of valid rank data with correct counts
	setupData := "0BBBBBB12222222233333444455556666777889M"
	msg := dto.LoadSetupMessage{SetupData: setupData}
	data, _ := json.Marshal(msg)

	err := handleLoadSetup(client, data)
	assert.NoError(t, err)
}

func TestHandleMove_Spectator(t *testing.T) {
	_, _, client := setupTest(t)
	client.SeatIndex = -1 // Spectator

	err := handleMove(client, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "spectators cannot make moves")
}

func TestHandlePause_Unauthorized(t *testing.T) {
	_, _, client := setupTest(t)
	client.SeatIndex = -1
	client.UserID = 999 // Not the creator

	err := handlePause(client, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")
}

func TestHandleGetValidMoves_Errors(t *testing.T) {
	_, _, client := setupTest(t)

	// Test spectator error
	client.SeatIndex = -1
	err := handleGetValidMoves(client, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "spectators cannot request valid moves")

	// Test unmarshal error
	client.SeatIndex = 0
	err = handleGetValidMoves(client, []byte("invalid-json"))
	assert.Error(t, err)

	// Test GetAvailableMoves error (e.g. position with no piece)
	reqMsg := dto.GetValidMovesMessage{
		Position: dto.PositionDTO{X: 0, Y: 4},
	}
	data, _ := json.Marshal(reqMsg)
	err = handleGetValidMoves(client, data)
	assert.Error(t, err)
}

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
