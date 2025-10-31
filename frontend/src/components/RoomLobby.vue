<template>
  <div class="room-lobby-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-left">
          <button @click="leaveRoom" class="back-button">
            <span class="back-icon">←</span>
            返回大厅
          </button>
        </div>

        <div class="header-center">
          <div class="room-title">
            <h1>{{ room?.name || '房间大厅' }}</h1>
            <div class="room-badge">
              <span class="room-id">ID: {{ room?.id }}</span>
              <div class="connection-status" :class="connectionStatus">
                <span class="status-dot"></span>
                {{ getConnectionStatusText() }}
              </div>
            </div>
          </div>
        </div>

        <div class="header-right">
          <div class="invite-section">
            <div class="invite-link-container">
              <input
                ref="inviteLinkInput"
                :value="inviteLink"
                readonly
                class="invite-link-input"
                placeholder="生成邀请链接..."
              />
              <button
                @click="copyInviteLink"
                class="copy-button"
                :class="{ 'copied': copySuccess }"
                :disabled="!inviteLink"
              >
                <span class="copy-icon">{{ copySuccess ? '✓' : '📋' }}</span>
                {{ copySuccess ? '已复制!' : '复制' }}
              </button>
            </div>
            <p class="invite-hint">分享链接邀请朋友加入</p>
          </div>
        </div>
      </div>
    </div>

    <!-- 主要内容区域 -->
    <div class="main-content">
      <!-- 左侧：游戏区域 -->
      <div class="game-area">
        <!-- 玩家列表 -->
        <div class="players-section">
          <div class="section-header">
            <div class="section-icon">👥</div>
            <h2>玩家列表</h2>
            <div class="player-count">
              {{ room?.players.length || 0 }}/{{ room?.maxPlayers || 2 }}
            </div>
          </div>

          <div class="players-grid">
            <div
              v-for="(player, index) in allPlayerSlots"
              :key="player?.id || `empty-${index}`"
              class="player-card"
              :class="{
                'current-player': player?.id === currentPlayer?.id,
                'empty': !player,
                'ready': player?.isReady
              }"
            >
              <div v-if="player" class="player-content">
                <div class="player-avatar">
                  <span class="avatar-text">{{ player.name.charAt(0).toUpperCase() }}</span>
                  <div v-if="player.id === currentPlayer?.id" class="current-badge">你</div>
                </div>
                <div class="player-info">
                  <div class="player-name">{{ player.name }}</div>
                  <div class="player-details">
                    <div class="player-status">
                      <span v-if="player.isReady" class="status ready">✓ 已准备</span>
                      <span v-else class="status not-ready">⏳ 未准备</span>
                    </div>
                    <div v-if="player.color" class="player-color">
                      <span class="color-indicator" :class="player.color"></span>
                      {{ player.color === 'black' ? '黑子' : '白子' }}
                    </div>
                  </div>
                </div>
              </div>
              <div v-else class="empty-slot">
                <div class="empty-avatar">
                  <span class="empty-icon">+</span>
                </div>
                <div class="empty-info">
                  <div class="empty-title">等待玩家</div>
                  <div class="empty-desc">加入房间开始对战</div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 游戏控制 -->
        <div class="game-controls">
          <div v-if="room?.status === 'waiting'" class="waiting-controls">
            <button
              @click="toggleReady"
              :disabled="isLoading"
              class="ready-button"
              :class="{ 'ready': currentPlayer?.isReady }"
            >
              <span class="button-icon">{{ currentPlayer?.isReady ? '✓' : '⚡' }}</span>
              {{ currentPlayer?.isReady ? '取消准备' : '准备游戏' }}
            </button>

            <div v-if="canStartGame" class="start-section">
              <button
                @click="startGame"
                :disabled="isLoading"
                class="start-button"
              >
                <span v-if="isLoading" class="loading-spinner"></span>
                {{ isLoading ? '开始中...' : '开始游戏' }}
              </button>
            </div>

            <div v-else-if="room?.players.length >= 2" class="status-message waiting">
              <span class="status-icon">⏰</span>
              等待所有玩家准备...
            </div>

            <div v-else class="status-message insufficient">
              <span class="status-icon">👤</span>
              需要至少2名玩家才能开始游戏
            </div>
          </div>

          <div v-else-if="room?.status === 'playing'" class="playing-controls">
            <div class="game-status playing">
              <span class="status-icon">🎮</span>
              <div class="status-text">
                <div class="status-title">游戏进行中</div>
                <div class="status-desc">精彩对局正在进行</div>
              </div>
            </div>
            <button @click="goToGame" class="action-button primary">
              <span class="button-icon">🎯</span>
              进入游戏
            </button>
          </div>

          <div v-else class="finished-controls">
            <div class="game-status finished">
              <span class="status-icon">🏁</span>
              <div class="status-text">
                <div class="status-title">游戏已结束</div>
                <div class="status-desc">查看比赛结果</div>
              </div>
            </div>
            <button @click="viewResult" class="action-button secondary">
              <span class="button-icon">📊</span>
              查看结果
            </button>
          </div>
        </div>
      </div>

      <!-- 右侧：聊天区域 -->
      <div class="chat-area">
        <div class="chat-section">
          <div class="section-header">
            <div class="section-icon">💬</div>
            <h2>聊天室</h2>
            <div class="chat-status" :class="{ 'online': isConnected, 'offline': !isConnected }">
              <span class="status-dot"></span>
              {{ isConnected ? '在线' : '离线' }}
            </div>
          </div>

          <div class="chat-messages" ref="chatMessagesRef">
            <div
              v-for="message in chatMessages"
              :key="message.id"
              class="chat-message"
              :class="{ 'own-message': message.playerId === currentPlayer?.id }"
            >
              <div class="message-avatar">
                {{ message.playerName.charAt(0).toUpperCase() }}
              </div>
              <div class="message-content">
                <div class="message-header">
                  <span class="sender-name">{{ message.playerName }}</span>
                  <span class="message-time">{{ formatTime(message.timestamp) }}</span>
                </div>
                <div class="message-text">{{ message.message }}</div>
              </div>
            </div>

            <div v-if="chatMessages.length === 0" class="empty-chat">
              <div class="empty-icon">💬</div>
              <div class="empty-text">
                <div class="empty-title">还没有聊天消息</div>
                <div class="empty-desc">发送第一条消息开始聊天吧！</div>
              </div>
            </div>
          </div>

          <div class="chat-input-container">
            <div class="chat-input">
              <input
                v-model="newMessage"
                @keyup.enter="sendMessage"
                :disabled="!isConnected"
                type="text"
                placeholder="输入消息..."
                maxlength="200"
                class="message-input"
              />
              <button
                @click="sendMessage"
                :disabled="!newMessage.trim() || !isConnected"
                class="send-button"
              >
                <span class="send-icon">➤</span>
              </button>
            </div>
            <div class="input-hint">按 Enter 发送消息</div>
          </div>
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
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { usePvpStore } from '../stores/pvp'
import { getGlobalWebSocketService } from '../services/websocket'
import type { Player } from '../types/pvp'

