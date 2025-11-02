<template>
  <div class="user-profile-container">
    <!-- 背景装饰 -->
    <div class="background-decoration">
      <div class="decoration-circle circle-1"></div>
      <div class="decoration-circle circle-2"></div>
      <div class="decoration-circle circle-3"></div>
    </div>

    <!-- 页面头部 -->
    <div class="profile-header">
      <div class="header-content">
        <div class="header-left">
          <div class="welcome-section">
            <h1 class="page-title">用户中心</h1>
            <p class="page-subtitle">管理您的个人信息和游戏数据</p>
          </div>
        </div>
        <div class="header-right">
          <button @click="goToHome" class="home-button">
            <svg class="home-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"></path>
            </svg>
            返回主页
          </button>
          <button @click="logout" class="logout-button">
            <svg class="logout-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"></path>
            </svg>
            退出登录
          </button>
        </div>
      </div>
    </div>

    <!-- 主要内容区域 -->
    <div class="main-content">
      <!-- 用户信息卡片 -->
      <div class="profile-card modern-card">
        <div class="card-header">
          <div class="header-info">
            <h2 class="card-title">个人信息</h2>
            <p class="card-subtitle">您的基本资料信息</p>
          </div>
          <div class="header-icon">
            <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
            </svg>
          </div>
        </div>

        <div class="user-profile-section">
          <div class="profile-avatar-section">
            <div class="avatar-container">
              <div class="avatar-circle">
                <span class="avatar-text">
                  {{ userStore.userDisplayName.charAt(0).toUpperCase() }}
                </span>
              </div>
              <div class="avatar-status"></div>
            </div>
            <div class="profile-info">
              <h3 class="profile-name">{{ userStore.userDisplayName }}</h3>
              <p class="profile-username">@{{ userStore.user?.username }}</p>
              <p class="profile-email">{{ userStore.user?.email }}</p>
            </div>
          </div>

          <div class="profile-details">
            <div class="detail-grid">
              <div class="detail-item">
                <div class="detail-label">用户名</div>
                <div class="detail-value">{{ userStore.user?.username }}</div>
              </div>
              <div class="detail-item">
                <div class="detail-label">昵称</div>
                <div class="detail-value">{{ userStore.user?.nickname || '未设置' }}</div>
              </div>
              <div class="detail-item">
                <div class="detail-label">邮箱</div>
                <div class="detail-value">{{ userStore.user?.email || '未设置' }}</div>
              </div>
              <div class="detail-item">
                <div class="detail-label">注册时间</div>
                <div class="detail-value">{{ formatDate(userStore.user?.createdAt) }}</div>
              </div>
            </div>
          </div>

          <div class="profile-actions">
            <button
              @click="showEditProfile = true"
              class="edit-profile-btn primary-btn"
            >
              <svg class="btn-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
              </svg>
              编辑资料
            </button>
            <button
              @click="showChangePassword = true"
              class="change-password-btn primary-btn"
            >
              <svg class="btn-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m0 0a2 2 0 012 2v6a2 2 0 01-2 2H7a2 2 0 01-2-2V9a2 2 0 012-2m0 0V7a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path>
              </svg>
              修改密码
            </button>
          </div>
        </div>
      </div>

      <!-- 游戏统计卡片 -->
      <div class="stats-card modern-card">
        <div class="card-header">
          <div class="header-info">
            <h2 class="card-title">游戏统计</h2>
            <p class="card-subtitle">您的游戏表现数据</p>
          </div>
          <div class="header-icon">
            <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
            </svg>
          </div>
        </div>
        <div class="stats-content">
          <div v-if="statsLoading" class="loading-state">
            <div class="loading-spinner"></div>
            <p class="loading-text">加载统计数据中...</p>
          </div>
          <div v-else-if="gameStats" class="stats-grid">
            <!-- AI游戏统计 -->
            <div class="stat-card stat-ai">
              <div class="stat-header">
                <div class="stat-icon">
                  <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"></path>
                  </svg>
                </div>
                <h3 class="stat-title">AI对战</h3>
              </div>
              <div class="stat-metrics">
                <div class="metric-item">
                  <span class="metric-label">总场次</span>
                  <span class="metric-value">{{ gameStats.ai.totalGames }}</span>
                </div>
                <div class="metric-item">
                  <span class="metric-label">胜利</span>
                  <span class="metric-value metric-win">{{ gameStats.ai.wins }}</span>
                </div>
                <div class="metric-item">
                  <span class="metric-label">失败</span>
                  <span class="metric-value metric-loss">{{ gameStats.ai.losses }}</span>
                </div>
                <div class="metric-item">
                  <span class="metric-label">胜率</span>
                  <span class="metric-value metric-rate">{{ gameStats.ai.winRate }}%</span>
                </div>
              </div>
            </div>

            <!-- LLM游戏统计 -->
            <div class="stat-card stat-llm">
              <div class="stat-header">
                <div class="stat-icon">
                  <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path>
                  </svg>
                </div>
                <h3 class="stat-title">LLM对战</h3>
              </div>
              <div class="stat-metrics">
                <div class="metric-item">
                  <span class="metric-label">总场次</span>
                  <span class="metric-value">{{ gameStats.llm.totalGames }}</span>
                </div>
                <div class="metric-item">
                  <span class="metric-label">胜利</span>
                  <span class="metric-value metric-win">{{ gameStats.llm.wins }}</span>
                </div>
                <div class="metric-item">
                  <span class="metric-label">失败</span>
                  <span class="metric-value metric-loss">{{ gameStats.llm.losses }}</span>
                </div>
                <div class="metric-item">
                  <span class="metric-label">胜率</span>
                  <span class="metric-value metric-rate">{{ gameStats.llm.winRate }}%</span>
                </div>
              </div>
            </div>

            <!-- PVP游戏统计 -->
            <div class="stat-card stat-pvp">
              <div class="stat-header">
                <div class="stat-icon">
                  <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"></path>
                  </svg>
                </div>
                <h3 class="stat-title">玩家对战</h3>
              </div>
              <div class="stat-metrics">
                <div class="metric-item">
                  <span class="metric-label">总场次</span>
                  <span class="metric-value">{{ gameStats.pvp.totalGames }}</span>
                </div>
                <div class="metric-item">
                  <span class="metric-label">胜利</span>
                  <span class="metric-value metric-win">{{ gameStats.pvp.wins }}</span>
                </div>
                <div class="metric-item">
                  <span class="metric-label">失败</span>
                  <span class="metric-value metric-loss">{{ gameStats.pvp.losses }}</span>
                </div>
                <div class="metric-item">
                  <span class="metric-label">胜率</span>
                  <span class="metric-value metric-rate">{{ gameStats.pvp.winRate }}%</span>
                </div>
              </div>
            </div>
          </div>
          <div v-else class="empty-state">
            <div class="empty-icon">
              <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
              </svg>
            </div>
            <p class="empty-text">暂无游戏统计数据</p>
            <p class="empty-subtext">开始您的第一场游戏吧！</p>
          </div>
        </div>
      </div>

      <!-- 最近游戏记录 -->
      <div class="recent-games-card modern-card">
        <div class="card-header">
          <div class="header-info">
            <h2 class="card-title">最近游戏</h2>
            <p class="card-subtitle">您的游戏历史记录</p>
          </div>
          <div class="header-actions">
            <button
              @click="loadGameHistory"
              class="refresh-btn"
            >
              <svg class="refresh-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
              </svg>
              刷新
            </button>
          </div>
        </div>
        
        <div class="games-content">
          <div v-if="historyLoading" class="loading-state">
            <div class="loading-spinner"></div>
            <p class="loading-text">加载游戏记录中...</p>
          </div>
          <div v-else-if="gameHistory.length > 0" class="games-list">
            <div
              v-for="game in gameHistory"
              :key="game.id"
              class="game-item"
            >
              <div class="game-status">
                <div
                  class="status-indicator"
                  :class="{
                    'status-win': game.result === 'win',
                    'status-loss': game.result === 'loss',
                    'status-draw': game.result === 'draw'
                  }"
                ></div>
              </div>
              
              <div class="game-info">
                <div class="game-details">
                  <div class="game-tags">
                    <span class="game-type-tag" :class="`tag-${game.gameType}`">
                      {{ getGameTypeName(game.gameType) }}
                    </span>
                    <span class="result-tag" :class="`result-${game.result}`">
                      {{ getResultText(game.result) }}
                    </span>
                  </div>
                  <div class="game-meta">
                    <span class="game-duration">
                      <svg class="meta-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                      </svg>
                      {{ formatDuration(game.duration) }}
                    </span>
                    <span class="game-date">
                      <svg class="meta-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
                      </svg>
                      {{ formatDate(game.createdAt) }}
                    </span>
                  </div>
                </div>
              </div>
              
              <div class="game-actions">
                <button class="view-details-btn">
                  <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
                  </svg>
                </button>
              </div>
            </div>
          </div>
          <div v-else class="empty-state">
            <div class="empty-icon">
              <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path>
              </svg>
            </div>
            <p class="empty-text">暂无游戏历史记录</p>
            <p class="empty-subtext">您还没有进行过游戏</p>
          </div>
        </div>
      </div>
    </div>

    <!-- 编辑资料弹窗 -->
    <EditProfileModal
      v-if="showEditProfile"
      @close="showEditProfile = false"
      @updated="handleProfileUpdated"
    />

    <!-- 修改密码弹窗 -->
    <ChangePasswordModal
      v-if="showChangePassword"
      @close="showChangePassword = false"
      @success="showChangePassword = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../../stores/user'
