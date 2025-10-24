// Package model defines the core data structures for the Gomoku game
// This package contains the game board, player information, and game state models
package model

import "time"

// Board represents the game board state and configuration
type Board struct {
	Grid          [][]int `json:"grid"`          // 15x15 grid: 0=empty, 1=player, 2=AI
	Size          int     `json:"size"`          // Board size (default: 15)
	CurrentPlayer int     `json:"currentPlayer"` // Current player: 1=human, 2=AI
	MoveCount     int     `json:"moveCount"`     // Number of moves made
}

// Player represents a game player (human or AI)
type Player struct {
	ID   int    `json:"id"`   // Player ID: 1=human, 2=AI
	Type string `json:"type"` // Player type: "human" or "ai"
	Name string `json:"name"` // Player display name
}

// Move represents a single move on the board
type Move struct {
	X      int `json:"x"`      // X coordinate (0-14)
	Y      int `json:"y"`      // Y coordinate (0-14)
	Player int `json:"player"` // Player who made the move
}

// GameState represents the current state of the game
type GameState struct {
	Status    string `json:"status"`    // Game status: "playing", "win", "lose", "draw"
	Winner    int    `json:"winner"`    // Winner: 0=none, 1=player, 2=AI
	MoveCount int    `json:"moveCount"` // Total moves made
	LastMove  *Move  `json:"lastMove"`  // Last move made
}

// AIMove represents an AI's move decision
type AIMove struct {
	X     int `json:"x"`     // X coordinate
	Y     int `json:"y"`     // Y coordinate
	Score int `json:"score"` // Move evaluation score
}

// GameRequest represents the request payload for AI move
type GameRequest struct {
	Board    [][]int `json:"board"`    // Current board state
	Player   int     `json:"player"`   // Current player
	LastMove Move    `json:"lastMove"` // Last move made
}

// GameResponse represents the response from AI move endpoint
type GameResponse struct {
	AIMove     AIMove `json:"aiMove"`     // AI's chosen move
	GameStatus string `json:"gameStatus"` // Updated game status
	Winner     int    `json:"winner"`     // Game winner if any
}

// MatchRoom represents an online match room (reserved for future PVP feature)
type MatchRoom struct {
	RoomID    string    `json:"roomId"`    // Unique room identifier
	Players   []Player  `json:"players"`   // Players in the room
	Board     *Board    `json:"board"`     // Current board state
	Status    string    `json:"status"`    // Room status: "waiting", "playing", "finished"
	CreatedAt time.Time `json:"createdAt"` // Room creation timestamp
}

// NewBoard creates and initializes a new empty game board
func NewBoard() *Board {
	grid := make([][]int, 15)
	for i := range grid {
		grid[i] = make([]int, 15)
	}
	return &Board{
		Grid:          grid,
		Size:          15,
		CurrentPlayer: 1, // Human player starts first
		MoveCount:     0,
	}
}

// IsValidMove checks if a move is valid on the current board
func (b *Board) IsValidMove(x, y int) bool {
	if x < 0 || x >= b.Size || y < 0 || y >= b.Size {
		return false
	}
	return b.Grid[y][x] == 0
}

// MakeMove places a piece on the board
func (b *Board) MakeMove(x, y, player int) bool {
	if !b.IsValidMove(x, y) {
		return false
	}
	b.Grid[y][x] = player
	b.MoveCount++
	b.CurrentPlayer = 3 - player // Switch between 1 and 2
	return true
}

// CheckWin checks if the specified player has won the game
func (b *Board) CheckWin(x, y, player int) bool {
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
			if nx < 0 || nx >= b.Size || ny < 0 || ny >= b.Size || b.Grid[ny][nx] != player {
				break
			}
			count++
		}
		
		// Check in negative direction
		for i := 1; i < 5; i++ {
			nx, ny := x-dir[0]*i, y-dir[1]*i
			if nx < 0 || nx >= b.Size || ny < 0 || ny >= b.Size || b.Grid[ny][nx] != player {
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

// IsBoardFull checks if the board is completely filled
func (b *Board) IsBoardFull() bool {
	return b.MoveCount >= b.Size*b.Size
}

// GetGameState returns the current game state
func (b *Board) GetGameState(lastMove *Move) GameState {
	if lastMove != nil && b.CheckWin(lastMove.X, lastMove.Y, lastMove.Player) {
		status := "win"
		if lastMove.Player == 2 {
			status = "lose"
		}
		return GameState{
			Status:    status,
			Winner:    lastMove.Player,
			MoveCount: b.MoveCount,
			LastMove:  lastMove,
		}
	}
	
	if b.IsBoardFull() {
		return GameState{
			Status:    "draw",
			Winner:    0,
			MoveCount: b.MoveCount,
			LastMove:  lastMove,
		}
	}
	
	return GameState{
		Status:    "playing",
		Winner:    0,
		MoveCount: b.MoveCount,
		LastMove:  lastMove,
	}
}