const router = useRouter()
const route = useRoute()
const pvpStore = usePvpStore()

// 响应式数据
const newMessage = ref('')
const chatMessagesRef = ref<HTMLElement>()
const loadingText = ref('加载中...')
const inviteLinkInput = ref<HTMLInputElement>()
const copySuccess = ref(false)

// 计算属性
const room = computed(() => pvpStore.currentRoom)
const currentPlayer = computed(() => pvpStore.currentPlayer)
const chatMessages = computed(() => pvpStore.chatMessages)
const isConnected = computed(() => pvpStore.isConnected)
const isLoading = computed(() => pvpStore.isLoading)
const errorMessage = computed(() => pvpStore.error)
const connectionStatus = computed(() => pvpStore.connectionStatus)
const canStartGame = computed(() => pvpStore.canStartGame)

const successMessage = ref('')

// 邀请链接计算属性
const inviteLink = computed(() => {
  if (!room.value?.id) return ''
  const baseUrl = window.location.origin
  return `${baseUrl}/invite/${room.value.id}`
})

// 生成所有玩家槽位（包括空槽位）
const allPlayerSlots = computed(() => {
  const slots: (Player | null)[] = []
  const players = room.value?.players || []
  const maxPlayers = room.value?.maxPlayers || 2
  
  // 添加现有玩家
  for (let i = 0; i < players.length; i++) {
    slots.push(players[i])
  }
  
  // 添加空槽位
  for (let i = players.length; i < maxPlayers; i++) {
    slots.push(null)
  }
  
  return slots
})

