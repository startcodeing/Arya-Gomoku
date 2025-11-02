package model

import (
	"time"
	"github.com/google/uuid"
)

// Room represents a game room for PVP matches
type Room struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Status      string    `json:"status"` // waiting, playing, finished
	MaxPlayers  int       `json:"maxPlayers"`
	Players     []*PVPPlayer `json:"players"`
	Game        *PVPGameSession  `json:"game,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	CreatorID   string    `json:"creatorId"`
}

// PVPPlayer represents a player in PVP mode
type PVPPlayer struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	RoomID       string    `json:"roomId"`
	PlayerNumber int       `json:"playerNumber"` // 1 or 2
	IsReady      bool      `json:"isReady"`
	IsOnline     bool      `json:"isOnline"`
	JoinedAt     time.Time `json:"joinedAt"`
	IsCreator    bool      `json:"isCreator"`
}

// PVPGameSession represents a PVP game session
type PVPGameSession struct {
	ID            string     `json:"id"`
	RoomID        string     `json:"roomId"`
	Status        string     `json:"status"`
	Board         [][]int    `json:"board"`
	CurrentPlayer string     `json:"currentPlayer"`
	Winner        string     `json:"winner"`
	MoveCount     int        `json:"moveCount"`
	Moves         []*PVPMove `json:"moves"`
	StartedAt     time.Time  `json:"startedAt"`
	EndedAt       *time.Time `json:"endedAt,omitempty"`
}

// PVPMove represents a move in PVP game
type PVPMove struct {
	ID         string    `json:"id"`
	GameID     string    `json:"gameId"`
	PlayerID   string    `json:"playerId"`
	X          int       `json:"x"`
	Y          int       `json:"y"`
	MoveNumber int       `json:"moveNumber"`
	CreatedAt  time.Time `json:"createdAt"`
}

// WSMessage represents WebSocket message structure
type WSMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// RoomUpdateData represents room update message data
type RoomUpdateData struct {
	Room   *Room      `json:"room"`
	Player *PVPPlayer `json:"player,omitempty"`
}

// GameUpdateData represents game update message data
type GameUpdateData struct {
	Game     *PVPGameSession `json:"game"`
	LastMove *PVPMove `json:"lastMove,omitempty"`
}

// ChatMessageData represents chat message data
type ChatMessageData struct {
	PlayerID   string    `json:"playerId"`
	PlayerName string    `json:"playerName"`
	Message    string    `json:"message"`
	Timestamp  time.Time `json:"timestamp"`
}

// CreateRoomRequest represents request to create a room
type CreateRoomRequest struct {
	RoomName   string `json:"roomName" binding:"required"`
	PlayerName string `json:"playerName" binding:"required"`
	MaxPlayers int    `json:"maxPlayers"`
}

// JoinRoomRequest represents request to join a room
type JoinRoomRequest struct {
	PlayerName string `json:"playerName" binding:"required"`
}

// MakeMoveRequest represents request to make a move
type MakeMoveRequest struct {
	X        int    `json:"x" binding:"required"`
	Y        int    `json:"y" binding:"required"`
	PlayerID string `json:"playerId" binding:"required"`
}

// NewRoom creates a new room
func NewRoom(name, creatorName string, maxPlayers int) *Room {
	roomID := uuid.New().String()
	playerID := uuid.New().String()
	
	creator := &PVPPlayer{
		ID:           playerID,
		Name:         creatorName,
		RoomID:       roomID,
		PlayerNumber: 1,
		IsReady:      false,
		IsOnline:     true,
		JoinedAt:     time.Now(),
		IsCreator:    true,
	}
	
	if maxPlayers <= 0 {
		maxPlayers = 2
	}
	
	return &Room{
		ID:         roomID,
		Name:       name,
		Status:     "waiting",
		MaxPlayers: maxPlayers,
		Players:    []*PVPPlayer{creator},
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		CreatorID:  playerID,
	}
}

// NewPVPGameSession creates a new PVP game
func NewPVPGameSession(room *Room) *PVPGameSession {
	gameID := uuid.New().String()
	
	// Initialize 15x15 board
	board := make([][]int, 15)
	for i := range board {
		board[i] = make([]int, 15)
	}
	
	// Set the first player (PlayerNumber 1) as the starting player
	var firstPlayerID string
	if len(room.Players) > 0 {
		// Find the player with PlayerNumber 1 (first player gets black pieces)
		for _, player := range room.Players {
			if player.PlayerNumber == 1 {
				firstPlayerID = player.ID
				break
			}
		}
	}
	
	return &PVPGameSession{
		ID:            gameID,
		RoomID:        room.ID,
		Status:        "playing",
		Board:         board,
		CurrentPlayer: firstPlayerID, // First player (black) starts first
		Winner:        "",
		MoveCount:     0,
		Moves:         []*PVPMove{},
		StartedAt:     time.Now(),
	}
}

// AddPlayer adds a player to the room
func (r *Room) AddPlayer(playerName string) *PVPPlayer {
	if len(r.Players) >= r.MaxPlayers {
		return nil
	}
	
	playerID := uuid.New().String()
	playerNumber := len(r.Players) + 1
	
	player := &PVPPlayer{
		ID:           playerID,
		Name:         playerName,
		RoomID:       r.ID,
		PlayerNumber: playerNumber,
		IsReady:      false,
		IsOnline:     true,
		JoinedAt:     time.Now(),
		IsCreator:    false,
	}
	
	r.Players = append(r.Players, player)
	r.UpdatedAt = time.Now()
	
	return player
}

// RemovePlayer removes a player from the room
func (r *Room) RemovePlayer(playerID string) bool {
	for i, player := range r.Players {
		if player.ID == playerID {
			r.Players = append(r.Players[:i], r.Players[i+1:]...)
			r.UpdatedAt = time.Now()
			return true
		}
	}
	return false
}

// GetPlayer gets a player by ID
func (r *Room) GetPlayer(playerID string) *PVPPlayer {
	for _, player := range r.Players {
		if player.ID == playerID {
			return player
		}
	}
	return nil
}

// CanStartGame checks if the game can be started
func (r *Room) CanStartGame() bool {
	if len(r.Players) < 2 {
		return false
	}
	
	for _, player := range r.Players {
		if !player.IsReady {
			return false
		}
	}
	
	return true
}

// IsValidMove checks if a move is valid
func (g *PVPGameSession) IsValidMove(x, y int) bool {
	if x < 0 || x >= 15 || y < 0 || y >= 15 {
		return false
	}
	return g.Board[y][x] == 0
}

// MakeMove makes a move in the game
func (g *PVPGameSession) MakeMove(x, y int, playerID string, playerNumber int, room *Room) *PVPMove {
	if !g.IsValidMove(x, y) || g.CurrentPlayer != playerID {
		return nil
	}
	
	moveID := uuid.New().String()
	move := &PVPMove{
		ID:         moveID,
		GameID:     g.ID,
		PlayerID:   playerID,
		X:          x,
		Y:          y,
		MoveNumber: g.MoveCount + 1,
		CreatedAt:  time.Now(),
	}
	
	g.Board[y][x] = playerNumber
	g.MoveCount++
	g.Moves = append(g.Moves, move)
	
	// Check for win
	if g.CheckWin(x, y, playerNumber) {
		g.Winner = playerID
		g.Status = "finished"
		now := time.Now()
		g.EndedAt = &now
	} else if g.IsBoardFull() {
		g.Status = "finished"
		now := time.Now()
		g.EndedAt = &now
	} else {
		// Switch to the other player
		for _, player := range room.Players {
			if player.ID != playerID {
				g.CurrentPlayer = player.ID
				break
			}
		}
	}
	
	return move
}

// CheckWin checks if the specified player has won
func (g *PVPGameSession) CheckWin(x, y, player int) bool {
	directions := [][]int{
		{1, 0},   // Horizontal
		{0, 1},   // Vertical
		{1, 1},   // Diagonal \
		{1, -1},  // Diagonal /
	}

	for _, dir := range directions {
		count := 1 // Count the current piece
		
		// Check in positive direction
		for i := 1; i < 5; i++ {
			nx, ny := x+dir[0]*i, y+dir[1]*i
			if nx < 0 || nx >= 15 || ny < 0 || ny >= 15 || g.Board[ny][nx] != player {
				break
			}
			count++
		}
		
		// Check in negative direction
		for i := 1; i < 5; i++ {
			nx, ny := x-dir[0]*i, y-dir[1]*i
			if nx < 0 || nx >= 15 || ny < 0 || ny >= 15 || g.Board[ny][nx] != player {
				break
			}
			count++
		}
		
		if count >= 5 {
			return true
		}
	}
	
	return false
}

// IsBoardFull checks if the board is full
func (g *PVPGameSession) IsBoardFull() bool {
	return g.MoveCount >= 15*15
}