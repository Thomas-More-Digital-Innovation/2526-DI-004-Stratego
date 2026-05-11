// Package core provides shared utilities for API.
package core

import (
	"digital-innovation/gostrategy/api/dto"
	"digital-innovation/gostrategy/auth"
	"digital-innovation/gostrategy/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SendError helper to be used when shifting from net/http to Gin
func SendError(c *gin.Context, message string, statusCode int) {
	c.JSON(statusCode, gin.H{dto.MsgTypeError: message})
}

// SendJSON helper
func SendJSON(c *gin.Context, data any, statusCode int) {
	c.JSON(statusCode, data)
}

// EnsureAuthenticated checks if a user is logged in, otherwise sends an error
func EnsureAuthenticated(c *gin.Context) *models.User {
	user := auth.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{dto.MsgTypeError: "Unauthorized: Please login"})
		c.Abort()
		return nil
	}
	return user
}

// ParseID extracts an integer ID from a path or query parameter
func ParseID(c *gin.Context, key string) (int, error) {
	idStr := c.Param(key)
	if idStr == "" {
		idStr = c.Query(key)
	}
	if idStr == "" {
		return 0, nil
	}
	return strconv.Atoi(idStr)
}
