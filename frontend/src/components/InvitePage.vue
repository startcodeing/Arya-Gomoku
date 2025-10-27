<template>
  <div class="invite-page-container">
    <div class="invite-card">
      <div class="invite-header">
        <h1>🎮 五子棋邀请</h1>
        <p class="invite-subtitle">你收到了一个游戏邀请</p>
      </div>

      <div v-if="isLoading" class="loading-section">
        <div class="loading-spinner"></div>
        <p>正在加载房间信息...</p>
      </div>

      <div v-else-if="error" class="error-section">
        <div class="error-icon">❌</div>
        <h3>邀请无效</h3>
        <p>{{ error }}</p>
        <button @click="goToLobby" class="primary-button">
          返回大厅
        </button>
      </div>

      <div v-else-if="roomInfo" class="room-info-section">
        <div class="room-details">
          <h2>{{ roomInfo.name }}</h2>
          <div class="room-meta">
            <div class="meta-item">
              <span class="meta-label">房间ID:</span>
              <span class="meta-value">{{ roomInfo.id }}</span>
            </div>
            <div class="meta-item">
              <span class="meta-label">玩家数量:</span>
              <span class="meta-value">{{ roomInfo.players.length }}/{{ roomInfo.maxPlayers }}</span>
            </div>
            <div class="meta-item">
              <span class="meta-label">房间状态:</span>
              <span class="meta-value status" :class="roomInfo.status">
                {{ getStatusText(roomInfo.status) }}
              </span>
            </div>
          </div>

          <div v-if="roomInfo.players.length > 0" class="players-preview">
            <h3>当前玩家</h3>
            <div class="players-list">
              <div 
                v-for="player in roomInfo.players" 
                :key="player.id"
                class="player-item"
              >
                <div class="player-avatar">
                  {{ player.name.charAt(0).toUpperCase() }}
                </div>
                <span class="player-name">{{ player.name }}</span>
              </div>
            </div>
          </div>
        </div>

        <div class="join-section">
          <div class="join-form">
            <input 
              v-model.trim="playerName" 
              placeholder="输入你的昵称" 
              class="name-input"
              @keyup.enter="joinRoom"
              :disabled="isJoining"
            />
            <button 
              @click="joinRoom" 
              :disabled="!canJoin || isJoining"
              class="join-button"
            >
              {{ isJoining ? '加入中...' : '加入游戏' }}
            </button>
          </div>
          <p class="join-hint">
            {{ roomInfo.status === 'waiting' ? '房间正在等待玩家，快来加入吧！' : 
               roomInfo.status === 'playing' ? '游戏正在进行中，你可以观战' : 
               '游戏已结束' }}
          </p>
        </div>
      </div>

      <div class="footer-actions">
        <button @click="goToLobby" class="secondary-button">
          浏览其他房间
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { usePvpStore } from '../stores/pvp'
import { pvpApi } from '../services/pvpApi'
import type { Room } from '../types/pvp'

const router = useRouter()
const route = useRoute()
const pvpStore = usePvpStore()

// 响应式数据
const roomInfo = ref<Room | null>(null)
const playerName = ref('')
const isLoading = ref(true)
const isJoining = ref(false)
const error = ref<string | null>(null)

// 计算属性
const roomId = computed(() => route.params.id as string)
const canJoin = computed(() => {
  return playerName.value.trim().length > 0 && 
         roomInfo.value?.status === 'waiting' &&
         roomInfo.value.players.length < roomInfo.value.maxPlayers
})

// 方法
async function loadRoomInfo() {
  if (!roomId.value) {
    error.value = '无效的房间ID'
    isLoading.value = false
    return
  }

  try {
    isLoading.value = true
    error.value = null
    
    // 获取房间信息
    const response = await pvpApi.getRoomInfo(roomId.value)
    roomInfo.value = response.room
    
    if (!roomInfo.value) {
      error.value = '房间不存在或已关闭'
    }
  } catch (err: any) {
    console.error('加载房间信息失败:', err)
    error.value = err?.message || '无法加载房间信息，请检查链接是否正确'
  } finally {
    isLoading.value = false
  }
}

async function joinRoom() {
  if (!canJoin.value || !roomInfo.value) return

  try {
    isJoining.value = true
    
    await pvpStore.joinRoom(roomInfo.value.id, { 
      playerName: playerName.value.trim() 
    })
    
    // 加入成功，跳转到房间大厅
    router.push(`/room/${roomInfo.value.id}`)
  } catch (err: any) {
    console.error('加入房间失败:', err)
    error.value = err?.message || '加入房间失败'
  } finally {
    isJoining.value = false
  }
}

