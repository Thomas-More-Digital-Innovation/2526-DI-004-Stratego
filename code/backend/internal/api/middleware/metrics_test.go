package middleware_test

import (
	"digital-innovation/gostrategy/internal/api/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func getMetricValue(metricName string, labels map[string]string) float64 {
	metricFamilies, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		return 0
	}
	for _, mf := range metricFamilies {
		if mf.GetName() != metricName {
			continue
		}
		for _, m := range mf.Metric {
			match := true
			for ln, lv := range labels {
				labelMatch := false
				for _, l := range m.Label {
					if l.GetName() == ln && l.GetValue() == lv {
						labelMatch = true
						break
					}
				}
				if !labelMatch {
					match = false
					break
				}
			}
			if match {
				if m.Counter != nil {
					return m.GetCounter().GetValue()
				}
				if m.Gauge != nil {
					return m.GetGauge().GetValue()
				}
			}
		}
	}
	return 0
}

func TestPrometheusMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	setupRouter := func() *gin.Engine {
		r := gin.New()
		r.Use(middleware.PrometheusMiddleware())
		r.GET("/test-metrics", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})
		r.POST("/test-error", func(c *gin.Context) {
			c.Status(http.StatusInternalServerError)
		})
		return r
	}

	t.Run("records successful request metrics", func(t *testing.T) {
		r := setupRouter()
		initialCount := getMetricValue("http_requests_total", map[string]string{
			"method":      "GET",
			"path":        "/test-metrics",
			"status_code": "200",
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test-metrics", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		finalCount := getMetricValue("http_requests_total", map[string]string{
			"method":      "GET",
			"path":        "/test-metrics",
			"status_code": "200",
		})
		assert.Equal(t, initialCount+1, finalCount)
	})

	t.Run("records error status code metrics", func(t *testing.T) {
		r := setupRouter()
		initialCount := getMetricValue("http_requests_total", map[string]string{
			"method":      "POST",
			"path":        "/test-error",
			"status_code": "500",
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test-error", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		finalCount := getMetricValue("http_requests_total", map[string]string{
			"method":      "POST",
			"path":        "/test-error",
			"status_code": "500",
		})
		assert.Equal(t, initialCount+1, finalCount)
	})

	t.Run("records unmatched routes", func(t *testing.T) {
		r := setupRouter()
		initialCount := getMetricValue("http_requests_total", map[string]string{
			"method":      "GET",
			"path":        "unmatched",
			"status_code": "404",
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/does-not-exist", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		finalCount := getMetricValue("http_requests_total", map[string]string{
			"method":      "GET",
			"path":        "unmatched",
			"status_code": "404",
		})
		assert.Equal(t, initialCount+1, finalCount)
	})
}

func TestWebSocketConnections(t *testing.T) {
	t.Run("increments and decrements active websocket gauge", func(t *testing.T) {
		initialVal := getMetricValue("websocket_connections_active", nil)

		middleware.WebSocketConnected()
		assert.Equal(t, initialVal+1, getMetricValue("websocket_connections_active", nil))

		middleware.WebSocketConnected()
		assert.Equal(t, initialVal+2, getMetricValue("websocket_connections_active", nil))

		middleware.WebSocketDisconnected()
		assert.Equal(t, initialVal+1, getMetricValue("websocket_connections_active", nil))

		middleware.WebSocketDisconnected()
		assert.Equal(t, initialVal, getMetricValue("websocket_connections_active", nil))
	})
}