import { statisticsApi } from '../../services/statistics'
import type { UserGameStats, GameSummary } from '../../services/statistics'
import EditProfileModal from './EditProfileModal.vue'
import ChangePasswordModal from './ChangePasswordModal.vue'

const router = useRouter()
const userStore = useUserStore()

const showEditProfile = ref(false)
const showChangePassword = ref(false)
const statsLoading = ref(false)
const historyLoading = ref(false)
const gameStats = ref<UserGameStats | null>(null)
const gameHistory = ref<GameSummary[]>([])

// 格式化日期
const formatDate = (dateString?: string) => {
  if (!dateString) return '未知'
  return new Date(dateString).toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

// 格式化游戏时长
const formatDuration = (duration: number) => {
  const minutes = Math.floor(duration / 60000)
  const seconds = Math.floor((duration % 60000) / 1000)
  return `${minutes}分${seconds}秒`
}

// 获取游戏类型名称
const getGameTypeName = (type: string) => {
  const typeMap: Record<string, string> = {
    'ai': 'AI对战',
    'llm': 'LLM对战',
    'pvp': '玩家对战'
  }
  return typeMap[type] || type
}

// 获取游戏类型样式
const getGameTypeClass = (type: string) => {
  const classMap: Record<string, string> = {
    'ai': 'bg-blue-100 text-blue-800',
    'llm': 'bg-green-100 text-green-800',
    'pvp': 'bg-purple-100 text-purple-800'
  }
  return classMap[type] || 'bg-gray-100 text-gray-800'
}

// 获取游戏类型文本
const getGameTypeText = (type: string) => {
  const typeMap: Record<string, string> = {
    'ai': 'AI对战',
    'llm': 'LLM对战', 
    'pvp': '玩家对战'
  }
  return typeMap[type] || type
}

// 退出登录
const logout = async () => {
  if (confirm('确定要退出登录吗？')) {
    try {
      await userStore.logout()
      router.push('/login')
    } catch (error) {
      console.error('Logout failed:', error)
    }
  }
}

// 处理资料更新
const handleProfileUpdated = () => {
  showEditProfile.value = false
  // 刷新用户信息
  userStore.fetchUserProfile().catch(console.error)
}

// 加载游戏统计数据
const loadGameStats = async () => {
  try {
    statsLoading.value = true
    const stats = await statisticsApi.getUserGameStats()
    gameStats.value = stats
  } catch (error) {
    console.error('Failed to load game stats:', error)
  } finally {
    statsLoading.value = false
  }
}

// 加载游戏历史记录
const loadGameHistory = async () => {
  try {
    historyLoading.value = true
    const history = await statisticsApi.getUserGameHistory()
    gameHistory.value = history
  } catch (error) {
    console.error('Failed to load game history:', error)
  } finally {
    historyLoading.value = false
  }
}

// 返回游戏主页
const goToHome = () => {
  router.push('/')
}

// 组件挂载时加载数据
onMounted(() => {
  loadGameStats()
  loadGameHistory()
})
</script>

<style scoped>
/* 主容器 */
.user-profile-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  position: relative;
  overflow: hidden;
}

