<template>
  <div class="room-lobby-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-text">
          <h1>双人对战大厅</h1>
          <p>创建房间或加入现有房间，与朋友一起享受五子棋对战的乐趣</p>
        </div>
        <button class="refresh-button" @click="refresh" :disabled="isLoading">
          <span class="refresh-icon">🔄</span>
          {{ isLoading ? '刷新中...' : '刷新房间' }}
        </button>
      </div>
    </div>

    <!-- 主要内容区域 -->
    <div class="main-content">
      <!-- 创建房间区域 -->
      <div class="create-room-section">
        <div class="section-header">
          <div class="section-icon">🎮</div>
          <h2>创建新房间</h2>
        </div>
        <div class="create-card">
          <div class="create-form">
            <div class="form-row">
              <div class="form-group">
                <label>房间名称</label>
                <input
                  v-model.trim="newRoom.name"
                  placeholder="给你的房间起个名字"
                  class="form-input"
                  maxlength="20"
                />
                <span class="input-hint">{{ newRoom.name.length }}/20</span>
              </div>
            </div>
            <div class="form-row">
              <div class="form-group">
                <label>玩家昵称</label>
                <input
                  v-model.trim="newRoom.playerName"
                  placeholder="输入你的昵称"
                  class="form-input"
                  maxlength="12"
                />
                <span class="input-hint">{{ newRoom.playerName.length }}/12</span>
              </div>
            </div>
            <div class="form-row">
              <div class="form-group">
                <label>房间人数</label>
                <select v-model.number="newRoom.maxPlayers" class="form-select">
                  <option :value="2">2人对战</option>
                </select>
              </div>
            </div>
            <button
              class="create-button"
              @click="createRoom"
              :disabled="isLoading || !canCreate"
            >
              <span v-if="isLoading" class="loading-spinner"></span>
              {{ isLoading ? '创建中...' : '创建房间' }}
            </button>
          </div>
          <div class="create-tips">
            <div class="tip-item">
              <span class="tip-icon">💡</span>
              <span>创建成功后将自动进入房间大厅</span>
            </div>
            <div class="tip-item">
              <span class="tip-icon">🔗</span>
              <span>可以分享邀请链接邀请朋友加入</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 房间列表区域 -->
      <div class="rooms-section">
        <div class="section-header">
          <div class="section-icon">🏠</div>
          <h2>可用房间</h2>
          <div class="room-count">
            共 {{ rooms.length }} 个房间
          </div>
        </div>

        <!-- 加入房间表单 -->
        <div class="join-card">
          <div class="join-form">
            <div class="form-group">
              <label>玩家昵称</label>
              <input
                v-model.trim="newRoom.playerName"
                placeholder="加入房间前请先输入你的昵称"
                class="form-input"
                maxlength="12"
              />
              <span class="input-hint">{{ newRoom.playerName.length }}/12</span>
            </div>
          </div>
        </div>

        <!-- 房间列表 -->
        <div class="rooms-container">
          <div v-if="rooms.length === 0" class="empty-state">
            <div class="empty-icon">🎯</div>
            <h3>暂无可用房间</h3>
            <p>创建一个新房间，邀请朋友开始对战吧！</p>
          </div>

          <div v-else class="rooms-grid">
            <div
              v-for="room in rooms"
              :key="room.id"
              class="room-card"
              :class="{ 'room-full': room.players.length >= room.maxPlayers }"
            >
              <div class="room-header">
                <div class="room-name">{{ room.name }}</div>
                <div class="room-status" :class="{ 'status-full': room.players.length >= room.maxPlayers }">
                  <span class="status-dot"></span>
                  {{ room.players.length >= room.maxPlayers ? '已满' : '等待中' }}
                </div>
              </div>

              <div class="room-body">
                <div class="room-info">
                  <div class="info-item">
                    <span class="info-label">玩家数量</span>
                    <span class="info-value">{{ room.players.length }}/{{ room.maxPlayers }}</span>
                  </div>
                  <div class="player-avatars">
                    <div
                      v-for="player in room.players.slice(0, 4)"
                      :key="player.id"
                      class="player-avatar"
                      :title="player.name"
                    >
                      {{ player.name.charAt(0).toUpperCase() }}
                    </div>
                    <div v-if="room.players.length < room.maxPlayers" class="empty-avatar" title="等待玩家">
                      +
                    </div>
                  </div>
                </div>
              </div>

              <div class="room-footer">
                <button
                  class="join-button"
                  @click="join(room.id)"
                  :disabled="isLoading || !newRoom.playerName.trim() || room.players.length >= room.maxPlayers"
                >
                  <span v-if="room.players.length >= room.maxPlayers">房间已满</span>
                  <span v-else-if="!newRoom.playerName.trim()">请输入昵称</span>
                  <span v-else>加入房间</span>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 错误提示 -->
    <div v-if="error" class="error-toast" @click="clearError">
      <div class="error-content">
        <span class="error-icon">⚠️</span>
        <span class="error-text">{{ error }}</span>
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
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { usePvpStore } from '../stores/pvp'
import type { Room } from '../types/pvp'

