package telemetry

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestIPFilterMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("allows none when env is empty", func(t *testing.T) {
		os.Clearenv()
		router := gin.New()
		router.Use(IPFilterMiddleware())
		router.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("X-Forwarded-For", "192.168.1.10")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("filters by specific allowed IPs", func(t *testing.T) {
		err := os.Setenv("TELEMETRY_ALLOWED_IPS", "192.168.1.100, 10.0.0.1")
		assert.NoError(t, err)
		defer os.Clearenv()

		router := gin.New()
		router.Use(IPFilterMiddleware())
		router.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		// allowed IP
		req1 := httptest.NewRequest(http.MethodGet, "/test", nil)
		req1.Header.Set("X-Forwarded-For", "192.168.1.100")
		w1 := httptest.NewRecorder()
		router.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusOK, w1.Code)

		// disallowed IP
		req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
		req2.Header.Set("X-Forwarded-For", "192.168.1.101")
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusForbidden, w2.Code)
	})

	t.Run("filters by CIDR block", func(t *testing.T) {
		err := os.Setenv("TELEMETRY_ALLOWED_IPS", "192.168.1.0/24")
		assert.NoError(t, err)
		defer os.Clearenv()

		router := gin.New()
		router.Use(IPFilterMiddleware())
		router.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		// allowed inside CIDR
		req1 := httptest.NewRequest(http.MethodGet, "/test", nil)
		req1.Header.Set("X-Forwarded-For", "192.168.1.50")
		w1 := httptest.NewRecorder()
		router.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusOK, w1.Code)

		// disallowed outside CIDR
		req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
		req2.Header.Set("X-Forwarded-For", "192.168.2.50")
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusForbidden, w2.Code)
	})
}
