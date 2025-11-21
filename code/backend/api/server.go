package api

import (
	AIhandler "digital-innovation/stratego/ai/handler"
	"digital-innovation/stratego/engine"
	"digital-innovation/stratego/game"
	"digital-innovation/stratego/models"
	"fmt"
	"log"
	"net/http"
	"sync"
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

// setCORSHeaders sets appropriate CORS headers for API requests
func setCORSHeaders(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	// Allow requests from frontend dev server and localhost
	if origin == "http://localhost:5173" || origin == "http://127.0.0.1:5173" || origin == "" {
		if origin == "" {
			origin = "http://localhost:5173"
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
	} else {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// CreateGame creates a new game session
func (s *GameServer) CreateGame(gameID string, gameType string, ai1, ai2 string) (*GameSessionHandler, error) {
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
		controller2 = AIhandler.CreateAI(ai1, &player2)

	case models.AiVsAi:
		player1 := engine.NewPlayer(0, "AI Red", "red")
		player2 := engine.NewPlayer(1, "AI Blue", "blue")
		controller1 = AIhandler.CreateAI(ai1, &player1)
		controller2 = AIhandler.CreateAI(ai2, &player2)

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

// GetSession returns a game session handler
func (s *GameServer) GetSession(gameID string) (*GameSessionHandler, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	handler, exists := s.sessions[gameID]
	return handler, exists
}

// StartServer starts the HTTP server
func (s *GameServer) StartServer(addr string) error {
	// User & Auth endpoints
	http.HandleFunc("/api/users/register", func(w http.ResponseWriter, r *http.Request) {
		setCORSHeaders(w, r)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		s.RegisterUserHandler(w, r)
	})

	http.HandleFunc("/api/users/login", func(w http.ResponseWriter, r *http.Request) {
		setCORSHeaders(w, r)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		s.LoginHandler(w, r)
	})

	http.HandleFunc("/api/users/logout", func(w http.ResponseWriter, r *http.Request) {
		setCORSHeaders(w, r)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		s.LogoutHandler(w, r)
	})

	http.HandleFunc("/api/users/me", func(w http.ResponseWriter, r *http.Request) {
		setCORSHeaders(w, r)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		s.GetCurrentUserHandler(w, r)
	})

	http.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		setCORSHeaders(w, r)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		s.GetUserHandler(w, r)
	})

	http.HandleFunc("/api/users/stats", func(w http.ResponseWriter, r *http.Request) {
		setCORSHeaders(w, r)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		s.GetUserStatsHandler(w, r)
	})

	// Board setup endpoints
	http.HandleFunc("/api/board-setups", func(w http.ResponseWriter, r *http.Request) {
		setCORSHeaders(w, r)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		switch r.Method {
		case http.MethodGet:
			if r.URL.Query().Get("id") != "" {
				s.GetBoardSetupHandler(w, r)
			} else {
				s.GetUserBoardSetupsHandler(w, r)
			}
		case http.MethodPost:
			s.CreateBoardSetupHandler(w, r)
		case http.MethodPut:
			s.UpdateBoardSetupHandler(w, r)
		case http.MethodDelete:
			s.DeleteBoardSetupHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Existing game endpoints
	http.HandleFunc("/api/games", func(w http.ResponseWriter, r *http.Request) {
		setCORSHeaders(w, r)
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
