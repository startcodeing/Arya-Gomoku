// Package controller handles HTTP requests and responses for the Gomoku API
// This package contains the REST API endpoints for AI gameplay and game management
package controller

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"gomoku-backend/internal/model"
	"gomoku-backend/internal/repository"
	"gomoku-backend/internal/service"
)

// AIController handles AI-related HTTP requests
type AIController struct {
	aiService         *service.AIService
	enhancedAIService *service.EnhancedAIService
	gameRepo          repository.GameRepository
}

// NewAIController creates a new AI controller instance
func NewAIController(gameRepo repository.GameRepository) *AIController {
	return &AIController{
		aiService:         service.NewAIService(),
		enhancedAIService: service.NewEnhancedAIService(),
		gameRepo:          gameRepo,
	}
}

// GetAIMove handles POST /api/ai/move requests
// This endpoint receives the current board state and returns the AI's next move
func (ac *AIController) GetAIMove(c *gin.Context) {
	var request model.GameRequest

	// Parse JSON request body
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Get difficulty level from query parameter (default: Medium)
	difficultyStr := c.DefaultQuery("difficulty", "medium")
	useEnhanced := c.DefaultQuery("enhanced", "true") == "true"

	var difficulty service.Difficulty
	switch difficultyStr {
	case "easy":
		difficulty = service.Easy
	case "medium":
		difficulty = service.Medium
	case "hard":
		difficulty = service.Hard
	case "expert":
		difficulty = service.Expert
	default:
		difficulty = service.Medium
	}

	// Validate board dimensions
	if len(request.Board) != 15 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Board must be 15x15",
		})
		return
	}

	for i, row := range request.Board {
		if len(row) != 15 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Board row " + string(rune(i+1)) + " must have 15 columns",
			})
			return
		}
	}

	// Validate player
	if request.Player != 1 && request.Player != 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Player must be 1 (human) or 2 (AI)",
		})
		return
	}

	// Validate last move coordinates
	if request.LastMove.X < 0 || request.LastMove.X >= 15 ||
	   request.LastMove.Y < 0 || request.LastMove.Y >= 15 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Last move coordinates must be within 0-14 range",
		})
		return
	}

	// Additional board state validation
	if err := ac.validateBoardState(request.Board); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Get AI move using enhanced or regular AI
	var aiMove model.AIMove
	if useEnhanced {
		aiMove = ac.enhancedAIService.GetAIMove(request.Board, request.LastMove, difficulty)
	} else {
		aiMove = ac.aiService.GetAIMove(request.Board, request.LastMove)
	}
	
	// Create a temporary board to check game state after AI move
	tempBoard := make([][]int, 15)
	for i := range tempBoard {
		tempBoard[i] = make([]int, 15)
		copy(tempBoard[i], request.Board[i])
	}
	
	// Apply AI move to temporary board
	tempBoard[aiMove.Y][aiMove.X] = 2
	
	// Create board model to check win condition
	board := &model.Board{
		Grid:          tempBoard,
		Size:          15,
		CurrentPlayer: 1, // Next turn would be human
		MoveCount:     ac.countMoves(tempBoard) + 1,
	}
	
	// Check game state after AI move
	aiMoveModel := &model.Move{X: aiMove.X, Y: aiMove.Y, Player: 2}
	gameState := board.GetGameState(aiMoveModel)
	
	// Prepare response
	response := model.GameResponse{
		AIMove:     aiMove,
		GameStatus: gameState.Status,
		Winner:     gameState.Winner,
	}

	// Add enhanced AI stats if using enhanced AI
	if useEnhanced {
		stats := ac.enhancedAIService.GetStats()
		c.JSON(http.StatusOK, gin.H{
			"aiMove":       response.AIMove,
			"gameStatus":   response.GameStatus,
			"winner":       response.Winner,
			"difficulty":   difficultyStr,
			"aiEngine":     "enhanced_minimax",
			"stats":        stats,
		})
	} else {
		c.JSON(http.StatusOK, response)
	}
}

// GetGameStatus handles GET /api/game/status requests (reserved for future use)
func (ac *AIController) GetGameStatus(c *gin.Context) {
	// This endpoint is reserved for future game state management
	c.JSON(http.StatusOK, gin.H{
		"status":  "available",
		"message": "AI service is running",
		"version": "1.0.0",
	})
}

// ResetGame handles POST /api/game/reset requests (reserved for future use)
func (ac *AIController) ResetGame(c *gin.Context) {
	// This endpoint is reserved for future game session management
	newBoard := model.NewBoard()
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Game reset successfully",
		"board":   newBoard,
	})
}

// GetAIStats handles GET /api/ai/stats requests
// Returns performance statistics of the enhanced AI
func (ac *AIController) GetAIStats(c *gin.Context) {
	stats := ac.enhancedAIService.GetStats()

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"stats":  stats,
	})
}

// ClearCache handles POST /api/ai/cache/clear requests
// Clears the transposition table
func (ac *AIController) ClearCache(c *gin.Context) {
	ac.enhancedAIService.ClearTranspositionTable()

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Transposition table cleared successfully",
	})
}