/* 背景装饰 */
.background-decoration {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  pointer-events: none;
  z-index: 0;
}

.decoration-circle {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
}

.circle-1 {
  width: 300px;
  height: 300px;
  top: -150px;
  right: -150px;
  animation: float 6s ease-in-out infinite;
}

.circle-2 {
  width: 200px;
  height: 200px;
  bottom: -100px;
  left: -100px;
  animation: float 8s ease-in-out infinite reverse;
}

.circle-3 {
  width: 150px;
  height: 150px;
  top: 50%;
  left: 10%;
  animation: float 10s ease-in-out infinite;
}

@keyframes float {
  0%, 100% { transform: translateY(0px) rotate(0deg); }
  50% { transform: translateY(-20px) rotate(180deg); }
}

/* 页面头部 */
.profile-header {
  position: relative;
  z-index: 1;
  padding: 1.5rem 0 0.75rem;
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 1.5rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.welcome-section .page-title {
  font-size: 2.25rem;
  font-weight: 700;
  color: white;
  margin: 0 0 0.25rem 0;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.welcome-section .page-subtitle {
  font-size: 1.1rem;
  color: rgba(255, 255, 255, 0.9);
  margin: 0;
}

/* 头部按钮样式 */
.header-right {
  display: flex;
  gap: 1rem;
}

.home-button {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  background: rgba(255, 255, 255, 0.2);
  color: white;
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 12px;
  font-weight: 500;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
}

.home-button:hover {
  background: rgba(255, 255, 255, 0.3);
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
}

.home-icon {
  width: 1.25rem;
  height: 1.25rem;
}

.logout-button {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  background: rgba(255, 255, 255, 0.2);
  color: white;
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 12px;
  font-weight: 500;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
}

.logout-button:hover {
  background: rgba(255, 255, 255, 0.3);
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
}

.logout-icon {
  width: 1.25rem;
  height: 1.25rem;
}

/* 主要内容区域 */
.main-content {
  position: relative;
  z-index: 1;
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 1.5rem 2rem;
  display: grid;
  gap: 1.25rem;
  grid-template-columns: 1fr;
}

/* 现代化卡片样式 */
.modern-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 20px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  overflow: hidden;
  transition: all 0.3s ease;
}

.modern-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.15);
}

