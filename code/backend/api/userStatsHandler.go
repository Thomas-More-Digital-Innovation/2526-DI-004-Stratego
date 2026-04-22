package api

import (
	"digital-innovation/stratego/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUserStatsHandler retrieves statistics for a specific user
// GetUserStatsHandler retrieves statistics for a specific user
// @Summary Get user statistics
// @Description Retrieve game statistics for a specific user or the user specified by "id" query/path param
// @Tags users
// @Produce json
// @Param id query int false "User ID"
// @Success 200 {object} models.UserStats
// @Failure 400 {object} map[string]string "Invalid or missing user ID"
// @Failure 404 {object} map[string]string "Stats not found"
// @Router /users/stats [get]
func (s *GameServer) GetUserStatsHandler(c *gin.Context) {
	userID, err := parseID(c, "id")
	if err != nil || userID == 0 {
		sendError(c, "Invalid or missing user ID", http.StatusBadRequest)
		return
	}

	stats, err := db.GetUserStats(userID)
	if err != nil {
		sendError(c, "Stats not found", http.StatusNotFound)
		return
	}

	sendJSON(c, stats, http.StatusOK)
}

// GetCurrentUserStatsHandler retrieves statistics for the currently logged-in user
// GetCurrentUserStatsHandler retrieves statistics for the currently logged-in user
// @Summary Get current user statistics
// @Description Retrieve game statistics for the authenticated user
// @Tags users
// @Produce json
// @Success 200 {object} models.UserStats
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Stats not found"
// @Router /users/me/stats [get]
func (s *GameServer) GetCurrentUserStatsHandler(c *gin.Context) {
	user := ensureAuthenticated(c)
	if user == nil {
		return
	}

	stats, err := db.GetUserStats(user.ID)
	if err != nil {
		sendError(c, "Stats not found", http.StatusNotFound)
		return
	}

	sendJSON(c, stats, http.StatusOK)
}
