package middleware_test

import (
	"digital-innovation/gostrategy/api/middleware"
	"digital-innovation/gostrategy/models"
	"digital-innovation/gostrategy/auth"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func TestRateLimiter(t *testing.T) {
	// Test core logic: 5 req per second, burst of 2
	rl := middleware.NewRateLimiter(rate.Limit(5), 2)
	key := "test-key"

	// First two should pass (burst)
	if !rl.GetLimiter(key).Allow() {
		t.Error("First request should be allowed")
	}
	if !rl.GetLimiter(key).Allow() {
		t.Error("Second request should be allowed")
	}

	// Third should fail
	if rl.GetLimiter(key).Allow() {
		t.Error("Third request should be rate limited")
	}

	// Different key should pass
	if !rl.GetLimiter("other-key").Allow() {
		t.Error("Request for different key should be allowed")
	}
}

func TestIPRateLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// 1 req/s, burst 1
	rl := middleware.NewRateLimiter(rate.Limit(1), 1)
	mw := middleware.IPRateLimitMiddleware(rl)

	t.Run("Allowed", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/", nil)
		c.Request.RemoteAddr = "192.168.1.1:1234"

		mw(c)

		if w.Code == http.StatusTooManyRequests {
			t.Error("Expected request to be allowed, but was rate limited")
		}
	})

	t.Run("Limited", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/", nil)
		c.Request.RemoteAddr = "192.168.1.1:1234"

		// Second request from same IP should fail (burst 1)
		mw(c)

		if w.Code != http.StatusTooManyRequests {
			t.Errorf("Expected status 429, got %d", w.Code)
		}
		if !c.IsAborted() {
			t.Error("Expected context to be aborted")
		}
	})
}

func TestUserRateLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// 1 req/s, burst 1
	rl := middleware.NewRateLimiter(rate.Limit(1), 1)
	mw := middleware.UserRateLimitMiddleware(rl)

	t.Run("Authenticated and Allowed", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/", nil)
		
		// Set user in context
		user := &models.User{ID: 100, Username: "testuser"}
		c.Set(auth.UserContextKey, user)

		mw(c)

		if w.Code == http.StatusTooManyRequests {
			t.Error("Expected request to be allowed")
		}
	})

	t.Run("Authenticated and Limited", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/", nil)
		
		// Same user ID as above
		user := &models.User{ID: 100, Username: "testuser"}
		c.Set(auth.UserContextKey, user)

		// Second request should fail
		mw(c)

		if w.Code != http.StatusTooManyRequests {
			t.Errorf("Expected status 429, got %d", w.Code)
		}
	})

	t.Run("Unauthenticated - Should Skip", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/", nil)
		
		// No user in context
		mw(c)

		if w.Code == http.StatusTooManyRequests {
			t.Error("Unauthenticated requests should not be limited by user-level middleware")
		}
		if c.IsAborted() {
			t.Error("Context should not be aborted for unauthenticated users")
		}
	})
}
