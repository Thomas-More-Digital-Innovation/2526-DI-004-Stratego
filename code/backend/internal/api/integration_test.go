// Package api provides integration tests for the API.
package api

import (
	"bytes"
	"digital-innovation/gostrategy/internal/db"
	"digital-innovation/gostrategy/internal/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const accessTokenCookie = "access_token"
const refreshTokenCookie = "refresh_token"
const contentTypeJSON = "application/json"

func setupTestServer(t *testing.T) (*GameServer, *httptest.Server) {
	gin.SetMode(gin.TestMode)
	db.SetupTestDB(t)

	server := NewGameServer()
	server.SetupRoutes()

	ts := httptest.NewServer(server.Router)
	return server, ts
}

func TestIntegration_UserRegistration(t *testing.T) {
	_, ts := setupTestServer(t)
	defer ts.Close()

	t.Run("Successful Registration", func(t *testing.T) {
		// #nosec G117
		reqBody, _ := json.Marshal(models.CreateUserRequest{
			Username: db.TestUser,
			Password: db.TestPassword,
		})
		resp, err := http.Post(ts.URL+"/users/register", contentTypeJSON, bytes.NewBuffer(reqBody))
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		var user models.User
		err = json.NewDecoder(resp.Body).Decode(&user)
		require.NoError(t, err)
		assert.Equal(t, "testuser", user.Username)

		// Check cookies
		cookies := resp.Cookies()
		hasSession := false
		hasRefresh := false
		for _, c := range cookies {
			if c.Name == accessTokenCookie {
				hasSession = true
			}
			if c.Name == refreshTokenCookie {
				hasRefresh = true
			}
		}
		assert.True(t, hasSession)
		assert.True(t, hasRefresh)
	})

	t.Run("Duplicate Username", func(t *testing.T) {
		// #nosec G117
		reqBody, _ := json.Marshal(models.CreateUserRequest{
			Username: "testuser_dup",
			Password: db.TestPassword,
		})
		// First one
		resp1, err := http.Post(ts.URL+"/users/register", contentTypeJSON, bytes.NewBuffer(reqBody))
		if err == nil {
			_ = resp1.Body.Close()
		}

		// Duplicate
		resp, err := http.Post(ts.URL+"/users/register", contentTypeJSON, bytes.NewBuffer(reqBody))
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	t.Run("Weak Password", func(t *testing.T) {
		// #nosec G117
		reqBody, _ := json.Marshal(models.CreateUserRequest{
			Username: "weakuser",
			Password: "123",
		})
		resp, err := http.Post(ts.URL+"/users/register", contentTypeJSON, bytes.NewBuffer(reqBody))
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestIntegration_LoginLogout(t *testing.T) {
	_, ts := setupTestServer(t)
	defer ts.Close()

	// Pre-register
	// #nosec G117
	regBody, _ := json.Marshal(models.CreateUserRequest{
		Username: db.TestUserLogin,
		Password: db.TestPassword,
	})
	respReg, err := http.Post(ts.URL+"/users/register", contentTypeJSON, bytes.NewBuffer(regBody))
	if err == nil {
		_ = respReg.Body.Close()
	}

	t.Run("Successful Login", func(t *testing.T) {
		// #nosec G117
		loginBody, _ := json.Marshal(models.LoginRequest{
			Username: db.TestUserLogin,
			Password: db.TestPassword,
		})
		resp, err := http.Post(ts.URL+"/users/login", contentTypeJSON, bytes.NewBuffer(loginBody))
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		cookies := resp.Cookies()
		assert.True(t, len(cookies) >= 2)
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		// #nosec G117
		loginBody, _ := json.Marshal(models.LoginRequest{
			Username: db.TestUserLogin,
			Password: "WrongPassword1",
		})
		resp, err := http.Post(ts.URL+"/users/login", contentTypeJSON, bytes.NewBuffer(loginBody))
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Logout", func(t *testing.T) {
		// Login first to get cookies
		// #nosec G117
		loginBody, _ := json.Marshal(models.LoginRequest{
			Username: db.TestUserLogin,
			Password: db.TestPassword,
		})
		resp, err := http.Post(ts.URL+"/users/login", contentTypeJSON, bytes.NewBuffer(loginBody))
		require.NoError(t, err)
		cookies := resp.Cookies()
		_ = resp.Body.Close()

		req, _ := http.NewRequest("POST", ts.URL+"/users/logout", nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}

		client := &http.Client{}
		respLogout, err := client.Do(req)
		require.NoError(t, err)
		defer func() { _ = respLogout.Body.Close() }()
		assert.Equal(t, http.StatusOK, respLogout.StatusCode)

		// Check if cookies are cleared (MaxAge < 0)
		logoutCookies := respLogout.Cookies()
		foundCleared := false
		for _, c := range logoutCookies {
			if c.Name == accessTokenCookie && c.MaxAge <= 0 {
				foundCleared = true
			}
		}
		assert.True(t, foundCleared)
	})
}

func TestIntegration_BoardSetups(t *testing.T) {
	_, ts := setupTestServer(t)
	defer ts.Close()

	// Pre-register and login
	// #nosec G117
	regBody, _ := json.Marshal(models.CreateUserRequest{
		Username: "boarduser",
		Password: db.TestPassword,
	})
	respReg, err := http.Post(ts.URL+"/users/register", contentTypeJSON, bytes.NewBuffer(regBody))
	require.NoError(t, err)
	cookies := respReg.Cookies()
	_ = respReg.Body.Close()

	t.Run("Create and List Board Setup", func(t *testing.T) {
		setupReq := models.CreateBoardSetupRequest{
			Name:      "My Setup",
			SetupData: `[{"type": "Marshall", "rank": "10", "ownerId": 0}]`,
		}
		reqBody, _ := json.Marshal(setupReq)
		req, _ := http.NewRequest("POST", ts.URL+"/board-setups", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", contentTypeJSON)
		for _, c := range cookies {
			req.AddCookie(c)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		_ = resp.Body.Close()

		// List setups
		reqList, _ := http.NewRequest("GET", ts.URL+"/board-setups", nil)
		for _, c := range cookies {
			reqList.AddCookie(c)
		}
		respList, err := client.Do(reqList)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, respList.StatusCode)

		var setups []models.BoardSetup
		err = json.NewDecoder(respList.Body).Decode(&setups)
		require.NoError(t, err)
		_ = respList.Body.Close()

		require.Len(t, setups, 1)
		assert.Equal(t, "My Setup", setups[0].Name)
	})

	t.Run("Board Setup Ownership", func(t *testing.T) {
		// Create another user
		// #nosec G117
		regBody, _ := json.Marshal(models.CreateUserRequest{
			Username: "otheruser",
			Password: db.TestPassword,
		})
		respOther, err := http.Post(ts.URL+"/users/register", contentTypeJSON, bytes.NewBuffer(regBody))
		require.NoError(t, err)
		otherCookies := respOther.Cookies()
		_ = respOther.Body.Close()

		// Try to list setups for otheruser (should be empty)
		reqList, _ := http.NewRequest("GET", ts.URL+"/board-setups", nil)
		for _, c := range otherCookies {
			reqList.AddCookie(c)
		}
		client := &http.Client{}
		respList, err := client.Do(reqList)
		require.NoError(t, err)
		defer func() { _ = respList.Body.Close() }()

		var setups []models.BoardSetup
		err = json.NewDecoder(respList.Body).Decode(&setups)
		require.NoError(t, err)
		assert.Empty(t, setups)
	})
}
