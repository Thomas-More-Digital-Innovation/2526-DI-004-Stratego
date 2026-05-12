package handlers

import (
	"bytes"
	"digital-innovation/gostrategy/api/core"
	"digital-innovation/gostrategy/api/ws"
	"digital-innovation/gostrategy/auth"
	"digital-innovation/gostrategy/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() (*gin.Engine, *Handler, *core.GameServer) {
	gin.SetMode(gin.TestMode)
	s := core.NewGameServer()
	h := NewHandler(s)
	r := gin.New()
	return r, h, s
}

func TestHandleCreateGame(t *testing.T) {
	r, h, _ := setupTestRouter()
	r.POST("/games", h.HandleCreateGame)

	t.Run("GuestCreatesGame", func(t *testing.T) {
		body := map[string]string{
			"gameType": models.HumanVsAi,
			"ai1":      models.Fafo,
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/games", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Contains(t, resp, "gameId")
	})

	t.Run("UserCreatesGame", func(t *testing.T) {
		user := &models.User{ID: 1, Username: "testuser"}

		rWithAuth := gin.New()
		rWithAuth.Use(func(c *gin.Context) {
			c.Set(auth.UserContextKey, user)
			c.Next()
		})
		rWithAuth.POST("/games", h.HandleCreateGame)

		body := map[string]string{
			"gameType": models.HumanVsAi,
			"ai1":      models.Fafo,
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/games", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()
		rWithAuth.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("UserInActiveGameBlocked", func(t *testing.T) {
		user := &models.User{ID: 2, Username: "busyuser"}

		// Setup existing active game
		handler, _ := h.CreateGame("busy-game", models.HumanVsAi, "busyuser", models.Fafo)
		handler.Session.Player1UserID = &user.ID

		// Register a client to make it "active"
		client := ws.NewTestClient()
		client.UserID = user.ID
		handler.Hub.Register() <- client
		// Give it a tiny bit of time to register
		time.Sleep(10 * time.Millisecond)

		rWithAuth := gin.New()
		rWithAuth.Use(func(c *gin.Context) {
			c.Set(auth.UserContextKey, user)
			c.Next()
		})
		rWithAuth.POST("/games", h.HandleCreateGame)

		body := map[string]string{
			"gameType": models.HumanVsAi,
			"ai1":      models.Fafo,
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/games", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()
		rWithAuth.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.Contains(t, w.Body.String(), "already in an active game")
	})
	t.Run("UserInStaleGameAllowed", func(t *testing.T) {
		user := &models.User{ID: 3, Username: "staleuser"}

		// Setup existing game that is "stale" (no one connected)
		handler, _ := h.CreateGame("stale-game", models.HumanVsAi, "staleuser", models.Fafo)
		handler.Session.Player1UserID = &user.ID

		rWithAuth := gin.New()
		rWithAuth.Use(func(c *gin.Context) {
			c.Set(auth.UserContextKey, user)
			c.Next()
		})
		rWithAuth.POST("/games", h.HandleCreateGame)

		body := map[string]string{
			"gameType": models.HumanVsAi,
			"ai1":      models.Fafo,
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/games", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()
		rWithAuth.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		// Stale game should have been removed
		_, exists := h.GetSession("stale-game")
		assert.False(t, exists)
	})
}

func TestAssociatePlayer(t *testing.T) {
	_, h, _ := setupTestRouter()

	t.Run("AllowUnassociatedUser", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		user := &models.User{ID: 10}
		userID := 10
		allowed := h.associatePlayer(c, &userID, nil, user, "some-game")
		assert.True(t, allowed)
	})

	t.Run("BlockWrongUser", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		currentUserID := 10
		associatedUserID := 20
		allowed := h.associatePlayer(c, &currentUserID, &associatedUserID, nil, "some-game")
		assert.False(t, allowed)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}
