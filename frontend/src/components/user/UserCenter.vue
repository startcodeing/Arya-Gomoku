<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-50 via-blue-50 to-indigo-100 relative overflow-hidden">
    <!-- 背景装饰 -->
    <div class="absolute inset-0 overflow-hidden pointer-events-none">
      <div class="absolute -top-40 -right-40 w-80 h-80 bg-gradient-to-br from-blue-400/20 to-purple-600/20 rounded-full blur-3xl"></div>
      <div class="absolute -bottom-40 -left-40 w-80 h-80 bg-gradient-to-tr from-indigo-400/20 to-pink-600/20 rounded-full blur-3xl"></div>
      <div class="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-96 h-96 bg-gradient-to-r from-cyan-400/10 to-blue-600/10 rounded-full blur-3xl"></div>
    </div>

    <!-- 浮动元素 -->
    <div class="absolute inset-0 overflow-hidden pointer-events-none">
      <div class="absolute top-20 left-20 w-2 h-2 bg-blue-400/60 rounded-full animate-pulse"></div>
      <div class="absolute top-40 right-32 w-1 h-1 bg-purple-400/60 rounded-full animate-ping"></div>
      <div class="absolute bottom-32 left-1/4 w-1.5 h-1.5 bg-indigo-400/60 rounded-full animate-pulse"></div>
      <div class="absolute bottom-20 right-20 w-2 h-2 bg-pink-400/60 rounded-full animate-ping"></div>
    </div>

    <!-- 导航栏 -->
    <nav class="relative bg-white/80 backdrop-blur-xl shadow-lg border-b border-white/20">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center">
            <router-link to="/" class="text-xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent hover:from-purple-600 hover:to-blue-600 transition-all duration-300">
              🎯 五子棋游戏
            </router-link>
          </div>
          
          <div class="flex items-center space-x-4">
            <div class="flex items-center space-x-3">
              <div class="relative">
                <div class="w-10 h-10 bg-gradient-to-r from-blue-500 to-purple-600 rounded-full flex items-center justify-center shadow-lg ring-2 ring-white/50">
                  <span class="text-white text-sm font-bold">
                    {{ userStore.userDisplayName.charAt(0).toUpperCase() }}
                  </span>
                </div>
                <div class="absolute -bottom-1 -right-1 w-4 h-4 bg-green-400 rounded-full border-2 border-white"></div>
              </div>
              <div>
                <span class="text-gray-800 font-medium">{{ userStore.userDisplayName }}</span>
                <div class="text-xs text-gray-500">在线</div>
              </div>
            </div>
            
            <button
              @click="handleLogout"
              class="group relative px-4 py-2 bg-gradient-to-r from-red-500 to-pink-600 text-white rounded-lg font-medium shadow-lg hover:shadow-xl transform hover:scale-105 transition-all duration-200 overflow-hidden"
            >
              <span class="relative z-10 flex items-center space-x-2">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"></path>
                </svg>
                <span>退出登录</span>
              </span>
              <div class="absolute inset-0 bg-gradient-to-r from-red-600 to-pink-700 opacity-0 group-hover:opacity-100 transition-opacity duration-200"></div>
            </button>
          </div>
        </div>
      </div>
    </nav>

    <div class="relative max-w-7xl mx-auto py-8 sm:px-6 lg:px-8">
      <!-- 页面标题 -->
      <div class="px-4 py-8 sm:px-0 text-center">
        <h1 class="text-4xl font-bold bg-gradient-to-r from-gray-900 via-blue-800 to-purple-800 bg-clip-text text-transparent mb-4">
          用户中心
        </h1>
        <p class="text-lg text-gray-600 max-w-2xl mx-auto">
          欢迎回来！管理您的游戏数据，查看统计信息，开始新的对战
        </p>
      </div>

      <!-- 统计卡片 -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-12">
        <div class="group relative bg-white/80 backdrop-blur-sm overflow-hidden shadow-xl rounded-2xl border border-white/20 hover:shadow-2xl transform hover:scale-105 transition-all duration-300">
          <div class="absolute inset-0 bg-gradient-to-br from-blue-500/10 to-purple-600/10 opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
          <div class="relative p-6">
            <div class="flex items-center">
              <div class="flex-shrink-0">
                <div class="w-12 h-12 bg-gradient-to-r from-blue-500 to-blue-600 rounded-xl flex items-center justify-center shadow-lg">
                  <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z"></path>
                  </svg>
                </div>
              </div>
              <div class="ml-5 w-0 flex-1">
                <dl>
                  <dt class="text-sm font-medium text-gray-500 truncate">
                    总游戏数
                  </dt>
                  <dd class="text-2xl font-bold text-gray-900">
                    {{ stats.totalGames }}
                  </dd>
                </dl>
              </div>
            </div>
          </div>
        </div>

        <div class="group relative bg-white/80 backdrop-blur-sm overflow-hidden shadow-xl rounded-2xl border border-white/20 hover:shadow-2xl transform hover:scale-105 transition-all duration-300">
          <div class="absolute inset-0 bg-gradient-to-br from-green-500/10 to-emerald-600/10 opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
          <div class="relative p-6">
            <div class="flex items-center">
              <div class="flex-shrink-0">
                <div class="w-12 h-12 bg-gradient-to-r from-green-500 to-emerald-600 rounded-xl flex items-center justify-center shadow-lg">
                  <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                  </svg>
                </div>
              </div>
              <div class="ml-5 w-0 flex-1">
                <dl>
                  <dt class="text-sm font-medium text-gray-500 truncate">
                    胜率
                  </dt>
                  <dd class="text-2xl font-bold text-gray-900">
                    {{ (stats.winRate * 100).toFixed(1) }}%
                  </dd>
                </dl>
              </div>
            </div>
          </div>
        </div>

        <div class="group relative bg-white/80 backdrop-blur-sm overflow-hidden shadow-xl rounded-2xl border border-white/20 hover:shadow-2xl transform hover:scale-105 transition-all duration-300">
          <div class="absolute inset-0 bg-gradient-to-br from-yellow-500/10 to-orange-600/10 opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
          <div class="relative p-6">
            <div class="flex items-center">
              <div class="flex-shrink-0">
                <div class="w-12 h-12 bg-gradient-to-r from-yellow-500 to-orange-600 rounded-xl flex items-center justify-center shadow-lg">
                  <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z"></path>
                  </svg>
                </div>
              </div>
              <div class="ml-5 w-0 flex-1">
                <dl>
                  <dt class="text-sm font-medium text-gray-500 truncate">
                    获胜场次
                  </dt>
                  <dd class="text-2xl font-bold text-gray-900">
                    {{ stats.winGames }}
                  </dd>
                </dl>
              </div>
            </div>
          </div>
        </div>

        <div class="group relative bg-white/80 backdrop-blur-sm overflow-hidden shadow-xl rounded-2xl border border-white/20 hover:shadow-2xl transform hover:scale-105 transition-all duration-300">
          <div class="absolute inset-0 bg-gradient-to-br from-purple-500/10 to-pink-600/10 opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
          <div class="relative p-6">
            <div class="flex items-center">
              <div class="flex-shrink-0">
                <div class="w-12 h-12 bg-gradient-to-r from-purple-500 to-pink-600 rounded-xl flex items-center justify-center shadow-lg">
                  <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                  </svg>
                </div>
              </div>
              <div class="ml-5 w-0 flex-1">
                <dl>
                  <dt class="text-sm font-medium text-gray-500 truncate">
                    平均时长
                  </dt>
                  <dd class="text-2xl font-bold text-gray-900">
                    {{ formatDuration(stats.avgDuration) }}
                  </dd>
                </dl>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 功能导航 -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
        <router-link
          to="/user/history"
          class="group relative bg-white/80 backdrop-blur-sm overflow-hidden shadow-xl rounded-2xl border border-white/20 hover:shadow-2xl transform hover:scale-105 transition-all duration-300"
        >
          <div class="absolute inset-0 bg-gradient-to-br from-blue-500/10 to-indigo-600/10 opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
          <div class="relative p-8">
            <div class="flex items-center">
              <div class="flex-shrink-0">
                <div class="w-14 h-14 bg-gradient-to-r from-blue-500 to-indigo-600 rounded-2xl flex items-center justify-center shadow-lg group-hover:scale-110 transition-transform duration-300">
                  <svg class="w-7 h-7 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
                  </svg>
                </div>
              </div>
              <div class="ml-6">
                <h3 class="text-xl font-bold text-gray-900 group-hover:text-blue-600 transition-colors duration-300">游戏历史</h3>
                <p class="text-sm text-gray-600 mt-1">查看您的游戏记录和详细数据</p>
              </div>
            </div>
          </div>
        </router-link>

        <router-link
          to="/user/statistics"
          class="group relative bg-white/80 backdrop-blur-sm overflow-hidden shadow-xl rounded-2xl border border-white/20 hover:shadow-2xl transform hover:scale-105 transition-all duration-300"
        >
          <div class="absolute inset-0 bg-gradient-to-br from-green-500/10 to-emerald-600/10 opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
          <div class="relative p-8">
            <div class="flex items-center">
              <div class="flex-shrink-0">
                <div class="w-14 h-14 bg-gradient-to-r from-green-500 to-emerald-600 rounded-2xl flex items-center justify-center shadow-lg group-hover:scale-110 transition-transform duration-300">
                  <svg class="w-7 h-7 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 12l3-3 3 3 4-4M8 21l4-4 4 4M3 4h18M4 4h16v12a1 1 0 01-1 1H5a1 1 0 01-1-1V4z"></path>
                  </svg>
                </div>
              </div>
              <div class="ml-6">
                <h3 class="text-xl font-bold text-gray-900 group-hover:text-green-600 transition-colors duration-300">数据统计</h3>
                <p class="text-sm text-gray-600 mt-1">分析您的游戏表现和趋势</p>
              </div>
            </div>
          </div>
        </router-link>

        <router-link
          to="/user/settings"
          class="group relative bg-white/80 backdrop-blur-sm overflow-hidden shadow-xl rounded-2xl border border-white/20 hover:shadow-2xl transform hover:scale-105 transition-all duration-300"
        >
          <div class="absolute inset-0 bg-gradient-to-br from-gray-500/10 to-slate-600/10 opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
          <div class="relative p-8">
            <div class="flex items-center">
              <div class="flex-shrink-0">
                <div class="w-14 h-14 bg-gradient-to-r from-gray-500 to-slate-600 rounded-2xl flex items-center justify-center shadow-lg group-hover:scale-110 transition-transform duration-300">
                  <svg class="w-7 h-7 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"></path>
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                  </svg>
                </div>
              </div>
              <div class="ml-6">
                <h3 class="text-xl font-bold text-gray-900 group-hover:text-gray-600 transition-colors duration-300">个人设置</h3>
                <p class="text-sm text-gray-600 mt-1">管理您的账户和偏好设置</p>
              </div>
            </div>
          </div>
        </router-link>

        <div class="group relative bg-white/80 backdrop-blur-sm overflow-hidden shadow-xl rounded-2xl border border-white/20 hover:shadow-2xl transform hover:scale-105 transition-all duration-300">
          <div class="absolute inset-0 bg-gradient-to-br from-indigo-500/10 to-purple-600/10 opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
          <div class="relative p-8">
            <div class="flex items-center">
              <div class="flex-shrink-0">
                <div class="w-14 h-14 bg-gradient-to-r from-indigo-500 to-purple-600 rounded-2xl flex items-center justify-center shadow-lg group-hover:scale-110 transition-transform duration-300">
                  <svg class="w-7 h-7 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path>
                  </svg>
                </div>
              </div>
              <div class="ml-6">
                <h3 class="text-xl font-bold text-gray-900 group-hover:text-indigo-600 transition-colors duration-300">快速开始</h3>
                <div class="mt-3 space-y-2">
                  <router-link
                    to="/ai-game"
                    class="block text-sm text-indigo-600 hover:text-indigo-800 font-medium transition-colors duration-200"
                  >
                    🤖 AI 对战
                  </router-link>
                  <router-link
                    to="/pvp"
                    class="block text-sm text-indigo-600 hover:text-indigo-800 font-medium transition-colors duration-200"
                  >
                    👥 双人对战
                  </router-link>
                  <router-link
                    to="/llm-battle"
                    class="block text-sm text-indigo-600 hover:text-indigo-800 font-medium transition-colors duration-200"
                  >
                    🧠 LLM 对战
                  </router-link>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="group relative bg-white/80 backdrop-blur-sm overflow-hidden shadow-xl rounded-2xl border border-white/20 hover:shadow-2xl transform hover:scale-105 transition-all duration-300">
          <div class="absolute inset-0 bg-gradient-to-br from-orange-500/10 to-red-600/10 opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
          <div class="relative p-8">
            <div class="flex items-center">
              <div class="flex-shrink-0">
                <div class="w-14 h-14 bg-gradient-to-r from-orange-500 to-red-600 rounded-2xl flex items-center justify-center shadow-lg group-hover:scale-110 transition-transform duration-300">
                  <svg class="w-7 h-7 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
                  </svg>
                </div>
              </div>
              <div class="ml-6">
                <h3 class="text-xl font-bold text-gray-900 group-hover:text-orange-600 transition-colors duration-300">最近活动</h3>
                <p class="text-sm text-gray-600 mt-1">
                  最后游戏时间：{{ formatDate(stats.lastGameTime) }}
                </p>
                <div class="mt-2 text-xs text-gray-500">
                  保持活跃，提升技能！
                </div>
              </div>
            </div>
          </div>
        </div>

        <router-link
          to="/user/data"
          class="group relative bg-white/80 backdrop-blur-sm overflow-hidden shadow-xl rounded-2xl border border-white/20 hover:shadow-2xl transform hover:scale-105 transition-all duration-300"
        >
          <div class="absolute inset-0 bg-gradient-to-br from-teal-500/10 to-cyan-600/10 opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
          <div class="relative p-8">
            <div class="flex items-center">
              <div class="flex-shrink-0">
                <div class="w-14 h-14 bg-gradient-to-r from-teal-500 to-cyan-600 rounded-2xl flex items-center justify-center shadow-lg group-hover:scale-110 transition-transform duration-300">
                  <svg class="w-7 h-7 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path>
                  </svg>
                </div>
              </div>
              <div class="ml-6">
                <h3 class="text-xl font-bold text-gray-900 group-hover:text-teal-600 transition-colors duration-300">数据管理</h3>
                <p class="text-sm text-gray-600 mt-1">导出数据、清理记录</p>
              </div>
            </div>
          </div>
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../../stores/user'
import { gameApi } from '../../services/gameApi'

