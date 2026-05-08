package api

import (
	"bytes"
	"digital-innovation/gostrategy/auth"
	"digital-innovation/gostrategy/db"
	"digital-innovation/gostrategy/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupTestServer(t *testing.T) (*GameServer, *httptest.Server) {
	gin.SetMode(gin.TestMode)
	db.SetupTestDB(t)

	server := NewGameServer()
	// Initialize routes (normally done in StartServer, but we want to avoid actual listening)
	server.router.Use(gin.Recovery())
	
	// User & Auth
	users := server.router.Group("/users")
	{
		users.POST("/register", server.RegisterUserHandler)
		users.POST("/login", server.LoginHandler)
		users.POST("/logout", server.LogoutHandler)
		users.POST("/refresh", server.RefreshTokenHandler)

		// Authenticated user routes
		me := users.Group("/me")
		me.Use(auth.RequireAuth())
		{
			me.POST("/password", server.ChangePasswordHandler)
		}
	}

	// Board Setups (Require Auth)
	setups := server.router.Group("/board-setups")
	setups.Use(auth.RequireAuth())
	{
		setups.POST("", server.CreateBoardSetupHandler)
		setups.GET("", server.GetUserBoardSetupsHandler)
		setups.GET("/:id", server.GetBoardSetupHandler)
		setups.PUT("/:id", server.UpdateBoardSetupHandler)
		setups.DELETE("/:id", server.DeleteBoardSetupHandler)
	}

	// Game Endpoints (Optional Auth)
	games := server.router.Group("/games")
	games.Use(auth.OptionalAuth())
	{
		games.POST("", server.HandleCreateGame)
		games.GET("", server.HandleListGames)
	}

	ts := httptest.NewServer(server.router)
	return server, ts
}

func TestIntegration_UserRegistration(t *testing.T) {
	_, ts := setupTestServer(t)
	defer ts.Close()

	t.Run("Successful Registration", func(t *testing.T) {
		reqBody, _ := json.Marshal(models.CreateUserRequest{
			Username: "testuser",
			Password: "StrongPassword1",
		})
		resp, err := http.Post(ts.URL+"/users/register", "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", resp.StatusCode)
		}

		var user models.User
		json.NewDecoder(resp.Body).Decode(&user)
		if user.Username != "testuser" {
			t.Errorf("Expected username testuser, got %s", user.Username)
		}

		// Check cookies
		cookies := resp.Cookies()
		hasSession := false
		hasRefresh := false
		for _, c := range cookies {
			if c.Name == "access_token" {
				hasSession = true
			}
			if c.Name == "refresh_token" {
				hasRefresh = true
			}
		}
		if !hasSession || !hasRefresh {
			t.Errorf("Missing cookies: access_token=%v, refresh_token=%v", hasSession, hasRefresh)
		}
	})

	t.Run("Duplicate Username", func(t *testing.T) {
		reqBody, _ := json.Marshal(models.CreateUserRequest{
			Username: "testuser_dup",
			Password: "StrongPassword1",
		})
		// First one
		http.Post(ts.URL+"/users/register", "application/json", bytes.NewBuffer(reqBody))
		
		// Duplicate
		resp, _ := http.Post(ts.URL+"/users/register", "application/json", bytes.NewBuffer(reqBody))
		if resp.StatusCode != http.StatusConflict {
			t.Errorf("Expected status 409 for duplicate, got %d", resp.StatusCode)
		}
	})

	t.Run("Weak Password", func(t *testing.T) {
		reqBody, _ := json.Marshal(models.CreateUserRequest{
			Username: "weakuser",
			Password: "123",
		})
		resp, _ := http.Post(ts.URL+"/users/register", "application/json", bytes.NewBuffer(reqBody))
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400 for weak password, got %d", resp.StatusCode)
		}
	})
}

