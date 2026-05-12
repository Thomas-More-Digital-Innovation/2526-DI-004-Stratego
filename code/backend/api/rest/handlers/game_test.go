package handlers

import (
	"context"
	"digital-innovation/gostrategy/auth"
	"digital-innovation/gostrategy/db"
	"digital-innovation/gostrategy/models"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandleGetGameHistory(t *testing.T) {
	db.SetupTestDB(t)
	r, h, _ := setupTestRouter()
	r.GET("/games/:id/history", h.HandleGetGameHistory)

	gameID := "test-history-game"
	// Save a mock game
	err := db.SaveGame(context.TODO(), gameID, nil, nil, models.HumanVsAi, map[string]any{"board": "setup"}, nil)
	assert.NoError(t, err)

	t.Run("GameFound", func(t *testing.T) {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/games/%s/history", gameID), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), gameID)
	})

	t.Run("GameNotFound", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/games/non-existent/history", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestHandleListMyGames(t *testing.T) {
	db.SetupTestDB(t)
	r, h, _ := setupTestRouter()

	user := &models.User{ID: 100, Username: "me"}
	r.Use(func(c *gin.Context) {
		c.Set(auth.UserContextKey, user)
		c.Next()
	})
	r.GET("/users/me/games", h.HandleListMyGames)

	// Save a game for this user
	userID := 100
	db.SaveGame(context.TODO(), "my-game", &userID, nil, models.HumanVsAi, map[string]any{}, nil)

	t.Run("ListGames", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users/me/games", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "my-game")
		assert.Contains(t, w.Body.String(), "\"total\":1")
	})
}

func TestMonitoringHandlers(t *testing.T) {
	db.SetupTestDB(t)
	r, h, _ := setupTestRouter()
	r.GET("/monitoring/users/count", h.UserCountHandler)
	r.GET("/monitoring/games/count", h.GamesPlayedCountHandler)

	t.Run("UserCount", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/monitoring/users/count", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GamesCount", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/monitoring/games/count", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestStatsHandlers(t *testing.T) {
	db.SetupTestDB(t)
	r, h, _ := setupTestRouter()
	r.GET("/users/:id/stats", h.GetUserStatsHandler)

	t.Run("GetStats", func(t *testing.T) {
		// Mock stats entry
		userID := 1
		db.DB.Create(&models.UserStats{UserID: userID, TotalGames: 10, Wins: 5})

		req, _ := http.NewRequest("GET", "/users/1/stats", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "\"wins\":5")
	})
}