/* 卡片头部 */
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem 1.5rem 0.75rem;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

.header-info .card-title {
  font-size: 1.5rem;
  font-weight: 700;
  color: #1a202c;
  margin: 0 0 0.25rem 0;
}

.header-info .card-subtitle {
  font-size: 0.9rem;
  color: #718096;
  margin: 0;
}

.header-icon {
  width: 2.5rem;
  height: 2.5rem;
  background: linear-gradient(135deg, #667eea, #764ba2);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.header-icon svg {
  width: 1.25rem;
  height: 1.25rem;
}

/* 用户资料部分 */
.user-profile-section {
  padding: 0.75rem 1.5rem 1.5rem;
}

.profile-avatar-section {
  display: flex;
  align-items: center;
  gap: 1.25rem;
  margin-bottom: 1.5rem;
}

.avatar-container {
  position: relative;
}

.avatar-circle {
  width: 5rem;
  height: 5rem;
  background: linear-gradient(135deg, #667eea, #764ba2);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8px 25px rgba(102, 126, 234, 0.3);
}

.avatar-text {
  font-size: 1.75rem;
  font-weight: 700;
  color: white;
}

.avatar-status {
  position: absolute;
  bottom: 0.25rem;
  right: 0.25rem;
  width: 1rem;
  height: 1rem;
  background: #48bb78;
  border: 2px solid white;
  border-radius: 50%;
}

.profile-info .profile-name {
  font-size: 1.5rem;
  font-weight: 700;
  color: #1a202c;
  margin: 0 0 0.25rem 0;
}

.profile-info .profile-username {
  font-size: 1rem;
  color: #667eea;
  margin: 0 0 0.25rem 0;
  font-weight: 500;
}

.profile-info .profile-email {
  font-size: 0.9rem;
  color: #718096;
  margin: 0;
}

/* 详细信息网格 */
.detail-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.detail-item {
  padding: 0.75rem;
  background: rgba(102, 126, 234, 0.05);
  border-radius: 12px;
  border-left: 4px solid #667eea;
}

.detail-label {
  font-size: 0.8rem;
  font-weight: 600;
  color: #718096;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: 0.5rem;
}

.detail-value {
  font-size: 1rem;
  font-weight: 600;
  color: #1a202c;
}

/* 按钮样式 */
.profile-actions {
  display: flex;
  gap: 0.75rem;
}

.primary-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border: none;
  border-radius: 12px;
  font-weight: 600;
  transition: all 0.3s ease;
  cursor: pointer;
}

.primary-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(102, 126, 234, 0.3);
}

