// Package ws provides WebSocket functionality for the game server.
package ws

import (
	"digital-innovation/gostrategy/api/dto"
	"encoding/json"
	"fmt"
	"time"
)

// MessageHandler is a function that handles a specific WebSocket message type
type MessageHandler func(c *Client, data json.RawMessage) error

var handlers = make(map[string]MessageHandler)

// RegisterHandler registers a new message handler for a given type
func RegisterHandler(msgType string, handler MessageHandler) {
	handlers[msgType] = handler
}

// handleMessage processes incoming WebSocket messages using the registered handlers
func (c *Client) handleMessage(message []byte) {
	var baseMsg dto.WSMessage
	if err := json.Unmarshal(message, &baseMsg); err != nil {
		c.sendError("Invalid message format")
		return
	}

	handler, exists := handlers[baseMsg.Type]
	if !exists {
		c.sendError(fmt.Sprintf("Unknown message type: %s", baseMsg.Type))
		return
	}

	// We need to re-marshal the data to json.RawMessage because WSMessage.Data is interface{}
	// In a more optimized version, we would use a struct with json.RawMessage for Data
	var rawData json.RawMessage
	if baseMsg.Data != nil {
		dataBytes, _ := json.Marshal(baseMsg.Data)
		rawData = json.RawMessage(dataBytes)
	}

	if err := handler(c, rawData); err != nil {
		c.sendError(err.Error())
	}
}

// sendError sends an error message to the client
func (c *Client) sendError(err string) {
	msg := dto.WSMessage{
		Type: dto.MsgTypeError,
		Data: dto.ErrorMessage{Error: err},
	}
	if jsonData, err := json.Marshal(msg); err == nil {
		c.Send(jsonData, 100*time.Millisecond)
	}
}
