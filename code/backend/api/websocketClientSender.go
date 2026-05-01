package api

import (
	"digital-innovation/stratego/logging"
	"encoding/json"
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
		logging.ErrorWithUser("Error marshaling move result", c.Username, c.UserID, err)
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
		logging.ErrorWithUser("Error marshaling error message", c.Username, c.UserID, err)
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
		logging.ErrorWithUser("Error marshaling pong", c.Username, c.UserID, err)
		return
	}

	c.send <- jsonData
}
