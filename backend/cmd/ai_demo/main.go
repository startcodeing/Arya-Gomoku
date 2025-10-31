// AI Demo - demonstrates enhanced AI capabilities
package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gomoku-backend/internal/model"
	"gomoku-backend/internal/service"
)

func main() {
	fmt.Println("🎮 Arya-Gomoku Enhanced AI Demo")
	fmt.Println("==============================")

	if len(os.Args) > 1 && os.Args[1] == "--help" {
		printUsage()
		return
	}

	// Initialize enhanced AI
	ai := service.NewEnhancedAIService()

	// Run different demo scenarios
	scenarios := []struct {
		name string
		fn   func(*service.EnhancedAIService)
	}{
		{"Empty Board Analysis", demoEmptyBoard},
		{"Win Detection Test", demoWinDetection},
		{"Block Detection Test", demoBlockDetection},
		{"Difficulty Comparison", demoDifficultyComparison},
		{"Performance Test", demoPerformance},
		{"Statistics Overview", demoStatistics},
	}

	for _, scenario := range scenarios {
		fmt.Printf("\n🔍 %s\n", scenario.name)
		fmt.Println(strings.Repeat("-", len(scenario.name)+4))
		scenario.fn(ai)
		fmt.Println()
	}

	fmt.Println("✅ Demo completed!")
	fmt.Println("\nTo run tests:")
	fmt.Println("  go test ./internal/service/ -v")
	fmt.Println("\nTo run benchmarks:")
	fmt.Println("  go test ./internal/service/ -bench=. -benchmem")
}

func printUsage() {
	fmt.Println("Usage: go run cmd/ai_demo/main.go [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --help    Show this help message")
	fmt.Println()
	fmt.Println("This demo showcases the enhanced AI capabilities including:")
	fmt.Println("  - Minimax algorithm with alpha-beta pruning")
	fmt.Println("  - Transposition table optimization")
	fmt.Println("  - Multiple difficulty levels")
	fmt.Println("  - Performance statistics")
}

func demoEmptyBoard(ai *service.EnhancedAIService) {
	board := createEmptyBoard()
	lastMove := model.Move{X: -1, Y: -1}

	fmt.Println("Testing AI move on empty board...")

	for _, difficulty := range []service.Difficulty{service.Easy, service.Medium, service.Hard, service.Expert} {
		start := time.Now()
		move := ai.GetAIMove(board, lastMove, difficulty)
		duration := time.Since(start)

		fmt.Printf("  %-6s: (%2d,%2d) Score:%6d Time:%v\n",
			difficultyString(difficulty), move.X, move.Y, move.Score, duration)
	}
}

func demoWinDetection(ai *service.EnhancedAIService) {
	board := createEmptyBoard()

	// Set up winning position for AI (4 pieces in a row)
	positions := [][2]int{{7, 5}, {7, 6}, {7, 7}, {7, 8}}
	for _, pos := range positions {
		board[pos[1]][pos[0]] = 2 // AI pieces
	}

	// Add a human piece nearby
	board[6][6] = 1

	lastMove := model.Move{X: 6, Y: 6}

	fmt.Println("AI has 4 pieces in a row at (7,5)-(7,8)")
	fmt.Println("Human played at (6,6)")
	fmt.Println("AI should find winning move...")

	start := time.Now()
	move := ai.GetAIMove(board, lastMove, service.Medium)
	duration := time.Since(start)

	expectedMoves := [][2]int{{7, 9}, {7, 4}}
	isCorrect := false
	for _, expected := range expectedMoves {
		if move.X == expected[0] && move.Y == expected[1] {
			isCorrect = true
			break
		}
	}

	status := "❌"
	if isCorrect {
		status = "✅"
	}

	fmt.Printf("  %s AI chose: (%d,%d) Score:%d Time:%v\n",
		status, move.X, move.Y, move.Score, duration)

	if isCorrect {
		fmt.Println("  Correct! AI found the winning move.")
	} else {
		fmt.Printf("  Expected move at (7,9) or (7,4)\n")
	}
}

func demoBlockDetection(ai *service.EnhancedAIService) {
	board := createEmptyBoard()

	// Set up winning position for Human (4 pieces in a row)
	positions := [][2]int{{7, 5}, {7, 6}, {7, 7}, {7, 8}}
	for _, pos := range positions {
		board[pos[1]][pos[0]] = 1 // Human pieces
	}

	// Add an AI piece nearby
	board[6][6] = 2

	lastMove := model.Move{X: 6, Y: 6}

	fmt.Println("Human has 4 pieces in a row at (7,5)-(7,8)")
	fmt.Println("AI played at (6,6)")
	fmt.Println("AI should block human's winning move...")

	start := time.Now()
	move := ai.GetAIMove(board, lastMove, service.Medium)
	duration := time.Since(start)

	expectedMoves := [][2]int{{7, 9}, {7, 4}}
	isCorrect := false
	for _, expected := range expectedMoves {
		if move.X == expected[0] && move.Y == expected[1] {
			isCorrect = true
			break
		}
	}

	status := "❌"
	if isCorrect {
		status = "✅"
	}

	fmt.Printf("  %s AI chose: (%d,%d) Score:%d Time:%v\n",
		status, move.X, move.Y, move.Score, duration)

	if isCorrect {
		fmt.Println("  Correct! AI blocked the winning move.")
	} else {
		fmt.Printf("  Expected block at (7,9) or (7,4)\n")
	}
}

