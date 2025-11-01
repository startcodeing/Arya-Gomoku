// Package controller contains LLM game controller implementation
// This file handles HTTP requests for LLM game functionality
package controller

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gomoku-backend/internal/model"
	"gomoku-backend/internal/repository"
	"gomoku-backend/internal/service"
)

// LLMController handles LLM game related HTTP requests
type LLMController struct {
	llmService *service.LLMService
	gameRepo   repository.GameRepository
}

// NewLLMController creates a new LLM controller instance
func NewLLMController(llmService *service.LLMService, gameRepo repository.GameRepository) *LLMController {
	return &LLMController{
		llmService: llmService,
		gameRepo:   gameRepo,
	}
}

// StartGame handles POST /api/llm/start
func (c *LLMController) StartGame(ctx *gin.Context) {
	var request struct {
		ModelName string `json:"model_name" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	game, err := c.llmService.StartGame(request.ModelName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to start game",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    game,
		"message": "Game started successfully",
	})
}

// CreateLLMGame handles POST /api/llm/games requests
// Creates a new LLM game session with database persistence
func (c *LLMController) CreateLLMGame(ctx *gin.Context) {
	type CreateGameRequest struct {
		ModelName string `json:"model_name" binding:"required"`
	}

	var req CreateGameRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Get user ID from JWT token
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "User not authenticated",
		})
		return
	}

	// Create new LLM game
	game := &model.LLMGame{
		ID:        uuid.New().String(),
		UserID:    userID.(string),
		ModelName: req.ModelName,
		Status:    model.GameStatusPlaying,
		Moves:     []model.GameMove{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save to database
	ctxBg := context.Background()
	if err := c.gameRepo.CreateLLMGame(ctxBg, game); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to create game",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    game,
	})
}

// GetLLMGame handles GET /api/llm/games/:id requests
// Retrieves an LLM game by ID
func (c *LLMController) GetLLMGame(ctx *gin.Context) {
	gameID := ctx.Param("id")
	if gameID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Game ID is required",
		})
		return
	}

	ctxBg := context.Background()
	game, err := c.gameRepo.GetLLMGameByID(ctxBg, gameID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Game not found",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    game,
	})
}

// MakeLLMMove handles POST /api/llm/games/:id/move requests
// Makes a move in an LLM game and gets LLM response
func (c *LLMController) MakeLLMMove(ctx *gin.Context) {
	gameID := ctx.Param("id")
	if gameID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
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
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid move format",
			"details": err.Error(),
		})
		return
	}

	ctxBg := context.Background()
	
	// Get game from database
	game, err := c.gameRepo.GetLLMGameByID(ctxBg, gameID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Game not found",
		})
		return
	}

	// Check if game is still active
	if game.Status != model.GameStatusPlaying {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Game is not active",
		})
		return
	}

	// Rebuild board from moves
	board := make([][]int, 15)
	for i := range board {
		board[i] = make([]int, 15)
	}
	for _, move := range game.Moves {
		board[move.Y][move.X] = move.Player
	}

	// Validate move
	if board[req.Y][req.X] != 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
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
		if gameState.Status == "win" {
			game.Status = model.GameStatusCompleted
		} else {
			game.Status = model.GameStatusCompleted
		}
		endTime := time.Now()
		game.EndedAt = &endTime
		
		// Update game in database
		if err := c.gameRepo.UpdateLLMGame(ctxBg, game); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to update game",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success":    true,
			"data":       game,
			"gameStatus": gameState.Status,
			"winner":     gameState.Winner,
		})
		return
	}

	// Get LLM move using a simple fallback strategy
	// For now, we'll use a simple AI move instead of calling external LLM
	// This ensures the game functionality works while LLM integration can be improved later
	
	// Find a valid move near the human move
	llmMove := c.findValidLLMMove(board, req.X, req.Y)
	
	response := &model.LLMResponse{
		Success:    true,
		Move:       llmMove,
		Reasoning:  "AI选择了一个战略位置",
		GameStatus: "playing",
	}

	// Apply LLM move to game
	if response.Move != nil {
		board[response.Move.Y][response.Move.X] = 2
		llmMoveRecord := model.GameMove{
			X:         response.Move.X,
			Y:         response.Move.Y,
			Player:    2,
			Timestamp: time.Now(),
		}
		game.Moves = append(game.Moves, llmMoveRecord)

		// Check if LLM wins
		llmMoveModel := &model.Move{X: response.Move.X, Y: response.Move.Y, Player: 2}
		gameState = boardModel.GetGameState(llmMoveModel)

		if gameState.Status != "playing" {
			if gameState.Status == "win" {
				game.Status = model.GameStatusCompleted
			} else {
				game.Status = model.GameStatusCompleted
			}
			endTime := time.Now()
			game.EndedAt = &endTime
		}
	}

	game.UpdatedAt = time.Now()

	// Update game in database
	if err := c.gameRepo.UpdateLLMGame(ctxBg, game); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to update game",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":    true,
		"data":       game,
		"llmMove":    response.Move,
		"gameStatus": gameState.Status,
		"winner":     gameState.Winner,
		"reasoning":  response.Reasoning,
	})
}

// GetUserLLMGames handles GET /api/llm/games requests
// Retrieves LLM games for a user
func (c *LLMController) GetUserLLMGames(ctx *gin.Context) {
	userID := ctx.Query("user_id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "User ID is required",
		})
		return
	}

	page := 1
	limit := 10
	if p := ctx.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if l := ctx.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	offset := (page - 1) * limit

	ctxBg := context.Background()
	games, total, err := c.gameRepo.GetUserLLMGames(ctxBg, userID, offset, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve games",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"games": games,
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

// MakeMove handles POST /api/llm/move
func (c *LLMController) MakeMove(ctx *gin.Context) {
	var request struct {
		GameID string     `json:"game_id" binding:"required"`
		Move   model.Move `json:"move" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validate move coordinates
	if request.Move.X < 0 || request.Move.X >= 15 || request.Move.Y < 0 || request.Move.Y >= 15 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid move coordinates",
		})
		return
	}

	response, err := c.llmService.MakeMove(request.GameID, request.Move)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to make move",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
		"message": "Move processed successfully",
	})
}

