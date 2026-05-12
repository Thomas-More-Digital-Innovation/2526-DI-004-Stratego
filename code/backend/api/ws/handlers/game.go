package handlers

import (
	"digital-innovation/gostrategy/api/dto"
	"digital-innovation/gostrategy/api/ws"
	"digital-innovation/gostrategy/engine"
	"digital-innovation/gostrategy/logging"
	"encoding/json"
	"fmt"
	"time"
)

func init() {
	ws.RegisterHandler(dto.MsgTypeMove, handleMove)
	ws.RegisterHandler(dto.MsgTypeGetValidMoves, handleGetValidMoves)
	ws.RegisterHandler(dto.MsgTypePause, handlePause)
	ws.RegisterHandler(dto.MsgTypeUnpause, handleUnpause)
	ws.RegisterHandler(dto.MsgTypeSetSpeed, handleSetSpeed)
}

func handleMove(c *ws.Client, data json.RawMessage) error {
	if c.GetSeatIndex() < 0 {
		return fmt.Errorf("spectators cannot make moves")
	}

	var moveMsg dto.MoveMessage
	if err := json.Unmarshal(data, &moveMsg); err != nil {
		return err
	}

	from := engine.NewPosition(moveMsg.From.X, moveMsg.From.Y)
	to := engine.NewPosition(moveMsg.To.X, moveMsg.To.Y)

	g := c.GetSession().GetGame()
	player := g.Players[c.GetSeatIndex()]
	move := engine.NewMove(from, to, player)

	err := c.GetSession().SubmitMove(c.GetSeatIndex(), move)
	if err != nil {
		sendMoveResult(c, false, err.Error())
		return nil
	}

	sendMoveResult(c, true, "")
	return nil
}

func handleGetValidMoves(c *ws.Client, data json.RawMessage) error {
	if c.GetSeatIndex() < 0 {
		return fmt.Errorf("spectators cannot request valid moves")
	}

	var reqMsg dto.GetValidMovesMessage
	if err := json.Unmarshal(data, &reqMsg); err != nil {
		return err
	}

	pos := engine.NewPosition(reqMsg.Position.X, reqMsg.Position.Y)
	moves, err := c.GetSession().GetAvailableMoves(c.GetSeatIndex(), pos)
	if err != nil {
		return err
	}

	validMoveDTOs := make([]dto.PositionDTO, len(moves))
	for i, move := range moves {
		validMoveDTOs[i] = dto.PositionToDTO(move.GetTo())
	}

	c.SendJSON(dto.MsgTypeValidMoves, dto.ValidMovesMessage{
		Position:   reqMsg.Position,
		ValidMoves: validMoveDTOs,
	})
	return nil
}

func handlePause(c *ws.Client, _ json.RawMessage) error {
	if !c.IsAuthorized() {
		return fmt.Errorf("unauthorized")
	}
	c.GetSession().Pause()
	logging.DebugWithUser(logging.TagWeb, c.GetUsername(), c.GetUserID(), "Game paused")
	c.GetHub().BroadcastGameState()
	return nil
}

func handleUnpause(c *ws.Client, _ json.RawMessage) error {
	if !c.IsAuthorized() {
		return fmt.Errorf("unauthorized")
	}
	c.GetSession().Unpause()
	logging.DebugWithUser(logging.TagWeb, c.GetUsername(), c.GetUserID(), "Game unpaused")
	c.GetHub().BroadcastGameState()
	return nil
}

func handleSetSpeed(c *ws.Client, data json.RawMessage) error {
	if !c.IsAuthorized() {
		return fmt.Errorf("unauthorized")
	}

	var msg dto.SetSpeedMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}

	speed := msg.SpeedMs
	if speed < 500 {
		speed = 500
	} else if speed > 5000 {
		speed = 5000
	}

	c.GetSession().SetTurnDelay(time.Duration(speed) * time.Millisecond)
	return nil
}

func sendMoveResult(c *ws.Client, success bool, err string) {
	c.SendJSON(dto.MsgTypeMoveResult, dto.MoveResultMessage{
		Success: success,
		Error:   err,
	})
}