func demoDifficultyComparison(ai *service.EnhancedAIService) {
	board := createComplexBoard()
	lastMove := model.Move{X: 7, Y: 7}

	fmt.Println("Comparing AI strengths on complex board...")
	fmt.Printf("Board has %d pieces placed\n", countPieces(board))

	results := make(map[service.Difficulty]model.AIMove)

	for _, difficulty := range []service.Difficulty{service.Easy, service.Medium, service.Hard, service.Expert} {
		start := time.Now()
		move := ai.GetAIMove(board, lastMove, difficulty)
		duration := time.Since(start)

		results[difficulty] = move

		fmt.Printf("  %-6s: (%2d,%2d) Score:%6d Time:%v\n",
			difficultyString(difficulty), move.X, move.Y, move.Score, duration)
	}

	// Check if different difficulties chose different moves
	uniqueMoves := make(map[string]bool)
	for _, move := range results {
		key := fmt.Sprintf("(%d,%d)", move.X, move.Y)
		uniqueMoves[key] = true
	}

	if len(uniqueMoves) > 1 {
		fmt.Printf("  🎯 Different difficulties chose different moves (%d unique positions)\n", len(uniqueMoves))
	} else {
		fmt.Printf("  📍 All difficulties agreed on the same move\n")
	}
}

func demoPerformance(ai *service.EnhancedAIService) {
	fmt.Println("Running performance test...")

	// Clear cache for fair test
	ai.ClearTranspositionTable()

	board := createEmptyBoard()
	lastMove := model.Move{X: -1, Y: -1}

	numMoves := 5
	totalTime := time.Duration(0)

	fmt.Printf("Making %d consecutive moves...\n", numMoves)

	for i := 0; i < numMoves; i++ {
		start := time.Now()
		move := ai.GetAIMove(board, lastMove, service.Medium)
		duration := time.Since(start)
		totalTime += duration

		fmt.Printf("  Move %d: (%d,%d) Score:%6d Time:%v\n",
			i+1, move.X, move.Y, move.Score, duration)

		// Apply move
		board[move.Y][move.X] = 2
		lastMove = model.Move{X: move.X, Y: move.Y}

		// Add human move for variety
		if i < numMoves-1 {
			humanMove := findNearbyEmptyCell(board, move.X, move.Y)
			board[humanMove.Y][humanMove.X] = 1
			lastMove = humanMove
		}
	}

	avgTime := totalTime / time.Duration(numMoves)
	fmt.Printf("\nAverage time per move: %v\n", avgTime)
	fmt.Printf("Total time: %v\n", totalTime)

	// Show stats
	stats := ai.GetStats()
	fmt.Printf("Nodes searched: %v\n", stats["nodes_searched"])
	fmt.Printf("Alpha-beta cutoffs: %v\n", stats["cutoffs"])
	fmt.Printf("Transposition table entries: %v\n", stats["table_entries"])
}

func demoStatistics(ai *service.EnhancedAIService) {
	fmt.Println("Current AI Statistics:")

	stats := ai.GetStats()

	fmt.Printf("  Nodes searched:       %v\n", stats["nodes_searched"])
	fmt.Printf("  Alpha-beta cutoffs:   %v\n", stats["cutoffs"])
	fmt.Printf("  Table entries:        %v\n", stats["table_entries"])
	fmt.Printf("  Search time:          %v\n", stats["search_time"])

	if pruningEfficiency, ok := stats["pruning_efficiency"]; ok {
		fmt.Printf("  Pruning efficiency:   %.1f%%\n", pruningEfficiency)
	}

	// Demonstrate cache clearing
	fmt.Println("\nClearing transposition table...")
	ai.ClearTranspositionTable()

	statsAfter := ai.GetStats()
	fmt.Printf("  Table entries after clear: %v\n", statsAfter["table_entries"])
}

// Helper functions

func createEmptyBoard() [][]int {
	board := make([][]int, 15)
	for i := range board {
		board[i] = make([]int, 15)
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

func countPieces(board [][]int) int {
	count := 0
	for _, row := range board {
		for _, cell := range row {
			if cell != 0 {
				count++
			}
		}
	}
	return count
}

func findNearbyEmptyCell(board [][]int, cx, cy int) model.Move {
	// Search in expanding squares around (cx,cy)
	for radius := 1; radius <= 5; radius++ {
		for dy := -radius; dy <= radius; dy++ {
			for dx := -radius; dx <= radius; dx++ {
				x, y := cx+dx, cy+dy
				if x >= 0 && x < 15 && y >= 0 && y < 15 && board[y][x] == 0 {
					return model.Move{X: x, Y: y}
				}
			}
		}
	}

	// Fallback to center
	return model.Move{X: 7, Y: 7}
}

func difficultyString(difficulty service.Difficulty) string {
	switch difficulty {
	case service.Easy:
		return "Easy"
	case service.Medium:
		return "Medium"
	case service.Hard:
		return "Hard"
	case service.Expert:
		return "Expert"
	default:
		return "Unknown"
	}
}