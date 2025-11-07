package api

import (
	"encoding/json"
	"log"
)

// sendMoveResult sends a move result message
func (c *WSClient) sendMoveResult(success bool, error string) {
	result := MoveResultMessage{
		Success: success,
		Error:   error,
	}

	msg := WSMessage{
		Type: MsgTypeMoveResult,
		Data: result,
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling move result: %v", err)
		return
	}

	c.send <- jsonData
}

// sendError sends an error message
func (c *WSClient) sendError(errMsg string) {
	msg := WSMessage{
		Type: MsgTypeError,
		Data: ErrorMessage{Error: errMsg},
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling error message: %v", err)
		return
	}

	c.send <- jsonData
}

// sendPong sends a pong response
func (c *WSClient) sendPong() {
	msg := WSMessage{
		Type: MsgTypePong,
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling pong: %v", err)
		return
	}

	c.send <- jsonData
}
