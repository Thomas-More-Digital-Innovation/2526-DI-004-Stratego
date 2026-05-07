// Package api provides the HTTP and WebSocket API for the GoStrategy game
package api

import (
	"digital-innovation/gostrategy/auth"
	"digital-innovation/gostrategy/logging"
	"digital-innovation/gostrategy/utils"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// CSRFMiddleware requires a custom header for non-safe methods
func CSRFMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Exempt login and register from CSRF validation to support "cold start"
		path := c.Request.URL.Path
		if path == "/users/login" || path == "/users/register" {
			c.Next()
			return
		}

		// Skip for safe methods
		if c.Request.Method == http.MethodGet || c.Request.Method == http.MethodHead || c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		// Allow skipping in development
		if !utils.IsProduction() {
			c.Next()
			return
		}

		// Require X-XSRF-TOKEN header and compare with cookie
		tokenHeader := c.GetHeader("X-XSRF-TOKEN")
		tokenCookie, err := c.Cookie("XSRF-TOKEN")

		if tokenHeader == "" || err != nil || tokenHeader != tokenCookie {
			user := auth.GetCurrentUser(c)
			username, userID, err := utils.TryGetUserOrError(user)
			if err != nil {
				logging.SecurityWarningWithIP("CSRF validation failed", "Path: "+path, c.ClientIP())
			} else {
				logging.SecurityWarning("CSRF validation failed", "Path: "+path, username, userID)
			}
			c.JSON(http.StatusForbidden, gin.H{MsgTypeError: "CSRF validation failed"})
			c.Abort()
			return
		}

		c.Next()
	}
}

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// IPRateLimiter manages rate limiters for different IP addresses
type IPRateLimiter struct {
	ips map[string]*visitor
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

// NewIPRateLimiter creates a new IPRateLimiter instance
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	i := &IPRateLimiter{
		ips: make(map[string]*visitor),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}

	go i.cleanupVisitors()

	return i
}

// GetLimiter returns the rate limiter for a specific IP address
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	v, exists := i.ips[ip]
	if !exists {
		// Hard cap at 10,000 IPs to prevent OOM
		if len(i.ips) >= 10000 {
			// Evict one random entry to make room for the new one
			// TODO: there may be a better way to handle this
			for ipToEvict := range i.ips {
				delete(i.ips, ipToEvict)
				break
			}
		}

		v = &visitor{
			limiter: rate.NewLimiter(i.r, i.b),
		}
		i.ips[ip] = v
	}

	v.lastSeen = time.Now()
	return v.limiter
}

func (i *IPRateLimiter) cleanupVisitors() {
	for {
		time.Sleep(1 * time.Hour)

		i.mu.Lock()
		for ip, v := range i.ips {
			if time.Since(v.lastSeen) > 3*time.Hour {
				delete(i.ips, ip)
			}
		}
		i.mu.Unlock()
	}
}

// RateLimitMiddleware limits requests per IP
func RateLimitMiddleware(limiter *IPRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.GetLimiter(ip).Allow() {
			user := auth.GetCurrentUser(c)
			username, userID, err := utils.TryGetUserOrError(user)
			if err != nil {
				logging.SecurityWarningWithIP("Rate limit triggered", "Too many requests from this IP", ip)
			} else {
				logging.SecurityWarning("Rate limit triggered", "Too many requests", username, userID)
			}
			c.JSON(http.StatusTooManyRequests, gin.H{MsgTypeError: "Too many requests"})
			c.Abort()
			return
		}
		c.Next()
	}
}

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

		// Skip logging for health check because it is used by docker compose health check
		if path == "/health" {
			c.Next()
			return
		}

		// Process request
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
