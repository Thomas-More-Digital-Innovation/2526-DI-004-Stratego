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
