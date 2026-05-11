package ws

import "time"

// Exported for testing only.
func (h *Hub) Register() chan<- *Client         { return h.register }
func (h *Hub) Unregister() chan<- *Client       { return h.unregister }
func (h *Hub) Broadcast() chan<- []byte         { return h.broadcast }
func (h *Hub) Clients() map[*Client]bool        { return h.clients }
func (h *Hub) Mutex() any                       { return &h.mutex }
func (h *Hub) SetCleanupPeriod(d time.Duration) { h.cleanupPeriod = d }

func (c *Client) SendChan() chan []byte { return c.send }
