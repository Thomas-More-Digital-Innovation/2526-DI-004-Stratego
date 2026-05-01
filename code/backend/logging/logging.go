package logging

import (
	"fmt"
	"log"
	"regexp"
	"runtime"
	"strings"
)

var nonPrintableRegex = regexp.MustCompile(`[[:cntrl:]]`)

const (
	TagAuth     = "🔐 [   AUTH   ]"
	TagGame     = "🎮 [   GAME   ]"
	TagWeb      = "🖥️ [   WEB    ]"
	TagError    = "❌ [  ERROR   ]"
	TagSecurity = "🛡️ [ SECURITY ]"
	TagDebug    = "{DEBUG}"
)

func init() {
	// Only set time flags, we will handle the location manually for custom ordering
	log.SetFlags(log.LstdFlags)
}

// sanitize removes control characters to prevent log injection and terminal manipulation
func sanitize(s string) string {
	s = nonPrintableRegex.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

// logOutput logs a message with a tag and location
// This custom logging impl is necessary to put the location after the tag
// It is used by all other log functions
// It is not intended to be used directly
func logOutput(tag string, message string) {
	location := getCallerLocation(3)
	log.Printf("%s %s: %s", tag, location, message)
}

// getCallerLocation returns the file:line of the caller at the specified depth
func getCallerLocation(depth int) string {
	_, file, line, ok := runtime.Caller(depth)
	if !ok {
		return "???:0"
	}
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	return fmt.Sprintf("%s:%d", short, line)
}

// FormatUser returns a standardized user identifier for logs
func FormatUser(username string, userID int) string {
	if username == "" && userID <= 0 {
		return "Guest"
	}
	return fmt.Sprintf("%s (ID: %d)", sanitize(username), userID)
}

// --------------------------------------------------------------------------------
// Auth
// --------------------------------------------------------------------------------

// Logs a registration event
func UserRegistered(username string, userID int) {
	logOutput(TagAuth, fmt.Sprintf("User registered: %s", FormatUser(username, userID)))
}

// Logs a login event
func UserLoggedIn(username string, userID int) {
	logOutput(TagAuth, fmt.Sprintf("User logged in: %s", FormatUser(username, userID)))
}

// --------------------------------------------------------------------------------
// Game
// --------------------------------------------------------------------------------

// Logs a game started event
func GameStarted(gameID, gameType string, username string, userID int) {
	logOutput(TagGame, fmt.Sprintf("Game started: %s (Type: %s) by %s", sanitize(gameID), sanitize(gameType), FormatUser(username, userID)))
}

// Logs a game finished event
func GameFinished(gameID string, winner, loser string, rounds int) {
	logOutput(TagGame, fmt.Sprintf("Game finished: %s | Winner: %s | Loser: %s | Rounds: %d", sanitize(gameID), sanitize(winner), sanitize(loser), rounds))
}

// Logs a game aborted event
func GameAborted(gameID string, reason string, username string, userID int) {
	logOutput(TagGame, fmt.Sprintf("Game aborted: %s (Reason: %s) by %s", sanitize(gameID), sanitize(reason), FormatUser(username, userID)))
}

// --------------------------------------------------------------------------------
// Connections
// --------------------------------------------------------------------------------

// Logs a connection error event
func ConnectionError(gameID string, username string, userID int, err error) {
	logOutput(TagWeb, fmt.Sprintf("WebSocket error [User: %s] in game %s: %v", FormatUser(username, userID), sanitize(gameID), err))
}

// Logs a connection closed event
func ConnectionClosed(gameID string, username string, userID int) {
	logOutput(TagWeb, fmt.Sprintf("Connection closed: Game: %s, User: %s", sanitize(gameID), FormatUser(username, userID)))
}

// --------------------------------------------------------------------------------
// Errors & Security
// --------------------------------------------------------------------------------

// Logs a generic error event
func Error(message string, err error) {
	logOutput(TagError, fmt.Sprintf("%s: %v", sanitize(message), err))
}

// Logs an error that includes a user in the message
func ErrorWithUser(message string, username string, userID int, err error) {
	logOutput(TagError, fmt.Sprintf("%s [User: %s]: %v", sanitize(message), FormatUser(username, userID), err))
}

// Logs an error that includes two users in the message
// Use this when 2 users are involved in the error
func ErrorWith2Users(message string, username1 string, userID1 int, username2 string, userID2 int, err error) {
	logOutput(TagError, fmt.Sprintf("%s [User 1: %s, User 2: %s]: %v", sanitize(message), FormatUser(username1, userID1), FormatUser(username2, userID2), err))
}

// Logs an error that includes an IP address in the message
func ErrorWithIP(message string, ip string, err error) {
	logOutput(TagError, fmt.Sprintf("%s [IP: %s]: %v", sanitize(message), sanitize(ip), err))
}

// Logs a security warning that includes a user in the message
func SecurityWarning(message string, details string, username string, userID int) {
	logOutput(TagSecurity, fmt.Sprintf("%s [User: %s] | Details: %s", sanitize(message), FormatUser(username, userID), sanitize(details)))
}

// Logs a security warning that includes an IP address in the message
func SecurityWarningWithIP(message string, details string, ip string) {
	logOutput(TagSecurity, fmt.Sprintf("%s [IP: %s] | Details: %s", sanitize(message), sanitize(ip), sanitize(details)))
}

// Fatal logs a message and exits the application
func Fatal(message string, err error) {
	logOutput(TagError, fmt.Sprintf("FATAL: %s: %v", sanitize(message), err))
	panic(err)
}

// Fatalf logs a formatted message and exits the application
func Fatalf(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	logOutput(TagError, fmt.Sprintf("FATAL: %s", msg))
	panic(msg)
}

// LogRaw prints a raw string to the log without any tags or sanitization.
// Use this only for pre-sanitized or structured data like JSON.
func LogRaw(s string) {
	log.Println(s)
}
