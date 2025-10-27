<template>
  <div class="game-result-container">
    <div class="result-header">
      <button @click="goBack" class="back-button">
        ← 返回
      </button>
      
      <h1>游戏结果</h1>
      
      <div class="header-actions">
        <button @click="shareResult" class="share-button">
          📤 分享
        </button>
      </div>
    </div>

    <div v-if="isLoading" class="loading-section">
      <div class="loading-spinner"></div>
      <p>加载游戏结果中...</p>
    </div>

    <div v-else-if="errorMessage" class="error-section">
      <div class="error-icon">⚠️</div>
      <h2>加载失败</h2>
      <p>{{ errorMessage }}</p>
      <button @click="retryLoad" class="retry-button">重试</button>
    </div>

    <div v-else-if="gameResult" class="result-content">
      <!-- 游戏结果概览 -->
      <div class="result-overview">
        <div class="result-icon">
          <span v-if="isWinner">🎉</span>
          <span v-else-if="gameResult.winner">😔</span>
          <span v-else>🤝</span>
        </div>
        
        <div class="result-title">
          <h2 v-if="isWinner">恭喜获胜！</h2>
          <h2 v-else-if="gameResult.winner">对手获胜</h2>
          <h2 v-else>平局</h2>
          
          <p class="result-subtitle">{{ getResultDescription() }}</p>
        </div>
        
        <div class="result-badge" :class="getResultBadgeClass()">
          {{ getResultBadgeText() }}
        </div>
      </div>

      <!-- 游戏统计 -->
      <div class="game-statistics">
        <h3>游戏统计</h3>
        
        <div class="stats-grid">
          <div class="stat-card">
            <div class="stat-icon">⏱️</div>
            <div class="stat-info">
              <div class="stat-label">游戏时长</div>
              <div class="stat-value">{{ formatDuration(gameResult.duration) }}</div>
            </div>
          </div>
          
          <div class="stat-card">
            <div class="stat-icon">🎯</div>
            <div class="stat-info">
              <div class="stat-label">总步数</div>
              <div class="stat-value">{{ gameResult.moves.length }}</div>
            </div>
          </div>
          
          <div class="stat-card">
            <div class="stat-icon">⚡</div>
            <div class="stat-info">
              <div class="stat-label">平均用时</div>
              <div class="stat-value">{{ getAverageTimePerMove() }}</div>
            </div>
          </div>
          
          <div class="stat-card">
            <div class="stat-icon">🏆</div>
            <div class="stat-info">
              <div class="stat-label">获胜方式</div>
              <div class="stat-value">{{ getWinMethodText() }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- 玩家表现 -->
      <div class="player-performance">
        <h3>玩家表现</h3>
        
        <div class="players-stats">
          <div
            v-for="player in playersWithStats"
            :key="player.id"
            class="player-stat-card"
            :class="{ 
              'winner': player.id === gameResult.winner,
              'current-player': player.id === currentPlayer?.id
            }"
          >
            <div class="player-header">
              <div class="player-avatar">
                <span class="avatar-text">{{ player.name.charAt(0).toUpperCase() }}</span>
              </div>
              <div class="player-info">
                <div class="player-name">{{ player.name }}</div>
                <div class="player-color">
                  <span class="color-indicator" :class="player.color"></span>
                  {{ player.color === 'black' ? '黑子' : '白子' }}
                </div>
              </div>
              <div v-if="player.id === gameResult.winner" class="winner-crown">👑</div>
            </div>
            
            <div class="player-stats">
              <div class="player-stat">
                <span class="stat-label">落子数</span>
                <span class="stat-value">{{ player.moveCount }}</span>
              </div>
              <div class="player-stat">
                <span class="stat-label">平均用时</span>
                <span class="stat-value">{{ player.averageTime }}</span>
              </div>
              <div class="player-stat">
                <span class="stat-label">最快落子</span>
                <span class="stat-value">{{ player.fastestMove }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 棋谱回放 -->
      <div class="game-replay">
        <h3>棋谱回放</h3>
        
        <div class="replay-container">
          <div class="replay-board">
            <Board
              :board="replayBoard"
              :current-player="replayCurrentPlayer"
              :last-move="replayLastMove"
              :can-move="false"
              :highlight-moves="replayHighlightMoves"
              :show-coordinates="true"
              size="small"
            />
          </div>
          
          <div class="replay-controls">
            <div class="replay-info">
              <div class="current-move">
                第 {{ currentReplayStep + 1 }} / {{ gameResult.moves.length }} 步
              </div>
              <div v-if="currentMove" class="move-details">
                {{ getCurrentMovePlayer() }} 在 ({{ currentMove.x + 1 }}, {{ currentMove.y + 1 }}) 落子
              </div>
            </div>
            
            <div class="replay-buttons">
              <button
                @click="goToFirstMove"
                :disabled="currentReplayStep <= 0"
                class="replay-button"
                title="第一步"
              >
                ⏮️
              </button>
              
              <button
                @click="previousMove"
                :disabled="currentReplayStep <= 0"
                class="replay-button"
                title="上一步"
              >
                ⏪
              </button>
              
              <button
                @click="toggleAutoReplay"
                class="replay-button auto-button"
                :class="{ active: isAutoReplaying }"
                title="自动播放"
              >
                {{ isAutoReplaying ? '⏸️' : '▶️' }}
              </button>
              
              <button
                @click="nextMove"
                :disabled="currentReplayStep >= gameResult.moves.length - 1"
                class="replay-button"
                title="下一步"
              >
                ⏩
              </button>
              
              <button
                @click="goToLastMove"
                :disabled="currentReplayStep >= gameResult.moves.length - 1"
                class="replay-button"
                title="最后一步"
              >
                ⏭️
              </button>
            </div>
            
            <div class="replay-speed">
              <label>播放速度:</label>
              <select v-model="autoReplaySpeed" @change="updateAutoReplaySpeed">
                <option value="2000">慢速 (2秒)</option>
                <option value="1000">正常 (1秒)</option>
                <option value="500">快速 (0.5秒)</option>
                <option value="200">极快 (0.2秒)</option>
              </select>
            </div>
          </div>
        </div>
      </div>

      <!-- 移动历史 -->
      <div class="move-history">
        <h3>移动历史</h3>
        
        <div class="history-table">
          <div class="history-header">
            <div class="col-step">步数</div>
            <div class="col-player">玩家</div>
            <div class="col-position">位置</div>
            <div class="col-time">用时</div>
            <div class="col-action">操作</div>
          </div>
          
          <div class="history-body">
            <div
              v-for="(move, index) in gameResult.moves"
              :key="index"
              class="history-row"
              :class="{ 
                'highlighted': index === currentReplayStep,
                'winning-move': isWinningMove(move, index)
              }"
            >
              <div class="col-step">{{ index + 1 }}</div>
              <div class="col-player">
                <span class="player-name">{{ getPlayerName(move.playerId) }}</span>
                <span class="color-indicator" :class="getPlayerColor(move.playerId)"></span>
              </div>
              <div class="col-position">({{ move.x + 1 }}, {{ move.y + 1 }})</div>
              <div class="col-time">{{ formatMoveTime(move.timestamp, index) }}</div>
              <div class="col-action">
                <button
                  @click="goToMove(index)"
                  class="goto-button"
                  title="跳转到此步"
                >
                  📍
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="result-actions">
        <button @click="downloadPGN" class="action-button download-button">
          💾 下载棋谱
        </button>
        
        <button @click="playAgain" class="action-button play-again-button">
          🔄 再来一局
        </button>
        
        <button @click="backToLobby" class="action-button lobby-button">
          🏠 返回大厅
        </button>
      </div>
    </div>

    <!-- 分享弹窗 -->
    <div v-if="showShareModal" class="modal-overlay" @click="closeShareModal">
      <div class="modal share-modal" @click.stop>
        <h3>分享游戏结果</h3>
        
        <div class="share-content">
          <div class="share-preview">
            <div class="preview-text">{{ getShareText() }}</div>
          </div>
          
          <div class="share-options">
            <button @click="copyShareText" class="share-option">
              📋 复制文本
            </button>
            <button @click="shareAsImage" class="share-option">
              🖼️ 生成图片
            </button>
            <button @click="downloadReplay" class="share-option">
              📹 下载回放
            </button>
          </div>
        </div>
        
        <div class="modal-actions">
          <button @click="closeShareModal" class="close-button">关闭</button>
        </div>
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
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { usePvpStore } from '../stores/pvp'

