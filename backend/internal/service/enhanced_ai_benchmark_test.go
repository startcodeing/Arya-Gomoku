// Benchmark tests for Enhanced AI Service
package service

import (
	"fmt"
	"testing"
	"time"

	"gomoku-backend/internal/model"
)

// BenchmarkGetAIMove benchmarks AI move generation for different difficulties
func BenchmarkGetAIMove(b *testing.B) {
	ai := NewEnhancedAIService()
	board := createEmptyBoard()
	lastMove := model.Move{X: -1, Y: -1}

	benchmarks := []struct {
		name      string
		difficulty Difficulty
	}{
		{"Easy", Easy},
		{"Medium", Medium},
		{"Hard", Hard},
		{"Expert", Expert},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ai.GetAIMove(board, lastMove, bm.difficulty)
			}
		})
	}
}

// BenchmarkGetAIMoveComplex benchmarks AI on a more complex board
func BenchmarkGetAIMoveComplex(b *testing.B) {
	ai := NewEnhancedAIService()
	board := createComplexBoard()
	lastMove := model.Move{X: 7, Y: 7}

	benchmarks := []struct {
		name      string
		difficulty Difficulty
	}{
		{"Easy", Easy},
		{"Medium", Medium},
		{"Hard", Hard},
		{"Expert", Expert},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ai.GetAIMove(board, lastMove, bm.difficulty)
			}
		})
	}
}

// BenchmarkEvaluateBoard benchmarks board evaluation
func BenchmarkEvaluateBoard(b *testing.B) {
	ai := NewEnhancedAIService()
	board := createComplexBoard()
	lastMove := model.Move{X: 7, Y: 7}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ai.evaluateBoard(board, lastMove)
	}
}

// BenchmarkCheckWin benchmarks win detection
func BenchmarkCheckWin(b *testing.B) {
	ai := NewEnhancedAIService()
	board := createHorizontalWinBoard(2)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ai.checkWin(board, 7, 7, 2)
	}
}

// BenchmarkTranspositionTable benchmarks transposition table operations
func BenchmarkTranspositionTable(b *testing.B) {
	ai := NewEnhancedAIService()
	board := createEmptyBoard()
	lastMove := model.Move{X: 7, Y: 7}

	// First, populate the table
	ai.GetAIMove(board, lastMove, Medium)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hash := ai.hashBoard(board)
		ai.tableMutex.RLock()
		_ = ai.transpositionTable[hash]
		ai.tableMutex.RUnlock()
	}
}

// BenchmarkHashBoard benchmarks board hashing
func BenchmarkHashBoard(b *testing.B) {
	ai := NewEnhancedAIService()
	board := createComplexBoard()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ai.hashBoard(board)
	}
}

// BenchmarkGetAvailableMoves benchmarks move generation
func BenchmarkGetAvailableMoves(b *testing.B) {
	ai := NewEnhancedAIService()
	boards := []struct {
		name  string
		board [][]int
		move  model.Move
	}{
		{"Empty", createEmptyBoard(), model.Move{X: -1, Y: -1}},
		{"Center", createCenterBoard(), model.Move{X: 7, Y: 7}},
		{"Complex", createComplexBoard(), model.Move{X: 7, Y: 7}},
	}

	for _, bm := range boards {
		b.Run(bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ai.getAvailableMoves(bm.board, bm.move)
			}
		})
	}
}

