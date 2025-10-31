package api

import (
	"digital-innovation/stratego/ai/fafo"
	"digital-innovation/stratego/engine"
	"digital-innovation/stratego/game"
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
	mu       sync.RWMutex
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
	case "human-vs-ai":
		player1 := engine.NewPlayer(0, "Human Player", "red")
		player2 := engine.NewPlayer(1, "AI Player", "blue")
		controller1 = engine.NewHumanPlayerController(&player1)
		controller2 = fafo.NewFafoAI(&player2)

	case "ai-vs-ai":
		player1 := engine.NewPlayer(0, "AI Red", "red")
		player2 := engine.NewPlayer(1, "AI Blue", "blue")
		controller1 = fafo.NewFafoAI(&player1)
		controller2 = fafo.NewFafoAI(&player2)

	case "human-vs-human":
		player1 := engine.NewPlayer(0, "Human Red", "red")
		player2 := engine.NewPlayer(1, "Human Blue", "blue")
		controller1 = engine.NewHumanPlayerController(&player1)
		controller2 = engine.NewHumanPlayerController(&player2)

	default:
		return nil, fmt.Errorf("unknown game type: %s", gameType)
	}

	session := game.NewGameSession(gameID, controller1, controller2)

	// Setup the game board with pieces (random placement)
	g := session.GetGame()
	player1Pieces := game.RandomSetup(g.Players[0])
	player2Pieces := game.RandomSetup(g.Players[1])

	// Place pieces on the board
	if err := game.SetupGame(g, player1Pieces, player2Pieces); err != nil {
		return nil, fmt.Errorf("failed to setup game: %v", err)
	}

	hub := NewWSHub(session)

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

	// Start the game
	if err := session.Start(); err != nil {
		return nil, err
	}

	return handler, nil
}

// monitorGame watches for game events and broadcasts them
func (s *GameServer) monitorGame(handler *GameSessionHandler, gameType string) {
	session := handler.Session
	hub := handler.Hub

	// Add delay for AI vs AI games to make them watchable
	moveDelay := 100 * time.Millisecond
	if gameType == "ai-vs-ai" {
		moveDelay = 500 * time.Millisecond // Half second per move for AI vs AI
	}

	lastMoveCount := 0
	ticker := time.NewTicker(moveDelay)
	defer ticker.Stop()

	for {
		<-ticker.C

		if !session.IsRunning() && session.GetGameState().IsGameOver {
			// Game over - broadcast final state with all pieces revealed
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
			return
		}

		// Check for new moves
		state := session.GetGameState()
		if state.MoveCount > lastMoveCount {
			lastMoveCount = state.MoveCount

			// Broadcast game state update
			hub.BroadcastMessage(MsgTypeGameState, GameStateMessage{
				Round:              state.Round,
				CurrentPlayerID:    state.CurrentPlayerID,
				CurrentPlayerName:  state.CurrentPlayerName,
				IsGameOver:         state.IsGameOver,
				WinnerID:           state.WinnerID,
				Player1Score:       state.Player1Score,
				Player2Score:       state.Player2Score,
				WaitingForInput:    state.WaitingForInput,
				MoveCount:          state.MoveCount,
				Player1AlivePieces: state.Player1AlivePieces,
				Player2AlivePieces: state.Player2AlivePieces,
			})

			// For AI vs AI, show all pieces to spectators
			if gameType == "ai-vs-ai" {
				s.broadcastBoardStateRevealed(hub)
			} else {
				// For other modes, send personalized board state to each client
				s.broadcastBoardStatePerClient(hub)
			}
		}
	}
}

// broadcastBoardState sends board state to all clients
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
		GameType string `json:"gameType"` // "human-vs-ai", "ai-vs-ai", "human-vs-human"
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.GameID == "" {
		req.GameID = fmt.Sprintf("game-%d-%d", time.Now().Unix(), time.Now().UnixNano()%1000000)
	}

	if req.GameType == "" {
		req.GameType = "human-vs-ai"
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
	json.NewEncoder(w).Encode(response)

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
	json.NewEncoder(w).Encode(games)
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
