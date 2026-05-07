package api

import (
	"digital-innovation/gostrategy/auth"
	"digital-innovation/gostrategy/db"
	"digital-innovation/gostrategy/logging"
	"digital-innovation/gostrategy/models"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
)

// isStrongPassword checks for password complexity
func isStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	var hasNumber, hasUpper, hasLower bool
	for _, char := range password {
		switch {
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		}
	}
	return hasNumber && hasUpper && hasLower
}

// RegisterUserHandler handles user registration
// @Summary Register a new user
// @Description Create a new user with username and password
// @Tags users
// @Accept json
// @Produce json
// @Param request body models.CreateUserRequest true "User registration details"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 409 {object} map[string]string "Username already exists"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/register [post]
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
		sendError(c, "Password must be at least 8 characters and contain at least one number, one uppercase and one lowercase letter", http.StatusBadRequest)
		return
	}

	user, err := db.CreateUser(c.Request.Context(), req.Username, req.Password, req.ProfilePicture)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			sendError(c, "Username already exists", http.StatusConflict)
			return
		}
		logging.ErrorWithIP("Failed to create user with username "+req.Username, c.ClientIP(), err)
		sendError(c, "Failed to create user", http.StatusInternalServerError)
		return
	}

	accessToken, err := auth.GenerateToken(user.ID, user.Username)
	if err != nil {
		logging.ErrorWithUser("Failed to generate access token", user.Username, user.ID, err)
		sendError(c, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		logging.ErrorWithUser("Failed to generate refresh token", user.Username, user.ID, err)
		sendError(c, "Failed to generate session", http.StatusInternalServerError)
		return
	}

	// Save refresh token to DB
	expiresAt := time.Now().Add(time.Duration(auth.MaxRefreshTokenAge) * time.Second)
	if err := db.SaveRefreshToken(c.Request.Context(), user.ID, refreshToken, expiresAt); err != nil {
		logging.ErrorWithUser("Failed to save refresh token", user.Username, user.ID, err)
		sendError(c, "Failed to persist session", http.StatusInternalServerError)
		return
	}

	auth.SetSessionCookie(c, accessToken)
	auth.SetRefreshTokenCookie(c, refreshToken)

	logging.UserRegistered(user.Username, user.ID)
	sendJSON(c, user, http.StatusCreated)
}

