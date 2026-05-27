// Package api provides helper functions for the API.
package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"digital-innovation/gostrategy/internal/api/core"
	"digital-innovation/gostrategy/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

func TestParseID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Valid Path Param", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Params = gin.Params{{Key: "id", Value: "123"}}

		id, err := core.ParseID(c, "id")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if id != 123 {
			t.Errorf("Expected 123, got %d", id)
		}
	})

	t.Run("Valid Query Param", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		req, _ := http.NewRequest("GET", "/?user_id=456", nil)
		c.Request = req

		id, err := core.ParseID(c, "user_id")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if id != 456 {
			t.Errorf("Expected 456, got %d", id)
		}
	})

	t.Run("Invalid ID", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Params = gin.Params{{Key: "id", Value: "abc"}}

		_, err := core.ParseID(c, "id")
		if err == nil {
			t.Error("Expected error for invalid ID, got nil")
		}
	})

	t.Run("Missing ID", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())

		id, err := core.ParseID(c, "id")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if id != 0 {
			t.Errorf("Expected 0 for missing ID, got %d", id)
		}
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
		if got := w.Header().Get(header); got != expectedValue {
			t.Errorf("Expected header %s to be %q, got %q", header, expectedValue, got)
		}
	}
}
