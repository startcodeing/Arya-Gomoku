<template>
  <div class="home-container">
    <!-- 用户导航栏 -->
    <div class="user-nav">
      <div v-if="userStore.isAuthenticated" class="user-info">
        <span class="welcome-text">欢迎, {{ userStore.userDisplayName }}</span>
        <button @click="goToProfile" class="profile-btn">个人中心</button>
        <button @click="logout" class="logout-btn">退出</button>
      </div>
      <div v-else class="auth-buttons">
        <button @click="goToLogin" class="login-btn">登录</button>
        <button @click="goToRegister" class="register-btn">注册</button>
      </div>
    </div>

    <div class="game-header">
      <h1>五子棋游戏</h1>
      <p class="subtitle">Gomoku Game</p>
      
      <!-- 用户统计信息 -->
      <div v-if="userStore.isAuthenticated && userStats" class="user-stats">
        <div class="stats-container">
          <div class="stat-item">
            <span class="stat-number">{{ userStats.totalGames }}</span>
            <span class="stat-label">总局数</span>
          </div>
          <div class="stat-item">
            <span class="stat-number">{{ userStats.winRate }}%</span>
            <span class="stat-label">胜率</span>
          </div>
          <div class="stat-item">
            <span class="stat-number">{{ userStats.wins }}</span>
            <span class="stat-label">胜利</span>
          </div>
          <div class="stat-item">
            <span class="stat-number">{{ userStats.losses }}</span>
            <span class="stat-label">失败</span>
          </div>
        </div>
      </div>
    </div>

    <div class="game-modes">
      <div class="mode-card" @click="goToAIGame">
        <div class="mode-icon">🤖</div>
        <h3>人机对战</h3>
        <p>与传统AI进行五子棋对战，支持4个难度级别</p>
        <div class="ai-features">
          <span class="feature-tag">初级</span>
          <span class="feature-tag">中级</span>
          <span class="feature-tag">高级</span>
          <span class="feature-tag">专家</span>
        </div>
        <button class="mode-button">开始游戏</button>
      </div>

      <div class="mode-card llm-card" @click="goToLLMBattle">
        <div class="mode-icon">🧠</div>
        <h3>与AI对战</h3>
        <p>与先进的大语言模型AI对战</p>
        <div class="llm-features">
          <span class="feature-tag">DeepSeek</span>
          <span class="feature-tag">ChatGPT</span>
          <span class="feature-tag">Ollama</span>
        </div>
        <button class="mode-button llm-button">挑战AI</button>
      </div>

      <div class="mode-card" @click="goToPVP">
        <div class="mode-icon">👥</div>
        <h3>双人对战</h3>
        <p>与其他玩家在线对战</p>
        <button class="mode-button">进入房间</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { ref, onMounted } from 'vue'
import { gameApi } from '../services/gameApi'

const router = useRouter()
const userStore = useUserStore()

// 用户统计数据
const userStats = ref<any>(null)
const loadingStats = ref(false)

// 加载用户统计数据
async function loadUserStats() {
  if (!userStore.isAuthenticated) return
  
  try {
    loadingStats.value = true
    userStats.value = await gameApi.getUserStats()
  } catch (error) {
    console.error('加载用户统计失败:', error)
  } finally {
    loadingStats.value = false
  }
}

onMounted(() => {
  loadUserStats()
})

// 游戏导航方法
function goToAIGame() {
  router.push('/ai-game')
}

function goToPVP() {
  router.push('/pvp')
}

function goToLLMBattle() {
  router.push('/llm-battle')
}

// 用户导航方法
function goToLogin() {
  router.push('/login')
}

function goToRegister() {
  router.push('/register')
}

function goToProfile() {
  router.push('/profile')
}

async function logout() {
  try {
    await userStore.logout()
    // 退出后刷新页面状态
  } catch (error) {
    console.error('退出登录失败:', error)
  }
}
</script>

<style scoped>
.home-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.user-nav {
  position: absolute;
  top: 20px;
  right: 20px;
  z-index: 1000;
  max-width: 350px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
  background: rgba(0, 0, 0, 0.7);
  padding: 14px 26px;
  border-radius: 35px;
  backdrop-filter: blur(20px);
  border: 2px solid rgba(255, 255, 255, 0.4);
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.3);
  transition: all 0.3s ease;
}

.user-info:hover {
  background: rgba(0, 0, 0, 0.8);
  transform: translateY(-3px);
  box-shadow: 0 15px 50px rgba(0, 0, 0, 0.4);
  border-color: rgba(255, 255, 255, 0.6);
}

.welcome-text {
  color: white;
  font-weight: 700;
  font-size: 1rem;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.8);
  letter-spacing: 0.5px;
  white-space: nowrap;
}

.auth-buttons {
  display: flex;
  gap: 12px;
}

.profile-btn,
.logout-btn,
.login-btn,
.register-btn {
  padding: 10px 20px;
  border: none;
  border-radius: 25px;
  cursor: pointer;
  font-weight: 600;
  font-size: 0.9rem;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
  letter-spacing: 0.3px;
}

.profile-btn {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.9), rgba(255, 255, 255, 0.8));
  color: #333;
  border: 2px solid rgba(255, 255, 255, 0.9);
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
  font-weight: 700;
  text-shadow: none;
}

.profile-btn:hover {
  background: linear-gradient(135deg, rgba(255, 255, 255, 1), rgba(255, 255, 255, 0.95));
  color: #222;
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.3);
  border-color: rgba(255, 255, 255, 1);
}

