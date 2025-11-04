package api

import (
	"digital-innovation/stratego/ai/fafo"
	"digital-innovation/stratego/engine"
	"digital-innovation/stratego/game"
	"digital-innovation/stratego/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// GameServer manages HTTP and WebSocket connections
type GameServer struct {
	sessions map[string]*GameSessionHandler
	mutex    sync.RWMutex
}

// GameSessionHandler wraps a game session with its WebSocket hub
type GameSessionHandler struct {
	Session  *game.GameSession
	Hub      *WSHub
	GameType string
}

func NewGameServer() *GameServer {
	return &GameServer{
		sessions: make(map[string]*GameSessionHandler),
	}
}

// CreateGame creates a new game session
func (s *GameServer) CreateGame(gameID string, gameType string) (*GameSessionHandler, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.sessions[gameID]; exists {
		return nil, fmt.Errorf("game %s already exists", gameID)
	}

	var controller1, controller2 engine.PlayerController

	switch gameType {
	case models.HumanVsAi:
		player1 := engine.NewPlayer(0, "Human Player", "red")
		player2 := engine.NewPlayer(1, "AI Player", "blue")
		controller1 = engine.NewHumanPlayerController(&player1)
		controller2 = fafo.NewFafoAI(&player2)

	case models.AiVsAi:
		player1 := engine.NewPlayer(0, "AI Red", "red")
		player2 := engine.NewPlayer(1, "AI Blue", "blue")
		controller1 = fafo.NewFafoAI(&player1)
		controller2 = fafo.NewFafoAI(&player2)

	case models.HumanVsHuman:
		player1 := engine.NewPlayer(0, "Human Red", "red")
		player2 := engine.NewPlayer(1, "Human Blue", "blue")
		controller1 = engine.NewHumanPlayerController(&player1)
		controller2 = engine.NewHumanPlayerController(&player2)

	default:
		return nil, fmt.Errorf("unknown game type: %s", gameType)
	}

	session := game.NewGameSession(gameID, controller1, controller2)

	// For AI vs AI, setup immediately and start
	if gameType == models.AiVsAi {
		g := session.GetGame()
		player1Pieces := game.RandomSetup(g.Players[0])
		player2Pieces := game.RandomSetup(g.Players[1])

		// Place pieces on the board
		if err := game.SetupGame(g, player1Pieces, player2Pieces); err != nil {
			return nil, fmt.Errorf("failed to setup game: %v", err)
		}

		// Exit setup phase for AI vs AI
		session.SetSetupPhaseComplete()
	}
	// For human vs AI, pieces are already generated in NewGameSession
	// Game will stay in setup phase until human confirms

	hub := NewWSHub(session, gameType)

	handler := &GameSessionHandler{
		Session:  session,
		Hub:      hub,
		GameType: gameType,
	}

	s.sessions[gameID] = handler

	// Start the hub
	go hub.Run()

	// Start game monitoring for broadcasting moves
	go s.monitorGame(handler, gameType)

	// Start the game only for AI vs AI
	if gameType == models.AiVsAi {
		if err := session.Start(); err != nil {
			return nil, err
		}
	}

	return handler, nil
}

// monitorGame watches for game events and broadcasts them
func (s *GameServer) monitorGame(handler *GameSessionHandler, gameType string) {
	session := handler.Session
	hub := handler.Hub
	log.Printf("Starting game monitor for %s (type: %s)", handler.Session.ID, gameType)

	// Send initial state to all connected clients
	time.Sleep(100 * time.Millisecond) // Brief delay for clients to connect
	s.broadcastFullState(hub, gameType)

	// WAIT IN SETUP PHASE - WebSocket handlers will broadcast when user acts
	for session.IsSetupPhase() {
		time.Sleep(100 * time.Millisecond)
	}

	log.Printf("GameMonitor %s: Exiting setup phase, game starting", session.ID)

	// NOW we enter the game loop
	for {
		// Wait for a move notification with timeout
		if !session.WaitForMoveNotification(5 * time.Second) {
			// Timeout - check if game is over
			if !session.IsRunning() && session.GetGameState().IsGameOver {
				s.handleGameOver(session, hub)
				return
			}
			continue
		}

		// Move was executed
		log.Printf("Move executed in game %s", session.ID)

		// Check if combat occurred
		combat := session.GetLastCombat()
		hasCombat := combat != nil && combat.Occurred

		if hasCombat {
			log.Printf("Combat detected! Broadcasting combat data and waiting for animation")

			// Broadcast combat message (with piece info)
			s.broadcastCombat(hub, combat, gameType)

			// Wait for frontend animation to complete (3 second timeout)
			session.WaitForAnimationComplete(3 * time.Second)

			log.Printf("Animation complete, broadcasting updated state")

			// Clear combat after animation
			session.ClearLastCombat()

			// NOW broadcast state after animation (winner moves to position, loser removed)
			s.broadcastFullState(hub, gameType)
		} else {
			// No combat - broadcast state immediately
			s.broadcastFullState(hub, gameType)
		}

		// Signal that move has been processed - GameRunner can continue
		session.AckMoveProcessed()

		// Check if game is over
		state := session.GetGameState()
		if state.IsGameOver {
			time.Sleep(500 * time.Millisecond) // Brief delay before game over message
			s.handleGameOver(session, hub)
			return
		}
	}
}

// broadcastFullState sends complete game state and board to all clients
func (s *GameServer) broadcastFullState(hub *WSHub, gameType string) {
	state := hub.session.GetGameState()

	var winnerName string
	var winCause string
	if state.WinnerID != nil {
		winner := hub.session.GetWinner()
		if winner != nil {
			winnerName = winner.GetName()
		}
		winCause = string(hub.session.GetWinCause())
	}

	// Broadcast game state
	hub.BroadcastMessage(MsgTypeGameState, GameStateMessage{
		Round:              state.Round,
		CurrentPlayerID:    state.CurrentPlayerID,
		CurrentPlayerName:  state.CurrentPlayerName,
		IsGameOver:         state.IsGameOver,
		WinnerID:           state.WinnerID,
		WinnerName:         winnerName,
		WinCause:           winCause,
		Player1Score:       state.Player1Score,
		Player2Score:       state.Player2Score,
		WaitingForInput:    state.WaitingForInput,
		MoveCount:          state.MoveCount,
		Player1AlivePieces: state.Player1AlivePieces,
		Player2AlivePieces: state.Player2AlivePieces,
		IsSetupPhase:       state.IsSetupPhase,
	})

	// Broadcast board state
	switch {
	case state.IsSetupPhase:
		s.broadcastSetupBoard(hub, gameType)
	case gameType == models.AiVsAi:
		s.broadcastBoardStateRevealed(hub)
	default:
		s.broadcastBoardStatePerClient(hub)
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

// broadcastBoardState sends board state to all clients
//
//lint:ignore U1000 Ignore unused function temporarily for debugging
func (s *GameServer) broadcastBoardState(hub *WSHub, viewerID int) {
	board := hub.session.GetBoard()
	field := board.GetField()

	boardDTO := make([][]PieceDTO, 10)
	for y := 0; y < 10; y++ {
		boardDTO[y] = make([]PieceDTO, 10)
		for x := 0; x < 10; x++ {
			piece := field[y][x]
			if piece != nil && piece.IsAlive() {
				dto := PieceToDTO(piece, viewerID)
				dto.Position = PositionDTO{X: x, Y: y}
				boardDTO[y][x] = dto
			}
		}
	}

	boardMsg := BoardStateMessage{
		Board:  boardDTO,
		Width:  10,
		Height: 10,
	}

	hub.BroadcastMessage(MsgTypeBoardState, boardMsg)
}

// broadcastSetupBoard sends the setup board state (pieces not yet placed on board)
func (s *GameServer) broadcastSetupBoard(hub *WSHub, gameType string) {
	session := hub.session

	// Create empty board
	boardDTO := make([][]PieceDTO, 10)
	for y := 0; y < 10; y++ {
		boardDTO[y] = make([]PieceDTO, 10)
	}

	// Place player 1 pieces in setup area (rows 6-9)
	player1Pieces := session.GetSetupPieces(0)
	idx := 0
	for y := 6; y <= 9; y++ {
		for x := 0; x < 10; x++ {
			if idx < len(player1Pieces) {
				piece := player1Pieces[idx]
				dto := PieceToDTO(piece, 0) // Player 0 can see their own pieces
				dto.Position = PositionDTO{X: x, Y: y}
				boardDTO[y][x] = dto
				idx++
			}
		}
	}

	// Place player 2 pieces in setup area (rows 0-3)
	// In human vs AI, player 2 is AI, so pieces are hidden
	player2Pieces := session.GetSetupPieces(1)
	idx = 0
	for y := 0; y <= 3; y++ {
		for x := 0; x < 10; x++ {
			if idx < len(player2Pieces) {
				piece := player2Pieces[idx]
				// For human vs AI, hide AI pieces during setup
				viewerID := -1
				if gameType == models.AiVsAi {
					viewerID = 1 // Show all pieces in AI vs AI
				}
				dto := PieceToDTO(piece, viewerID)
				dto.Position = PositionDTO{X: x, Y: y}
				boardDTO[y][x] = dto
				idx++
			}
		}
	}

	boardMsg := BoardStateMessage{
		Board:  boardDTO,
		Width:  10,
		Height: 10,
	}

	hub.BroadcastMessage(MsgTypeBoardState, boardMsg)
}

// broadcastBoardStatePerClient sends personalized board state to each connected client
func (s *GameServer) broadcastBoardStatePerClient(hub *WSHub) {
	hub.mutex.RLock()
	clients := make([]*WSClient, 0, len(hub.clients))
	for client := range hub.clients {
		clients = append(clients, client)
	}
	hub.mutex.RUnlock()

	// Send personalized board state to each client
	for _, client := range clients {
		hub.sendBoardState(client)
	}
}

// broadcastCombat sends combat information to all clients
func (s *GameServer) broadcastCombat(hub *WSHub, combat *game.CombatResult, gameType string) {
	if combat == nil || !combat.Occurred {
		return
	}

	attacker := combat.AttackerPiece
	defender := combat.DefenderPiece

	// For AI vs AI or spectators, reveal both pieces
	// For player games, reveal based on ownership
	attackerDTO := PieceToDTO(attacker, attacker.GetOwner().GetID())
	attackerDTO.Position = PositionToDTO(combat.AttackerPosition)
	attackerDTO.Revealed = true

	defenderDTO := PieceToDTO(defender, defender.GetOwner().GetID())
	defenderDTO.Position = PositionToDTO(combat.DefenderPosition)
	defenderDTO.Revealed = true

	combatMsg := CombatMessage{
		Attacker:     attackerDTO,
		Defender:     defenderDTO,
		AttackerWon:  attacker.IsAlive(),
		DefenderWon:  defender.IsAlive(),
		AttackerDied: !attacker.IsAlive(),
		DefenderDied: !defender.IsAlive(),
	}

	hub.BroadcastMessage(MsgTypeCombat, combatMsg)
	log.Printf("Combat message sent: %+v", combatMsg)
}

// broadcastBoardStateRevealed sends board state with all pieces revealed (for AI vs AI spectating)
func (s *GameServer) broadcastBoardStateRevealed(hub *WSHub) {
	board := hub.session.GetBoard()
	field := board.GetField()

	boardDTO := make([][]PieceDTO, 10)
	for y := 0; y < 10; y++ {
		boardDTO[y] = make([]PieceDTO, 10)
		for x := 0; x < 10; x++ {
			piece := field[y][x]
			if piece != nil && piece.IsAlive() {
				// Force reveal all pieces for spectators
				dto := PieceToDTO(piece, piece.GetOwner().GetID())
				dto.Position = PositionDTO{X: x, Y: y}
				dto.Revealed = true
				boardDTO[y][x] = dto
			}
		}
	}

	boardMsg := BoardStateMessage{
		Board:  boardDTO,
		Width:  10,
		Height: 10,
	}

	hub.BroadcastMessage(MsgTypeBoardState, boardMsg)
}

// GetSession returns a game session handler
func (s *GameServer) GetSession(gameID string) (*GameSessionHandler, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	handler, exists := s.sessions[gameID]
	return handler, exists
}

// HTTP Handlers

// HandleCreateGame handles POST /api/games
func (s *GameServer) HandleCreateGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		GameID   string `json:"gameId"`
		GameType string `json:"gameType"` // "human-vs-ai", models.AiVsAi, models.HumanVsHuman
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

	_, err := s.CreateGame(req.GameID, req.GameType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
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

	if playerIDStr == "0" {
		playerID = 0
	} else if playerIDStr == "1" {
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

// StartServer starts the HTTP server
func (s *GameServer) StartServer(addr string) error {
	http.HandleFunc("/api/games", func(w http.ResponseWriter, r *http.Request) {
		// Enable CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method == http.MethodPost {
			s.HandleCreateGame(w, r)
		} else {
			s.HandleListGames(w, r)
		}
	})

	http.HandleFunc("/ws/game/", func(w http.ResponseWriter, r *http.Request) {
		s.HandleWebSocketConnection(w, r)
	})

	log.Printf("Starting game server on %s", addr)
	return http.ListenAndServe(addr, nil)
}
