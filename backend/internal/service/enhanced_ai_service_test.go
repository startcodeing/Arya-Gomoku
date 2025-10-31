// Unit tests for Enhanced AI Service
package service

import (
	"testing"
	"time"

	"gomoku-backend/internal/model"
)

func TestNewEnhancedAIService(t *testing.T) {
	ai := NewEnhancedAIService()

	if ai.boardSize != 15 {
		t.Errorf("Expected board size 15, got %d", ai.boardSize)
	}

	if len(ai.maxDepth) != 4 {
		t.Errorf("Expected 4 difficulty levels, got %d", len(ai.maxDepth))
	}

	if ai.maxDepth[Easy] != 2 {
		t.Errorf("Expected Easy depth 2, got %d", ai.maxDepth[Easy])
	}

	if ai.timeLimit != 5*time.Second {
		t.Errorf("Expected 5s time limit, got %v", ai.timeLimit)
	}
}

func TestGetAIMove_EmptyBoard(t *testing.T) {
	ai := NewEnhancedAIService()
	board := createEmptyBoard()
	lastMove := model.Move{X: -1, Y: -1}

	// Test all difficulty levels
	difficulties := []Difficulty{Easy, Medium, Hard, Expert}
	expectedMoves := []struct {
		minX, minY, maxX, maxY int
	}{
		{6, 6, 8, 8}, // Easy: around center
		{6, 6, 8, 8}, // Medium: around center
		{6, 6, 8, 8}, // Hard: around center
		{6, 6, 8, 8}, // Expert: around center
	}

	for i, difficulty := range difficulties {
		move := ai.GetAIMove(board, lastMove, difficulty)
		expected := expectedMoves[i]

		if move.X < expected.minX || move.X > expected.maxX ||
			move.Y < expected.minY || move.Y > expected.maxY {
			t.Errorf("Difficulty %v: expected move near center (7,7), got (%d,%d)",
				difficulty, move.X, move.Y)
		}

		if move.Score < 0 {
			t.Errorf("Difficulty %v: expected positive score, got %d", difficulty, move.Score)
		}
	}
}

