package api

import (
	"bytes"
	"context"
	"digital-innovation/gostrategy/db"
	"digital-innovation/gostrategy/models"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
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

	ts := httptest.NewServer(server.router)
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

func TestIntegration_TokenManagement(t *testing.T) {
	_, ts := setupTestServer(t)
	defer ts.Close()

	// 1. Token Refresh
	t.Run("Token Refresh", func(t *testing.T) {
		// #nosec G117
		regBody, _ := json.Marshal(models.CreateUserRequest{
			Username: "refreshuser",
			Password: db.TestPassword,
		})
		respReg, _ := http.Post(ts.URL+"/users/register", contentTypeJSON, bytes.NewBuffer(regBody))
		cookies := respReg.Cookies()
		_ = respReg.Body.Close()

		// Call refresh
		req, _ := http.NewRequest("POST", ts.URL+"/users/refresh", nil)
		for _, c := range cookies {
			if c.Name == "refresh_token" {
				req.AddCookie(c)
			}
		}

		client := &http.Client{}
		respRefresh, _ := client.Do(req)
		if respRefresh.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200 for refresh, got %d", respRefresh.StatusCode)
		}

		// Check if new access_token cookie is set
		newCookies := respRefresh.Cookies()
		found := false
		for _, c := range newCookies {
			if c.Name == accessTokenCookie {
				found = true
			}
		}
		if !found {
			t.Error("New access_token cookie not found after refresh")
		}
		_ = respRefresh.Body.Close()
	})

	// 2. Password Change
	t.Run("Change Password", func(t *testing.T) {
		// #nosec G117
		regBody, _ := json.Marshal(models.CreateUserRequest{
			Username: db.TestUserPassChange,
			Password: db.TestPassword,
		})
		respReg, _ := http.Post(ts.URL+"/users/register", contentTypeJSON, bytes.NewBuffer(regBody))
		cookies := respReg.Cookies()
		_ = respReg.Body.Close()

		changeReq := models.ChangePasswordRequest{
			OldPassword:     db.TestPassword,
			NewPassword:     db.TestPasswordAlt,
			ConfirmPassword: db.TestPasswordAlt,
		}
		reqBody, _ := json.Marshal(changeReq)
		req, _ := http.NewRequest("POST", ts.URL+"/users/me/password", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", contentTypeJSON)
		for _, c := range cookies {
			req.AddCookie(c)
		}

		client := &http.Client{}
		respChange, _ := client.Do(req)
		if respChange.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200 for password change, got %d", respChange.StatusCode)
		}
		_ = respChange.Body.Close()

		// Try to login with OLD password (should fail)
		// #nosec G117
		loginBody, _ := json.Marshal(models.LoginRequest{
			Username: db.TestUserPassChange,
			Password: db.TestPassword,
		})
		respOldLogin, _ := http.Post(ts.URL+"/users/login", contentTypeJSON, bytes.NewBuffer(loginBody))
		if respOldLogin.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401 with old password, got %d", respOldLogin.StatusCode)
		}
		_ = respOldLogin.Body.Close()

		// Try to login with NEW password (should succeed)
		// #nosec G117
		loginBodyNew, _ := json.Marshal(models.LoginRequest{
			Username: db.TestUserPassChange,
			Password: db.TestPasswordAlt,
		})
		respNewLogin, _ := http.Post(ts.URL+"/users/login", contentTypeJSON, bytes.NewBuffer(loginBodyNew))
		if respNewLogin.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200 with new password, got %d", respNewLogin.StatusCode)
		}
		_ = respNewLogin.Body.Close()
	})
}

