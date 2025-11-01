<template>
  <div class="ai-game-container">
    <div class="game-header">
      <button @click="backToHome" class="back-button">
        ← 返回首页
      </button>

      <div class="game-info">
        <h1>人机对战</h1>
        <div class="game-status">
          <span v-if="isGameActive" class="status playing">
            🎮 游戏进行中
          </span>
          <span v-else class="status waiting">
            🎯 准备开始
          </span>
        </div>
      </div>

      <div class="user-status">
        <div v-if="userStore.isAuthenticated" class="logged-in-user">
          <span class="user-icon">👤</span>
          <span class="username">{{ userStore.user?.username }}</span>
          <span class="save-indicator" title="游戏记录将自动保存">💾</span>
        </div>
        <div v-else class="guest-user">
          <span class="guest-icon">👤</span>
          <span class="guest-text">游客模式</span>
          <button @click="goToLogin" class="login-prompt">登录保存记录</button>
        </div>
      </div>
    </div>

    <!-- 难度选择区域 -->
    <div v-if="!isGameActive" class="difficulty-selection">
      <div class="difficulty-selector-card">
        <h2>选择游戏难度</h2>
        <div class="difficulty-description">
          <p>选择适合您的难度级别，享受不同强度的挑战！</p>
        </div>

        <div class="difficulty-options">
          <div
            v-for="level in difficultyLevels"
            :key="level.value"
            class="difficulty-card"
            :class="{ 'selected': selectedDifficulty === level.value }"
            @click="selectDifficulty(level.value)"
          >
            <div class="difficulty-icon">{{ level.icon }}</div>
            <h3>{{ level.name }}</h3>
            <p class="difficulty-desc">{{ level.description }}</p>
            <div class="difficulty-stats">
              <span class="stat">{{ level.responseTime }}</span>
              <span class="stat">{{ level.strength }}</span>
            </div>
          </div>
        </div>

        <button
          @click="startNewGame"
          :disabled="!selectedDifficulty"
          class="btn btn-primary start-game-btn"
        >
          开始游戏
        </button>
      </div>
    </div>

    <!-- 游戏区域 -->
    <div v-if="isGameActive" class="game-content">
      <!-- 游戏头部信息 -->
      <div class="game-header-info">
        <div class="difficulty-badge">
          <span class="difficulty-label">当前难度:</span>
          <span class="difficulty-value">{{ getCurrentDifficultyName() }}</span>
        </div>
        <button @click="backToDifficulty" class="btn-change-difficulty">
          更换难度
        </button>
      </div>

      <!-- 主要游戏区域 -->
      <div class="main-game-area">
        <!-- 左侧：棋盘区域 -->
        <div class="board-section">
          <div class="board-container">
            <!-- 当前回合提示 -->
            <div v-if="gameState.gameStatus === 'playing'" class="turn-indicator-banner">
              <div v-if="gameState.currentPlayer === Player.HUMAN" class="my-turn-banner">
                <div class="turn-icon">🎯</div>
                <div class="turn-text">
                  <h3>轮到你了！</h3>
                  <p v-if="moveCount === 0">你执黑子，请先落子</p>
                  <p v-else>点击棋盘空白处落子</p>
                </div>
              </div>
              <div v-else class="opponent-turn-banner">
                <div class="turn-icon">🤖</div>
                <div class="turn-text">
                  <h3>AI正在思考...</h3>
                  <p>AI正在分析最佳落子位置</p>
                </div>
              </div>
            </div>

            <div class="board-wrapper">
              <Board
                :board="gameState.board"
                :current-player="gameState.currentPlayer"
                :last-move="gameState.lastMove"
                :can-move="canMove"
                @move="handlePlayerMove"
              />
            </div>

            <!-- 游戏控制按钮 -->
            <div class="game-controls">
              <button
                v-if="gameState.gameStatus === 'playing'"
                @click="restartGame"
                class="control-button restart-button"
              >
                重新开始
              </button>

              <button
                v-if="gameState.gameStatus !== 'playing'"
                @click="backToDifficulty"
                class="control-button new-game-button"
              >
                新游戏
              </button>
            </div>
          </div>
        </div>

        <!-- 右侧：统计信息面板 -->
        <div class="stats-sidebar">
          <!-- 游戏状态面板 -->
          <div class="status-panel">
            <h4>游戏状态</h4>
            <div class="status-content">
              <ControlPanel
                :game-status="gameState.gameStatus"
                :current-player="gameState.currentPlayer"
                :last-move="gameState.lastMove"
                :move-count="moveCount"
                :is-processing="isProcessing"
                :is-ai-thinking="isAiThinking"
                :statistics="statistics"
                @restart="restartGame"
              />
            </div>
          </div>

          <!-- AI性能统计 -->
          <div v-if="aiStats" class="ai-stats-panel">
            <h4>AI性能统计</h4>
            <div class="stats-grid">
              <div class="stat-item">
                <span class="stat-label">搜索节点:</span>
                <span class="stat-value">{{ formatNumber(aiStats.nodes_searched) }}</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">剪枝次数:</span>
                <span class="stat-value">{{ formatNumber(aiStats.cutoffs) }}</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">剪枝效率:</span>
                <span class="stat-value">{{ aiStats.pruning_efficiency?.toFixed(1) || 0 }}%</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">搜索时间:</span>
                <span class="stat-value">{{ aiStats.search_time || 'N/A' }}</span>
              </div>
            </div>
          </div>

          <!-- 游戏统计 -->
          <div class="game-stats-panel">
            <h4>战绩统计</h4>
            <div class="overall-stats">
              <div class="stat-row">
                <span class="stat-label">总场次:</span>
                <span class="stat-value">{{ statistics.totalGames }}</span>
              </div>
              <div class="stat-row">
                <span class="stat-label">胜利:</span>
                <span class="stat-value win">{{ statistics.humanWins }}</span>
              </div>
              <div class="stat-row">
                <span class="stat-label">失败:</span>
                <span class="stat-value lose">{{ statistics.aiWins }}</span>
              </div>
              <div class="stat-row">
                <span class="stat-label">平局:</span>
                <span class="stat-value draw">{{ statistics.draws }}</span>
              </div>
            </div>

            <div v-if="selectedDifficulty" class="difficulty-stats">
              <h5>当前难度</h5>
              <div class="stat-row">
                <span class="stat-label">胜利:</span>
                <span class="stat-value win">
                  {{ statistics.difficultyStats[selectedDifficulty]?.humanWins || 0 }}
                </span>
              </div>
              <div class="stat-row">
                <span class="stat-label">失败:</span>
                <span class="stat-value lose">
                  {{ statistics.difficultyStats[selectedDifficulty]?.aiWins || 0 }}
                </span>
              </div>
              <div class="stat-row">
                <span class="stat-label">平局:</span>
                <span class="stat-value draw">
                  {{ statistics.difficultyStats[selectedDifficulty]?.draws || 0 }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 游戏结果弹窗 -->
    <div v-if="showGameResult" class="game-result-overlay" @click="closeGameResult">
      <div class="game-result-modal" @click.stop>
        <div class="result-header">
          <h2>游戏结束</h2>
          <button @click="closeGameResult" class="close-btn">&times;</button>
        </div>

        <div class="result-content">
          <div class="result-icon">
            <span v-if="gameState.gameStatus === 'human_win'">🎉</span>
            <span v-else-if="gameState.gameStatus === 'ai_win'">🤖</span>
            <span v-else>🤝</span>
          </div>

          <div class="result-message">
            <h3 v-if="gameState.gameStatus === 'human_win'">恭喜获胜！</h3>
            <h3 v-else-if="gameState.gameStatus === 'ai_win'">AI获胜</h3>
            <h3 v-else>平局</h3>

            <p v-if="gameState.gameStatus === 'human_win'">
              您在{{ getCurrentDifficultyName() }}难度下成功击败了AI！
            </p>
            <p v-else-if="gameState.gameStatus === 'ai_win'">
              AI在{{ getCurrentDifficultyName() }}难度下展现了强大的棋力
            </p>
            <p v-else>
              您在{{ getCurrentDifficultyName() }}难度下与AI棋力相当
            </p>
          </div>

          <div class="result-stats">
            <div class="stat-item">
              <span class="label">总步数:</span>
              <span class="value">{{ moveCount }}</span>
            </div>
            <div class="stat-item">
              <span class="label">游戏难度:</span>
              <span class="value">{{ getCurrentDifficultyName() }}</span>
            </div>
          </div>

          <div class="result-actions">
            <button @click="restartGame" class="btn btn-primary">再来一局</button>
            <button @click="backToDifficulty" class="btn btn-secondary">更换难度</button>
            <button @click="backToHome" class="btn btn-info">返回首页</button>
          </div>
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
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Board from './Board.vue'
import ControlPanel from './ControlPanel.vue'
import { Player, GameStatus, type Position, type BoardState } from '../types/game'
import { aiApi, gameApi } from '../services/api'
import { useUserStore } from '../stores/user'
import {
  createInitialGameState,
  makeMove,
  getGameStatus,
  BOARD_SIZE
} from '../utils/gameLogic'

const router = useRouter()
const userStore = useUserStore()

// 响应式数据
const isGameActive = ref(false)
const selectedDifficulty = ref<'easy' | 'medium' | 'hard' | 'expert' | ''>('')
const showGameResult = ref(false)
const gameState = reactive<BoardState>(createInitialGameState())
const isProcessing = ref(false)
const isAiThinking = ref(false)
const errorMessage = ref('')
const successMessage = ref('')
const moveHistory = ref<Position[]>([])
const aiStats = ref<any>(null)

// 难度级别定义
const difficultyLevels = [
  {
    value: 'easy',
    name: '初级',
    icon: '🌱',
    description: '适合初学者，AI会使用简单启发式算法',
    responseTime: '< 100ms',
    strength: '⭐⭐'
  },
  {
    value: 'medium',
    name: '中级',
    icon: '🌿',
    description: '平衡的游戏体验，适合大多数玩家',
    responseTime: '100-500ms',
    strength: '⭐⭐⭐'
  },
  {
    value: 'hard',
    name: '高级',
    icon: '🌳',
    description: '具有挑战性的AI，需要仔细思考',
    responseTime: '500ms-2s',
    strength: '⭐⭐⭐⭐'
  },
  {
    value: 'expert',
    name: '专家',
    icon: '🌲',
    description: '最强AI，具备深度分析能力',
    responseTime: '2-5s',
    strength: '⭐⭐⭐⭐⭐'
  }
]

// 统计数据
const statistics = reactive({
  humanWins: 0,
  aiWins: 0,
  draws: 0,
  totalGames: 0,
  // 按难度统计
  difficultyStats: {
    easy: { humanWins: 0, aiWins: 0, draws: 0 },
    medium: { humanWins: 0, aiWins: 0, draws: 0 },
    hard: { humanWins: 0, aiWins: 0, draws: 0 },
    expert: { humanWins: 0, aiWins: 0, draws: 0 }
  }
})

// 计算属性
const moveCount = computed(() => moveHistory.value.length)
const canMove = computed(() =>
  gameState.gameStatus === GameStatus.PLAYING &&
  !isProcessing.value &&
  !isAiThinking.value
)

// 方法
function selectDifficulty(level: 'easy' | 'medium' | 'hard' | 'expert') {
  selectedDifficulty.value = level
}

function startNewGame() {
  if (!selectedDifficulty.value) return

  isGameActive.value = true
  restartGame()
  showSuccess(`已选择${getCurrentDifficultyName()}难度，祝您游戏愉快！`)
}

function backToHome() {
  router.push('/')
}

function goToLogin() {
  router.push('/login')
}

function backToDifficulty() {
  isGameActive.value = false
  showGameResult.value = false
  restartGame()
}

function getCurrentDifficultyName(): string {
  if (!selectedDifficulty.value) return '未选择'
  const level = difficultyLevels.find(l => l.value === selectedDifficulty.value)
  return level ? level.name : '未知'
}

function formatNumber(num: any): string {
  if (!num) return '0'
  const number = Number(num)
  return number.toLocaleString()
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
    aiStats.value = null // 清空之前的统计数据

    // 准备AI请求数据
    const request = {
      board: gameState.board,
      player: Player.AI,
      lastMove: gameState.lastMove || { x: 0, y: 0 }
    }

    // 调用AI接口，传递难度参数
    const response = await aiApi.getMove(
      request,
      selectedDifficulty.value as 'easy' | 'medium' | 'hard' | 'expert',
      true // 使用增强AI
    )

    // 更新AI统计信息
    if (response.stats) {
      aiStats.value = response.stats
    }

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
    console.error('AI move error:', error)

    // 如果是网络错误或API错误，尝试使用原始AI
    if (error.message?.includes('fetch') || error.message?.includes('API')) {
      showSuccess('使用经典AI模式')
      try {
        const request = {
          board: gameState.board,
          player: Player.AI,
          lastMove: gameState.lastMove || { x: 0, y: 0 }
        }

        const response = await aiApi.getMove(request, selectedDifficulty.value as any, false)

        if (!isValidAIMove(response.aiMove.x, response.aiMove.y)) {
          throw new Error('备用AI也返回了无效的移动')
        }

        if (!makeMove(gameState.board, response.aiMove.x, response.aiMove.y, Player.AI)) {
          throw new Error('备用AI移动执行失败')
        }

        const aiMove: Position = { x: response.aiMove.x, y: response.aiMove.y }
        moveHistory.value.push(aiMove)
        gameState.lastMove = aiMove

        gameState.gameStatus = getGameStatus(gameState.board, aiMove)

        if (gameState.gameStatus !== GameStatus.PLAYING) {
          handleGameEnd()
          return
        }

        gameState.currentPlayer = Player.HUMAN
      } catch (fallbackError: any) {
        throw new Error('AI服务不可用，请稍后重试')
      }
    } else {
      throw error
    }
  } finally {
    isAiThinking.value = false
  }
}

