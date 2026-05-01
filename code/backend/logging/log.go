package logging

import (
	"digital-innovation/stratego/utils"
	"log"
)

const (
	TagAuth     = "🔐 [   AUTH   ]"
	TagGame     = "🎮 [   GAME   ]"
	TagWeb      = "🖥️ [   WEB    ]"
	TagError    = "❌ [  ERROR   ]"
	TagSecurity = "🛡️ [ SECURITY ]"
	TagDebug    = "[ DEBUG ]"
)

var isDebugEnabled = utils.GetEnv("STRATEGO_DEBUG", "false") == "true"

// Debug logs are only shown if STRATEGO_DEBUG environment variable is set to "true"
func Debug(tag string, format string, v ...any) {
	if !isDebugEnabled {
		return
	}
	// Combine TagDebug with the provided category tag
	msg := log.Prefix() + TagDebug + " " + tag + " " + format
	log.Printf(msg, v...)
}

// Auth
func UserRegistered(username string, userID int) {
	log.Printf("%s User registered: %s (ID: %d)", TagAuth, username, userID)
}

func UserLoggedIn(username string, userID int) {
	log.Printf("%s User logged in: %s (ID: %d)", TagAuth, username, userID)
}

// Game
func GameStarted(gameID, gameType string, userID int) {
	log.Printf("%s Game started: %s (Type: %s) by User: %d", TagGame, gameID, gameType, userID)
}

func GameFinished(gameID string, winnerName string, rounds int) {
	log.Printf("%s Game finished: %s | Winner: %s | Rounds: %d", TagGame, gameID, winnerName, rounds)
}

func GameAborted(gameID string, reason string) {
	log.Printf("%s Game aborted: %s (Reason: %s)", TagGame, gameID, reason)
}

// Connections
func ConnectionError(gameID string, err error) {
	log.Printf("%s WebSocket error in game %s: %v", TagWeb, gameID, err)
}

func ConnectionClosed(gameID string, clientID string) {
	log.Printf("%s Connection closed: Game: %s, Client: %s", TagWeb, gameID, clientID)
}

// Errors & Security
func Error(message string, err error) {
	log.Printf("%s %s: %v", TagError, message, err)
}

func SecurityWarning(message string, details string) {
	log.Printf("%s %s | Details: %s", TagSecurity, message, details)
}
