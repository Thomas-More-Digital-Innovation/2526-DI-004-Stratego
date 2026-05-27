// Package handlers implements WebSocket message handlers for the API.
package handlers

import (
	"digital-innovation/gostrategy/internal/api/dto"
	"digital-innovation/gostrategy/internal/api/ws"
	"encoding/json"
	"fmt"
)

func init() {
	ws.RegisterHandler(dto.MsgTypePing, handlePing)
	ws.RegisterHandler(dto.MsgTypeAnimationComplete, handleAnimationComplete)
	ws.RegisterHandler(dto.MsgTypeStep, handleStep)
}

func handlePing(c *ws.Client, _ json.RawMessage) error {
	c.SendJSON(dto.MsgTypePong, nil)
	return nil
}

func handleAnimationComplete(c *ws.Client, _ json.RawMessage) error {
	c.GetSession().SignalAnimationComplete()
	return nil
}

func handleStep(c *ws.Client, _ json.RawMessage) error {
	if !c.IsAuthorized() {
		return fmt.Errorf("unauthorized")
	}
	if c.GetSession().StepAI() {
		c.GetHub().BroadcastGameState()
	}
	return nil
}
