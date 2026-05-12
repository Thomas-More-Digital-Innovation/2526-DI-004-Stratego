// Package ws provides WebSocket functionality for the game server.
package ws

import (
	"digital-innovation/gostrategy/api/dto"
	"digital-innovation/gostrategy/game"
	"digital-innovation/gostrategy/logging"
	"digital-innovation/gostrategy/models"
	"encoding/json"
	"time"

	"sync"

	"github.com/gorilla/websocket"
	"golang.org/x/time/rate"
)

// Client represents a WebSocket client connection
type Client struct {
	conn      *websocket.Conn
	send      chan []byte
	Session   *game.Session
	SeatIndex int // -1 for spectator, 0 or 1 for player
	Hub       *Hub
	Username  string
	UserID    int
	limiter   *rate.Limiter
	mu        sync.Mutex
	closed    bool
}

// GetSeatIndex returns the client's seat index
func (c *Client) GetSeatIndex() int { return c.SeatIndex }

// GetSession returns the client's game session
func (c *Client) GetSession() *game.Session { return c.Session }

// GetHub returns the client's hub
func (c *Client) GetHub() *Hub { return c.Hub }

// GetUsername returns the client's username
func (c *Client) GetUsername() string { return c.Username }

// GetUserID returns the client's user ID
func (c *Client) GetUserID() int { return c.UserID }

// SendJSON helper to send a JSON message
func (c *Client) SendJSON(msgType string, data any) {
	msg := dto.WSMessage{
		Type: msgType,
		Data: data,
	}
	if jsonData, err := json.Marshal(msg); err == nil {
		c.Send(jsonData, 100*time.Millisecond)
	}
}

// IsAuthorized checks if the client is a player or the creator of an AI game
func (c *Client) IsAuthorized() bool {
	p1ID, _ := c.Session.GetPlayerIDs()
	isCreator := p1ID == nil || c.UserID == *p1ID
	return c.SeatIndex >= 0 || (c.Hub.GameType == models.HumanVsAi && isCreator) || (c.Hub.GameType == models.AiVsAi && isCreator)
}

// Close marks the client as closed and stops sending messages
func (c *Client) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.closed {
		c.closed = true
		if c.send != nil {
			close(c.send)
		}
	}
}

// IsClosed returns whether the client connection is closed
func (c *Client) IsClosed() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.closed
}

// Send attempts to send data to the client with a timeout
func (c *Client) Send(data []byte, timeout time.Duration) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return false
	}

	select {
	case c.send <- data:
		return true
	case <-time.After(timeout):
		return false
	}
}

// readPump pumps messages from the websocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		select {
		case c.Hub.unregister <- c:
		default:
			// Hub is already stopped or busy; avoid blocking here to prevent goroutine leak
		}
		_ = c.conn.Close()
	}()

	err := c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	if err != nil {
		logging.ConnectionError(c.Session.ID, c.Username, c.UserID, err)
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
				logging.ConnectionError(c.Session.ID, c.Username, c.UserID, err)
			} else {
				logging.ConnectionClosed(c.Session.ID, c.Username, c.UserID)
			}
			break
		}

		// Rate limiting: 20 messages per second with a burst of 40
		if c.limiter == nil {
			c.limiter = rate.NewLimiter(rate.Limit(20), 40)
		}

		if !c.limiter.Allow() {
			logging.SecurityWarning("WebSocket rate limit exceeded", "Dropping message from user", c.Username, c.UserID)
			continue
		}

		c.handleMessage(message)
	}
}

// writePump pumps messages from the hub to the websocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()

	}()

	for {
		select {
		case message, ok := <-c.send:
			err := c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err != nil {
				logging.ConnectionError(c.Session.ID, c.Username, c.UserID, err)
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
