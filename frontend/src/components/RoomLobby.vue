<template>
  <div class="room-lobby-container">
    <div class="header">
      <button @click="leaveRoom" class="back-button">
        ← 离开房间
      </button>
      <div class="room-info">
        <h1>{{ room?.name || '房间大厅' }}</h1>
        <p class="room-id">房间ID: {{ room?.id }}</p>
      </div>
      <div class="connection-status" :class="connectionStatus">
        <span class="status-dot"></span>
        {{ getConnectionStatusText() }}
      </div>
    </div>

    <div class="content">
      <!-- 左侧：玩家列表和游戏设置 -->
      <div class="left-panel">
        <!-- 玩家列表 -->
        <div class="players-section">
          <h2>玩家列表 ({{ room?.players.length || 0 }}/{{ room?.maxPlayers || 2 }})</h2>
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
                </div>
                <div class="player-info">
                  <div class="player-name">{{ player.name }}</div>
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
              <div v-else class="empty-slot">
                <div class="empty-avatar">+</div>
                <div class="empty-text">等待玩家加入</div>
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
              <span v-if="currentPlayer?.isReady">✓ 已准备</span>
              <span v-else>准备游戏</span>
            </button>
            
            <button
              v-if="canStartGame"
              @click="startGame"
              :disabled="isLoading"
              class="start-button"
            >
              <span v-if="isLoading">开始中...</span>
              <span v-else>开始游戏</span>
            </button>
            
            <div v-else-if="room?.players.length >= 2" class="start-hint">
              等待所有玩家准备...
            </div>
            
            <div v-else class="start-hint">
              等待更多玩家加入...
            </div>
          </div>
          
          <div v-else-if="room?.status === 'playing'" class="playing-controls">
            <div class="game-status">
              <span class="status-icon">🎮</span>
              游戏进行中
            </div>
            <button @click="goToGame" class="join-game-button">
              进入游戏
            </button>
          </div>
          
          <div v-else class="finished-controls">
            <div class="game-status">
              <span class="status-icon">🏁</span>
              游戏已结束
            </div>
            <button @click="viewResult" class="view-result-button">
              查看结果
            </button>
          </div>
        </div>
      </div>

      <!-- 右侧：聊天区域 -->
      <div class="right-panel">
        <div class="chat-section">
          <h2>聊天室</h2>
          
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
              <p>还没有聊天消息</p>
              <p>发送第一条消息开始聊天吧！</p>
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
import type { Player } from '../types/pvp'

const router = useRouter()
const route = useRoute()
const pvpStore = usePvpStore()

// 响应式数据
const newMessage = ref('')
const chatMessagesRef = ref<HTMLElement>()
const loadingText = ref('加载中...')

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

// 生命周期
onMounted(async () => {
  const roomId = route.params.id as string
  
  if (!roomId) {
    router.push('/pvp')
    return
  }
  
  // 如果不在房间中，尝试获取房间信息
  if (!room.value || room.value.id !== roomId) {
    try {
      loadingText.value = '加载房间信息...'
      const roomData = await pvpStore.getRoom(roomId)
      
      // 检查当前玩家是否在房间中
      const isPlayerInRoom = roomData.players.some(p => p.id === currentPlayer.value?.id)
      
      if (!isPlayerInRoom) {
        // 如果玩家不在房间中，跳转到房间列表
        router.push('/pvp')
        return
      }
      
    } catch (error) {
      // 如果获取房间失败，跳转到房间列表
      router.push('/pvp')
      return
    }
  }
  
  scrollChatToBottom()
})

onUnmounted(() => {
  // 清理工作在组件销毁时进行
})
</script>

