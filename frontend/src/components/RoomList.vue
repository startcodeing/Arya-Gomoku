<template>
  <div class="room-lobby-container">
    <div class="header">
      <h1>双人对战大厅</h1>
      <button class="refresh-button" @click="refresh" :disabled="isLoading">刷新</button>
    </div>

    <div class="content-grid">
      <!-- 创建房间卡片 -->
      <div class="card create-card">
        <h2>创建新房间</h2>
        <div class="form-grid">
          <input v-model.trim="newRoom.name" placeholder="房间名称" />
          <input v-model.trim="newRoom.playerName" placeholder="玩家昵称" />
          <select v-model.number="newRoom.maxPlayers">
            <option :value="2">2人</option>
          </select>
          <button class="primary-button" @click="createRoom" :disabled="isLoading || !canCreate">
            {{ isLoading ? '创建中...' : '创建房间' }}
          </button>
        </div>
        <p class="helper-text">创建后将自动进入房间大厅。</p>
      </div>

      <!-- 房间列表卡片 -->
      <div class="card rooms-card">
        <h2>可用房间</h2>
        <div v-if="rooms.length === 0" class="empty">暂无可用房间，创建一个房间开始对战吧!</div>
        <ul v-else class="rooms-list">
          <li v-for="room in rooms" :key="room.id" class="room-item">
            <div class="room-info">
              <span class="room-name">{{ room.name }}</span>
              <span class="room-capacity">{{ room.players.length }}/{{ room.maxPlayers }}</span>
            </div>
            <button class="secondary-button" @click="join(room.id)" :disabled="isLoading">加入</button>
          </li>
        </ul>
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
  try {
    await pvpStore.joinRoom(roomId, { playerName: newRoom.value.playerName.trim() })
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
  padding: 24px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  color: white;
}

.header h1 { font-size: 2rem; }

.refresh-button {
  background: rgba(255,255,255,0.2);
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 20px;
  cursor: pointer;
  transition: all 0.3s ease;
}
.refresh-button:hover { background: rgba(255,255,255,0.3); }
.refresh-button:disabled { opacity: 0.6; cursor: not-allowed; }

.content-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
}

.card {
  background: white;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 10px 24px rgba(0,0,0,0.1);
}

.create-card h2, .rooms-card h2 { color: #333; margin-bottom: 16px; }

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr auto auto;
  gap: 12px;
}

input, select {
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 8px;
  outline: none;
}
input:focus, select:focus { border-color: #667eea; }

.primary-button {
  background: #667eea;
  color: white;
  border: none;
  padding: 10px 16px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}
.primary-button:hover { background: #5a6fe0; }
.primary-button:disabled { opacity: 0.6; cursor: not-allowed; }

.helper-text { color: #666; margin-top: 8px; font-size: 0.9rem; }

.rooms-list { list-style: none; padding: 0; margin: 0; display: flex; flex-direction: column; gap: 12px; }
.room-item { display: flex; align-items: center; justify-content: space-between; padding: 12px; border: 1px solid #eee; border-radius: 12px; }
.room-info { display: flex; gap: 12px; align-items: center; }
.room-name { font-weight: 600; color: #333; }
.room-capacity { color: #666; }

.secondary-button {
  background: #f0f4ff;
  color: #333;
  border: none;
  padding: 8px 14px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}
.secondary-button:hover { background: #e4e8ff; }
.secondary-button:disabled { opacity: 0.6; cursor: not-allowed; }

.empty { color: #888; padding: 24px; text-align: center; }

/* Toast */
.error-toast, .success-toast {
  position: fixed;
  top: 20px;
  right: 20px;
  z-index: 1000;
}
.error-content, .success-content {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 14px;
  border-radius: 10px;
  box-shadow: 0 6px 16px rgba(0,0,0,0.15);
}
.error-content { background: #ffe6e6; color: #cc0000; }
.success-content { background: #e6ffed; color: #0f8a2c; }
.error-close, .success-close { background: transparent; border: none; color: inherit; font-size: 18px; cursor: pointer; }
</style>