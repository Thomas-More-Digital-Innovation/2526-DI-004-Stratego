package api

import (
	"bytes"
	"context"
	"digital-innovation/gostrategy/internal/db"
	"digital-innovation/gostrategy/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

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
