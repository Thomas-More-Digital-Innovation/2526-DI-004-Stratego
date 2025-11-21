package auth

import (
	"digital-innovation/stratego/models"
	"net/http"
)

// Middleware checks if user is authenticated
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "Unauthorized: Please login", http.StatusUnauthorized)
			return
		}

		session, exists := Store.GetSession(cookie.Value)
		if !exists {
			http.Error(w, "Unauthorized: Invalid or expired session", http.StatusUnauthorized)
			return
		}

		// Add session to context for handlers to use
		r.Header.Set("X-User-ID", string(rune(session.UserID)))
		r.Header.Set("X-Username", session.Username)

		next(w, r)
	}
}

// OptionalAuth allows guests but identifies logged-in users
func OptionalAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err == nil {
			if session, exists := Store.GetSession(cookie.Value); exists {
				r.Header.Set("X-User-ID", string(rune(session.UserID)))
				r.Header.Set("X-Username", session.Username)
			}
		}
		next(w, r)
	}
}

// GetCurrentUser extracts user info from request (after auth middleware)
func GetCurrentUser(r *http.Request) *models.User {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil
	}

	session, exists := Store.GetSession(cookie.Value)
	if !exists {
		return nil
	}

	return &models.User{
		ID:       session.UserID,
		Username: session.Username,
	}
}

// SetSessionCookie sets the session cookie in response
func SetSessionCookie(w http.ResponseWriter, sessionID string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60, // 7 days
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		// Secure: true, // Enable in production with HTTPS
	})
}

// ClearSessionCookie removes the session cookie
func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
}