// TestPerformanceComparison tests and reports performance characteristics
func TestPerformanceComparison(t *testing.T) {
	ai := NewEnhancedAIService()

	// Clear any existing data
	ai.ClearTranspositionTable()

	tests := []struct {
		name      string
		difficulty Difficulty
		board     [][]int
		lastMove  model.Move
		maxTime   time.Duration
	}{
		{
			name:      "Easy_Empty",
			difficulty: Easy,
			board:     createEmptyBoard(),
			lastMove:  model.Move{X: -1, Y: -1},
			maxTime:   100 * time.Millisecond,
		},
		{
			name:      "Medium_Empty",
			difficulty: Medium,
			board:     createEmptyBoard(),
			lastMove:  model.Move{X: -1, Y: -1},
			maxTime:   500 * time.Millisecond,
		},
		{
			name:      "Hard_Empty",
			difficulty: Hard,
			board:     createEmptyBoard(),
			lastMove:  model.Move{X: -1, Y: -1},
			maxTime:   2 * time.Second,
		},
		{
			name:      "Medium_Complex",
			difficulty: Medium,
			board:     createComplexBoard(),
			lastMove:  model.Move{X: 7, Y: 7},
			maxTime:   1 * time.Second,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Clear table for fair comparison
			ai.ClearTranspositionTable()

			start := time.Now()
			move := ai.GetAIMove(test.board, test.lastMove, test.difficulty)
			duration := time.Since(start)

			t.Logf("Move: (%d,%d), Score: %d", move.X, move.Y, move.Score)
			t.Logf("Time: %v", duration)

			// Check performance expectations
			if duration > test.maxTime {
				t.Errorf("Expected move within %v, took %v", test.maxTime, duration)
			}

			// Get stats
			stats := ai.GetStats()
			t.Logf("Nodes searched: %v", stats["nodes_searched"])
			t.Logf("Cutoffs: %v", stats["cutoffs"])
			t.Logf("Table entries: %v", stats["table_entries"])

			// Validate move
			if move.X < 0 || move.X >= 15 || move.Y < 0 || move.Y >= 15 {
				t.Errorf("Invalid move: (%d,%d)", move.X, move.Y)
			}

			// For complex boards, higher difficulty should search more nodes
			if test.name == "Medium_Complex" {
				if nodes := stats["nodes_searched"]; nodes == 0 {
					t.Error("Expected some nodes to be searched")
				}
			}
		})
	}
}

// TestStressTest performs stress testing with multiple rapid moves
func TestStressTest(t *testing.T) {
	ai := NewEnhancedAIService()
	board := createEmptyBoard()
	lastMove := model.Move{X: -1, Y: -1}

	// Simulate a rapid game
	numMoves := 20
	totalTime := time.Duration(0)

	for i := 0; i < numMoves; i++ {
		player := 2 // Always AI for this test
		difficulty := Medium

		start := time.Now()
		move := ai.GetAIMove(board, lastMove, difficulty)
		duration := time.Since(start)
		totalTime += duration

		t.Logf("Move %d: (%d,%d) in %v", i+1, move.X, move.Y, duration)

		// Apply move
		board[move.Y][move.X] = player
		lastMove = model.Move{X: move.X, Y: move.Y}

		// Validate move
		if move.X < 0 || move.X >= 15 || move.Y < 0 || move.Y >= 15 {
			t.Errorf("Invalid move %d: (%d,%d)", i+1, move.X, move.Y)
			break
		}

		// Check for reasonable time (should be fast for consecutive moves due to cache)
		if duration > 1*time.Second {
			t.Errorf("Move %d took too long: %v", i+1, duration)
		}
	}

	avgTime := totalTime / time.Duration(numMoves)
	t.Logf("Average time per move: %v", avgTime)
	t.Logf("Total time for %d moves: %v", numMoves, totalTime)

	// Get final stats
	stats := ai.GetStats()
	t.Logf("Final stats: %+v", stats)

	// Performance should improve over time due to transposition table
	if avgTime > 200*time.Millisecond {
		t.Errorf("Expected average time under 200ms, got %v", avgTime)
	}
}

