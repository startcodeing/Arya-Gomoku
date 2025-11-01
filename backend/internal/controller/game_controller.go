// Package controller handles HTTP requests and responses for the Gomoku API
// This file contains game room management endpoints for PVP feature
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"gomoku-backend/internal/model"
	"gomoku-backend/internal/repository"
	"gomoku-backend/internal/service"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

// GameController handles multiplayer game room requests for PVP feature
type GameController struct {
	gameService *service.GameService
	hub         *service.Hub
}

// NewGameController creates a new game controller instance
func NewGameController(gameRepo repository.GameRepository) *GameController {
	gameService := service.NewGameService(gameRepo)
	hub := service.NewHub(gameService)
	go hub.Run()
	
	return &GameController{
		gameService: gameService,
		hub:         hub,
	}
}

// CreateRoom handles POST /api/rooms requests
func (gc *GameController) CreateRoom(c *gin.Context) {
	var request model.CreateRoomRequest
	
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}
	
	// Set default max players if not specified
	if request.MaxPlayers == 0 {
		request.MaxPlayers = 2
	}
	
	room, err := gc.gameService.CreateRoom(request.RoomName, request.PlayerName, request.MaxPlayers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create room",
			"details": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"room": room,
	})
}

// JoinRoom handles POST /api/rooms/:id/join requests
func (gc *GameController) JoinRoom(c *gin.Context) {
	roomID := c.Param("id")
	var request model.JoinRoomRequest
	
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}
	
	room, player, err := gc.gameService.JoinRoom(roomID, request.PlayerName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to join room",
			"details": err.Error(),
		})
		return
	}
	
	// First respond to the HTTP request
	c.JSON(http.StatusOK, gin.H{
		"room":   room,
		"player": player,
	})
	
	// Then broadcast room update to all clients in the room (async to avoid deadlock)
	go func() {
		gc.hub.BroadcastToRoom(roomID, model.WSMessage{
			Type: "room_update",
			Data: model.RoomUpdateData{
				Room: room,
			},
		})
	}()
}

// GetRoom handles GET /api/rooms/:id requests
func (gc *GameController) GetRoom(c *gin.Context) {
	roomID := c.Param("id")
	
	room := gc.gameService.GetRoom(roomID)
	if room == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Room not found",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"room": room,
	})
}

// StartGame handles POST /api/rooms/:id/start requests
func (gc *GameController) StartGame(c *gin.Context) {
	roomID := c.Param("id")
	
	err := gc.gameService.StartGame(roomID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to start game",
			"details": err.Error(),
		})
		return
	}
	
	room := gc.gameService.GetRoom(roomID)
	if room == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Room not found",
		})
		return
	}
	
	// Broadcast game start to all clients in the room
	gc.hub.BroadcastToRoom(roomID, model.WSMessage{
		Type: "game_start",
		Data: model.GameUpdateData{
			Game: room.Game,
		},
	})
	
	c.JSON(http.StatusOK, gin.H{
		"room": room,
	})
}

// LeaveRoom handles POST /api/rooms/:id/leave requests
func (gc *GameController) LeaveRoom(c *gin.Context) {
	roomID := c.Param("id")
	
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
	
	// Broadcast room update to remaining clients
	room := gc.gameService.GetRoom(roomID)
	if room != nil {
		gc.hub.BroadcastToRoom(roomID, model.WSMessage{
			Type: "room_update",
			Data: model.RoomUpdateData{
				Room: room,
			},
		})
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully left the room",
	})
}

// GetActiveRooms handles GET /api/rooms requests
func (gc *GameController) GetActiveRooms(c *gin.Context) {
	rooms := gc.gameService.GetActiveRooms()
	
	c.JSON(http.StatusOK, gin.H{
		"rooms": rooms,
	})
}

// HandleWebSocket handles WebSocket connections for real-time gameplay
func (gc *GameController) HandleWebSocket(c *gin.Context) {
	roomID := c.Query("roomId")
	playerID := c.Query("playerId")
	
	if roomID == "" || playerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "roomId and playerId are required",
		})
		return
	}
	
	// Use the hub's ServeWS method to handle the WebSocket connection
	gc.hub.ServeWS(c.Writer, c.Request, roomID, playerID)
}

// MakeMove handles POST /api/rooms/:id/move requests
func (gc *GameController) MakeMove(c *gin.Context) {
	roomID := c.Param("id")
	var request model.MakeMoveRequest
	
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}
	
	room, move, err := gc.gameService.MakeMove(roomID, request.PlayerID, request.X, request.Y)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to make move",
			"details": err.Error(),
		})
		return
	}
	
	// Broadcast move to all clients in the room
	gc.hub.BroadcastToRoom(roomID, model.WSMessage{
		Type: "game_update",
		Data: model.GameUpdateData{
			Game:     room.Game,
			LastMove: move,
		},
	})
	
	c.JSON(http.StatusOK, gin.H{
		"room": room,
		"move": move,
	})
}

// SetPlayerReady handles POST /api/rooms/:id/ready requests
func (gc *GameController) SetPlayerReady(c *gin.Context) {
	roomID := c.Param("id")
	
	var request struct {
		PlayerID string `json:"playerId" binding:"required"`
		Ready    bool   `json:"ready" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}
	
	err := gc.gameService.SetPlayerReady(roomID, request.PlayerID, request.Ready)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to set player ready status",
			"details": err.Error(),
		})
		return
	}
	
	room := gc.gameService.GetRoom(roomID)
	if room == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Room not found",
		})
		return
	}
	
	// Broadcast room update to all clients in the room
	gc.hub.BroadcastToRoom(roomID, model.WSMessage{
		Type: "room_update",
		Data: model.RoomUpdateData{
			Room: room,
		},
	})
	
	c.JSON(http.StatusOK, gin.H{
		"room": room,
	})
}