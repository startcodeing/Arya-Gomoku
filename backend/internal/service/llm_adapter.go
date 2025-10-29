// Package service contains LLM adapter interfaces and implementations
// This file defines the adapter pattern for different LLM providers
package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gomoku-backend/internal/model"
)

// LLMAdapter interface for different LLM providers
type LLMAdapter interface {
	GetMove(board [][]int, lastMove model.Move, config model.LLMConfig) (*model.LLMMove, error)
	ValidateConfig(config model.LLMConfig) error
	GetModelInfo() model.LLMModel
}

// DeepSeekAdapter implements LLMAdapter for DeepSeek API
type DeepSeekAdapter struct {
	httpClient *http.Client
}

// NewDeepSeekAdapter creates a new DeepSeek adapter
func NewDeepSeekAdapter() *DeepSeekAdapter {
	return &DeepSeekAdapter{
		httpClient: &http.Client{
			Timeout: 90 * time.Second, // 增加超时时间以适应大模型响应时间
		},
	}
}

// GetMove implements LLMAdapter interface for DeepSeek
func (d *DeepSeekAdapter) GetMove(board [][]int, lastMove model.Move, config model.LLMConfig) (*model.LLMMove, error) {
	if config.APIKey == "" {
		return nil, errors.New("DeepSeek API key is required")
	}

	// Build the prompt for DeepSeek
	prompt := d.buildGamePrompt(board, lastMove)

	// Prepare request payload
	requestBody := map[string]interface{}{
		"model": "deepseek-chat",
		"messages": []map[string]interface{}{
			{
				"role":    "system",
				"content": "你是一个五子棋专家。请分析当前棋局并选择最佳落子位置。你的回答必须严格按照JSON格式：{\"x\": 数字, \"y\": 数字, \"reasoning\": \"你的分析过程\"}",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": 0.7,
		"max_tokens":  1000,
	}

	// Add custom parameters if provided
	if params, ok := config.Parameters["temperature"]; ok {
		requestBody["temperature"] = params
	}
	if params, ok := config.Parameters["max_tokens"]; ok {
		requestBody["max_tokens"] = params
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	// Create HTTP request
	endpoint := config.Endpoint
	if endpoint == "" {
		endpoint = "https://api.deepseek.com/v1/chat/completions"
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.APIKey)

	// Send request
	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if len(response.Choices) == 0 {
		return nil, errors.New("no response from DeepSeek API")
	}

	// Parse the move from response
	move, err := d.parseMove(response.Choices[0].Message.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse move: %v", err)
	}

	// Validate move
	if !d.isValidMove(board, move.X, move.Y) {
		// If invalid, find a random valid move
		validMove := d.findValidMove(board)
		move.X = validMove.X
		move.Y = validMove.Y
		move.Reasoning = "原始选择无效，选择了一个有效位置：" + move.Reasoning
	}

	move.Player = 2 // LLM is player 2
	move.Timestamp = time.Now()
	move.GameID = "" // Will be set by service layer

	return move, nil
}

// ValidateConfig validates DeepSeek configuration
func (d *DeepSeekAdapter) ValidateConfig(config model.LLMConfig) error {
	if config.APIKey == "" {
		return errors.New("API key is required for DeepSeek")
	}
	if !strings.HasPrefix(config.APIKey, "sk-") {
		return errors.New("invalid DeepSeek API key format")
	}
	return nil
}

// GetModelInfo returns DeepSeek model information
func (d *DeepSeekAdapter) GetModelInfo() model.LLMModel {
	return model.LLMModel{
		Name:           "deepseek",
		DisplayName:    "DeepSeek",
		Provider:       "deepseek",
		RequiresAPIKey: true,
		DefaultParams: map[string]interface{}{
			"temperature": 0.7,
			"max_tokens":  1000,
		},
		Status: "available",
	}
}

// buildGamePrompt creates a prompt describing the current game state
func (d *DeepSeekAdapter) buildGamePrompt(board [][]int, lastMove model.Move) string {
	var prompt strings.Builder

	prompt.WriteString("当前五子棋棋局状态：\n")
	prompt.WriteString("棋盘大小：15x15\n")
	prompt.WriteString("玩家标记：1=人类玩家(黑子), 2=AI(白子), 0=空位\n\n")

	// Add board state
	prompt.WriteString("棋盘状态：\n")
	for y := 0; y < 15; y++ {
		for x := 0; x < 15; x++ {
			if x > 0 {
				prompt.WriteString(" ")
			}
			prompt.WriteString(strconv.Itoa(board[y][x]))
		}
		prompt.WriteString("\n")
	}

	prompt.WriteString(fmt.Sprintf("\n上一步棋：人类玩家在位置(%d, %d)下了黑子\n", lastMove.X, lastMove.Y))
	prompt.WriteString("现在轮到你(AI)下白子。\n\n")
	prompt.WriteString("请分析棋局并选择最佳落子位置。考虑因素包括：\n")
	prompt.WriteString("1. 是否能形成五连获胜\n")
	prompt.WriteString("2. 是否需要阻止对手形成五连\n")
	prompt.WriteString("3. 是否能形成活三、活四等威胁\n")
	prompt.WriteString("4. 整体战略布局\n\n")
	prompt.WriteString("请返回JSON格式：{\"x\": 横坐标(0-14), \"y\": 纵坐标(0-14), \"reasoning\": \"你的分析过程\"}")

	return prompt.String()
}

// parseMove parses the LLM response to extract move information
func (d *DeepSeekAdapter) parseMove(content string) (*model.LLMMove, error) {
	// Try to find JSON in the response
	start := strings.Index(content, "{")
	end := strings.LastIndex(content, "}")

	if start == -1 || end == -1 || start >= end {
		return nil, errors.New("no valid JSON found in response")
	}

	jsonStr := content[start : end+1]

	var moveData struct {
		X         int    `json:"x"`
		Y         int    `json:"y"`
		Reasoning string `json:"reasoning"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &moveData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return &model.LLMMove{
		X:          moveData.X,
		Y:          moveData.Y,
		Reasoning:  moveData.Reasoning,
		Confidence: 0.8, // Default confidence for DeepSeek
	}, nil
}

// isValidMove checks if a move is valid
func (d *DeepSeekAdapter) isValidMove(board [][]int, x, y int) bool {
	if x < 0 || x >= 15 || y < 0 || y >= 15 {
		return false
	}
	return board[y][x] == 0
}

// findValidMove finds a random valid move as fallback
func (d *DeepSeekAdapter) findValidMove(board [][]int) model.Move {
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

// ChatGPTAdapter implements LLMAdapter for OpenAI ChatGPT (placeholder)
type ChatGPTAdapter struct {
	httpClient *http.Client
}

// NewChatGPTAdapter creates a new ChatGPT adapter
func NewChatGPTAdapter() *ChatGPTAdapter {
	return &ChatGPTAdapter{
		httpClient: &http.Client{
			Timeout: 90 * time.Second, // 增加超时时间以适应大模型响应时间
		},
	}
}

// GetMove implements LLMAdapter interface for ChatGPT (placeholder implementation)
func (c *ChatGPTAdapter) GetMove(board [][]int, lastMove model.Move, config model.LLMConfig) (*model.LLMMove, error) {
	return nil, errors.New("ChatGPT adapter not implemented yet")
}

// ValidateConfig validates ChatGPT configuration
func (c *ChatGPTAdapter) ValidateConfig(config model.LLMConfig) error {
	if config.APIKey == "" {
		return errors.New("API key is required for ChatGPT")
	}
	return nil
}

// GetModelInfo returns ChatGPT model information
func (c *ChatGPTAdapter) GetModelInfo() model.LLMModel {
	return model.LLMModel{
		Name:           "chatgpt",
		DisplayName:    "ChatGPT",
		Provider:       "openai",
		RequiresAPIKey: true,
		DefaultParams: map[string]interface{}{
			"model":       "gpt-3.5-turbo",
			"temperature": 0.7,
			"max_tokens":  1000,
		},
		Status: "unavailable", // Not implemented yet
	}
}

// OllamaAdapter implements LLMAdapter for Ollama local models (placeholder)
type OllamaAdapter struct {
	httpClient *http.Client
}

// NewOllamaAdapter creates a new Ollama adapter
func NewOllamaAdapter() *OllamaAdapter {
	return &OllamaAdapter{
		httpClient: &http.Client{
			Timeout: 60 * time.Second, // Longer timeout for local models
		},
	}
}

// GetMove implements LLMAdapter interface for Ollama
func (o *OllamaAdapter) GetMove(board [][]int, lastMove model.Move, config model.LLMConfig) (*model.LLMMove, error) {
	// Build the prompt for Ollama
	prompt := o.buildGamePrompt(board, lastMove)

	// Prepare request payload for Ollama API
	requestBody := map[string]interface{}{
		"model":  config.ModelName,
		"prompt": prompt,
		"stream": false,
		"options": map[string]interface{}{
			"temperature": 0.7,
			"num_predict": 1000,
		},
	}

	// Add custom parameters if provided
	if params, ok := config.Parameters["temperature"]; ok {
		requestBody["options"].(map[string]interface{})["temperature"] = params
	}
	if params, ok := config.Parameters["num_predict"]; ok {
		requestBody["options"].(map[string]interface{})["num_predict"] = params
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	// Create HTTP request
	endpoint := config.Endpoint
	if endpoint == "" {
		endpoint = "http://localhost:11434/api/generate"
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := o.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var response struct {
		Response string `json:"response"`
		Done     bool   `json:"done"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if !response.Done {
		return nil, errors.New("incomplete response from Ollama API")
	}

	// Parse the move from response
	move, err := o.parseMove(response.Response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse move: %v", err)
	}

	// Validate move
	if !o.isValidMove(board, move.X, move.Y) {
		// If invalid, find a random valid move
		validMove := o.findValidMove(board)
		move.X = validMove.X
		move.Y = validMove.Y
		move.Reasoning = "原始选择无效，选择了一个有效位置：" + move.Reasoning
	}

	move.Player = 2 // LLM is player 2
	move.Timestamp = time.Now()
	move.GameID = "" // Will be set by service layer

	return move, nil
}

// ValidateConfig validates Ollama configuration
func (o *OllamaAdapter) ValidateConfig(config model.LLMConfig) error {
	// Ollama doesn't require API key, but needs model name
	if config.ModelName == "" {
		return errors.New("model name is required for Ollama")
	}
	// Endpoint is optional, will use default if not provided
	return nil
}

// GetModelInfo returns Ollama model information
func (o *OllamaAdapter) GetModelInfo() model.LLMModel {
	return model.LLMModel{
		Name:           "ollama",
		DisplayName:    "Ollama Local",
		Provider:       "ollama",
		RequiresAPIKey: true,
		DefaultParams: map[string]interface{}{
			"model":       "llama2",
			"temperature": 0.7,
		},
		Status: "available", // Now implemented
	}
}

// buildGamePrompt creates a prompt describing the current game state for Ollama
func (o *OllamaAdapter) buildGamePrompt(board [][]int, lastMove model.Move) string {
	var prompt strings.Builder

	prompt.WriteString("当前五子棋棋局状态：\n")
	prompt.WriteString("棋盘大小：15x15\n")
	prompt.WriteString("玩家标记：1=人类玩家(黑子), 2=AI(白子), 0=空位\n\n")

	// Add board state
	prompt.WriteString("棋盘状态：\n")
	for y := 0; y < 15; y++ {
		for x := 0; x < 15; x++ {
			if x > 0 {
				prompt.WriteString(" ")
			}
			prompt.WriteString(strconv.Itoa(board[y][x]))
		}
		prompt.WriteString("\n")
	}

	prompt.WriteString(fmt.Sprintf("\n上一步棋：人类玩家在位置(%d, %d)下了黑子\n", lastMove.X, lastMove.Y))
	prompt.WriteString("现在轮到你(AI)下白子。\n\n")
	prompt.WriteString("请分析棋局并选择最佳落子位置。考虑因素包括：\n")
	prompt.WriteString("1. 是否能形成五连获胜\n")
	prompt.WriteString("2. 是否需要阻止对手形成五连\n")
	prompt.WriteString("3. 是否能形成活三、活四等威胁\n")
	prompt.WriteString("4. 整体战略布局\n\n")
	prompt.WriteString("请返回JSON格式：{\"x\": 横坐标(0-14), \"y\": 纵坐标(0-14), \"reasoning\": \"你的分析过程\"}")

	return prompt.String()
}

// parseMove parses the Ollama response to extract move information
func (o *OllamaAdapter) parseMove(content string) (*model.LLMMove, error) {
	// Try to find JSON in the response
	start := strings.Index(content, "{")
	end := strings.LastIndex(content, "}")

	if start == -1 || end == -1 || start >= end {
		return nil, errors.New("no valid JSON found in response")
	}

	jsonStr := content[start : end+1]

	var moveData struct {
		X         int    `json:"x"`
		Y         int    `json:"y"`
		Reasoning string `json:"reasoning"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &moveData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return &model.LLMMove{
		X:          moveData.X,
		Y:          moveData.Y,
		Reasoning:  moveData.Reasoning,
		Confidence: 0.8, // Default confidence for Ollama
	}, nil
}

// isValidMove checks if a move is valid
func (o *OllamaAdapter) isValidMove(board [][]int, x, y int) bool {
	if x < 0 || x >= 15 || y < 0 || y >= 15 {
		return false
	}
	return board[y][x] == 0
}

// findValidMove finds a random valid move as fallback
func (o *OllamaAdapter) findValidMove(board [][]int) model.Move {
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
