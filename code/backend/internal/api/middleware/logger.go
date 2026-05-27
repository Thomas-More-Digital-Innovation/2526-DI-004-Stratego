// Package middleware provides middleware for API handlers.
package middleware

import (
	"digital-innovation/gostrategy/internal/auth"
	"digital-innovation/gostrategy/internal/logging"
	"digital-innovation/gostrategy/internal/utils"
	"encoding/json"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// maskSensitiveParams redacts values of potentially sensitive query parameters
func maskSensitiveParams(query string) string {
	if query == "" {
		return ""
	}
	params := strings.Split(query, "&")
	for i, p := range params {
		kv := strings.SplitN(p, "=", 2)
		if len(kv) == 2 {
			key := strings.ToLower(kv[0])
			if strings.Contains(key, "token") || strings.Contains(key, "password") || strings.Contains(key, "secret") || strings.Contains(key, "key") {
				params[i] = kv[0] + "=[REDACTED]"
			}
		}
	}
	return strings.Join(params, "&")
}

// JSONLoggerMiddleware logs requests in JSON format
func JSONLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := maskSensitiveParams(c.Request.URL.RawQuery)

		if path == "/health" {
			c.Next()
			return
		}

		c.Next()

		fullPath := path
		if raw != "" {
			fullPath = path + "?" + raw
		}

		user := auth.GetCurrentUser(c)
		username, userID := utils.TryGetUser(user)

		logData := map[string]any{
			"time":    time.Now().Format(time.RFC3339),
			"latency": time.Since(start).String(),
			"ip":      c.ClientIP(),
			"method":  c.Request.Method,
			"path":    fullPath,
			"status":  c.Writer.Status(),
			"user":    logging.FormatUser(username, userID),
			"agent":   c.Request.UserAgent(),
		}

		if jsonBytes, err := json.Marshal(logData); err == nil {
			logging.LogRaw(string(jsonBytes))
		}
	}
}
