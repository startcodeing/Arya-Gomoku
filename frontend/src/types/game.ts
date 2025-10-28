// 游戏相关类型定义

// 玩家类型
export enum Player {
  NONE = 0,
  HUMAN = 1,
  AI = 2
}

// 游戏状态
export enum GameStatus {
  PLAYING = 'playing',
  HUMAN_WIN = 'human_win',
  AI_WIN = 'ai_win',
  DRAW = 'draw'
}

// 棋子位置
export interface Position {
  x: number
  y: number
}

// 移动信息
export interface Move {
  x: number
  y: number
  player: Player
}

// 棋盘状态
export interface BoardState {
  board: number[][]
  currentPlayer: Player
  gameStatus: GameStatus
  winner: Player
  lastMove?: Position
}

// AI移动请求
export interface AIRequest {
  board: number[][]
  player: number
  lastMove: Position
}

// AI移动响应
export interface AIResponse {
  aiMove: {
    x: number
    y: number
    score: number
  }
  gameStatus: string
  winner: number
}

// API响应基础类型
export interface ApiResponse<T = any> {
  success?: boolean
  data?: T
  error?: string
  message?: string
}

// ===== LLM 相关类型定义 =====

// LLM游戏状态
export enum LLMGameStatus {
  PLAYING = 'playing',
  HUMAN_WIN = 'human_win',
  AI_WIN = 'ai_win',
  DRAW = 'draw'
}

// LLM移动信息
export interface LLMMove {
  x: number
  y: number
  player: number
  reasoning: string
  confidence: number
  timestamp: string
  game_id: string
}

// LLM游戏信息
export interface LLMGame {
  game_id: string
  model_name: string
  status: LLMGameStatus
  board: number[][]
  moves: Move[]
  start_time: string
  end_time?: string
}

// LLM模型信息
export interface LLMModel {
  name: string
  display_name: string
  provider: string
  requires_api_key: boolean
  default_params: Record<string, any>
  status: 'available' | 'not_configured' | 'unavailable'
}

// LLM配置
export interface LLMConfig {
  model_name: string
  api_key: string
  endpoint: string
  parameters: Record<string, any>
}

// LLM配置请求
export interface LLMConfigRequest {
  api_key: string
  endpoint?: string
  parameters?: Record<string, any>
}

// LLM配置响应
export interface LLMConfigResponse {
  model_name: string
  endpoint: string
  parameters: Record<string, any>
  has_api_key: boolean
}

// LLM移动请求
export interface LLMMoveRequest {
  game_id: string
  move: {
    x: number
    y: number
  }
}

// LLM移动响应
export interface LLMMoveResponse {
  move?: LLMMove
  game_status: LLMGameStatus
  message: string
}

// LLM开始游戏请求
export interface LLMStartGameRequest {
  model_name: string
}

// LLM游戏历史响应
export interface LLMGameHistoryResponse {
  game_id: string
  moves: Move[]
  total: number
  status: LLMGameStatus
  model_name: string
}

// LLM统计信息
export interface LLMStats {
  total_games: number
  active_games: number
  human_wins: number
  ai_wins: number
  draws: number
  popular_model: string
}

// LLM健康检查响应
export interface LLMHealthResponse {
  status: string
  available_models: number
  configured_models: number
  timestamp: Record<string, any>
}