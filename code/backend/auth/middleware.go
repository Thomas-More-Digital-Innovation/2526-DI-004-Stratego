package auth

import (
	"digital-innovation/stratego/models"
	"digital-innovation/stratego/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

const UserContextKey = "user"

// RequireAuth checks if user is authenticated
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("session_id")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Please login"})
			c.Abort()
			return
		}

		user, err := VerifyToken(cookie)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized Token"})
			c.Abort()
			return
		}

		// Add user to context for handlers to use
		c.Set(UserContextKey, user)

		c.Next()
	}
}

// OptionalAuth allows guests but identifies logged-in users
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("session_id")
		if err == nil {
			if user, err := VerifyToken(cookie); err == nil {
				c.Set(UserContextKey, user)
			}
		}
		c.Next()
	}
}

// GetCurrentUser extracts user info from Gin context
func GetCurrentUser(c *gin.Context) *models.User {
	val, exists := c.Get(UserContextKey)
	if !exists {
		return nil
	}
	user, ok := val.(*models.User)
	if !ok {
		return nil
	}
	return user
}

var cookieSecure = utils.IsProduction()

// SetSessionCookie sets the session cookie in response
func SetSessionCookie(c *gin.Context, sessionID string) {
	c.SetCookie("session_id", sessionID, maxCookieAge, "/", "", cookieSecure, true)
}

// ClearSessionCookie removes the session cookie
func ClearSessionCookie(c *gin.Context) {
	c.SetCookie("session_id", "", -1, "/", "", cookieSecure, true)
}