const route = useRoute()
const router = useRouter()
const pvpStore = usePvpStore()

// 响应式数据
const gameResult = ref<GameResult | null>(null)
const isLoading = ref(true)
const errorMessage = ref('')
const successMessage = ref('')
const showShareModal = ref(false)

// 回放相关
const currentReplayStep = ref(-1) // -1 表示初始状态
const isAutoReplaying = ref(false)
const autoReplaySpeed = ref(1000)
const autoReplayInterval = ref<number | null>(null)
const replayHighlightMoves = ref<{x: number, y: number}[]>([])

// 计算属性
const currentPlayer = computed(() => pvpStore.currentPlayer)

const isWinner = computed(() => {
  return gameResult.value?.winner === currentPlayer.value?.id
})

const currentMove = computed(() => {
  if (!gameResult.value || currentReplayStep.value < 0) return null
  return gameResult.value.moves[currentReplayStep.value]
})

const replayBoard = computed(() => {
  if (!gameResult.value) {
    return Array(BOARD_SIZE).fill(null).map(() => Array(BOARD_SIZE).fill(0))
  }
  
  // 根据当前回放步数重建棋盘
  const board = Array(BOARD_SIZE).fill(null).map(() => Array(BOARD_SIZE).fill(0))
  
  for (let i = 0; i <= currentReplayStep.value; i++) {
    const move = gameResult.value.moves[i]
    if (move) {
      // 简单的颜色分配：奇数步黑子，偶数步白子
      board[move.y][move.x] = (i % 2 === 0) ? PLAYER_COLORS.BLACK : PLAYER_COLORS.WHITE
    }
  }
  
  return board
})

