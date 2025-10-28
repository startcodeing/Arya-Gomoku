// Package service contains LLM game service implementation
// This file implements the core LLM game logic and model management
package service

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"gomoku-backend/internal/model"
)

// LLMService manages LLM games and model interactions
type LLMService struct {
	adapters map[string]LLMAdapter
	games    map[string]*model.LLMGame
	configs  map[string]model.LLMConfig
	mutex    sync.RWMutex
}

// NewLLMService creates a new LLM service instance
func NewLLMService() *LLMService {
	service := &LLMService{
		adapters: make(map[string]LLMAdapter),
		games:    make(map[string]*model.LLMGame),
		configs:  make(map[string]model.LLMConfig),
	}

	// Register available adapters
	service.registerAdapters()
	
	// Set default configurations
	service.setDefaultConfigs()

	return service
}

// registerAdapters registers all available LLM adapters
func (s *LLMService) registerAdapters() {
	s.adapters["deepseek"] = NewDeepSeekAdapter()
	s.adapters["chatgpt"] = NewChatGPTAdapter()
	s.adapters["ollama"] = NewOllamaAdapter()
}

// setDefaultConfigs sets default configurations for each model
func (s *LLMService) setDefaultConfigs() {
	// DeepSeek default config
	s.configs["deepseek"] = model.LLMConfig{
		ModelName:  "deepseek",
		APIKey:     "", // To be set by user
		Endpoint:   "https://api.deepseek.com/v1/chat/completions",
		Parameters: map[string]interface{}{
			"temperature": 0.7,
			"max_tokens":  1000,
		},
	}

	// ChatGPT default config
	s.configs["chatgpt"] = model.LLMConfig{
		ModelName:  "chatgpt",
		APIKey:     "", // To be set by user
		Endpoint:   "https://api.openai.com/v1/chat/completions",
		Parameters: map[string]interface{}{
			"model":       "gpt-3.5-turbo",
			"temperature": 0.7,
			"max_tokens":  1000,
		},
	}

	// Ollama default config
	s.configs["ollama"] = model.LLMConfig{
		ModelName:  "ollama",
		APIKey:     "", // Not required for Ollama
		Endpoint:   "http://localhost:11434/api/generate",
		Parameters: map[string]interface{}{
			"model":       "llama2",
			"temperature": 0.7,
		},
	}
}

// StartGame creates a new LLM game
func (s *LLMService) StartGame(modelName string) (*model.LLMGame, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Validate model
	if _, exists := s.adapters[modelName]; !exists {
		return nil, fmt.Errorf("unsupported model: %s", modelName)
	}

	// Check if model is configured
	config, exists := s.configs[modelName]
	if !exists {
		return nil, fmt.Errorf("model %s is not configured", modelName)
	}

	// Validate configuration
	adapter := s.adapters[modelName]
	if err := adapter.ValidateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid configuration for %s: %v", modelName, err)
	}

	// Create new game
	game := model.NewLLMGame(modelName, "medium")
	s.games[game.ID] = game

	return game, nil
}

// MakeMove processes a human move and gets LLM response
func (s *LLMService) MakeMove(gameID string, humanMove model.Move) (*model.LLMResponse, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Get game
	game, exists := s.games[gameID]
	if !exists {
		return nil, errors.New("game not found")
	}

	if game.Status != "playing" {
		return nil, errors.New("game is not in playing state")
	}

	// Validate human move
	if !game.Board.IsValidMove(humanMove.X, humanMove.Y) {
		return nil, errors.New("invalid move position")
	}

	// Make human move
	humanMove.Player = 1
	llmMove := model.LLMMove{
		X:         humanMove.X,
		Y:         humanMove.Y,
		Player:    1,
		Timestamp: time.Now(),
	}
	game.AddMove(llmMove)
	game.Board.MakeMove(humanMove.X, humanMove.Y, 1)

	// Check if human wins
	if game.Board.CheckWin(humanMove.X, humanMove.Y, 1) {
		game.Status = "human_win"
		game.UpdatedAt = time.Now()
		return &model.LLMResponse{
			Move:       nil,
			GameStatus: game.Status,
			Reasoning:  "恭喜！你获胜了！",
		}, nil
	}

	// Check if board is full
	if game.Board.IsBoardFull() {
		game.Status = "draw"
		game.UpdatedAt = time.Now()
		return &model.LLMResponse{
			Move:       nil,
			GameStatus: game.Status,
			Reasoning:  "平局！",
		}, nil
	}

	// Get LLM move
	config := s.configs[game.ModelName]
	adapter := s.adapters[game.ModelName]
	
	llmMovePtr, err := adapter.GetMove(game.Board.Grid, humanMove, config)
	if err != nil {
		return nil, fmt.Errorf("failed to get LLM move: %v", err)
	}

	// Validate LLM move
	if !game.Board.IsValidMove(llmMovePtr.X, llmMovePtr.Y) {
		// Find a valid move as fallback
		validMove := s.findValidMove(game.Board.Grid)
		llmMovePtr.X = validMove.X
		llmMovePtr.Y = validMove.Y
		llmMovePtr.Reasoning = "原始选择无效，选择了一个有效位置"
	}

	// Make LLM move
	llmMovePtr.Player = 2
	llmMovePtr.GameID = gameID
	llmMovePtr.Timestamp = time.Now()
	
	game.AddMove(*llmMovePtr)
	game.Board.MakeMove(llmMovePtr.X, llmMovePtr.Y, 2)

	// Check if LLM wins
	if game.Board.CheckWin(llmMovePtr.X, llmMovePtr.Y, 2) {
		game.Status = "ai_win"
		game.UpdatedAt = time.Now()
		return &model.LLMResponse{
			Move:       llmMovePtr,
			GameStatus: game.Status,
			Reasoning:  "AI获胜！",
		}, nil
	}

	// Check if board is full after LLM move
	if game.Board.IsBoardFull() {
		game.Status = "draw"
		game.UpdatedAt = time.Now()
		return &model.LLMResponse{
			Move:       llmMovePtr,
			GameStatus: game.Status,
			Reasoning:  "平局！",
		}, nil
	}

	// Game continues
	return &model.LLMResponse{
		Move:       llmMovePtr,
		GameStatus: game.Status,
		Reasoning:  "游戏继续",
	}, nil
}