const router = useRouter()
const pvpStore = usePvpStore()

// 使用 store 的响应式数据
const rooms = computed<Room[]>(() => pvpStore.rooms)
const isLoading = computed(() => pvpStore.isLoading)
const error = computed(() => pvpStore.error)

const successMessage = ref('')

const newRoom = ref({
  name: '',
  playerName: '',
  maxPlayers: 2,
})

const canCreate = computed(() => !!newRoom.value.name.trim() && !!newRoom.value.playerName.trim())

function showSuccess(message: string) {
  successMessage.value = message
  setTimeout(() => { successMessage.value = '' }, 2500)
}

function clearError() {
  pvpStore.error = null
}

function clearSuccess() {
  successMessage.value = ''
}

async function refresh() {
  await pvpStore.fetchRooms()
}

async function createRoom() {
  if (!canCreate.value) return
  const payload = {
    roomName: newRoom.value.name.trim(),
    playerName: newRoom.value.playerName.trim(),
    maxPlayers: Number(newRoom.value.maxPlayers),
  }
  try {
    await pvpStore.createRoom(payload)
    showSuccess('房间创建成功，正在进入...')
    if (pvpStore.currentRoom) {
      router.push(`/room/${pvpStore.currentRoom.id}`)
    }
  } catch (e) {
    // 错误由 store 设置到 error
  }
}

async function join(roomId: string) {
  // 验证玩家名称
  const playerName = newRoom.value.playerName.trim()
  if (!playerName) {
    pvpStore.error = '请输入玩家名称'
    return
  }
  
  try {
    await pvpStore.joinRoom(roomId, { playerName })
    showSuccess('加入成功，正在进入房间...')
    if (pvpStore.currentRoom) {
      router.push(`/room/${pvpStore.currentRoom.id}`)
    }
  } catch (e) {
    // 错误由 store 设置到 error
  }
}

onMounted(() => {
  refresh()
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
  padding: 30px 0;
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 30px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-text h1 {
  color: white;
  font-size: 2.5rem;
  margin: 0 0 10px 0;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.3);
}

.header-text p {
  color: rgba(255, 255, 255, 0.9);
  font-size: 1.1rem;
  margin: 0;
  max-width: 600px;
}

.refresh-button {
  background: rgba(255, 255, 255, 0.2);
  color: white;
  border: 1px solid rgba(255, 255, 255, 0.3);
  padding: 12px 24px;
  border-radius: 30px;
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 1rem;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 8px;
  backdrop-filter: blur(10px);
}

.refresh-button:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.3);
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.2);
}

.refresh-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.refresh-icon {
  font-size: 1.1rem;
}

/* 主要内容区域 */
.main-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 40px 30px;
  display: grid;
  grid-template-columns: 420px 1fr;
  gap: 40px;
  align-items: start;
}

/* 区域标题 */
.section-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 25px;
}

.section-icon {
  font-size: 2rem;
  background: rgba(255, 255, 255, 0.9);
  width: 50px;
  height: 50px;
  border-radius: 15px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
}

.section-header h2 {
  color: white;
  font-size: 1.5rem;
  margin: 0;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.3);
}

/* 创建房间区域 */
.create-room-section {
  position: sticky;
  top: 30px;
}

.create-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 20px;
  padding: 30px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.create-form {
  margin-bottom: 20px;
}

.form-row {
  margin-bottom: 20px;
}

.form-group {
  position: relative;
}

.form-group label {
  display: block;
  color: #374151;
  font-weight: 600;
  margin-bottom: 8px;
  font-size: 0.95rem;
}