.profile-btn:active {
  transform: translateY(0);
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
}

.logout-btn {
  background: linear-gradient(135deg, #ff6b6b, #ee5a52);
  color: white;
  box-shadow: 0 4px 15px rgba(238, 90, 82, 0.3);
}

.logout-btn:hover {
  background: linear-gradient(135deg, #ff5252, #e53935);
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(238, 90, 82, 0.4);
}

.logout-btn:active {
  transform: translateY(0);
  box-shadow: 0 4px 15px rgba(238, 90, 82, 0.3);
}

.login-btn {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.25), rgba(255, 255, 255, 0.15));
  color: white;
  border: 1px solid rgba(255, 255, 255, 0.3);
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
}

.login-btn:hover {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.35), rgba(255, 255, 255, 0.25));
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
}

.login-btn:active {
  transform: translateY(0);
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
}

.register-btn {
  background: linear-gradient(135deg, #4caf50, #45a049);
  color: white;
  box-shadow: 0 4px 15px rgba(76, 175, 80, 0.3);
}

.register-btn:hover {
  background: linear-gradient(135deg, #45a049, #388e3c);
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(76, 175, 80, 0.4);
}

.register-btn:active {
  transform: translateY(0);
  box-shadow: 0 4px 15px rgba(76, 175, 80, 0.3);
}

.game-header {
  text-align: center;
  margin-bottom: 40px;
  color: white;
}

.game-header h1 {
  font-size: 3rem;
  margin-bottom: 10px;
  text-shadow: 2px 2px 4px rgba(0,0,0,0.3);
}

.subtitle {
  font-size: 1.2rem;
  opacity: 0.9;
}

.user-stats {
  margin-top: 40px;
  padding: 0 20px;
  max-width: calc(100% - 400px);
  margin-left: auto;
  margin-right: auto;
}

.stats-container {
  display: flex;
  justify-content: center;
  gap: 30px;
  flex-wrap: wrap;
  margin-top: 20px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  background: rgba(255, 255, 255, 0.1);
  padding: 15px 20px;
  border-radius: 15px;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  min-width: 80px;
}

.stat-number {
  font-size: 1.8rem;
  font-weight: bold;
  color: white;
  margin-bottom: 5px;
}

.stat-label {
  font-size: 0.9rem;
  color: rgba(255, 255, 255, 0.8);
  font-weight: 500;
}

.game-modes {
  display: flex;
  justify-content: center;
  gap: 40px;
  margin-bottom: 40px;
  flex-wrap: wrap;
}

.mode-card {
  background: white;
  border-radius: 20px;
  padding: 40px 30px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 10px 30px rgba(0,0,0,0.2);
  min-width: 250px;
}

.mode-card:hover {
  transform: translateY(-10px);
  box-shadow: 0 20px 40px rgba(0,0,0,0.3);
}

.mode-icon {
  font-size: 4rem;
  margin-bottom: 20px;
}

.mode-card h3 {
  font-size: 1.5rem;
  margin-bottom: 15px;
  color: #333;
}

.mode-card p {
  color: #666;
  margin-bottom: 25px;
  line-height: 1.5;
}

.mode-button {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 12px 30px;
  border-radius: 25px;
  font-size: 1rem;
  cursor: pointer;
  transition: all 0.3s ease;
}

.mode-button:hover {
  transform: scale(1.05);
  box-shadow: 0 5px 15px rgba(0,0,0,0.2);
}

.llm-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: 2px solid #ffd700;
  position: relative;
  overflow: hidden;
}

.llm-card::before {
  content: '';
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: linear-gradient(45deg, transparent, rgba(255, 255, 255, 0.1), transparent);
  transform: rotate(45deg);
  animation: shimmer 3s infinite;
}

.llm-card h3 {
  color: white;
}

.llm-card p {
  color: rgba(255, 255, 255, 0.9);
}

.llm-features, .ai-features {
  display: flex;
  gap: 8px;
  justify-content: center;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.feature-tag {
  background: rgba(255, 255, 255, 0.2);
  color: white;
  padding: 4px 12px;
  border-radius: 15px;
  font-size: 0.8rem;
  font-weight: 500;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.llm-button {
  background: linear-gradient(135deg, #ffd700 0%, #ffed4e 100%);
  color: #333;
  font-weight: 600;
  border: none;
  box-shadow: 0 4px 15px rgba(255, 215, 0, 0.3);
}

.llm-button:hover {
  background: linear-gradient(135deg, #ffed4e 0%, #ffd700 100%);
  transform: scale(1.05);
  box-shadow: 0 6px 20px rgba(255, 215, 0, 0.4);
}

@keyframes shimmer {
  0% {
    transform: translateX(-100%) translateY(-100%) rotate(45deg);
  }
  100% {
    transform: translateX(100%) translateY(100%) rotate(45deg);
  }
}

@media (max-width: 768px) {
  .game-modes {
    flex-direction: column;
    align-items: center;
  }
  
  .user-nav {
    position: relative;
    top: 0;
    right: 0;
    margin-bottom: 20px;
    max-width: 100%;
    display: flex;
    justify-content: center;
  }
  
  .user-stats {
    max-width: 100%;
    padding: 0 10px;
  }
  
  .stats-container {
    gap: 15px;
  }
  
  .welcome-text {
    font-size: 0.9rem;
  }
  
  .profile-btn,
  .logout-btn {
    padding: 8px 16px;
    font-size: 0.85rem;
  }
}
</style>