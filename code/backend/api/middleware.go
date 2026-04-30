package api

import (
	"digital-innovation/stratego/auth"
	"digital-innovation/stratego/utils"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// CSRFMiddleware requires a custom header for non-safe methods
func CSRFMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
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
			log.Printf("CSRF validation failed: header=%s, cookie=%s, err=%v", tokenHeader, tokenCookie, err)
			c.JSON(http.StatusForbidden, gin.H{"error": "CSRF validation failed"})
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

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	v, exists := i.ips[ip]
	if !exists {
		// Hard cap at 10,000 IPs to prevent OOM
		if len(i.ips) >= 10000 {
			// If we're at capacity, return a strict one-time limiter 
			// or we could evict a random entry. For simplicity, we just 
			// don't track the new IP until the next cleanup.
			return rate.NewLimiter(i.r, i.b)
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
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// JSONLoggerMiddleware logs requests in JSON format
func JSONLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Skip logging for health check because it is used by docker compose health check
		if path == "/health" {
			c.Next()
			return
		}

		// Process request
		c.Next()

		if raw != "" {
			path = path + "?" + raw
		}

		user := auth.GetCurrentUser(c)
		userName := "anonymous"
		if user != nil {
			userName = user.Username
		}

		// Use standard log package which we've redirected to both stdout and file
		// Gin will already log its own format, but this gives us a clean JSON object for Loki
		log.Printf(`{"time":"%s", "latency":"%s", "ip":"%s", "method":"%s", "path":"%s", "status":%d, "user":"%s", "agent":"%s"}`+"\n",
			time.Now().Format(time.RFC3339),
			time.Since(start),
			c.ClientIP(),
			c.Request.Method,
			path,
			c.Writer.Status(),
			userName,
			c.Request.UserAgent(),
		)
	}
}