func TestGetAIMove_WinDetection(t *testing.T) {
	ai := NewEnhancedAIService()

	// Test scenario: AI has 4 pieces in a row, should win
	board := createEmptyBoard()
	// Place 4 AI pieces horizontally
	board[7][5] = 2
	board[7][6] = 2
	board[7][7] = 2
	board[7][8] = 2
	// Place a human piece nearby
	board[6][6] = 1

	lastMove := model.Move{X: 6, Y: 6}

	move := ai.GetAIMove(board, lastMove, Medium)

	// Should play winning move at (7,9) or (7,4)
	expectedMoves := [][2]int{{7, 9}, {7, 4}}
	found := false
	for _, expected := range expectedMoves {
		if move.X == expected[0] && move.Y == expected[1] {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected winning move at (7,9) or (7,4), got (%d,%d)", move.X, move.Y)
	}

	if move.Score < 100000 {
		t.Errorf("Expected high score for winning move, got %d", move.Score)
	}
}

func TestGetAIMove_BlockWin(t *testing.T) {
	ai := NewEnhancedAIService()

	// Test scenario: Human has 4 pieces in a row, AI should block
	board := createEmptyBoard()
	// Place 4 human pieces horizontally
	board[7][5] = 1
	board[7][6] = 1
	board[7][7] = 1
	board[7][8] = 1
	// Place AI piece nearby
	board[6][6] = 2

	lastMove := model.Move{X: 6, Y: 6}

	move := ai.GetAIMove(board, lastMove, Medium)

	// Should block at (7,9) or (7,4)
	expectedMoves := [][2]int{{7, 9}, {7, 4}}
	found := false
	for _, expected := range expectedMoves {
		if move.X == expected[0] && move.Y == expected[1] {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected blocking move at (7,9) or (7,4), got (%d,%d)", move.X, move.Y)
	}
}

func TestGetAIMove_TimeLimit(t *testing.T) {
	ai := NewEnhancedAIService()
	// Set very short time limit for testing
	ai.timeLimit = 10 * time.Millisecond

	board := createComplexBoard() // Create a board with many pieces
	lastMove := model.Move{X: 7, Y: 7}

	start := time.Now()
	move := ai.GetAIMove(board, lastMove, Expert) // Use Expert to trigger deep search
	duration := time.Since(start)

	// Should complete within reasonable time (allow some margin for test execution)
	if duration > 100*time.Millisecond {
		t.Errorf("Expected move within 100ms, took %v", duration)
	}

	// Should still return a valid move
	if move.X < 0 || move.X >= 15 || move.Y < 0 || move.Y >= 15 {
		t.Errorf("Expected valid move coordinates, got (%d,%d)", move.X, move.Y)
	}
}

func TestCheckWin(t *testing.T) {
	ai := NewEnhancedAIService()

	tests := []struct {
		name     string
		board    [][]int
		x, y     int
		player   int
		expected bool
	}{
		{
			name:     "Horizontal win",
			board:    createHorizontalWinBoard(2),
			x:        7,
			y:        7,
			player:   2,
			expected: true,
		},
		{
			name:     "Vertical win",
			board:    createVerticalWinBoard(1),
			x:        7,
			y:        7,
			player:   1,
			expected: true,
		},
		{
			name:     "Diagonal win",
			board:    createDiagonalWinBoard(2),
			x:        7,
			y:        7,
			player:   2,
			expected: true,
		},
		{
			name:     "No win",
			board:    createEmptyBoard(),
			x:        7,
			y:        7,
			player:   2,
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ai.checkWin(test.board, test.x, test.y, test.player)
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestEvaluatePositionAdvanced(t *testing.T) {
	ai := NewEnhancedAIService()

	// Test evaluation of a position with potential
	board := createEmptyBoard()
	board[7][6] = 2 // AI piece
	board[7][8] = 2 // AI piece

	score := ai.evaluatePositionAdvanced(board, 7, 7, 2) // Check position between two AI pieces

	if score <= 0 {
		t.Errorf("Expected positive score for connecting position, got %d", score)
	}
}

func TestGetAvailableMoves(t *testing.T) {
	ai := NewEnhancedAIService()

	// Empty board should return center positions
	board := createEmptyBoard()
	lastMove := model.Move{X: -1, Y: -1}

	moves := ai.getAvailableMoves(board, lastMove)

	if len(moves) == 0 {
		t.Error("Expected at least one available move")
	}

	// All moves should be near center for empty board
	centerX, centerY := 7, 7
	for _, move := range moves {
		if abs(move.X-centerX) > 1 || abs(move.Y-centerY) > 1 {
			t.Errorf("Expected move near center, got (%d,%d)", move.X, move.Y)
		}
	}
}

func TestTranspositionTable(t *testing.T) {
	ai := NewEnhancedAIService()

	// Clear table
	ai.ClearTranspositionTable()
	stats := ai.GetStats()
	if entries := stats["table_entries"]; entries != 0 {
		t.Errorf("Expected empty transposition table, got %d entries", entries)
	}

	// Make a move to populate table
	board := createEmptyBoard()
	lastMove := model.Move{X: 7, Y: 7}
	ai.GetAIMove(board, lastMove, Medium)

	stats = ai.GetStats()
	if entries := stats["table_entries"]; entries == 0 {
		t.Error("Expected transposition table to be populated")
	}

	// Clear table again
	ai.ClearTranspositionTable()
	stats = ai.GetStats()
	if entries := stats["table_entries"]; entries != 0 {
		t.Errorf("Expected empty table after clear, got %d entries", entries)
	}
}

func TestGetStats(t *testing.T) {
	ai := NewEnhancedAIService()

	// Initial stats
	stats := ai.GetStats()
	if _, ok := stats["nodes_searched"]; !ok {
		t.Error("Expected nodes_searched in stats")
	}
	if _, ok := stats["cutoffs"]; !ok {
		t.Error("Expected cutoffs in stats")
	}
	if _, ok := stats["table_entries"]; !ok {
		t.Error("Expected table_entries in stats")
	}

	// Make a move and check stats change
	board := createEmptyBoard()
	lastMove := model.Move{X: -1, Y: -1}
	ai.GetAIMove(board, lastMove, Medium)

	statsAfter := ai.GetStats()
	if nodesBefore := stats["nodes_searched"]; nodesBefore != statsAfter["nodes_searched"] {
		// This should change - we expect some nodes to be searched
	} else {
		t.Error("Expected nodes_searched to increase after making a move")
	}
}

func TestDifficulties(t *testing.T) {
	ai := NewEnhancedAIService()
	board := createEmptyBoard()
	lastMove := model.Move{X: -1, Y: -1}

	// Test that different difficulties produce different results (usually)
	easyMove := ai.GetAIMove(board, lastMove, Easy)
	hardMove := ai.GetAIMove(board, lastMove, Hard)

	// Note: This test might occasionally fail if both algorithms choose the same move
	// In that case, we just verify they produce valid moves
	if easyMove.X == hardMove.X && easyMove.Y == hardMove.Y {
		t.Logf("Easy and Hard chose the same move (%d,%d) - this is acceptable", easyMove.X, easyMove.Y)
	}

	// Both moves should be valid
	if easyMove.X < 0 || easyMove.X >= 15 || easyMove.Y < 0 || easyMove.Y >= 15 {
		t.Errorf("Easy move invalid: (%d,%d)", easyMove.X, easyMove.Y)
	}

	if hardMove.X < 0 || hardMove.X >= 15 || hardMove.Y < 0 || hardMove.Y >= 15 {
		t.Errorf("Hard move invalid: (%d,%d)", hardMove.X, hardMove.Y)
	}
}

// Helper functions for creating test boards

func createEmptyBoard() [][]int {
	board := make([][]int, 15)
	for i := range board {
		board[i] = make([]int, 15)
	}
	return board
}

func createHorizontalWinBoard(player int) [][]int {
	board := createEmptyBoard()
	// Place 5 pieces horizontally
	for x := 5; x < 10; x++ {
		board[7][x] = player
	}
	return board
}

func createVerticalWinBoard(player int) [][]int {
	board := createEmptyBoard()
	// Place 5 pieces vertically
	for y := 5; y < 10; y++ {
		board[y][7] = player
	}
	return board
}

func createDiagonalWinBoard(player int) [][]int {
	board := createEmptyBoard()
	// Place 5 pieces diagonally
	for i := 0; i < 5; i++ {
		board[5+i][5+i] = player
	}
	return board
}

func createComplexBoard() [][]int {
	board := createEmptyBoard()

	// Add various pieces to create a complex position
	pieces := [][2]int{
		{7, 7}, {7, 8}, {6, 7}, {6, 8}, {8, 7},
		{5, 5}, {5, 6}, {6, 5}, {6, 6},
		{8, 8}, {8, 9}, {9, 8}, {9, 9},
		{3, 3}, {3, 4}, {4, 3}, {4, 4},
		{10, 10}, {10, 11}, {11, 10}, {11, 11},
	}

	for i, pos := range pieces {
		player := (i % 2) + 1 // Alternate between 1 and 2
		board[pos[1]][pos[0]] = player
	}

	return board
}
