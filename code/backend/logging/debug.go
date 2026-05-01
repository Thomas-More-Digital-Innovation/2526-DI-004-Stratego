package logging

import (
	"digital-innovation/stratego/utils"
	"fmt"
	"log"
)

// Debug logs are only shown if STRATEGO_DEBUG environment variable is set to "true"
func Debug(tag string, format string, v ...any) {
	if utils.GetEnv("STRATEGO_DEBUG", "true") != "true" {
		return
	}

	msg := fmt.Sprintf(TagDebug+" "+format, v...)
	log.Printf("%s %s", tag, sanitize(msg))
}

// DebugWithUser logs a debug message with standardized user information
func DebugWithUser(tag string, username string, userID int, format string, v ...any) {
	if utils.GetEnv("STRATEGO_DEBUG", "true") != "true" {
		return
	}

	userPart := fmt.Sprintf("[User: %s] ", FormatUser(username, userID))
	msg := fmt.Sprintf(TagDebug+" "+userPart+format, v...)
	log.Printf("%s %s", tag, sanitize(msg))
}
