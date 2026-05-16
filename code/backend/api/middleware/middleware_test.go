package middleware_test

import (
	"digital-innovation/gostrategy/api/middleware"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

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

	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
}

func TestJSONLoggerMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.JSONLoggerMiddleware())
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test?token=secret", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCSRFMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Save/Restore APP_ENV
	oldEnv := os.Getenv("APP_ENV")
	defer func() { _ = os.Setenv("APP_ENV", oldEnv) }()

	setup := func() *gin.Engine {
		r := gin.New()
		r.Use(middleware.CSRFMiddleware())
		r.POST("/test", func(c *gin.Context) { c.Status(http.StatusOK) })
		r.POST("/users/login", func(c *gin.Context) { c.Status(http.StatusOK) })
		return r
	}

	t.Run("Safe Methods pass", func(t *testing.T) {
		_ = os.Setenv("APP_ENV", "production")
		r := setup()
		r.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Exempt Paths pass", func(t *testing.T) {
		_ = os.Setenv("APP_ENV", "production")
		r := setup()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/users/login", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Development bypass", func(t *testing.T) {
		_ = os.Setenv("APP_ENV", "development")
		r := setup()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Production fails without token", func(t *testing.T) {
		_ = os.Setenv("APP_ENV", "production")
		r := setup()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Production succeeds with valid token", func(t *testing.T) {
		_ = os.Setenv("APP_ENV", "production")
		r := setup()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", nil)
		req.Header.Set("X-XSRF-TOKEN", "secret")
		req.AddCookie(&http.Cookie{Name: "XSRF-TOKEN", Value: "secret", HttpOnly: true, Secure: true, SameSite: http.SameSiteStrictMode})
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
