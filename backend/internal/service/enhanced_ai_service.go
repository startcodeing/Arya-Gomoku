// Package service contains enhanced AI algorithms for the Gomoku game
// This package implements minimax algorithm with alpha-beta pruning and transposition table
package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"gomoku-backend/internal/model"
)

// Difficulty levels for AI
type Difficulty int

const (
	Easy Difficulty = iota
	Medium
	Hard
	Expert
)

// TranspositionTableEntry stores evaluated board positions
type TranspositionTableEntry struct {
	Score    int
	Depth    int
	Flag     EntryFlag
	BestMove model.Move
}

// EntryFlag indicates the type of score stored
type EntryFlag int

const (
	Exact EntryFlag = iota
	LowerBound
	UpperBound
)

// EnhancedAIService handles AI move generation with advanced algorithms
type EnhancedAIService struct {
	boardSize       int
	maxDepth        map[Difficulty]int
	transpositionTable map[string]*TranspositionTableEntry
	tableMutex      sync.RWMutex
	searchStartTime time.Time
	timeLimit       time.Duration
	nodesSearched   uint64
	cutoffs         uint64
}

// NewEnhancedAIService creates a new enhanced AI service instance
func NewEnhancedAIService() *EnhancedAIService {
	return &EnhancedAIService{
		boardSize:    15,
		maxDepth: map[Difficulty]int{
			Easy:   2,  // Quick look-ahead
			Medium: 4,  // Moderate depth
			Hard:   6,  // Deep analysis
			Expert: 8,  // Maximum depth with time limit
		},
		transpositionTable: make(map[string]*TranspositionTableEntry),
		timeLimit:          5 * time.Second, // 5 second thinking time
	}
}

// GetAIMove generates the best move for the AI using minimax with alpha-beta pruning
func (ai *EnhancedAIService) GetAIMove(board [][]int, lastMove model.Move, difficulty Difficulty) model.AIMove {
	ai.searchStartTime = time.Now()
	ai.nodesSearched = 0
	ai.cutoffs = 0

	// Get available moves
	moves := ai.getAvailableMoves(board, lastMove)
	if len(moves) == 0 {
		return model.AIMove{X: 7, Y: 7, Score: -1}
	}

	// For easy difficulty, use simple heuristic
	if difficulty == Easy {
		return ai.getHeuristicMove(board, lastMove, moves)
	}

	// For medium and above, use minimax with alpha-beta
	maxDepth := ai.maxDepth[difficulty]

	var bestMove model.Move
	var bestScore int = math.MinInt32

	// Iterative deepening with time limit
	for depth := 1; depth <= maxDepth; depth++ {
		if time.Since(ai.searchStartTime) > ai.timeLimit {
			break
		}

		score, move := ai.minimax(board, depth, math.MinInt32, math.MaxInt32, true, lastMove)

		if score > bestScore {
			bestScore = score
			bestMove = move
		}

		// If we found a winning move, no need to search deeper
		if score >= 100000 {
			break
		}
	}

	// Fallback to center if no move found
	if bestMove.X == -1 && len(moves) > 0 {
		bestMove = moves[0]
	}

	return model.AIMove{
		X:     bestMove.X,
		Y:     bestMove.Y,
		Score: bestScore,
	}
}

