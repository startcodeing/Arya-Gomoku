<template>
  <div class="user-profile">
    <!-- 页面头部 -->
    <div class="profile-header">
      <div class="header-content">
        <div class="header-title">
          <h1>用户中心</h1>
          <span class="welcome-text">欢迎回来，{{ userStore.user?.username || '玩家' }}！</span>
        </div>
        <button @click="logout" class="logout-btn">
          退出登录
        </button>
      </div>
    </div>

    <!-- 主要内容 -->
    <div class="profile-content">
      <!-- 第一行：用户信息和游戏统计 -->
      <div class="content-row">
        <!-- 用户信息卡片 -->
        <div class="profile-card">
          <div class="card-header">
            <h2>个人信息</h2>
            <div class="header-icon">👤</div>
          </div>
          <div class="user-info">
            <div class="info-item">
              <span class="info-label">用户名:</span>
              <span class="info-value">{{ userStore.user?.username }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">昵称:</span>
              <span v-if="!editingProfile" class="info-value">{{ userStore.user?.nickname || '未设置' }}</span>
              <input
                v-else
                v-model="profileForm.nickname"
                type="text"
                placeholder="请输入昵称"
                class="edit-input"
              >
            </div>
            <div class="info-item">
              <span class="info-label">邮箱:</span>
              <span v-if="!editingProfile" class="info-value">{{ userStore.user?.email || '未设置' }}</span>
              <input
                v-else
                v-model="profileForm.email"
                type="email"
                placeholder="请输入邮箱"
                class="edit-input"
              >
            </div>
            <div class="info-item">
              <span class="info-label">注册时间:</span>
              <span class="info-value">{{ formatDate(userStore.user?.created_at) }}</span>
            </div>
          </div>

          <div class="profile-actions">
            <button
              v-if="!editingProfile"
              @click="startEditProfile"
              class="edit-btn"
            >
              <span class="btn-icon">✏️</span>
              编辑资料
            </button>
            <template v-else>
              <button @click="saveProfile" class="save-btn" :disabled="profileLoading">
                <span class="btn-icon">💾</span>
                {{ profileLoading ? '保存中...' : '保存' }}
              </button>
              <button @click="cancelEditProfile" class="cancel-btn">
                <span class="btn-icon">❌</span>
                取消
              </button>
            </template>
          </div>
        </div>

        <!-- 游戏统计卡片 -->
        <div class="stats-card">
          <div class="card-header">
            <h2>游戏统计</h2>
            <div class="header-icon">📊</div>
          </div>
          <div v-if="statsLoading" class="loading">
            <div class="loading-spinner"></div>
            <span>加载中...</span>
          </div>
          <div v-else-if="gameStats" class="stats-grid">
            <div class="stat-item primary">
              <div class="stat-icon">🎮</div>
              <div class="stat-content">
                <div class="stat-value">{{ gameStats.total_games || 0 }}</div>
                <div class="stat-label">总游戏数</div>
              </div>
            </div>
            <div class="stat-item success">
              <div class="stat-icon">🏆</div>
              <div class="stat-content">
                <div class="stat-value">{{ gameStats.wins || 0 }}</div>
                <div class="stat-label">胜利次数</div>
              </div>
            </div>
            <div class="stat-item danger">
              <div class="stat-icon">💔</div>
              <div class="stat-content">
                <div class="stat-value">{{ gameStats.losses || 0 }}</div>
                <div class="stat-label">失败次数</div>
              </div>
            </div>
            <div class="stat-item info">
              <div class="stat-icon">📈</div>
              <div class="stat-content">
                <div class="stat-value">{{ winRate }}%</div>
                <div class="stat-label">胜率</div>
              </div>
            </div>
            <div class="stat-item warning">
              <div class="stat-icon">🎯</div>
              <div class="stat-content">
                <div class="stat-value">{{ gameStats.avg_moves || 0 }}</div>
                <div class="stat-label">平均步数</div>
              </div>
            </div>
            <div class="stat-item secondary">
              <div class="stat-icon">⏱️</div>
              <div class="stat-content">
                <div class="stat-value">{{ formatDuration(gameStats.avg_duration) }}</div>
                <div class="stat-label">平均时长</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 第二行：游戏历史 -->
      <div class="content-row full-width">
        <!-- 游戏历史卡片 -->
        <div class="history-card">
          <div class="card-header">
            <h2>游戏历史</h2>
            <div class="header-icon">📋</div>
          </div>
          <div v-if="historyLoading" class="loading">
            <div class="loading-spinner"></div>
            <span>加载中...</span>
          </div>
          <div v-else-if="gameHistory.length > 0" class="history-list">
            <div
              v-for="game in gameHistory"
              :key="game.id"
              class="history-item"
            >
              <div class="game-info">
                <div class="game-type">{{ getGameTypeText(game.game_type) }}</div>
                <div class="game-result" :class="game.result">
                  {{ getResultText(game.result) }}
                </div>
              </div>
              <div class="game-details">
                <span class="detail-item">🎯 步数: {{ game.moves }}</span>
                <span class="detail-item">⏱️ 时长: {{ formatDuration(game.duration) }}</span>
                <span v-if="game.difficulty" class="detail-item">📊 难度: {{ getDifficultyText(game.difficulty) }}</span>
              </div>
              <div class="game-date">
                📅 {{ formatDate(game.created_at) }}
              </div>
            </div>

            <!-- 分页 -->
            <div class="pagination">
              <button
                @click="loadPreviousPage"
                :disabled="currentPage <= 1"
                class="page-btn"
              >
                <span class="btn-icon">⬅️</span>
                上一页
              </button>
              <span class="page-info">第 {{ currentPage }} 页</span>
              <button
                @click="loadNextPage"
                :disabled="gameHistory.length < pageSize"
                class="page-btn"
              >
                <span class="btn-icon">➡️</span>
                下一页
              </button>
            </div>
          </div>
          <div v-else class="no-data">
            <div class="no-data-icon">🎮</div>
            <div>暂无游戏记录</div>
            <div class="no-data-hint">开始你的第一局游戏吧！</div>
          </div>
        </div>
      </div>

      <!-- 第三行：密码修改 -->
      <div class="content-row">
        <!-- 密码修改卡片 -->
        <div class="password-card">
          <div class="card-header">
            <h2>修改密码</h2>
            <div class="header-icon">🔐</div>
          </div>
          <form @submit.prevent="changePassword" class="password-form">
            <div class="form-group">
              <label class="form-label">
                <span class="label-icon">🔑</span>
                当前密码:
              </label>
              <input
                v-model="passwordForm.oldPassword"
                type="password"
                required
                class="form-input"
                placeholder="请输入当前密码"
              >
            </div>
            <div class="form-group">
              <label class="form-label">
                <span class="label-icon">🔒</span>
                新密码:
              </label>
              <input
                v-model="passwordForm.newPassword"
                type="password"
                required
                minlength="8"
                class="form-input"
                placeholder="请输入新密码（至少8位）"
              >
            </div>
            <div class="form-group">
              <label class="form-label">
                <span class="label-icon">🔓</span>
                确认新密码:
              </label>
              <input
                v-model="passwordForm.confirmPassword"
                type="password"
                required
                class="form-input"
                placeholder="请再次输入新密码"
              >
            </div>
            <div class="password-strength">
              <div class="strength-info">
                <span class="strength-text">密码要求：</span>
                <span class="requirement">至少8个字符</span>
              </div>
            </div>
            <button
              type="submit"
              class="change-password-btn"
              :disabled="passwordLoading"
            >
              <span class="btn-icon">🔄</span>
              {{ passwordLoading ? '修改中...' : '修改密码' }}
            </button>
          </form>
        </div>
      </div>
    </div>

    <!-- 错误提示 -->
    <transition name="slide-in">
      <div v-if="error" class="error-message">
        <span class="message-icon">⚠️</span>
        {{ error }}
        <button @click="error = null" class="close-btn">×</button>
      </div>
    </transition>

    <!-- 成功提示 -->
    <transition name="slide-in">
      <div v-if="successMessage" class="success-message">
        <span class="message-icon">✅</span>
        {{ successMessage }}
        <button @click="successMessage = null" class="close-btn">×</button>
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { userApi } from '../services/api'

const router = useRouter()
const userStore = useUserStore()

// 响应式数据
const editingProfile = ref(false)
const profileLoading = ref(false)
const statsLoading = ref(false)
const historyLoading = ref(false)
const passwordLoading = ref(false)
const error = ref<string | null>(null)
const successMessage = ref<string | null>(null)

const gameStats = ref<any>(null)
const gameHistory = ref<any[]>([])
const currentPage = ref(1)
const pageSize = ref(10)

// 表单数据
const profileForm = ref({
  nickname: '',
  email: ''
})

const passwordForm = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 计算属性
const winRate = computed(() => {
  if (!gameStats.value || !gameStats.value.total_games) return 0
  return Math.round((gameStats.value.wins / gameStats.value.total_games) * 100)
})

// 方法
const startEditProfile = () => {
  profileForm.value.nickname = userStore.user?.nickname || ''
  profileForm.value.email = userStore.user?.email || ''
  editingProfile.value = true
}

const cancelEditProfile = () => {
  editingProfile.value = false
  profileForm.value.nickname = ''
  profileForm.value.email = ''
}

const saveProfile = async () => {
  try {
    profileLoading.value = true
    error.value = null
    
    await userStore.updateProfile({
      nickname: profileForm.value.nickname,
      email: profileForm.value.email
    })
    
    editingProfile.value = false
    successMessage.value = '资料更新成功'
    setTimeout(() => {
      successMessage.value = null
    }, 3000)
  } catch (err: any) {
    error.value = err.message || '更新资料失败'
  } finally {
    profileLoading.value = false
  }
}

const changePassword = async () => {
  try {
    error.value = null
    
    if (passwordForm.value.newPassword !== passwordForm.value.confirmPassword) {
      error.value = '新密码和确认密码不匹配'
      return
    }
    
    if (passwordForm.value.newPassword.length < 8) {
      error.value = '新密码长度至少8个字符'
      return
    }
    
    passwordLoading.value = true
    
    await userStore.changePassword(
      passwordForm.value.oldPassword,
      passwordForm.value.newPassword
    )
    
    passwordForm.value.oldPassword = ''
    passwordForm.value.newPassword = ''
    passwordForm.value.confirmPassword = ''
    
    successMessage.value = '密码修改成功'
    setTimeout(() => {
      successMessage.value = null
    }, 3000)
  } catch (err: any) {
    error.value = err.message || '修改密码失败'
  } finally {
    passwordLoading.value = false
  }
}

const loadGameStats = async () => {
  try {
    statsLoading.value = true
    const response = await userApi.getGameStats()
    gameStats.value = response.data
  } catch (err: any) {
    console.error('加载游戏统计失败:', err)
  } finally {
    statsLoading.value = false
  }
}

const loadGameHistory = async (page: number = 1) => {
  try {
    historyLoading.value = true
    const response = await userApi.getGameHistory(page, pageSize.value)
    gameHistory.value = response.data || []
    currentPage.value = page
  } catch (err: any) {
    console.error('加载游戏历史失败:', err)
  } finally {
    historyLoading.value = false
  }
}

const loadPreviousPage = () => {
  if (currentPage.value > 1) {
    loadGameHistory(currentPage.value - 1)
  }
}

const loadNextPage = () => {
  loadGameHistory(currentPage.value + 1)
}

const logout = async () => {
  try {
    await userStore.logout()
    router.push('/login')
  } catch (err: any) {
    console.error('退出登录失败:', err)
  }
}

// 工具函数
const formatDate = (dateString: string | undefined) => {
  if (!dateString) return '未知'
  return new Date(dateString).toLocaleString('zh-CN')
}

const formatDuration = (seconds: number | undefined) => {
  if (!seconds) return '0秒'
  const minutes = Math.floor(seconds / 60)
  const remainingSeconds = seconds % 60
  return minutes > 0 ? `${minutes}分${remainingSeconds}秒` : `${remainingSeconds}秒`
}

const getGameTypeText = (type: string) => {
  const types: Record<string, string> = {
    'ai': 'AI对战',
    'pvp': '玩家对战'
  }
  return types[type] || type
}

const getResultText = (result: string) => {
  const results: Record<string, string> = {
    'win': '胜利',
    'loss': '失败',
    'draw': '平局'
  }
  return results[result] || result
}

const getDifficultyText = (difficulty: string) => {
  const difficulties: Record<string, string> = {
    'easy': '简单',
    'medium': '中等',
    'hard': '困难',
    'expert': '专家'
  }
  return difficulties[difficulty] || difficulty
}

// 生命周期
onMounted(() => {
  loadGameStats()
  loadGameHistory()
})
</script>

<style scoped>
.user-profile {
  max-width: 1400px;
  margin: 0 auto;
  padding: 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  min-height: 100vh;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
}

/* 页面头部 */
.profile-header {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 20px;
  padding: 32px;
  margin-bottom: 28px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.2);
  position: relative;
  overflow: hidden;
}

