package ws

import (
	"digital-innovation/gostrategy/api/dto"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDispatcher_HandleMessage(t *testing.T) {
	client := NewTestClient()

	// Track if handler was called
	called := false
	var receivedData json.RawMessage

	testHandler := func(_ *Client, data json.RawMessage) error {
		called = true
		receivedData = data
		return nil
	}

	// Register a test handler
	msgType := "test-msg"
	RegisterHandler(msgType, testHandler)

	t.Run("successful dispatch", func(t *testing.T) {
		called = false
		msg := dto.WSMessage{
			Type: msgType,
			Data: map[string]string{"foo": "bar"},
		}
		jsonData, _ := json.Marshal(msg)

		client.handleMessage(jsonData)

		assert.True(t, called)
		assert.Contains(t, string(receivedData), "foo")
		assert.Contains(t, string(receivedData), "bar")
	})

	t.Run("invalid json", func(t *testing.T) {
		client.handleMessage([]byte("{invalid"))

		// Should send error message
		msg := <-client.send
		assert.Contains(t, string(msg), "error")
		assert.Contains(t, string(msg), "Invalid message format")
	})

	t.Run("unknown message type", func(t *testing.T) {
		msg := dto.WSMessage{
			Type: "unknown",
			Data: "data",
		}
		jsonData, _ := json.Marshal(msg)

		client.handleMessage(jsonData)

		// Should send error message
		msgBytes := <-client.send
		assert.Contains(t, string(msgBytes), "error")
		assert.Contains(t, string(msgBytes), "Unknown message type")
	})

	t.Run("handler error", func(t *testing.T) {
		RegisterHandler("error-msg", func(_ *Client, _ json.RawMessage) error {
			return assert.AnError
		})

		msg := dto.WSMessage{
			Type: "error-msg",
			Data: "data",
		}
		jsonData, _ := json.Marshal(msg)

		client.handleMessage(jsonData)

		// Should send error message from handler
		msgBytes := <-client.send
		assert.Contains(t, string(msgBytes), "error")
		assert.Contains(t, string(msgBytes), assert.AnError.Error())
	})
}
