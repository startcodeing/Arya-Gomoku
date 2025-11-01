<template>
  <div class="user-profile">
    <!-- 页面头部 -->
    <div class="profile-header">
      <div class="header-content">
        <h1>用户中心</h1>
        <button @click="logout" class="logout-btn">
          退出登录
        </button>
      </div>
    </div>

    <!-- 主要内容 -->
    <div class="profile-content">
      <!-- 用户信息卡片 -->
      <div class="profile-card">
        <h2>个人信息</h2>
        <div class="user-info">
          <div class="info-item">
            <label>用户名:</label>
            <span>{{ userStore.user?.username }}</span>
          </div>
          <div class="info-item">
            <label>昵称:</label>
            <span v-if="!editingProfile">{{ userStore.user?.nickname || '未设置' }}</span>
            <input 
              v-else 
              v-model="profileForm.nickname" 
              type="text" 
              placeholder="请输入昵称"
              class="edit-input"
            >
          </div>
          <div class="info-item">
            <label>邮箱:</label>
            <span v-if="!editingProfile">{{ userStore.user?.email || '未设置' }}</span>
            <input 
              v-else 
              v-model="profileForm.email" 
              type="email" 
              placeholder="请输入邮箱"
              class="edit-input"
            >
          </div>
          <div class="info-item">
            <label>注册时间:</label>
            <span>{{ formatDate(userStore.user?.created_at) }}</span>
          </div>
        </div>
        
        <div class="profile-actions">
          <button 
            v-if="!editingProfile" 
            @click="startEditProfile" 
            class="edit-btn"
          >
            编辑资料
          </button>
          <template v-else>
            <button @click="saveProfile" class="save-btn" :disabled="profileLoading">
              {{ profileLoading ? '保存中...' : '保存' }}
            </button>
            <button @click="cancelEditProfile" class="cancel-btn">
              取消
            </button>
          </template>
        </div>
      </div>

      <!-- 游戏统计卡片 -->
      <div class="stats-card">
        <h2>游戏统计</h2>
        <div v-if="statsLoading" class="loading">
          加载中...
        </div>
        <div v-else-if="gameStats" class="stats-grid">
          <div class="stat-item">
            <div class="stat-value">{{ gameStats.total_games || 0 }}</div>
            <div class="stat-label">总游戏数</div>
          </div>
          <div class="stat-item">
            <div class="stat-value">{{ gameStats.wins || 0 }}</div>
            <div class="stat-label">胜利次数</div>
          </div>
          <div class="stat-item">
            <div class="stat-value">{{ gameStats.losses || 0 }}</div>
            <div class="stat-label">失败次数</div>
          </div>
          <div class="stat-item">
            <div class="stat-value">{{ winRate }}%</div>
            <div class="stat-label">胜率</div>
          </div>
          <div class="stat-item">
            <div class="stat-value">{{ gameStats.avg_moves || 0 }}</div>
            <div class="stat-label">平均步数</div>
          </div>
          <div class="stat-item">
            <div class="stat-value">{{ formatDuration(gameStats.avg_duration) }}</div>
            <div class="stat-label">平均时长</div>
          </div>
        </div>
      </div>

      <!-- 游戏历史卡片 -->
      <div class="history-card">
        <h2>游戏历史</h2>
        <div v-if="historyLoading" class="loading">
          加载中...
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
              <span>步数: {{ game.moves }}</span>
              <span>时长: {{ formatDuration(game.duration) }}</span>
              <span v-if="game.difficulty">难度: {{ getDifficultyText(game.difficulty) }}</span>
            </div>
            <div class="game-date">
              {{ formatDate(game.created_at) }}
            </div>
          </div>
          
          <!-- 分页 -->
          <div class="pagination">
            <button 
              @click="loadPreviousPage" 
              :disabled="currentPage <= 1"
              class="page-btn"
            >
              上一页
            </button>
            <span class="page-info">第 {{ currentPage }} 页</span>
            <button 
              @click="loadNextPage" 
              :disabled="gameHistory.length < pageSize"
              class="page-btn"
            >
              下一页
            </button>
          </div>
        </div>
        <div v-else class="no-data">
          暂无游戏记录
        </div>
      </div>

      <!-- 密码修改卡片 -->
      <div class="password-card">
        <h2>修改密码</h2>
        <form @submit.prevent="changePassword" class="password-form">
          <div class="form-group">
            <label>当前密码:</label>
            <input 
              v-model="passwordForm.oldPassword" 
              type="password" 
              required
              class="form-input"
            >
          </div>
          <div class="form-group">
            <label>新密码:</label>
            <input 
              v-model="passwordForm.newPassword" 
              type="password" 
              required
              minlength="8"
              class="form-input"
            >
          </div>
          <div class="form-group">
            <label>确认新密码:</label>
            <input 
              v-model="passwordForm.confirmPassword" 
              type="password" 
              required
              class="form-input"
            >
          </div>
          <button 
            type="submit" 
            class="change-password-btn"
            :disabled="passwordLoading"
          >
            {{ passwordLoading ? '修改中...' : '修改密码' }}
          </button>
        </form>
      </div>
    </div>

    <!-- 错误提示 -->
    <div v-if="error" class="error-message">
      {{ error }}
    </div>

    <!-- 成功提示 -->
    <div v-if="successMessage" class="success-message">
      {{ successMessage }}
    </div>
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
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  background: #f5f5f5;
  min-height: 100vh;
}

