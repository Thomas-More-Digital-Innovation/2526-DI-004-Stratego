package handlers

import (
	"digital-innovation/gostrategy/internal/api/core"
	"digital-innovation/gostrategy/internal/api/ws"
	"digital-innovation/gostrategy/internal/auth"
	"digital-innovation/gostrategy/internal/logging"
	"digital-innovation/gostrategy/internal/models"
	"digital-innovation/gostrategy/internal/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	keyGameType = "gameType"
	keyAI1      = "ai1"
)

// Handler provides REST API handlers with access to the game server
type Handler struct {
	*core.GameServer
}

// NewHandler creates a new REST handler
func NewHandler(s *core.GameServer) *Handler {
	return &Handler{s}
}

// HandleCreateGame handles game creation
// @Summary Create a new game
// @Description Initialize a new game session with specified type and players
// @Tags games
// @Accept json
// @Produce json
// @Param request body object{gameId=string,gameType=string,ai1=string,ai2=string} true "Game creation parameters"
// @Success 200 {object} map[string]string "Game created successfully"
// @Failure 400 {object} map[string]string "Invalid request body or parameters"
// @Router /games [post]
func (h *Handler) HandleCreateGame(c *gin.Context) {
	var req struct {
		GameID   string `json:"gameId"`
		GameType string `json:"gameType"`
		AI1      string `json:"ai1"`
		AI2      string `json:"ai2"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		core.SendError(c, "Invalid request body", http.StatusBadRequest)
		return
	}

	user := auth.GetCurrentUser(c)
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

	// Limiter: A player can only play one game at a time
	if userID != -1 {
		if existingHandler, found := h.IsUserInActiveGame(userID); found {
			if !existingHandler.IsWaitingForCleanup(userID) {
				core.SendError(c, fmt.Sprintf("You are already in an active game (%s). Please finish it before starting a new one.", existingHandler.Session.ID), http.StatusConflict)
				return
			}
			// stale/waiting for cleanup: clean up directly and allow new game
			logging.Debug(logging.TagWeb, "User %d starting new game, force cleaning up stale session %s", userID, existingHandler.Session.ID)
			h.RemoveSession(existingHandler.Session.ID, "Replaced by new session")
		}
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

	handler, err := h.CreateGame(req.GameID, req.GameType, name1, name2)
	if err != nil {
		core.SendError(c, err.Error(), http.StatusBadRequest)
		return
	}

	if userID != -1 {
		if req.GameType != models.AiVsAi {
			handler.Session.SetPlayer1Associate(userID, username)
		} else {
			// For AI vs AI, just set the ID without changing the AI name
			handler.Session.SetPlayer1Associate(userID, name1)
		}
	}

	core.SendJSON(c, gin.H{
		"gameId":    req.GameID,
		keyGameType: req.GameType,
		"wsUrl":     fmt.Sprintf("/game/%s", req.GameID),
	}, http.StatusOK)

	logging.GameStarted(req.GameID, req.GameType, username, userID)
}

// HandleWebSocketConnection handles WebSocket connections
// @Summary Connect to game via WebSocket
// @Description Establish a WebSocket connection to an active game session
// @Tags games
// @Param gameID path string true "Game ID"
// @Param player query string false "Player index (0 or 1)"
// @Success 101 {string} string "Switching Protocols"
// @Failure 400 {object} map[string]string "Game ID required or invalid player index"
// @Failure 403 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Game not found"
// @Router /game/{gameID} [get]
func (h *Handler) HandleWebSocketConnection(c *gin.Context) {
	gameID := c.Param("gameID")
	if gameID == "" {
		core.SendError(c, "Game ID required", http.StatusBadRequest)
		return
	}

	handler, exists := h.GetSession(gameID)
	if !exists {
		core.SendError(c, "Game not found", http.StatusNotFound)
		return
	}

	playerIDStr := c.Query("player")
	playerID := -1
	switch playerIDStr {
	case "0":
		playerID = 0
	case "1":
		playerID = 1
	}

	user := auth.GetCurrentUser(c)
	var currentUserID *int
	if user != nil {
		currentUserID = &user.ID
	}

	p1ID, p2ID := handler.Session.GetPlayerIDs()
	switch playerID {
	case 0:
		if !h.associatePlayer(c, currentUserID, p1ID, user, gameID) {
			return
		}
		if p1ID == nil && user != nil && handler.GameType != models.AiVsAi {
			handler.Session.SetPlayer1Associate(user.ID, user.Username)
		}
	case 1:
		if !h.associatePlayer(c, currentUserID, p2ID, user, gameID) {
			return
		}
		if p2ID == nil && user != nil && handler.GameType != models.AiVsAi {
			handler.Session.SetPlayer2Associate(user.ID, user.Username)
		}
	}

	username, userID := utils.TryGetUser(user)
	ws.HandleWebSocket(c.Writer, c.Request, handler.Session, handler.Hub, playerID, username, userID)
}

func (h *Handler) associatePlayer(c *gin.Context, currentUserID *int, associatedUserID *int, user *models.User, gameID string) bool {
	if associatedUserID != nil {
		if currentUserID == nil || *currentUserID != *associatedUserID {
			core.SendError(c, "Unauthorized: You are not logged in as the correct player.", http.StatusForbidden)
			return false
		}
		return true // Already associated and matches
	}

	if user != nil {
		// Associate if slot empty
		if existingHandler, found := h.IsUserInActiveGame(user.ID); found && existingHandler.Session.ID != gameID {
			if !existingHandler.IsWaitingForCleanup(user.ID) {
				core.SendError(c, fmt.Sprintf("You are already in another active game (%s).", existingHandler.Session.ID), http.StatusConflict)
				return false
			}
			// stale/waiting for cleanup
			logging.Debug(logging.TagWeb, "User %d joining game %s, force cleaning up stale session %s", user.ID, gameID, existingHandler.Session.ID)
			h.RemoveSession(existingHandler.Session.ID, "Replaced by new session")
		}
	}

	return true // Allowed to connect (either as guest or new associate)
}

// HandleListGames returns a list of active games
// @Summary List active games
// @Description Retrieve a list of all currently active game sessions
// @Tags games
// @Produce json
// @Success 200 {array} models.GameSummary
// @Router /games [get]
func (h *Handler) HandleListGames(c *gin.Context) {
	h.Mutex.RLock()
	handlers := make([]*core.SessionHandler, 0, len(h.Sessions))
	for _, handler := range h.Sessions {
		handlers = append(handlers, handler)
	}
	h.Mutex.RUnlock()

	summaries := make([]models.GameSummary, 0, len(handlers))
	for _, handler := range handlers {
		summaries = append(summaries, handler.Session.GetGameSummary(handler.GameType))
	}

	core.SendJSON(c, summaries, http.StatusOK)
}
