package handlers_test

import (
	"bytes"
	"context"
	"digital-innovation/gostrategy/internal/db"
	"digital-innovation/gostrategy/internal/models"
	"digital-innovation/gostrategy/internal/testutils"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const (
	oldPassword     = "OldPassword123"
	newPassword     = "NewPassword123"
	weakPassword    = "weak"
	invalidPassword = "Invalid<Password"
)

func TestUserHandlers(t *testing.T) {
	_, h, _ := testutils.SetupHandlerTest(t)

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

	t.Run("ChangePassword - Passwords Mismatch", func(t *testing.T) {
		u, _ := db.CreateUser(context.Background(), "user_mismatch", oldPassword, "")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		changeReq := models.ChangePasswordRequest{
			OldPassword:     oldPassword,
			NewPassword:     newPassword,
			ConfirmPassword: "DifferentPassword123",
		}
		jsonBody, _ := json.Marshal(changeReq)
		c.Request, _ = http.NewRequest("POST", "/users/me/password", bytes.NewBuffer(jsonBody))
		c.Set("user", u)
		h.ChangePasswordHandler(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ChangePassword - Invalid JSON Body", func(t *testing.T) {
		u, _ := db.CreateUser(context.Background(), "user_invalid_body", oldPassword, "")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/users/me/password", bytes.NewBufferString("invalid-json"))
		c.Set("user", u)
		h.ChangePasswordHandler(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ChangePassword - Invalid Characters in Password", func(t *testing.T) {
		u, _ := db.CreateUser(context.Background(), "user_invalid_chars", oldPassword, "")

		// Case 1: Invalid char in old password
		w1 := httptest.NewRecorder()
		c1, _ := gin.CreateTestContext(w1)
		changeReq1 := models.ChangePasswordRequest{
			OldPassword:     invalidPassword,
			NewPassword:     newPassword,
			ConfirmPassword: newPassword,
		}
		jsonBody1, _ := json.Marshal(changeReq1)
		c1.Request, _ = http.NewRequest("POST", "/users/me/password", bytes.NewBuffer(jsonBody1))
		c1.Set("user", u)
		h.ChangePasswordHandler(c1)
		assert.Equal(t, http.StatusBadRequest, w1.Code)

		// Case 2: Invalid char in new password
		w2 := httptest.NewRecorder()
		c2, _ := gin.CreateTestContext(w2)
		changeReq2 := models.ChangePasswordRequest{
			OldPassword:     oldPassword,
			NewPassword:     invalidPassword,
			ConfirmPassword: invalidPassword,
		}
		jsonBody2, _ := json.Marshal(changeReq2)
		c2.Request, _ = http.NewRequest("POST", "/users/me/password", bytes.NewBuffer(jsonBody2))
		c2.Set("user", u)
		h.ChangePasswordHandler(c2)
		assert.Equal(t, http.StatusBadRequest, w2.Code)
	})
}