.profile-header::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, #667eea 0%, #764ba2 50%, #f093fb 100%);
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 24px;
}

.header-title {
  flex: 1;
}

.header-title h1 {
  margin: 0 0 6px 0;
  color: #1a202c;
  font-size: 2.2rem;
  font-weight: 800;
  letter-spacing: -0.02em;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.welcome-text {
  color: #64748b;
  font-size: 1.1rem;
  font-weight: 500;
  letter-spacing: 0.01em;
}

.logout-btn {
  background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
  color: white;
  border: none;
  padding: 14px 28px;
  border-radius: 12px;
  cursor: pointer;
  font-weight: 600;
  font-size: 1rem;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 8px;
  box-shadow: 0 4px 16px rgba(239, 68, 68, 0.3);
  position: relative;
  overflow: hidden;
}

.logout-btn::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  transition: left 0.5s;
}

.logout-btn:hover::before {
  left: 100%;
}

.logout-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(239, 68, 68, 0.4);
}

/* 主要内容区域 */
.profile-content {
  display: flex;
  flex-direction: column;
  gap: 28px;
}

.content-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 28px;
}

.content-row.full-width {
  grid-template-columns: 1fr;
}

/* 卡片样式 */
.profile-card,
.stats-card,
.history-card,
.password-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 20px;
  padding: 28px;
  box-shadow: 0 16px 48px rgba(0, 0, 0, 0.12);
  border: 1px solid rgba(255, 255, 255, 0.2);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
}

