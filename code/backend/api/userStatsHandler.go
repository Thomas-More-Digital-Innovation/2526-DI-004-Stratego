package api

import (
	"digital-innovation/stratego/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUserStatsHandler retrieves statistics for a specific user
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
