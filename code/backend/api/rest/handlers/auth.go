// Package handlers provides handlers for the REST API.
package handlers

import (
	"digital-innovation/gostrategy/api/core"
	"digital-innovation/gostrategy/api/dto"
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

// Handler wraps the GameServer to provide REST handlers

// IsStrongPassword checks for password complexity
func IsStrongPassword(password string) bool {
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
// @Description Create a new user account with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "User registration details"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string "Invalid request body or validation error"
// @Failure 409 {object} map[string]string "Username already exists"
// @Failure 500 {object} map[string]string "Failed to create user"
// @Router /users/register [post]
func (h *Handler) RegisterUserHandler(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.SendError(c, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Username) < 3 || len(req.Username) > 50 {
		core.SendError(c, "Username must be 3-50 characters", http.StatusBadRequest)
		return
	}
	if !IsStrongPassword(req.Password) {
		core.SendError(c, "Password must be at least 8 characters and contain at least one number, one uppercase and one lowercase letter", http.StatusBadRequest)
		return
	}

	user, err := db.CreateUser(c.Request.Context(), req.Username, req.Password, req.ProfilePicture)
	if err != nil {
		errLower := strings.ToLower(err.Error())
		if strings.Contains(errLower, "duplicate") || strings.Contains(errLower, "unique") {
			core.SendError(c, "Username already exists", http.StatusConflict)
			return
		}
		logging.ErrorWithIP("Failed to create user with username "+req.Username, c.ClientIP(), err)
		core.SendError(c, "Failed to create user", http.StatusInternalServerError)
		return
	}

	accessToken, _ := auth.GenerateToken(user.ID, user.Username)
	refreshToken, _ := auth.GenerateRefreshToken()
	expiresAt := time.Now().Add(time.Duration(auth.MaxRefreshTokenAge) * time.Second)
	_ = db.SaveRefreshToken(c.Request.Context(), user.ID, refreshToken, expiresAt)

	auth.SetSessionCookie(c, accessToken)
	auth.SetRefreshTokenCookie(c, refreshToken)

	logging.UserRegistered(user.Username, user.ID)
	core.SendJSON(c, user, http.StatusCreated)
}

// LoginHandler handles user login
// @Summary User login
// @Description Authenticate a user and return user details with session cookies
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 401 {object} map[string]string "Invalid username or password"
// @Router /users/login [post]
func (h *Handler) LoginHandler(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.SendError(c, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := db.AuthenticateUser(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		core.SendError(c, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	accessToken, _ := auth.GenerateToken(user.ID, user.Username)
	refreshToken, _ := auth.GenerateRefreshToken()
	expiresAt := time.Now().Add(time.Duration(auth.MaxRefreshTokenAge) * time.Second)
	_ = db.SaveRefreshToken(c.Request.Context(), user.ID, refreshToken, expiresAt)

	auth.SetSessionCookie(c, accessToken)
	auth.SetRefreshTokenCookie(c, refreshToken)

	logging.UserLoggedIn(user.Username, user.ID)
	core.SendJSON(c, user, http.StatusOK)
}

// LogoutHandler handles user logout
// @Summary User logout
// @Description Clear user session and delete refresh token
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]string "Logged out successfully"
// @Router /users/logout [post]
func (h *Handler) LogoutHandler(c *gin.Context) {
	if token, err := c.Cookie("refresh_token"); err == nil {
		_ = db.DeleteRefreshToken(c.Request.Context(), token)
	}

	auth.ClearSessionCookie(c)
	core.SendJSON(c, gin.H{dto.MsgTypeMessage: "Logged out successfully"}, http.StatusOK)
}

// RefreshTokenHandler handles access token renewal
// @Summary Refresh access token
// @Description Renew the access token using a valid refresh token from cookies
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]string "Token refreshed"
// @Failure 401 {object} map[string]string "Refresh token missing or invalid"
// @Router /users/refresh [post]
func (h *Handler) RefreshTokenHandler(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		core.SendError(c, "Refresh token missing", http.StatusUnauthorized)
		return
	}

	userID, err := db.GetUserIDByRefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		core.SendError(c, "Invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	user, _ := db.GetUserByID(c.Request.Context(), userID)
	accessToken, _ := auth.GenerateToken(user.ID, user.Username)

	auth.SetSessionCookie(c, accessToken)
	core.SendJSON(c, gin.H{dto.MsgTypeMessage: "Token refreshed"}, http.StatusOK)
}
