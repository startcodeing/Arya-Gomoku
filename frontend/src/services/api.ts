import axios from 'axios'
import type { AIRequest, AIResponse, ApiResponse } from '../types/game'

// 创建axios实例
const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    console.log('发送请求:', config.method?.toUpperCase(), config.url, config.data)
    return config
  },
  (error) => {
    console.error('请求错误:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    console.log('收到响应:', response.status, response.data)
    return response
  },
  (error) => {
    console.error('响应错误:', error.response?.status, error.response?.data || error.message)
    return Promise.reject(error)
  }
)

// AI相关API
export const aiApi = {
  // 获取AI移动
  async getMove(request: AIRequest): Promise<AIResponse> {
    try {
      const response = await api.post<AIResponse>('/ai/move', request)
      return response.data
    } catch (error: any) {
      console.error('获取AI移动失败:', error)
      throw new Error(error.response?.data?.error || '获取AI移动失败')
    }
  },

  // 获取游戏状态
  async getStatus(): Promise<ApiResponse> {
    try {
      const response = await api.get('/ai/status')
      return response.data
    } catch (error: any) {
      console.error('获取游戏状态失败:', error)
      throw new Error(error.response?.data?.error || '获取游戏状态失败')
    }
  },

  // 重置游戏
  async resetGame(): Promise<ApiResponse> {
    try {
      const response = await api.post('/ai/reset')
      return response.data
    } catch (error: any) {
      console.error('重置游戏失败:', error)
      throw new Error(error.response?.data?.error || '重置游戏失败')
    }
  }
}

// 匹配相关API（预留给未来的PVP功能）
export const matchApi = {
  // 开始匹配
  async startMatch(playerId: string, gameMode: 'ai' | 'pvp' = 'ai'): Promise<ApiResponse> {
    try {
      const response = await api.post('/match/start', { playerId, gameMode })
      return response.data
    } catch (error: any) {
      console.error('开始匹配失败:', error)
      throw new Error(error.response?.data?.error || '开始匹配失败')
    }
  },

  // 加入房间
  async joinRoom(roomId: string, playerId: string): Promise<ApiResponse> {
    try {
      const response = await api.post('/match/join', { roomId, playerId })
      return response.data
    } catch (error: any) {
      console.error('加入房间失败:', error)
      throw new Error(error.response?.data?.error || '加入房间失败')
    }
  },

  // 获取房间状态
  async getRoomStatus(roomId: string): Promise<ApiResponse> {
    try {
      const response = await api.get(`/match/status/${roomId}`)
      return response.data
    } catch (error: any) {
      console.error('获取房间状态失败:', error)
      throw new Error(error.response?.data?.error || '获取房间状态失败')
    }
  },

  // 在房间中移动
  async makeMove(roomId: string, playerId: string, x: number, y: number): Promise<ApiResponse> {
    try {
      const response = await api.post(`/match/${roomId}/move`, { playerId, x, y })
      return response.data
    } catch (error: any) {
      console.error('移动失败:', error)
      throw new Error(error.response?.data?.error || '移动失败')
    }
  },

  // 离开房间
  async leaveRoom(roomId: string, playerId: string): Promise<ApiResponse> {
    try {
      const response = await api.post(`/match/${roomId}/leave`, { playerId })
      return response.data
    } catch (error: any) {
      console.error('离开房间失败:', error)
      throw new Error(error.response?.data?.error || '离开房间失败')
    }
  },

  // 获取活跃房间数
  async getActiveRooms(): Promise<ApiResponse> {
    try {
      const response = await api.get('/match/rooms')
      return response.data
    } catch (error: any) {
      console.error('获取活跃房间失败:', error)
      throw new Error(error.response?.data?.error || '获取活跃房间失败')
    }
  }
}

export default api