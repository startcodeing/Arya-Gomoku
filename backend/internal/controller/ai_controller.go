// Package controller handles HTTP requests and responses for the Gomoku API
// This package contains the REST API endpoints for AI gameplay and game management
package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"gomoku-backend/internal/model"
	"gomoku-backend/internal/service"
)

// AIController handles AI-related HTTP requests
type AIController struct {
	aiService *service.AIService
}

// NewAIController creates a new AI controller instance
func NewAIController() *AIController {
	return &AIController{
		aiService: service.NewAIService(),
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
				"error": "Board row " + string(rune(i)) + " must have 15 columns",
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
	
	// Get AI move
	aiMove := ac.aiService.GetAIMove(request.Board, request.LastMove)
	
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
	
	c.JSON(http.StatusOK, response)
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