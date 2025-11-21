package api

import (
	"digital-innovation/stratego/auth"
	"digital-innovation/stratego/db"
	"digital-innovation/stratego/game"
	"digital-innovation/stratego/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// HTTP Handlers

// HandleCreateGame handles POST /api/games
func (s *GameServer) HandleCreateGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		GameID   string `json:"gameId"`
		GameType string `json:"gameType"`
		AI1      string `json:"ai1"`
		AI2      string `json:"ai2"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.GameID == "" {
		req.GameID = fmt.Sprintf("game-%d-%d", time.Now().Unix(), time.Now().UnixNano()%1000000)
	}

	if req.GameType == "" {
		req.GameType = models.HumanVsAi
	}

	_, err := s.CreateGame(req.GameID, req.GameType, req.AI1, req.AI2) // TODO: build logic for frontend to select AI type
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]any{
		"gameId":   req.GameID,
		"gameType": req.GameType,
		"wsUrl":    fmt.Sprintf("/ws/game/%s", req.GameID),
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	log.Printf("Created game %s (type: %s)", req.GameID, req.GameType)
}

// HandleWebSocketConnection handles WebSocket connections
// GET /ws/game/{gameID}?player={0|1|spectator}
func (s *GameServer) HandleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	// Extract game ID from path
	gameID := r.URL.Path[len("/ws/game/"):]
	if gameID == "" {
		http.Error(w, "Game ID required", http.StatusBadRequest)
		return
	}

	handler, exists := s.GetSession(gameID)
	if !exists {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	// Get player ID from query parameter
	playerIDStr := r.URL.Query().Get("player")
	playerID := -1 // Default to spectator

	switch playerIDStr {
	case "0":
		playerID = 0
	case "1":
		playerID = 1
	}

	// Try to get user from session cookie and associate with game
	user := auth.GetCurrentUser(r)
	if user != nil && handler.Session.Player1UserID == nil && playerID == 0 {
		// First player connecting with a logged-in account
		handler.Session.Player1UserID = &user.ID
		log.Printf("Associated game %s player 0 with user ID %d (%s)", gameID, user.ID, user.Username)
	} else if user != nil && handler.Session.Player2UserID == nil && playerID == 1 {
		// Second player connecting with a logged-in account
		handler.Session.Player2UserID = &user.ID
		log.Printf("Associated game %s player 1 with user ID %d (%s)", gameID, user.ID, user.Username)
	}

	log.Printf("WebSocket connection for game %s (player %d)", gameID, playerID)

	HandleWebSocket(w, r, handler.Session, handler.Hub, playerID)
}

// HandleListGames handles GET /api/games
func (s *GameServer) HandleListGames(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	games := make([]map[string]interface{}, 0, len(s.sessions))
	for gameID, handler := range s.sessions {
		state := handler.Session.GetGameState()
		games = append(games, map[string]interface{}{
			"gameId":     gameID,
			"round":      state.Round,
			"isRunning":  handler.Session.IsRunning(),
			"isGameOver": state.IsGameOver,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(games); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// handleGameOver broadcasts final game state and saves stats
func (s *GameServer) handleGameOver(session *game.GameSession, hub *WSHub) {
	state := session.GetGameState()
	winner := session.GetWinner()
	var winnerID *int
	var winnerName string
	if winner != nil {
		id := winner.GetID()
		winnerID = &id
		winnerName = winner.GetName()
	}

	gameOverMsg := GameOverMessage{
		WinnerID:   winnerID,
		WinnerName: winnerName,
		WinCause:   string(session.GetWinCause()),
		Round:      state.Round,
	}

	hub.BroadcastMessage(MsgTypeGameOver, gameOverMsg)

	// Broadcast final board state with all pieces revealed
	s.broadcastBoardStateRevealed(hub)

	// Save game stats to database
	go s.saveGameStats(session, winnerID)

	// Wait longer before stopping monitoring so users can see results
	time.Sleep(30 * time.Second)
}

// saveGameStats saves game statistics to the database
func (s *GameServer) saveGameStats(session *game.GameSession, winnerID *int) {
	duration := time.Since(session.StartTime).Seconds()
	state := session.GetGameState()

	// Track stats for player 1 if they have a user ID
	if session.Player1UserID != nil {
		userID := *session.Player1UserID
		won := winnerID != nil && *winnerID == 0

		if err := db.UpdateUserStats(userID, won, state.MoveCount, duration); err != nil {
			log.Printf("Failed to update stats for user %d: %v", userID, err)
		} else {
			log.Printf("Updated stats for user %d (won=%v)", userID, won)
		}
	}

	// Track stats for player 2 if they have a user ID
	if session.Player2UserID != nil {
		userID := *session.Player2UserID
		won := winnerID != nil && *winnerID == 1

		if err := db.UpdateUserStats(userID, won, state.MoveCount, duration); err != nil {
			log.Printf("Failed to update stats for user %d: %v", userID, err)
		} else {
			log.Printf("Updated stats for user %d (won=%v)", userID, won)
		}
	}
}