// CreateAIGame handles POST /api/ai/games requests
// Creates a new AI game session
func (ac *AIController) CreateAIGame(c *gin.Context) {
	type CreateGameRequest struct {
		Difficulty string `json:"difficulty"`
	}

	var req CreateGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Get user ID from JWT token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "User not authenticated",
		})
		return
	}

	// Set default difficulty
	if req.Difficulty == "" {
		req.Difficulty = "medium"
	}

	// Convert difficulty string to Difficulty type
	var difficulty model.Difficulty
	switch req.Difficulty {
	case "easy":
		difficulty = model.DifficultyEasy
	case "medium":
		difficulty = model.DifficultyMedium
	case "hard":
		difficulty = model.DifficultyHard
	case "expert":
		difficulty = model.DifficultyExpert
	default:
		difficulty = model.DifficultyMedium
	}

	// Create new AI game
	game := &model.AIGame{
		ID:         uuid.New().String(),
		UserID:     userID.(string),
		Difficulty: difficulty,
		Status:     model.GameStatusPlaying,
		Moves:      []model.GameMove{},
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Save to database
	ctx := context.Background()
	if err := ac.gameRepo.CreateAIGame(ctx, game); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to create game",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    game,
	})
}

// GetAIGame handles GET /api/ai/games/:id requests
// Retrieves an AI game by ID
func (ac *AIController) GetAIGame(c *gin.Context) {
	gameID := c.Param("id")
	if gameID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Game ID is required",
		})
		return
	}

	ctx := context.Background()
	game, err := ac.gameRepo.GetAIGameByID(ctx, gameID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Game not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    game,
	})
}

// MakeAIMove handles POST /api/ai/games/:id/move requests
// Makes a move in an AI game and gets AI response
func (ac *AIController) MakeAIMove(c *gin.Context) {
	gameID := c.Param("id")
	if gameID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Game ID is required",
		})
		return
	}

	type MoveRequest struct {
		X int `json:"x" binding:"required,min=0,max=14"`
		Y int `json:"y" binding:"required,min=0,max=14"`
	}

	var req MoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid move format",
			"details": err.Error(),
		})
		return
	}

	ctx := context.Background()
	
	// Get game from database
	game, err := ac.gameRepo.GetAIGameByID(ctx, gameID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Game not found",
		})
		return
	}

	// Check if game is still active
	if game.Status != model.GameStatusPlaying {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Game is not active",
		})
		return
	}

	// Reconstruct board from moves
	board := make([][]int, 15)
	for i := range board {
		board[i] = make([]int, 15)
	}
	for _, move := range game.Moves {
		board[move.Y][move.X] = move.Player
	}

	// Validate move
	if board[req.Y][req.X] != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Position already occupied",
		})
		return
	}

	// Make human move
	board[req.Y][req.X] = 1
	humanMove := model.GameMove{
		X:         req.X,
		Y:         req.Y,
		Player:    1,
		Timestamp: time.Now(),
	}
	game.Moves = append(game.Moves, humanMove)

	// Check if human wins
	boardModel := &model.Board{Grid: board, Size: 15}
	humanMoveModel := &model.Move{X: req.X, Y: req.Y, Player: 1}
	gameState := boardModel.GetGameState(humanMoveModel)

	if gameState.Status != "playing" {
		if gameState.Status == "win" && gameState.Winner == 1 {
			game.Status = model.GameStatusCompleted
		} else if gameState.Status == "win" && gameState.Winner == 2 {
			game.Status = model.GameStatusCompleted
		} else {
			game.Status = model.GameStatusCompleted
		}
		endTime := time.Now()
		game.EndedAt = &endTime
		
		// Update game in database
		if err := ac.gameRepo.UpdateAIGame(ctx, game); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to update game",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success":    true,
			"data":       game,
			"gameStatus": gameState.Status,
			"winner":     gameState.Winner,
		})
		return
	}

	// Get AI move
	var difficulty service.Difficulty
	switch game.Difficulty {
	case "easy":
		difficulty = service.Easy
	case "medium":
		difficulty = service.Medium
	case "hard":
		difficulty = service.Hard
	case "expert":
		difficulty = service.Expert
	default:
		difficulty = service.Medium
	}

	lastMove := model.Move{X: req.X, Y: req.Y}
	aiMove := ac.enhancedAIService.GetAIMove(board, lastMove, difficulty)

	// Make AI move
	board[aiMove.Y][aiMove.X] = 2
	aiMoveRecord := model.GameMove{
		X:         aiMove.X,
		Y:         aiMove.Y,
		Player:    2,
		Timestamp: time.Now(),
	}
	game.Moves = append(game.Moves, aiMoveRecord)

	// Check if AI wins
	aiMoveModel := &model.Move{X: aiMove.X, Y: aiMove.Y, Player: 2}
	gameState = boardModel.GetGameState(aiMoveModel)

	if gameState.Status != "playing" {
		if gameState.Status == "win" && gameState.Winner == 2 {
			game.Status = model.GameStatusCompleted
		} else if gameState.Status == "win" && gameState.Winner == 1 {
			game.Status = model.GameStatusCompleted
		} else {
			game.Status = model.GameStatusCompleted
		}
		endTime := time.Now()
		game.EndedAt = &endTime
	}

	game.UpdatedAt = time.Now()

	// Update game in database
	if err := ac.gameRepo.UpdateAIGame(ctx, game); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to update game",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"data":       game,
		"aiMove":     aiMove,
		"gameStatus": gameState.Status,
		"winner":     gameState.Winner,
	})
}

