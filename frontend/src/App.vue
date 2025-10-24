<template>
  <div id="app">
    <div class="game-container">
      <div class="game-header">
        <h1>五子棋人机对战</h1>
        <p class="subtitle">Gomoku AI vs Human</p>
      </div>

      <div class="game-content">
        <div class="board-section">
          <Board
            :board="gameState.board"
            :current-player="gameState.currentPlayer"
            :last-move="gameState.lastMove"
            :can-move="canMove"
            @move="handlePlayerMove"
          />
        </div>

        <div class="control-section">
          <ControlPanel
            :game-status="gameState.gameStatus"
            :current-player="gameState.currentPlayer"
            :last-move="gameState.lastMove"
            :move-count="moveCount"
            :is-processing="isProcessing"
            :is-ai-thinking="isAiThinking"
            :statistics="statistics"
            @restart="handleRestart"
          />
        </div>
      </div>

      <div class="game-footer">
        <p>&copy; 2024 五子棋人机对战系统</p>
      </div>
    </div>

    <!-- 错误提示 -->
    <div v-if="errorMessage" class="error-toast" @click="clearError">
      <div class="error-content">
        <span class="error-icon">⚠️</span>
        <span class="error-text">{{ errorMessage }}</span>
        <button class="error-close">&times;</button>
      </div>
    </div>

    <!-- 成功提示 -->
    <div v-if="successMessage" class="success-toast" @click="clearSuccess">
      <div class="success-content">
        <span class="success-icon">✅</span>
        <span class="success-text">{{ successMessage }}</span>
        <button class="success-close">&times;</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import Board from './components/Board.vue'
import ControlPanel from './components/ControlPanel.vue'
import { Player, GameStatus, type Position, type BoardState } from './types/game'
import { aiApi } from './services/api'
import {
  createInitialGameState,
  makeMove,
  getGameStatus,
  copyBoard,
  BOARD_SIZE
} from './utils/gameLogic'

// 响应式数据
const gameState = reactive<BoardState>(createInitialGameState())
const isProcessing = ref(false)
const isAiThinking = ref(false)
const errorMessage = ref('')
const successMessage = ref('')
const moveHistory = ref<Position[]>([])

// 统计数据
const statistics = reactive({
  humanWins: 0,
  aiWins: 0,
  draws: 0,
  totalGames: 0
})

// 计算属性
const moveCount = computed(() => moveHistory.value.length)
const canMove = computed(() => 
  gameState.gameStatus === GameStatus.PLAYING && 
  !isProcessing.value && 
  !isAiThinking.value
)

// 方法
async function handlePlayerMove(x: number, y: number) {
  if (!canMove.value || gameState.currentPlayer !== Player.HUMAN) return

  try {
    isProcessing.value = true

    // 执行玩家移动
    if (!makeMove(gameState.board, x, y, Player.HUMAN)) {
      showError('无效的移动位置')
      return
    }

    // 记录移动
    const move: Position = { x, y }
    moveHistory.value.push(move)
    gameState.lastMove = move

    // 检查游戏状态
    gameState.gameStatus = getGameStatus(gameState.board, move)
    
    if (gameState.gameStatus !== GameStatus.PLAYING) {
      handleGameEnd()
      return
    }

    // 切换到AI回合
    gameState.currentPlayer = Player.AI
    
    // 获取AI移动
    await getAIMove()

  } catch (error: any) {
    showError(error.message || '移动失败')
  } finally {
    isProcessing.value = false
  }
}

async function getAIMove() {
  try {
    isAiThinking.value = true

    // 准备AI请求数据
    const request = {
      board: gameState.board,
      player: Player.AI,
      lastMove: gameState.lastMove || { x: 0, y: 0 }
    }

    // 调用AI接口
    const response = await aiApi.getMove(request)

    // 验证AI移动
    if (!isValidAIMove(response.aiMove.x, response.aiMove.y)) {
      throw new Error('AI返回了无效的移动')
    }

    // 执行AI移动
    if (!makeMove(gameState.board, response.aiMove.x, response.aiMove.y, Player.AI)) {
      throw new Error('AI移动执行失败')
    }

    // 记录AI移动
    const aiMove: Position = { x: response.aiMove.x, y: response.aiMove.y }
    moveHistory.value.push(aiMove)
    gameState.lastMove = aiMove

    // 检查游戏状态
    gameState.gameStatus = getGameStatus(gameState.board, aiMove)
    
    if (gameState.gameStatus !== GameStatus.PLAYING) {
      handleGameEnd()
      return
    }

    // 切换回玩家回合
    gameState.currentPlayer = Player.HUMAN

  } catch (error: any) {
    showError(error.message || 'AI移动失败')
    // 如果AI移动失败，切换回玩家回合
    gameState.currentPlayer = Player.HUMAN
  } finally {
    isAiThinking.value = false
  }
}