// 方法
async function leaveRoom() {
  try {
    loadingText.value = '离开房间中...'
    await pvpStore.leaveRoom()
    router.push('/pvp')
  } catch (error: any) {
    // Error is handled by the store
  }
}

async function toggleReady() {
  try {
    pvpStore.toggleReady()
    showSuccess(currentPlayer.value?.isReady ? '已准备' : '取消准备')
  } catch (error: any) {
    // Error is handled by the store
  }
}

async function startGame() {
  try {
    loadingText.value = '开始游戏中...'
    await pvpStore.startGame()
    showSuccess('游戏开始！')
  } catch (error: any) {
    // Error is handled by the store
  }
}

function goToGame() {
  if (room.value) {
    router.push(`/game/${room.value.id}`)
  }
}

function viewResult() {
  if (room.value) {
    router.push(`/result/${room.value.id}`)
  }
}

function sendMessage() {
  if (!newMessage.value.trim()) return
  
  pvpStore.sendChatMessage(newMessage.value)
  newMessage.value = ''
}

async function copyInviteLink() {
  if (!inviteLink.value) return
  
  try {
    await navigator.clipboard.writeText(inviteLink.value)
    copySuccess.value = true
    showSuccess('邀请链接已复制到剪贴板')
    
    // 3秒后重置复制状态
    setTimeout(() => {
      copySuccess.value = false
    }, 3000)
  } catch (error) {
    // 如果剪贴板API不可用，使用传统方法
    if (inviteLinkInput.value) {
      inviteLinkInput.value.select()
      inviteLinkInput.value.setSelectionRange(0, 99999) // 移动端兼容
      try {
        document.execCommand('copy')
        copySuccess.value = true
        showSuccess('邀请链接已复制到剪贴板')
        setTimeout(() => {
          copySuccess.value = false
        }, 3000)
      } catch (err) {
        showSuccess('复制失败，请手动复制链接')
      }
    }
  }
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

// 自动滚动聊天消息到底部
function scrollChatToBottom() {
  nextTick(() => {
    if (chatMessagesRef.value) {
      chatMessagesRef.value.scrollTop = chatMessagesRef.value.scrollHeight
    }
  })
}

// 监听聊天消息变化，自动滚动到底部
watch(chatMessages, () => {
  scrollChatToBottom()
}, { deep: true })

// 监听房间状态变化
watch(() => room.value?.status, (newStatus, oldStatus) => {
  if (newStatus === 'playing' && oldStatus === 'waiting') {
    showSuccess('游戏开始了！')
  }
})

// 页面关闭时的清理函数
const handleBeforeUnload = (event: BeforeUnloadEvent) => {
  // 使用 sendBeacon API 确保在页面关闭时能可靠地发送离开房间请求
  if (currentRoom.value && currentPlayer.value) {
    const url = `http://localhost:8080/api/rooms/${currentRoom.value.id}/leave`
    const data = JSON.stringify({ playerId: currentPlayer.value.id })
    
    // 尝试使用 sendBeacon API（更可靠）
    if (navigator.sendBeacon) {
      const blob = new Blob([data], { type: 'application/json' })
      navigator.sendBeacon(url, blob)
    } else {
      // 降级到同步 XMLHttpRequest（作为备选方案）
      try {
        const xhr = new XMLHttpRequest()
        xhr.open('POST', url, false) // 同步请求
        xhr.setRequestHeader('Content-Type', 'application/json')
        xhr.send(data)
      } catch (error) {
        console.warn('页面关闭时离开房间失败:', error)
      }
    }
    
    // 立即断开WebSocket连接
    const ws = getGlobalWebSocketService()
    ws.disconnect()
  }
}

// 生命周期
onMounted(async () => {
  const roomId = route.params.id as string
  
  if (!roomId) {
    router.push('/pvp')
    return
  }
  
  // 首先从localStorage初始化Pinia store数据
  pvpStore.initializeFromLocalStorage()
  
  // 总是尝试重新获取房间信息，确保数据是最新的
  try {
    loadingText.value = '加载房间信息...'
    const roomData = await pvpStore.getRoom(roomId)
    
    // 检查当前玩家是否在房间中
    const isPlayerInRoom = roomData.players.some(p => p.id === currentPlayer.value?.id)
    
    if (!isPlayerInRoom) {
      // 如果玩家不在房间中，跳转到房间列表
      console.warn('当前玩家不在房间中，跳转到房间列表')
      router.push('/pvp')
      return
    }
    
    // 确保WebSocket连接正常
    const ws = getGlobalWebSocketService()
    
    // 设置WebSocket事件处理器，确保连接状态正确更新到Pinia store
    pvpStore.setupWebSocketEventHandlers(ws)
    
    if (!ws.isConnected() && currentPlayer.value) {
      console.log('WebSocket未连接，尝试重新连接...')
      try {
        await ws.connect(roomId, currentPlayer.value.id)
        console.log('WebSocket重连成功')
      } catch (wsError) {
        console.error('WebSocket连接失败:', wsError)
        // WebSocket连接失败不阻止页面显示，但会影响实时更新
      }
    }
    
  } catch (error) {
    console.error('获取房间信息失败:', error)
    // 如果获取房间失败，跳转到房间列表
    router.push('/pvp')
    return
  }
  
  // 监听页面关闭事件
  window.addEventListener('beforeunload', handleBeforeUnload)
  
  scrollChatToBottom()
})

onUnmounted(() => {
  // 移除页面关闭事件监听器
  window.removeEventListener('beforeunload', handleBeforeUnload)
  
  // 清理工作在组件销毁时进行
})
</script>

<style scoped>
.room-lobby-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 0;
}

