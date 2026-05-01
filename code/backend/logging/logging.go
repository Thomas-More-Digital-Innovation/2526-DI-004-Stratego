package logging

import (
	"digital-innovation/stratego/utils"
	"fmt"
	"log"
)

const (
	TagAuth     = "🔐 [   AUTH   ]"
	TagGame     = "🎮 [   GAME   ]"
	TagWeb      = "🖥️ [   WEB    ]"
	TagError    = "❌ [  ERROR   ]"
	TagSecurity = "🛡️ [ SECURITY ]"
	TagDebug    = "{DEBUG}"
)

var isDebugEnabled = utils.GetEnv("STRATEGO_DEBUG", "false") == "true"

// FormatUser returns a standardized user identifier for logs
func FormatUser(username string, userID int) string {
	if username == "" && userID <= 0 {
		return "Guest"
	}
	return fmt.Sprintf("%s (ID: %d)", username, userID)
}

// Debug logs are only shown if STRATEGO_DEBUG environment variable is set to "true"
func Debug(tag string, format string, v ...any) {
	if !isDebugEnabled {
		return
	}
	// Combine category tag with debug tag
	fullFormat := tag + " " + TagDebug + " " + format
	log.Printf(fullFormat, v...)
}

// Auth
func UserRegistered(username string, userID int) {
	log.Printf("%s User registered: %s", TagAuth, FormatUser(username, userID))
}

func UserLoggedIn(username string, userID int) {
	log.Printf("%s User logged in: %s", TagAuth, FormatUser(username, userID))
}

// Game
func GameStarted(gameID, gameType string, username string, userID int) {
	log.Printf("%s Game started: %s (Type: %s) by %s", TagGame, gameID, gameType, FormatUser(username, userID))
}

func GameFinished(gameID string, winner, loser string, rounds int) {
	log.Printf("%s Game finished: %s | Winner: %s | Loser: %s | Rounds: %d", TagGame, gameID, winner, loser, rounds)
}

func GameAborted(gameID string, reason string, username string, userID int) {
	log.Printf("%s Game aborted: %s (Reason: %s) by %s", TagGame, gameID, reason, FormatUser(username, userID))
}

// Connections
func ConnectionError(gameID string, username string, userID int, err error) {
	log.Printf("%s WebSocket error [User: %s] in game %s: %v", TagWeb, FormatUser(username, userID), gameID, err)
}

func ConnectionClosed(gameID string, username string, userID int) {
	log.Printf("%s Connection closed: Game: %s, User: %s", TagWeb, gameID, FormatUser(username, userID))
}

// Errors & Security
func Error(message string, err error) {
	log.Printf("%s %s: %v", TagError, message, err)
}
func ErrorWithUser(message string, username string, userID int, err error) {
	log.Printf("%s %s [User: %s]: %v", TagError, message, FormatUser(username, userID), err)
}

func ErrorWith2Users(message string, username1 string, userID1 int, username2 string, userID2 int, err error) {
	log.Printf("%s %s [User 1: %s, User 2: %s]: %v", TagError, message, FormatUser(username1, userID1), FormatUser(username2, userID2), err)
}

func ErrorWithIP(message string, ip string, err error) {
	log.Printf("%s %s [IP: %s]: %v", TagError, message, ip, err)
}

func SecurityWarning(message string, details string, username string, userID int) {
	log.Printf("%s %s [User: %s] | Details: %s", TagSecurity, message, FormatUser(username, userID), details)
}

func SecurityWarningWithIP(message string, details string, ip string) {
	log.Printf("%s %s [IP: %s] | Details: %s", TagSecurity, message, ip, details)
}