function isValidAIMove(x: number, y: number): boolean {
  return x >= 0 && x < BOARD_SIZE &&
         y >= 0 && y < BOARD_SIZE &&
         gameState.board[y][x] === Player.NONE
}

async function handleGameEnd() {
  // 更新总统计数据
  statistics.totalGames++

  // 更新按难度的统计数据
  if (selectedDifficulty.value) {
    const difficulty = selectedDifficulty.value as keyof typeof statistics.difficultyStats
    statistics.difficultyStats[difficulty].totalGames = (statistics.difficultyStats[difficulty].totalGames || 0) + 1
  }

  let winner: Player = Player.NONE
  switch (gameState.gameStatus) {
    case GameStatus.HUMAN_WIN:
      statistics.humanWins++
      winner = Player.HUMAN

      if (selectedDifficulty.value) {
        const difficulty = selectedDifficulty.value as keyof typeof statistics.difficultyStats
        statistics.difficultyStats[difficulty].humanWins++
      }
      break
    case GameStatus.AI_WIN:
      statistics.aiWins++
      winner = Player.AI

      if (selectedDifficulty.value) {
        const difficulty = selectedDifficulty.value as keyof typeof statistics.difficultyStats
        statistics.difficultyStats[difficulty].aiWins++
      }
      break
    case GameStatus.DRAW:
      statistics.draws++

      if (selectedDifficulty.value) {
        const difficulty = selectedDifficulty.value as keyof typeof statistics.difficultyStats
        statistics.difficultyStats[difficulty].draws++
      }
      break
  }

  gameState.winner = winner

  // 如果用户已登录，保存游戏记录
  if (userStore.isAuthenticated) {
    try {
      await saveGameRecord(winner)
    } catch (error) {
      console.error('保存游戏记录失败:', error)
      // 不显示错误给用户，避免影响游戏体验
    }
  }

  // 显示结果
  setTimeout(() => {
    showGameResult.value = true
  }, 1000)
}