.profile-card::before,
.stats-card::before,
.history-card::before,
.password-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
  transform: scaleX(0);
  transition: transform 0.3s ease;
}

.profile-card:hover,
.stats-card:hover,
.history-card:hover,
.password-card:hover {
  transform: translateY(-6px) scale(1.02);
  box-shadow: 0 24px 72px rgba(0, 0, 0, 0.18);
}

.profile-card:hover::before,
.stats-card:hover::before,
.history-card:hover::before,
.password-card:hover::before {
  transform: scaleX(1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 28px;
  padding-bottom: 20px;
  border-bottom: 2px solid rgba(226, 232, 240, 0.6);
}

.card-header h2 {
  margin: 0;
  color: #1a202c;
  font-size: 1.6rem;
  font-weight: 700;
  letter-spacing: -0.01em;
}

.header-icon {
  font-size: 1.6rem;
  opacity: 0.9;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
}

/* 用户信息样式 */
.user-info {
  margin-bottom: 24px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid #e2e8f0;
}

.info-item:last-child {
  border-bottom: none;
}

.info-label {
  font-weight: 600;
  color: #4a5568;
  font-size: 0.95rem;
}

.info-value {
  color: #2d3748;
  font-weight: 500;
}

.edit-input {
  padding: 8px 12px;
  border: 2px solid #e2e8f0;
  border-radius: 8px;
  width: 100%;
  font-size: 0.95rem;
  transition: border-color 0.3s ease;
}

.edit-input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

/* 按钮样式 */
.profile-actions {
  display: flex;
  gap: 12px;
}

.edit-btn,
.save-btn,
.cancel-btn,
.change-password-btn {
  padding: 12px 20px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 600;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 8px;
}

.btn-icon {
  font-size: 1.1rem;
}

.edit-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.edit-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(102, 126, 234, 0.4);
}