// minimax implements the minimax algorithm with alpha-beta pruning
func (ai *EnhancedAIService) minimax(board [][]int, depth, alpha, beta int, isMaximizing bool, lastMove model.Move) (int, model.Move) {
	ai.nodesSearched++

	// Check time limit
	if time.Since(ai.searchStartTime) > ai.timeLimit {
		return 0, model.Move{X: -1, Y: -1}
	}

	// Check terminal conditions
	if score := ai.evaluateTerminal(board, lastMove); score != math.MinInt32 {
		return score, model.Move{X: -1, Y: -1}
	}

	// Check depth limit
	if depth == 0 {
		return ai.evaluateBoard(board, lastMove), model.Move{X: -1, Y: -1}
	}

	// Check transposition table
	boardHash := ai.hashBoard(board)
	ai.tableMutex.RLock()
	if entry, exists := ai.transpositionTable[boardHash]; exists && entry.Depth >= depth {
		ai.tableMutex.RUnlock()

		switch entry.Flag {
		case Exact:
			return entry.Score, entry.BestMove
		case LowerBound:
			if entry.Score > alpha {
				alpha = entry.Score
			}
		case UpperBound:
			if entry.Score < beta {
				beta = entry.Score
			}
		}

		if alpha >= beta {
			ai.cutoffs++
			return entry.Score, entry.BestMove
		}
	} else {
		ai.tableMutex.RUnlock()
	}

	moves := ai.getAvailableMoves(board, lastMove)
	if len(moves) == 0 {
		return 0, model.Move{X: -1, Y: -1}
	}

	var bestMove model.Move
	var bestScore int

	if isMaximizing {
		bestScore = math.MinInt32
		for _, move := range moves {
			// Make move
			board[move.Y][move.X] = 2 // AI piece

			score, _ := ai.minimax(board, depth-1, alpha, beta, false, move)

			// Undo move
			board[move.Y][move.X] = 0

			if score > bestScore {
				bestScore = score
				bestMove = move
			}

			alpha = max(alpha, bestScore)
			if beta <= alpha {
				ai.cutoffs++
				break // Beta cutoff
			}
		}
	} else {
		bestScore = math.MaxInt32
		for _, move := range moves {
			// Make move
			board[move.Y][move.X] = 1 // Human piece

			score, _ := ai.minimax(board, depth-1, alpha, beta, true, move)

			// Undo move
			board[move.Y][move.X] = 0

			if score < bestScore {
				bestScore = score
				bestMove = move
			}

			beta = min(beta, bestScore)
			if beta <= alpha {
				ai.cutoffs++
				break // Alpha cutoff
			}
		}
	}

	// Store result in transposition table
	flag := Exact
	if bestScore <= alpha {
		flag = UpperBound
	} else if bestScore >= beta {
		flag = LowerBound
	}

	ai.tableMutex.Lock()
	ai.transpositionTable[boardHash] = &TranspositionTableEntry{
		Score:    bestScore,
		Depth:    depth,
		Flag:     flag,
		BestMove: bestMove,
	}
	ai.tableMutex.Unlock()

	return bestScore, bestMove
}

// evaluateTerminal checks for terminal game states (win/loss/draw)
func (ai *EnhancedAIService) evaluateTerminal(board [][]int, lastMove model.Move) int {
	// Check if last move created a win
	if lastMove.X >= 0 && lastMove.Y >= 0 {
		player := board[lastMove.Y][lastMove.X]
		if player != 0 && ai.checkWin(board, lastMove.X, lastMove.Y, player) {
			if player == 2 { // AI wins
				return 100000
			} else { // Human wins
				return -100000
			}
		}
	}

	// Check for draw
	if ai.isBoardFull(board) {
		return 0
	}

	return math.MinInt32 // Not a terminal state
}

// evaluateBoard provides a comprehensive evaluation of the board state
func (ai *EnhancedAIService) evaluateBoard(board [][]int, lastMove model.Move) int {
	score := 0

	// Evaluate all positions
	for y := 0; y < ai.boardSize; y++ {
		for x := 0; x < ai.boardSize; x++ {
			if board[y][x] != 0 {
				pieceScore := ai.evaluatePositionAdvanced(board, x, y, board[y][x])
				if board[y][x] == 2 { // AI piece
					score += pieceScore
				} else { // Human piece
					score -= pieceScore
				}
			}
		}
	}

	// Add positional bonuses
	centerX, centerY := ai.boardSize/2, ai.boardSize/2
	if lastMove.X >= 0 && lastMove.Y >= 0 {
		distanceToCenter := abs(centerX-lastMove.X) + abs(centerY-lastMove.Y)
		score += (7 - distanceToCenter) * 2 // Prefer center positions
	}

	return score
}

