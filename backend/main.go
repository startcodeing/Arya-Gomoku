// Package main is the entry point for the Gomoku backend server
// This file sets up the Gin web server, configures CORS, and initializes routes
package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"gomoku-backend/internal/controller"
)

func main() {
	// Initialize Gin router
	r := gin.Default()

	// Configure CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173", "http://192.168.110.136:5173"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// Initialize controllers
	aiController := controller.NewAIController()
	gameController := controller.NewGameController()

	// Setup routes
	api := r.Group("/api")
	{
		// AI endpoints
		api.POST("/ai/move", aiController.GetAIMove)
		api.GET("/ai/status", aiController.GetGameStatus)
		api.POST("/ai/reset", aiController.ResetGame)

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

	// Start server on port 8080
	log.Println("Starting Gomoku backend server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
