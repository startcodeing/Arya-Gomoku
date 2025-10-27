import axios from 'axios'
import type {
  CreateRoomRequest,
  CreateRoomResponse,
  JoinRoomRequest,
  JoinRoomResponse,
  RoomListResponse,
  RoomResponse,
  GameResponse
} from '../types/pvp'

const apiClient = axios.create({
  baseURL: '/api',
})

export const pvpApi = {
  async getActiveRooms(): Promise<RoomListResponse> {
    const { data } = await apiClient.get('/rooms')
    return data
  },

  async createRoom(payload: CreateRoomRequest): Promise<CreateRoomResponse> {
    const { data } = await apiClient.post('/rooms', payload)
    return data
  },

  async joinRoom(roomId: string, payload: JoinRoomRequest): Promise<JoinRoomResponse> {
    const { data } = await apiClient.post(`/rooms/${roomId}/join`, payload)
    return data
  },

  async getRoom(roomId: string): Promise<RoomResponse> {
    const { data } = await apiClient.get(`/rooms/${roomId}`)
    return data
  },

  async getRoomInfo(roomId: string): Promise<RoomResponse> {
    const { data } = await apiClient.get(`/rooms/${roomId}`)
    return data
  },

  async startGame(roomId: string): Promise<RoomResponse> {
    const { data } = await apiClient.post(`/rooms/${roomId}/start`)
    return data
  },

  async makeMove(roomId: string, payload: { x: number, y: number, playerId: string }): Promise<RoomResponse> {
    const { data } = await apiClient.post(`/rooms/${roomId}/move`, payload)
    return data
  },

  async getGame(roomId: string): Promise<GameResponse> {
    const { data } = await apiClient.get(`/rooms/${roomId}/game`)
    return data
  },
}