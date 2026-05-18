package handlers_test

import (
	"bytes"
	"context"
	"digital-innovation/gostrategy/api/core"
	"digital-innovation/gostrategy/api/rest/handlers"
	"digital-innovation/gostrategy/db"
	"digital-innovation/gostrategy/models"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestBoardSetupHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db.SetupTestDB(t)
	server := core.NewGameServer()
	h := handlers.NewHandler(server)

	user, _ := db.CreateUser(context.Background(), "alice", "StrongPassword1", "")
	var createdSetup models.BoardSetup

	t.Run("Create", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		setupReq := models.BoardSetup{
			Name:        "Test Setup",
			Description: "A test description",
			SetupData:   "some-binary-data",
		}
		jsonBody, _ := json.Marshal(setupReq)
		c.Request, _ = http.NewRequest("POST", "/users/me/board-setups", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user", user)
		h.CreateBoardSetupHandler(c)
		assert.Equal(t, http.StatusCreated, w.Code)

		err := json.Unmarshal(w.Body.Bytes(), &createdSetup)
		assert.NoError(t, err)
	})

	t.Run("List", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/users/me/board-setups", nil)
		c.Set("user", user)
		h.GetUserBoardSetupsHandler(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Get", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", fmt.Sprintf("/users/me/board-setups/%d", createdSetup.ID), nil)
		c.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%d", createdSetup.ID)}}
		c.Set("user", user)
		h.GetBoardSetupHandler(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Update", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		setupReq := models.BoardSetup{
			Name: "Updated Name",
		}
		jsonBodyUpdate, _ := json.Marshal(setupReq)
		c.Request, _ = http.NewRequest("PUT", fmt.Sprintf("/users/me/board-setups/%d", createdSetup.ID), bytes.NewBuffer(jsonBodyUpdate))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%d", createdSetup.ID)}}
		c.Set("user", user)
		h.UpdateBoardSetupHandler(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Delete", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("DELETE", fmt.Sprintf("/users/me/board-setups/%d", createdSetup.ID), nil)
		c.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%d", createdSetup.ID)}}
		c.Set("user", user)
		h.DeleteBoardSetupHandler(c)
		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}