.btn-icon {
  width: 1rem;
  height: 1rem;
}

.change-password-btn {
  background: linear-gradient(135deg, #f6ad55, #ed8936) !important;
}

.change-password-btn:hover {
  background: linear-gradient(135deg, #ed8936, #dd6b20) !important;
  box-shadow: 0 8px 25px rgba(237, 137, 54, 0.3) !important;
}

/* 统计卡片 */
.stats-content {
  padding: 0.75rem 1.5rem 1.5rem;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1rem;
}

.stat-card {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.9), rgba(255, 255, 255, 0.7));
  border-radius: 16px;
  padding: 1.25rem;
  border: 1px solid rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(10px);
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 12px 30px rgba(0, 0, 0, 0.1);
}

.stat-ai {
  border-left: 4px solid #3182ce;
}

.stat-llm {
  border-left: 4px solid #38a169;
}

.stat-pvp {
  border-left: 4px solid #805ad5;
}

.stat-header {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 0.75rem;
}

.stat-icon {
  width: 2rem;
  height: 2rem;
  background: linear-gradient(135deg, #667eea, #764ba2);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.stat-icon svg {
  width: 1rem;
  height: 1rem;
}

.stat-title {
  font-size: 1.1rem;
  font-weight: 700;
  color: #1a202c;
  margin: 0;
}

.stat-metrics {
  display: grid;
  gap: 0.75rem;
}

.metric-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem 0;
}

.metric-label {
  font-size: 0.9rem;
  color: #718096;
  font-weight: 500;
}

.metric-value {
  font-size: 1.1rem;
  font-weight: 700;
  color: #1a202c;
}

.metric-win {
  color: #38a169;
}

.metric-loss {
  color: #e53e3e;
}

.metric-rate {
  color: #667eea;
}

/* 游戏记录 */
.games-content {
  padding: 0.75rem 1.5rem 1.5rem;
}

.games-list {
  display: grid;
  gap: 0.75rem;
}

.game-item {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  background: rgba(255, 255, 255, 0.7);
  border-radius: 12px;
  border: 1px solid rgba(0, 0, 0, 0.05);
  transition: all 0.3s ease;
}

.game-item:hover {
  background: rgba(255, 255, 255, 0.9);
  transform: translateX(5px);
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
}

.status-indicator {
  width: 0.75rem;
  height: 0.75rem;
  border-radius: 50%;
  flex-shrink: 0;
}

.status-win {
  background: #38a169;
  box-shadow: 0 0 0 3px rgba(56, 161, 105, 0.2);
}

.status-lose {
  background: #e53e3e;
  box-shadow: 0 0 0 3px rgba(229, 62, 62, 0.2);
}

.status-draw {
  background: #d69e2e;
  box-shadow: 0 0 0 3px rgba(214, 158, 46, 0.2);
}

.game-info {
  flex: 1;
}

.game-details .game-opponent {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 0.25rem;
}

.opponent-name {
  font-weight: 600;
  color: #1a202c;
}

.game-type {
  padding: 0.25rem 0.5rem;
  border-radius: 6px;
  font-size: 0.75rem;
  font-weight: 500;
}

.type-ai {
  background: rgba(49, 130, 206, 0.1);
  color: #3182ce;
}

.type-llm {
  background: rgba(56, 161, 105, 0.1);
  color: #38a169;
}

.type-pvp {
  background: rgba(128, 90, 213, 0.1);
  color: #805ad5;
}

.game-time {
  font-size: 0.85rem;
  color: #718096;
}

.result-badge {
  padding: 0.5rem 1rem;
  border-radius: 8px;
  font-size: 0.85rem;
  font-weight: 600;
}

