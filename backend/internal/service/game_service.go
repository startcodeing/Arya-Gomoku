// Package service contains the business logic for the Gomoku game
// This file implements game room management and online match services for PVP feature
package service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"gomoku-backend/internal/model"
	"gomoku-backend/internal/repository"
)

// GameService manages game rooms and online matches for PVP feature
type GameService struct {
	rooms    map[string]*model.Room
	mutex    sync.RWMutex
	gameRepo repository.GameRepository
}

// NewGameService creates a new game service
func NewGameService(gameRepo repository.GameRepository) *GameService {
	return &GameService{
		rooms:    make(map[string]*model.Room),
		gameRepo: gameRepo,
	}
}

// CreateRoom creates a new match room for PVP feature
func (gs *GameService) CreateRoom(roomName, playerName string, maxPlayers int) (*model.Room, error) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	
	log.Printf("开始创建房间: roomName=%s, playerName=%s, maxPlayers=%d", roomName, playerName, maxPlayers)
	
	room := model.NewRoom(roomName, playerName, maxPlayers)
	log.Printf("房间创建成功: roomID=%s, roomName=%s", room.ID, room.Name)
	
	gs.rooms[room.ID] = room
	log.Printf("房间已存储到内存，当前房间总数: %d", len(gs.rooms))
	
	// 验证房间是否正确存储
	if storedRoom := gs.rooms[room.ID]; storedRoom != nil {
		log.Printf("验证成功: 房间 %s 已正确存储", room.ID)
	} else {
		log.Printf("警告: 房间 %s 存储失败", room.ID)
	}
	
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
	
	log.Printf("尝试获取房间: roomID=%s", roomID)
	log.Printf("当前内存中的房间总数: %d", len(gs.rooms))
	
	// 列出所有房间ID用于调试
	if len(gs.rooms) > 0 {
		log.Printf("当前存储的房间ID列表:")
		for id, room := range gs.rooms {
			log.Printf("  - 房间ID: %s, 房间名: %s, 状态: %s", id, room.Name, room.Status)
		}
	}
	
	room := gs.rooms[roomID]
	if room != nil {
		log.Printf("房间获取成功: roomID=%s, roomName=%s, status=%s", room.ID, room.Name, room.Status)
	} else {
		log.Printf("房间未找到: roomID=%s", roomID)
	}
	
	return room
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
	if room.Game.CurrentPlayer != playerID {
		return nil, nil, fmt.Errorf("not your turn")
	}
	
	// Make the move
	move := room.Game.MakeMove(x, y, playerID, player.PlayerNumber, room)
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
	
	room.Game = model.NewPVPGameSession(room)
	room.Status = "playing"
	
	return nil
}

// SaveGameResult 保存游戏结果到数据库
func (gs *GameService) SaveGameResult(room *model.Room) error {
	if room.Game == nil {
		return fmt.Errorf("游戏会话不存在")
	}

	// 创建PVP游戏记录
	game := &model.PVPGame{
		RoomID:    room.ID,
		Status:    model.GameStatusCompleted,
		BoardSize: 15, // 默认棋盘大小
		StartedAt: &room.Game.StartedAt,
		EndedAt:   room.Game.EndedAt,
		MoveCount: room.Game.MoveCount,
	}

	// 设置胜者
	if room.Game.Winner != "" {
		// 查找获胜玩家的UUID
		for _, player := range room.Players {
			if player.ID == room.Game.Winner {
				// 这里需要将string ID转换为UUID，但目前Player使用的是string ID
				// 暂时跳过WinnerID设置，因为需要用户系统的UUID
				break
			}
		}
	}

	// 保存到数据库
	return gs.gameRepo.CreatePVPGame(context.Background(), game)
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
	
	// If game is in progress and player disconnects, end the game
	if room.Game != nil && room.Game.Status == "playing" {
		room.Game.Status = "finished"
		// Set the other player as winner
		for _, p := range room.Players {
			if p.ID != playerID {
				room.Game.Winner = p.ID
				break
			}
		}
		now := time.Now()
		room.Game.EndedAt = &now

		// 保存游戏结果到数据库
		if err := gs.SaveGameResult(room); err != nil {
			log.Printf("Error saving game result on disconnect: %v", err)
		}
	}
	
	// Remove player from room to allow reuse
	room.RemovePlayer(playerID)
	room.UpdatedAt = time.Now()
	
	// If room is empty, delete it
	if len(room.Players) == 0 {
		delete(gs.rooms, roomID)
	} else {
		// Reset room status to waiting if there are remaining players
		room.Status = "waiting"
		// Reset game if it exists
		room.Game = nil
	}
	
	return nil
}