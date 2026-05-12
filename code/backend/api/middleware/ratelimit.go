// Package middleware provides middleware for API handlers.
package middleware

import (
	"digital-innovation/gostrategy/api/dto"
	"digital-innovation/gostrategy/auth"
	"digital-innovation/gostrategy/logging"
	"digital-innovation/gostrategy/utils"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

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
		if len(i.ips) >= 10000 {
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
			c.JSON(http.StatusTooManyRequests, gin.H{dto.MsgTypeError: "Too many requests"})
			c.Abort()
			return
		}
		c.Next()
	}
}