.badge-win {
  background: rgba(56, 161, 105, 0.1);
  color: #38a169;
}

.badge-lose {
  background: rgba(229, 62, 62, 0.1);
  color: #e53e3e;
}

.badge-draw {
  background: rgba(214, 158, 46, 0.1);
  color: #d69e2e;
}

/* 空状态 */
.empty-state {
  text-align: center;
  padding: 3rem 1rem;
}

.empty-icon {
  width: 4rem;
  height: 4rem;
  margin: 0 auto 1rem;
  background: rgba(113, 128, 150, 0.1);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #718096;
}

.empty-icon svg {
  width: 2rem;
  height: 2rem;
}

.empty-text {
  font-size: 1.1rem;
  font-weight: 600;
  color: #4a5568;
  margin: 0 0 0.5rem 0;
}

.empty-subtext {
  font-size: 0.9rem;
  color: #718096;
  margin: 0;
}

/* 加载状态 */
.loading-state {
  text-align: center;
  padding: 3rem 1rem;
}

.loading-spinner {
  width: 2rem;
  height: 2rem;
  border: 3px solid rgba(102, 126, 234, 0.2);
  border-top: 3px solid #667eea;
  border-radius: 50%;
  margin: 0 auto 1rem;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.loading-text {
  font-size: 0.9rem;
  color: #718096;
  margin: 0;
}

/* 刷新按钮 */
.header-actions {
  display: flex;
  gap: 0.5rem;
}

.refresh-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
  border: 1px solid rgba(102, 126, 234, 0.2);
  border-radius: 8px;
  font-size: 0.85rem;
  font-weight: 500;
  transition: all 0.3s ease;
  cursor: pointer;
}

.refresh-btn:hover {
  background: rgba(102, 126, 234, 0.2);
  transform: translateY(-1px);
}

.refresh-icon {
  width: 1rem;
  height: 1rem;
}

/* 游戏历史特殊样式 */
.game-tags {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 0.5rem;
}

.game-type-tag, .result-tag {
  padding: 0.25rem 0.5rem;
  border-radius: 6px;
  font-size: 0.75rem;
  font-weight: 500;
}

.tag-ai {
  background: rgba(49, 130, 206, 0.1);
  color: #3182ce;
}

.tag-llm {
  background: rgba(56, 161, 105, 0.1);
  color: #38a169;
}

.tag-pvp {
  background: rgba(128, 90, 213, 0.1);
  color: #805ad5;
}

.result-win {
  background: rgba(56, 161, 105, 0.1);
  color: #38a169;
}

.result-loss {
  background: rgba(229, 62, 62, 0.1);
  color: #e53e3e;
}

.result-draw {
  background: rgba(214, 158, 46, 0.1);
  color: #d69e2e;
}

.game-meta {
  display: flex;
  gap: 1rem;
  font-size: 0.8rem;
  color: #718096;
}

.game-duration, .game-date {
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.meta-icon {
  width: 0.875rem;
  height: 0.875rem;
}

.game-actions {
  display: flex;
  align-items: center;
}

.view-details-btn {
  padding: 0.5rem;
  background: transparent;
  border: none;
  color: #718096;
  border-radius: 6px;
  transition: all 0.3s ease;
  cursor: pointer;
}

.view-details-btn:hover {
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
}

.view-details-btn svg {
  width: 1rem;
  height: 1rem;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .header-content {
    flex-direction: column;
    gap: 1rem;
    text-align: center;
  }

  .main-content {
    padding: 0 1rem 2rem;
  }

  .profile-avatar-section {
    flex-direction: column;
    text-align: center;
  }

  .detail-grid {
    grid-template-columns: 1fr;
  }

  .stats-grid {
    grid-template-columns: 1fr;
  }

  .game-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.75rem;
  }

  .game-info {
    width: 100%;
  }

  .welcome-section .page-title {
    font-size: 2rem;
  }
}

@media (max-width: 480px) {
  .card-header {
    padding: 1.5rem 1rem 0.75rem;
  }

  .user-profile-section,
  .stats-content,
  .games-content {
    padding: 1rem;
  }

  .welcome-section .page-title {
    font-size: 1.75rem;
  }
}
</style>