package api

import (
	"digital-innovation/stratego/auth"
	"digital-innovation/stratego/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// parseID extracts an integer ID from a query parameter
func parseID(r *http.Request, key string) (int, error) {
	idStr := r.URL.Query().Get(key)
	if idStr == "" {
		return 0, fmt.Errorf("missing %s", key)
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid %s", key)
	}

	return id, nil
}

// sendJSON encodes data to JSON and sends it with the specified status code
func sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// If encoding fails, we've already sent the header, so we can't send a different status
		// But we can log it
		fmt.Printf("Failed to encode JSON response: %v\n", err)
	}
}

// sendNoContent sends a 204 No Content response
func sendNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// sendError sends a plain text error message with the specified status code
func sendError(w http.ResponseWriter, message string, statusCode int) {
	http.Error(w, message, statusCode)
}

// ensureAuthenticated checks if a user is logged in, otherwise sends an error
func ensureAuthenticated(w http.ResponseWriter, r *http.Request) *models.User {
	user := auth.GetCurrentUser(r)
	if user == nil {
		sendError(w, "Unauthorized: Please login", http.StatusUnauthorized)
		return nil
	}
	return user
}
