// Package service contains the business logic for the Gomoku game
// This file implements game room management and online match services for PVP feature
package service

import (
	"fmt"
	"sync"
	"time"

	"gomoku-backend/internal/model"
)

// GameService manages game rooms and online matches for PVP feature
type GameService struct {
	rooms map[string]*model.Room
	mutex sync.RWMutex
}

// NewGameService creates a new game service instance
func NewGameService() *GameService {
	return &GameService{
		rooms: make(map[string]*model.Room),
	}
}

// CreateRoom creates a new match room for PVP feature
func (gs *GameService) CreateRoom(roomName, playerName string, maxPlayers int) (*model.Room, error) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	
	room := model.NewRoom(roomName, playerName, maxPlayers)
	gs.rooms[room.ID] = room
	
	return room, nil
}

// JoinRoom allows a player to join an existing room
func (gs *GameService) JoinRoom(roomID, playerName string) (*model.Room, *model.PVPPlayer, error) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	
	room, exists := gs.rooms[roomID]
	if !exists {
		return nil, nil, fmt.Errorf("room not found")
	}
	
	if len(room.Players) >= room.MaxPlayers {
		return nil, nil, fmt.Errorf("room is full")
	}
	
	if room.Status != "waiting" {
		return nil, nil, fmt.Errorf("room is not accepting new players")
	}
	
	player := room.AddPlayer(playerName)
	if player == nil {
		return nil, nil, fmt.Errorf("failed to add player to room")
	}
	
	return room, player, nil
}

// GetRoom retrieves a room by ID
func (gs *GameService) GetRoom(roomID string) *model.Room {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	
	return gs.rooms[roomID]
}

// MakeMove processes a player's move in a room
func (gs *GameService) MakeMove(roomID, playerID string, x, y int) (*model.Room, *model.PVPMove, error) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	
	room, exists := gs.rooms[roomID]
	if !exists {
		return nil, nil, fmt.Errorf("room not found")
	}
	
	if room.Game == nil {
		return nil, nil, fmt.Errorf("game not started")
	}
	
	if room.Game.Status != "playing" {
		return nil, nil, fmt.Errorf("game is not in progress")
	}
	
	// Validate player
	player := room.GetPlayer(playerID)
	if player == nil {
		return nil, nil, fmt.Errorf("player not found in room")
	}
	
	// Validate turn
	if room.Game.CurrentPlayer != player.PlayerNumber {
		return nil, nil, fmt.Errorf("not your turn")
	}
	
	// Make the move
	move := room.Game.MakeMove(x, y, playerID, player.PlayerNumber)
	if move == nil {
		return nil, nil, fmt.Errorf("invalid move")
	}
	
	return room, move, nil
}

// CleanupRooms removes inactive rooms
func (gs *GameService) CleanupRooms() {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	
	now := time.Now()
	for roomID, room := range gs.rooms {
		// Remove rooms that have been inactive for more than 1 hour
		if now.Sub(room.CreatedAt) > time.Hour {
			delete(gs.rooms, roomID)
		}
	}
}

// GetActiveRooms returns all active rooms
func (gs *GameService) GetActiveRooms() []*model.Room {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	
	var activeRooms []*model.Room
	for _, room := range gs.rooms {
		if room.Status == "waiting" || room.Status == "playing" {
			activeRooms = append(activeRooms, room)
		}
	}
	
	return activeRooms
}

// StartGame starts the game in a room
func (gs *GameService) StartGame(roomID string) error {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	
	room, exists := gs.rooms[roomID]
	if !exists {
		return fmt.Errorf("room not found")
	}
	
	if !room.CanStartGame() {
		return fmt.Errorf("cannot start game: not enough players")
	}
	
	room.Game = model.NewPVPGame(roomID)
	room.Status = "playing"
	
	return nil
}

// LeaveRoom removes a player from a room
func (gs *GameService) LeaveRoom(roomID, playerID string) error {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	
	room, exists := gs.rooms[roomID]
	if !exists {
		return fmt.Errorf("room not found")
	}
	
	room.RemovePlayer(playerID)
	
	// If room is empty, delete it
	if len(room.Players) == 0 {
		delete(gs.rooms, roomID)
	}
	
	return nil
}

// GetRoomStatus returns the current status of a room
func (gs *GameService) GetRoomStatus(roomID string) (string, error) {
	room := gs.GetRoom(roomID)
	if room == nil {
		return "", fmt.Errorf("room not found")
	}
	
	return room.Status, nil
}

// SetPlayerReady sets a player's ready status
func (gs *GameService) SetPlayerReady(roomID, playerID string, ready bool) error {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	
	room, exists := gs.rooms[roomID]
	if !exists {
		return fmt.Errorf("room not found")
	}
	
	player := room.GetPlayer(playerID)
	if player == nil {
		return fmt.Errorf("player not found in room")
	}
	
	player.IsReady = ready
	room.UpdatedAt = time.Now()
	
	return nil
}

// HandlePlayerDisconnect handles when a player disconnects
func (gs *GameService) HandlePlayerDisconnect(roomID, playerID string) error {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	
	room, exists := gs.rooms[roomID]
	if !exists {
		return fmt.Errorf("room not found")
	}
	
	player := room.GetPlayer(playerID)
	if player == nil {
		return fmt.Errorf("player not found in room")
	}
	
	player.IsOnline = false
	room.UpdatedAt = time.Now()
	
	// If game is in progress and player disconnects, end the game
	if room.Game != nil && room.Game.Status == "playing" {
		room.Game.Status = "finished"
		// Set the other player as winner
		for _, p := range room.Players {
			if p.ID != playerID {
				room.Game.Winner = p.PlayerNumber
				break
			}
		}
		now := time.Now()
		room.Game.EndedAt = &now
	}
	
	return nil
}