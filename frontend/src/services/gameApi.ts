import api from './authApi'

export interface UserGameStats {
  userId: string
  totalGames: number
  winGames: number
  loseGames: number
  drawGames: number
  winRate: number
  aiGames: number
  llmGames: number
  pvpGames: number
  avgDuration: number
  lastGameTime: string
}

export interface GameStats {
  totalGames: number
  aiGames: number
  llmGames: number
  pvpGames: number
  completedRate: number
  avgDuration: number
}

export interface PlayerRanking {
  userId: string
  username: string
  nickname: string
  winGames: number
  winRate: number
  rank: number
}

export interface GameSummary {
  id: string
  type: string
  status: string
  winner: string
  duration: number
  createdAt: string
  userId: string
  username: string
}

export interface GameHistory {
  games: GameSummary[]
  total: number
  page: number
  pageSize: number
}

export interface AIGame {
  id: string
  userId: string
  difficulty: string
  aiType: string
  status: string
  boardSize: number
  moves: any[]
  moveCount: number
  totalTimeMs: number
  startedAt: string
  endedAt?: string
  createdAt: string
  updatedAt: string
}

export interface LLMGame {
  id: string
  userId: string
  modelName: string
  difficulty: string
  status: string
  boardSize: number
  moves: any[]
  moveCount: number
  totalTimeMs: number
  startedAt: string
  endedAt?: string
  createdAt: string
  updatedAt: string
}

export interface PVPGame {
  id: string
  roomId: string
  status: string
  boardSize: number
  moves: any[]
  moveCount: number
  currentPlayerId?: string
  winnerId?: string
  startedAt: string
  endedAt?: string
  createdAt: string
  updatedAt: string
  players: any[]
}