func TestIntegration_GameManagement(t *testing.T) {
	_, ts := setupTestServer(t)
	defer ts.Close()

	// Pre-register and login
	// #nosec G117
	regBody, _ := json.Marshal(models.CreateUserRequest{
		Username: "gameuser",
		Password: db.TestPassword,
	})
	respReg, _ := http.Post(ts.URL+"/users/register", contentTypeJSON, bytes.NewBuffer(regBody))
	cookies := respReg.Cookies()
	_ = respReg.Body.Close()

	t.Run("Create Game", func(t *testing.T) {
		gameReq := map[string]string{
			"gameType": models.HumanVsAi,
			"ai1":      models.Fafo,
		}
		reqBody, _ := json.Marshal(gameReq)
		req, _ := http.NewRequest("POST", ts.URL+"/games", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", contentTypeJSON)
		for _, c := range cookies {
			req.AddCookie(c)
		}

		client := &http.Client{}
		resp, _ := client.Do(req)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		var respData map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		_ = resp.Body.Close()
		if respData["gameId"] == "" {
			t.Error("Expected gameId in response")
		}
		if !strings.Contains(respData["wsUrl"].(string), "/game/") {
			t.Errorf("Expected valid wsUrl, got %v", respData["wsUrl"])
		}
	})

	t.Run("List Games", func(t *testing.T) {
		resp, _ := http.Get(ts.URL + "/games")
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		var games []models.GameSummary
		if err := json.NewDecoder(resp.Body).Decode(&games); err != nil {
			t.Fatalf("Failed to decode games: %v", err)
		}
		_ = resp.Body.Close()
		// Should have at least one game from previous subtest
		if len(games) < 1 {
			t.Error("Expected at least one game in listing")
		}
	})
}

func TestIntegration_GameHistory(t *testing.T) {
	_, ts := setupTestServer(t)
	defer ts.Close()

	// 1. Create a user and some games
	// #nosec G117
	regBody, _ := json.Marshal(models.CreateUserRequest{
		Username: db.TestUserHistory,
		Password: db.TestPassword,
	})
	respReg, _ := http.Post(ts.URL+"/users/register", contentTypeJSON, bytes.NewBuffer(regBody))
	cookies := respReg.Cookies()

	var user models.User
	if err := json.NewDecoder(respReg.Body).Decode(&user); err != nil {
		t.Fatalf("Failed to decode user: %v", err)
	}
	_ = respReg.Body.Close()

	// Save some fake games directly to DB for testing history
	ctx := context.Background()
	gameID := "hist-1"
	if err := db.SaveGame(ctx, gameID, &user.ID, nil, "ranked", map[string]string{"board": "state"}, nil); err != nil {
		t.Fatalf("Failed to save game: %v", err)
	}
	if err := db.SaveMove(ctx, gameID, models.HistoricalMove{MoveIndex: 1, PlayerID: 0, Result: models.ResultMove}); err != nil {
		t.Fatalf("Failed to save move: %v", err)
	}

	t.Run("Get Single Game History", func(t *testing.T) {
		resp, _ := http.Get(ts.URL + "/games/" + gameID + "/history")
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		var history models.GameHistory
		if err := json.NewDecoder(resp.Body).Decode(&history); err != nil {
			t.Fatalf("Failed to decode history: %v", err)
		}
		_ = resp.Body.Close()
		if history.GameID != gameID {
			t.Errorf("Expected gameId %s, got %s", gameID, history.GameID)
		}
		if len(history.Moves) != 1 {
			t.Errorf("Expected 1 move, got %d", len(history.Moves))
		}
	})

	t.Run("List My Games (Paged)", func(t *testing.T) {
		req, _ := http.NewRequest("GET", ts.URL+"/users/me/games?limit=5&offset=0", nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}

		client := &http.Client{}
		resp, _ := client.Do(req)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode result: %v", err)
		}
		_ = resp.Body.Close()

		games := result["games"].([]interface{})
		if len(games) != 1 {
			t.Errorf("Expected 1 game, got %d", len(games))
		}
		if int(result["total"].(float64)) != 1 {
			t.Errorf("Expected total 1, got %v", result["total"])
		}
	})

	t.Run("List User Games (Public)", func(t *testing.T) {
		resp, _ := http.Get(fmt.Sprintf("%s/users/%d/games", ts.URL, user.ID))
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode result: %v", err)
		}
		_ = resp.Body.Close()
		if int(result["total"].(float64)) != 1 {
			t.Errorf("Expected total 1, got %v", result["total"])
		}
	})
}
