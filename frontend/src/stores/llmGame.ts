// LLM对战状态管理 Store
// 使用 Pinia 管理LLM游戏状态

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { llmApi, llmGameUtils } from '../services/llmApi'
import type {
  LLMGame,
  LLMModel,
  LLMMove,
  LLMGameStatus,
  LLMConfigRequest,
  LLMConfigResponse,
  Move
} from '../types/game'

export const useLLMGameStore = defineStore('llmGame', () => {
  // ===== 状态定义 =====
  
  // 当前游戏信息
  const currentGame = ref<LLMGame | null>(null)
  
  // 可用模型列表
  const availableModels = ref<LLMModel[]>([])
  
  // 当前选择的模型
  const selectedModel = ref<string>('deepseek')
  
  // 游戏状态
  const isGameActive = ref(false)
  const isLoading = ref(false)
  const isThinking = ref(false)
  
  // 错误信息
  const error = ref<string | null>(null)
  
  // 最后一次AI移动信息
  const lastAIMove = ref<LLMMove | null>(null)
  
  // 模型配置状态
  const modelConfigs = ref<Record<string, LLMConfigResponse>>({})
  
  // ===== 计算属性 =====
  
  // 当前棋盘状态
  const board = computed(() => {
    return currentGame.value?.board || Array(15).fill(null).map(() => Array(15).fill(0))
  })
  
  // 游戏状态
  const gameStatus = computed(() => {
    return currentGame.value?.status || 'playing'
  })
  
  // 是否游戏结束
  const isGameFinished = computed(() => {
    return llmGameUtils.isGameFinished(gameStatus.value)
  })
  
  // 游戏结果消息
  const gameResultMessage = computed(() => {
    return llmGameUtils.getGameResultMessage(gameStatus.value)
  })
  
  // 可用的已配置模型
  const configuredModels = computed(() => {
    return availableModels.value.filter(model => model.status === 'available')
  })
  
  // 当前选择的模型信息
  const selectedModelInfo = computed(() => {
    return availableModels.value.find(model => model.name === selectedModel.value)
  })
  
  // 移动历史
  const moveHistory = computed(() => {
    return currentGame.value?.moves || []
  })
  
  // 游戏统计
  const gameStats = computed(() => {
    const moves = moveHistory.value
    return {
      totalMoves: moves.length,
      humanMoves: moves.filter(move => move.player === 1).length,
      aiMoves: moves.filter(move => move.player === 2).length,
      gameDuration: currentGame.value ? 
        llmGameUtils.calculateGameDuration(currentGame.value.start_time, currentGame.value.end_time) : 
        '0秒'
    }
  })
  
  // ===== Actions =====
  
  // 加载可用模型
  async function loadAvailableModels() {
    try {
      isLoading.value = true
      error.value = null
      
      const response = await llmApi.getModels()
      if (response.success && response.data) {
        availableModels.value = response.data
        
        // 如果当前选择的模型不可用，选择第一个可用的模型
        if (!configuredModels.value.find(model => model.name === selectedModel.value)) {
          const firstAvailable = configuredModels.value[0]
          if (firstAvailable) {
            selectedModel.value = firstAvailable.name
          }
        }
      } else {
        throw new Error(response.error || '加载模型列表失败')
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : '加载模型列表失败'
      console.error('Failed to load models:', err)
    } finally {
      isLoading.value = false
    }
  }
  
  // 开始新游戏
  async function startNewGame(modelName?: string) {
    try {
      isLoading.value = true
      error.value = null
      
      const model = modelName || selectedModel.value
      
      // 检查模型是否可用
      const modelInfo = availableModels.value.find(m => m.name === model)
      if (!modelInfo || modelInfo.status !== 'available') {
        throw new Error(`模型 ${model} 不可用，请先配置API密钥`)
      }
      
      const response = await llmApi.startGame(model)
      if (response.success && response.data) {
        currentGame.value = response.data
        selectedModel.value = model
        isGameActive.value = true
        lastAIMove.value = null
      } else {
        throw new Error(response.error || '开始游戏失败')
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : '开始游戏失败'
      console.error('Failed to start game:', err)
    } finally {
      isLoading.value = false
    }
  }
  
  // 进行移动
  async function makeMove(x: number, y: number) {
    if (!currentGame.value || isGameFinished.value || isThinking.value) {
      return
    }
    
    try {
      isThinking.value = true
      error.value = null
      
      const response = await llmApi.makeMove(currentGame.value.game_id, x, y)
      if (response.success && response.data) {
        const { move, game_status, message } = response.data
        
        // 更新棋盘状态
        if (currentGame.value) {
          // 添加人类移动
          const humanMove: Move = { x, y, player: 1 }
          currentGame.value.moves.push(humanMove)
          currentGame.value.board[y][x] = 1
          
          // 添加AI移动（如果有）
          if (move) {
            const aiMove: Move = { x: move.x, y: move.y, player: 2 }
            currentGame.value.moves.push(aiMove)
            currentGame.value.board[move.y][move.x] = 2
            lastAIMove.value = move
          }
          
          // 更新游戏状态
          currentGame.value.status = game_status
          
          // 如果游戏结束，设置结束时间
          if (llmGameUtils.isGameFinished(game_status)) {
            currentGame.value.end_time = new Date().toISOString()
            isGameActive.value = false
          }
        }
      } else {
        throw new Error(response.error || '移动失败')
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : '移动失败'
      console.error('Failed to make move:', err)
    } finally {
      isThinking.value = false
    }
  }
  
  // 结束当前游戏
  async function endGame() {
    if (currentGame.value) {
      try {
        await llmApi.deleteGame(currentGame.value.game_id)
      } catch (err) {
        console.warn('Failed to delete game:', err)
      }
    }
    
    currentGame.value = null
    isGameActive.value = false
    isThinking.value = false
    lastAIMove.value = null
    error.value = null
  }
  
  // 重新开始游戏
  async function restartGame() {
    const model = selectedModel.value
    await endGame()
    await startNewGame(model)
  }
  
  // 更新模型配置
  async function updateModelConfig(modelName: string, config: LLMConfigRequest) {
    try {
      isLoading.value = true
      error.value = null
      
      const response = await llmApi.updateConfig(modelName, config)
      if (response.success) {
        // 重新加载模型列表以更新状态
        await loadAvailableModels()
        
        // 加载更新后的配置
        await loadModelConfig(modelName)
      } else {
        throw new Error(response.error || '更新配置失败')
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : '更新配置失败'
      console.error('Failed to update config:', err)
      throw err
    } finally {
      isLoading.value = false
    }
  }
  
  // 加载模型配置
  async function loadModelConfig(modelName: string) {
    try {
      const response = await llmApi.getConfig(modelName)
      if (response.success && response.data) {
        modelConfigs.value[modelName] = response.data
      }
    } catch (err) {
      console.warn(`Failed to load config for ${modelName}:`, err)
    }
  }
  
  // 选择模型
  function selectModel(modelName: string) {
    const modelInfo = availableModels.value.find(m => m.name === modelName)
    if (modelInfo && modelInfo.status === 'available') {
      selectedModel.value = modelName
    }
  }
  
  // 清除错误
  function clearError() {
    error.value = null
  }
  
  // 检查位置是否可以下棋
  function canMakeMove(x: number, y: number): boolean {
    if (!currentGame.value || isGameFinished.value || isThinking.value) {
      return false
    }
    
    // 检查坐标是否有效
    if (x < 0 || x >= 15 || y < 0 || y >= 15) {
      return false
    }
    
    // 检查位置是否为空
    return currentGame.value.board[y][x] === 0
  }
  
  // 获取最后一步移动
  function getLastMove(): Move | null {
    const moves = moveHistory.value
    return moves.length > 0 ? moves[moves.length - 1] : null
  }
  
  // 获取指定位置的棋子
  function getPieceAt(x: number, y: number): number {
    if (x < 0 || x >= 15 || y < 0 || y >= 15) {
      return 0
    }
    return board.value[y][x]
  }
  
  // 初始化store
  async function initialize() {
    await loadAvailableModels()
  }
  
  return {
    // 状态
    currentGame,
    availableModels,
    selectedModel,
    isGameActive,
    isLoading,
    isThinking,
    error,
    lastAIMove,
    modelConfigs,
    
    // 计算属性
    board,
    gameStatus,
    isGameFinished,
    gameResultMessage,
    configuredModels,
    selectedModelInfo,
    moveHistory,
    gameStats,
    
    // 方法
    loadAvailableModels,
    startNewGame,
    makeMove,
    endGame,
    restartGame,
    updateModelConfig,
    loadModelConfig,
    selectModel,
    clearError,
    canMakeMove,
    getLastMove,
    getPieceAt,
    initialize
  }
})