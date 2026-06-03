package handlers_test

import (
	"bytes"
	"digital-innovation/gostrategy/internal/api/ws"
	"digital-innovation/gostrategy/internal/models"
	"digital-innovation/gostrategy/internal/testutils"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandleCreateGame(t *testing.T) {
	_, h, server := testutils.SetupHandlerTest(t)

	t.Run("ValidCreation", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		reqBody := map[string]string{
			"gameType": models.HumanVsAi,
			"ai1":      models.Fafo,
		}
		jsonBody, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest("POST", "/games", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")
		h.HandleCreateGame(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("ConflictWithActiveGame", func(t *testing.T) {
		user := &models.User{ID: 1, Username: "alice"}
		reqBody := map[string]string{
			"gameType": models.HumanVsAi,
			"ai1":      models.Fafo,
		}
		jsonBody, _ := json.Marshal(reqBody)

		// First game
		w1 := httptest.NewRecorder()
		c1, _ := gin.CreateTestContext(w1)
		c1.Request, _ = http.NewRequest("POST", "/games", bytes.NewBuffer(jsonBody))
		c1.Request.Header.Set("Content-Type", "application/json")
		c1.Set("user", user)
		h.HandleCreateGame(c1)
		assert.Equal(t, http.StatusOK, w1.Code)

		// Make user active
		var resp struct {
			GameID string `json:"gameId"`
		}
		err := json.Unmarshal(w1.Body.Bytes(), &resp)
		assert.NoError(t, err)
		handler, _ := server.GetSession(resp.GameID)
		client := ws.NewTestClient()
		client.UserID = user.ID
		handler.Hub.Register() <- client
		time.Sleep(50 * time.Millisecond)

		// Second game (conflict)
		w2 := httptest.NewRecorder()
		c2, _ := gin.CreateTestContext(w2)
		c2.Request, _ = http.NewRequest("POST", "/games", bytes.NewBuffer(jsonBody))
		c2.Request.Header.Set("Content-Type", "application/json")
		c2.Set("user", user)
		h.HandleCreateGame(c2)
		assert.Equal(t, http.StatusConflict, w2.Code)
	})
}

func TestHandleListGames(t *testing.T) {
	_, h, server := testutils.SetupHandlerTest(t)
	_, err := server.CreateGame("game-1", models.HumanVsAi, "A", models.Fafo)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	h.HandleListGames(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandleWebSocketConnection_Errors(t *testing.T) {
	_, h, _ := testutils.SetupHandlerTest(t)

	t.Run("GameNotFound", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "gameID", Value: "non-existent"}}
		h.HandleWebSocketConnection(c)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("MissingGameID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		h.HandleWebSocketConnection(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
