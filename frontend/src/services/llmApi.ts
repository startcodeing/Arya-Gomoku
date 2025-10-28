// LLM API 客户端服务
// 实现与后端LLM接口的通信

import api from './api'
import type {
  ApiResponse,
  LLMGame,
  LLMModel,
  LLMConfig,
  LLMConfigRequest,
  LLMConfigResponse,
  LLMMoveRequest,
  LLMMoveResponse,
  LLMStartGameRequest,
  LLMGameHistoryResponse,
  LLMStats,
  LLMHealthResponse
} from '../types/game'

export const llmApi = {
  // 开始新的LLM游戏
  async startGame(modelName: string): Promise<ApiResponse<LLMGame>> {
    const request: LLMStartGameRequest = {
      model_name: modelName
    }
    
    const response = await api.post('/llm/start', request)
    return response.data
  },

  // 进行移动
  async makeMove(gameId: string, x: number, y: number): Promise<ApiResponse<LLMMoveResponse>> {
    const request: LLMMoveRequest = {
      game_id: gameId,
      move: { x, y }
    }
    
    const response = await api.post('/llm/move', request)
    return response.data
  },

  // 获取游戏信息
  async getGame(gameId: string): Promise<ApiResponse<LLMGame>> {
    const response = await api.get(`/llm/game/${gameId}`)
    return response.data
  },

  // 删除游戏
  async deleteGame(gameId: string): Promise<ApiResponse<void>> {
    const response = await api.delete(`/llm/game/${gameId}`)
    return response.data
  },

  // 获取游戏历史
  async getGameHistory(gameId: string, limit?: number): Promise<ApiResponse<LLMGameHistoryResponse>> {
    const params = limit ? { limit: limit.toString() } : {}
    const response = await api.get(`/llm/game/${gameId}/history`, { params })
    return response.data
  },

  // 获取可用模型列表
  async getModels(): Promise<ApiResponse<LLMModel[]>> {
    const response = await api.get('/llm/models')
    return response.data
  },

  // 更新模型配置
  async updateConfig(modelName: string, config: LLMConfigRequest): Promise<ApiResponse<void>> {
    const response = await api.put(`/llm/config/${modelName}`, config)
    return response.data
  },

  // 获取模型配置
  async getConfig(modelName: string): Promise<ApiResponse<LLMConfigResponse>> {
    const response = await api.get(`/llm/config/${modelName}`)
    return response.data
  },

  // 获取统计信息
  async getStats(): Promise<ApiResponse<LLMStats>> {
    const response = await api.get('/llm/stats')
    return response.data
  },

  // 健康检查
  async healthCheck(): Promise<ApiResponse<LLMHealthResponse>> {
    const response = await api.get('/llm/health')
    return response.data
  }
}

// LLM 游戏工具函数
export const llmGameUtils = {
  // 检查游戏是否结束
  isGameFinished(status: string): boolean {
    return ['human_win', 'ai_win', 'draw'].includes(status)
  },

  // 获取游戏结果描述
  getGameResultMessage(status: string): string {
    switch (status) {
      case 'human_win':
        return '恭喜！你获胜了！'
      case 'ai_win':
        return 'AI获胜！'
      case 'draw':
        return '平局！'
      case 'playing':
        return '游戏进行中...'
      default:
        return '未知状态'
    }
  },

  // 获取模型状态描述
  getModelStatusMessage(status: string): string {
    switch (status) {
      case 'available':
        return '可用'
      case 'not_configured':
        return '未配置'
      case 'unavailable':
        return '不可用'
      default:
        return '未知'
    }
  },

  // 获取模型状态颜色
  getModelStatusColor(status: string): string {
    switch (status) {
      case 'available':
        return 'text-green-600'
      case 'not_configured':
        return 'text-yellow-600'
      case 'unavailable':
        return 'text-red-600'
      default:
        return 'text-gray-600'
    }
  },

  // 验证API密钥格式
  validateApiKey(provider: string, apiKey: string): { valid: boolean; message?: string } {
    if (!apiKey.trim()) {
      return { valid: false, message: 'API密钥不能为空' }
    }

    switch (provider) {
      case 'deepseek':
        if (!apiKey.startsWith('sk-')) {
          return { valid: false, message: 'DeepSeek API密钥应以 "sk-" 开头' }
        }
        break
      case 'openai':
        if (!apiKey.startsWith('sk-')) {
          return { valid: false, message: 'OpenAI API密钥应以 "sk-" 开头' }
        }
        break
      case 'ollama':
        // Ollama 不需要API密钥
        return { valid: true }
      default:
        break
    }

    return { valid: true }
  },

  // 格式化时间戳
  formatTimestamp(timestamp: string): string {
    try {
      const date = new Date(timestamp)
      return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit'
      })
    } catch {
      return '无效时间'
    }
  },

  // 计算游戏时长
  calculateGameDuration(startTime: string, endTime?: string): string {
    try {
      const start = new Date(startTime)
      const end = endTime ? new Date(endTime) : new Date()
      const duration = Math.floor((end.getTime() - start.getTime()) / 1000)
      
      const hours = Math.floor(duration / 3600)
      const minutes = Math.floor((duration % 3600) / 60)
      const seconds = duration % 60
      
      if (hours > 0) {
        return `${hours}小时${minutes}分钟${seconds}秒`
      } else if (minutes > 0) {
        return `${minutes}分钟${seconds}秒`
      } else {
        return `${seconds}秒`
      }
    } catch {
      return '未知'
    }
  },

  // 获取置信度颜色
  getConfidenceColor(confidence: number): string {
    if (confidence >= 0.8) {
      return 'text-green-600'
    } else if (confidence >= 0.6) {
      return 'text-yellow-600'
    } else {
      return 'text-red-600'
    }
  },

  // 格式化置信度
  formatConfidence(confidence: number): string {
    return `${Math.round(confidence * 100)}%`
  }
}