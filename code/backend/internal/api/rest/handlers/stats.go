package handlers

import (
	"digital-innovation/gostrategy/internal/api/core"
	"digital-innovation/gostrategy/internal/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUserStatsHandler returns statistics for a user
// @Summary Get user statistics
// @Description Retrieve win/loss/draw statistics for a specific user
// @Tags stats
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.UserStats
// @Failure 404 {object} map[string]string "Stats not found"
// @Router /users/{id}/stats [get]
func (h *Handler) GetUserStatsHandler(c *gin.Context) {
	userID, _ := core.ParseID(c, "id")
	stats, err := db.GetUserStats(c.Request.Context(), userID)
	if err != nil {
		core.SendError(c, "Stats not found", http.StatusNotFound)
		return
	}
	core.SendJSON(c, stats, http.StatusOK)
}

// GetCurrentUserStatsHandler returns statistics for the authenticated user
// @Summary Get current user statistics
// @Description Retrieve win/loss/draw statistics for the currently logged-in user
// @Tags stats
// @Produce json
// @Success 200 {object} models.UserStats
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Stats not found"
// @Router /users/me/stats [get]
func (h *Handler) GetCurrentUserStatsHandler(c *gin.Context) {
	user := core.EnsureAuthenticated(c)
	if user == nil {
		return
	}
	stats, err := db.GetUserStats(c.Request.Context(), user.ID)
	if err != nil {
		core.SendError(c, "Stats not found", http.StatusNotFound)
		return
	}
	core.SendJSON(c, stats, http.StatusOK)
}
