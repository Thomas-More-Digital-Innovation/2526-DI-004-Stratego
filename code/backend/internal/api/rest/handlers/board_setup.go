package handlers

import (
	"digital-innovation/gostrategy/internal/api/core"
	"digital-innovation/gostrategy/internal/db"
	"digital-innovation/gostrategy/internal/logging"
	"digital-innovation/gostrategy/internal/models"
	"digital-innovation/gostrategy/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUserBoardSetupsHandler returns all board setups for the authenticated user
// @Summary List board setups
// @Description Retrieve all saved board setups for the authenticated user
// @Tags board-setups
// @Produce json
// @Success 200 {array} models.BoardSetup
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /users/me/board-setups [get]
func (h *Handler) GetUserBoardSetupsHandler(c *gin.Context) {
	user := core.EnsureAuthenticated(c)
	if user == nil {
		return
	}

	setups, err := db.GetUserBoardSetups(c.Request.Context(), user.ID)
	if err != nil {
		core.SendError(c, "Failed to retrieve board setups", http.StatusInternalServerError)
		return
	}

	core.SendJSON(c, setups, http.StatusOK)
}

// CreateBoardSetupHandler creates a new board setup
// @Summary Create board setup
// @Description Save a new board setup configuration
// @Tags board-setups
// @Accept json
// @Produce json
// @Param setup body models.BoardSetup true "Board setup configuration"
// @Success 201 {object} models.BoardSetup
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /users/me/board-setups [post]
func (h *Handler) CreateBoardSetupHandler(c *gin.Context) {
	user := core.EnsureAuthenticated(c)
	if user == nil {
		return
	}

	var req models.BoardSetup
	if err := c.ShouldBindJSON(&req); err != nil {
		core.SendError(c, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := utils.ValidateGeneric(req.Name, "name"); err != nil {
		core.SendError(c, err.Error(), http.StatusBadRequest)
		return
	}
	if err := utils.ValidateGeneric(req.Description, "description"); err != nil {
		core.SendError(c, err.Error(), http.StatusBadRequest)
		return
	}

	req.UserID = user.ID
	setup, err := db.CreateBoardSetup(c.Request.Context(), user.ID, req.Name, req.Description, req.SetupData, req.IsDefault)
	if err != nil {
		core.SendError(c, "Failed to create board setup", http.StatusInternalServerError)
		return
	}

	logging.DebugWithUser(logging.TagGame, user.Username, user.ID, "Created board setup: %s", setup.Name)
	core.SendJSON(c, setup, http.StatusCreated)
}

// UpdateBoardSetupHandler updates an existing board setup
// @Summary Update board setup
// @Description Update an existing board setup configuration by ID
// @Tags board-setups
// @Accept json
// @Produce json
// @Param id path int true "Board Setup ID"
// @Param setup body models.BoardSetup true "Updated board setup configuration"
// @Success 200 {object} models.BoardSetup
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /users/me/board-setups/{id} [put]
func (h *Handler) UpdateBoardSetupHandler(c *gin.Context) {
	user := core.EnsureAuthenticated(c)
	if user == nil {
		return
	}

	id, _ := core.ParseID(c, "id")
	var req models.BoardSetup
	if err := c.ShouldBindJSON(&req); err != nil {
		core.SendError(c, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := utils.ValidateGeneric(req.Name, "name"); err != nil {
		core.SendError(c, err.Error(), http.StatusBadRequest)
		return
	}
	if err := utils.ValidateGeneric(req.Description, "description"); err != nil {
		core.SendError(c, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.UpdateBoardSetup(c.Request.Context(), int(id), user.ID, req.Name, req.Description, req.SetupData, req.IsDefault); err != nil {
		core.SendError(c, "Failed to update board setup", http.StatusInternalServerError)
		return
	}

	core.SendJSON(c, req, http.StatusOK)
}

// DeleteBoardSetupHandler deletes a board setup
// @Summary Delete board setup
// @Description Delete a specific board setup configuration by ID
// @Tags board-setups
// @Param id path int true "Board Setup ID"
// @Success 204 "No Content"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /users/me/board-setups/{id} [delete]
func (h *Handler) DeleteBoardSetupHandler(c *gin.Context) {
	user := core.EnsureAuthenticated(c)
	if user == nil {
		return
	}

	id, _ := core.ParseID(c, "id")
	if err := db.DeleteBoardSetup(c.Request.Context(), int(id), user.ID); err != nil {
		core.SendError(c, "Failed to delete board setup", http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

// GetBoardSetupHandler returns a specific board setup
// @Summary Get board setup
// @Description Retrieve a specific board setup configuration by ID
// @Tags board-setups
// @Produce json
// @Param id path int true "Board Setup ID"
// @Success 200 {object} models.BoardSetup
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Board setup not found"
// @Router /users/me/board-setups/{id} [get]
func (h *Handler) GetBoardSetupHandler(c *gin.Context) {
	user := core.EnsureAuthenticated(c)
	if user == nil {
		return
	}

	id, _ := core.ParseID(c, "id")
	setup, err := db.GetBoardSetup(c.Request.Context(), int(id), user.ID)
	if err != nil {
		core.SendError(c, "Board setup not found", http.StatusNotFound)
		return
	}

	core.SendJSON(c, setup, http.StatusOK)
}
