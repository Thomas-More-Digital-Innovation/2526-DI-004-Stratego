// Package core provides shared utilities for API.
package core

import (
	"context"
	"digital-innovation/gostrategy/api/ws"
	"digital-innovation/gostrategy/game"
	"sync"

	"github.com/gin-gonic/gin"
)

// GameServer manages HTTP and WebSocket connections
type GameServer struct {
	Sessions map[string]*SessionHandler
	Mutex    sync.RWMutex
	Router   *gin.Engine
	Ctx      context.Context
	Cancel   context.CancelFunc
}

// SessionHandler wraps a game session with its WebSocket hub
type SessionHandler struct {
	Session  *game.Session
	Hub      *ws.Hub
	GameType string
}
