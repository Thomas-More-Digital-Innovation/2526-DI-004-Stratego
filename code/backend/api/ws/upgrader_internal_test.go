package ws

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpgrader_CheckOrigin(t *testing.T) {
	// Save current env and restore after test
	oldProd := os.Getenv("APP_ENV")
	oldAllowed := os.Getenv("ALLOWED_ORIGINS")
	defer func() {
		_ = os.Setenv("APP_ENV", oldProd)
		_ = os.Setenv("ALLOWED_ORIGINS", oldAllowed)
	}()

	t.Run("allow empty origin in dev", func(t *testing.T) {
		_ = os.Setenv("APP_ENV", "development")
		req := httptest.NewRequest(http.MethodGet, "/ws", nil)
		assert.True(t, upgrader.CheckOrigin(req))
	})

	t.Run("reject empty origin in production", func(t *testing.T) {
		_ = os.Setenv("APP_ENV", "production")
		req := httptest.NewRequest(http.MethodGet, "/ws", nil)
		assert.False(t, upgrader.CheckOrigin(req))
	})

	t.Run("allow matching origin", func(t *testing.T) {
		_ = os.Setenv("ALLOWED_ORIGINS", "http://localhost:3000,http://example.com")
		req := httptest.NewRequest(http.MethodGet, "/ws", nil)
		req.Header.Set("Origin", "http://example.com")
		assert.True(t, upgrader.CheckOrigin(req))
	})

	t.Run("reject non-matching origin", func(t *testing.T) {
		_ = os.Setenv("ALLOWED_ORIGINS", "http://localhost:3000")
		req := httptest.NewRequest(http.MethodGet, "/ws", nil)
		req.Header.Set("Origin", "http://evil.com")
		assert.False(t, upgrader.CheckOrigin(req))
	})

	t.Run("allow any origin in dev if not set", func(t *testing.T) {
		_ = os.Setenv("APP_ENV", "development")
		_ = os.Setenv("ALLOWED_ORIGINS", "")
		req := httptest.NewRequest(http.MethodGet, "/ws", nil)
		req.Header.Set("Origin", "http://any-origin.com")
		assert.True(t, upgrader.CheckOrigin(req))
	})
}
