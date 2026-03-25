package api

import (
	"digital-innovation/stratego/db"
	"net/http"
)

// GetUserStatsHandler retrieves user statistics
func (s *GameServer) GetUserStatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := parseID(r, "user_id")
	if err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	stats, err := db.GetUserStats(userID)
	if err != nil {
		sendError(w, "Stats not found", http.StatusNotFound)
		return
	}

	sendJSON(w, stats, http.StatusOK)
}

// GetCurrentUserStatsHandler retrieves statistics for the currently authenticated user
func (s *GameServer) GetCurrentUserStatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := ensureAuthenticated(w, r)
	if user == nil {
		return
	}

	stats, err := db.GetUserStats(user.ID)
	if err != nil {
		sendError(w, "Stats not found", http.StatusNotFound)
		return
	}

	sendJSON(w, stats, http.StatusOK)
}
