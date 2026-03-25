package api

import (
	"digital-innovation/stratego/auth"
	"digital-innovation/stratego/db"
	"digital-innovation/stratego/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// RegisterUserHandler handles user registration
func (s *GameServer) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Username) < 3 || len(req.Username) > 50 {
		sendError(w, "Username must be 3-50 characters", http.StatusBadRequest)
		return
	}
	if len(req.Password) < 6 {
		sendError(w, "Password must be at least 6 characters", http.StatusBadRequest)
		return
	}

	user, err := db.CreateUser(req.Username, req.Password, req.ProfilePicture)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			sendError(w, "Username already exists", http.StatusConflict)
			return
		}
		sendError(w, fmt.Sprintf("Failed to create user: %v", err), http.StatusInternalServerError)
		return
	}

	session, err := auth.Store.CreateSession(user.ID, user.Username)
	if err != nil {
		sendError(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	auth.SetSessionCookie(w, session.ID)

	sendJSON(w, user, http.StatusCreated)
}

// LoginHandler handles user login
func (s *GameServer) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := db.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		sendError(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	session, err := auth.Store.CreateSession(user.ID, user.Username)
	if err != nil {
		sendError(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	auth.SetSessionCookie(w, session.ID)

	sendJSON(w, user, http.StatusOK)
}

// LogoutHandler handles user logout
func (s *GameServer) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err == nil {
		auth.Store.DeleteSession(cookie.Value)
	}

	auth.ClearSessionCookie(w)

	sendJSON(w, map[string]string{"message": "Logged out successfully"}, http.StatusOK)
}

// GetCurrentUserHandler returns the currently logged-in user
func (s *GameServer) GetCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := ensureAuthenticated(w, r)
	if user == nil {
		return
	}

	fullUser, err := db.GetUserByID(user.ID)
	if err != nil {
		sendError(w, "User not found", http.StatusNotFound)
		return
	}

	sendJSON(w, fullUser, http.StatusOK)
}

// GetUserHandler retrieves user information
func (s *GameServer) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := parseID(r, "id")
	if err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := db.GetUserByID(userID)
	if err != nil {
		sendError(w, "User not found", http.StatusNotFound)
		return
	}

	sendJSON(w, user, http.StatusOK)
}
