// Package core provides shared utilities for API.
package core

import (
	"context"
	AIhandler "digital-innovation/gostrategy/ai/handler"
	"digital-innovation/gostrategy/api/ws"
	"digital-innovation/gostrategy/db"
	"digital-innovation/gostrategy/engine"
	"digital-innovation/gostrategy/game"
	"digital-innovation/gostrategy/logging"
	"digital-innovation/gostrategy/models"
	"digital-innovation/gostrategy/utils"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

const maxGames = 500

// NewGameServer creates and initializes a new GameServer instance
func NewGameServer() *GameServer {
	if utils.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &GameServer{
		Sessions: make(map[string]*SessionHandler),
		Router:   gin.New(),
		Ctx:      ctx,
		Cancel:   cancel,
	}
}

// Stop stops the game server and releases resources
func (s *GameServer) Stop() {
	if s.Cancel != nil {
		s.Cancel()
	}

	s.Mutex.Lock()
	gameIDs := make([]string, 0, len(s.Sessions))
	for gameID := range s.Sessions {
		gameIDs = append(gameIDs, gameID)
	}
	s.Mutex.Unlock()

	for _, gameID := range gameIDs {
		s.RemoveSession(gameID)
	}
}

// CreateGame creates a new game session
func (s *GameServer) CreateGame(gameID string, gameType string, name1, name2 string) (*SessionHandler, error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	if len(s.Sessions) >= maxGames {
		return nil, fmt.Errorf("server busy: maximum number of concurrent games reached")
	}

	if _, exists := s.Sessions[gameID]; exists {
		return nil, fmt.Errorf("game %s already exists", gameID)
	}

	var controller1, controller2 engine.PlayerController
	switch gameType {
	case models.HumanVsAi:
		player1 := engine.NewPlayer(0, name1, "red")
		player2 := engine.NewPlayer(1, name2, "blue")
		controller1 = engine.NewHumanPlayerController(&player1)
		controller2 = AIhandler.CreateAI(name2, &player2)

	case models.AiVsAi:
		aiType1 := name1
		aiType2 := name2
		if name1 == name2 {
			name1 += " 1"
			name2 += " 2"
		}
		player1 := engine.NewPlayer(0, name1, "red")
		player2 := engine.NewPlayer(1, name2, "blue")
		controller1 = AIhandler.CreateAI(aiType1, &player1)
		controller2 = AIhandler.CreateAI(aiType2, &player2)

	case models.HumanVsHuman:
		player1 := engine.NewPlayer(0, name1, "red")
		player2 := engine.NewPlayer(1, name2, "blue")
		controller1 = engine.NewHumanPlayerController(&player1)
		controller2 = engine.NewHumanPlayerController(&player2)

	default:
		return nil, fmt.Errorf("unknown game type: %s", gameType)
	}

	session := game.NewSession(gameID, controller1, controller2)

	hub := ws.NewHub(session, gameType)
	hub.OnCleanup = func() {
		s.RemoveSession(gameID)
	}

	handler := &SessionHandler{
		Session:  session,
		Hub:      hub,
		GameType: gameType,
	}

	s.Sessions[gameID] = handler

	go hub.Run()
	go s.monitorGame(handler)

	return handler, nil
}

// GetSession returns a game session handler
func (s *GameServer) GetSession(gameID string) (*SessionHandler, bool) {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()
	handler, exists := s.Sessions[gameID]
	return handler, exists
}

// IsUserInActiveGame checks if a user is already in a game session (active or stale)
func (s *GameServer) IsUserInActiveGame(userID int) (*SessionHandler, bool) {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	for _, handler := range s.Sessions {
		p1, p2 := handler.Session.GetPlayerIDs()
		if (p1 != nil && *p1 == userID) || (p2 != nil && *p2 == userID) {
			return handler, true
		}
	}

	return nil, false
}

// GetUserActiveGameSeat retrieves the session handler and the player seat index (0 or 1) if the user is in an active game session.
func (s *GameServer) GetUserActiveGameSeat(userID int) (*SessionHandler, int, bool) {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	for _, handler := range s.Sessions {
		p1, p2 := handler.Session.GetPlayerIDs()
		if p1 != nil && *p1 == userID {
			return handler, 0, true
		}
		if p2 != nil && *p2 == userID {
			return handler, 1, true
		}
	}

	return nil, -1, false
}

// IsWaitingForCleanup checks if the user is waiting for cleanup in this session (game over or user disconnected)
func (sh *SessionHandler) IsWaitingForCleanup(userID int) bool {
	return sh.Session.GetGameState().IsGameOver || !sh.Hub.IsUserConnected(userID)
}

// RemoveSession removes a game session from the server and stops its resources
func (s *GameServer) RemoveSession(gameID string) {
	s.Mutex.Lock()
	handler, exists := s.Sessions[gameID]
	if exists {
		delete(s.Sessions, gameID)
	}
	s.Mutex.Unlock()

	if exists && handler != nil {
		logging.Debug(logging.TagWeb, "Removed session %s from GameServer and stopping resources", gameID)
		if handler.Hub != nil {
			handler.Hub.Stop()
		}
		if handler.Session != nil {
			handler.Session.Stop()
		}
	}
}

// monitorGame watches for game events and broadcasts them
func (s *GameServer) monitorGame(handler *SessionHandler) {
	session := handler.Session
	hub := handler.Hub

	time.Sleep(100 * time.Millisecond)
	hub.BroadcastFullState()

	select {
	case <-session.GetSetupCompleteChan():
	case <-session.IsAbortedChan():
		logging.Debug(logging.TagWeb, "Game aborted during setup: %s", session.ID)
		s.RemoveSession(session.ID)
		return
	case <-s.Ctx.Done():
		return
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.Ctx.Done():
			return
		case <-session.IsAbortedChan():
			logging.Debug(logging.TagWeb, "Game aborted during gameplay: %s", session.ID)
			s.RemoveSession(session.ID)
			return
		case <-session.GetMoveNotifyChan():
			isHeadless := session.IsHeadless()
			if !isHeadless {
				logging.Debug(logging.TagGame, "Move executed in game %s", session.ID)
			}

			if isHeadless {
				state := session.GetGameState()
				if state.IsGameOver {
					time.Sleep(100 * time.Millisecond)
					s.handleGameOver(session, hub)
					return
				}
				session.AckMoveProcessed()
				continue
			}

			combat := session.GetLastCombat()
			if combat != nil && combat.Occurred {
				hub.BroadcastCombat(combat)
				session.WaitForAnimationComplete(3 * time.Second)
				session.ClearLastCombat()
				hub.BroadcastFullState()
			} else {
				hub.BroadcastFullState()
			}

			session.AckMoveProcessed()

			state := session.GetGameState()
			if state.IsGameOver {
				time.Sleep(500 * time.Millisecond)
				s.handleGameOver(session, hub)
				return
			}
		case <-ticker.C:
			if !session.IsRunning() && session.GetGameState().IsGameOver {
				s.handleGameOver(session, hub)
				return
			}
		}
	}
}