.save-btn {
  background: linear-gradient(135deg, #48bb78 0%, #38a169 100%);
  color: white;
}

.save-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(72, 187, 120, 0.4);
}

.cancel-btn {
  background: linear-gradient(135deg, #a0aec0 0%, #718096 100%);
  color: white;
}

.cancel-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(160, 174, 192, 0.4);
}

.save-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

/* 统计卡片样式 */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  border-radius: 16px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid rgba(226, 232, 240, 0.6);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

.stat-item::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 100%;
  background: linear-gradient(135deg, var(--gradient-start) 0%, var(--gradient-end) 100%);
  opacity: 0.1;
  transition: opacity 0.3s ease;
}

.stat-item:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.15);
}

.stat-item:hover::before {
  opacity: 0.15;
}

.stat-item.primary {
  --gradient-start: #667eea;
  --gradient-end: #764ba2;
  border-left: 4px solid #667eea;
}

.stat-item.success {
  --gradient-start: #10b981;
  --gradient-end: #059669;
  border-left: 4px solid #10b981;
}

.stat-item.danger {
  --gradient-start: #ef4444;
  --gradient-end: #dc2626;
  border-left: 4px solid #ef4444;
}

.stat-item.info {
  --gradient-start: #3b82f6;
  --gradient-end: #2563eb;
  border-left: 4px solid #3b82f6;
}

