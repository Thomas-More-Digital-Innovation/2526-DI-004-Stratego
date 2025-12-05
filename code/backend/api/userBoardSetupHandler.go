package api

import (
	"digital-innovation/stratego/auth"
	"digital-innovation/stratego/db"
	"digital-innovation/stratego/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// CreateBoardSetupHandler creates a new board setup
func (s *GameServer) CreateBoardSetupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := auth.GetCurrentUser(r)
	if user == nil {
		http.Error(w, "Unauthorized: Please login", http.StatusUnauthorized)
		return
	}

	var req models.CreateBoardSetupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Setup name is required", http.StatusBadRequest)
		return
	}

	setup, err := db.CreateBoardSetup(user.ID, req.Name, req.Description, req.SetupData, req.IsDefault)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create board setup: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(setup)
}

// GetUserBoardSetupsHandler retrieves all setups for a user
func (s *GameServer) GetUserBoardSetupsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := auth.GetCurrentUser(r)
	if user == nil {
		http.Error(w, "Unauthorized: Please login", http.StatusUnauthorized)
		return
	}

	setups, err := db.GetUserBoardSetups(user.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get board setups: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(setups)
}

// GetBoardSetupHandler retrieves a single board setup
func (s *GameServer) GetBoardSetupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	setupIDStr := r.URL.Query().Get("id")
	if setupIDStr == "" {
		http.Error(w, "Missing setup ID", http.StatusBadRequest)
		return
	}

	setupID, err := strconv.Atoi(setupIDStr)
	if err != nil {
		http.Error(w, "Invalid setup ID", http.StatusBadRequest)
		return
	}

	setup, err := db.GetBoardSetup(setupID)
	if err != nil {
		http.Error(w, "Setup not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(setup)
}

// UpdateBoardSetupHandler updates an existing board setup
func (s *GameServer) UpdateBoardSetupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	setupIDStr := r.URL.Query().Get("id")
	if setupIDStr == "" {
		http.Error(w, "Missing setup ID", http.StatusBadRequest)
		return
	}

	setupID, err := strconv.Atoi(setupIDStr)
	if err != nil {
		http.Error(w, "Invalid setup ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateBoardSetupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = db.UpdateBoardSetup(setupID, req.Name, req.Description, req.SetupData, req.IsDefault)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update board setup: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteBoardSetupHandler deletes a board setup
func (s *GameServer) DeleteBoardSetupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	setupIDStr := r.URL.Query().Get("id")
	if setupIDStr == "" {
		http.Error(w, "Missing setup ID", http.StatusBadRequest)
		return
	}

	setupID, err := strconv.Atoi(setupIDStr)
	if err != nil {
		http.Error(w, "Invalid setup ID", http.StatusBadRequest)
		return
	}

	// Get user from session
	user := auth.GetCurrentUser(r)
	if user == nil {
		http.Error(w, "Unauthorized: Please login", http.StatusUnauthorized)
		return
	}

	err = db.DeleteBoardSetup(setupID, user.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete board setup: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
