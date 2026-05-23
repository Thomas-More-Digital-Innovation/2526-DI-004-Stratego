package handlers

import (
	"digital-innovation/gostrategy/api/core"
	"digital-innovation/gostrategy/db"
	"digital-innovation/gostrategy/logging"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleGetGameHistory returns the history of a specific game
// @Summary Get game history
// @Description Retrieve the move history and final state of a specific game
// @Tags games
// @Produce json
// @Param id path string true "Game ID"
// @Success 200 {object} models.Game
// @Failure 400 {object} map[string]string "Game ID required"
// @Failure 404 {object} map[string]string "Game history not found"
// @Router /games/{id}/history [get]
func (h *Handler) HandleGetGameHistory(c *gin.Context) {
	gameID := c.Param("id")
	if gameID == "" {
		core.SendError(c, "Game ID required", http.StatusBadRequest)
		return
	}

	gameData, err := db.GetGameHistory(c.Request.Context(), gameID)
	if err != nil {
		core.SendError(c, "Game history not found", http.StatusNotFound)
		return
	}

	core.SendJSON(c, gameData, http.StatusOK)
}

// HandleListUserGames returns a list of games for a specific user
// @Summary List user games
// @Description Retrieve a list of games played by a specific user (paged)
// @Tags games
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string "Failed to retrieve games"
// @Router /users/{id}/games [get]
func (h *Handler) HandleListUserGames(c *gin.Context) {
	userID, _ := core.ParseID(c, "id")
	games, err := db.GetGamesForUserPaged(c.Request.Context(), userID, 50, 0)
	if err != nil {
		core.SendError(c, "Failed to retrieve games", http.StatusInternalServerError)
		return
	}

	total, _ := db.GetGamesCountForUser(c.Request.Context(), userID)

	core.SendJSON(c, gin.H{
		"games": games,
		"total": total,
	}, http.StatusOK)
}

// HandleListMyGames returns a list of games for the authenticated user
// @Summary List current user games
// @Description Retrieve a list of games played by the currently logged-in user
// @Tags games
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Failed to retrieve games"
// @Router /users/me/games [get]
func (h *Handler) HandleListMyGames(c *gin.Context) {
	user := core.EnsureAuthenticated(c)
	if user == nil {
		return
	}
	games, err := db.GetGamesForUserPaged(c.Request.Context(), user.ID, 50, 0)
	if err != nil {
		core.SendError(c, "Failed to retrieve games", http.StatusInternalServerError)
		return
	}

	total, _ := db.GetGamesCountForUser(c.Request.Context(), user.ID)

	logging.Debug(logging.TagWeb, "Listing games for user %d: found %d games", user.ID, total)

	core.SendJSON(c, gin.H{
		"games": games,
		"total": total,
	}, http.StatusOK)
}

// HandleGetReconnectableGame handles querying for a reconnectable game session.
// @Summary Get current user's reconnectable game
// @Description Retrieve a game session in progress for the authenticated user
// @Tags games
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /users/me/reconnectable [get]
func (h *Handler) HandleGetReconnectableGame(c *gin.Context) {
	user := core.EnsureAuthenticated(c)
	if user == nil {
		return
	}

	handler, seatIndex, found := h.GetUserActiveGameSeat(user.ID)
	if !found {
		core.SendJSON(c, gin.H{"hasGame": false}, http.StatusOK)
		return
	}

	if handler.Session.GetGameState().IsGameOver {
		core.SendJSON(c, gin.H{"hasGame": false}, http.StatusOK)
		return
	}

	core.SendJSON(c, gin.H{
		"hasGame":   true,
		"gameId":    handler.Session.ID,
		"gameType":  handler.GameType,
		"seatIndex": seatIndex,
	}, http.StatusOK)
}
