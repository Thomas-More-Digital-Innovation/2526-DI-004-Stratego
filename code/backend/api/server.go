package api

import (
	AIhandler "digital-innovation/stratego/ai/handler"
	"digital-innovation/stratego/auth"
	"digital-innovation/stratego/engine"
	"digital-innovation/stratego/game"
	"digital-innovation/stratego/models"
	"digital-innovation/stratego/utils"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// GameServer manages HTTP and WebSocket connections
type GameServer struct {
	sessions map[string]*GameSessionHandler
	mutex    sync.RWMutex
	router   *gin.Engine
}

// GameSessionHandler wraps a game session with its WebSocket hub
type GameSessionHandler struct {
	Session  *game.GameSession
	Hub      *WSHub
	GameType string
}

func NewGameServer() *GameServer {
	// Set Gin mode based on environment
	if utils.GetEnv("APP_ENV", "development") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	return &GameServer{
		sessions: make(map[string]*GameSessionHandler),
		router:   gin.Default(),
	}
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
	// Configure CORS
	corsConfig := cors.DefaultConfig()
	// Whitelist allowed origins
	allowedOrigins := utils.GetEnv("ALLOWED_ORIGINS", "")
	corsConfig.AllowOrigins = strings.Split(allowedOrigins, ",")
	corsConfig.AllowCredentials = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Content-Type", "Authorization", "X-Requested-With", "X-XSRF-TOKEN"}
	s.router.Use(cors.New(corsConfig))

	// CSRF Protection
	s.router.Use(CSRFMiddleware())

	// Rate Limiting (5 requests per second per IP, burst of 10)
	limiter := NewIPRateLimiter(rate.Limit(5), 10)
	s.router.Use(RateLimitMiddleware(limiter))

	// Health check
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// User & Auth endpoints
	users := s.router.Group("/users")
	{
		users.POST("/register", s.RegisterUserHandler)
		users.POST("/login", s.LoginHandler)
		users.POST("/logout", s.LogoutHandler)

		// Authenticated user routes
		me := users.Group("/me")
		me.Use(auth.RequireAuth())
		{
			me.GET("", s.GetCurrentUserHandler)
			me.GET("/stats", s.GetCurrentUserStatsHandler)
		}

		// Public/Optional user info
		users.GET("/:id", s.GetUserHandler)
		users.GET("/stats", s.GetUserStatsHandler)
	}

	// Board setup endpoints (all require auth)
	setups := s.router.Group("/board-setups")
	setups.Use(auth.RequireAuth())
	{
		setups.GET("", s.GetUserBoardSetupsHandler)
		setups.GET("/:id", s.GetBoardSetupHandler)
		setups.POST("", s.CreateBoardSetupHandler)
		setups.PUT("/:id", s.UpdateBoardSetupHandler)
		setups.DELETE("/:id", s.DeleteBoardSetupHandler)
	}

	// Game endpoints
	games := s.router.Group("/games")
	{
		games.POST("", s.HandleCreateGame)
		games.GET("", s.HandleListGames)
	}

	// WebSocket endpoint
	s.router.GET("/game/:gameID", s.HandleWebSocketConnection)

	log.Printf("Starting game server on %s", addr)
	return s.router.Run(addr)
}
