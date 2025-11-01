// Package main is the entry point for the Gomoku backend server
// This file sets up the Gin web server, configures CORS, and initializes routes
package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"gomoku-backend/internal/config"
	"gomoku-backend/internal/controller"
	"gomoku-backend/internal/database"
	"gomoku-backend/internal/middleware"
	"gomoku-backend/internal/repository"
	"gomoku-backend/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize database
	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Auto migrate database
	if err := db.AutoMigrate(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Create indexes
	if err := db.CreateIndexes(); err != nil {
		log.Fatal("Failed to create indexes:", err)
	}

	// Initialize data
	if err := db.InitializeData(); err != nil {
		log.Fatal("Failed to initialize data:", err)
	}

	// Initialize repositories
	_ = repository.NewUserRepository(db.DB) // TODO: 将来可能需要用到
	gameRepo := repository.NewGameRepository(db.DB)
	statsRepo := repository.NewStatisticsRepository(db.DB)

	// Initialize services
	llmService := service.NewLLMService()
	authService := service.NewAuthService(db.DB)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWT.SecretKey, db.DB)

	// Initialize Gin router
	r := gin.Default()

	// Configure CORS middleware
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{
		"http://localhost:5173",
		"http://0.0.0.0:5173",
		"http://192.168.0.109:5173",
	}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(corsConfig))

	// Initialize controllers
	aiController := controller.NewAIController(gameRepo)
	llmController := controller.NewLLMController(llmService, gameRepo)
	gameController := controller.NewGameController(gameRepo)
	authController := controller.NewAuthController(authService, authMiddleware)
	statisticsController := controller.NewStatisticsController(gameRepo, statsRepo)

	// Setup routes
	api := r.Group("/api")
	{
		// Authentication endpoints (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
			auth.POST("/refresh", authController.RefreshToken)
			auth.POST("/validate", authController.ValidateToken)
		}

		// User endpoints (protected)
		user := api.Group("/user")
		user.Use(authMiddleware.RequireAuth())
		{
			user.GET("/profile", authController.GetProfile)
			user.PUT("/profile", authController.UpdateProfile)
			user.POST("/change-password", authController.ChangePassword)
			user.POST("/logout", authController.Logout)
			user.GET("/info/:id", authController.GetUserInfo) // Public user info
		}

		// AI endpoints (optionally protected)
		ai := api.Group("/ai")
		ai.Use(authMiddleware.OptionalAuth())
		{
			// Legacy AI endpoints
			ai.POST("/move", aiController.GetAIMove)
			ai.GET("/status", aiController.GetGameStatus)
			ai.POST("/reset", aiController.ResetGame)
			ai.GET("/stats", aiController.GetAIStats)
			ai.POST("/cache/clear", aiController.ClearCache)
			ai.GET("/difficulties", aiController.GetDifficultyLevels)
			ai.POST("/benchmark", aiController.BenchmarkAI)
			
			// New AI game management endpoints (protected)
			ai.POST("/games", authMiddleware.RequireAuth(), aiController.CreateAIGame)
			ai.GET("/games/:id", authMiddleware.RequireAuth(), aiController.GetAIGame)
			ai.POST("/games/:id/move", authMiddleware.RequireAuth(), aiController.MakeAIMove)
			ai.GET("/games", authMiddleware.RequireAuth(), aiController.GetUserAIGames)
		}

		// LLM endpoints (optionally protected)
		llm := api.Group("/llm")
		llm.Use(authMiddleware.OptionalAuth())
		{
			llm.POST("/start", llmController.StartGame)
			llm.POST("/move", llmController.MakeMove)
			llm.GET("/game/:id", llmController.GetGame)
			llm.DELETE("/game/:id", llmController.DeleteGame)
			llm.GET("/game/:id/history", llmController.GetGameHistory)
			llm.GET("/models", llmController.GetModels)
			llm.PUT("/config/:model", llmController.UpdateConfig)
			llm.GET("/config/:model", llmController.GetConfig)
			llm.GET("/stats", llmController.GetGameStats)
			llm.GET("/health", llmController.HealthCheck)
			llm.GET("/cache/stats", llmController.GetCacheStats)
			llm.DELETE("/cache", llmController.ClearCache)
			
			// New LLM game management endpoints
			llm.POST("/games", authMiddleware.RequireAuth(), llmController.CreateLLMGame)
			llm.GET("/games/:id", authMiddleware.RequireAuth(), llmController.GetLLMGame)
			llm.POST("/games/:id/move", authMiddleware.RequireAuth(), llmController.MakeLLMMove)
			llm.GET("/user/games", authMiddleware.RequireAuth(), llmController.GetUserLLMGames)
		}

		// PVP Room endpoints (protected)
		rooms := api.Group("/rooms")
		rooms.Use(authMiddleware.RequireAuth())
		{
			rooms.POST("", gameController.CreateRoom)
			rooms.GET("", gameController.GetActiveRooms)
			rooms.GET("/:id", gameController.GetRoom)
			rooms.POST("/:id/join", gameController.JoinRoom)
			rooms.POST("/:id/start", gameController.StartGame)
			rooms.POST("/:id/move", gameController.MakeMove)
			rooms.POST("/:id/leave", gameController.LeaveRoom)
			rooms.POST("/:id/ready", gameController.SetPlayerReady)
		}

		// Statistics endpoints (protected)
		stats := api.Group("/statistics")
		stats.Use(authMiddleware.RequireAuth())
		{
			stats.GET("/games", statisticsController.GetUserGameHistory)
			stats.GET("/user", statisticsController.GetUserGameStats)
			stats.GET("/system", authMiddleware.RequireAdmin(), statisticsController.GetSystemStatistics)
			stats.GET("/date-range", statisticsController.GetGameStatsByDateRange)
			stats.GET("/top-players", statisticsController.GetTopPlayers)
			stats.GET("/search", statisticsController.SearchGames)
			stats.GET("/game-types", statisticsController.GetGameTypeStatistics)
			stats.GET("/difficulties", statisticsController.GetDifficultyStatistics)
			stats.GET("/export", statisticsController.ExportGameData)
			stats.DELETE("/games/:id", statisticsController.DeleteGameRecord)
		}

		// WebSocket endpoint (protected)
		api.GET("/ws", authMiddleware.RequireAuth(), gameController.HandleWebSocket)
	}

	// Start server on port 8081, bind to all interfaces for LAN access
	log.Println("Starting Gomoku backend server on 0.0.0.0:8080")
	if err := http.ListenAndServe("0.0.0.0:8080", r); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