<style scoped>
.room-lobby-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.header {
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

.room-info {
  text-align: center;
  flex: 1;
}

.room-info h1 {
  font-size: 2.5rem;
  margin-bottom: 5px;
  text-shadow: 2px 2px 4px rgba(0,0,0,0.3);
}

.room-id {
  font-size: 1rem;
  opacity: 0.8;
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

.content {
  max-width: 1200px;
  margin: 0 auto;
  display: grid;
  grid-template-columns: 1fr 400px;
  gap: 30px;
}

.left-panel, .right-panel {
  background: white;
  border-radius: 20px;
  padding: 30px;
  box-shadow: 0 10px 30px rgba(0,0,0,0.2);
}

.players-section h2, .chat-section h2 {
  margin-bottom: 20px;
  color: #333;
  font-size: 1.5rem;
}

.players-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.player-card {
  border: 2px solid #e1e5e9;
  border-radius: 15px;
  padding: 20px;
  transition: all 0.3s ease;
  background: #fafbfc;
}

.player-card.current-player {
  border-color: #667eea;
  background: #f0f4ff;
}

.player-card.ready {
  border-color: #28a745;
  background: #f8fff9;
}

.player-card.empty {
  border-style: dashed;
  border-color: #dee2e6;
  background: #f8f9fa;
}

.player-content {
  display: flex;
  align-items: center;
  gap: 15px;
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
  font-weight: 600;
  color: #333;
  margin-bottom: 5px;
}

.player-status {
  margin-bottom: 5px;
}

.status {
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 0.8rem;
  font-weight: 600;
}

.status.ready {
  background: #d4edda;
  color: #155724;
}

.status.not-ready {
  background: #fff3cd;
  color: #856404;
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

.empty-slot {
  text-align: center;
  color: #adb5bd;
}

.empty-avatar {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  border: 2px dashed #dee2e6;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  margin: 0 auto 10px;
}

.empty-text {
  font-size: 0.9rem;
}

.game-controls {
  border-top: 1px solid #e1e5e9;
  padding-top: 20px;
}

.ready-button, .start-button, .join-game-button, .view-result-button {
  width: 100%;
  padding: 15px;
  border: none;
  border-radius: 10px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  margin-bottom: 10px;
}

.ready-button {
  background: #6c757d;
  color: white;
}

.ready-button.ready {
  background: #28a745;
}

.start-button, .join-game-button {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.view-result-button {
  background: linear-gradient(135deg, #ffc107 0%, #fd7e14 100%);
  color: white;
}

.ready-button:hover:not(:disabled),
.start-button:hover:not(:disabled),
.join-game-button:hover:not(:disabled),
.view-result-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(0,0,0,0.2);
}

.ready-button:disabled,
.start-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.start-hint {
  text-align: center;
  color: #666;
  font-size: 0.9rem;
  padding: 10px;
  background: #f8f9fa;
  border-radius: 8px;
}

.game-status {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 15px;
  background: #f8f9fa;
  border-radius: 10px;
  margin-bottom: 15px;
  font-weight: 600;
  color: #333;
}

.status-icon {
  font-size: 1.2rem;
}

.chat-section {
  height: 600px;
  display: flex;
  flex-direction: column;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  border: 1px solid #e1e5e9;
  border-radius: 10px;
  padding: 15px;
  margin-bottom: 15px;
  background: #fafbfc;
}

.chat-message {
  margin-bottom: 15px;
  padding: 10px;
  border-radius: 10px;
  background: white;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
}

.chat-message.own-message {
  background: #e3f2fd;
  margin-left: 20px;
}

.message-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 5px;
}

.sender-name {
  font-weight: 600;
  color: #333;
  font-size: 0.9rem;
}

.message-time {
  font-size: 0.8rem;
  color: #666;
}

.message-content {
  color: #333;
  line-height: 1.4;
}

.empty-chat {
  text-align: center;
  color: #adb5bd;
  padding: 40px 20px;
}

.empty-icon {
  font-size: 3rem;
  margin-bottom: 15px;
}

.chat-input {
  display: flex;
  gap: 10px;
}

.chat-input input {
  flex: 1;
  padding: 12px;
  border: 2px solid #e1e5e9;
  border-radius: 25px;
  font-size: 1rem;
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
  padding: 12px 20px;
  border-radius: 25px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.send-button:hover:not(:disabled) {
  transform: scale(1.05);
  box-shadow: 0 5px 15px rgba(0,0,0,0.2);
}

.send-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
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

@media (max-width: 1024px) {
  .content {
    grid-template-columns: 1fr;
    gap: 20px;
  }
  
  .right-panel {
    order: -1;
  }
  
  .chat-section {
    height: 400px;
  }
}

@media (max-width: 768px) {
  .room-lobby-container {
    padding: 15px;
  }
  
  .header {
    flex-direction: column;
    gap: 15px;
    text-align: center;
  }
  
  .room-info h1 {
    font-size: 2rem;
  }
  
  .players-grid {
    grid-template-columns: 1fr;
  }
}
</style>