const replayCurrentPlayer = computed(() => {
  if (currentReplayStep.value < 0) return PLAYER_COLORS.BLACK
  return (currentReplayStep.value % 2 === 0) ? PLAYER_COLORS.WHITE : PLAYER_COLORS.BLACK
})

const replayLastMove = computed(() => {
  if (!currentMove.value) return null
  return { x: currentMove.value.x, y: currentMove.value.y }
})

const playersWithStats = computed(() => {
  if (!gameResult.value) return []
  
  const players = gameResult.value.players || []
  
  return players.map(player => {
    const playerMoves = gameResult.value!.moves.filter(move => move.playerId === player.id)
    const moveCount = playerMoves.length
    
    // 计算平均用时（简化计算）
    const averageTime = moveCount > 0 ? 
      Math.round(gameResult.value!.duration / moveCount) : 0
    
    // 找到最快落子时间（简化为固定值）
    const fastestMove = Math.max(1, Math.round(averageTime * 0.6))
    
    return {
      ...player,
      color: getPlayerColor(player.id),
      moveCount,
      averageTime: formatDuration(averageTime),
      fastestMove: formatDuration(fastestMove)
    }
  })
})

// 方法
async function loadGameResult() {
  const roomId = route.params.id as string
  
  if (!roomId) {
    errorMessage.value = '无效的游戏ID'
    isLoading.value = false
    return
  }
  
  try {
    isLoading.value = true
    errorMessage.value = ''
    
    // 尝试从store获取游戏结果
    const room = await pvpStore.getRoom(roomId)
    
    if (!room.game || room.game.status !== 'finished') {
      errorMessage.value = '游戏尚未结束或不存在'
      isLoading.value = false
      return
    }
    
    // 构建游戏结果
    gameResult.value = {
      winner: room.game.winner,
      winnerName: room.game.winner ? getPlayerName(room.game.winner) : undefined,
      reason: 'win', // 简化处理
      finalBoard: room.game.board,
      moves: room.game.moves,
      duration: calculateGameDuration(room.game),
      players: room.players
    }
    
    // 初始化回放到最后一步
    currentReplayStep.value = gameResult.value.moves.length - 1
    
  } catch (error: any) {
    errorMessage.value = error.message || '加载游戏结果失败'
  } finally {
    isLoading.value = false
  }
}

function retryLoad() {
  loadGameResult()
}

function goBack() {
  router.go(-1)
}

function getResultDescription(): string {
  if (!gameResult.value) return ''
  
  if (gameResult.value.winner) {
    return `${gameResult.value.winnerName} 获得胜利`
  } else {
    return '双方平局'
  }
}

function getResultBadgeClass(): string {
  if (!gameResult.value) return ''
  
  if (isWinner.value) return 'winner'
  if (gameResult.value.winner) return 'loser'
  return 'draw'
}