// evaluatePositionAdvanced provides advanced position evaluation
func (ai *EnhancedAIService) evaluatePositionAdvanced(board [][]int, x, y, player int) int {
	directions := [][]int{
		{1, 0},   // Horizontal
		{0, 1},   // Vertical
		{1, 1},   // Diagonal \
		{1, -1},  // Diagonal /
	}

	totalScore := 0
	opponent := 3 - player

	for _, dir := range directions {
		pattern := ai.getPattern(board, x, y, dir[0], dir[1], player)
		score := ai.evaluatePattern(pattern, player, opponent)
		totalScore += score
	}

	return totalScore
}

// getPattern extracts the pattern around a position in a specific direction
func (ai *EnhancedAIService) getPattern(board [][]int, x, y, dx, dy, player int) string {
	pattern := make([]rune, 9) // 4 spaces + piece + 4 spaces

	// Center is current position
	pattern[4] = rune(player + '0')

	// Build pattern in both directions
	for i := 1; i <= 4; i++ {
		// Positive direction
		px, py := x+dx*i, y+dy*i
		if px >= 0 && px < ai.boardSize && py >= 0 && py < ai.boardSize {
			pattern[4+i] = rune(board[py][px] + '0')
		} else {
			pattern[4+i] = 'X' // Out of bounds
		}

		// Negative direction
		nx, ny := x-dx*i, y-dy*i
		if nx >= 0 && nx < ai.boardSize && ny >= 0 && ny < ai.boardSize {
			pattern[4-i] = rune(board[ny][nx] + '0')
		} else {
			pattern[4-i] = 'X' // Out of bounds
		}
	}

	return string(pattern)
}

// evaluatePattern evaluates a pattern and returns a score
func (ai *EnhancedAIService) evaluatePattern(pattern string, player, opponent int) int {
	// Winning patterns
	if strings.Contains(pattern, fmt.Sprintf("%d%d%d%d%d", player, player, player, player, player)) {
		return 100000
	}

	// Four in a row with open end
	if strings.Contains(pattern, fmt.Sprintf("0%d%d%d%d0", player, player, player, player)) ||
	   strings.Contains(pattern, fmt.Sprintf("%d%d%d%d0", player, player, player, player)) ||
	   strings.Contains(pattern, fmt.Sprintf("0%d%d%d%d", player, player, player, player)) {
		return 10000
	}

	// Three in a row with open ends
	if strings.Contains(pattern, fmt.Sprintf("0%d%d%d0", player, player, player)) {
		return 1000
	}

	// Three in a row with one open end
	if strings.Contains(pattern, fmt.Sprintf("%d%d%d0", player, player, player)) ||
	   strings.Contains(pattern, fmt.Sprintf("0%d%d%d", player, player, player)) {
		return 100
	}

	// Two in a row with potential
	if strings.Contains(pattern, fmt.Sprintf("0%d%d0", player, player)) {
		return 50
	}

	// Blocking opponent's patterns
	opponentStr := fmt.Sprintf("%d", opponent)
	if strings.Contains(pattern, fmt.Sprintf("0%s%s%s%s0", opponentStr, opponentStr, opponentStr, opponentStr)) {
		return 5000 // Block opponent's four
	}

	return 10
}

// getAvailableMoves gets reasonable moves to consider (pruning the search space)
func (ai *EnhancedAIService) getAvailableMoves(board [][]int, lastMove model.Move) []model.Move {
	var moves []model.Move
	visited := make(map[[2]int]bool)

	// Consider positions near existing pieces
	searchRadius := 2

	for y := 0; y < ai.boardSize; y++ {
		for x := 0; x < ai.boardSize; x++ {
			if board[y][x] != 0 {
				// Check positions around this piece
				for dy := -searchRadius; dy <= searchRadius; dy++ {
					for dx := -searchRadius; dx <= searchRadius; dx++ {
						nx, ny := x+dx, y+dy
						key := [2]int{nx, ny}

						if nx >= 0 && nx < ai.boardSize && ny >= 0 && ny < ai.boardSize &&
						   board[ny][nx] == 0 && !visited[key] {
							moves = append(moves, model.Move{X: nx, Y: ny})
							visited[key] = true
						}
					}
				}
			}
		}
	}

	// If no moves near existing pieces (likely empty board), add center positions
	if len(moves) == 0 {
		center := ai.boardSize / 2
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				nx, ny := center+dx, center+dy
				if nx >= 0 && nx < ai.boardSize && ny >= 0 && ny < ai.boardSize &&
				   board[ny][nx] == 0 {
					moves = append(moves, model.Move{X: nx, Y: ny})
				}
			}
		}
	}

	return moves
}