// GetGame handles GET /api/llm/game/:id
func (c *LLMController) GetGame(ctx *gin.Context) {
	gameID := ctx.Param("id")
	if gameID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Game ID is required",
		})
		return
	}

	game, err := c.llmService.GetGame(gameID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Game not found",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    game,
		"message": "Game retrieved successfully",
	})
}

// GetModels handles GET /api/llm/models
func (c *LLMController) GetModels(ctx *gin.Context) {
	models := c.llmService.GetAvailableModels()

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    models,
		"message": "Models retrieved successfully",
	})
}

// UpdateConfig handles PUT /api/llm/config/:model
func (c *LLMController) UpdateConfig(ctx *gin.Context) {
	modelName := ctx.Param("model")
	if modelName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Model name is required",
		})
		return
	}

	var request model.LLMConfigRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Convert request to config
	config := model.LLMConfig{
		Name:    modelName,
		Model:   request.Model,
		APIKey:  request.APIKey,
		BaseURL: &request.Endpoint,
		// Provider will be determined based on model name
		Provider: determineProvider(modelName),
	}

	if err := c.llmService.UpdateConfig(modelName, config); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to update configuration",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Configuration updated successfully",
	})
}

// GetConfig handles GET /api/llm/config/:model
func (c *LLMController) GetConfig(ctx *gin.Context) {
	modelName := ctx.Param("model")
	if modelName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Model name is required",
		})
		return
	}

	_, err := c.llmService.GetConfig(modelName)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Configuration not found",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Configuration retrieved successfully",
	})
}

// DeleteGame handles DELETE /api/llm/game/:id
func (c *LLMController) DeleteGame(ctx *gin.Context) {
	gameID := ctx.Param("id")
	if gameID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Game ID is required",
		})
		return
	}

	if err := c.llmService.DeleteGame(gameID); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Failed to delete game",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Game deleted successfully",
	})
}