// LoginHandler handles user login
// @Summary User login
// @Description Authenticate user and create session
// @Tags users
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login details"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 401 {object} map[string]string "Invalid username or password"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/login [post]
func (s *GameServer) LoginHandler(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendError(c, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := db.AuthenticateUser(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		sendError(c, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	accessToken, err := auth.GenerateToken(user.ID, user.Username)
	if err != nil {
		logging.ErrorWithUser("Failed to generate access token", user.Username, user.ID, err)
		sendError(c, "Failed to create session", http.StatusInternalServerError)
		return
	}

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		logging.ErrorWithUser("Failed to generate refresh token", user.Username, user.ID, err)
		sendError(c, "Failed to create session", http.StatusInternalServerError)
		return
	}

	// Save refresh token to DB
	expiresAt := time.Now().Add(time.Duration(auth.MaxRefreshTokenAge) * time.Second)
	if err := db.SaveRefreshToken(c.Request.Context(), user.ID, refreshToken, expiresAt); err != nil {
		logging.ErrorWithUser("Failed to save refresh token", user.Username, user.ID, err)
		sendError(c, "Failed to persist session", http.StatusInternalServerError)
		return
	}

	auth.SetSessionCookie(c, accessToken)
	auth.SetRefreshTokenCookie(c, refreshToken)

	logging.UserLoggedIn(user.Username, user.ID)
	sendJSON(c, user, http.StatusOK)
}

// LogoutHandler handles user logout
// @Summary User logout
// @Description Delete user session and clear cookie
// @Tags users
// @Produce json
// @Success 200 {object} map[string]string "Logged out successfully"
// @Router /users/logout [post]
func (s *GameServer) LogoutHandler(c *gin.Context) {
	// Try to delete refresh token from DB if it exists
	if token, err := c.Cookie("refresh_token"); err == nil {
		_ = db.DeleteRefreshToken(c.Request.Context(), token)
	}

	auth.ClearSessionCookie(c)
	sendJSON(c, gin.H{MsgTypeMessage: "Logged out successfully"}, http.StatusOK)
}

// RefreshTokenHandler handles access token renewal using a refresh token
// @Summary Refresh access token
// @Description Use a refresh token to get a new access token
// @Tags users
// @Produce json
// @Success 200 {object} map[string]string "Success"
// @Failure 401 {object} map[string]string "Invalid refresh token"
// @Router /users/refresh [post]
func (s *GameServer) RefreshTokenHandler(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		sendError(c, "Refresh token missing", http.StatusUnauthorized)
		return
	}

	userID, err := db.GetUserIDByRefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		sendError(c, "Invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	user, err := db.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		sendError(c, "User no longer exists", http.StatusUnauthorized)
		return
	}

	accessToken, err := auth.GenerateToken(user.ID, user.Username)
	if err != nil {
		sendError(c, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	auth.SetSessionCookie(c, accessToken)
	sendJSON(c, gin.H{MsgTypeMessage: "Token refreshed"}, http.StatusOK)
}

// ChangePasswordHandler handles password updates
// @Summary Change user password
// @Description Update the authenticated user's password and revoke all sessions
// @Tags users
// @Accept json
// @Produce json
// @Param request body models.ChangePasswordRequest true "Password change details"
// @Success 200 {object} map[string]string "Password updated"
// @Failure 400 {object} map[string]string "Invalid request or weak password"
// @Failure 401 {object} map[string]string "Invalid old password"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/me/password [post]
func (s *GameServer) ChangePasswordHandler(c *gin.Context) {
	user := ensureAuthenticated(c)
	if user == nil {
		return
	}

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendError(c, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Verify old password
	_, err := db.AuthenticateUser(c.Request.Context(), user.Username, req.OldPassword)
	if err != nil {
		sendError(c, "Invalid old password", http.StatusUnauthorized)
		return
	}

	// Validate passwords match
	if req.NewPassword != req.ConfirmPassword {
		sendError(c, "New passwords do not match", http.StatusBadRequest)
		return
	}

	// Validate new password complexity
	if !isStrongPassword(req.NewPassword) {
		sendError(c, "New password must be at least 8 characters and contain at least one number, one uppercase and one lowercase letter", http.StatusBadRequest)
		return
	}

	if req.OldPassword == req.NewPassword {
		sendError(c, "New password must be different from the old one", http.StatusBadRequest)
		return
	}

	// Update password in DB
	if err := db.UpdateUserPassword(c.Request.Context(), user.ID, req.NewPassword); err != nil {
		logging.ErrorWithUser("Failed to update password", user.Username, user.ID, err)
		sendError(c, "Failed to update password", http.StatusInternalServerError)
		return
	}

	// Revoke all refresh tokens for this user (logs out all other devices)
	if err := db.DeleteAllUserRefreshTokens(c.Request.Context(), user.ID); err != nil {
		logging.ErrorWithUser("Failed to revoke refresh tokens after password change", user.Username, user.ID, err)
	}

	// Generate new session tokens for the current device
	accessToken, err := auth.GenerateToken(user.ID, user.Username)
	if err != nil {
		logging.ErrorWithUser("Failed to generate new access token after password change", user.Username, user.ID, err)
		sendError(c, "Password updated, but failed to refresh session. Please login again.", http.StatusInternalServerError)
		return
	}

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		logging.ErrorWithUser("Failed to generate new refresh token after password change", user.Username, user.ID, err)
		sendError(c, "Password updated, but failed to refresh session. Please login again.", http.StatusInternalServerError)
		return
	}

	// Save new refresh token
	expiresAt := time.Now().Add(time.Duration(auth.MaxRefreshTokenAge) * time.Second)
	if err := db.SaveRefreshToken(c.Request.Context(), user.ID, refreshToken, expiresAt); err != nil {
		logging.ErrorWithUser("Failed to save new refresh token after password change", user.Username, user.ID, err)
		sendError(c, "Password updated, but failed to persist session. Please login again.", http.StatusInternalServerError)
		return
	}

	// Set new cookies
	auth.SetSessionCookie(c, accessToken)
	auth.SetRefreshTokenCookie(c, refreshToken)

	logging.DebugWithUser(logging.TagAuth, user.Username, user.ID, "Password changed and session refreshed")
	sendJSON(c, gin.H{MsgTypeMessage: "Password updated successfully"}, http.StatusOK)
}

// GetCurrentUserHandler returns the currently logged-in user
// GetCurrentUserHandler returns the currently logged-in user
// @Summary Get current user
// @Description Retrieve profile of the authenticated user
// @Tags users
// @Produce json
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "User not found"
// @Router /users/me [get]
func (s *GameServer) GetCurrentUserHandler(c *gin.Context) {
	user := ensureAuthenticated(c)
	if user == nil {
		return
	}

	fullUser, err := db.GetUserByID(c.Request.Context(), user.ID)
	if err != nil {
		sendError(c, "User not found", http.StatusNotFound)
		return
	}

	sendJSON(c, fullUser, http.StatusOK)
}

// GetUserHandler retrieves user information
// GetUserHandler retrieves user information
// @Summary Get user profile
// @Description Retrieve profile of a user by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string "Invalid or missing user ID"
// @Failure 404 {object} map[string]string "User not found"
// @Router /users/{id} [get]
func (s *GameServer) GetUserHandler(c *gin.Context) {
	userID, err := parseID(c, "id")
	if err != nil || userID == 0 {
		sendError(c, "Invalid or missing user ID", http.StatusBadRequest)
		return
	}

	user, err := db.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		sendError(c, "User not found", http.StatusNotFound)
		return
	}

	sendJSON(c, user, http.StatusOK)
}