.stat-item.warning {
  --gradient-start: #f59e0b;
  --gradient-end: #d97706;
  border-left: 4px solid #f59e0b;
}

.stat-item.secondary {
  --gradient-start: #6b7280;
  --gradient-end: #4b5563;
  border-left: 4px solid #6b7280;
}

.stat-icon {
  font-size: 2rem;
  opacity: 0.9;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
  background: linear-gradient(135deg, var(--gradient-start) 0%, var(--gradient-end) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.stat-content {
  flex: 1;
  z-index: 1;
}

.stat-value {
  font-size: 2rem;
  font-weight: 800;
  margin-bottom: 6px;
  letter-spacing: -0.02em;
  background: linear-gradient(135deg, #1a202c 0%, #475569 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.stat-label {
  font-size: 0.85rem;
  opacity: 0.7;
  font-weight: 500;
  letter-spacing: 0.02em;
  text-transform: uppercase;
}

/* 历史记录样式 */
.history-list {
  max-height: 500px;
  overflow-y: auto;
  margin-bottom: 24px;
  padding-right: 8px;
}

.history-list::-webkit-scrollbar {
  width: 8px;
}

.history-list::-webkit-scrollbar-track {
  background: rgba(226, 232, 240, 0.3);
  border-radius: 4px;
}

.history-list::-webkit-scrollbar-thumb {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 4px;
  opacity: 0.7;
}

.history-list::-webkit-scrollbar-thumb:hover {
  opacity: 0.9;
}

.history-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border: 1px solid rgba(226, 232, 240, 0.6);
  border-radius: 16px;
  margin-bottom: 16px;
  background: rgba(255, 255, 255, 0.8);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
}

.history-item::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 4px;
  background: var(--result-color);
  transform: scaleY(0);
  transition: transform 0.3s ease;
}

.history-item:hover {
  transform: translateX(8px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
  border-color: var(--result-color);
}

.history-item:hover::before {
  transform: scaleY(1);
}

.history-item.win {
  --result-color: #10b981;
}

.history-item.loss {
  --result-color: #ef4444;
}

.history-item.draw {
  --result-color: #f59e0b;
}

.game-info {
  display: flex;
  flex-direction: column;
  gap: 12px;
  flex: 1;
}

.game-type {
  font-weight: 700;
  color: #1a202c;
  font-size: 1.1rem;
  letter-spacing: -0.01em;
}

.game-result {
  padding: 6px 16px;
  border-radius: 24px;
  font-size: 0.8rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  position: relative;
  z-index: 1;
}

.game-result.win {
  background: linear-gradient(135deg, #d1fae5 0%, #a7f3d0 100%);
  color: #064e3b;
  box-shadow: 0 2px 8px rgba(16, 185, 129, 0.2);
}

.game-result.loss {
  background: linear-gradient(135deg, #fee2e2 0%, #fecaca 100%);
  color: #7f1d1d;
  box-shadow: 0 2px 8px rgba(239, 68, 68, 0.2);
}

.game-result.draw {
  background: linear-gradient(135deg, #fef3c7 0%, #fde68a 100%);
  color: #78350f;
  box-shadow: 0 2px 8px rgba(245, 158, 11, 0.2);
}

.game-details {
  display: flex;
  flex-direction: column;
  gap: 8px;
  font-size: 0.9rem;
  color: #64748b;
  margin-left: 24px;
}

.detail-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 8px;
  background: rgba(248, 250, 252, 0.8);
  border-radius: 8px;
  border: 1px solid rgba(226, 232, 240, 0.4);
}

.game-date {
  font-size: 0.8rem;
  color: #94a3b8;
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: rgba(241, 245, 249, 0.8);
  border-radius: 12px;
  border: 1px solid rgba(226, 232, 240, 0.4);
}

/* 分页样式 */
.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 20px;
  margin-top: 28px;
  padding: 20px;
  background: rgba(255, 255, 255, 0.6);
  border-radius: 16px;
  border: 1px solid rgba(226, 232, 240, 0.4);
}

.page-btn {
  padding: 12px 24px;
  border: 2px solid rgba(226, 232, 240, 0.6);
  background: rgba(255, 255, 255, 0.9);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  font-size: 0.95rem;
  position: relative;
  overflow: hidden;
}

.page-btn::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(102, 126, 234, 0.1), transparent);
  transition: left 0.5s;
}

.page-btn:hover::before {
  left: 100%;
}

.page-btn:hover:not(:disabled) {
  background: rgba(102, 126, 234, 0.1);
  border-color: #667eea;
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(102, 126, 234, 0.2);
}

.page-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}

