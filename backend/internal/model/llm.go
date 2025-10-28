// Package model defines the core data structures for LLM battle functionality
// This file contains LLM-specific models for AI battle features
package model

import (
	"time"
)

// LLMGame represents a game session with LLM
type LLMGame struct {
	ID            string    `json:"id"`            // Unique game identifier
	ModelName     string    `json:"modelName"`     // LLM model name (deepseek, chatgpt, ollama)
	Difficulty    string    `json:"difficulty"`    // Difficulty level: easy, medium, hard
	Status        string    `json:"status"`        // Game status: playing, finished, error
	CurrentPlayer int       `json:"currentPlayer"` // Current player: 1=human, 2=LLM
	Board         *Board    `json:"board"`         // Current board state
	Moves         []LLMMove `json:"moves"`         // Move history
	CreatedAt     time.Time `json:"createdAt"`     // Game creation timestamp
	UpdatedAt     time.Time `json:"updatedAt"`     // Last update timestamp
}

// LLMMove represents a move made by LLM with reasoning
type LLMMove struct {
	ID         string    `json:"id"`                   // Unique move identifier
	GameID     string    `json:"gameId"`               // Associated game ID
	X          int       `json:"x"`                    // X coordinate (0-14)
	Y          int       `json:"y"`                    // Y coordinate (0-14)
	Player     int       `json:"player"`               // Player who made the move
	Reasoning  string    `json:"reasoning,omitempty"`  // LLM's reasoning process
	Confidence float64   `json:"confidence"`           // Move confidence score (0-1)
	Timestamp  time.Time `json:"timestamp"`            // Move timestamp
}

// LLMConfig represents LLM model configuration
type LLMConfig struct {
	ModelName  string                 `json:"modelName"`            // Model identifier
	APIKey     string                 `json:"apiKey,omitempty"`     // API key (encrypted)
	Endpoint   string                 `json:"endpoint,omitempty"`   // Custom API endpoint
	Parameters map[string]interface{} `json:"parameters,omitempty"` // Model-specific parameters
	Enabled    bool                   `json:"enabled"`              // Whether model is enabled
}

// LLMModel represents available LLM model information
type LLMModel struct {
	Name           string                 `json:"name"`           // Model identifier
	DisplayName    string                 `json:"displayName"`    // Human-readable name
	Provider       string                 `json:"provider"`       // Provider name (deepseek, openai, ollama)
	RequiresAPIKey bool                   `json:"requiresApiKey"` // Whether API key is required
	DefaultParams  map[string]interface{} `json:"defaultParams"`  // Default parameters
	Status         string                 `json:"status"`         // Model status: available, unavailable, error
}

// LLMRequest represents request payload for LLM move
type LLMRequest struct {
	Board      [][]int `json:"board"`                // Current board state
	LastMove   Move    `json:"lastMove"`             // Last move made
	Model      string  `json:"model"`                // LLM model to use
	Difficulty string  `json:"difficulty,omitempty"` // Difficulty level
}

// LLMResponse represents response from LLM move endpoint
type LLMResponse struct {
	Success    bool     `json:"success"`              // Request success status
	Move       *LLMMove `json:"move,omitempty"`       // LLM's chosen move
	Reasoning  string   `json:"reasoning,omitempty"`  // LLM's reasoning process
	GameStatus string   `json:"gameStatus"`           // Updated game status
	Error      string   `json:"error,omitempty"`      // Error message if any
}

// LLMModelsResponse represents response for available models
type LLMModelsResponse struct {
	Success bool       `json:"success"` // Request success status
	Models  []LLMModel `json:"models"`  // Available LLM models
	Error   string     `json:"error,omitempty"` // Error message if any
}

// LLMConfigRequest represents request for model configuration
type LLMConfigRequest struct {
	Model      string                 `json:"model"`                // Model name
	APIKey     string                 `json:"apiKey,omitempty"`     // API key
	Endpoint   string                 `json:"endpoint,omitempty"`   // Custom endpoint
	Parameters map[string]interface{} `json:"parameters,omitempty"` // Model parameters
}

// LLMConfigResponse represents response for model configuration
type LLMConfigResponse struct {
	Success bool   `json:"success"`           // Configuration success status
	Message string `json:"message,omitempty"` // Success/error message
	Error   string `json:"error,omitempty"`   // Error details if any
}

// NewLLMGame creates a new LLM game instance
func NewLLMGame(modelName, difficulty string) *LLMGame {
	return &LLMGame{
		ID:            generateGameID(),
		ModelName:     modelName,
		Difficulty:    difficulty,
		Status:        "playing",
		CurrentPlayer: 1, // Human starts first
		Board:         NewBoard(),
		Moves:         make([]LLMMove, 0),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// AddMove adds a move to the game history
func (g *LLMGame) AddMove(move LLMMove) {
	g.Moves = append(g.Moves, move)
	g.UpdatedAt = time.Now()
}

// GetLastMove returns the last move made in the game
func (g *LLMGame) GetLastMove() *LLMMove {
	if len(g.Moves) == 0 {
		return nil
	}
	return &g.Moves[len(g.Moves)-1]
}

// IsGameFinished checks if the game is finished
func (g *LLMGame) IsGameFinished() bool {
	return g.Status == "finished" || g.Status == "error"
}

// generateGameID generates a unique game identifier
func generateGameID() string {
	return "llm_" + time.Now().Format("20060102150405") + "_" + randomString(6)
}

// randomString generates a random string of specified length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(result)
}