async function saveGameRecord(winner: Player) {
  if (!userStore.isAuthenticated || !selectedDifficulty.value) return

  try {
    const gameRecord = {
      gameType: 'ai' as const,
      difficulty: selectedDifficulty.value,
      result: winner === Player.HUMAN ? 'win' : winner === Player.AI ? 'lose' : 'draw',
      moves: moveHistory.value,
      moveCount: moveCount.value,
      aiStats: aiStats.value
    }

    await gameApi.saveGameRecord(gameRecord)
  } catch (error) {
    console.error('保存游戏记录失败:', error)
    throw error
  }
}

function restartGame() {
  // 重置游戏状态
  Object.assign(gameState, createInitialGameState())
  moveHistory.value = []
  isProcessing.value = false
  isAiThinking.value = false
  aiStats.value = null
  showGameResult.value = false
  clearMessages()
}

function closeGameResult() {
  showGameResult.value = false
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

// 生命周期
onMounted(() => {
  // 设置默认难度
  selectedDifficulty.value = 'medium'
})
</script>

<style scoped>
.ai-game-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.game-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 30px;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  border-radius: 15px;
  padding: 20px;
}

.user-status {
  display: flex;
  align-items: center;
  gap: 10px;
}

.logged-in-user {
  display: flex;
  align-items: center;
  gap: 8px;
  background: rgba(255, 255, 255, 0.1);
  padding: 8px 12px;
  border-radius: 8px;
  backdrop-filter: blur(10px);
}

