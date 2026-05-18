package handlers_test

import (
	"bytes"
	"context"
	"digital-innovation/gostrategy/api/core"
	"digital-innovation/gostrategy/api/rest/handlers"
	"digital-innovation/gostrategy/db"
	"digital-innovation/gostrategy/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const strongPassword = "StrongPassword1"

func TestIsStrongPassword(t *testing.T) {
	assert.True(t, handlers.IsStrongPassword("Strong123"))
	assert.False(t, handlers.IsStrongPassword("weak"))
	assert.False(t, handlers.IsStrongPassword("NoNumber"))
	assert.False(t, handlers.IsStrongPassword("nonupper1"))
	assert.False(t, handlers.IsStrongPassword("NOLOWER1"))
}

func TestRegisterUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db.SetupTestDB(t)

	server := core.NewGameServer()
	h := handlers.NewHandler(server)

	// Test valid registration
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	reqBody := models.CreateUserRequest{
		Username: "alice",
		Password: strongPassword,
	}
	jsonBody, _ := json.Marshal(reqBody) // nolint:gosec // G117: marshaling password in test
	c.Request, _ = http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	h.RegisterUserHandler(c)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Test duplicate registration
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request, _ = http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
	c2.Request.Header.Set("Content-Type", "application/json")
	h.RegisterUserHandler(c2)
	assert.Equal(t, http.StatusConflict, w2.Code)
}

func TestLoginHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db.SetupTestDB(t)
	server := core.NewGameServer()
	h := handlers.NewHandler(server)

	// Register a user first
	reqBody := models.CreateUserRequest{
		Username: "bob",
		Password: strongPassword,
	}
	jsonBody, _ := json.Marshal(reqBody) // nolint:gosec // G117: marshaling password in test
	wReg := httptest.NewRecorder()
	cReg, _ := gin.CreateTestContext(wReg)
	cReg.Request, _ = http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
	cReg.Request.Header.Set("Content-Type", "application/json")
	h.RegisterUserHandler(cReg)

	// Test valid login
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	loginBody := models.LoginRequest{
		Username: "bob",
		Password: strongPassword,
	}
	jsonLogin, _ := json.Marshal(loginBody) // nolint:gosec // G117: marshaling password in test
	c.Request, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(jsonLogin))
	c.Request.Header.Set("Content-Type", "application/json")

	h.LoginHandler(c)
	assert.Equal(t, http.StatusOK, w.Code)

	// Test invalid login
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	loginBody.Password = "wrong"
	jsonLogin2, _ := json.Marshal(loginBody) // nolint:gosec // G117: marshaling password in test
	c2.Request, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(jsonLogin2))
	c2.Request.Header.Set("Content-Type", "application/json")
	h.LoginHandler(c2)
	assert.Equal(t, http.StatusUnauthorized, w2.Code)
}

func TestLogoutHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db.SetupTestDB(t)
	server := core.NewGameServer()
	h := handlers.NewHandler(server)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/logout", nil)

	h.LogoutHandler(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRefreshTokenHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db.SetupTestDB(t)
	server := core.NewGameServer()
	h := handlers.NewHandler(server)

	// Create a user and a refresh token
	user, _ := db.CreateUser(context.Background(), "alice", strongPassword, "")
	refreshToken := "some-refresh-token"
	expiresAt := time.Now().Add(24 * time.Hour)
	ctx := db.WithUserID(context.Background(), user.ID)
	err := db.SaveRefreshToken(ctx, user.ID, refreshToken, expiresAt)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/refresh", nil)
	c.Request.AddCookie(&http.Cookie{Name: "refresh_token", Value: refreshToken, HttpOnly: true, Secure: true, SameSite: http.SameSiteStrictMode})

	h.RefreshTokenHandler(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRefreshTokenHandler_Errors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db.SetupTestDB(t)
	server := core.NewGameServer()
	h := handlers.NewHandler(server)

	// Case 1: Missing refresh token cookie
	w1 := httptest.NewRecorder()
	c1, _ := gin.CreateTestContext(w1)
	c1.Request, _ = http.NewRequest("POST", "/refresh", nil)
	h.RefreshTokenHandler(c1)
	assert.Equal(t, http.StatusUnauthorized, w1.Code)

	// Case 2: Invalid/non-existent refresh token
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request, _ = http.NewRequest("POST", "/refresh", nil)
	c2.Request.AddCookie(&http.Cookie{Name: "refresh_token", Value: "invalid-token"})
	h.RefreshTokenHandler(c2)
	assert.Equal(t, http.StatusUnauthorized, w2.Code)
}

func TestRegisterUserHandler_Errors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db.SetupTestDB(t)
	server := core.NewGameServer()
	h := handlers.NewHandler(server)

	// Case 1: Invalid request body JSON
	w1 := httptest.NewRecorder()
	c1, _ := gin.CreateTestContext(w1)
	c1.Request, _ = http.NewRequest("POST", "/register", bytes.NewBufferString("invalid-json"))
	h.RegisterUserHandler(c1)
	assert.Equal(t, http.StatusBadRequest, w1.Code)

	// Case 2: Too short username
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	req2 := models.CreateUserRequest{Username: "a", Password: strongPassword}
	json2, _ := json.Marshal(req2)
	c2.Request, _ = http.NewRequest("POST", "/register", bytes.NewBuffer(json2))
	c2.Request.Header.Set("Content-Type", "application/json")
	h.RegisterUserHandler(c2)
	assert.Equal(t, http.StatusBadRequest, w2.Code)

	// Case 3: Weak password
	w3 := httptest.NewRecorder()
	c3, _ := gin.CreateTestContext(w3)
	req3 := models.CreateUserRequest{Username: "validuser", Password: "weak"}
	json3, _ := json.Marshal(req3)
	c3.Request, _ = http.NewRequest("POST", "/register", bytes.NewBuffer(json3))
	c3.Request.Header.Set("Content-Type", "application/json")
	h.RegisterUserHandler(c3)
	assert.Equal(t, http.StatusBadRequest, w3.Code)
}

func TestLoginHandler_Errors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db.SetupTestDB(t)
	server := core.NewGameServer()
	h := handlers.NewHandler(server)

	// Case 1: Invalid JSON body
	w1 := httptest.NewRecorder()
	c1, _ := gin.CreateTestContext(w1)
	c1.Request, _ = http.NewRequest("POST", "/login", bytes.NewBufferString("invalid-json"))
	h.LoginHandler(c1)
	assert.Equal(t, http.StatusBadRequest, w1.Code)

	// Case 2: Invalid username format (e.g. invalid chars)
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	req2 := models.LoginRequest{Username: "invalid@user", Password: strongPassword}
	json2, _ := json.Marshal(req2)
	c2.Request, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(json2))
	c2.Request.Header.Set("Content-Type", "application/json")
	h.LoginHandler(c2)
	assert.Equal(t, http.StatusBadRequest, w2.Code)
}
