package api

import (
	"digital-innovation/stratego/db"
	"encoding/json"
	"net/http"
	"strconv"
)

// GetUserStatsHandler retrieves user statistics
func (s *GameServer) GetUserStatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	stats, err := db.GetUserStats(userID)
	if err != nil {
		http.Error(w, "Stats not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