.user-icon, .guest-icon {
  font-size: 1.2rem;
}

.username {
  font-weight: 500;
  font-size: 0.9rem;
}

.save-indicator {
  font-size: 1rem;
  opacity: 0.8;
}

.guest-user {
  display: flex;
  align-items: center;
  gap: 8px;
}

.guest-text {
  font-size: 0.9rem;
  opacity: 0.8;
}

.login-prompt {
  background: rgba(255, 255, 255, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.3);
  color: white;
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 0.8rem;
  cursor: pointer;
  transition: all 0.2s ease;
}

.login-prompt:hover {
  background: rgba(255, 255, 255, 0.3);
  transform: translateY(-1px);
}

.back-button {
  background: rgba(255, 255, 255, 0.2);
  border: none;
  color: white;
  padding: 10px 20px;
  border-radius: 8px;
  cursor: pointer;
  font-size: 16px;
  transition: all 0.3s ease;
}

.back-button:hover {
  background: rgba(255, 255, 255, 0.3);
  transform: translateX(-2px);
}

.game-info h1 {
  color: white;
  font-size: 2.5rem;
  margin: 0;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.3);
}

.game-status {
  margin-top: 10px;
}

.status {
  padding: 8px 16px;
  border-radius: 20px;
  font-weight: 500;
  font-size: 14px;
}

