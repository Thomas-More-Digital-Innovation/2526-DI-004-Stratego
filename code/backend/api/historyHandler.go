package api

import (
	"digital-innovation/stratego/db"
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

	history, err := db.GetGameHistory(gameID)
	if err != nil {
		sendError(c, "Game history not found or error retrieving it", http.StatusNotFound)
		return
	}

	sendJSON(c, history, http.StatusOK)
}
