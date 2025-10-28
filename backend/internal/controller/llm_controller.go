// Package controller contains LLM game controller implementation
// This file handles HTTP requests for LLM game functionality
package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gomoku-backend/internal/model"
	"gomoku-backend/internal/service"
)

// LLMController handles LLM game related HTTP requests
type LLMController struct {
	llmService *service.LLMService
}

// NewLLMController creates a new LLM controller instance
func NewLLMController(llmService *service.LLMService) *LLMController {
	return &LLMController{
		llmService: llmService,
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
		ModelName:  modelName,
		APIKey:     request.APIKey,
		Endpoint:   request.Endpoint,
		Parameters: request.Parameters,
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

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"status":            "healthy",
			"available_models":  availableCount,
			"configured_models": configuredCount,
			"timestamp":         gin.H{},
		},
		"message": "LLM service is healthy",
	})
}