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

// 请求拦截器 - 添加认证token
api.interceptors.request.use(
  (config) => {
    // 添加认证token
    const token = localStorage.getItem('auth_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    
    console.log('发送请求:', config.method?.toUpperCase(), config.url, config.data)
    return config
  },
  (error) => {
    console.error('请求错误:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器 - 处理认证错误
api.interceptors.response.use(
  (response) => {
    console.log('收到响应:', response.status, response.data)
    return response
  },
  async (error) => {
    const originalRequest = error.config
    
    // 处理401错误（未授权）
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true
      
      // 尝试刷新token
      const refreshToken = localStorage.getItem('refresh_token')
      if (refreshToken) {
        try {
          const response = await axios.post('/api/auth/refresh', {
            refresh_token: refreshToken
          })
          
          const { token } = response.data
          localStorage.setItem('auth_token', token)
          
          // 重新发送原始请求
          originalRequest.headers.Authorization = `Bearer ${token}`
          return api(originalRequest)
        } catch (refreshError) {
          // 刷新失败，清除认证信息并跳转到登录页
          localStorage.removeItem('auth_token')
          localStorage.removeItem('refresh_token')
          localStorage.removeItem('user_info')
          window.location.href = '/login'
          return Promise.reject(refreshError)
        }
      } else {
        // 没有refresh token，直接跳转到登录页
        localStorage.removeItem('auth_token')
        localStorage.removeItem('refresh_token')
        localStorage.removeItem('user_info')
        window.location.href = '/login'
      }
    }
    
    console.error('响应错误:', error.response?.status, error.response?.data || error.message)
    return Promise.reject(error)
  }
)

// AI相关API
export const aiApi = {
  // 获取AI移动
  async getMove(request: AIRequest, difficulty: 'easy' | 'medium' | 'hard' | 'expert' = 'medium', useEnhanced: boolean = true): Promise<AIResponse> {
    try {
      const params = new URLSearchParams({
        difficulty: difficulty,
        enhanced: useEnhanced.toString()
      })

      const response = await api.post<AIResponse>(`/ai/move?${params}`, request)
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
  },

  // 获取AI性能统计
  async getStats(): Promise<ApiResponse> {
    try {
      const response = await api.get('/ai/stats')
      return response.data
    } catch (error: any) {
      console.error('获取AI统计失败:', error)
      throw new Error(error.response?.data?.error || '获取AI统计失败')
    }
  },

  // 清除AI缓存
  async clearCache(): Promise<ApiResponse> {
    try {
      const response = await api.post('/ai/cache/clear')
      return response.data
    } catch (error: any) {
      console.error('清除AI缓存失败:', error)
      throw new Error(error.response?.data?.error || '清除AI缓存失败')
    }
  },

  // 获取难度等级
  async getDifficultyLevels(): Promise<ApiResponse> {
    try {
      const response = await api.get('/ai/difficulties')
      return response.data
    } catch (error: any) {
      console.error('获取难度等级失败:', error)
      throw new Error(error.response?.data?.error || '获取难度等级失败')
    }
  },

  // AI性能基准测试
  async benchmarkAI(difficulty: string = 'medium', moveCount: number = 10): Promise<ApiResponse> {
    try {
      const response = await api.post('/ai/benchmark', {
        difficulty,
        moveCount
      })
      return response.data
    } catch (error: any) {
      console.error('AI基准测试失败:', error)
      throw new Error(error.response?.data?.error || 'AI基准测试失败')
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
      const response = await api.post(`/rooms/${roomId}/leave`, { playerId })
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

// 用户相关API
export const userApi = {
  // 获取用户资料
  async getProfile(): Promise<ApiResponse> {
    try {
      const response = await api.get('/user/profile')
      return response.data
    } catch (error: any) {
      console.error('获取用户资料失败:', error)
      throw new Error(error.response?.data?.error || '获取用户资料失败')
    }
  },

  // 更新用户资料
  async updateProfile(data: { nickname?: string; email?: string }): Promise<ApiResponse> {
    try {
      const response = await api.put('/user/profile', data)
      return response.data
    } catch (error: any) {
      console.error('更新用户资料失败:', error)
      throw new Error(error.response?.data?.error || '更新用户资料失败')
    }
  },

  // 修改密码
  async changePassword(data: { old_password: string; new_password: string }): Promise<ApiResponse> {
    try {
      const response = await api.put('/user/password', data)
      return response.data
    } catch (error: any) {
      console.error('修改密码失败:', error)
      throw new Error(error.response?.data?.error || '修改密码失败')
    }
  },

  // 获取用户游戏统计
  async getGameStats(): Promise<ApiResponse> {
    try {
      const response = await api.get('/user/stats')
      return response.data
    } catch (error: any) {
      console.error('获取游戏统计失败:', error)
      throw new Error(error.response?.data?.error || '获取游戏统计失败')
    }
  },

  // 获取用户游戏历史
  async getGameHistory(page: number = 1, limit: number = 10): Promise<ApiResponse> {
    try {
      const response = await api.get(`/user/games?page=${page}&limit=${limit}`)
      return response.data
    } catch (error: any) {
      console.error('获取游戏历史失败:', error)
      throw new Error(error.response?.data?.error || '获取游戏历史失败')
    }
  }
}

// 游戏相关API
export const gameApi = {
  // 保存游戏记录
  async saveGame(gameData: {
    game_type: string
    result: string
    moves: number
    duration: number
    difficulty?: string
    board_state?: string
  }): Promise<ApiResponse> {
    try {
      const response = await api.post('/games', gameData)
      return response.data
    } catch (error: any) {
      console.error('保存游戏记录失败:', error)
      throw new Error(error.response?.data?.error || '保存游戏记录失败')
    }
  },

  // 获取游戏详情
  async getGame(gameId: string): Promise<ApiResponse> {
    try {
      const response = await api.get(`/games/${gameId}`)
      return response.data
    } catch (error: any) {
      console.error('获取游戏详情失败:', error)
      throw new Error(error.response?.data?.error || '获取游戏详情失败')
    }
  }
}

export default api