func TestIntegration_LoginLogout(t *testing.T) {
	_, ts := setupTestServer(t)
	defer ts.Close()

	// Pre-register
	regBody, _ := json.Marshal(models.CreateUserRequest{
		Username: "loginuser",
		Password: "StrongPassword1",
	})
	http.Post(ts.URL+"/users/register", "application/json", bytes.NewBuffer(regBody))

	t.Run("Successful Login", func(t *testing.T) {
		loginBody, _ := json.Marshal(models.LoginRequest{
			Username: "loginuser",
			Password: "StrongPassword1",
		})
		resp, _ := http.Post(ts.URL+"/users/login", "application/json", bytes.NewBuffer(loginBody))
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		cookies := resp.Cookies()
		if len(cookies) < 2 {
			t.Errorf("Expected at least 2 cookies, got %d", len(cookies))
		}
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		loginBody, _ := json.Marshal(models.LoginRequest{
			Username: "loginuser",
			Password: "WrongPassword1",
		})
		resp, _ := http.Post(ts.URL+"/users/login", "application/json", bytes.NewBuffer(loginBody))
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("Logout", func(t *testing.T) {
		// Login first to get cookies
		loginBody, _ := json.Marshal(models.LoginRequest{
			Username: "loginuser",
			Password: "StrongPassword1",
		})
		resp, _ := http.Post(ts.URL+"/users/login", "application/json", bytes.NewBuffer(loginBody))
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
			if c.Name == "access_token" && c.MaxAge <= 0 {
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
	regBody, _ := json.Marshal(models.CreateUserRequest{
		Username: "boarduser",
		Password: "StrongPassword1",
	})
	respReg, _ := http.Post(ts.URL+"/users/register", "application/json", bytes.NewBuffer(regBody))
	cookies := respReg.Cookies()

	t.Run("Create and List Board Setup", func(t *testing.T) {
		setupReq := models.CreateBoardSetupRequest{
			Name:      "My Setup",
			SetupData: `[{"type": "Marshall", "rank": "10", "ownerId": 0}]`,
		}
		reqBody, _ := json.Marshal(setupReq)
		req, _ := http.NewRequest("POST", ts.URL+"/board-setups", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		for _, c := range cookies {
			req.AddCookie(c)
		}

		client := &http.Client{}
		resp, _ := client.Do(req)
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("Expected status 201, got %d", resp.StatusCode)
		}

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
		json.NewDecoder(respList.Body).Decode(&setups)
		if len(setups) != 1 {
			t.Errorf("Expected 1 setup, got %d", len(setups))
		}
		if setups[0].Name != "My Setup" {
			t.Errorf("Expected name 'My Setup', got '%s'", setups[0].Name)
		}
	})

	t.Run("Board Setup Ownership", func(t *testing.T) {
		// Create another user
		regBody, _ := json.Marshal(models.CreateUserRequest{
			Username: "otheruser",
			Password: "StrongPassword1",
		})
		respOther, _ := http.Post(ts.URL+"/users/register", "application/json", bytes.NewBuffer(regBody))
		otherCookies := respOther.Cookies()

		// Try to list setups for otheruser (should be empty)
		reqList, _ := http.NewRequest("GET", ts.URL+"/board-setups", nil)
		for _, c := range otherCookies {
			reqList.AddCookie(c)
		}
		client := &http.Client{}
		respList, _ := client.Do(reqList)
		var setups []models.BoardSetup
		json.NewDecoder(respList.Body).Decode(&setups)
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
		regBody, _ := json.Marshal(models.CreateUserRequest{
			Username: "refreshuser",
			Password: "StrongPassword1",
		})
		respReg, _ := http.Post(ts.URL+"/users/register", "application/json", bytes.NewBuffer(regBody))
		cookies := respReg.Cookies()

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
			if c.Name == "access_token" {
				found = true
			}
		}
		if !found {
			t.Error("New access_token cookie not found after refresh")
		}
	})

	// 2. Password Change
	t.Run("Change Password", func(t *testing.T) {
		regBody, _ := json.Marshal(models.CreateUserRequest{
			Username: "passuser",
			Password: "StrongPassword1",
		})
		respReg, _ := http.Post(ts.URL+"/users/register", "application/json", bytes.NewBuffer(regBody))
		cookies := respReg.Cookies()

		changeReq := models.ChangePasswordRequest{
			OldPassword:     "StrongPassword1",
			NewPassword:     "NewStrongPassword2",
			ConfirmPassword: "NewStrongPassword2",
		}
		reqBody, _ := json.Marshal(changeReq)
		req, _ := http.NewRequest("POST", ts.URL+"/users/me/password", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		for _, c := range cookies {
			req.AddCookie(c)
		}

		client := &http.Client{}
		respChange, _ := client.Do(req)
		if respChange.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200 for password change, got %d", respChange.StatusCode)
		}

		// Try to login with OLD password (should fail)
		loginBody, _ := json.Marshal(models.LoginRequest{
			Username: "passuser",
			Password: "StrongPassword1",
		})
		respOldLogin, _ := http.Post(ts.URL+"/users/login", "application/json", bytes.NewBuffer(loginBody))
		if respOldLogin.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401 with old password, got %d", respOldLogin.StatusCode)
		}

		// Try to login with NEW password (should succeed)
		loginBodyNew, _ := json.Marshal(models.LoginRequest{
			Username: "passuser",
			Password: "NewStrongPassword2",
		})
		respNewLogin, _ := http.Post(ts.URL+"/users/login", "application/json", bytes.NewBuffer(loginBodyNew))
		if respNewLogin.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200 with new password, got %d", respNewLogin.StatusCode)
		}
	})
}

func TestIntegration_GameManagement(t *testing.T) {
	_, ts := setupTestServer(t)
	defer ts.Close()

	// Pre-register and login
	regBody, _ := json.Marshal(models.CreateUserRequest{
		Username: "gameuser",
		Password: "StrongPassword1",
	})
	respReg, _ := http.Post(ts.URL+"/users/register", "application/json", bytes.NewBuffer(regBody))
	cookies := respReg.Cookies()

	t.Run("Create Game", func(t *testing.T) {
		gameReq := map[string]string{
			"gameType": models.HumanVsAi,
			"ai1":      models.Fafo,
		}
		reqBody, _ := json.Marshal(gameReq)
		req, _ := http.NewRequest("POST", ts.URL+"/games", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		for _, c := range cookies {
			req.AddCookie(c)
		}

		client := &http.Client{}
		resp, _ := client.Do(req)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		var respData map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&respData)
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
		json.NewDecoder(resp.Body).Decode(&games)
		// Should have at least one game from previous subtest
		if len(games) < 1 {
			t.Error("Expected at least one game in listing")
		}
	})
}
