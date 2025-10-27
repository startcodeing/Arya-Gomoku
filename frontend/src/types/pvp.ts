// PVP Game Types

export interface Player {
  id: string
  name: string
  isReady: boolean
  color?: 'black' | 'white'
}

export interface Room {
  id: string
  name: string
  status: 'waiting' | 'playing' | 'finished'
  players: Player[]
  maxPlayers: number
  createdAt: string
  game?: PVPGame
  // creatorId is present from backend but optional in TS
  creatorId?: string
}

export interface PVPGame {
  id: string
  roomId: string
  board: number[][]
  currentPlayer: string
  status: 'playing' | 'finished'
  winner?: string
  moves: Move[]
  startedAt: string
  finishedAt?: string
}

export interface Move {
  x: number
  y: number
  playerId: string
  timestamp: string
}

export interface CreateRoomRequest {
  roomName: string
  playerName: string
  maxPlayers: number
}

export interface JoinRoomRequest {
  playerName: string
}

export interface MakeMoveRequest {
  x: number
  y: number
  playerId: string
}

export interface CreateRoomResponse {
  room: Room
  // player may be returned by backend in some routes; optional here
  player?: Player
}

export interface JoinRoomResponse {
  room: Room
  player: Player
}

export interface WebSocketMessage {
  type: 'room_update' | 'game_start' | 'game_update' | 'game_end' | 'player_joined' | 'player_left' | 'chat_message' | 'error'
  data: any
  timestamp: string
}

export interface ChatMessage {
  id: string
  playerId: string
  playerName: string
  message: string
  timestamp: string
}

export interface GameResult {
  winner?: string
  winnerName?: string
  reason: 'win' | 'draw' | 'disconnect'
  finalBoard: number[][]
  moves: Move[]
  duration: number
}

// WebSocket Events
export interface RoomUpdateEvent {
  room: Room
}

export interface GameStartEvent {
  game: PVPGame
}

export interface GameUpdateEvent {
  game: PVPGame
  lastMove: Move
}

export interface GameEndEvent {
  result: GameResult
}

export interface PlayerJoinedEvent {
  player: Player
  room: Room
}

export interface PlayerLeftEvent {
  playerId: string
  room: Room
}

export interface ChatMessageEvent {
  message: ChatMessage
}

export interface ErrorEvent {
  message: string
  code?: string
}

// Store State
export interface PVPState {
  currentRoom: Room | null
  currentPlayer: Player | null
  isConnected: boolean
  isLoading: boolean
  error: string | null
  chatMessages: ChatMessage[]
}

// API Response Types
export interface ApiResponse<T = any> {
  success: boolean
  data?: T
  message?: string
  error?: string
}

export interface RoomListResponse {
  rooms: Room[]
}

export interface RoomResponse {
  room: Room
}

export interface GameResponse {
  game: PVPGame
}

// Connection Status
export type ConnectionStatus = 'disconnected' | 'connecting' | 'connected' | 'reconnecting' | 'error'

// Game Constants
export const BOARD_SIZE = 15
export const WIN_CONDITION = 5

// Player Colors
export const PLAYER_COLORS = {
  BLACK: 1,
  WHITE: 2
} as const

export type PlayerColor = typeof PLAYER_COLORS[keyof typeof PLAYER_COLORS]

// Room Status
export const ROOM_STATUS = {
  WAITING: 'waiting',
  PLAYING: 'playing',
  FINISHED: 'finished'
} as const

export type RoomStatus = typeof ROOM_STATUS[keyof typeof ROOM_STATUS]

// Game Status
export const GAME_STATUS = {
  PLAYING: 'playing',
  FINISHED: 'finished'
} as const

export type GameStatus = typeof GAME_STATUS[keyof typeof GAME_STATUS]

// WebSocket Message Types
export const WS_MESSAGE_TYPES = {
  ROOM_UPDATE: 'room_update',
  GAME_START: 'game_start',
  GAME_UPDATE: 'game_update',
  GAME_END: 'game_end',
  PLAYER_JOINED: 'player_joined',
  PLAYER_LEFT: 'player_left',
  CHAT_MESSAGE: 'chat_message',
  ERROR: 'error'
} as const

export type WSMessageType = typeof WS_MESSAGE_TYPES[keyof typeof WS_MESSAGE_TYPES]