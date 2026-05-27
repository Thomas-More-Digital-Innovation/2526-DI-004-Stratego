// Package api_test provides integration tests for the API.
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
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		var user models.User
		if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
			t.Fatalf("Failed to decode user: %v", err)
		}
		if user.Username != "testuser" {
			t.Errorf("Expected username testuser, got %s", user.Username)
		}

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
		if !hasSession || !hasRefresh {
			t.Errorf("Missing cookies: access_token=%v, refresh_token=%v", hasSession, hasRefresh)
		}
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
		resp, _ := http.Post(ts.URL+"/users/register", contentTypeJSON, bytes.NewBuffer(reqBody))
		if resp.StatusCode != http.StatusConflict {
			t.Errorf("Expected status 409 for duplicate, got %d", resp.StatusCode)
		}
	})

	t.Run("Weak Password", func(t *testing.T) {
		// #nosec G117
		reqBody, _ := json.Marshal(models.CreateUserRequest{
			Username: "weakuser",
			Password: "123",
		})
		resp, _ := http.Post(ts.URL+"/users/register", contentTypeJSON, bytes.NewBuffer(reqBody))
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400 for weak password, got %d", resp.StatusCode)
		}
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
		resp, _ := http.Post(ts.URL+"/users/login", contentTypeJSON, bytes.NewBuffer(loginBody))
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		cookies := resp.Cookies()
		if len(cookies) < 2 {
			t.Errorf("Expected at least 2 cookies, got %d", len(cookies))
		}
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		// #nosec G117
		loginBody, _ := json.Marshal(models.LoginRequest{
			Username: db.TestUserLogin,
			Password: "WrongPassword1",
		})
		resp, _ := http.Post(ts.URL+"/users/login", contentTypeJSON, bytes.NewBuffer(loginBody))
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("Logout", func(t *testing.T) {
		// Login first to get cookies
		// #nosec G117
		loginBody, _ := json.Marshal(models.LoginRequest{
			Username: db.TestUserLogin,
			Password: db.TestPassword,
		})
		resp, _ := http.Post(ts.URL+"/users/login", contentTypeJSON, bytes.NewBuffer(loginBody))
		cookies := resp.Cookies()

		req, _ := http.NewRequest("POST", ts.URL+"/users/logout", nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}

		client := &http.Client{}
		respLogout, _ := client.Do(req)
		if respLogout.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200 for logout, got %d", respLogout.StatusCode)
		}

		// Check if cookies are cleared (MaxAge < 0)
		logoutCookies := respLogout.Cookies()
		foundCleared := false
		for _, c := range logoutCookies {
			if c.Name == accessTokenCookie && c.MaxAge <= 0 {
				foundCleared = true
			}
		}
		if !foundCleared {
			t.Errorf("Access token cookie not cleared")
		}
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
	respReg, _ := http.Post(ts.URL+"/users/register", contentTypeJSON, bytes.NewBuffer(regBody))
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
		resp, _ := client.Do(req)
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("Expected status 201, got %d", resp.StatusCode)
		}
		_ = resp.Body.Close()

		// List setups
		reqList, _ := http.NewRequest("GET", ts.URL+"/board-setups", nil)
		for _, c := range cookies {
			reqList.AddCookie(c)
		}
		respList, _ := client.Do(reqList)
		if respList.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200 for list, got %d", respList.StatusCode)
		}

		var setups []models.BoardSetup
		if err := json.NewDecoder(respList.Body).Decode(&setups); err != nil {
			t.Fatalf("Failed to decode setups: %v", err)
		}
		_ = respList.Body.Close()
		if len(setups) != 1 {
			t.Errorf("Expected 1 setup, got %d", len(setups))
		}
		if setups[0].Name != "My Setup" {
			t.Errorf("Expected name 'My Setup', got '%s'", setups[0].Name)
		}
	})

	t.Run("Board Setup Ownership", func(t *testing.T) {
		// Create another user
		// #nosec G117
		regBody, _ := json.Marshal(models.CreateUserRequest{
			Username: "otheruser",
			Password: db.TestPassword,
		})
		respOther, _ := http.Post(ts.URL+"/users/register", contentTypeJSON, bytes.NewBuffer(regBody))
		otherCookies := respOther.Cookies()
		_ = respOther.Body.Close()

		// Try to list setups for otheruser (should be empty)
		reqList, _ := http.NewRequest("GET", ts.URL+"/board-setups", nil)
		for _, c := range otherCookies {
			reqList.AddCookie(c)
		}
		client := &http.Client{}
		respList, _ := client.Do(reqList)
		var setups []models.BoardSetup
		if err := json.NewDecoder(respList.Body).Decode(&setups); err != nil {
			t.Fatalf("Failed to decode setups: %v", err)
		}
		_ = respList.Body.Close()
		if len(setups) != 0 {
			t.Errorf("Expected 0 setups for otheruser, got %d", len(setups))
		}
	})
}