/* 页面头部 */
.page-header {
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(20px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.2);
  padding: 25px 0;
}

.header-content {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 30px;
  display: grid;
  grid-template-columns: 1fr 2fr 1fr;
  gap: 30px;
  align-items: center;
}

.header-left, .header-center, .header-right {
  display: flex;
  align-items: center;
}

.header-center {
  justify-content: center;
}

.header-right {
  justify-content: flex-end;
}

.back-button {
  background: rgba(255, 255, 255, 0.2);
  color: white;
  border: 1px solid rgba(255, 255, 255, 0.3);
  padding: 12px 20px;
  border-radius: 25px;
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 1rem;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 8px;
  backdrop-filter: blur(10px);
}

.back-button:hover {
  background: rgba(255, 255, 255, 0.3);
  transform: translateX(-5px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.2);
}

.back-icon {
  font-size: 1.1rem;
}

.room-title {
  text-align: center;
}

.room-title h1 {
  color: white;
  font-size: 2.2rem;
  margin: 0 0 10px 0;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.3);
}

.room-badge {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 15px;
  flex-wrap: wrap;
}

.room-id {
  background: rgba(255, 255, 255, 0.2);
  color: white;
  padding: 6px 14px;
  border-radius: 20px;
  font-size: 0.9rem;
  font-weight: 500;
  backdrop-filter: blur(10px);
}

.connection-status {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.2);
  font-size: 0.9rem;
  font-weight: 500;
  backdrop-filter: blur(10px);
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

.invite-section {
  max-width: 350px;
}

.invite-link-container {
  display: flex;
  gap: 8px;
  align-items: center;
  margin-bottom: 8px;
}

.invite-link-input {
  flex: 1;
  padding: 10px 14px;
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.1);
  color: white;
  font-size: 0.85rem;
  text-align: center;
  backdrop-filter: blur(10px);
}

