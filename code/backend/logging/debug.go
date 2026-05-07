// Package logging provides standardized logging utilities for the application
package logging

import (
	"digital-innovation/gostrategy/utils"
	"fmt"
	"log"
)

// Debug logs are only shown if GOSTRATEGY_DEBUG environment variable is set to "true"
func Debug(tag string, format string, v ...any) {
	if utils.GetEnv("GOSTRATEGY_DEBUG", "true") != "true" {
		return
	}

	msg := fmt.Sprintf(format, v...)
	log.Printf("%s %s %s", tag, TagDebug, sanitize(msg))
}

// DebugWithUser logs a debug message with standardized user information
func DebugWithUser(tag string, username string, userID int, format string, v ...any) {
	if utils.GetEnv("GOSTRATEGY_DEBUG", "true") != "true" {
		return
	}

	userPart := fmt.Sprintf("[User: %s] ", FormatUser(username, userID))
	msg := fmt.Sprintf(userPart+format, v...)
	log.Printf("%s %s %s", tag, TagDebug, sanitize(msg))
}
