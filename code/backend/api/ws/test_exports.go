package ws

import "time"

//! Exported for testing only.

// Register exposes the register channel for testing.
func (h *Hub) Register() chan<- *Client { return h.register }

// Unregister exposes the unregister channel for testing.
func (h *Hub) Unregister() chan<- *Client { return h.unregister }

// Broadcast exposes the broadcast channel for testing.
func (h *Hub) Broadcast() chan<- []byte { return h.broadcast }

// Mutex exposes the hub mutex for testing.
func (h *Hub) Mutex() any { return &h.mutex }

// SetCleanupPeriod sets the cleanup period for testing.
func (h *Hub) SetCleanupPeriod(d time.Duration) { h.cleanupPeriod = d }

// SendChan returns the send channel for testing.
func (c *Client) SendChan() chan []byte { return c.send }

// NewTestClient creates a new test client.
func NewTestClient() *Client {
	return &Client{
		send:     make(chan []byte, 256),
		stopChan: make(chan struct{}),
	}
}
