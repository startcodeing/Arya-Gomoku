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
  id: string
  gameId: string
  x: number
  y: number
  player: number
  reasoning: string
  confidence: number
  timestamp: string
}

// LLM游戏信息
export interface LLMGame {
  id: string
  modelName: string
  difficulty: string
  status: string
  currentPlayer: number
  board: number[][]
  moves: LLMMove[]
  createdAt: string
  updatedAt: string
}

// LLM模型信息
export interface LLMModel {
  name: string
  displayName: string
  provider: string
  requiresApiKey: boolean
  defaultParams: Record<string, any>
  status: 'available' | 'unavailable' | 'error'
}

// LLM配置
export interface LLMConfig {
  modelName: string
  apiKey: string
  endpoint: string
  parameters: Record<string, any>
  enabled: boolean
}

// LLM配置请求
export interface LLMConfigRequest {
  model: string
  apiKey?: string
  endpoint?: string
  parameters?: Record<string, any>
}

// LLM配置响应
export interface LLMConfigResponse {
  success: boolean
  message?: string
  error?: string
}

// LLM移动请求
export interface LLMMoveRequest {
  gameId: string
  move: {
    x: number
    y: number
  }
}

// LLM移动响应
export interface LLMMoveResponse {
  success: boolean
  move?: LLMMove
  reasoning?: string
  gameStatus: string
  error?: string
}

// LLM开始游戏请求
export interface LLMStartGameRequest {
  modelName: string
  difficulty?: string
}

// LLM游戏历史响应
export interface LLMGameHistoryResponse {
  gameId: string
  moves: LLMMove[]
  total: number
  status: string
  modelName: string
}

// LLM统计信息
export interface LLMStats {
  totalGames: number
  activeGames: number
  humanWins: number
  aiWins: number
  draws: number
  popularModel: string
}

// LLM健康检查响应
export interface LLMHealthResponse {
  status: string
  availableModels: number
  configuredModels: number
  timestamp: Record<string, any>
}