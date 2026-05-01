package auth

import (
	"digital-innovation/stratego/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequireAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Valid Token", func(t *testing.T) {
		userID := 1
		username := "authuser"
		token, _ := GenerateToken(userID, username)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/", nil)
		c.Request.AddCookie(&http.Cookie{Name: "access_token", Value: token})

		handler := RequireAuth()
		handler(c)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", w.Code)
		}
		if c.IsAborted() {
			t.Error("Expected context not to be aborted")
		}

		user := GetCurrentUser(c)
		if user == nil {
			t.Fatal("Expected user in context, got nil")
		}
		if user.ID != userID {
			t.Errorf("Expected userID %d, got %d", userID, user.ID)
		}
		if user.Username != username {
			t.Errorf("Expected username %s, got %s", username, user.Username)
		}
	})

	t.Run("Missing Cookie", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/", nil)

		handler := RequireAuth()
		handler(c)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status Unauthorized, got %d", w.Code)
		}
		if !c.IsAborted() {
			t.Error("Expected context to be aborted")
		}
	})

	t.Run("Invalid Token", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/", nil)
		c.Request.AddCookie(&http.Cookie{Name: "access_token", Value: "invalid-token"})

		handler := RequireAuth()
		handler(c)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status Unauthorized, got %d", w.Code)
		}
		if !c.IsAborted() {
			t.Error("Expected context to be aborted")
		}
	})
}

func TestOptionalAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("With Valid Token", func(t *testing.T) {
		token, _ := GenerateToken(1, "test")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/", nil)
		c.Request.AddCookie(&http.Cookie{Name: "access_token", Value: token})

		handler := OptionalAuth()
		handler(c)

		if GetCurrentUser(c) == nil {
			t.Error("Expected user in context, got nil")
		}
		if c.IsAborted() {
			t.Error("Expected context not to be aborted")
		}
	})

	t.Run("Without Token", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/", nil)

		handler := OptionalAuth()
		handler(c)

		if GetCurrentUser(c) != nil {
			t.Error("Expected no user in context, but found one")
		}
		if c.IsAborted() {
			t.Error("Expected context not to be aborted")
		}
	})
}

func TestGetCurrentUser(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	// Test empty context
	if GetCurrentUser(c) != nil {
		t.Error("Expected nil for empty context")
	}

	// Test with valid user
	testUser := &models.User{ID: 1, Username: "test"}
	c.Set(UserContextKey, testUser)
	if GetCurrentUser(c) != testUser {
		t.Error("Expected to retrieve the same user set in context")
	}

	// Test with invalid type in context
	c.Set(UserContextKey, "not-a-user")
	if GetCurrentUser(c) != nil {
		t.Error("Expected nil for invalid type in context")
	}
}
