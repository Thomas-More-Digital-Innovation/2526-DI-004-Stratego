package api

import (
	"context"
	"digital-innovation/stratego/auth"
	"digital-innovation/stratego/db"
	"digital-innovation/stratego/game"
	"digital-innovation/stratego/logging"
	"digital-innovation/stratego/models"
	"digital-innovation/stratego/utils"
	"fmt"
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
	username := "Guest"
	if user != nil {
		userID = user.ID
		username = user.Username
	}

	if req.GameID == "" {
		req.GameID = fmt.Sprintf("game-%d-%d", time.Now().Unix(), time.Now().UnixNano()%1000000)
	}

	if req.GameType == "" {
		req.GameType = models.HumanVsAi
	}

	name1 := username
	name2 := req.AI1
	switch req.GameType {
	case models.AiVsAi:
		name1 = req.AI1
		name2 = req.AI2
	case models.HumanVsHuman:
		name1 = "Human Red"
		name2 = "Human Blue"
	}

	handler, err := s.CreateGame(req.GameID, req.GameType, name1, name2)
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

	username = "Guest"
	if user != nil {
		username = user.Username
		if req.GameType != models.AiVsAi {
			handler.Session.Player1Username = username
		}
	}
	logging.GameStarted(req.GameID, req.GameType, username, userID)
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

	p1ID, p2ID := handler.Session.GetPlayerIDs()
	switch playerID {
	case 0:
		if p1ID != nil {
			if currentUserID == nil || *currentUserID != *p1ID {
				sendError(c, "Unauthorized: You are not Player 1", http.StatusForbidden)
				return
			}
		} else if currentUserID != nil {
			handler.Session.SetPlayer1Associate(*currentUserID, user.Username)
		}
	case 1:
		if p2ID != nil {
			if currentUserID == nil || *currentUserID != *p2ID {
				sendError(c, "Unauthorized: You are not Player 2", http.StatusForbidden)
				return
			}
		} else if currentUserID != nil {
			handler.Session.SetPlayer2Associate(*currentUserID, user.Username)
		}
	}

	username, userID := utils.TryGetUser(user)
	logging.GameStarted(gameID, "WebSocket Join", username, userID)

	HandleWebSocket(c.Writer, c.Request, handler.Session, handler.Hub, playerID, username, userID)
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
	handlers := make([]*SessionHandler, 0, len(s.sessions))
	for _, handler := range s.sessions {
		handlers = append(handlers, handler)
	}
	s.mutex.RUnlock()

	// Build summaries outside of the global lock
	summaries := make([]models.GameSummary, 0, len(handlers))
	for _, handler := range handlers {
		summaries = append(summaries, handler.Session.GetGameSummary(handler.GameType))
	}

	sendJSON(c, summaries, http.StatusOK)
}

// handleGameOver broadcasts final game state and saves stats
func (s *GameServer) handleGameOver(session *game.Session, hub *WSHub) {
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

	// Determine labels for logging
	p1 := logging.FormatUser(session.Player1Username, utils.GetIntSafe(session.Player1UserID))
	p2 := logging.FormatUser(session.Player2Username, utils.GetIntSafe(session.Player2UserID))

	winnerLabel := "Draw"
	loserLabel := "Draw"

	if winner != nil {
		if winner.GetID() == 0 {
			winnerLabel = p1
			loserLabel = p2
		} else {
			winnerLabel = p2
			loserLabel = p1
		}
	}

	logging.GameFinished(session.ID, winnerLabel, loserLabel, state.Round)

	// Start mandatory cleanup (10 minutes) to prevent zombie sessions
	// If players are still watching the board, they have 10 minutes before the session is killed
	hub.StartGameOverCleanup(10 * time.Minute)
}

// saveGameStats saves game statistics to the database
func (s *GameServer) saveGameStats(session *game.Session, winnerID *int, gameType string) {
	duration := time.Since(session.StartTime).Seconds()
	state := session.GetGameState()

	// Save game history (metadata and moves)
	g := session.GetGame()
	initialState := g.GetInitialBoardState()

	// Use a background context with timeout for stats saving to avoid leaking goroutines if DB hangs
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
	defer cancel()

	if err := db.SaveGame(ctx, session.ID, session.Player1UserID, session.Player2UserID, gameType, initialState, winnerID); err != nil {
		logging.ErrorWith2Users(fmt.Sprintf("Failed to save game metadata for %s", session.ID), session.Player1Username, utils.GetIntSafe(session.Player1UserID), session.Player2Username, utils.GetIntSafe(session.Player2UserID), err)
	} else {
		if err := db.SaveMoves(ctx, session.ID, g.HistoricalHistory); err != nil {
			logging.ErrorWith2Users(fmt.Sprintf("Failed to save moves for game %s", session.ID), session.Player1Username, utils.GetIntSafe(session.Player1UserID), session.Player2Username, utils.GetIntSafe(session.Player2UserID), err)
		}
		logging.DebugWithUser(logging.TagGame, session.Player1Username, utils.GetIntSafe(session.Player1UserID), "Saved full history for game %s (%d moves)", session.ID, len(g.HistoricalHistory))
	}

	// Track stats for player 1 if they have a user ID
	if session.Player1UserID != nil {
		userID := *session.Player1UserID
		won := winnerID != nil && *winnerID == 0

		if err := db.UpdateUserStats(ctx, userID, won, state.MoveCount, duration); err != nil {
			logging.ErrorWithUser("Failed to update stats", session.Player1Username, userID, err)
		} else {
			logging.DebugWithUser(logging.TagGame, session.Player1Username, userID, "Updated stats (won=%v)", won)
		}

	}

	// Track stats for player 2 if they have a user ID
	if session.Player2UserID != nil {
		userID := *session.Player2UserID
		won := winnerID != nil && *winnerID == 1

		if err := db.UpdateUserStats(ctx, userID, won, state.MoveCount, duration); err != nil {
			logging.ErrorWithUser("Failed to update stats", session.Player2Username, userID, err)
		} else {
			logging.DebugWithUser(logging.TagGame, session.Player2Username, userID, "Updated stats (won=%v)", won)
		}

	}
}
