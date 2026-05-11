package api

import (
	"digital-innovation/gostrategy/api/core"
	"digital-innovation/gostrategy/logging"
)

// GameServer wraps the core GameServer to add routing and server management
type GameServer struct {
	*core.GameServer
}

// NewGameServer creates and initializes a new GameServer instance
func NewGameServer() *GameServer {
	return &GameServer{
		GameServer: core.NewGameServer(),
	}
}

// StartServer starts the HTTP server
func (s *GameServer) StartServer(addr string) error {
	s.SetupRoutes()
	s.PrintRoutes()

	logging.Debug(logging.TagWeb, "Starting game server on %s", addr)
	return s.Router.Run(addr)
}
