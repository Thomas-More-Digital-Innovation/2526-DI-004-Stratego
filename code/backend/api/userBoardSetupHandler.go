package api

import (
	"digital-innovation/stratego/db"
	"digital-innovation/stratego/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateBoardSetupHandler creates a new board setup
// @Summary Create board setup
// @Description Save a new piece configuration for the user
// @Tags board-setups
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body models.CreateBoardSetupRequest true "Setup details"
// @Success 201 {object} models.BoardSetup
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /board-setups [post]
func (s *GameServer) CreateBoardSetupHandler(c *gin.Context) {
	user := ensureAuthenticated(c)
	if user == nil {
		return
	}

	var req models.CreateBoardSetupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendError(c, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		sendError(c, "Setup name is required", http.StatusBadRequest)
		return
	}

	setup, err := db.CreateBoardSetup(user.ID, req.Name, req.Description, req.SetupData, req.IsDefault)
	if err != nil {
		log.Printf("Failed to create board setup: %v", err)
		sendError(c, "Failed to create board setup", http.StatusInternalServerError)
		return
	}

	sendJSON(c, setup, http.StatusCreated)
}

// GetUserBoardSetupsHandler retrieves all setups for a user
// @Summary List board setups
// @Description Retrieve all board setups for the authenticated user
// @Tags board-setups
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} models.BoardSetup
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /board-setups [get]
func (s *GameServer) GetUserBoardSetupsHandler(c *gin.Context) {
	user := ensureAuthenticated(c)
	if user == nil {
		return
	}

	setups, err := db.GetUserBoardSetups(user.ID)
	if err != nil {
		log.Printf("Failed to get board setups: %v", err)
		sendError(c, "Failed to get board setups", http.StatusInternalServerError)
		return
	}

	sendJSON(c, setups, http.StatusOK)
}

// GetBoardSetupHandler retrieves a single board setup
// @Summary Get board setup
// @Description Retrieve a specific board setup by ID
// @Tags board-setups
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Setup ID"
// @Success 200 {object} models.BoardSetup
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Setup not found"
// @Router /board-setups/{id} [get]
func (s *GameServer) GetBoardSetupHandler(c *gin.Context) {
	user := ensureAuthenticated(c)
	if user == nil {
		return
	}

	setupID, err := parseID(c, "id")
	if err != nil || setupID == 0 {
		sendError(c, "Invalid or missing setup ID", http.StatusBadRequest)
		return
	}

	setup, err := db.GetBoardSetup(setupID, user.ID)
	if err != nil {
		sendError(c, "Setup not found or not owned by user", http.StatusNotFound)
		return
	}

	sendJSON(c, setup, http.StatusOK)
}

// UpdateBoardSetupHandler updates an existing board setup
// @Summary Update board setup
// @Description Modify an existing board setup by ID
// @Tags board-setups
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Setup ID"
// @Param request body models.UpdateBoardSetupRequest true "Updated details"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /board-setups/{id} [put]
func (s *GameServer) UpdateBoardSetupHandler(c *gin.Context) {
	user := ensureAuthenticated(c)
	if user == nil {
		return
	}

	setupID, err := parseID(c, "id")
	if err != nil || setupID == 0 {
		sendError(c, "Invalid or missing setup ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateBoardSetupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendError(c, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = db.UpdateBoardSetup(setupID, user.ID, req.Name, req.Description, req.SetupData, req.IsDefault)
	if err != nil {
		log.Printf("Failed to update board setup: %v", err)
		sendError(c, "Failed to update board setup", http.StatusInternalServerError)
		return
	}

	sendNoContent(c)
}

// DeleteBoardSetupHandler deletes a board setup
// @Summary Delete board setup
// @Description Remove a board setup by ID
// @Tags board-setups
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Setup ID"
// @Success 204 "No Content"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /board-setups/{id} [delete]
func (s *GameServer) DeleteBoardSetupHandler(c *gin.Context) {
	user := ensureAuthenticated(c)
	if user == nil {
		return
	}

	setupID, err := parseID(c, "id")
	if err != nil || setupID == 0 {
		sendError(c, "Invalid or missing setup ID", http.StatusBadRequest)
		return
	}

	err = db.DeleteBoardSetup(setupID, user.ID)
	if err != nil {
		log.Printf("Failed to delete board setup: %v", err)
		sendError(c, "Failed to delete board setup", http.StatusInternalServerError)
		return
	}

	sendNoContent(c)
}
