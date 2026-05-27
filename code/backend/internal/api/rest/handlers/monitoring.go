package handlers

import (
	"digital-innovation/gostrategy/internal/api/core"
	"digital-innovation/gostrategy/internal/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserCountHandler returns the total number of registered users
// @Summary Get total user count
// @Description Retrieve the total number of registered users
// @Tags monitoring
// @Produce json
// @Success 200 {object} map[string]int64
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /monitoring/users/count [get]
func (h *Handler) UserCountHandler(c *gin.Context) {
	count, err := db.GetTotalUserCount(c.Request.Context())
	if err != nil {
		core.SendError(c, "Failed to retrieve user count", http.StatusInternalServerError)
		return
	}
	core.SendJSON(c, gin.H{"count": count}, http.StatusOK)
}

// GamesPlayedCountHandler returns the total number of games played
// @Summary Get total games played count
// @Description Retrieve the total number of games played in the system
// @Tags monitoring
// @Produce json
// @Success 200 {object} map[string]int64
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /monitoring/games/count [get]
func (h *Handler) GamesPlayedCountHandler(c *gin.Context) {
	count, err := db.GetTotalGamesPlayedCount(c.Request.Context())
	if err != nil {
		core.SendError(c, "Failed to retrieve games count", http.StatusInternalServerError)
		return
	}
	core.SendJSON(c, gin.H{"count": count}, http.StatusOK)
}
