// Package auth provides authentication and authorization functionality
package auth

import (
	"digital-innovation/gostrategy/internal/db"
	"digital-innovation/gostrategy/internal/models"
	"digital-innovation/gostrategy/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserContextKey is the key used to store the user object in the Gin context
const UserContextKey = "user"

// RequireAuth checks if user is authenticated
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie(AccessTokenCookieName)
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

		// Inject into request context for DB RLS
		c.Request = c.Request.WithContext(db.WithUserID(c.Request.Context(), user.ID))

		c.Next()
	}
}

// OptionalAuth allows guests but identifies logged-in users
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie(AccessTokenCookieName)
		if err == nil {
			if user, err := VerifyToken(cookie); err == nil {
				c.Set(UserContextKey, user)
				// Inject into request context for DB RLS
				c.Request = c.Request.WithContext(db.WithUserID(c.Request.Context(), user.ID))
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

// SetSessionCookie sets the access token cookie in response
func SetSessionCookie(c *gin.Context, accessToken string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(AccessTokenCookieName, accessToken, MaxCookieAge, "/", "", cookieSecure, true)

	// Also ensure user has a CSRF token
	SetCSRFCookie(c)
}

// SetRefreshTokenCookie sets the refresh token cookie
func SetRefreshTokenCookie(c *gin.Context, refreshToken string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(RefreshTokenCookieName, refreshToken, MaxRefreshTokenAge, "/", "", cookieSecure, true)
}

// SetCSRFCookie generates and sets a new CSRF cookie
func SetCSRFCookie(c *gin.Context) string {
	csrfToken, _ := GenerateRefreshToken() // Reuse the random string generator
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(XSRFTokenCookieName, csrfToken, MaxCookieAge, "/", "", cookieSecure, false)
	return csrfToken
}

// ClearSessionCookie removes all auth-related cookies
func ClearSessionCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(AccessTokenCookieName, "", -1, "/", "", cookieSecure, true)
	c.SetCookie(RefreshTokenCookieName, "", -1, "/", "", cookieSecure, true)
	c.SetCookie(XSRFTokenCookieName, "", -1, "/", "", cookieSecure, false)
}
