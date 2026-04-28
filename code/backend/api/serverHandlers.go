package api

import (
	"digital-innovation/stratego/auth"
	"digital-innovation/stratego/db"
	"digital-innovation/stratego/game"
	"digital-innovation/stratego/models"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HTTP Handlers

// HandleCreateGame handles game creation
// @Summary Create a new game
// @Description Initialize a new game session with specified type and AIs
// @Tags games
// @Accept json
// @Produce json
// @Param request body map[string]string true "Game creation details (id, type, ai1, ai2)"
// @Success 201 {object} map[string]string "Game created"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Router /games [post]
func (s *GameServer) HandleCreateGame(c *gin.Context) {
	var req struct {
		GameID   string `json:"gameId"`
		GameType string `json:"gameType"`
		AI1      string `json:"ai1"`
		AI2      string `json:"ai2"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		sendError(c, "Invalid request body", http.StatusBadRequest)
		return
	}

	user := auth.GetCurrentUser(c)
	// We allow guests to create games too, but if logged in, we track them
	userID := -1
	if user != nil {
		userID = user.ID
	}

	if req.GameID == "" {
		req.GameID = fmt.Sprintf("game-%d-%d", time.Now().Unix(), time.Now().UnixNano()%1000000)
	}

	if req.GameType == "" {
		req.GameType = models.HumanVsAi
	}

	handler, err := s.CreateGame(req.GameID, req.GameType, req.AI1, req.AI2)
	if err != nil {
		sendError(c, err.Error(), http.StatusBadRequest)
		return
	}

	// Set creator as Player 1
	if userID != -1 {
		handler.Session.Player1UserID = &userID
	}

	response := gin.H{
		"gameId":   req.GameID,
		"gameType": req.GameType,
		"wsUrl":    fmt.Sprintf("/game/%s", req.GameID),
	}

	sendJSON(c, response, http.StatusOK)

	log.Printf("Created game %s (type: %s) by user %d", req.GameID, req.GameType, userID)
}

// HandleWebSocketConnection handles WebSocket connections
// @Summary Game WebSocket
// @Description Real-time game connection. Use `player` query param to join as 0 (Red), 1 (Blue), or anything else (Spectator)
// @Tags games
// @Param gameID path string true "Game ID"
// @Param player query string false "Player role (0, 1, or spec)"
// @Router /game/{gameID} [get]
func (s *GameServer) HandleWebSocketConnection(c *gin.Context) {
	gameID := c.Param("gameID")
	if gameID == "" {
		sendError(c, "Game ID required", http.StatusBadRequest)
		return
	}

	handler, exists := s.GetSession(gameID)
	if !exists {
		sendError(c, "Game not found", http.StatusNotFound)
		return
	}

	// Get player ID from query parameter
	playerIDStr := c.Query("player")
	playerID := -1 // Default to spectator

	switch playerIDStr {
	case "0":
		playerID = 0
	case "1":
		playerID = 1
	}

	// Security Check: Verify user session against authorized player IDs
	user := auth.GetCurrentUser(c)
	var currentUserID *int
	if user != nil {
		currentUserID = &user.ID
	}

	switch playerID {
	case 0:
		if handler.Session.Player1UserID != nil {
			if currentUserID == nil || *currentUserID != *handler.Session.Player1UserID {
				sendError(c, "Unauthorized: You are not Player 1", http.StatusForbidden)
				return
			}
		} else if currentUserID != nil {
			// Associate if vacant
			handler.Session.Player1UserID = currentUserID
		}
	case 1:
		if handler.Session.Player2UserID != nil {
			if currentUserID == nil || *currentUserID != *handler.Session.Player2UserID {
				sendError(c, "Unauthorized: You are not Player 2", http.StatusForbidden)
				return
			}
		} else if currentUserID != nil {
			// Associate if vacant
			handler.Session.Player2UserID = currentUserID
		}
	}

	log.Printf("WebSocket connection for game %s (player %d, user %v)", gameID, playerID, currentUserID)

	HandleWebSocket(c.Writer, c.Request, handler.Session, handler.Hub, playerID)
}

// HandleListGames handles GET /games
// @Summary List active games
// @Description Retrieve a list of all currently active game sessions
// @Tags games
// @Produce json
// @Success 200 {array} map[string]interface{} "List of games"
// @Router /games [get]
func (s *GameServer) HandleListGames(c *gin.Context) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	games := make([]gin.H, 0, len(s.sessions))
	for gameID, handler := range s.sessions {
		state := handler.Session.GetGameState()
		games = append(games, gin.H{
			"gameId":     gameID,
			"round":      state.Round,
			"isRunning":  handler.Session.IsRunning(),
			"isGameOver": state.IsGameOver,
		})
	}

	sendJSON(c, games, http.StatusOK)
}

// handleGameOver broadcasts final game state and saves stats
func (s *GameServer) handleGameOver(session *game.GameSession, hub *WSHub) {
	hub.BroadcastGameState()

	state := session.GetGameState()
	winner := session.GetWinner()
	var winnerID *int
	var winnerName string
	if winner != nil {
		id := winner.GetID()
		winnerID = &id
		winnerName = winner.GetName()
	}

	gameOverMsg := GameOverMessage{
		WinnerID:   winnerID,
		WinnerName: winnerName,
		WinCause:   string(session.GetWinCause()),
		Round:      state.Round,
	}

	hub.BroadcastMessage(MsgTypeGameOver, gameOverMsg)

	// Broadcast final board state with all pieces revealed
	s.broadcastBoardStateRevealed(hub)
	hub.BroadcastMoveHistory()

	// Save game stats to database
	go s.saveGameStats(session, winnerID, hub.gameType)

	// Wait longer before stopping monitoring so users can see results
	time.Sleep(30 * time.Second)
}

// saveGameStats saves game statistics to the database
func (s *GameServer) saveGameStats(session *game.GameSession, winnerID *int, gameType string) {
	duration := time.Since(session.StartTime).Seconds()
	state := session.GetGameState()

	// Save game history (metadata and moves)
	g := session.GetGame()
	initialState := g.GetInitialBoardState()

	if err := db.SaveGame(session.ID, session.Player1UserID, session.Player2UserID, gameType, initialState, winnerID); err != nil {
		log.Printf("Failed to save game metadata for %s: %v", session.ID, err)
	} else {
		for _, m := range g.HistoricalHistory {
			if err := db.SaveMove(session.ID, m); err != nil {
				log.Printf("Failed to save move %d for game %s: %v", m.MoveIndex, session.ID, err)
			}
		}
		log.Printf("Saved full history for game %s (%d moves)", session.ID, len(g.HistoricalHistory))
	}
	// Track stats for player 1 if they have a user ID
	if session.Player1UserID != nil {
		userID := *session.Player1UserID
		won := winnerID != nil && *winnerID == 0

		if err := db.UpdateUserStats(userID, won, state.MoveCount, duration); err != nil {
			log.Printf("Failed to update stats for user %d: %v", userID, err)
		} else {
			log.Printf("Updated stats for user %d (won=%v)", userID, won)
		}
	}

	// Track stats for player 2 if they have a user ID
	if session.Player2UserID != nil {
		userID := *session.Player2UserID
		won := winnerID != nil && *winnerID == 1

		if err := db.UpdateUserStats(userID, won, state.MoveCount, duration); err != nil {
			log.Printf("Failed to update stats for user %d: %v", userID, err)
		} else {
			log.Printf("Updated stats for user %d (won=%v)", userID, won)
		}
	}
}
