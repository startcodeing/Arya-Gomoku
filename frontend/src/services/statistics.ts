import api from './authApi'

// 统计数据类型定义
export interface UserGameStats {
  totalGames: number
  winCount: number
  lossCount: number
  drawCount: number
  winRate: number
  aiGames: number
  llmGames: number
  pvpGames: number
  avgDuration: number
  lastGameTime: string
}

export interface GameSummary {
  id: string
  type: 'ai' | 'llm' | 'pvp'
  status: string
  result: 'win' | 'loss' | 'draw' | 'unknown'
  userId: string
  username: string
  startTime: string
  endTime: string
  duration: number
  moves: any[]
}

export interface GameHistoryResponse {
  games: GameSummary[]
  total: number
  page: number
  pageSize: number
  hasMore: boolean
}

export interface SystemStatistics {
  totalUsers: number
  activeUsers: number
  totalGames: number
  aiGames: number
  llmGames: number
  pvpGames: number
  avgGameDuration: number
  completedGames: number
  completedRate: number
}

export interface PlayerRanking {
  userId: string
  username: string
  nickname?: string
  avatar?: string
  totalGames: number
  winCount: number
  winRate: number
  rating: number
  rank: number
}

export interface GameTypeStats {
  ai: number
  llm: number
  pvp: number
}

export interface DifficultyStats {
  easy: number
  medium: number
  hard: number
  expert: number
}

// 重试配置
const RETRY_CONFIG = {
  maxRetries: 3,
  retryDelay: 1000,
  backoffMultiplier: 2
}

// 重试函数
async function withRetry<T>(
  fn: () => Promise<T>,
  retries = RETRY_CONFIG.maxRetries,
  delay = RETRY_CONFIG.retryDelay
): Promise<T> {
  try {
    return await fn()
  } catch (error: any) {
    if (retries > 0 && error.response?.status >= 500) {
      await new Promise(resolve => setTimeout(resolve, delay))
      return withRetry(fn, retries - 1, delay * RETRY_CONFIG.backoffMultiplier)
    }
    throw error
  }
}

export const statisticsApi = {
  // 获取用户游戏统计
  async getUserStats(): Promise<UserGameStats> {
    return withRetry(async () => {
      try {
        const response = await api.get('/statistics/user')
        
        if (!response.data.success) {
          throw new Error(response.data.message || '获取用户统计失败')
        }
        
        return response.data.data
      } catch (error: any) {
        throw new Error(error.response?.data?.message || '获取用户统计失败')
      }
    })
  },

  // 获取用户游戏历史
  async getUserGameHistory(page = 1, pageSize = 10): Promise<GameHistoryResponse> {
    return withRetry(async () => {
      try {
        const response = await api.get('/statistics/games', {
          params: { page, page_size: pageSize }
        })
        
        if (!response.data.success) {
          throw new Error(response.data.message || '获取游戏历史失败')
        }
        
        const data = response.data.data
        return {
          games: data.games || [],
          total: data.total || 0,
          page: data.page || 1,
          pageSize: data.page_size || pageSize,
          hasMore: data.has_more || false
        }
      } catch (error: any) {
        throw new Error(error.response?.data?.message || '获取游戏历史失败')
      }
    })
  },

  // 获取系统统计（需要管理员权限）
  async getSystemStats(): Promise<SystemStatistics> {
    return withRetry(async () => {
      try {
        const response = await api.get('/statistics/system')
        
        if (!response.data.success) {
          throw new Error(response.data.message || '获取系统统计失败')
        }
        
        return response.data.data
      } catch (error: any) {
        throw new Error(error.response?.data?.message || '获取系统统计失败')
      }
    })
  },

  // 获取指定日期范围的游戏统计
  async getGameStatsByDateRange(startDate: string, endDate: string): Promise<any> {
    return withRetry(async () => {
      try {
        const response = await api.get('/statistics/date-range', {
          params: { start_date: startDate, end_date: endDate }
        })
        
        if (!response.data.success) {
          throw new Error(response.data.message || '获取日期范围统计失败')
        }
        
        return response.data.data
      } catch (error: any) {
        throw new Error(error.response?.data?.message || '获取日期范围统计失败')
      }
    })
  },

  // 获取排行榜
  async getTopPlayers(gameType = 'all', limit = 10): Promise<PlayerRanking[]> {
    return withRetry(async () => {
      try {
        const response = await api.get('/statistics/top-players', {
          params: { game_type: gameType, limit }
        })
        
        if (!response.data.success) {
          throw new Error(response.data.message || '获取排行榜失败')
        }
        
        return response.data.data || []
      } catch (error: any) {
        throw new Error(error.response?.data?.message || '获取排行榜失败')
      }
    })
  },

  // 搜索游戏记录
  async searchGames(keyword: string, page = 1, pageSize = 10): Promise<GameHistoryResponse> {
    return withRetry(async () => {
      try {
        const response = await api.get('/statistics/search', {
          params: { keyword, page, page_size: pageSize }
        })
        
        if (!response.data.success) {
          throw new Error(response.data.message || '搜索游戏记录失败')
        }
        
        const data = response.data.data
        return {
          games: data.games || [],
          total: data.total || 0,
          page: data.page || 1,
          pageSize: data.page_size || pageSize,
          hasMore: data.has_more || false
        }
      } catch (error: any) {
        throw new Error(error.response?.data?.message || '搜索游戏记录失败')
      }
    })
  },

  // 获取游戏类型统计
  async getGameTypeStats(): Promise<GameTypeStats> {
    return withRetry(async () => {
      try {
        const response = await api.get('/statistics/game-types')
        
        if (!response.data.success) {
          throw new Error(response.data.message || '获取游戏类型统计失败')
        }
        
        return response.data.data
      } catch (error: any) {
        throw new Error(error.response?.data?.message || '获取游戏类型统计失败')
      }
    })
  },

  // 获取难度统计
  async getDifficultyStats(): Promise<DifficultyStats> {
    return withRetry(async () => {
      try {
        const response = await api.get('/statistics/difficulties')
        
        if (!response.data.success) {
          throw new Error(response.data.message || '获取难度统计失败')
        }
        
        return response.data.data
      } catch (error: any) {
        throw new Error(error.response?.data?.message || '获取难度统计失败')
      }
    })
  },

  // 导出游戏数据
  async exportGameData(format = 'json', startDate?: string, endDate?: string): Promise<Blob> {
    return withRetry(async () => {
      try {
        const params: any = { format }
        if (startDate) params.start_date = startDate
        if (endDate) params.end_date = endDate
        
        const response = await api.get('/statistics/export', {
          params,
          responseType: 'blob'
        })
        
        return response.data
      } catch (error: any) {
        throw new Error(error.response?.data?.message || '导出游戏数据失败')
      }
    })
  },

  // 删除游戏记录
  async deleteGameRecord(gameId: string): Promise<void> {
    return withRetry(async () => {
      try {
        const response = await api.delete(`/statistics/games/${gameId}`)
        
        if (!response.data.success) {
          throw new Error(response.data.message || '删除游戏记录失败')
        }
      } catch (error: any) {
        throw new Error(error.response?.data?.message || '删除游戏记录失败')
      }
    })
  }
}

export default statisticsApi