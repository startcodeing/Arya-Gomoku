import { defineStore } from 'pinia'
import { ref, computed, reactive } from 'vue'
import { llmApi } from '../services/llmApi'
import { Player, GameStatus, type Position, type BoardState } from '../types/game'
import { createInitialGameState, makeMove, getGameStatus, BOARD_SIZE } from '../utils/gameLogic'

export interface LLMModel {
  id: string
  name: string
  provider: string
  description?: string
}

export interface LLMConfig {
  temperature: number
  maxTokens: number
  systemPrompt: string
}

export interface LLMGameSession {
  id: string
  status: 'playing' | 'finished'
  board: number[][]
  moves: Position[]
  currentPlayer: Player
  winner?: Player
  startedAt: string
  finishedAt?: string
}

export const useLLMGameStore = defineStore('llmGame', () => {
  // 基础状态 - 简化版本，参考Home.vue模式
  const gameState = reactive<BoardState>(createInitialGameState())
  const currentSession = ref<LLMGameSession | null>(null)
  const availableModels = ref<LLMModel[]>([])
  const selectedModel = ref<LLMModel | null>(null)
  const modelConfig = ref<LLMConfig>({
    temperature: 0.7,
    maxTokens: 1000,
    systemPrompt: '你是一个五子棋专家，请分析当前棋局并选择最佳落子位置。'
  })
  
  // 状态标志
  const isLoading = ref(false)
  const isAiThinking = ref(false)
  const error = ref<string | null>(null)
  const moveHistory = ref<Position[]>([])
  
  // 统计数据
  const statistics = reactive({
    humanWins: 0,
    aiWins: 0,
    draws: 0,
    totalGames: 0
  })

  // 计算属性
  const board = computed(() => gameState.board)
  const gameStatus = computed(() => gameState.gameStatus)
  const currentPlayer = computed(() => gameState.currentPlayer)
  const lastMove = computed(() => gameState.lastMove)
  const moveCount = computed(() => moveHistory.value.length)
  
  const canMakeMove = computed(() => 
    gameState.gameStatus === GameStatus.PLAYING && 
    gameState.currentPlayer === Player.HUMAN &&
    !isLoading.value && 
    !isAiThinking.value
  )

  const isGameActive = computed(() => gameState.gameStatus === GameStatus.PLAYING)

  // Actions
  async function loadAvailableModels() {
    if (isLoading.value) return
    
    isLoading.value = true
    error.value = null
    
    try {
      const response = await llmApi.getModels()
      availableModels.value = response.models || []
      
      // 自动选择第一个模型
      if (availableModels.value.length > 0 && !selectedModel.value) {
        selectedModel.value = availableModels.value[0]
      }
    } catch (err: any) {
      error.value = err.message || '加载模型列表失败'
      console.error('加载模型失败:', err)
    } finally {
      isLoading.value = false
    }
  }

  async function startNewGame() {
    if (!selectedModel.value) {
      error.value = '请先选择一个AI模型'
      return
    }

    isLoading.value = true
    error.value = null

    try {
      const response = await llmApi.startGame({
        modelId: selectedModel.value.id,
        config: modelConfig.value
      })

      // 重置游戏状态
      Object.assign(gameState, createInitialGameState())
      moveHistory.value = []
      currentSession.value = {
        id: response.gameId,
        status: 'playing',
        board: gameState.board,
        moves: [],
        currentPlayer: Player.HUMAN,
        startedAt: new Date().toISOString()
      }

      clearError()
    } catch (err: any) {
      error.value = err.message || '开始游戏失败'
      console.error('开始游戏失败:', err)
    } finally {
      isLoading.value = false
    }
  }

  async function makePlayerMove(x: number, y: number) {
    if (!canMakeMove.value || !currentSession.value) return false

    try {
      isLoading.value = true

      // 执行玩家移动
      if (!makeMove(gameState.board, x, y, Player.HUMAN)) {
        error.value = '无效的移动位置'
        return false
      }

      // 记录移动
      const move: Position = { x, y }
      moveHistory.value.push(move)
      gameState.lastMove = move

      // 检查游戏状态
      gameState.gameStatus = getGameStatus(gameState.board, move)
      
      if (gameState.gameStatus !== GameStatus.PLAYING) {
        await handleGameEnd()
        return true
      }

      // 切换到AI回合并获取AI移动
      gameState.currentPlayer = Player.AI
      await getAIMove()
      
      return true
    } catch (err: any) {
      error.value = err.message || '移动失败'
      console.error('玩家移动失败:', err)
      return false
    } finally {
      isLoading.value = false
    }
  }

  async function getAIMove() {
    if (!currentSession.value) return

    try {
      isAiThinking.value = true

      const response = await llmApi.makeMove(currentSession.value.id, {
        x: gameState.lastMove?.x || 0,
        y: gameState.lastMove?.y || 0
      })

      // 后端返回的数据结构是 { data: { move: {...}, gameStatus: "...", reasoning: "..." } }
      const llmResponse = response.data || response
      const aiMoveData = llmResponse.move || llmResponse.Move

      // 验证AI移动数据
      if (!aiMoveData || typeof aiMoveData.x !== 'number' || typeof aiMoveData.y !== 'number') {
        throw new Error('AI返回的移动数据格式无效')
      }

      // 验证AI移动位置
      if (!isValidMove(aiMoveData.x, aiMoveData.y)) {
        throw new Error('AI返回了无效的移动位置')
      }

      // 执行AI移动
      if (!makeMove(gameState.board, aiMoveData.x, aiMoveData.y, Player.AI)) {
        throw new Error('AI移动执行失败')
      }

      // 记录AI移动
      const aiMove: Position = { x: aiMoveData.x, y: aiMoveData.y }
      moveHistory.value.push(aiMove)
      gameState.lastMove = aiMove

      // 更新游戏状态（使用后端返回的状态）
      const backendGameStatus = llmResponse.gameStatus || llmResponse.GameStatus
      if (backendGameStatus) {
        switch (backendGameStatus) {
          case 'playing':
            gameState.gameStatus = GameStatus.PLAYING
            break
          case 'ai_win':
            gameState.gameStatus = GameStatus.AI_WIN
            break
          case 'human_win':
            gameState.gameStatus = GameStatus.HUMAN_WIN
            break
          case 'draw':
            gameState.gameStatus = GameStatus.DRAW
            break
          default:
            // 如果后端状态未知，使用前端逻辑检查
            gameState.gameStatus = getGameStatus(gameState.board, aiMove)
        }
      } else {
        // 如果后端没有返回状态，使用前端逻辑检查
        gameState.gameStatus = getGameStatus(gameState.board, aiMove)
      }
      
      if (gameState.gameStatus !== GameStatus.PLAYING) {
        await handleGameEnd()
        return
      }

      // 切换回玩家回合
      gameState.currentPlayer = Player.HUMAN

    } catch (err: any) {
      error.value = err.message || 'AI移动失败'
      console.error('AI移动失败:', err)
      // 如果AI移动失败，切换回玩家回合
      gameState.currentPlayer = Player.HUMAN
    } finally {
      isAiThinking.value = false
    }
  }

  function isValidMove(x: number, y: number): boolean {
    return x >= 0 && x < BOARD_SIZE && 
           y >= 0 && y < BOARD_SIZE && 
           gameState.board[y][x] === Player.NONE
  }

  async function handleGameEnd() {
    if (!currentSession.value) return

    // 更新统计数据
    statistics.totalGames++
    
    switch (gameState.gameStatus) {
      case GameStatus.HUMAN_WIN:
        statistics.humanWins++
        gameState.winner = Player.HUMAN
        break
      case GameStatus.AI_WIN:
        statistics.aiWins++
        gameState.winner = Player.AI
        break
      case GameStatus.DRAW:
        statistics.draws++
        break
    }

    // 更新会话状态
    currentSession.value.status = 'finished'
    currentSession.value.finishedAt = new Date().toISOString()
    currentSession.value.winner = gameState.winner
  }

  function restartGame() {
    // 重置游戏状态
    Object.assign(gameState, createInitialGameState())
    moveHistory.value = []
    currentSession.value = null
    clearError()
  }

  function selectModel(model: LLMModel) {
    selectedModel.value = model
  }

  function updateModelConfig(config: Partial<LLMConfig>) {
    Object.assign(modelConfig.value, config)
  }

  function clearError() {
    error.value = null
  }

  // 辅助方法
  function getPieceAt(x: number, y: number): Player {
    if (x < 0 || x >= BOARD_SIZE || y < 0 || y >= BOARD_SIZE) {
      return Player.NONE
    }
    return gameState.board[y][x]
  }

  function getLastMove(): Position | null {
    return gameState.lastMove || null
  }

  // 初始化
  function initialize() {
    loadAvailableModels()
  }

  return {
    // 状态
    gameState,
    currentSession,
    availableModels,
    selectedModel,
    modelConfig,
    isLoading,
    isAiThinking,
    error,
    moveHistory,
    statistics,
    
    // 计算属性
    board,
    gameStatus,
    currentPlayer,
    lastMove,
    moveCount,
    canMakeMove,
    isGameActive,
    
    // Actions
    loadAvailableModels,
    startNewGame,
    makePlayerMove,
    restartGame,
    selectModel,
    updateModelConfig,
    clearError,
    getPieceAt,
    getLastMove,
    initialize
  }
})