.invite-link-input::placeholder {
  color: rgba(255, 255, 255, 0.6);
}

.copy-button {
  background: rgba(255, 255, 255, 0.2);
  color: white;
  border: 1px solid rgba(255, 255, 255, 0.3);
  padding: 10px 14px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 0.85rem;
  font-weight: 500;
  white-space: nowrap;
  display: flex;
  align-items: center;
  gap: 6px;
  backdrop-filter: blur(10px);
}

.copy-button:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.3);
  transform: translateY(-1px);
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
}

.copy-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.copy-button.copied {
  background: rgba(40, 167, 69, 0.8);
  border-color: rgba(40, 167, 69, 0.8);
}

.copy-icon {
  font-size: 0.9rem;
}

.invite-hint {
  font-size: 0.8rem;
  color: rgba(255, 255, 255, 0.7);
  margin: 0;
  text-align: center;
}

/* 主要内容区域 */
.main-content {
  max-width: 1400px;
  margin: 0 auto;
  padding: 40px 30px;
  display: grid;
  grid-template-columns: 1fr 400px;
  gap: 30px;
  align-items: start;
}

/* 区域标题 */
.section-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 25px;
  position: relative;
}

.section-icon {
  font-size: 1.8rem;
  background: rgba(255, 255, 255, 0.9);
  width: 45px;
  height: 45px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
}

.section-header h2 {
  color: white;
  font-size: 1.4rem;
  margin: 0;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.3);
  flex: 1;
}

.player-count, .chat-status {
  background: rgba(255, 255, 255, 0.2);
  color: white;
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 0.9rem;
  font-weight: 500;
  backdrop-filter: blur(10px);
  display: flex;
  align-items: center;
  gap: 6px;
}

.chat-status.online .status-dot {
  background: #28a745;
}

.chat-status.offline .status-dot {
  background: #dc3545;
}