.status.playing {
  background: rgba(34, 197, 94, 0.2);
  color: #16a34a;
  border: 1px solid rgba(34, 197, 94, 0.3);
}

.status.waiting {
  background: rgba(251, 191, 36, 0.2);
  color: #d97706;
  border: 1px solid rgba(251, 191, 36, 0.3);
}

/* 难度选择区域 */
.difficulty-selection {
  max-width: 900px;
  margin: 0 auto;
}

.difficulty-selector-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  padding: 40px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
}

.difficulty-selector-card h2 {
  text-align: center;
  color: #374151;
  margin-bottom: 20px;
  font-size: 1.8rem;
}

.difficulty-description {
  text-align: center;
  margin-bottom: 30px;
}

.difficulty-description p {
  color: #6b7280;
  font-size: 1.1rem;
  margin: 0;
}

.difficulty-options {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.difficulty-card {
  border: 2px solid #e5e7eb;
  border-radius: 15px;
  padding: 25px 20px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
  background: white;
  position: relative;
  overflow: hidden;
}

.difficulty-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, #667eea, #764ba2);
  transform: scaleX(0);
  transition: transform 0.3s ease;
}

.difficulty-card:hover {
  border-color: #667eea;
  transform: translateY(-5px);
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
}

.difficulty-card:hover::before {
  transform: scaleX(1);
}

