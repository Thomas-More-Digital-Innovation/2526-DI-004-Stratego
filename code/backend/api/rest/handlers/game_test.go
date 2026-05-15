package handlers_test

import (
	"digital-innovation/gostrategy/api/core"
	"digital-innovation/gostrategy/api/rest/handlers"
	"digital-innovation/gostrategy/db"
	"digital-innovation/gostrategy/models"
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
