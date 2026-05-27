// Package middleware provides middleware for API handlers.
package middleware

import (
	"digital-innovation/gostrategy/internal/api/dto"
	"digital-innovation/gostrategy/internal/auth"
	"digital-innovation/gostrategy/internal/logging"
	"digital-innovation/gostrategy/internal/utils"
	"net/http"
	"sync"
	"time"

	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimiter manages rate limiters for different keys (IP, User ID, etc.)
type RateLimiter struct {
	visitors map[string]*visitor
	mu       *sync.RWMutex
	r        rate.Limit
	b        int
	stopChan chan struct{}
}

// NewRateLimiter creates a new RateLimiter instance
func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		mu:       &sync.RWMutex{},
		r:        r,
		b:        b,
		stopChan: make(chan struct{}),
	}

	go rl.cleanupVisitors()

	return rl
}

// GetLimiter returns the rate limiter for a specific key
func (rl *RateLimiter) GetLimiter(key string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[key]
	if !exists {
		if len(rl.visitors) >= 10000 {
			// Basic eviction: clear all if limit reached
			// In a production app, we might use LRU or more granular eviction
			for k := range rl.visitors {
				delete(rl.visitors, k)
				break
			}
		}

		v = &visitor{
			limiter: rate.NewLimiter(rl.r, rl.b),
		}
		rl.visitors[key] = v
	}

	v.lastSeen = time.Now()
	return v.limiter
}

// Stop stops the rate limiter background cleaner goroutine.
func (rl *RateLimiter) Stop() {
	close(rl.stopChan)
}

func (rl *RateLimiter) cleanupVisitors() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.mu.Lock()
			for key, v := range rl.visitors {
				if time.Since(v.lastSeen) > 3*time.Hour {
					delete(rl.visitors, key)
				}
			}
			rl.mu.Unlock()
		case <-rl.stopChan:
			return
		}
	}
}

// IPRateLimitMiddleware limits requests per IP
func IPRateLimitMiddleware(limiter *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.GetLimiter(ip).Allow() {
			user := auth.GetCurrentUser(c)
			username, userID, err := utils.TryGetUserOrError(user)
			if err != nil {
				logging.SecurityWarningWithIP("IP Rate limit triggered", "Too many requests from this IP", ip)
			} else {
				logging.SecurityWarning("IP Rate limit triggered", "Too many requests from this IP", username, userID)
			}
			c.JSON(http.StatusTooManyRequests, gin.H{dto.MsgTypeError: "Too many requests from this IP"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// UserRateLimitMiddleware limits requests per authenticated user
func UserRateLimitMiddleware(limiter *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := auth.GetCurrentUser(c)
		if user == nil {
			c.Next()
			return
		}

		key := strconv.Itoa(user.ID)
		if !limiter.GetLimiter(key).Allow() {
			logging.SecurityWarning("User rate limit triggered", "Too many requests for this user account", user.Username, user.ID)
			c.JSON(http.StatusTooManyRequests, gin.H{dto.MsgTypeError: "Too many requests for this user account"})
			c.Abort()
			return
		}
		c.Next()
	}
}