.difficulty-card.selected {
  border-color: #667eea;
  background: linear-gradient(135deg, #f0f4ff 0%, #e8f2ff 100%);
  box-shadow: 0 10px 25px rgba(102, 126, 234, 0.2);
}

.difficulty-card.selected::before {
  transform: scaleX(1);
}

.difficulty-icon {
  font-size: 3rem;
  margin-bottom: 15px;
}

.difficulty-card h3 {
  margin: 0 0 10px 0;
  color: #374151;
  font-size: 1.3rem;
  font-weight: 600;
}

.difficulty-desc {
  color: #6b7280;
  font-size: 0.9rem;
  line-height: 1.4;
  margin: 0 0 15px 0;
}

.difficulty-stats {
  display: flex;
  justify-content: center;
  gap: 10px;
  flex-wrap: wrap;
}

.stat {
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 0.8rem;
  font-weight: 500;
}

.start-game-btn {
  width: 100%;
  padding: 15px;
  font-size: 1.1rem;
  font-weight: 600;
}

/* 游戏区域 */
.game-content {
  max-width: 1400px;
  margin: 0 auto;
}

.game-header-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 25px;
  padding: 18px 25px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 15px;
  border: 1px solid #e2e8f0;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.05);
}

.difficulty-badge {
  display: flex;
  align-items: center;
  gap: 10px;
}

.difficulty-label {
  color: #64748b;
  font-size: 0.95rem;
  font-weight: 500;
}

.difficulty-value {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 8px 16px;
  border-radius: 20px;
  font-size: 0.95rem;
  font-weight: 600;
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.3);
}

.btn-change-difficulty {
  background: white;
  color: #667eea;
  border: 2px solid #667eea;
  padding: 10px 20px;
  border-radius: 25px;
  font-size: 0.95rem;
  cursor: pointer;
  transition: all 0.3s ease;
  font-weight: 600;
}

.btn-change-difficulty:hover {
  background: #667eea;
  color: white;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

/* 主要游戏区域布局 */
.main-game-area {
  display: grid;
  grid-template-columns: 1fr 320px;
  gap: 25px;
  align-items: start;
}

.board-section {
  display: flex;
  justify-content: center;
  align-items: flex-start;
}

.board-container {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  padding: 25px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 600px;
}

.turn-indicator-banner {
  margin-bottom: 20px;
  text-align: center;
}

.my-turn-banner, .opponent-turn-banner {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 15px;
  padding: 15px;
  border-radius: 12px;
  margin-bottom: 20px;
}

.my-turn-banner {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  color: white;
}

.opponent-turn-banner {
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
  color: white;
}

.turn-icon {
  font-size: 2rem;
}

.turn-text h3 {
  margin: 0 0 5px 0;
  font-size: 1.2rem;
}

.turn-text p {
  margin: 0;
  opacity: 0.9;
}

.board-wrapper {
  margin-bottom: 20px;
}

.game-controls {
  display: flex;
  justify-content: center;
  gap: 15px;
}

.control-button {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 25px;
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.3s ease;
  font-weight: 500;
}

.control-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.2);
}