// getHeuristicMove provides a simple heuristic move for easy difficulty
func (ai *EnhancedAIService) getHeuristicMove(board [][]int, lastMove model.Move, moves []model.Move) model.AIMove {
	// Priority 1: Check if AI can win
	for _, move := range moves {
		board[move.Y][move.X] = 2
		if ai.checkWin(board, move.X, move.Y, 2) {
			board[move.Y][move.X] = 0
			return model.AIMove{X: move.X, Y: move.Y, Score: 1000}
		}
		board[move.Y][move.X] = 0
	}

	// Priority 2: Block opponent's win
	for _, move := range moves {
		board[move.Y][move.X] = 1
		if ai.checkWin(board, move.X, move.Y, 1) {
			board[move.Y][move.X] = 0
			return model.AIMove{X: move.X, Y: move.Y, Score: 900}
		}
		board[move.Y][move.X] = 0
	}

	// Priority 3: Best positional score
	bestScore := math.MinInt32
	var bestMove model.Move

	for _, move := range moves {
		score := ai.evaluatePositionAdvanced(board, move.X, move.Y, 2)
		if score > bestScore {
			bestScore = score
			bestMove = move
		}
	}

	return model.AIMove{X: bestMove.X, Y: bestMove.Y, Score: bestScore}
}

// checkWin checks if a player has won
func (ai *EnhancedAIService) checkWin(board [][]int, x, y, player int) bool {
	directions := [][]int{
		{1, 0},   // Horizontal
		{0, 1},   // Vertical
		{1, 1},   // Diagonal \
		{1, -1},  // Diagonal /
	}

	for _, dir := range directions {
		count := 1

		// Check positive direction
		for i := 1; i < 5; i++ {
			nx, ny := x+dir[0]*i, y+dir[1]*i
			if nx < 0 || nx >= ai.boardSize || ny < 0 || ny >= ai.boardSize || board[ny][nx] != player {
				break
			}
			count++
		}

		// Check negative direction
		for i := 1; i < 5; i++ {
			nx, ny := x-dir[0]*i, y-dir[1]*i
			if nx < 0 || nx >= ai.boardSize || ny < 0 || ny >= ai.boardSize || board[ny][nx] != player {
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

// isBoardFull checks if the board is completely filled
func (ai *EnhancedAIService) isBoardFull(board [][]int) bool {
	for _, row := range board {
		for _, cell := range row {
			if cell == 0 {
				return false
			}
		}
	}
	return true
}

// hashBoard creates a hash of the board state for transposition table
func (ai *EnhancedAIService) hashBoard(board [][]int) string {
	var sb strings.Builder
	for _, row := range board {
		for _, cell := range row {
			sb.WriteRune(rune(cell + '0'))
		}
	}

	hash := sha256.Sum256([]byte(sb.String()))
	return hex.EncodeToString(hash[:])
}

// GetStats returns performance statistics
func (ai *EnhancedAIService) GetStats() map[string]interface{} {
	ai.tableMutex.RLock()
	defer ai.tableMutex.RUnlock()

	return map[string]interface{}{
		"nodes_searched":       ai.nodesSearched,
		"cutoffs":             ai.cutoffs,
		"table_entries":       len(ai.transpositionTable),
		"search_time":         time.Since(ai.searchStartTime).String(),
		"pruning_efficiency":  float64(ai.cutoffs) / float64(ai.nodesSearched) * 100,
	}
}

// ClearTranspositionTable clears the transposition table
func (ai *EnhancedAIService) ClearTranspositionTable() {
	ai.tableMutex.Lock()
	defer ai.tableMutex.Unlock()

	ai.transpositionTable = make(map[string]*TranspositionTableEntry)
}

// Helper functions
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}