.page-info {
  font-weight: 600;
  color: #1a202c;
  font-size: 1.1rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

/* 密码修改样式 */
.password-form {
  max-width: 440px;
}

.form-group {
  margin-bottom: 24px;
}

.form-label {
  display: block;
  margin-bottom: 12px;
  font-weight: 700;
  color: #374151;
  font-size: 1rem;
  letter-spacing: -0.01em;
}

.label-icon {
  font-size: 1.2rem;
  margin-right: 8px;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
}

.form-input {
  width: 100%;
  padding: 16px 20px;
  border: 2px solid rgba(226, 232, 240, 0.6);
  border-radius: 12px;
  font-size: 1rem;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  background: rgba(255, 255, 255, 0.9);
  position: relative;
}

.form-input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 4px rgba(102, 126, 234, 0.1), 0 8px 24px rgba(102, 126, 234, 0.1);
  transform: translateY(-1px);
}

.password-strength {
  margin-bottom: 28px;
}

.strength-info {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: rgba(248, 250, 252, 0.8);
  border-radius: 12px;
  border: 1px solid rgba(226, 232, 240, 0.4);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.strength-text {
  font-size: 0.9rem;
  color: #475569;
  font-weight: 600;
  letter-spacing: 0.01em;
}

.requirement {
  font-size: 0.9rem;
  color: #667eea;
  font-weight: 700;
  letter-spacing: 0.01em;
}

.change-password-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 16px 28px;
  border-radius: 12px;
  cursor: pointer;
  font-weight: 700;
  font-size: 1rem;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  width: 100%;
  justify-content: center;
  position: relative;
  overflow: hidden;
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.3);
}

.change-password-btn::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  transition: left 0.5s;
}

.change-password-btn:hover::before {
  left: 100%;
}

.change-password-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 12px 32px rgba(102, 126, 234, 0.4);
}

.change-password-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
  transform: none;
}

/* 加载和无数据样式 */
.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 40px;
  color: #64748b;
}