.control-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

/* 右侧统计面板 */
.stats-sidebar {
  display: flex;
  flex-direction: column;
  gap: 20px;
  width: 100%;
  max-width: 320px;
}

/* 状态面板 */
.status-panel {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 15px;
  padding: 20px;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.status-panel h4 {
  margin: 0 0 15px 0;
  color: #334155;
  font-size: 1.1rem;
  font-weight: 600;
  text-align: center;
  padding-bottom: 10px;
  border-bottom: 1px solid #e2e8f0;
}

.status-content {
  padding-top: 5px;
}

/* AI统计面板 */
.ai-stats-panel {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 15px;
  padding: 20px;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.ai-stats-panel h4 {
  margin: 0 0 15px 0;
  color: #334155;
  font-size: 1.1rem;
  font-weight: 600;
  text-align: center;
  padding-bottom: 10px;
  border-bottom: 1px solid #e2e8f0;
}

.stats-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 12px 8px;
  background: white;
  border-radius: 10px;
  border: 1px solid #e2e8f0;
  transition: all 0.2s ease;
}

.stat-item:hover {
  border-color: #667eea;
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.1);
}

.stat-label {
  color: #64748b;
  font-size: 0.8rem;
  font-weight: 500;
  margin-bottom: 4px;
  text-align: center;
}

.stat-value {
  color: #334155;
  font-size: 0.95rem;
  font-weight: 600;
  text-align: center;
}

/* 游戏统计面板 */
.game-stats-panel {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 15px;
  padding: 20px;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.game-stats-panel h4 {
  margin: 0 0 15px 0;
  color: #334155;
  font-size: 1.1rem;
  font-weight: 600;
  text-align: center;
  padding-bottom: 10px;
  border-bottom: 1px solid #e2e8f0;
}

.overall-stats, .difficulty-stats {
  margin-bottom: 15px;
}

.difficulty-stats h5 {
  margin: 15px 0 10px 0;
  color: #64748b;
  font-size: 0.95rem;
  font-weight: 600;
  text-align: center;
  padding-bottom: 5px;
  border-bottom: 1px solid #f1f5f9;
}

.stat-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 12px;
  border-bottom: 1px solid #f1f5f9;
  border-radius: 6px;
  transition: background-color 0.2s ease;
}

.stat-row:hover {
  background-color: #f8fafc;
}

.stat-row:last-child {
  border-bottom: none;
}

.stat-row .stat-label {
  color: #64748b;
  font-size: 0.9rem;
  font-weight: 500;
}

.stat-row .stat-value {
  font-weight: 600;
  font-size: 0.95rem;
}

.stat-row .stat-value.win {
  color: #10b981;
}

.stat-row .stat-value.lose {
  color: #ef4444;
}

.stat-row .stat-value.draw {
  color: #f59e0b;
}

/* 游戏结果弹窗 */
.game-result-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.game-result-modal {
  background: white;
  border-radius: 20px;
  padding: 30px;
  max-width: 500px;
  width: 90%;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.2);
  animation: slideUp 0.3s ease;
}

