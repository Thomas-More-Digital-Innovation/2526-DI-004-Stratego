package api

import (
	"digital-innovation/gostrategy/internal/api/middleware"
	"digital-innovation/gostrategy/internal/api/rest/handlers"
	"digital-innovation/gostrategy/internal/auth"
	"digital-innovation/gostrategy/internal/utils"
	"strings"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/time/rate"

	_ "digital-innovation/gostrategy/docs" // Required for Swagger UI

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"digital-innovation/gostrategy/internal/telemetry"
)

// SetupRoutes registers all API routes and middleware
func (s *GameServer) SetupRoutes() {
	telemetry.Init()
	h := &handlers.Handler{GameServer: s.GameServer}

	s.Router.Use(gin.Recovery())
	s.Router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/health"},
	}))

	// Configure CORS
	corsConfig := cors.DefaultConfig()
	allowedOrigins, err := utils.GetEnvOrError("ALLOWED_ORIGINS")
	if err == nil {
		corsConfig.AllowOrigins = strings.Split(allowedOrigins, ",")
	} else {
		if !utils.IsProduction() {
			corsConfig.AllowOrigins = []string{"http://localhost:5000"}
		} else {
			corsConfig.AllowOrigins = []string{"https://gostrategy.dotsem.be"}
		}
	}
	corsConfig.AllowCredentials = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Content-Type", "Authorization", "X-Requested-With", "X-XSRF-TOKEN"}
	s.Router.Use(cors.New(corsConfig))

	// Middlewares
	s.Router.Use(middleware.PrometheusMiddleware())
	s.Router.Use(middleware.SecurityMiddleware())
	s.Router.Use(middleware.JSONLoggerMiddleware())
	s.Router.Use(middleware.CSRFMiddleware())

	// Rate limiters
	globalLimiter := middleware.NewRateLimiter(rate.Limit(10), 20)
	authLimiter := middleware.NewRateLimiter(rate.Every(time.Hour/5), 5)       // 5 per hour
	actionLimiter := middleware.NewRateLimiter(rate.Every(time.Minute/10), 10) // 10 per minute

	s.Router.Use(middleware.IPRateLimitMiddleware(globalLimiter))

	// Health & metrics
	s.Router.GET("/health", h.HealthHandler)
	s.Router.GET("/csrf-token", h.GetCSRFToken)

	metricsGroup := s.Router.Group("/metrics")
	metricsGroup.Use(telemetry.IPFilterMiddleware())
	metricsGroup.GET("", gin.WrapH(promhttp.Handler()))

	// Debug and Documentation (Dev only)
	if !utils.IsProduction() {
		s.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		s.Router.Static("/asyncapi", "./docs/asyncapi-docs")

		// s.Router.GET("/debug/stats", h.DebugStats) // Add if needed
	}

	// User & Auth endpoints
	users := s.Router.Group("/users")
	{
		users.POST("/register", h.RegisterUserHandler)
		users.POST("/login", h.LoginHandler)
		users.POST("/logout", h.LogoutHandler)
		users.POST("/refresh", h.RefreshTokenHandler)

		// Authenticated user routes
		me := users.Group("/me")
		me.Use(auth.RequireAuth())
		{
			me.GET("", h.GetCurrentUserHandler)
			me.GET("/stats", h.GetCurrentUserStatsHandler)
			me.GET("/games", h.HandleListMyGames)
			me.GET("/reconnectable", h.HandleGetReconnectableGame)
			me.POST("/password", middleware.UserRateLimitMiddleware(authLimiter), h.ChangePasswordHandler)
		}

		// Public info
		users.GET("/count", h.UserCountHandler)
		users.GET("/:id", h.GetUserHandler)
		users.GET("/:id/games", h.HandleListUserGames)
		users.GET("/stats", h.GetUserStatsHandler)
	}

	// Board setup endpoints (all require auth)
	setups := s.Router.Group("/board-setups")
	setups.Use(auth.RequireAuth())
	setups.Use(middleware.UserRateLimitMiddleware(actionLimiter))
	{
		setups.GET("", h.GetUserBoardSetupsHandler)
		setups.GET("/:id", h.GetBoardSetupHandler)
		setups.POST("", h.CreateBoardSetupHandler)
		setups.PUT("/:id", h.UpdateBoardSetupHandler)
		setups.DELETE("/:id", h.DeleteBoardSetupHandler)
	}

	// Game endpoints
	games := s.Router.Group("/games")
	games.Use(auth.OptionalAuth())
	{
		games.POST("", middleware.UserRateLimitMiddleware(actionLimiter), h.HandleCreateGame)
		games.GET("", h.HandleListGames)
		games.GET("/:id/history", h.HandleGetGameHistory)
		games.DELETE("/:id", auth.RequireAuth(), h.HandleDeleteGame)
		games.GET("/count", h.GamesPlayedCountHandler)
	}

	// WebSocket endpoint
	s.Router.GET("/game/:gameID", auth.OptionalAuth(), h.HandleWebSocketConnection)
}
