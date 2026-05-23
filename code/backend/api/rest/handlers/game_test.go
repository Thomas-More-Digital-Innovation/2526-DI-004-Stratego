package handlers_test

import (
	"digital-innovation/gostrategy/api/core"
	"digital-innovation/gostrategy/api/rest/handlers"
	"digital-innovation/gostrategy/db"
	"digital-innovation/gostrategy/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandleGetGameHistory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db.SetupTestDB(t)
	server := core.NewGameServer()
	h := handlers.NewHandler(server)

	// Test non-existent game
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/games/non-existent/history", nil)
	c.Params = []gin.Param{{Key: "id", Value: "non-existent"}}
	h.HandleGetGameHistory(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestHandleListUserGames(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db.SetupTestDB(t)
	server := core.NewGameServer()
	h := handlers.NewHandler(server)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/users/123/games", nil)
	c.Params = []gin.Param{{Key: "id", Value: "123"}}
	h.HandleListUserGames(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandleListMyGames(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db.SetupTestDB(t)
	server := core.NewGameServer()
	h := handlers.NewHandler(server)

	// Unauthenticated
	w1 := httptest.NewRecorder()
	c1, _ := gin.CreateTestContext(w1)
	c1.Request, _ = http.NewRequest("GET", "/users/me/games", nil)
	h.HandleListMyGames(c1)
	assert.Equal(t, http.StatusUnauthorized, w1.Code)

	// Authenticated
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request, _ = http.NewRequest("GET", "/users/me/games", nil)
	c2.Set("user", &models.User{ID: 1, Username: "test"})
	h.HandleListMyGames(c2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestHandleGetReconnectableGame(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db.SetupTestDB(t)
	server := core.NewGameServer()
	h := handlers.NewHandler(server)

	t.Run("Unauthenticated", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/users/me/reconnectable", nil)
		h.HandleGetReconnectableGame(c)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("NoActiveGame", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/users/me/reconnectable", nil)
		c.Set("user", &models.User{ID: 1, Username: "test"})
		h.HandleGetReconnectableGame(c)
		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, false, resp["hasGame"])
	})

	t.Run("ActiveGameExists", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/users/me/reconnectable", nil)
		user := &models.User{ID: 1, Username: "test"}
		c.Set("user", user)

		handler, err := h.CreateGame("reconnect-game", models.HumanVsAi, "test", models.Fafo)
		assert.NoError(t, err)
		handler.Session.SetPlayer1Associate(user.ID, "test")

		h.HandleGetReconnectableGame(c)
		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]any
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, true, resp["hasGame"])
		assert.Equal(t, "reconnect-game", resp["gameId"])
		assert.Equal(t, models.HumanVsAi, resp["gameType"])
	})
}