function isValidAIMove(x: number, y: number): boolean {
  return x >= 0 && x < BOARD_SIZE && 
         y >= 0 && y < BOARD_SIZE && 
         gameState.board[y][x] === Player.NONE
}

function handleGameEnd() {
  // 更新统计数据
  statistics.totalGames++
  
  switch (gameState.gameStatus) {
    case GameStatus.HUMAN_WIN:
      statistics.humanWins++
      gameState.winner = Player.HUMAN
      showSuccess('恭喜！您获得了胜利！')
      break
    case GameStatus.AI_WIN:
      statistics.aiWins++
      gameState.winner = Player.AI
      showSuccess('AI获得了胜利！再接再厉！')
      break
    case GameStatus.DRAW:
      statistics.draws++
      gameState.winner = Player.NONE
      showSuccess('平局！势均力敌！')
      break
  }

  // 保存统计数据到本地存储
  saveStatistics()
}

function handleRestart() {
  // 重置游戏状态
  Object.assign(gameState, createInitialGameState())
  moveHistory.value = []
  isProcessing.value = false
  isAiThinking.value = false
  clearMessages()
  
  showSuccess('游戏已重新开始')
}

function showError(message: string) {
  errorMessage.value = message
  setTimeout(clearError, 5000)
}

function showSuccess(message: string) {
  successMessage.value = message
  setTimeout(clearSuccess, 3000)
}

function clearError() {
  errorMessage.value = ''
}

function clearSuccess() {
  successMessage.value = ''
}

function clearMessages() {
  clearError()
  clearSuccess()
}

function saveStatistics() {
  try {
    localStorage.setItem('gomoku-statistics', JSON.stringify(statistics))
  } catch (error) {
    console.warn('无法保存统计数据到本地存储')
  }
}

function loadStatistics() {
  try {
    const saved = localStorage.getItem('gomoku-statistics')
    if (saved) {
      const data = JSON.parse(saved)
      Object.assign(statistics, data)
    }
  } catch (error) {
    console.warn('无法从本地存储加载统计数据')
  }
}

// 生命周期
onMounted(() => {
  loadStatistics()
  showSuccess('欢迎来到五子棋人机对战！')
})
</script>

<style scoped>
#app {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

.game-container {
  background: rgba(255, 255, 255, 0.95);
  border-radius: 15px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
  padding: 30px;
  max-width: 1200px;
  width: 100%;
}

.game-header {
  text-align: center;
  margin-bottom: 30px;
}

.game-header h1 {
  margin: 0 0 10px 0;
  color: #333;
  font-size: 2.5em;
  font-weight: 700;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.1);
}

.subtitle {
  margin: 0;
  color: #666;
  font-size: 1.1em;
  font-style: italic;
}

.game-content {
  display: flex;
  gap: 30px;
  align-items: flex-start;
  justify-content: center;
  margin-bottom: 30px;
}

.board-section {
  flex: 1;
  display: flex;
  justify-content: center;
  max-width: 600px;
}

.control-section {
  flex: 0 0 auto;
}

.game-footer {
  text-align: center;
  color: #666;
  font-size: 0.9em;
  border-top: 1px solid #eee;
  padding-top: 20px;
}

/* 提示消息样式 */
.error-toast,
.success-toast {
  position: fixed;
  top: 20px;
  right: 20px;
  z-index: 1000;
  max-width: 400px;
  cursor: pointer;
  animation: slideIn 0.3s ease;
}

.error-toast {
  background: #ffebee;
  border-left: 4px solid #f44336;
}

.success-toast {
  background: #e8f5e8;
  border-left: 4px solid #4caf50;
}

.error-content,
.success-content {
  display: flex;
  align-items: center;
  padding: 15px;
  border-radius: 5px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.error-icon,
.success-icon {
  font-size: 1.2em;
  margin-right: 10px;
}

.error-text,
.success-text {
  flex: 1;
  font-weight: 500;
}

.error-text {
  color: #d32f2f;
}

.success-text {
  color: #2e7d32;
}

.error-close,
.success-close {
  background: none;
  border: none;
  font-size: 1.5em;
  cursor: pointer;
  margin-left: 10px;
  opacity: 0.7;
  transition: opacity 0.2s;
}

.error-close:hover,
.success-close:hover {
  opacity: 1;
}

@keyframes slideIn {
  from {
    transform: translateX(100%);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .game-content {
    flex-direction: column;
    align-items: center;
  }
  
  .control-section {
    width: 100%;
    max-width: 350px;
  }
}

@media (max-width: 768px) {
  #app {
    padding: 10px;
  }
  
  .game-container {
    padding: 20px;
  }
  
  .game-header h1 {
    font-size: 2em;
  }
  
  .game-content {
    gap: 20px;
  }
}

@media (max-width: 480px) {
  .game-container {
    padding: 15px;
  }
  
  .game-header h1 {
    font-size: 1.8em;
  }
  
  .error-toast,
  .success-toast {
    left: 10px;
    right: 10px;
    max-width: none;
  }
}
</style>