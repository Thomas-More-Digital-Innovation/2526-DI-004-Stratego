package api

import (
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

// handleGameOver broadcasts final game state
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

	// Wait longer before stopping monitoring so users can see results
	time.Sleep(30 * time.Second)
}