// GetGame retrieves a game by ID
func (s *LLMService) GetGame(gameID string) (*model.LLMGame, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	game, exists := s.games[gameID]
	if !exists {
		return nil, errors.New("game not found")
	}

	return game, nil
}

// GetAvailableModels returns list of available LLM models
func (s *LLMService) GetAvailableModels() []model.LLMModel {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var models []model.LLMModel
	for name, adapter := range s.adapters {
		modelInfo := adapter.GetModelInfo()
		
		// Check if model is configured
		config, exists := s.configs[name]
		if exists && config.APIKey != "" {
			modelInfo.Status = "available"
		} else if !modelInfo.RequiresAPIKey {
			modelInfo.Status = "available"
		} else {
			modelInfo.Status = "not_configured"
		}
		
		models = append(models, modelInfo)
	}

	return models
}

// UpdateConfig updates configuration for a specific model
func (s *LLMService) UpdateConfig(modelName string, config model.LLMConfig) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	adapter, exists := s.adapters[modelName]
	if !exists {
		return fmt.Errorf("unsupported model: %s", modelName)
	}

	// Validate configuration
	if err := adapter.ValidateConfig(config); err != nil {
		return fmt.Errorf("invalid configuration: %v", err)
	}

	// Update configuration
	config.ModelName = modelName
	s.configs[modelName] = config

	return nil
}

// GetConfig retrieves configuration for a specific model
func (s *LLMService) GetConfig(modelName string) (*model.LLMConfig, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	config, exists := s.configs[modelName]
	if !exists {
		return nil, fmt.Errorf("model %s not found", modelName)
	}

	// Don't expose API key in response
	safeCopy := config
	safeCopy.APIKey = ""
	
	return &safeCopy, nil
}

// DeleteGame removes a game from memory
func (s *LLMService) DeleteGame(gameID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.games[gameID]; !exists {
		return errors.New("game not found")
	}

	delete(s.games, gameID)
	return nil
}

// Helper methods

// isValidMove checks if a move is valid
func (s *LLMService) isValidMove(board [][]int, x, y int) bool {
	if x < 0 || x >= 15 || y < 0 || y >= 15 {
		return false
	}
	return board[y][x] == 0
}

// findValidMove finds a random valid move
func (s *LLMService) findValidMove(board [][]int) model.Move {
	for y := 0; y < 15; y++ {
		for x := 0; x < 15; x++ {
			if board[y][x] == 0 {
				return model.Move{X: x, Y: y}
			}
		}
	}
	// Should never reach here in a valid game
	return model.Move{X: 7, Y: 7}
}

// checkWin checks if a player has won
func (s *LLMService) checkWin(board [][]int, x, y, player int) bool {
	directions := [][2]int{{1, 0}, {0, 1}, {1, 1}, {1, -1}}
	
	for _, dir := range directions {
		count := 1
		
		// Check positive direction
		for i := 1; i < 5; i++ {
			nx, ny := x+dir[0]*i, y+dir[1]*i
			if nx < 0 || nx >= 15 || ny < 0 || ny >= 15 || board[ny][nx] != player {
				break
			}
			count++
		}
		
		// Check negative direction
		for i := 1; i < 5; i++ {
			nx, ny := x-dir[0]*i, y-dir[1]*i
			if nx < 0 || nx >= 15 || ny < 0 || ny >= 15 || board[ny][nx] != player {
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

// isBoardFull checks if the board is full
func (s *LLMService) isBoardFull(board [][]int) bool {
	for y := 0; y < 15; y++ {
		for x := 0; x < 15; x++ {
			if board[y][x] == 0 {
				return false
			}
		}
	}
	return true
}