// GetUserAIGames handles GET /api/ai/games requests
// Retrieves AI games for a user
func (ac *AIController) GetUserAIGames(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "User ID is required",
		})
		return
	}

	page := 1
	limit := 10
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	offset := (page - 1) * limit

	ctx := context.Background()
	games, total, err := ac.gameRepo.GetUserAIGames(ctx, userID, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve games",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"games": games,
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

// GetDifficultyLevels handles GET /api/ai/difficulties requests
// Returns available difficulty levels
func (ac *AIController) GetDifficultyLevels(c *gin.Context) {
	difficulties := []map[string]interface{}{
		{
			"level":        "easy",
			"description":  "Simple heuristic moves, quick response",
			"maxDepth":     2,
			"estimatedTime": "< 100ms",
		},
		{
			"level":        "medium",
			"description":  "Balanced gameplay with moderate analysis",
			"maxDepth":     4,
			"estimatedTime": "100-500ms",
		},
		{
			"level":        "hard",
			"description":  "Deep analysis with strong play",
			"maxDepth":     6,
			"estimatedTime": "500ms-2s",
		},
		{
			"level":        "expert",
			"description":  "Maximum depth analysis with time limit",
			"maxDepth":     8,
			"estimatedTime": "2-5s",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       "success",
		"difficulties": difficulties,
	})
}

// BenchmarkAI handles POST /api/ai/benchmark requests
// Runs a benchmark of the AI performance
func (ac *AIController) BenchmarkAI(c *gin.Context) {
	type BenchmarkRequest struct {
		Difficulty string `json:"difficulty"`
		MoveCount  int    `json:"moveCount"`
	}

	var req BenchmarkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Set defaults
	if req.Difficulty == "" {
		req.Difficulty = "medium"
	}
	if req.MoveCount == 0 {
		req.MoveCount = 10
	}

	var difficulty service.Difficulty
	switch req.Difficulty {
	case "easy":
		difficulty = service.Easy
	case "medium":
		difficulty = service.Medium
	case "hard":
		difficulty = service.Hard
	case "expert":
		difficulty = service.Expert
	default:
		difficulty = service.Medium
	}

	// Create test board
	board := make([][]int, 15)
	for i := range board {
		board[i] = make([]int, 15)
	}

	// Simulate some moves
	board[7][7] = 1 // Human center
	board[7][8] = 2 // AI response
	lastMove := model.Move{X: 7, Y: 8}

	// Run benchmark
	startTime := time.Now()
	totalNodes := uint64(0)
	totalCutoffs := uint64(0)

	for i := 0; i < req.MoveCount; i++ {
		aiMove := ac.enhancedAIService.GetAIMove(board, lastMove, difficulty)
		stats := ac.enhancedAIService.GetStats()

		totalNodes += stats["nodes_searched"].(uint64)
		totalCutoffs += stats["cutoffs"].(uint64)

		// Apply move to continue game
		board[aiMove.Y][aiMove.X] = 2
		lastMove = model.Move{X: aiMove.X, Y: aiMove.Y}

		// Simulate human response (random valid move)
		if i < req.MoveCount-1 {
			board[6][6+i] = 1
			lastMove = model.Move{X: 6 + i, Y: 6}
		}
	}

	duration := time.Since(startTime)

	c.JSON(http.StatusOK, gin.H{
		"status":       "success",
		"difficulty":   req.Difficulty,
		"moveCount":    req.MoveCount,
		"totalTime":    duration.String(),
		"avgTime":      (duration / time.Duration(req.MoveCount)).String(),
		"totalNodes":   totalNodes,
		"totalCutoffs": totalCutoffs,
		"nodesPerMs":   float64(totalNodes) / float64(duration.Milliseconds()),
	})
}

// countMoves counts the total number of pieces on the board
func (ac *AIController) countMoves(board [][]int) int {
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

// validateBoardState performs additional validation on the board state
func (ac *AIController) validateBoardState(board [][]int) error {
	playerCount := 0
	aiCount := 0

	for _, row := range board {
		for _, cell := range row {
			switch cell {
			case 1:
				playerCount++
			case 2:
				aiCount++
			case 0:
				// Empty cell, valid
			default:
				return errors.New("Invalid cell value. Must be 0, 1, or 2")
			}
		}
	}

	// Check if move counts are reasonable (AI should not have more pieces than player + 1)
	if aiCount > playerCount || playerCount > aiCount+1 {
		return errors.New("Invalid board state: unrealistic piece distribution")
	}

	return nil
}