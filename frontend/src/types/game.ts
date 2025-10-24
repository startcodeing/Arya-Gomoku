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