function getResultBadgeText(): string {
  if (!gameResult.value) return ''
  
  if (isWinner.value) return '胜利'
  if (gameResult.value.winner) return '失败'
  return '平局'
}

function getWinMethodText(): string {
  if (!gameResult.value?.winner) return '平局'
  
  switch (gameResult.value.reason) {
    case 'win':
      return '五子连珠'
    case 'draw':
      return '和棋'
    case 'disconnect':
      return '对手断线'
    default:
      return '获胜'
  }
}

function getAverageTimePerMove(): string {
  if (!gameResult.value || gameResult.value.moves.length === 0) return '0秒'
  
  const averageSeconds = Math.round(gameResult.value.duration / gameResult.value.moves.length)
  return formatDuration(averageSeconds)
}

function formatDuration(seconds: number): string {
  if (seconds < 60) {
    return `${seconds}秒`
  }
  
  const minutes = Math.floor(seconds / 60)
  const remainingSeconds = seconds % 60
  
  if (minutes < 60) {
    return `${minutes}分${remainingSeconds}秒`
  }
  
  const hours = Math.floor(minutes / 60)
  const remainingMinutes = minutes % 60
  return `${hours}时${remainingMinutes}分${remainingSeconds}秒`
}

function calculateGameDuration(game: any): number {
  // 简化计算：假设每步平均30秒
  return game.moves.length * 30
}

function getPlayerName(playerId: string): string {
  const player = gameResult.value?.players?.find(p => p.id === playerId)
  return player?.name || '未知玩家'
}

function getPlayerColor(playerId: string): 'black' | 'white' {
  if (!gameResult.value?.players) return 'black'
  
  const playerIndex = gameResult.value.players.findIndex(p => p.id === playerId)
  return playerIndex === 0 ? 'black' : 'white'
}

function getCurrentMovePlayer(): string {
  if (!currentMove.value) return ''
  return getPlayerName(currentMove.value.player)
}

function formatMoveTime(timestamp: string, index: number): string {
  // 简化处理：显示相对时间
  const seconds = (index + 1) * 30 // 假设每步30秒
  return formatDuration(seconds)
}

function isWinningMove(move: Move, index: number): boolean {
  // 简化处理：最后一步且有获胜者
  return index === (gameResult.value?.moves.length || 0) - 1 && !!gameResult.value?.winner
}

// 回放控制
function goToFirstMove() {
  currentReplayStep.value = -1
  stopAutoReplay()
}

function goToLastMove() {
  if (gameResult.value) {
    currentReplayStep.value = gameResult.value.moves.length - 1
  }
  stopAutoReplay()
}

function previousMove() {
  if (currentReplayStep.value > -1) {
    currentReplayStep.value--
  }
  stopAutoReplay()
}

function nextMove() {
  if (gameResult.value && currentReplayStep.value < gameResult.value.moves.length - 1) {
    currentReplayStep.value++
  }
  
  // 如果到达最后一步，停止自动播放
  if (gameResult.value && currentReplayStep.value >= gameResult.value.moves.length - 1) {
    stopAutoReplay()
  }
}

function goToMove(index: number) {
  currentReplayStep.value = index
  stopAutoReplay()
  
  // 高亮当前移动
  if (gameResult.value && index >= 0 && index < gameResult.value.moves.length) {
    const move = gameResult.value.moves[index]
    replayHighlightMoves.value = [{ x: move.x, y: move.y }]
    
    setTimeout(() => {
      replayHighlightMoves.value = []
    }, 2000)
  }
}

function toggleAutoReplay() {
  if (isAutoReplaying.value) {
    stopAutoReplay()
  } else {
    startAutoReplay()
  }
}

function startAutoReplay() {
  if (!gameResult.value) return
  
  isAutoReplaying.value = true
  
  autoReplayInterval.value = window.setInterval(() => {
    if (currentReplayStep.value >= gameResult.value!.moves.length - 1) {
      stopAutoReplay()
    } else {
      nextMove()
    }
  }, autoReplaySpeed.value)
}

function stopAutoReplay() {
  isAutoReplaying.value = false
  
  if (autoReplayInterval.value) {
    clearInterval(autoReplayInterval.value)
    autoReplayInterval.value = null
  }
}

function updateAutoReplaySpeed() {
  if (isAutoReplaying.value) {
    stopAutoReplay()
    startAutoReplay()
  }
}

// 分享功能
function shareResult() {
  showShareModal.value = true
}

