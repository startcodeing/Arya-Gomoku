import { api } from './api'
import type {
  CreateRoomRequest,
  CreateRoomResponse,
  JoinRoomRequest,
  JoinRoomResponse,
  RoomListResponse,
  RoomResponse,
  GameResponse
} from '../types/pvp'

export const pvpApi = {
  async getActiveRooms(): Promise<RoomListResponse> {
    const { data } = await api.get('/rooms')
    return data
  },

  async createRoom(payload: CreateRoomRequest): Promise<CreateRoomResponse> {
    const { data } = await api.post('/rooms', payload)
    return data
  },

  async joinRoom(roomId: string, payload: JoinRoomRequest): Promise<JoinRoomResponse> {
    const { data } = await api.post(`/rooms/${roomId}/join`, payload)
    return data
  },

  async getRoom(roomId: string): Promise<RoomResponse> {
    const { data } = await api.get(`/rooms/${roomId}`)
    return data
  },

  async getRoomInfo(roomId: string): Promise<RoomResponse> {
    const { data } = await api.get(`/rooms/${roomId}`)
    return data
  },

  async startGame(roomId: string): Promise<RoomResponse> {
    const { data } = await api.post(`/rooms/${roomId}/start`)
    return data
  },

  async makeMove(roomId: string, payload: { x: number, y: number, playerId: string }): Promise<RoomResponse> {
    const { data } = await api.post(`/rooms/${roomId}/move`, payload)
    return data
  },

  async getGame(roomId: string): Promise<GameResponse> {
    const { data } = await api.get(`/rooms/${roomId}/game`)
    return data
  },

  async leaveRoom(roomId: string, playerId: string): Promise<void> {
    await api.post(`/rooms/${roomId}/leave`, { playerId })
  },
}