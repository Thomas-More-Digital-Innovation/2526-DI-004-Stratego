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
		Password: "StrongPassword1",
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
		Password: "StrongPassword1",
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
		Password: "StrongPassword1",
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
	user, _ := db.CreateUser(context.Background(), "alice", "StrongPassword1", "")
	refreshToken := "some-refresh-token"
	expiresAt := time.Now().Add(24 * time.Hour)
	err := db.SaveRefreshToken(context.Background(), user.ID, refreshToken, expiresAt)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/refresh", nil)
	c.Request.AddCookie(&http.Cookie{Name: "refresh_token", Value: refreshToken})

	h.RefreshTokenHandler(c)
	assert.Equal(t, http.StatusOK, w.Code)
}