/* 游戏区域 */
.game-area {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 20px;
  padding: 30px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.players-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.player-card {
  border: 2px solid #e5e7eb;
  border-radius: 16px;
  padding: 20px;
  transition: all 0.3s ease;
  background: #fafbfc;
  position: relative;
  overflow: hidden;
}

.player-card::before {
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

.player-card:hover {
  border-color: #667eea;
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1);
}

.player-card:hover::before {
  transform: scaleX(1);
}

.player-card.current-player {
  border-color: #667eea;
  background: linear-gradient(135deg, #f0f4ff 0%, #e8f2ff 100%);
  box-shadow: 0 8px 25px rgba(102, 126, 234, 0.15);
}

.player-card.current-player::before {
  transform: scaleX(1);
}

.player-card.ready {
  border-color: #10b981;
  background: #f0fdf4;
}

.player-card.empty {
  border-style: dashed;
  border-color: #d1d5db;
  background: #f9fafb;
  opacity: 0.7;
}

.player-content {
  display: flex;
  align-items: center;
  gap: 15px;
}

.player-avatar {
  position: relative;
  width: 60px;
  height: 60px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: bold;
  font-size: 1.4rem;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
  border: 3px solid white;
}

.current-badge {
  position: absolute;
  top: -5px;
  right: -5px;
  background: #10b981;
  color: white;
  font-size: 0.7rem;
  font-weight: bold;
  padding: 3px 8px;
  border-radius: 10px;
  border: 2px solid white;
  box-shadow: 0 2px 8px rgba(16, 185, 129, 0.4);
}

.player-info {
  flex: 1;
}

.player-name {
  font-weight: 700;
  color: #1f2937;
  font-size: 1.1rem;
  margin-bottom: 8px;
}

.player-details {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.player-status .status {
  padding: 3px 10px;
  border-radius: 12px;
  font-size: 0.8rem;
  font-weight: 600;
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.status.ready {
  background: #dcfce7;
  color: #16a34a;
}

.status.not-ready {
  background: #fef3c7;
  color: #d97706;
}

.player-color {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.85rem;
  color: #6b7280;
  font-weight: 500;
}

.color-indicator {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  border: 2px solid #e5e7eb;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.color-indicator.black {
  background: #1f2937;
}

.color-indicator.white {
  background: #f9fafb;
}

.empty-slot {
  text-align: center;
  padding: 20px 0;
}

.empty-avatar {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  border: 3px dashed #d1d5db;
  background: #f3f4f6;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 12px;
}

.empty-icon {
  font-size: 1.8rem;
  color: #9ca3af;
}

.empty-info .empty-title {
  font-weight: 600;
  color: #6b7280;
  margin-bottom: 4px;
}

.empty-info .empty-desc {
  font-size: 0.85rem;
  color: #9ca3af;
}

/* 游戏控制 */
.game-controls {
  border-top: 1px solid #e5e7eb;
  padding-top: 25px;
}

.ready-button, .start-button, .action-button {
  width: 100%;
  padding: 16px;
  border: none;
  border-radius: 12px;
  font-size: 1.1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  margin-bottom: 15px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
}

.ready-button {
  background: #6b7280;
  color: white;
  box-shadow: 0 4px 15px rgba(107, 114, 128, 0.3);
}

.ready-button.ready {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  box-shadow: 0 4px 15px rgba(16, 185, 129, 0.3);
}

.start-button, .action-button.primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
}

.action-button.secondary {
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
  color: white;
  box-shadow: 0 4px 15px rgba(245, 158, 11, 0.3);
}

.ready-button:hover:not(:disabled),
.start-button:hover:not(:disabled),
.action-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.2);
}

.ready-button:disabled,
.start-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.button-icon {
  font-size: 1.1rem;
}

.loading-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top: 2px solid white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.status-message {
  text-align: center;
  padding: 20px;
  border-radius: 12px;
  margin-bottom: 15px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  font-weight: 600;
}

.status-message.waiting {
  background: #fef3c7;
  color: #d97706;
}

.status-message.insufficient {
  background: #fee2e2;
  color: #dc2626;
}

.status-icon {
  font-size: 1.2rem;
}

.game-status {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 20px;
  border-radius: 12px;
  margin-bottom: 20px;
  background: #f8fafc;
}

.game-status.playing {
  background: linear-gradient(135deg, #dcfce7 0%, #d1fae5 100%);
  border: 1px solid #bbf7d0;
}

.game-status.finished {
  background: linear-gradient(135deg, #fef3c7 0%, #fed7aa 100%);
  border: 1px solid #fde68a;
}

.status-text {
  flex: 1;
}

.status-title {
  font-weight: 700;
  color: #1f2937;
  font-size: 1.1rem;
  margin-bottom: 2px;
}

.status-desc {
  color: #6b7280;
  font-size: 0.9rem;
}

/* 聊天区域 */
.chat-area {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 20px;
  padding: 30px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  height: fit-content;
  max-height: 700px;
  display: flex;
  flex-direction: column;
}

.chat-section {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 500px;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 20px;
  background: #fafbfc;
  max-height: 350px;
}

.chat-message {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
  padding: 12px;
  border-radius: 12px;
  background: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  transition: all 0.2s ease;
}

.chat-message:hover {
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
}

.chat-message.own-message {
  background: linear-gradient(135deg, #dbeafe 0%, #bfdbfe 100%);
  margin-left: 20px;
  border: 1px solid #93c5fd;
}

.message-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  font-weight: bold;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.9rem;
  flex-shrink: 0;
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.3);
}

.message-content {
  flex: 1;
}

.message-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}

.sender-name {
  font-weight: 600;
  color: #1f2937;
  font-size: 0.9rem;
}

.message-time {
  font-size: 0.75rem;
  color: #9ca3af;
  font-weight: 500;
}

.message-text {
  color: #374151;
  line-height: 1.5;
  font-size: 0.95rem;
}

.empty-chat {
  text-align: center;
  color: #9ca3af;
  padding: 40px 20px;
}

.empty-icon {
  font-size: 3rem;
  margin-bottom: 16px;
  opacity: 0.7;
}

.empty-text .empty-title {
  font-weight: 600;
  font-size: 1.1rem;
  margin-bottom: 6px;
  color: #6b7280;
}

.empty-text .empty-desc {
  font-size: 0.9rem;
  color: #9ca3af;
}

.chat-input-container {
  border-top: 1px solid #e5e7eb;
  padding-top: 20px;
}

.chat-input {
  display: flex;
  gap: 10px;
  margin-bottom: 8px;
}

.message-input {
  flex: 1;
  padding: 14px 16px;
  border: 2px solid #e5e7eb;
  border-radius: 25px;
  font-size: 1rem;
  transition: all 0.3s ease;
  background: white;
}

.message-input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.message-input:disabled {
  background: #f8fafc;
  color: #9ca3af;
}

.send-button {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 14px 20px;
  border-radius: 50%;
  cursor: pointer;
  transition: all 0.3s ease;
  width: 50px;
  height: 50px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
}

.send-button:hover:not(:disabled) {
  transform: scale(1.05);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4);
}

