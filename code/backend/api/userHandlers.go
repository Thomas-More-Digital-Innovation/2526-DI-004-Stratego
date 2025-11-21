package api

import (
	"digital-innovation/stratego/auth"
	"digital-innovation/stratego/db"
	"digital-innovation/stratego/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// RegisterUserHandler handles user registration
func (s *GameServer) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Username) < 3 || len(req.Username) > 50 {
		http.Error(w, "Username must be 3-50 characters", http.StatusBadRequest)
		return
	}
	if len(req.Password) < 6 {
		http.Error(w, "Password must be at least 6 characters", http.StatusBadRequest)
		return
	}

	user, err := db.CreateUser(req.Username, req.Password, req.ProfilePicture)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			http.Error(w, "Username already exists", http.StatusConflict)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to create user: %v", err), http.StatusInternalServerError)
		return
	}

	session, err := auth.Store.CreateSession(user.ID, user.Username)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	auth.SetSessionCookie(w, session.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// LoginHandler handles user login
func (s *GameServer) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := db.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	session, err := auth.Store.CreateSession(user.ID, user.Username)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	auth.SetSessionCookie(w, session.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// LogoutHandler handles user logout
func (s *GameServer) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err == nil {
		auth.Store.DeleteSession(cookie.Value)
	}

	auth.ClearSessionCookie(w)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

// GetCurrentUserHandler returns the currently logged-in user
func (s *GameServer) GetCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := auth.GetCurrentUser(r)
	if user == nil {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	fullUser, err := db.GetUserByID(user.ID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fullUser)
}

// GetUserHandler retrieves user information
func (s *GameServer) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.URL.Query().Get("id")
	if userIDStr == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := db.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

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
