package api

import (
	"digital-innovation/stratego/db"
	"digital-innovation/stratego/models"
	"encoding/json"
	"fmt"
	"net/http"
)

// CreateBoardSetupHandler creates a new board setup
func (s *GameServer) CreateBoardSetupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := ensureAuthenticated(w, r)
	if user == nil {
		return
	}

	var req models.CreateBoardSetupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		sendError(w, "Setup name is required", http.StatusBadRequest)
		return
	}

	setup, err := db.CreateBoardSetup(user.ID, req.Name, req.Description, req.SetupData, req.IsDefault)
	if err != nil {
		sendError(w, fmt.Sprintf("Failed to create board setup: %v", err), http.StatusInternalServerError)
		return
	}

	sendJSON(w, setup, http.StatusCreated)
}

// GetUserBoardSetupsHandler retrieves all setups for a user
func (s *GameServer) GetUserBoardSetupsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := ensureAuthenticated(w, r)
	if user == nil {
		return
	}

	setups, err := db.GetUserBoardSetups(user.ID)
	if err != nil {
		sendError(w, fmt.Sprintf("Failed to get board setups: %v", err), http.StatusInternalServerError)
		return
	}

	sendJSON(w, setups, http.StatusOK)
}

// GetBoardSetupHandler retrieves a single board setup
func (s *GameServer) GetBoardSetupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	setupID, err := parseID(r, "id")
	if err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	setup, err := db.GetBoardSetup(setupID)
	if err != nil {
		sendError(w, "Setup not found", http.StatusNotFound)
		return
	}

	sendJSON(w, setup, http.StatusOK)
}

// UpdateBoardSetupHandler updates an existing board setup
func (s *GameServer) UpdateBoardSetupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	setupID, err := parseID(r, "id")
	if err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req models.UpdateBoardSetupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = db.UpdateBoardSetup(setupID, req.Name, req.Description, req.SetupData, req.IsDefault)
	if err != nil {
		sendError(w, fmt.Sprintf("Failed to update board setup: %v", err), http.StatusInternalServerError)
		return
	}

	sendNoContent(w)
}

// DeleteBoardSetupHandler deletes a board setup
func (s *GameServer) DeleteBoardSetupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	setupID, err := parseID(r, "id")
	if err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := ensureAuthenticated(w, r)
	if user == nil {
		return
	}

	err = db.DeleteBoardSetup(setupID, user.ID)
	if err != nil {
		sendError(w, fmt.Sprintf("Failed to delete board setup: %v", err), http.StatusInternalServerError)
		return
	}

	sendNoContent(w)
}