function closeShareModal() {
  showShareModal.value = false
}

function getShareText(): string {
  if (!gameResult.value) return ''
  
  const winner = gameResult.value.winner ? getPlayerName(gameResult.value.winner) : '平局'
  const duration = formatDuration(gameResult.value.duration)
  const moves = gameResult.value.moves.length
  
  return `🎮 五子棋对战结果
🏆 获胜者: ${winner}
⏱️ 游戏时长: ${duration}
🎯 总步数: ${moves}步
📅 时间: ${new Date().toLocaleDateString('zh-CN')}`
}

async function copyShareText() {
  try {
    await navigator.clipboard.writeText(getShareText())
    showSuccess('分享文本已复制到剪贴板')
    closeShareModal()
  } catch (error) {
    showSuccess('复制失败，请手动复制')
  }
}

function shareAsImage() {
  // TODO: 实现生成图片功能
  showSuccess('图片生成功能开发中...')
  closeShareModal()
}

function downloadReplay() {
  // TODO: 实现下载回放功能
  showSuccess('回放下载功能开发中...')
  closeShareModal()
}

// 其他操作
function downloadPGN() {
  if (!gameResult.value) return
  
  // 生成简单的PGN格式
  let pgn = `[Event "五子棋对战"]
[Date "${new Date().toISOString().split('T')[0]}"]
[Result "${gameResult.value.winner ? '1-0' : '1/2-1/2'}"]

`
  
  gameResult.value.moves.forEach((move, index) => {
    const moveNumber = Math.floor(index / 2) + 1
    const player = index % 2 === 0 ? '1.' : '2.'
    pgn += `${moveNumber}${player} ${String.fromCharCode(97 + move.x)}${move.y + 1} `
    
    if (index % 4 === 3) pgn += '\n'
  })
  
  // 下载文件
  const blob = new Blob([pgn], { type: 'text/plain' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `gomoku-game-${new Date().toISOString().split('T')[0]}.pgn`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
  
  showSuccess('棋谱已下载')
}

function playAgain() {
  router.push('/pvp')
}

function backToLobby() {
  router.push('/pvp')
}

function showSuccess(message: string) {
  successMessage.value = message
  setTimeout(() => {
    successMessage.value = ''
  }, 3000)
}

function clearSuccess() {
  successMessage.value = ''
}

// 生命周期
onMounted(() => {
  loadGameResult()
})

onUnmounted(() => {
  stopAutoReplay()
})

// 监听回放步数变化
watch(currentReplayStep, () => {
  // 可以在这里添加音效或其他反馈
})
</script>

<style scoped>
.game-result-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.result-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  color: white;
}

.back-button, .share-button {
  background: rgba(255,255,255,0.2);
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 25px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.back-button:hover, .share-button:hover {
  background: rgba(255,255,255,0.3);
  transform: translateY(-2px);
}

.result-header h1 {
  font-size: 2.5rem;
  text-shadow: 2px 2px 4px rgba(0,0,0,0.3);
}

.loading-section, .error-section {
  text-align: center;
  color: white;
  padding: 60px 20px;
}

.loading-spinner {
  width: 50px;
  height: 50px;
  border: 4px solid rgba(255,255,255,0.3);
  border-top: 4px solid white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 20px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.error-icon {
  font-size: 4rem;
  margin-bottom: 20px;
}

.retry-button {
  background: linear-gradient(135deg, #28a745 0%, #20c997 100%);
  color: white;
  border: none;
  padding: 12px 24px;
  border-radius: 25px;
  cursor: pointer;
  margin-top: 20px;
  transition: all 0.3s ease;
}

.retry-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(0,0,0,0.2);
}

.result-content {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 30px;
}

.result-overview {
  background: white;
  border-radius: 20px;
  padding: 40px;
  text-align: center;
  box-shadow: 0 10px 30px rgba(0,0,0,0.2);
  display: flex;
  align-items: center;
  gap: 30px;
}

.result-icon {
  font-size: 5rem;
}

.result-title {
  flex: 1;
}

.result-title h2 {
  font-size: 2.5rem;
  margin-bottom: 10px;
  color: #333;
}

.result-subtitle {
  font-size: 1.2rem;
  color: #666;
}

.result-badge {
  padding: 15px 25px;
  border-radius: 25px;
  font-weight: 600;
  font-size: 1.1rem;
}

.result-badge.winner {
  background: linear-gradient(135deg, #28a745 0%, #20c997 100%);
  color: white;
}

.result-badge.loser {
  background: linear-gradient(135deg, #dc3545 0%, #c82333 100%);
  color: white;
}

.result-badge.draw {
  background: linear-gradient(135deg, #ffc107 0%, #fd7e14 100%);
  color: white;
}

.game-statistics, .player-performance, .game-replay, .move-history {
  background: white;
  border-radius: 20px;
  padding: 30px;
  box-shadow: 0 10px 30px rgba(0,0,0,0.2);
}

.game-statistics h3, .player-performance h3, .game-replay h3, .move-history h3 {
  font-size: 1.5rem;
  margin-bottom: 25px;
  color: #333;
  text-align: center;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
}

.stat-card {
  background: #f8f9fa;
  border-radius: 15px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 15px;
  transition: transform 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-5px);
}

.stat-icon {
  font-size: 2rem;
  width: 50px;
  text-align: center;
}

.stat-info {
  flex: 1;
}

.stat-label {
  display: block;
  font-size: 0.9rem;
  color: #666;
  margin-bottom: 5px;
}

.stat-value {
  display: block;
  font-size: 1.3rem;
  font-weight: 600;
  color: #333;
}

.players-stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
}

.player-stat-card {
  background: #f8f9fa;
  border-radius: 15px;
  padding: 20px;
  border: 2px solid transparent;
  transition: all 0.3s ease;
}

.player-stat-card.winner {
  border-color: #28a745;
  background: #f8fff9;
}

.player-stat-card.current-player {
  border-color: #667eea;
  background: #f0f4ff;
}

.player-header {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 15px;
  position: relative;
}

.player-avatar {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: bold;
  font-size: 1.2rem;
}

.player-info {
  flex: 1;
}

.player-name {
  font-size: 1.1rem;
  font-weight: 600;
  color: #333;
  margin-bottom: 5px;
}

.player-color {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 0.9rem;
  color: #666;
}

.color-indicator {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  border: 1px solid #ccc;
}

.color-indicator.black {
  background: #333;
}

.color-indicator.white {
  background: #fff;
}

.winner-crown {
  font-size: 1.5rem;
  position: absolute;
  top: -5px;
  right: 0;
}

.player-stats {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
}

.player-stat {
  text-align: center;
  padding: 10px;
  background: white;
  border-radius: 8px;
}

.player-stat .stat-label {
  display: block;
  font-size: 0.8rem;
  color: #666;
  margin-bottom: 3px;
}

.player-stat .stat-value {
  display: block;
  font-size: 1rem;
  font-weight: 600;
  color: #333;
}

.replay-container {
  display: grid;
  grid-template-columns: 1fr 400px;
  gap: 30px;
  align-items: start;
}

.replay-board {
  background: #f8f9fa;
  border-radius: 15px;
  padding: 20px;
  display: flex;
  justify-content: center;
}

.replay-controls {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.replay-info {
  text-align: center;
  padding: 15px;
  background: #f8f9fa;
  border-radius: 10px;
}

.current-move {
  font-size: 1.1rem;
  font-weight: 600;
  color: #333;
  margin-bottom: 5px;
}

.move-details {
  font-size: 0.9rem;
  color: #666;
}

.replay-buttons {
  display: flex;
  justify-content: center;
  gap: 10px;
  flex-wrap: wrap;
}

.replay-button {
  background: #667eea;
  color: white;
  border: none;
  padding: 10px 15px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 1rem;
}

.replay-button:hover:not(:disabled) {
  background: #5a6fd8;
  transform: translateY(-2px);
}

.replay-button:disabled {
  background: #ccc;
  cursor: not-allowed;
  transform: none;
}

.auto-button.active {
  background: #28a745;
}

.replay-speed {
  display: flex;
  align-items: center;
  gap: 10px;
  justify-content: center;
}

.replay-speed label {
  font-size: 0.9rem;
  color: #666;
}

.replay-speed select {
  padding: 5px 10px;
  border: 1px solid #ddd;
  border-radius: 5px;
  font-size: 0.9rem;
}

.history-table {
  background: #f8f9fa;
  border-radius: 10px;
  overflow: hidden;
}

.history-header {
  display: grid;
  grid-template-columns: 60px 1fr 100px 80px 60px;
  gap: 10px;
  padding: 15px;
  background: #667eea;
  color: white;
  font-weight: 600;
  font-size: 0.9rem;
}

.history-body {
  max-height: 400px;
  overflow-y: auto;
}

.history-row {
  display: grid;
  grid-template-columns: 60px 1fr 100px 80px 60px;
  gap: 10px;
  padding: 12px 15px;
  border-bottom: 1px solid #e1e5e9;
  transition: background-color 0.2s ease;
  font-size: 0.9rem;
}

.history-row:hover {
  background: #fff;
}

.history-row.highlighted {
  background: #fff3cd;
  border-color: #ffc107;
}

.history-row.winning-move {
  background: #d4edda;
  border-color: #28a745;
}

.col-step {
  text-align: center;
  font-weight: 600;
}

.col-player {
  display: flex;
  align-items: center;
  gap: 8px;
}

.col-position {
  font-family: monospace;
  text-align: center;
}

.col-time {
  text-align: center;
  color: #666;
}

.col-action {
  text-align: center;
}

.goto-button {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 1rem;
  padding: 2px;
  border-radius: 3px;
  transition: background-color 0.2s ease;
}

.goto-button:hover {
  background: #e1e5e9;
}

.result-actions {
  display: flex;
  justify-content: center;
  gap: 20px;
  flex-wrap: wrap;
  margin-top: 20px;
}

.action-button {
  padding: 15px 30px;
  border: none;
  border-radius: 25px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 1rem;
}

.download-button {
  background: linear-gradient(135deg, #6c757d 0%, #495057 100%);
  color: white;
}

.play-again-button {
  background: linear-gradient(135deg, #28a745 0%, #20c997 100%);
  color: white;
}

.lobby-button {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.action-button:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 20px rgba(0,0,0,0.2);
}

/* Share Modal */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  background: white;
  border-radius: 20px;
  padding: 30px;
  max-width: 500px;
  width: 90%;
  box-shadow: 0 20px 40px rgba(0,0,0,0.3);
}

.share-modal h3 {
  text-align: center;
  margin-bottom: 25px;
  color: #333;
}

.share-content {
  margin-bottom: 25px;
}

.share-preview {
  background: #f8f9fa;
  border-radius: 10px;
  padding: 20px;
  margin-bottom: 20px;
}

.preview-text {
  white-space: pre-line;
  font-family: monospace;
  font-size: 0.9rem;
  color: #333;
}

.share-options {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
  gap: 10px;
}

.share-option {
  background: #667eea;
  color: white;
  border: none;
  padding: 12px;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 0.9rem;
}

.share-option:hover {
  background: #5a6fd8;
  transform: translateY(-2px);
}

.modal-actions {
  display: flex;
  justify-content: center;
}

.close-button {
  background: #6c757d;
  color: white;
  border: none;
  padding: 12px 24px;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.close-button:hover {
  background: #5a6268;
  transform: translateY(-2px);
}

/* Success Toast */
.success-toast {
  position: fixed;
  top: 20px;
  right: 20px;
  z-index: 1000;
  cursor: pointer;
  animation: slideIn 0.3s ease;
}

.success-content {
  background: white;
  padding: 15px 20px;
  border-radius: 10px;
  box-shadow: 0 5px 15px rgba(0,0,0,0.2);
  display: flex;
  align-items: center;
  gap: 10px;
  max-width: 300px;
  border-left: 4px solid #2ed573;
}

.success-close {
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

@media (max-width: 1200px) {
  .replay-container {
    grid-template-columns: 1fr;
    gap: 20px;
  }
  
  .replay-controls {
    order: -1;
  }
}

@media (max-width: 768px) {
  .game-result-container {
    padding: 15px;
  }
  
  .result-header {
    flex-direction: column;
    gap: 15px;
    text-align: center;
  }
  
  .result-header h1 {
    font-size: 2rem;
  }
  
  .result-overview {
    flex-direction: column;
    text-align: center;
  }
  
  .result-title h2 {
    font-size: 2rem;
  }
  
  .stats-grid {
    grid-template-columns: 1fr;
  }
  
  .players-stats {
    grid-template-columns: 1fr;
  }
  
  .history-header, .history-row {
    grid-template-columns: 40px 1fr 80px 60px 40px;
    font-size: 0.8rem;
  }
  
  .result-actions {
    flex-direction: column;
    align-items: center;
  }
  
  .action-button {
    width: 100%;
    max-width: 300px;
  }
}
</style>