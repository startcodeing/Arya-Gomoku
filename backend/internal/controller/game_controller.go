// Package controller handles HTTP requests and responses for the Gomoku API
// This file contains game room management endpoints (reserved for future PVP feature)
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gomoku-backend/internal/service"
)

// GameController handles multiplayer game room requests (reserved for future PVP feature)
type GameController struct {
	gameService *service.GameService
}

// NewGameController creates a new game controller instance
func NewGameController() *GameController {
	return &GameController{
		gameService: service.NewGameService(),
	}
}

// StartMatch handles POST /api/match/start requests (reserved for future implementation)
func (gc *GameController) StartMatch(c *gin.Context) {
	// This endpoint will be used to start online PVP matches
	// Currently returns a placeholder response
	
	var request struct {
		PlayerID string `json:"playerId" binding:"required"`
		GameMode string `json:"gameMode"` // "pvp" or "ai"
	}
	
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}
	
	// For now, only AI mode is supported
	if request.GameMode == "ai" {
		c.JSON(http.StatusOK, gin.H{
			"message":  "AI match started",
			"gameMode": "ai",
			"playerId": request.PlayerID,
		})
		return
	}
	
	// PVP mode - reserved for future implementation
	if request.GameMode == "pvp" {
		room, err := gc.gameService.CreateRoom(request.PlayerID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to create room",
				"details": err.Error(),
			})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{
			"message":  "Room created, waiting for opponent",
			"roomId":   room.RoomID,
			"gameMode": "pvp",
			"status":   room.Status,
		})
		return
	}
	
	c.JSON(http.StatusBadRequest, gin.H{
		"error": "Invalid game mode. Use 'ai' or 'pvp'",
	})
}

// JoinMatch handles POST /api/match/join requests (reserved for future implementation)
func (gc *GameController) JoinMatch(c *gin.Context) {
	var request struct {
		RoomID   string `json:"roomId" binding:"required"`
		PlayerID string `json:"playerId" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}
	
	room, err := gc.gameService.JoinRoom(request.RoomID, request.PlayerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to join room",
			"details": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully joined room",
		"room":    room,
	})
}

// GetMatchStatus handles GET /api/match/status/:roomId requests (reserved for future implementation)
func (gc *GameController) GetMatchStatus(c *gin.Context) {
	roomID := c.Param("roomId")
	
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Room ID is required",
		})
		return
	}
	
	room, err := gc.gameService.GetRoom(roomID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Room not found",
			"details": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"room": room,
	})
}

// MakeMove handles POST /api/match/:roomId/move requests (reserved for future implementation)
func (gc *GameController) MakeMove(c *gin.Context) {
	roomID := c.Param("roomId")
	
	var request struct {
		PlayerID string `json:"playerId" binding:"required"`
		X        int    `json:"x" binding:"required"`
		Y        int    `json:"y" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}
	
	// Validate coordinates
	if request.X < 0 || request.X >= 15 || request.Y < 0 || request.Y >= 15 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Coordinates must be within 0-14 range",
		})
		return
	}
	
	room, err := gc.gameService.MakeMove(roomID, request.PlayerID, request.X, request.Y)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to make move",
			"details": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Move made successfully",
		"room":    room,
	})
}

// LeaveMatch handles POST /api/match/:roomId/leave requests (reserved for future implementation)
func (gc *GameController) LeaveMatch(c *gin.Context) {
	roomID := c.Param("roomId")
	
	var request struct {
		PlayerID string `json:"playerId" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}
	
	err := gc.gameService.LeaveRoom(roomID, request.PlayerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to leave room",
			"details": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully left the room",
	})
}

// GetActiveRooms handles GET /api/match/rooms requests (reserved for future implementation)
func (gc *GameController) GetActiveRooms(c *gin.Context) {
	activeCount := gc.gameService.GetActiveRooms()
	
	c.JSON(http.StatusOK, gin.H{
		"activeRooms": activeCount,
		"message":     "Active rooms count retrieved",
	})
}

// HandleWebSocket handles WebSocket connections for real-time gameplay (reserved for future implementation)
func (gc *GameController) HandleWebSocket(c *gin.Context) {
	// This method will handle WebSocket connections for real-time PVP gameplay
	// Implementation will be added when PVP feature is developed
	
	c.JSON(http.StatusNotImplemented, gin.H{
		"error":   "WebSocket support not implemented yet",
		"message": "This feature is reserved for future PVP implementation",
	})
}