package middleware_test

import (
	"digital-innovation/gostrategy/internal/api/middleware"
	"digital-innovation/gostrategy/internal/auth"
	"digital-innovation/gostrategy/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

func TestRateLimiter(t *testing.T) {
	// Test core logic: 5 req per second, burst of 2
	rl := middleware.NewRateLimiter(rate.Limit(5), 2)
	defer rl.Stop()
	key := "test-key"

	t.Run("Burst Allowance", func(t *testing.T) {
		// First two should pass (burst)
		assert.True(t, rl.GetLimiter(key).Allow(), "First request should be allowed")
		assert.True(t, rl.GetLimiter(key).Allow(), "Second request should be allowed")
	})

	t.Run("Rate Limit Block", func(t *testing.T) {
		// Third should fail
		assert.False(t, rl.GetLimiter(key).Allow(), "Third request should be rate limited")
	})

	t.Run("Key Isolation", func(t *testing.T) {
		// Different key should pass
		assert.True(t, rl.GetLimiter("other-key").Allow(), "Request for different key should be allowed")
	})
}

func TestIPRateLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 1 req/s, burst 1
	rl := middleware.NewRateLimiter(rate.Limit(1), 1)
	defer rl.Stop()
	mw := middleware.IPRateLimitMiddleware(rl)

	t.Run("Allowed", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/", nil)
		c.Request.RemoteAddr = "192.168.1.1:1234"

		mw(c)

		assert.NotEqual(t, http.StatusTooManyRequests, w.Code, "Expected request to be allowed, but was rate limited")
	})

	t.Run("Limited", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/", nil)
		c.Request.RemoteAddr = "192.168.1.1:1234"

		// Second request from same IP should fail (burst 1)
		mw(c)

		assert.Equal(t, http.StatusTooManyRequests, w.Code, "Expected status 429")
		assert.True(t, c.IsAborted(), "Expected context to be aborted")
	})
}

func TestUserRateLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 1 req/s, burst 1
	rl := middleware.NewRateLimiter(rate.Limit(1), 1)
	defer rl.Stop()
	mw := middleware.UserRateLimitMiddleware(rl)

	t.Run("Authenticated and Allowed", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/", nil)

		// Set user in context
		user := &models.User{ID: 100, Username: "testuser"}
		c.Set(auth.UserContextKey, user)

		mw(c)

		assert.NotEqual(t, http.StatusTooManyRequests, w.Code, "Expected request to be allowed")
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

		assert.Equal(t, http.StatusTooManyRequests, w.Code, "Expected status 429")
	})

	t.Run("Unauthenticated - Should Skip", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/", nil)

		// No user in context
		mw(c)

		assert.NotEqual(t, http.StatusTooManyRequests, w.Code, "Unauthenticated requests should not be limited by user-level middleware")
		assert.False(t, c.IsAborted(), "Context should not be aborted for unauthenticated users")
	})
}