.send-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.send-icon {
  font-size: 1.1rem;
  font-weight: bold;
}

.input-hint {
  font-size: 0.8rem;
  color: #9ca3af;
  text-align: center;
}

/* Loading overlay */
.loading-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  color: white;
  backdrop-filter: blur(5px);
}

.loading-spinner {
  width: 60px;
  height: 60px;
  border: 4px solid rgba(255, 255, 255, 0.3);
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
  padding: 16px 20px;
  border-radius: 12px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.15);
  display: flex;
  align-items: center;
  gap: 12px;
  max-width: 350px;
  backdrop-filter: blur(10px);
}

.error-content {
  border-left: 4px solid #ef4444;
}

.success-content {
  border-left: 4px solid #10b981;
}

.error-icon, .success-icon {
  font-size: 1.2rem;
}

.error-text, .success-text {
  flex: 1;
  color: #1f2937;
  font-weight: 500;
}

.error-close, .success-close {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #6b7280;
  opacity: 0.7;
  transition: opacity 0.2s ease;
}

.error-close:hover, .success-close:hover {
  opacity: 1;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .main-content {
    grid-template-columns: 1fr 380px;
    gap: 25px;
  }

  .header-content {
    gap: 20px;
  }

  .invite-section {
    max-width: 300px;
  }
}

@media (max-width: 1024px) {
  .main-content {
    grid-template-columns: 1fr;
    gap: 25px;
  }

  .chat-area {
    order: -1;
    max-height: 500px;
  }

  .players-grid {
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  }
}

@media (max-width: 768px) {
  .page-header {
    padding: 20px 0;
  }

  .header-content {
    grid-template-columns: 1fr;
    gap: 20px;
    padding: 0 20px;
    text-align: center;
  }

  .room-title h1 {
    font-size: 1.8rem;
  }

  .header-left, .header-center, .header-right {
    justify-content: center;
  }

  .main-content {
    padding: 30px 20px;
  }

  .players-grid {
    grid-template-columns: 1fr;
  }

  .game-area, .chat-area {
    padding: 20px;
  }

  .section-header {
    margin-bottom: 20px;
  }

  .chat-section {
    min-height: 400px;
  }
}

@media (max-width: 480px) {
  .page-header {
    padding: 15px 0;
  }

  .header-content {
    padding: 0 15px;
  }

  .room-title h1 {
    font-size: 1.5rem;
  }

  .room-badge {
    flex-direction: column;
    gap: 10px;
  }

  .invite-link-container {
    flex-direction: column;
    gap: 10px;
  }

  .main-content {
    padding: 20px 15px;
  }

  .game-area, .chat-area {
    padding: 15px;
  }

  .players-grid {
    grid-template-columns: 1fr;
    gap: 15px;
  }

  .player-card {
    padding: 15px;
  }

  .player-avatar {
    width: 50px;
    height: 50px;
    font-size: 1.2rem;
  }

  .empty-avatar {
    width: 50px;
    height: 50px;
  }
}
</style>