// handleGameOver handles the end of a game
func (s *GameServer) handleGameOver(session *game.Session, hub *ws.Hub) {
	state := session.GetGameState()
	logging.Debug(logging.TagGame, "Game %s over. Winner: %v, Cause: %s", session.ID, state.WinnerID, session.GetWinCause())

	// Save game to database
	g := session.GetGame()

	// Assuming a context for DB operations, we use a background context with timeout
	dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.SaveGame(dbCtx, session.ID, session.Player1UserID, session.Player2UserID, hub.GetGameType(), g.InitialState, state.WinnerID); err != nil {
		logging.Error("Failed to save game to database", err)
	}

	if err := db.SaveMoves(dbCtx, session.ID, g.HistoricalHistory); err != nil {
		logging.Error("Failed to save game moves to database", err)
	}

	// Broadcast final state, full board, and move history
	hub.BroadcastFullState()
	hub.BroadcastMoveHistory()

	// Start cleanup timer for the hub (e.g., 5 minutes to let players see the result)
	hub.StartGameOverCleanup(5 * time.Minute)
}

// PrintRoutes prints an overview of all registered routes
func (s *GameServer) PrintRoutes() {
	fmt.Println("\n=== Registered Endpoints ===")
	for _, route := range s.Router.Routes() {
		fmt.Printf("  %-8s %-30s %s\n", route.Method, route.Path, route.Handler)
	}
	fmt.Printf("  %-8s %-30s %s\n", "WS", "/game/:gameID?player={0|1|spec}", "WebSocket connection")
	fmt.Println("============================")
	fmt.Println()
}