// TestMemoryUsage checks that memory usage stays reasonable
func TestMemoryUsage(t *testing.T) {
	ai := NewEnhancedAIService()

	// Clear initial state
	ai.ClearTranspositionTable()
	initialStats := ai.GetStats()
	initialEntries := initialStats["table_entries"]

	// Generate many moves to populate transposition table
	board := createEmptyBoard()
	lastMove := model.Move{X: -1, Y: -1}

	numMoves := 50
	for i := 0; i < numMoves; i++ {
		// Use different difficulties to generate varied positions
		difficulty := Difficulty(i%4 + 1)
		move := ai.GetAIMove(board, lastMove, difficulty)

		// Apply move
		board[move.Y][move.X] = 2
		lastMove = model.Move{X: move.X, Y: move.Y}

		// Add some random human moves to vary the board
		if i%2 == 0 {
			humanMove := findRandomEmptyCell(board)
			board[humanMove.Y][humanMove.X] = 1
			lastMove = humanMove
		}
	}

	finalStats := ai.GetStats()
	finalEntries := finalStats["table_entries"]

	t.Logf("Initial transposition entries: %v", initialEntries)
	t.Logf("Final transposition entries: %v", finalEntries)
	t.Logf("Entries added: %v", finalEntries.(int)-initialEntries.(int))

	// Check that memory usage is reasonable (should not grow excessively)
	if finalEntries.(int) > 10000 {
		t.Errorf("Transposition table grew too large: %d entries", finalEntries)
	}

	// Test clearing
	ai.ClearTranspositionTable()
	clearedStats := ai.GetStats()
	clearedEntries := clearedStats["table_entries"]

	if clearedEntries != 0 {
		t.Errorf("Expected empty table after clear, got %d entries", clearedEntries)
	}
}

// Helper function to find a random empty cell
func findRandomEmptyCell(board [][]int) model.Move {
	for y := 0; y < 15; y++ {
		for x := 0; x < 15; x++ {
			if board[y][x] == 0 {
				return model.Move{X: x, Y: y}
			}
		}
	}
	return model.Move{X: 7, Y: 7} // Fallback
}

// BenchmarkComparison compares old vs new AI performance (placeholder for old AI)
func BenchmarkComparison(b *testing.B) {
	// This benchmark demonstrates how to compare the new AI with the old one
	// You would need to import and use the old AIService here

	b.Run("Enhanced_Medium", func(b *testing.B) {
		ai := NewEnhancedAIService()
		board := createEmptyBoard()
		lastMove := model.Move{X: -1, Y: -1}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ai.GetAIMove(board, lastMove, Medium)
		}
	})

	// Old AI benchmark would go here
	// b.Run("Original_Heuristic", func(b *testing.B) {
	//     oldAI := service.NewAIService()
	//     // ... benchmark original AI
	// })
}

// TestAIStrength provides a basic strength test
func TestAIStrength(t *testing.T) {
	ai := NewEnhancedAIService()

	// Test 1: Should find immediate win
	board := createEmptyBoard()
	// Set up winning position
	board[7][3] = 2
	board[7][4] = 2
	board[7][5] = 2
	board[7][6] = 2
	lastMove := model.Move{X: 6, Y: 7}

	move := ai.GetAIMove(board, lastMove, Easy) // Even easy should find this

	// Should play at (7,7) to win
	if !(move.X == 7 && move.Y == 7) && !(move.X == 2 && move.Y == 7) {
		t.Errorf("Expected winning move at (7,7) or (7,2), got (%d,%d)", move.X, move.Y)
	}

	if move.Score < 100000 {
		t.Errorf("Expected winning score, got %d", move.Score)
	}

	// Test 2: Should block opponent's win
	board = createEmptyBoard()
	// Set up human winning position
	board[6][5] = 1
	board[6][6] = 1
	board[6][7] = 1
	board[6][8] = 1
	lastMove = model.Move{X: 8, Y: 6}

	move = ai.GetAIMove(board, lastMove, Medium)

	// Should block at (6,9) or (6,4)
	if !(move.X == 9 && move.Y == 6) && !(move.X == 4 && move.Y == 6) {
		t.Errorf("Expected blocking move at (6,9) or (6,4), got (%d,%d)", move.X, move.Y)
	}
}

// Example of how to run the benchmarks
func ExampleBenchmark() {
	fmt.Println("To run benchmarks:")
	fmt.Println("go test -bench=. ./internal/service/")
	fmt.Println("go test -bench=BenchmarkGetAIMove ./internal/service/")
	fmt.Println("go test -bench=BenchmarkGetAIMoveComplex ./internal/service/")
	fmt.Println("go test -v ./internal/service/ -run TestPerformance")
}