.loading-spinner {
  width: 32px;
  height: 32px;
  border: 4px solid rgba(226, 232, 240, 0.6);
  border-top: 4px solid #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 20px;
  box-shadow: 0 4px 16px rgba(102, 126, 234, 0.2);
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.no-data {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 40px;
  text-align: center;
  color: #94a3b8;
  background: rgba(255, 255, 255, 0.6);
  border-radius: 20px;
  border: 1px solid rgba(226, 232, 240, 0.4);
}

.no-data-icon {
  font-size: 5rem;
  margin-bottom: 24px;
  opacity: 0.8;
  filter: drop-shadow(0 4px 8px rgba(0, 0, 0, 0.1));
}

.no-data-hint {
  font-size: 1.1rem;
  color: #64748b;
  margin-top: 12px;
  font-weight: 500;
}

/* 消息提示样式 */
.error-message,
.success-message {
  position: fixed;
  top: 24px;
  right: 24px;
  padding: 20px 24px;
  border-radius: 16px;
  border: 1px solid;
  z-index: 1000;
  max-width: 380px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: 0 16px 48px rgba(0, 0, 0, 0.15);
  backdrop-filter: blur(20px);
  transform: translateX(400px);
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.error-message.show,
.success-message.show {
  transform: translateX(0);
}

.error-message {
  background: rgba(239, 68, 68, 0.1);
  border-color: rgba(239, 68, 68, 0.3);
  color: #dc2626;
  box-shadow: 0 16px 48px rgba(239, 68, 68, 0.2);
}

.success-message {
  background: rgba(16, 185, 129, 0.1);
  border-color: rgba(16, 185, 129, 0.3);
  color: #059669;
  box-shadow: 0 16px 48px rgba(16, 185, 129, 0.2);
}

.message-icon {
  font-size: 1.4rem;
  margin-right: 12px;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
}

.close-btn {
  background: none;
  border: none;
  font-size: 1.4rem;
  cursor: pointer;
  color: inherit;
  padding: 6px;
  border-radius: 8px;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.close-btn:hover {
  background: rgba(0, 0, 0, 0.1);
  transform: rotate(90deg);
}

/* 过渡动画 */
.slide-in-enter-active,
.slide-in-leave-active {
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.slide-in-enter-from {
  transform: translateX(100%);
  opacity: 0;
}

.slide-in-leave-to {
  transform: translateX(100%);
  opacity: 0;
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .user-profile {
    padding: 16px;
  }

  .content-row {
    grid-template-columns: 1fr;
    gap: 24px;
  }

  .stats-grid {
    grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
    gap: 16px;
  }

  .profile-card,
  .stats-card,
  .history-card,
  .password-card {
    padding: 24px;
  }

  .header-title h1 {
    font-size: 2rem;
  }

  .card-header h2 {
    font-size: 1.5rem;
  }

  .stat-value {
    font-size: 1.8rem;
  }

  .stat-label {
    font-size: 0.8rem;
  }
}

@media (max-width: 768px) {
  .user-profile {
    padding: 12px;
  }

  .profile-header {
    padding: 24px;
    margin-bottom: 24px;
  }

  .header-content {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }

  .header-title h1 {
    font-size: 1.8rem;
  }

  .welcome-text {
    font-size: 1rem;
  }

  .logout-btn {
    align-self: stretch;
    justify-content: center;
    padding: 14px 24px;
  }

  .content-row {
    gap: 20px;
  }

  .profile-card,
  .stats-card,
  .history-card,
  .password-card {
    padding: 20px;
    border-radius: 16px;
  }

  .card-header {
    margin-bottom: 20px;
    padding-bottom: 16px;
  }

  .card-header h2 {
    font-size: 1.3rem;
  }

  .stats-grid {
    grid-template-columns: 1fr;
    gap: 12px;
  }

  .stat-item {
    padding: 16px;
    border-radius: 12px;
  }

  .stat-value {
    font-size: 1.6rem;
  }

  .history-list {
    max-height: 400px;
  }

  .history-item {
    padding: 16px;
    margin-bottom: 12px;
    border-radius: 12px;
  }

  .game-details {
    flex-direction: row;
    flex-wrap: wrap;
    gap: 8px;
    margin-left: 16px;
  }

  .detail-item {
    padding: 6px 12px;
    font-size: 0.8rem;
  }

  .game-date {
    align-self: flex-start;
    font-size: 0.75rem;
    padding: 4px 8px;
  }

  .pagination {
    flex-wrap: wrap;
    gap: 12px;
    padding: 16px;
  }

  .page-btn {
    flex: 1;
    min-width: 100px;
    padding: 10px 16px;
    font-size: 0.9rem;
  }

  .page-info {
    font-size: 1rem;
  }

  .password-form {
    max-width: 100%;
  }

  .form-group {
    margin-bottom: 20px;
  }

  .form-input {
    padding: 14px 16px;
    font-size: 0.95rem;
  }

  .change-password-btn {
    padding: 14px 24px;
    font-size: 0.95rem;
  }

  .error-message,
  .success-message {
    top: 16px;
    right: 16px;
    left: 16px;
    max-width: none;
    transform: translateY(-100px);
  }

  .error-message.show,
  .success-message.show {
    transform: translateY(0);
  }
}

@media (max-width: 480px) {
  .user-profile {
    padding: 10px;
  }

  .profile-header {
    padding: 20px;
    border-radius: 12px;
  }

  .header-title h1 {
    font-size: 1.5rem;
  }

  .profile-card,
  .stats-card,
  .history-card,
  .password-card {
    padding: 16px;
    border-radius: 12px;
  }

  .card-header h2 {
    font-size: 1.2rem;
  }

  .stat-value {
    font-size: 1.4rem;
  }

  .stat-label {
    font-size: 0.75rem;
  }

  .history-item {
    padding: 12px;
    margin-bottom: 8px;
  }

  .game-details {
    gap: 6px;
  }

  .detail-item {
    padding: 4px 8px;
    font-size: 0.75rem;
  }

  .game-date {
    font-size: 0.7rem;
  }

  .loading {
    padding: 40px 20px;
  }

  .no-data {
    padding: 60px 20px;
  }

  .no-data-icon {
    font-size: 4rem;
    margin-bottom: 16px;
  }

  .no-data-hint {
    font-size: 1rem;
  }
}

  </style>