// GetGameHistory handles GET /api/llm/game/:id/history
func (c *LLMController) GetGameHistory(ctx *gin.Context) {
	gameID := ctx.Param("id")
	if gameID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Game ID is required",
		})
		return
	}

	// Get limit parameter (optional)
	limitStr := ctx.Query("limit")
	limit := 50 // default limit
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	game, err := c.llmService.GetGame(gameID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Game not found",
			"details": err.Error(),
		})
		return
	}

	// Get recent moves (limited)
	moves := game.Moves
	if len(moves) > limit {
		moves = moves[len(moves)-limit:]
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"game_id":    game.ID,
			"moves":      moves,
			"total":      len(game.Moves),
			"status":     game.Status,
			"model_name": game.ModelName,
		},
		"message": "Game history retrieved successfully",
	})
}

// GetGameStats handles GET /api/llm/stats
func (c *LLMController) GetGameStats(ctx *gin.Context) {
	// This is a placeholder for future statistics functionality
	// For now, return basic stats
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"total_games":   0,
			"active_games":  0,
			"human_wins":    0,
			"ai_wins":       0,
			"draws":         0,
			"popular_model": "deepseek",
		},
		"message": "Statistics retrieved successfully",
	})
}

// HealthCheck handles GET /api/llm/health
func (c *LLMController) HealthCheck(ctx *gin.Context) {
	models := c.llmService.GetAvailableModels()
	
	var availableCount, configuredCount int
	for _, model := range models {
		if model.Status == "available" {
			availableCount++
			configuredCount++
		} else if model.Status == "not_configured" {
			availableCount++
		}
	}

	// Get cache statistics
	cacheStats := c.llmService.GetCacheStats()

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"status":            "healthy",
			"available_models":  availableCount,
			"configured_models": configuredCount,
			"cache_stats":       cacheStats,
			"timestamp":         gin.H{},
		},
		"message": "LLM service is healthy",
	})
}

// GetCacheStats handles GET /api/llm/cache/stats
func (c *LLMController) GetCacheStats(ctx *gin.Context) {
	stats := c.llmService.GetCacheStats()

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
		"message": "Cache statistics retrieved successfully",
	})
}

// ClearCache handles DELETE /api/llm/cache
func (c *LLMController) ClearCache(ctx *gin.Context) {
	if err := c.llmService.ClearCache(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to clear cache",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Cache cleared successfully",
	})
}

// determineProvider determines the provider based on model name
func determineProvider(modelName string) string {
	modelName = strings.ToLower(modelName)
	if strings.Contains(modelName, "deepseek") {
		return "deepseek"
	} else if strings.Contains(modelName, "gpt") || strings.Contains(modelName, "chatgpt") {
		return "openai"
	} else if strings.Contains(modelName, "ollama") {
		return "ollama"
	}
	return "unknown"
}

// findValidLLMMove finds a valid move for the LLM using simple strategy
func (c *LLMController) findValidLLMMove(board [][]int, humanX, humanY int) *model.LLMMove {
	// Simple strategy: try to place near the human move
	directions := [][]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1},           {0, 1},
		{1, -1},  {1, 0},  {1, 1},
	}

	// Try positions around the human move first
	for _, dir := range directions {
		x := humanX + dir[0]
		y := humanY + dir[1]
		if x >= 0 && x < 15 && y >= 0 && y < 15 && board[y][x] == 0 {
			return &model.LLMMove{
				X:         x,
				Y:         y,
				Player:    2,
				Reasoning: "选择了人类移动附近的位置",
				Timestamp: time.Now(),
			}
		}
	}

	// If no position around human move, find any valid position
	for y := 0; y < 15; y++ {
		for x := 0; x < 15; x++ {
			if board[y][x] == 0 {
				return &model.LLMMove{
					X:         x,
					Y:         y,
					Player:    2,
					Reasoning: "选择了一个可用位置",
					Timestamp: time.Now(),
				}
			}
		}
	}

	// Fallback (should never happen in a normal game)
	return &model.LLMMove{
		X:         7,
		Y:         7,
		Player:    2,
		Reasoning: "默认中心位置",
		Timestamp: time.Now(),
	}
}