.profile-header {
  background: white;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-content h1 {
  margin: 0;
  color: #333;
}

.logout-btn {
  background: #dc3545;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.logout-btn:hover {
  background: #c82333;
}

.profile-content {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

@media (max-width: 768px) {
  .profile-content {
    grid-template-columns: 1fr;
  }
}

.profile-card,
.stats-card,
.history-card,
.password-card {
  background: white;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.history-card {
  grid-column: 1 / -1;
}

.profile-card h2,
.stats-card h2,
.history-card h2,
.password-card h2 {
  margin: 0 0 20px 0;
  color: #333;
  border-bottom: 2px solid #007bff;
  padding-bottom: 10px;
}

.user-info {
  margin-bottom: 20px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid #eee;
}

.info-item:last-child {
  border-bottom: none;
}

.info-item label {
  font-weight: bold;
  color: #555;
}

.edit-input {
  padding: 4px 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
  width: 200px;
}

.profile-actions {
  display: flex;
  gap: 10px;
}

.edit-btn,
.save-btn,
.cancel-btn {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.edit-btn {
  background: #007bff;
  color: white;
}

.edit-btn:hover {
  background: #0056b3;
}

.save-btn {
  background: #28a745;
  color: white;
}

.save-btn:hover:not(:disabled) {
  background: #1e7e34;
}

.cancel-btn {
  background: #6c757d;
  color: white;
}

.cancel-btn:hover {
  background: #545b62;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
  gap: 15px;
}

.stat-item {
  text-align: center;
  padding: 15px;
  background: #f8f9fa;
  border-radius: 8px;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #007bff;
  margin-bottom: 5px;
}

.stat-label {
  font-size: 14px;
  color: #666;
}

.history-list {
  max-height: 400px;
  overflow-y: auto;
}

.history-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px;
  border: 1px solid #eee;
  border-radius: 8px;
  margin-bottom: 10px;
  background: #f8f9fa;
}

.game-info {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.game-type {
  font-weight: bold;
  color: #333;
}

.game-result {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: bold;
}

.game-result.win {
  background: #d4edda;
  color: #155724;
}

.game-result.loss {
  background: #f8d7da;
  color: #721c24;
}

.game-result.draw {
  background: #fff3cd;
  color: #856404;
}

.game-details {
  display: flex;
  flex-direction: column;
  gap: 5px;
  font-size: 14px;
  color: #666;
}

.game-date {
  font-size: 12px;
  color: #999;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 15px;
  margin-top: 20px;
}

.page-btn {
  padding: 8px 16px;
  border: 1px solid #ddd;
  background: white;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.page-btn:hover:not(:disabled) {
  background: #f8f9fa;
}

.page-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.page-info {
  font-weight: bold;
  color: #333;
}

.password-form {
  max-width: 300px;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
  color: #555;
}

.form-input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.change-password-btn {
  background: #007bff;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.change-password-btn:hover:not(:disabled) {
  background: #0056b3;
}

.loading {
  text-align: center;
  padding: 20px;
  color: #666;
}

.no-data {
  text-align: center;
  padding: 40px;
  color: #999;
  font-style: italic;
}

.error-message {
  position: fixed;
  top: 20px;
  right: 20px;
  background: #f8d7da;
  color: #721c24;
  padding: 15px 20px;
  border-radius: 4px;
  border: 1px solid #f5c6cb;
  z-index: 1000;
  max-width: 300px;
}

.success-message {
  position: fixed;
  top: 20px;
  right: 20px;
  background: #d4edda;
  color: #155724;
  padding: 15px 20px;
  border-radius: 4px;
  border: 1px solid #c3e6cb;
  z-index: 1000;
  max-width: 300px;
}
</style>