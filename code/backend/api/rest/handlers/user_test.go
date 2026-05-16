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

const (
	oldPassword  = "OldPassword123"
	newPassword  = "NewPassword123"
	weakPassword = "weak"
)

func TestUserHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db.SetupTestDB(t)
	server := core.NewGameServer()
	h := handlers.NewHandler(server)

	user, _ := db.CreateUser(context.Background(), "alice", oldPassword, "")

	t.Run("GetCurrentUser", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/users/me", nil)
		c.Set("user", user)
		h.GetCurrentUserHandler(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetCurrentUser - Unauthorized", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/users/me", nil)
		h.GetCurrentUserHandler(c)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("GetUserByID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", fmt.Sprintf("/users/%d", user.ID), nil)
		c.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%d", user.ID)}}
		h.GetUserHandler(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetUserByID - Invalid ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/users/abc", nil)
		c.Params = []gin.Param{{Key: "id", Value: "abc"}}
		h.GetUserHandler(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GetUserByID - Not Found", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/users/9999", nil)
		c.Params = []gin.Param{{Key: "id", Value: "9999"}}
		h.GetUserHandler(c)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("ChangePassword - Success", func(t *testing.T) {
		u, _ := db.CreateUser(context.Background(), "user_success", oldPassword, "")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		changeReq := models.ChangePasswordRequest{
			OldPassword:     oldPassword,
			NewPassword:     newPassword,
			ConfirmPassword: newPassword,
		}
		jsonBody, _ := json.Marshal(changeReq)
		c.Request, _ = http.NewRequest("POST", "/users/me/password", bytes.NewBuffer(jsonBody))
		c.Set("user", u)
		h.ChangePasswordHandler(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("ChangePassword - Wrong Current Password", func(t *testing.T) {
		u, _ := db.CreateUser(context.Background(), "user_wrong_pass", "SomeOtherPassword123", "")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		changeReq := models.ChangePasswordRequest{
			OldPassword:     oldPassword,
			NewPassword:     newPassword,
			ConfirmPassword: newPassword,
		}
		jsonBody, _ := json.Marshal(changeReq)
		c.Request, _ = http.NewRequest("POST", "/users/me/password", bytes.NewBuffer(jsonBody))
		c.Set("user", u)
		h.ChangePasswordHandler(c)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("ChangePassword - Weak New Password", func(t *testing.T) {
		u, _ := db.CreateUser(context.Background(), "user_weak_pass", oldPassword, "")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		changeReq := models.ChangePasswordRequest{
			OldPassword:     oldPassword,
			NewPassword:     weakPassword,
			ConfirmPassword: weakPassword,
		}
		jsonBody, _ := json.Marshal(changeReq)
		c.Request, _ = http.NewRequest("POST", "/users/me/password", bytes.NewBuffer(jsonBody))
		c.Set("user", u)
		h.ChangePasswordHandler(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
