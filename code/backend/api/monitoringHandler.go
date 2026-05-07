package api

import (
	"digital-innovation/gostrategy/auth"
	"digital-innovation/gostrategy/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler returns the server status
// @Summary Health check
// @Description Confirm the server is running
// @Tags monitoring
// @Produce json
// @Success 200 {object} map[string]string "Status OK"
// @Router /health [get]
func (s *GameServer) HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// UserCountHandler returns the total number of users
// @Summary User count
// @Description Get the total number of users
// @Tags monitoring
// @Produce json
// @Success 200 {object} map[string]int "User count"
// @Router /users/count [get]
func (s *GameServer) UserCountHandler(c *gin.Context) {
	count, err := db.GetTotalUserCount(c.Request.Context())
	if err != nil {
		sendError(c, "Failed to get user count", http.StatusInternalServerError)
		return
	}
	sendJSON(c, gin.H{"count": count}, http.StatusOK)
}

// GamesPlayedCountHandler returns the total number of games played
// @Summary Games played count
// @Description Get the total number of games played
// @Tags monitoring
// @Produce json
// @Success 200 {object} map[string]int "Games played count"
// @Router /games/count [get]
func (s *GameServer) GamesPlayedCountHandler(c *gin.Context) {
	count, err := db.GetTotalGamesPlayedCount(c.Request.Context())
	if err != nil {
		sendError(c, "Failed to get games played count", http.StatusInternalServerError)
		return
	}
	sendJSON(c, gin.H{"count": count}, http.StatusOK)
}

// GetCSRFToken provides a new CSRF token to guests
// @Summary Get CSRF token
// @Description Set the XSRF-TOKEN cookie and return the token
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]string "CSRF token"
// @Router /csrf-token [get]
func (s *GameServer) GetCSRFToken(c *gin.Context) {
	token := auth.SetCSRFCookie(c)
	c.JSON(http.StatusOK, gin.H{"csrf_token": token})
}
