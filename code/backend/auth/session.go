package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sync"
	"time"
)

// Session represents a user session
type Session struct {
	ID        string
	UserID    int
	Username  string
	CreatedAt time.Time
	ExpiresAt time.Time
}

// SessionStore manages active sessions
type SessionStore struct {
	sessions map[string]*Session
	mutex    sync.RWMutex
}

var Store = &SessionStore{
	sessions: make(map[string]*Session),
}

// CreateSession creates a new session for a user
func (s *SessionStore) CreateSession(userID int, username string) (*Session, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return nil, err
	}

	session := &Session{
		ID:        sessionID,
		UserID:    userID,
		Username:  username,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
	}

	s.mutex.Lock()
	s.sessions[sessionID] = session
	s.mutex.Unlock()

	return session, nil
}

// GetSession retrieves a session by ID
func (s *SessionStore) GetSession(sessionID string) (*Session, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		return nil, false
	}

	// Check if session expired
	if time.Now().After(session.ExpiresAt) {
		go s.DeleteSession(sessionID)
		return nil, false
	}

	return session, true
}

// DeleteSession removes a session
func (s *SessionStore) DeleteSession(sessionID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.sessions, sessionID)
}

// CleanupExpiredSessions removes expired sessions
func (s *SessionStore) CleanupExpiredSessions() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now()
	for id, session := range s.sessions {
		if now.After(session.ExpiresAt) {
			delete(s.sessions, id)
		}
	}
}

// StartCleanupRoutine starts a background routine to cleanup expired sessions
func (s *SessionStore) StartCleanupRoutine() {
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			s.CleanupExpiredSessions()
		}
	}()
}

// generateSessionID generates a random session ID
func generateSessionID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate session ID: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