function goToLobby() {
  router.push('/pvp')
}

function getStatusText(status: string): string {
  switch (status) {
    case 'waiting':
      return '等待中'
    case 'playing':
      return '游戏中'
    case 'finished':
      return '已结束'
    default:
      return '未知'
  }
}

// 生命周期
onMounted(() => {
  loadRoomInfo()
})
</script>

<style scoped>
.invite-page-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.invite-card {
  background: white;
  border-radius: 20px;
  padding: 40px;
  max-width: 500px;
  width: 100%;
  box-shadow: 0 20px 40px rgba(0,0,0,0.1);
  text-align: center;
}

.invite-header h1 {
  font-size: 2.5rem;
  margin-bottom: 10px;
  color: #333;
}

.invite-subtitle {
  color: #666;
  font-size: 1.1rem;
  margin-bottom: 30px;
}

.loading-section {
  padding: 40px 0;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #f3f3f3;
  border-top: 4px solid #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 20px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.error-section {
  padding: 40px 0;
}

.error-icon {
  font-size: 3rem;
  margin-bottom: 20px;
}

.error-section h3 {
  color: #dc3545;
  margin-bottom: 15px;
}

.error-section p {
  color: #666;
  margin-bottom: 25px;
}

.room-info-section {
  text-align: left;
}

.room-details h2 {
  color: #333;
  margin-bottom: 20px;
  text-align: center;
  font-size: 1.8rem;
}

.room-meta {
  background: #f8f9fa;
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 25px;
}

.meta-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 10px;
}

.meta-item:last-child {
  margin-bottom: 0;
}

.meta-label {
  color: #666;
  font-weight: 500;
}

.meta-value {
  color: #333;
  font-weight: 600;
}

.meta-value.status.waiting {
  color: #28a745;
}

.meta-value.status.playing {
  color: #ffc107;
}

.meta-value.status.finished {
  color: #dc3545;
}

.players-preview h3 {
  color: #333;
  margin-bottom: 15px;
  font-size: 1.2rem;
}

.players-list {
  display: flex;
  gap: 15px;
  flex-wrap: wrap;
}

.player-item {
  display: flex;
  align-items: center;
  gap: 10px;
  background: #f8f9fa;
  padding: 10px 15px;
  border-radius: 25px;
}

.player-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: #667eea;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  font-size: 0.9rem;
}

.player-name {
  color: #333;
  font-weight: 500;
}

.join-section {
  margin-top: 30px;
  text-align: center;
}

.join-form {
  display: flex;
  gap: 15px;
  margin-bottom: 15px;
}

.name-input {
  flex: 1;
  padding: 12px 16px;
  border: 2px solid #e9ecef;
  border-radius: 10px;
  font-size: 1rem;
  outline: none;
  transition: border-color 0.3s ease;
}

.name-input:focus {
  border-color: #667eea;
}

.join-button {
  background: #667eea;
  color: white;
  border: none;
  padding: 12px 24px;
  border-radius: 10px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  white-space: nowrap;
}

.join-button:hover:not(:disabled) {
  background: #5a6fe0;
  transform: translateY(-2px);
}

.join-button:disabled {
  background: #ccc;
  cursor: not-allowed;
  transform: none;
}

.join-hint {
  color: #666;
  font-size: 0.9rem;
  margin: 0;
}

.footer-actions {
  margin-top: 30px;
  padding-top: 20px;
  border-top: 1px solid #e9ecef;
}

.primary-button, .secondary-button {
  padding: 12px 24px;
  border-radius: 10px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  border: none;
}

.primary-button {
  background: #667eea;
  color: white;
}

.primary-button:hover {
  background: #5a6fe0;
  transform: translateY(-2px);
}

.secondary-button {
  background: #f8f9fa;
  color: #666;
  border: 2px solid #e9ecef;
}

.secondary-button:hover {
  background: #e9ecef;
  transform: translateY(-2px);
}

@media (max-width: 600px) {
  .invite-card {
    padding: 30px 20px;
    margin: 10px;
  }
  
  .join-form {
    flex-direction: column;
  }
  
  .players-list {
    justify-content: center;
  }
}
</style>