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

func TestUserHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db.SetupTestDB(t)
	server := core.NewGameServer()
	h := handlers.NewHandler(server)

	user, _ := db.CreateUser(context.Background(), "alice", "StrongPassword1", "")

	t.Run("GetCurrentUser", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/users/me", nil)
		c.Set("user", user)
		h.GetCurrentUserHandler(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetUserByID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", fmt.Sprintf("/users/%d", user.ID), nil)
		c.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%d", user.ID)}}
		h.GetUserHandler(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("ChangePassword", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		changeReq := models.ChangePasswordRequest{
			OldPassword:     "StrongPassword1",
			NewPassword:     "NewStrong123",
			ConfirmPassword: "NewStrong123",
		}
		jsonBody, _ := json.Marshal(changeReq)
		c.Request, _ = http.NewRequest("POST", "/users/me/password", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user", user)
		h.ChangePasswordHandler(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
