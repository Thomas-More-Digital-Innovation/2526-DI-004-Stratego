package api

import (
	"digital-innovation/stratego/auth"
	"digital-innovation/stratego/db"
	"digital-innovation/stratego/models"
	"log"
	"net/http"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
)

// isStrongPassword checks for password complexity
func isStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	var hasNumber bool
	for _, char := range password {
		if unicode.IsDigit(char) {
			hasNumber = true
			break
		}
	}
	return hasNumber
}

// RegisterUserHandler handles user registration
func (s *GameServer) RegisterUserHandler(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendError(c, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Username) < 3 || len(req.Username) > 50 {
		sendError(c, "Username must be 3-50 characters", http.StatusBadRequest)
		return
	}
	if !isStrongPassword(req.Password) {
		sendError(c, "Password must be at least 8 characters and contain at least one number", http.StatusBadRequest)
		return
	}

	user, err := db.CreateUser(req.Username, req.Password, req.ProfilePicture)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			sendError(c, "Username already exists", http.StatusConflict)
			return
		}
		log.Printf("Failed to create user: %v", err)
		sendError(c, "Failed to create user", http.StatusInternalServerError)
		return
	}

	session, err := auth.Store.CreateSession(user.ID, user.Username)
	if err != nil {
		log.Printf("Failed to create session: %v", err)
		sendError(c, "Failed to create session", http.StatusInternalServerError)
		return
	}

	auth.SetSessionCookie(c, session.ID)

	sendJSON(c, user, http.StatusCreated)
}

// LoginHandler handles user login
func (s *GameServer) LoginHandler(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendError(c, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := db.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		sendError(c, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	session, err := auth.Store.CreateSession(user.ID, user.Username)
	if err != nil {
		log.Printf("Failed to create session: %v", err)
		sendError(c, "Failed to create session", http.StatusInternalServerError)
		return
	}

	auth.SetSessionCookie(c, session.ID)

	sendJSON(c, user, http.StatusOK)
}

// LogoutHandler handles user logout
func (s *GameServer) LogoutHandler(c *gin.Context) {
	cookie, err := c.Cookie("session_id")
	if err == nil {
		auth.Store.DeleteSession(cookie)
	}

	auth.ClearSessionCookie(c)

	sendJSON(c, gin.H{"message": "Logged out successfully"}, http.StatusOK)
}

// GetCurrentUserHandler returns the currently logged-in user
func (s *GameServer) GetCurrentUserHandler(c *gin.Context) {
	user := ensureAuthenticated(c)
	if user == nil {
		return
	}

	fullUser, err := db.GetUserByID(user.ID)
	if err != nil {
		sendError(c, "User not found", http.StatusNotFound)
		return
	}

	sendJSON(c, fullUser, http.StatusOK)
}

// GetUserHandler retrieves user information
func (s *GameServer) GetUserHandler(c *gin.Context) {
	userID, err := parseID(c, "id")
	if err != nil || userID == 0 {
		sendError(c, "Invalid or missing user ID", http.StatusBadRequest)
		return
	}

	user, err := db.GetUserByID(userID)
	if err != nil {
		sendError(c, "User not found", http.StatusNotFound)
		return
	}

	sendJSON(c, user, http.StatusOK)
}
