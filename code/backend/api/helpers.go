package api

import (
	"digital-innovation/stratego/auth"
	"digital-innovation/stratego/models"
	"digital-innovation/stratego/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// parseID extracts an integer ID from a path or query parameter
func parseID(c *gin.Context, key string) (int, error) {
	// First check path parameter
	idStr := c.Param(key)
	if idStr == "" {
		// Then check query parameter
		idStr = c.Query(key)
	}

	if idStr == "" {
		return 0, nil // Optional path or query
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// SecurityMiddleware adds security headers to the response
func SecurityMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")

		if utils.IsProduction() {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; font-src 'self' https://fonts.gstatic.com; connect-src 'self' ws: wss:; img-src 'self' data:;")
		}

		c.Next()
	}
}

// ensureAuthenticated checks if a user is logged in, otherwise sends an error
func ensureAuthenticated(c *gin.Context) *models.User {
	user := auth.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Please login"})
		c.Abort()
		return nil
	}
	return user
}

// sendError helper to be used when shifting from net/http to Gin (optional, can just use c.JSON)
// But for now, we'll keep it to maintain similar API structure if needed
func sendError(c *gin.Context, message string, statusCode int) {
	c.JSON(statusCode, gin.H{"error": message})
}

// sendJSON helper (optional, can just use c.JSON)
func sendJSON(c *gin.Context, data interface{}, statusCode int) {
	c.JSON(statusCode, data)
}

// sendNoContent helper
func sendNoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