const router = useRouter()
const userStore = useUserStore()

interface UserStats {
  totalGames: number
  winGames: number
  loseGames: number
  drawGames: number
  winRate: number
  avgDuration: number
  lastGameTime: string
}

const stats = ref<UserStats>({
  totalGames: 0,
  winGames: 0,
  loseGames: 0,
  drawGames: 0,
  winRate: 0,
  avgDuration: 0,
  lastGameTime: ''
})

const loading = ref(false)

const fetchUserStats = async () => {
  if (!userStore.isAuthenticated) return

  loading.value = true
  try {
    const userStats = await gameApi.getUserStats()
    stats.value = {
      totalGames: userStats.totalGames,
      winGames: userStats.winGames,
      loseGames: userStats.loseGames,
      drawGames: userStats.drawGames,
      winRate: userStats.winRate,
      avgDuration: userStats.avgDuration,
      lastGameTime: userStats.lastGameTime
    }
  } catch (error) {
    console.error('Failed to fetch user stats:', error)
  } finally {
    loading.value = false
  }
}

const handleLogout = async () => {
  try {
    await userStore.logout()
    router.push('/login')
  } catch (error) {
    console.error('Logout failed:', error)
  }
}

const formatDuration = (seconds: number): string => {
  if (seconds < 60) {
    return `${Math.round(seconds)}秒`
  } else if (seconds < 3600) {
    return `${Math.round(seconds / 60)}分钟`
  } else {
    return `${Math.round(seconds / 3600)}小时`
  }
}

const formatDate = (dateString: string): string => {
  if (!dateString) return '暂无'
  
  const date = new Date(dateString)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))
  
  if (diffDays === 0) {
    return '今天'
  } else if (diffDays === 1) {
    return '昨天'
  } else if (diffDays < 7) {
    return `${diffDays}天前`
  } else {
    return date.toLocaleDateString('zh-CN')
  }
}

onMounted(() => {
  fetchUserStats()
})
</script>