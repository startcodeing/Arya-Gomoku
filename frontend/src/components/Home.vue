<template>
  <div class="home-container">
    <div class="game-header">
      <h1>五子棋游戏</h1>
      <p class="subtitle">Gomoku Game</p>
    </div>

    <div class="game-modes">
      <div class="mode-card" @click="startAIGame">
        <div class="mode-icon">🤖</div>
        <h3>人机对战</h3>
        <p>与AI进行五子棋对战</p>
        <button class="mode-button">开始游戏</button>
      </div>

      <div class="mode-card" @click="goToPVP">
        <div class="mode-icon">👥</div>
        <h3>双人对战</h3>
        <p>与其他玩家在线对战</p>
        <button class="mode-button">进入房间</button>
      </div>
    </div>

    <!-- AI Game Section -->
    <div v-if="showAIGame" class="ai-game-section">
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
            @back="backToHome"
          />
        </div>
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
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import Board from './Board.vue'
import ControlPanel from './ControlPanel.vue'
import { Player, GameStatus, type Position, type BoardState } from '../types/game'
import { aiApi } from '../services/api'
import {
  createInitialGameState,
  makeMove,
  getGameStatus,
  BOARD_SIZE
} from '../utils/gameLogic'

const router = useRouter()

// 响应式数据
const showAIGame = ref(false)
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
function startAIGame() {
  showAIGame.value = true
  handleRestart()
}

function backToHome() {
  showAIGame.value = false
  handleRestart()
}

function goToPVP() {
  router.push('/pvp')
}

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
      showSuccess('AI获得了胜利！再试一次吧！')
      break
    case GameStatus.DRAW:
      statistics.draws++
      showSuccess('平局！棋力相当！')
      break
  }
}

function handleRestart() {
  // 重置游戏状态
  Object.assign(gameState, createInitialGameState())
  moveHistory.value = []
  isProcessing.value = false
  isAiThinking.value = false
  clearMessages()
}

function showError(message: string) {
  errorMessage.value = message
  setTimeout(() => {
    errorMessage.value = ''
  }, 3000)
}

function showSuccess(message: string) {
  successMessage.value = message
  setTimeout(() => {
    successMessage.value = ''
  }, 3000)
}

function clearError() {
  errorMessage.value = ''
}

function clearSuccess() {
  successMessage.value = ''
}

function clearMessages() {
  errorMessage.value = ''
  successMessage.value = ''
}
</script>

<style scoped>
.home-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.game-header {
  text-align: center;
  margin-bottom: 40px;
  color: white;
}

.game-header h1 {
  font-size: 3rem;
  margin-bottom: 10px;
  text-shadow: 2px 2px 4px rgba(0,0,0,0.3);
}

.subtitle {
  font-size: 1.2rem;
  opacity: 0.9;
}

.game-modes {
  display: flex;
  justify-content: center;
  gap: 40px;
  margin-bottom: 40px;
  flex-wrap: wrap;
}

.mode-card {
  background: white;
  border-radius: 20px;
  padding: 40px 30px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 10px 30px rgba(0,0,0,0.2);
  min-width: 250px;
}

.mode-card:hover {
  transform: translateY(-10px);
  box-shadow: 0 20px 40px rgba(0,0,0,0.3);
}

.mode-icon {
  font-size: 4rem;
  margin-bottom: 20px;
}

.mode-card h3 {
  font-size: 1.5rem;
  margin-bottom: 15px;
  color: #333;
}

.mode-card p {
  color: #666;
  margin-bottom: 25px;
  line-height: 1.5;
}

.mode-button {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 12px 30px;
  border-radius: 25px;
  font-size: 1rem;
  cursor: pointer;
  transition: all 0.3s ease;
}

.mode-button:hover {
  transform: scale(1.05);
  box-shadow: 0 5px 15px rgba(0,0,0,0.2);
}

.ai-game-section {
  background: white;
  border-radius: 20px;
  padding: 30px;
  box-shadow: 0 10px 30px rgba(0,0,0,0.2);
}

.game-content {
  display: flex;
  gap: 30px;
  justify-content: center;
  align-items: flex-start;
  flex-wrap: wrap;
}

.board-section {
  flex: 1;
  max-width: 600px;
}

.control-section {
  flex: 0 0 300px;
  min-width: 280px;
}

/* Toast styles */
.error-toast, .success-toast {
  position: fixed;
  top: 20px;
  right: 20px;
  z-index: 1000;
  cursor: pointer;
  animation: slideIn 0.3s ease;
}

.error-content, .success-content {
  background: white;
  padding: 15px 20px;
  border-radius: 10px;
  box-shadow: 0 5px 15px rgba(0,0,0,0.2);
  display: flex;
  align-items: center;
  gap: 10px;
  max-width: 300px;
}

.error-content {
  border-left: 4px solid #ff4757;
}

.success-content {
  border-left: 4px solid #2ed573;
}

.error-close, .success-close {
  background: none;
  border: none;
  font-size: 1.2rem;
  cursor: pointer;
  margin-left: auto;
  opacity: 0.7;
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

@media (max-width: 768px) {
  .game-modes {
    flex-direction: column;
    align-items: center;
  }
  
  .game-content {
    flex-direction: column;
  }
  
  .control-section {
    flex: none;
  }
}
</style>