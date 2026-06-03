// Package testutils provides utility functions for testing the Stratego game server.
package testutils

import (
	"bytes"
	"digital-innovation/gostrategy/internal/api/core"
	"digital-innovation/gostrategy/internal/api/rest/handlers"
	"digital-innovation/gostrategy/internal/db"
	"digital-innovation/gostrategy/pkg/game"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	Player1Name = "Piet"
	Player2Name = "Bob"
)

// SetupHandlerTest registers a cleanup hook to guarantee gameserver terminates
func SetupHandlerTest(t *testing.T) (*gin.Engine, *handlers.Handler, *core.GameServer) {
	gin.SetMode(gin.TestMode)
	db.SetupTestDB(t)
	server := core.NewGameServer()
	t.Cleanup(server.Stop)
	h := handlers.NewHandler(server)
	r := gin.New()
	return r, h, server
}

// PerformRequest parses body to json and records the response
func PerformRequest(r *gin.Engine, method, path string, body any, headers map[string]string) *httptest.ResponseRecorder {
	var bodyReader *bytes.Reader
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		bodyReader = bytes.NewReader(jsonBody)
	} else {
		bodyReader = bytes.NewReader([]byte{})
	}

	req, _ := http.NewRequest(method, path, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// SetupDBTest returns test database instance and restores old connection pointer afterwards
func SetupDBTest(t *testing.T) *gorm.DB {
	testDB := db.SetupTestDB(t)
	oldDB := db.DB
	db.DB = testDB
	t.Cleanup(func() {
		db.DB = oldDB
	})
	return testDB
}

// SetupTestGame initializes the game with players and human controllers
func SetupTestGame() (*game.Game, *game.Player, *game.Player) {
	player1 := game.NewPlayer(0, Player1Name, "red")
	controller1 := game.NewHumanPlayerController(&player1)
	player2 := game.NewPlayer(1, Player2Name, "blue")
	controller2 := game.NewHumanPlayerController(&player2)
	g := game.NewGame(controller1, controller2)
	return g, &player1, &player2
}
