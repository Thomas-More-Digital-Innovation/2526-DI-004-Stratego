package api

import (
	"digital-innovation/gostrategy/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleGetGameHistory handles GET /games/:id/history
// @Summary Get game history
// @Description Retrieve the full move history and initial setup of a finished game
// @Tags games
// @Produce json
// @Param id path string true "Game ID"
// @Success 200 {object} models.GameHistory
// @Failure 404 {object} map[string]string "Game not found"
// @Router /games/{id}/history [get]
func (s *GameServer) HandleGetGameHistory(c *gin.Context) {
	gameID := c.Param("id")
	if gameID == "" {
		sendError(c, "Game ID required", http.StatusBadRequest)
		return
	}

	history, err := db.GetGameHistory(c.Request.Context(), gameID)
	if err != nil {
		sendError(c, "Game history not found or error retrieving it", http.StatusNotFound)
		return
	}

	sendJSON(c, history, http.StatusOK)
}

// HandleListUserGames handles GET /users/:id/games?limit=10&offset=0
// @Summary List user's finished games
// @Description Retrieve a paged list of finished games for a specific user
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Param limit query int false "Maximum number of games to return" default(10)
// @Param offset query int false "Number of games to skip" default(0)
// @Success 200 {object} map[string]interface{} "List of games and total count"
// @Router /users/{id}/games [get]
func (s *GameServer) HandleListUserGames(c *gin.Context) {
	userID, err := parseID(c, "id")
	if err != nil || userID == 0 {
		sendError(c, "Invalid or missing user ID", http.StatusBadRequest)
		return
	}

	limit := getQueryInt(c, "limit", 10)
	offset := getQueryInt(c, "offset", 0)

	games, err := db.GetGamesForUserPaged(c.Request.Context(), userID, limit, offset)
	if err != nil {
		sendError(c, "Failed to retrieve game history", http.StatusInternalServerError)
		return
	}

	total, _ := db.GetGamesCountForUser(c.Request.Context(), userID)

	sendJSON(c, gin.H{
		"games":  games,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	}, http.StatusOK)
}

// HandleListMyGames handles GET /users/me/games
// @Summary List current user's finished games
// @Description Retrieve a paged list of finished games for the authenticated user
// @Tags users
// @Produce json
// @Param limit query int false "Maximum number of games to return" default(10)
// @Param offset query int false "Number of games to skip" default(0)
// @Success 200 {object} map[string]interface{} "List of games and total count"
// @Router /users/me/games [get]
func (s *GameServer) HandleListMyGames(c *gin.Context) {
	user := ensureAuthenticated(c)
	if user == nil {
		return
	}

	limit := getQueryInt(c, "limit", 10)
	offset := getQueryInt(c, "offset", 0)

	games, err := db.GetGamesForUserPaged(c.Request.Context(), user.ID, limit, offset)
	if err != nil {
		sendError(c, "Failed to retrieve game history", http.StatusInternalServerError)
		return
	}

	total, _ := db.GetGamesCountForUser(c.Request.Context(), user.ID)

	sendJSON(c, gin.H{
		"games":  games,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	}, http.StatusOK)
}
