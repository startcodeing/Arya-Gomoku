// Package main is the entry point for the Gomoku backend server
// This file sets up the Gin web server, configures CORS, and initializes routes
package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"gomoku-backend/internal/controller"
	"gomoku-backend/internal/service"
)

func main() {
	// Initialize Gin router
	r := gin.Default()

	// Configure CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:5173",
		"http://0.0.0.0:5173",
		"http://192.168.0.109:5173",
	}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// Initialize services
	llmService := service.NewLLMService()

	// Initialize controllers
	aiController := controller.NewAIController()
	gameController := controller.NewGameController()
	llmController := controller.NewLLMController(llmService)

	// Setup routes
	api := r.Group("/api")
	{
		// AI endpoints
		api.POST("/ai/move", aiController.GetAIMove)
		api.GET("/ai/status", aiController.GetGameStatus)
		api.POST("/ai/reset", aiController.ResetGame)

		// LLM endpoints
		api.POST("/llm/start", llmController.StartGame)
		api.POST("/llm/move", llmController.MakeMove)
		api.GET("/llm/game/:id", llmController.GetGame)
		api.DELETE("/llm/game/:id", llmController.DeleteGame)
		api.GET("/llm/game/:id/history", llmController.GetGameHistory)
		api.GET("/llm/models", llmController.GetModels)
		api.PUT("/llm/config/:model", llmController.UpdateConfig)
		api.GET("/llm/config/:model", llmController.GetConfig)
		api.GET("/llm/stats", llmController.GetGameStats)
		api.GET("/llm/health", llmController.HealthCheck)

		// PVP Room endpoints
		api.POST("/rooms", gameController.CreateRoom)
		api.GET("/rooms", gameController.GetActiveRooms)
		api.GET("/rooms/:id", gameController.GetRoom)
		api.POST("/rooms/:id/join", gameController.JoinRoom)
		api.POST("/rooms/:id/start", gameController.StartGame)
		api.POST("/rooms/:id/move", gameController.MakeMove)
		api.POST("/rooms/:id/leave", gameController.LeaveRoom)
		api.POST("/rooms/:id/ready", gameController.SetPlayerReady)

		// WebSocket endpoint
		api.GET("/ws", gameController.HandleWebSocket)
	}

	// Start server on port 8081, bind to all interfaces for LAN access
	log.Println("Starting Gomoku backend server on 0.0.0.0:8081")
	if err := http.ListenAndServe("0.0.0.0:8081", r); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
