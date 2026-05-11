package handlers

import (
	"digital-innovation/gostrategy/api/core"
	"digital-innovation/gostrategy/api/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler returns the health status of the server
// @Summary Health check
// @Description Check if the API server is up and running
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *Handler) HealthHandler(c *gin.Context) {
	core.SendJSON(c, gin.H{"status": "UP"}, http.StatusOK)
}

// GetCSRFToken returns a message for CSRF token retrieval
// @Summary Get CSRF token
// @Description Retrieve a CSRF token for subsequent state-changing requests
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]string
// @Router /csrf [get]
func (h *Handler) GetCSRFToken(c *gin.Context) {
	// The CSRF token is already set in a cookie by the middleware
	core.SendJSON(c, gin.H{dto.MsgTypeMessage: "CSRF cookie set"}, http.StatusOK)
}
