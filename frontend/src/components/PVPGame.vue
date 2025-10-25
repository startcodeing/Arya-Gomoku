<template>
  <div class="pvp-game-container">
    <div class="game-header">
      <button @click="leaveGame" class="back-button">
        ← 离开游戏
      </button>
      
      <div class="game-info">
        <h1>{{ room?.name || '五子棋对战' }}</h1>
        <div class="game-status">
          <span v-if="game?.status === 'playing'" class="status playing">
            🎮 游戏进行中
          </span>
          <span v-else-if="game?.status === 'finished'" class="status finished">
            🏁 游戏结束
          </span>
          <span v-else class="status waiting">
            ⏳ 等待开始
          </span>
        </div>
      </div>
      
      <div class="connection-status" :class="connectionStatus">
        <span class="status-dot"></span>
        {{ getConnectionStatusText() }}
      </div>
    </div>

    <div class="game-content">
      <!-- 左侧：游戏棋盘 -->
      <div class="board-section">
        <div class="board-container">
          <div class="board-wrapper">
            <Board
              :board="gameBoard"
              :current-player="currentPlayerColor"
              :last-move="lastMove"
              :can-move="canMakeMove"
              :highlight-moves="highlightMoves"
              @move="handleMove"
            />
          </div>
          
          <!-- 游戏控制按钮 -->
          <div class="game-controls">
            <button
              v-if="game?.status === 'playing'"
              @click="requestDraw"
              :disabled="isLoading"
              class="control-button draw-button"
            >
              求和
            </button>
            
            <button
              v-if="game?.status === 'playing'"
              @click="surrender"
              :disabled="isLoading"
              class="control-button surrender-button"
            >
              认输
            </button>
            
            <button
              v-if="game?.status === 'finished'"
              @click="viewResult"
              class="control-button result-button"
            >
              查看结果
            </button>
          </div>
        </div>
      </div>

      <!-- 右侧：游戏信息和聊天 -->
      <div class="info-section">
        <!-- 玩家信息 -->
        <div class="players-info">
          <h3>对战玩家</h3>
          <div class="players-list">
            <div
              v-for="player in gamePlayers"
              :key="player.id"
              class="player-info"
              :class="{ 
                'current-turn': isPlayerTurn(player),
                'current-player': player.id === currentPlayer?.id
              }"
            >
              <div class="player-avatar">
                <span class="avatar-text">{{ player.name.charAt(0).toUpperCase() }}</span>
              </div>
              <div class="player-details">
                <div class="player-name">{{ player.name }}</div>
                <div class="player-color">
                  <span class="color-indicator" :class="player.color"></span>
                  {{ player.color === 'black' ? '黑子' : '白子' }}
                </div>
                <div v-if="isPlayerTurn(player)" class="turn-indicator">
                  ⏰ 轮到你了
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 游戏统计 -->
        <div class="game-stats">
          <h3>游戏统计</h3>
          <div class="stats-grid">
            <div class="stat-item">
              <span class="stat-label">回合数</span>
              <span class="stat-value">{{ moveCount }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">游戏时长</span>
              <span class="stat-value">{{ formatGameDuration() }}</span>
            </div>
            <div v-if="game?.winner" class="stat-item">
              <span class="stat-label">获胜者</span>
              <span class="stat-value">{{ getWinnerName() }}</span>
            </div>
          </div>
        </div>

        <!-- 移动历史 -->
        <div class="move-history">
          <h3>移动历史</h3>
          <div class="history-list">
            <div
              v-for="(move, index) in recentMoves"
              :key="index"
              class="history-item"
              @click="highlightMove(move)"
            >
              <span class="move-number">{{ index + 1 }}</span>
              <span class="move-player">{{ getPlayerName(move.player) }}</span>
              <span class="move-position">({{ move.x + 1 }}, {{ move.y + 1 }})</span>
            </div>
            
            <div v-if="moves.length === 0" class="empty-history">
              还没有移动记录
            </div>
          </div>
        </div>

        <!-- 聊天区域 -->
        <div class="chat-section">
          <h3>聊天</h3>
          
          <div class="chat-messages" ref="chatMessagesRef">
            <div
              v-for="message in chatMessages"
              :key="message.id"
              class="chat-message"
              :class="{ 'own-message': message.playerId === currentPlayer?.id }"
            >
              <div class="message-header">
                <span class="sender-name">{{ message.playerName }}</span>
                <span class="message-time">{{ formatTime(message.timestamp) }}</span>
              </div>
              <div class="message-content">{{ message.message }}</div>
            </div>
            
            <div v-if="chatMessages.length === 0" class="empty-chat">
              <div class="empty-icon">💬</div>
              <p>开始聊天吧！</p>
            </div>
          </div>
          
          <div class="chat-input">
            <input
              v-model="newMessage"
              @keyup.enter="sendMessage"
              :disabled="!isConnected"
              type="text"
              placeholder="输入消息..."
              maxlength="200"
            />
            <button
              @click="sendMessage"
              :disabled="!newMessage.trim() || !isConnected"
              class="send-button"
            >
              发送
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 游戏结果弹窗 -->
    <div v-if="showGameResult" class="modal-overlay" @click="closeGameResult">
      <div class="modal game-result-modal" @click.stop>
        <h2>游戏结束</h2>
        
        <div class="result-content">
          <div class="result-icon">
            <span v-if="gameResult?.winner === currentPlayer?.id">🎉</span>
            <span v-else-if="gameResult?.winner">😔</span>
            <span v-else>🤝</span>
          </div>
          
          <div class="result-text">
            <h3 v-if="gameResult?.winner === currentPlayer?.id">恭喜获胜！</h3>
            <h3 v-else-if="gameResult?.winner">对手获胜</h3>
            <h3 v-else>平局</h3>
            
            <p v-if="gameResult?.reason === 'win'">五子连珠获胜</p>
            <p v-else-if="gameResult?.reason === 'draw'">双方同意平局</p>
            <p v-else-if="gameResult?.reason === 'disconnect'">对手断线</p>
          </div>
          
          <div class="result-stats">
            <div class="stat">
              <span class="label">游戏时长:</span>
              <span class="value">{{ formatDuration(gameResult?.duration || 0) }}</span>
            </div>
            <div class="stat">
              <span class="label">总步数:</span>
              <span class="value">{{ gameResult?.moves.length || 0 }}</span>
            </div>
          </div>
        </div>
        
        <div class="modal-actions">
          <button @click="closeGameResult" class="close-button">关闭</button>
          <button @click="backToLobby" class="lobby-button">返回大厅</button>
        </div>
      </div>
    </div>

    <!-- 加载遮罩 -->
    <div v-if="isLoading" class="loading-overlay">
      <div class="loading-spinner"></div>
      <p>{{ loadingText }}</p>
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
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { usePvpStore } from '../stores/pvp'
import type { Player, Room } from '../types/pvp'

const route = useRoute()
const router = useRouter()
const pvpStore = usePvpStore()

// 响应式数据
const newMessage = ref('')
const chatMessagesRef = ref<HTMLElement>()
const loadingText = ref('加载中...')
const successMessage = ref('')
const showGameResult = ref(false)
const gameResult = ref<GameResult | null>(null)
const highlightMoves = ref<{x: number, y: number}[]>([])
const gameStartTime = ref<Date | null>(null)
const gameDurationInterval = ref<number | null>(null)

// 计算属性
const room = computed(() => pvpStore.currentRoom)
const currentPlayer = computed(() => pvpStore.currentPlayer)
const game = computed(() => pvpStore.currentGame)
const chatMessages = computed(() => pvpStore.chatMessages)
const isConnected = computed(() => pvpStore.isConnected)
const isLoading = computed(() => pvpStore.isLoading)
const errorMessage = computed(() => pvpStore.error)
const connectionStatus = computed(() => pvpStore.connectionStatus)
const isMyTurn = computed(() => pvpStore.isMyTurn)

// 游戏相关计算属性
const gameBoard = computed(() => {
  if (!game.value?.board) {
    // 创建空棋盘
    return Array(BOARD_SIZE).fill(null).map(() => Array(BOARD_SIZE).fill(0))
  }
  return game.value.board
})

const moves = computed(() => game.value?.moves || [])
const moveCount = computed(() => moves.value.length)
const lastMove = computed(() => {
  const lastMoveData = moves.value[moves.value.length - 1]
  return lastMoveData ? { x: lastMoveData.x, y: lastMoveData.y } : null
})

const recentMoves = computed(() => {
  return moves.value.slice(-10).reverse() // 显示最近10步，倒序显示
})

const gamePlayers = computed(() => {
  if (!room.value?.players) return []
  
  return room.value.players.map(player => ({
    ...player,
    color: getPlayerColor(player.id)
  }))
})

const currentPlayerColor = computed(() => {
  if (!game.value?.currentPlayer) return PLAYER_COLORS.BLACK
  
  const currentGamePlayer = gamePlayers.value.find(p => p.id === game.value?.currentPlayer)
  return currentGamePlayer?.color === 'white' ? PLAYER_COLORS.WHITE : PLAYER_COLORS.BLACK
})

const canMakeMove = computed(() => {
  return game.value?.status === 'playing' && 
         isMyTurn.value && 
         isConnected.value && 
         !isLoading.value
})

// 方法
async function handleMove(x: number, y: number) {
  if (!canMakeMove.value) return
  
  try {
    loadingText.value = '移动中...'
    await pvpStore.makeMove(x, y)
    showSuccess('移动成功')
  } catch (error: any) {
    // Error is handled by the store
  }
}

async function leaveGame() {
  try {
    loadingText.value = '离开游戏中...'
    await pvpStore.leaveRoom()
    router.push('/pvp')
  } catch (error: any) {
    // Error is handled by the store
  }
}

function requestDraw() {
  // TODO: Implement draw request functionality
  showSuccess('求和请求已发送')
}

function surrender() {
  // TODO: Implement surrender functionality
  showSuccess('已认输')
}

function viewResult() {
  if (room.value) {
    router.push(`/result/${room.value.id}`)
  }
}

function backToLobby() {
  closeGameResult()
  if (room.value) {
    router.push(`/room/${room.value.id}`)
  }
}

function closeGameResult() {
  showGameResult.value = false
  gameResult.value = null
}

function sendMessage() {
  if (!newMessage.value.trim()) return
  
  pvpStore.sendChatMessage(newMessage.value)
  newMessage.value = ''
}

function highlightMove(move: Move) {
  highlightMoves.value = [{ x: move.x, y: move.y }]
  
  // 清除高亮
  setTimeout(() => {
    highlightMoves.value = []
  }, 2000)
}

function getPlayerColor(playerId: string): 'black' | 'white' {
  // 简单的颜色分配逻辑：第一个玩家黑子，第二个玩家白子
  if (!room.value?.players) return 'black'
  
  const playerIndex = room.value.players.findIndex(p => p.id === playerId)
  return playerIndex === 0 ? 'black' : 'white'
}

function isPlayerTurn(player: Player): boolean {
  return game.value?.currentPlayer === player.id
}

function getPlayerName(playerId: string): string {
  const player = room.value?.players.find(p => p.id === playerId)
  return player?.name || '未知玩家'
}

function getWinnerName(): string {
  if (!game.value?.winner) return ''
  return getPlayerName(game.value.winner)
}

function getConnectionStatusText(): string {
  switch (connectionStatus.value) {
    case 'connected':
      return '已连接'
    case 'connecting':
      return '连接中'
    case 'reconnecting':
      return '重连中'
    case 'disconnected':
      return '已断开'
    case 'error':
      return '连接错误'
    default:
      return '未知状态'
  }
}

function formatTime(timestamp: string): string {
  const date = new Date(timestamp)
  return date.toLocaleTimeString('zh-CN', { 
    hour: '2-digit', 
    minute: '2-digit' 
  })
}

function formatGameDuration(): string {
  if (!gameStartTime.value) return '00:00'
  
  const now = new Date()
  const duration = Math.floor((now.getTime() - gameStartTime.value.getTime()) / 1000)
  return formatDuration(duration)
}

function formatDuration(seconds: number): string {
  const minutes = Math.floor(seconds / 60)
  const remainingSeconds = seconds % 60
  return `${minutes.toString().padStart(2, '0')}:${remainingSeconds.toString().padStart(2, '0')}`
}

function showSuccess(message: string) {
  successMessage.value = message
  setTimeout(() => {
    successMessage.value = ''
  }, 3000)
}

function clearError() {
  pvpStore.clearError()
}

function clearSuccess() {
  successMessage.value = ''
}

function scrollChatToBottom() {
  nextTick(() => {
    if (chatMessagesRef.value) {
      chatMessagesRef.value.scrollTop = chatMessagesRef.value.scrollHeight
    }
  })
}

function startGameTimer() {
  if (game.value?.startedAt) {
    gameStartTime.value = new Date(game.value.startedAt)
  } else {
    gameStartTime.value = new Date()
  }
  
  // 每秒更新游戏时长
  gameDurationInterval.value = window.setInterval(() => {
    // 触发重新计算
  }, 1000)
}

function stopGameTimer() {
  if (gameDurationInterval.value) {
    clearInterval(gameDurationInterval.value)
    gameDurationInterval.value = null
  }
}

// 监听聊天消息变化
watch(chatMessages, () => {
  scrollChatToBottom()
}, { deep: true })

// 监听游戏状态变化
watch(() => game.value?.status, (newStatus, oldStatus) => {
  if (newStatus === 'playing' && oldStatus !== 'playing') {
    startGameTimer()
    showSuccess('游戏开始了！')
  } else if (newStatus === 'finished' && oldStatus === 'playing') {
    stopGameTimer()
    // 显示游戏结果
    if (game.value) {
      gameResult.value = {
        winner: game.value.winner,
        winnerName: game.value.winner ? getPlayerName(game.value.winner) : undefined,
        reason: 'win', // 默认为获胜
        finalBoard: game.value.board,
        moves: game.value.moves,
        duration: gameStartTime.value ? 
          Math.floor((new Date().getTime() - gameStartTime.value.getTime()) / 1000) : 0
      }
      showGameResult.value = true
    }
  }
})

// 生命周期
onMounted(async () => {
  const roomId = route.params.id as string
  
  if (!roomId) {
    router.push('/pvp')
    return
  }
  
  // 如果不在房间中或房间不匹配，尝试获取房间信息
  if (!room.value || room.value.id !== roomId) {
    try {
      loadingText.value = '加载游戏信息...'
      const roomData = await pvpStore.getRoom(roomId)
      
      // 检查当前玩家是否在房间中
      const isPlayerInRoom = roomData.players.some(p => p.id === currentPlayer.value?.id)
      
      if (!isPlayerInRoom) {
        router.push('/pvp')
        return
      }
      
      // 如果游戏还没开始，跳转到房间大厅
      if (roomData.status === 'waiting') {
        router.push(`/room/${roomId}`)
        return
      }
      
    } catch (error) {
      router.push('/pvp')
      return
    }
  }
  
  // 如果游戏正在进行，启动计时器
  if (game.value?.status === 'playing') {
    startGameTimer()
  }
  
  scrollChatToBottom()
})

onUnmounted(() => {
  stopGameTimer()
})
</script>

<style scoped>
.pvp-game-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.game-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  color: white;
}

