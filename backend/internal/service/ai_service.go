// Package service contains the business logic for the Gomoku game
// This package implements the AI algorithms and game management services
package service

import (
	"math/rand"
	"time"

	"gomoku-backend/internal/model"
)

// AIService handles AI move generation and game logic
type AIService struct {
	rand *rand.Rand
}

// NewAIService creates a new AI service instance
func NewAIService() *AIService {
	return &AIService{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GetAIMove generates the best move for the AI using heuristic algorithm
// Priority: Win > Block opponent's four-in-a-row > Create three-in-a-row > Random
func (ai *AIService) GetAIMove(board [][]int, lastMove model.Move) model.AIMove {
	size := len(board)
	
	// Priority 1: Check if AI can win in one move
	if move := ai.findWinningMove(board, 2, size); move != nil {
		return *move
	}
	
	// Priority 2: Block opponent's winning move
	if move := ai.findWinningMove(board, 1, size); move != nil {
		return *move
	}
	
	// Priority 3: Create threats (three-in-a-row)
	if move := ai.findThreateningMove(board, 2, size); move != nil {
		return *move
	}
	
	// Priority 4: Block opponent's threats
	if move := ai.findThreateningMove(board, 1, size); move != nil {
		return *move
	}
	
	// Priority 5: Strategic positioning near last move
	if move := ai.findStrategicMove(board, lastMove, size); move != nil {
		return *move
	}
	
	// Fallback: Random valid move
	return ai.findRandomMove(board, size)
}

// findWinningMove looks for a move that creates five-in-a-row for the specified player
func (ai *AIService) findWinningMove(board [][]int, player, size int) *model.AIMove {
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if board[y][x] == 0 {
				// Temporarily place the piece
				board[y][x] = player
				if ai.checkWin(board, x, y, player, size) {
					board[y][x] = 0 // Restore
					return &model.AIMove{X: x, Y: y, Score: 1000}
				}
				board[y][x] = 0 // Restore
			}
		}
	}
	return nil
}

// findThreateningMove looks for a move that creates three-in-a-row with potential to win
func (ai *AIService) findThreateningMove(board [][]int, player, size int) *model.AIMove {
	bestMove := &model.AIMove{Score: -1}
	
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if board[y][x] == 0 {
				score := ai.evaluatePosition(board, x, y, player, size)
				if score > bestMove.Score && score >= 50 { // Threshold for threatening moves
					bestMove.X = x
					bestMove.Y = y
					bestMove.Score = score
				}
			}
		}
	}
	
	if bestMove.Score > -1 {
		return bestMove
	}
	return nil
}

// findStrategicMove finds a good move near the last opponent move
func (ai *AIService) findStrategicMove(board [][]int, lastMove model.Move, size int) *model.AIMove {
	// Search in a 3x3 area around the last move
	directions := [][]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1},           {0, 1},
		{1, -1},  {1, 0},  {1, 1},
	}
	
	bestMove := &model.AIMove{Score: -1}
	
	for _, dir := range directions {
		x := lastMove.X + dir[0]
		y := lastMove.Y + dir[1]
		
		if x >= 0 && x < size && y >= 0 && y < size && board[y][x] == 0 {
			score := ai.evaluatePosition(board, x, y, 2, size)
			if score > bestMove.Score {
				bestMove.X = x
				bestMove.Y = y
				bestMove.Score = score
			}
		}
	}
	
	if bestMove.Score > -1 {
		return bestMove
	}
	return nil
}

// findRandomMove finds a random valid move as fallback
func (ai *AIService) findRandomMove(board [][]int, size int) model.AIMove {
	var validMoves []model.AIMove
	
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if board[y][x] == 0 {
				validMoves = append(validMoves, model.AIMove{X: x, Y: y, Score: 1})
			}
		}
	}
	
	if len(validMoves) > 0 {
		return validMoves[ai.rand.Intn(len(validMoves))]
	}
	
	// Should never reach here in a valid game
	return model.AIMove{X: 7, Y: 7, Score: 1}
}

// evaluatePosition calculates the strategic value of a position
func (ai *AIService) evaluatePosition(board [][]int, x, y, player, size int) int {
	directions := [][]int{
		{1, 0},   // Horizontal
		{0, 1},   // Vertical
		{1, 1},   // Diagonal \
		{1, -1},  // Diagonal /
	}
	
	totalScore := 0
	
	for _, dir := range directions {
		score := ai.evaluateDirection(board, x, y, dir[0], dir[1], player, size)
		totalScore += score
	}
	
	return totalScore
}

// evaluateDirection evaluates the potential in a specific direction
func (ai *AIService) evaluateDirection(board [][]int, x, y, dx, dy, player, size int) int {
	count := 1 // Count the piece we're placing
	openEnds := 0
	
	// Check positive direction
	for i := 1; i < 5; i++ {
		nx, ny := x+dx*i, y+dy*i
		if nx < 0 || nx >= size || ny < 0 || ny >= size {
			break
		}
		if board[ny][nx] == player {
			count++
		} else if board[ny][nx] == 0 {
			openEnds++
			break
		} else {
			break
		}
	}
	
	// Check negative direction
	for i := 1; i < 5; i++ {
		nx, ny := x-dx*i, y-dy*i
		if nx < 0 || nx >= size || ny < 0 || ny >= size {
			break
		}
		if board[ny][nx] == player {
			count++
		} else if board[ny][nx] == 0 {
			openEnds++
			break
		} else {
			break
		}
	}
	
	// Score based on count and open ends
	switch count {
	case 4:
		return 100 // Four in a row - very strong
	case 3:
		if openEnds >= 1 {
			return 50 // Three in a row with open end
		}
		return 10
	case 2:
		if openEnds >= 2 {
			return 20 // Two in a row with both ends open
		}
		if openEnds >= 1 {
			return 5
		}
		return 2
	default:
		return 1
	}
}

// checkWin checks if placing a piece at (x,y) creates five-in-a-row
func (ai *AIService) checkWin(board [][]int, x, y, player, size int) bool {
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
			if nx < 0 || nx >= size || ny < 0 || ny >= size || board[ny][nx] != player {
				break
			}
			count++
		}
		
		// Check in negative direction
		for i := 1; i < 5; i++ {
			nx, ny := x-dir[0]*i, y-dir[1]*i
			if nx < 0 || nx >= size || ny < 0 || ny >= size || board[ny][nx] != player {
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