.form-input, .form-select {
  width: 100%;
  padding: 14px 16px;
  border: 2px solid #e5e7eb;
  border-radius: 12px;
  font-size: 1rem;
  outline: none;
  transition: all 0.3s ease;
  background: white;
}

.form-input:focus, .form-select:focus {
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.form-input::placeholder {
  color: #9ca3af;
}

.input-hint {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  color: #9ca3af;
  font-size: 0.85rem;
  background: white;
  padding: 2px 6px;
  border-radius: 4px;
}

.create-button {
  width: 100%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 16px;
  border-radius: 12px;
  font-size: 1.1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
}

.create-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(102, 126, 234, 0.4);
}

.create-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.loading-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top: 2px solid white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.create-tips {
  padding-top: 20px;
  border-top: 1px solid #e5e7eb;
}

.tip-item {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
  color: #6b7280;
  font-size: 0.9rem;
}

.tip-icon {
  font-size: 1.1rem;
}

/* 房间列表区域 */
.rooms-section {
  flex: 1;
}

.room-count {
  background: rgba(255, 255, 255, 0.2);
  color: white;
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 0.9rem;
  font-weight: 500;
  backdrop-filter: blur(10px);
}

.join-card {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  border-radius: 15px;
  padding: 20px;
  margin-bottom: 25px;
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.join-form .form-group {
  margin-bottom: 0;
}

/* 房间网格 */
.rooms-container {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 20px;
  padding: 25px;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: white;
}

.empty-icon {
  font-size: 4rem;
  margin-bottom: 20px;
  opacity: 0.8;
}

.empty-state h3 {
  font-size: 1.5rem;
  margin: 0 0 10px 0;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.3);
}

.empty-state p {
  font-size: 1.1rem;
  margin: 0;
  opacity: 0.9;
}

.rooms-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.room-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 16px;
  padding: 20px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
}

.room-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 30px rgba(0, 0, 0, 0.15);
  border-color: rgba(255, 255, 255, 0.5);
}

.room-card.room-full {
  opacity: 0.7;
  border-color: rgba(239, 68, 68, 0.3);
}

.room-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 15px;
}

.room-name {
  font-size: 1.2rem;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 5px;
  flex: 1;
  margin-right: 10px;
}

.room-status {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  border-radius: 20px;
  font-size: 0.8rem;
  font-weight: 600;
  white-space: nowrap;
}

.room-status:not(.status-full) {
  background: #dcfce7;
  color: #16a34a;
}

.room-status.status-full {
  background: #fee2e2;
  color: #dc2626;
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: currentColor;
}

.room-body {
  margin-bottom: 15px;
}

.room-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.info-label {
  font-size: 0.85rem;
  color: #6b7280;
  font-weight: 500;
}

.info-value {
  font-size: 1.1rem;
  font-weight: 600;
  color: #1f2937;
}

.player-avatars {
  display: flex;
  gap: 6px;
}

.player-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  font-weight: bold;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.9rem;
  border: 2px solid white;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.empty-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  border: 2px dashed #d1d5db;
  background: #f9fafb;
  color: #9ca3af;
  font-weight: bold;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.9rem;
}

.room-footer {
  margin-top: auto;
}

.join-button {
  width: 100%;
  padding: 12px;
  border: none;
  border-radius: 10px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
}

.join-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4);
}

.join-button:disabled {
  background: #e5e7eb;
  color: #9ca3af;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

/* Toast 样式 */
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
    grid-template-columns: 380px 1fr;
    gap: 30px;
  }

  .rooms-grid {
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  }
}

@media (max-width: 968px) {
  .main-content {
    grid-template-columns: 1fr;
    gap: 30px;
  }

  .create-room-section {
    position: static;
  }

  .header-content {
    flex-direction: column;
    gap: 20px;
    text-align: center;
  }

  .header-text p {
    max-width: none;
  }
}

@media (max-width: 640px) {
  .page-header {
    padding: 20px 0;
  }

  .header-content {
    padding: 0 20px;
  }

  .header-text h1 {
    font-size: 2rem;
  }

  .main-content {
    padding: 30px 20px;
  }

  .rooms-grid {
    grid-template-columns: 1fr;
  }

  .refresh-button {
    padding: 10px 20px;
    font-size: 0.9rem;
  }

  .create-card, .join-card {
    padding: 20px;
  }

  .rooms-container {
    padding: 20px;
  }
}
</style>