.back-button {
  background: rgba(255,255,255,0.2);
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 25px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.back-button:hover {
  background: rgba(255,255,255,0.3);
  transform: translateX(-5px);
}

.game-info {
  text-align: center;
  flex: 1;
}

.game-info h1 {
  font-size: 2.5rem;
  margin-bottom: 10px;
  text-shadow: 2px 2px 4px rgba(0,0,0,0.3);
}

.game-status {
  font-size: 1.2rem;
}

.status {
  padding: 8px 16px;
  border-radius: 20px;
  background: rgba(255,255,255,0.2);
}

.status.playing {
  background: rgba(40, 167, 69, 0.8);
}

.status.finished {
  background: rgba(255, 193, 7, 0.8);
}

.status.waiting {
  background: rgba(108, 117, 125, 0.8);
}

.connection-status {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  border-radius: 20px;
  background: rgba(255,255,255,0.2);
  font-size: 0.9rem;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #dc3545;
}

.connection-status.connected .status-dot {
  background: #28a745;
}

.connection-status.connecting .status-dot,
.connection-status.reconnecting .status-dot {
  background: #ffc107;
  animation: pulse 1s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.game-content {
  max-width: 1400px;
  margin: 0 auto;
  display: grid;
  grid-template-columns: 1fr 400px;
  gap: 30px;
}

.board-section {
  background: white;
  border-radius: 20px;
  padding: 30px;
  box-shadow: 0 10px 30px rgba(0,0,0,0.2);
  display: flex;
  flex-direction: column;
  align-items: center;
}

.board-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
}

.board-wrapper {
  max-width: 600px;
  width: 100%;
}

.game-controls {
  display: flex;
  gap: 15px;
  flex-wrap: wrap;
  justify-content: center;
}

.control-button {
  padding: 12px 24px;
  border: none;
  border-radius: 25px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.draw-button {
  background: linear-gradient(135deg, #ffc107 0%, #fd7e14 100%);
  color: white;
}

.surrender-button {
  background: linear-gradient(135deg, #dc3545 0%, #c82333 100%);
  color: white;
}

.result-button {
  background: linear-gradient(135deg, #28a745 0%, #20c997 100%);
  color: white;
}

.control-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(0,0,0,0.2);
}

.control-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.info-section {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.players-info, .game-stats, .move-history, .chat-section {
  background: white;
  border-radius: 15px;
  padding: 20px;
  box-shadow: 0 5px 15px rgba(0,0,0,0.1);
}

.players-info h3, .game-stats h3, .move-history h3, .chat-section h3 {
  margin-bottom: 15px;
  color: #333;
  font-size: 1.2rem;
}

.players-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.player-info {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-radius: 10px;
  border: 2px solid transparent;
  transition: all 0.3s ease;
}

.player-info.current-player {
  background: #f0f4ff;
  border-color: #667eea;
}

.player-info.current-turn {
  background: #fff3cd;
  border-color: #ffc107;
}

.player-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: bold;
}

.player-details {
  flex: 1;
}

.player-name {
  font-weight: 600;
  color: #333;
  margin-bottom: 4px;
}

.player-color {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 0.9rem;
  color: #666;
  margin-bottom: 4px;
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

.turn-indicator {
  font-size: 0.8rem;
  color: #856404;
  font-weight: 600;
}

.stats-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 15px;
}

.stat-item {
  text-align: center;
  padding: 12px;
  background: #f8f9fa;
  border-radius: 8px;
}

.stat-label {
  display: block;
  font-size: 0.9rem;
  color: #666;
  margin-bottom: 5px;
}

.stat-value {
  display: block;
  font-size: 1.2rem;
  font-weight: 600;
  color: #333;
}

.history-list {
  max-height: 200px;
  overflow-y: auto;
}

.history-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px;
  border-radius: 6px;
  cursor: pointer;
  transition: background-color 0.2s ease;
  font-size: 0.9rem;
}

.history-item:hover {
  background: #f8f9fa;
}

.move-number {
  width: 24px;
  height: 24px;
  background: #667eea;
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.8rem;
  font-weight: 600;
}

.move-player {
  flex: 1;
  font-weight: 500;
}

.move-position {
  color: #666;
  font-family: monospace;
}

.empty-history {
  text-align: center;
  color: #adb5bd;
  padding: 20px;
  font-style: italic;
}

.chat-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 300px;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  border: 1px solid #e1e5e9;
  border-radius: 10px;
  padding: 15px;
  margin-bottom: 15px;
  background: #fafbfc;
  max-height: 250px;
}

.chat-message {
  margin-bottom: 12px;
  padding: 8px;
  border-radius: 8px;
  background: white;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
}

.chat-message.own-message {
  background: #e3f2fd;
  margin-left: 15px;
}

.message-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.sender-name {
  font-weight: 600;
  color: #333;
  font-size: 0.85rem;
}

.message-time {
  font-size: 0.75rem;
  color: #666;
}

.message-content {
  color: #333;
  line-height: 1.3;
  font-size: 0.9rem;
}

.empty-chat {
  text-align: center;
  color: #adb5bd;
  padding: 30px 15px;
}

.empty-icon {
  font-size: 2rem;
  margin-bottom: 10px;
}

.chat-input {
  display: flex;
  gap: 8px;
}

.chat-input input {
  flex: 1;
  padding: 10px;
  border: 2px solid #e1e5e9;
  border-radius: 20px;
  font-size: 0.9rem;
  transition: border-color 0.3s ease;
}

.chat-input input:focus {
  outline: none;
  border-color: #667eea;
}

.chat-input input:disabled {
  background: #f8f9fa;
  color: #6c757d;
}

.send-button {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 10px 16px;
  border-radius: 20px;
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 0.9rem;
}

.send-button:hover:not(:disabled) {
  transform: scale(1.05);
  box-shadow: 0 3px 10px rgba(0,0,0,0.2);
}

.send-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

/* Game Result Modal */
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

.game-result-modal h2 {
  text-align: center;
  margin-bottom: 25px;
  color: #333;
}

.result-content {
  text-align: center;
  margin-bottom: 25px;
}

.result-icon {
  font-size: 4rem;
  margin-bottom: 20px;
}

.result-text h3 {
  margin-bottom: 10px;
  color: #333;
}

.result-text p {
  color: #666;
  margin-bottom: 20px;
}

.result-stats {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 15px;
  margin-top: 20px;
}

.stat {
  padding: 15px;
  background: #f8f9fa;
  border-radius: 10px;
}

.stat .label {
  display: block;
  font-size: 0.9rem;
  color: #666;
  margin-bottom: 5px;
}

.stat .value {
  display: block;
  font-size: 1.2rem;
  font-weight: 600;
  color: #333;
}

.modal-actions {
  display: flex;
  gap: 15px;
}

.close-button, .lobby-button {
  flex: 1;
  padding: 12px;
  border: none;
  border-radius: 10px;
  cursor: pointer;
  font-weight: 600;
  transition: all 0.3s ease;
}

.close-button {
  background: #6c757d;
  color: white;
}

.lobby-button {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.close-button:hover, .lobby-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(0,0,0,0.2);
}

/* Loading overlay */
.loading-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  color: white;
}

.loading-spinner {
  width: 50px;
  height: 50px;
  border: 4px solid rgba(255,255,255,0.3);
  border-top: 4px solid white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 20px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
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

@media (max-width: 1200px) {
  .game-content {
    grid-template-columns: 1fr;
    gap: 20px;
  }
  
  .info-section {
    order: -1;
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 15px;
  }
  
  .chat-section {
    grid-column: 1 / -1;
  }
}

@media (max-width: 768px) {
  .pvp-game-container {
    padding: 15px;
  }
  
  .game-header {
    flex-direction: column;
    gap: 15px;
    text-align: center;
  }
  
  .game-info h1 {
    font-size: 2rem;
  }
  
  .info-section {
    grid-template-columns: 1fr;
  }
  
  .stats-grid {
    grid-template-columns: 1fr;
  }
}
</style>