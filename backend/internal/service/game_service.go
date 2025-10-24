// Package service contains the business logic for the Gomoku game
// This file implements game room management and online match services (reserved for future PVP feature)
package service

import (
	"fmt"
	"sync"
	"time"

	"gomoku-backend/internal/model"
)

// GameService manages game rooms and online matches (reserved for future implementation)
type GameService struct {
	rooms map[string]*model.MatchRoom
	mutex sync.RWMutex
}

// NewGameService creates a new game service instance
func NewGameService() *GameService {
	return &GameService{
		rooms: make(map[string]*model.MatchRoom),
	}
}

// CreateRoom creates a new match room (reserved for future PVP feature)
func (gs *GameService) CreateRoom(playerID string) (*model.MatchRoom, error) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	
	roomID := gs.generateRoomID()
	
	room := &model.MatchRoom{
		RoomID: roomID,
		Players: []model.Player{
			{
				ID:   1,
				Type: "human",
				Name: playerID,
			},
		},
		Board:     model.NewBoard(),
		Status:    "waiting",
		CreatedAt: time.Now(),
	}
	
	gs.rooms[roomID] = room
	return room, nil
}

// JoinRoom allows a player to join an existing room (reserved for future PVP feature)
func (gs *GameService) JoinRoom(roomID, playerID string) (*model.MatchRoom, error) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	
	room, exists := gs.rooms[roomID]
	if !exists {
		return nil, fmt.Errorf("room not found")
	}
	
	if len(room.Players) >= 2 {
		return nil, fmt.Errorf("room is full")
	}
	
	if room.Status != "waiting" {
		return nil, fmt.Errorf("room is not accepting new players")
	}
	
	// Add second player
	room.Players = append(room.Players, model.Player{
		ID:   2,
		Type: "human",
		Name: playerID,
	})
	
	// Start the game
	room.Status = "playing"
	
	return room, nil
}

// GetRoom retrieves a room by ID (reserved for future PVP feature)
func (gs *GameService) GetRoom(roomID string) (*model.MatchRoom, error) {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	
	room, exists := gs.rooms[roomID]
	if !exists {
		return nil, fmt.Errorf("room not found")
	}
	
	return room, nil
}

// MakeMove processes a move in a multiplayer room (reserved for future PVP feature)
func (gs *GameService) MakeMove(roomID string, playerID string, x, y int) (*model.MatchRoom, error) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	
	room, exists := gs.rooms[roomID]
	if !exists {
		return nil, fmt.Errorf("room not found")
	}
	
	if room.Status != "playing" {
		return nil, fmt.Errorf("game is not in progress")
	}
	
	// Find player
	var playerNum int
	for i, player := range room.Players {
		if player.Name == playerID {
			playerNum = i + 1
			break
		}
	}
	
	if playerNum == 0 {
		return nil, fmt.Errorf("player not found in room")
	}
	
	if room.Board.CurrentPlayer != playerNum {
		return nil, fmt.Errorf("not your turn")
	}
	
	// Make the move
	if !room.Board.MakeMove(x, y, playerNum) {
		return nil, fmt.Errorf("invalid move")
	}
	
	// Check game state
	lastMove := &model.Move{X: x, Y: y, Player: playerNum}
	gameState := room.Board.GetGameState(lastMove)
	
	if gameState.Status != "playing" {
		room.Status = "finished"
	}
	
	return room, nil
}

// CleanupRooms removes old inactive rooms (reserved for future implementation)
func (gs *GameService) CleanupRooms() {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	
	cutoff := time.Now().Add(-2 * time.Hour) // Remove rooms older than 2 hours
	
	for roomID, room := range gs.rooms {
		if room.CreatedAt.Before(cutoff) {
			delete(gs.rooms, roomID)
		}
	}
}

// GetActiveRooms returns the number of active rooms (reserved for future implementation)
func (gs *GameService) GetActiveRooms() int {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	
	return len(gs.rooms)
}

// generateRoomID creates a unique room identifier
func (gs *GameService) generateRoomID() string {
	// Simple room ID generation - in production, use a more robust method
	return fmt.Sprintf("room_%d", time.Now().UnixNano())
}

// BroadcastToRoom sends a message to all players in a room (reserved for WebSocket implementation)
func (gs *GameService) BroadcastToRoom(roomID string, message interface{}) error {
	// TODO: Implement WebSocket broadcasting when PVP feature is added
	// This method will be used to send real-time updates to all players in a room
	return fmt.Errorf("WebSocket broadcasting not implemented yet")
}

// GetRoomStatus returns the current status of a room (reserved for future implementation)
func (gs *GameService) GetRoomStatus(roomID string) (string, error) {
	room, err := gs.GetRoom(roomID)
	if err != nil {
		return "", err
	}
	
	return room.Status, nil
}

// LeaveRoom removes a player from a room (reserved for future implementation)
func (gs *GameService) LeaveRoom(roomID, playerID string) error {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	
	room, exists := gs.rooms[roomID]
	if !exists {
		return fmt.Errorf("room not found")
	}
	
	// Remove player from room
	for i, player := range room.Players {
		if player.Name == playerID {
			room.Players = append(room.Players[:i], room.Players[i+1:]...)
			break
		}
	}
	
	// If room becomes empty, delete it
	if len(room.Players) == 0 {
		delete(gs.rooms, roomID)
	} else if room.Status == "playing" {
		// If game was in progress, mark as finished
		room.Status = "finished"
	}
	
	return nil
}