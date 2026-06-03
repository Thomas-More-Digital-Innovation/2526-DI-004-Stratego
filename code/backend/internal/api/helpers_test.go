// Package api provides helper functions for the API.
package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"digital-innovation/gostrategy/internal/api/core"
	"digital-innovation/gostrategy/internal/api/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestParseID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Valid Path Param", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Params = gin.Params{{Key: "id", Value: "123"}}

		id, err := core.ParseID(c, "id")
		assert.NoError(t, err, "Unexpected error: %v", err)
		assert.Equal(t, id, 123, "Expected 123, got %d", id)
	})

	t.Run("Valid Query Param", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		req, _ := http.NewRequest("GET", "/?user_id=456", nil)
		c.Request = req

		id, err := core.ParseID(c, "user_id")
		assert.NoError(t, err)
		assert.Equal(t, 456, id)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Params = gin.Params{{Key: "id", Value: "abc"}}

		_, err := core.ParseID(c, "id")
		assert.Error(t, err)
	})

	t.Run("Missing ID", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())

		id, err := core.ParseID(c, "id")
		assert.NoError(t, err)
		assert.Equal(t, 0, id)
	})
}

func TestSecurityMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.SecurityMiddleware())
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	expectedHeaders := map[string]string{
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options":        "DENY",
		"X-XSS-Protection":       "1; mode=block",
	}

	for header, expectedValue := range expectedHeaders {
		assert.Equal(t, expectedValue, w.Header().Get(header))
	}
}
