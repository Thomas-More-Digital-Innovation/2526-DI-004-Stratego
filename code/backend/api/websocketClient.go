package api

import (
	"digital-innovation/stratego/game"
	"digital-innovation/stratego/logging"
	"time"

	"github.com/gorilla/websocket"
)

// WSClient represents a WebSocket client connection
type WSClient struct {
	conn      *websocket.Conn
	send      chan []byte
	session   *game.GameSession
	seatIndex int // -1 for spectator, 0 or 1 for player
	hub       *WSHub
	Username  string
	UserID    int
}

// readPump pumps messages from the websocket connection to the hub
func (c *WSClient) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	err := c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	if err != nil {
		logging.ConnectionError(c.session.ID, c.Username, c.UserID, err)
		return
	}
	c.conn.SetPongHandler(func(string) error {
		err := c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		if err != nil {
			logging.DebugWithUser(logging.TagWeb, c.Username, c.UserID, "Error setting read deadline: %v", err)
		}
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logging.ConnectionError(c.session.ID, c.Username, c.UserID, err)
			} else {
				logging.ConnectionClosed(c.session.ID, c.Username, c.UserID)
			}
			break
		}

		c.handleMessage(message)
	}
}

// writePump pumps messages from the hub to the websocket connection
func (c *WSClient) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			err := c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err != nil {
				logging.ConnectionError(c.session.ID, c.Username, c.UserID, err)
				return
			}
			if !ok {
				err := c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					logging.DebugWithUser(logging.TagWeb, c.Username, c.UserID, "Error writing close message: %v", err)
				}
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			err := c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err != nil {
				logging.DebugWithUser(logging.TagWeb, c.Username, c.UserID, "Error setting write deadline: %v", err)
				return
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