export const gameApi = {
  // 获取用户游戏统计
  async getUserStats(): Promise<UserGameStats> {
    try {
      const response = await api.get('/statistics/user')
      
      if (!response.data.success) {
        throw new Error(response.data.message || '获取用户统计失败')
      }
      
      const data = response.data.data
      return {
        userId: data.user_id || '',
        totalGames: (data.ai_games || 0) + (data.llm_games || 0) + (data.pvp_games || 0),
        winGames: data.win_count || 0,
        loseGames: data.loss_count || 0,
        drawGames: data.draw_count || 0,
        winRate: data.win_rate || 0,
        aiGames: data.ai_games || 0,
        llmGames: data.llm_games || 0,
        pvpGames: data.pvp_games || 0,
        avgDuration: data.avg_duration || 0,
        lastGameTime: data.last_game_time || ''
      }
    } catch (error: any) {
      throw new Error(error.response?.data?.message || '获取用户统计失败')
    }
  },

  // 获取游戏统计（按日期范围）
  async getGameStats(startDate: string, endDate: string): Promise<GameStats> {
    try {
      const response = await api.get('/statistics/date-range', {
        params: { start_date: startDate, end_date: endDate }
      })
      
      if (!response.data.success) {
        throw new Error(response.data.message || '获取游戏统计失败')
      }
      
      return response.data.data
    } catch (error: any) {
      throw new Error(error.response?.data?.message || '获取游戏统计失败')
    }
  },

  // 获取排行榜
  async getTopPlayers(gameType: string, limit: number = 10): Promise<PlayerRanking[]> {
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
  },

  // 获取用户游戏历史
  async getUserGameHistory(
    page: number = 1,
    limit: number = 20,
    gameType?: string
  ): Promise<GameHistory> {
    try {
      const params: any = { page, limit }
      if (gameType) {
        params.type = gameType
      }
      
      const response = await api.get('/statistics/games', { params })
      
      if (!response.data.success) {
        throw new Error(response.data.message || '获取游戏历史失败')
      }
      
      const data = response.data.data
      return {
        games: data.games || [],
        total: data.pagination?.total || 0,
        page: data.pagination?.page || page,
        pageSize: data.pagination?.limit || limit
      }
    } catch (error: any) {
      throw new Error(error.response?.data?.message || '获取游戏历史失败')
    }
  },

  // 获取最近游戏
  async getRecentGames(limit: number = 10): Promise<GameSummary[]> {
    try {
      const response = await api.get('/statistics/games', {
        params: { limit }
      })
      
      if (!response.data.success) {
        throw new Error(response.data.message || '获取最近游戏失败')
      }
      
      const data = response.data.data
      return data.games || []
    } catch (error: any) {
      throw new Error(error.response?.data?.message || '获取最近游戏失败')
    }
  },

  // 搜索游戏
  async searchGames(
    keyword: string,
    page: number = 1,
    pageSize: number = 20
  ): Promise<GameHistory> {
    try {
      const response = await api.get('/statistics/search', {
        params: { keyword, page, pageSize }
      })
      return response.data
    } catch (error: any) {
      throw new Error(error.response?.data?.message || '搜索游戏失败')
    }
  },

  // AI 游戏相关
  async getAIGame(gameId: string): Promise<AIGame> {
    try {
      const response = await api.get(`/games/ai/${gameId}`)
      return response.data
    } catch (error: any) {
      throw new Error(error.response?.data?.message || '获取AI游戏失败')
    }
  },

  async getUserAIGames(
    userId: string,
    page: number = 1,
    pageSize: number = 20
  ): Promise<{ games: AIGame[], total: number }> {
    try {
      const response = await api.get(`/games/ai/users/${userId}`, {
        params: { page, pageSize }
      })
      return response.data
    } catch (error: any) {
      throw new Error(error.response?.data?.message || '获取AI游戏列表失败')
    }
  },

  // LLM 游戏相关
  async getLLMGame(gameId: string): Promise<LLMGame> {
    try {
      const response = await api.get(`/games/llm/${gameId}`)
      return response.data
    } catch (error: any) {
      throw new Error(error.response?.data?.message || '获取LLM游戏失败')
    }
  },

  async getUserLLMGames(
    userId: string,
    page: number = 1,
    pageSize: number = 20
  ): Promise<{ games: LLMGame[], total: number }> {
    try {
      const response = await api.get(`/games/llm/users/${userId}`, {
        params: { page, pageSize }
      })
      return response.data
    } catch (error: any) {
      throw new Error(error.response?.data?.message || '获取LLM游戏列表失败')
    }
  },

  // PVP 游戏相关
  async getPVPGame(gameId: string): Promise<PVPGame> {
    try {
      const response = await api.get(`/games/pvp/${gameId}`)
      return response.data
    } catch (error: any) {
      throw new Error(error.response?.data?.message || '获取PVP游戏失败')
    }
  },

  async getUserPVPGames(
    userId: string,
    page: number = 1,
    pageSize: number = 20
  ): Promise<{ games: PVPGame[], total: number }> {
    try {
      const response = await api.get(`/games/pvp/users/${userId}`, {
        params: { page, pageSize }
      })
      return response.data
    } catch (error: any) {
      throw new Error(error.response?.data?.message || '获取PVP游戏列表失败')
    }
  },

  // 数据导出
  async exportUserData(userId: string, format: 'json' | 'csv' = 'json'): Promise<Blob> {
    try {
      const response = await api.get(`/games/users/${userId}/export`, {
        params: { format },
        responseType: 'blob'
      })
      return response.data
    } catch (error: any) {
      throw new Error(error.response?.data?.message || '导出数据失败')
    }
  },

  // 删除游戏记录
  async deleteGame(gameId: string, gameType: 'ai' | 'llm' | 'pvp'): Promise<void> {
    try {
      await api.delete(`/games/${gameType}/${gameId}`)
    } catch (error: any) {
      throw new Error(error.response?.data?.message || '删除游戏记录失败')
    }
  },

  // 批量删除游戏记录
  async batchDeleteGames(gameIds: string[], gameType: 'ai' | 'llm' | 'pvp'): Promise<void> {
    try {
      await api.post(`/games/${gameType}/batch-delete`, { gameIds })
    } catch (error: any) {
      throw new Error(error.response?.data?.message || '批量删除游戏记录失败')
    }
  },

  // 保存游戏记录
  async saveGameRecord(gameRecord: {
    gameType: 'ai' | 'llm' | 'pvp'
    difficulty?: string
    llmModel?: string
    opponentName?: string
    result: 'win' | 'lose' | 'draw'
    moves: any[]
    moveCount: number
    duration?: number
    aiStats?: any
    modelConfig?: any
  }): Promise<void> {
    try {
      // 游戏记录通过各自的游戏类型端点自动保存
      // AI游戏：通过 /api/ai/games/:id/move 自动更新
      // LLM游戏：通过 /api/llm/games/:id/move 自动更新  
      // PVP游戏：通过房间系统自动保存
      
      // 这里只是一个占位符方法，实际的游戏记录保存
      // 是在游戏进行过程中自动完成的
      console.log('游戏记录已自动保存:', gameRecord)
      
      // 如果需要手动触发保存，可以调用对应的游戏更新API
      // 但通常情况下，游戏记录会在游戏过程中自动保存
    } catch (error: any) {
      throw new Error(error.response?.data?.message || '保存游戏记录失败')
    }
  }
}