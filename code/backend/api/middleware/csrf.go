// Package middleware provides middleware for API handlers.
package middleware

import (
	"digital-innovation/gostrategy/api/dto"
	"digital-innovation/gostrategy/auth"
	"digital-innovation/gostrategy/logging"
	"digital-innovation/gostrategy/utils"
	"net/http"

	"github.com/gin-gonic/gin"
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
			c.JSON(http.StatusForbidden, gin.H{dto.MsgTypeError: "CSRF validation failed"})
			c.Abort()
			return
		}

		c.Next()
	}
}