@keyframes slideUp {
  from {
    transform: translateY(50px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

.result-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.result-header h2 {
  color: #374151;
  margin: 0;
}

.close-btn {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #6b7280;
}

.result-content {
  text-align: center;
}

.result-icon {
  font-size: 4rem;
  margin-bottom: 20px;
}

.result-message h3 {
  color: #374151;
  margin-bottom: 10px;
  font-size: 1.5rem;
}

.result-message p {
  color: #6b7280;
  margin-bottom: 20px;
}

.result-stats {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 15px;
  margin-bottom: 30px;
  padding: 20px;
  background: #f9fafb;
  border-radius: 10px;
}

.result-actions {
  display: flex;
  gap: 10px;
  justify-content: center;
}

.btn {
  padding: 10px 20px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.3s ease;
  text-align: center;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-primary {
  background: #3b82f6;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: #2563eb;
  transform: translateY(-1px);
}

.btn-secondary {
  background: #6b7280;
  color: white;
}

.btn-secondary:hover:not(:disabled) {
  background: #4b5563;
}

.btn-info {
  background: #06b6d4;
  color: white;
}

.btn-info:hover:not(:disabled) {
  background: #0891b2;
}

/* Toast样式 */
.error-toast, .success-toast {
  position: fixed;
  top: 20px;
  right: 20px;
  z-index: 1000;
  cursor: pointer;
  animation: slideIn 0.3s ease;
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

.error-content, .success-content {
  background: white;
  padding: 15px 20px;
  border-radius: 10px;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.2);
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

/* 响应式设计 */
@media (max-width: 1200px) {
  .main-game-area {
    grid-template-columns: 1fr 280px;
    gap: 20px;
  }

  .stats-sidebar {
    max-width: 280px;
  }

  .board-container {
    max-width: 500px;
  }
}

@media (max-width: 1024px) {
  .main-game-area {
    grid-template-columns: 1fr;
    gap: 20px;
  }

  .stats-sidebar {
    order: -1;
    max-width: 100%;
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 15px;
  }

  .board-container {
    max-width: 550px;
    margin: 0 auto;
  }
}

@media (max-width: 768px) {
  .ai-game-container {
    padding: 15px;
  }

  .game-header {
    flex-direction: column;
    gap: 15px;
    text-align: center;
    padding: 15px;
  }

  .game-info h1 {
    font-size: 2rem;
  }

  .game-header-info {
    flex-direction: column;
    gap: 12px;
    padding: 15px;
    text-align: center;
  }

  .difficulty-badge {
    justify-content: center;
  }

  .board-container {
    padding: 20px 15px;
  }

  .turn-indicator-banner {
    margin-bottom: 15px;
  }

  .my-turn-banner, .opponent-turn-banner {
    padding: 12px;
    flex-direction: column;
    gap: 8px;
  }

  .turn-icon {
    font-size: 1.5rem;
  }

  .turn-text h3 {
    font-size: 1rem;
  }

  .turn-text p {
    font-size: 0.85rem;
  }

  .difficulty-options {
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 15px;
  }

  .difficulty-selector-card {
    padding: 25px 20px;
  }

  .stats-sidebar {
    grid-template-columns: 1fr;
    gap: 15px;
  }

  .stats-grid {
    grid-template-columns: 1fr 1fr;
    gap: 8px;
  }

  .stat-item {
    padding: 10px 6px;
  }

  .result-actions {
    flex-direction: column;
    gap: 10px;
  }

  .result-stats {
    grid-template-columns: 1fr;
    gap: 10px;
  }
}

@media (max-width: 480px) {
  .ai-game-container {
    padding: 10px;
  }

  .game-header {
    padding: 10px;
  }

  .game-info h1 {
    font-size: 1.75rem;
  }

  .game-header-info {
    padding: 12px 10px;
  }

  .board-container {
    padding: 15px 10px;
  }

  .difficulty-selector-card {
    padding: 20px 15px;
  }

  .difficulty-options {
    grid-template-columns: 1fr;
    gap: 12px;
  }

  .difficulty-card {
    padding: 20px 15px;
  }

  .stats-grid {
    grid-template-columns: 1fr;
  }

  .game-stats-panel, .ai-stats-panel, .status-panel {
    padding: 15px;
  }

  .stat-row {
    padding: 8px 10px;
  }

  .control-button {
    padding: 12px 18px;
    font-size: 0.85rem;
  }

  .btn-change-difficulty {
    padding: 8px 16px;
    font-size: 0.85rem;
  }
}
</style>