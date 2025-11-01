<template>
  <div class="min-h-screen bg-gray-50 py-8">
    <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
      <!-- 页面标题 -->
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900">用户中心</h1>
        <p class="mt-2 text-gray-600">管理您的个人信息和游戏数据</p>
      </div>

      <!-- 用户信息卡片 -->
      <div class="bg-white shadow rounded-lg mb-8">
        <div class="px-6 py-4 border-b border-gray-200">
          <h2 class="text-lg font-medium text-gray-900">个人信息</h2>
        </div>
        <div class="p-6">
          <div class="flex items-center space-x-6 mb-6">
            <div class="flex-shrink-0">
              <div class="h-20 w-20 rounded-full bg-indigo-100 flex items-center justify-center">
                <span class="text-2xl font-bold text-indigo-600">
                  {{ userStore.userDisplayName.charAt(0).toUpperCase() }}
                </span>
              </div>
            </div>
            <div class="flex-1">
              <h3 class="text-xl font-bold text-gray-900">{{ userStore.userDisplayName }}</h3>
              <p class="text-gray-600">@{{ userStore.user?.username }}</p>
              <p class="text-sm text-gray-500">{{ userStore.user?.email }}</p>
            </div>
            <div class="flex-shrink-0">
              <button
                @click="showEditProfile = true"
                class="bg-indigo-600 text-white px-4 py-2 rounded-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500"
              >
                编辑资料
              </button>
            </div>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <label class="block text-sm font-medium text-gray-700">用户名</label>
              <p class="mt-1 text-sm text-gray-900">{{ userStore.user?.username }}</p>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700">昵称</label>
              <p class="mt-1 text-sm text-gray-900">{{ userStore.user?.nickname }}</p>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700">邮箱</label>
              <p class="mt-1 text-sm text-gray-900">{{ userStore.user?.email }}</p>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700">注册时间</label>
              <p class="mt-1 text-sm text-gray-900">{{ formatDate(userStore.user?.createdAt) }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- 游戏统计卡片 -->
      <div class="bg-white shadow rounded-lg mb-8">
        <div class="px-6 py-4 border-b border-gray-200">
          <h2 class="text-lg font-medium text-gray-900">游戏统计</h2>
        </div>
        <div class="p-6">
          <div v-if="statsLoading" class="flex justify-center py-8">
            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
          </div>
          <div v-else-if="gameStats" class="grid grid-cols-1 md:grid-cols-3 gap-6">
            <!-- AI游戏统计 -->
            <div class="bg-blue-50 rounded-lg p-4">
              <h3 class="text-lg font-semibold text-blue-900 mb-3">AI对战</h3>
              <div class="space-y-2">
                <div class="flex justify-between">
                  <span class="text-blue-700">总场次</span>
                  <span class="font-medium text-blue-900">{{ gameStats.ai.totalGames }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-blue-700">胜利</span>
                  <span class="font-medium text-green-600">{{ gameStats.ai.wins }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-blue-700">失败</span>
                  <span class="font-medium text-red-600">{{ gameStats.ai.losses }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-blue-700">胜率</span>
                  <span class="font-medium text-blue-900">{{ gameStats.ai.winRate }}%</span>
                </div>
              </div>
            </div>

            <!-- LLM游戏统计 -->
            <div class="bg-green-50 rounded-lg p-4">
              <h3 class="text-lg font-semibold text-green-900 mb-3">LLM对战</h3>
              <div class="space-y-2">
                <div class="flex justify-between">
                  <span class="text-green-700">总场次</span>
                  <span class="font-medium text-green-900">{{ gameStats.llm.totalGames }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-green-700">胜利</span>
                  <span class="font-medium text-green-600">{{ gameStats.llm.wins }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-green-700">失败</span>
                  <span class="font-medium text-red-600">{{ gameStats.llm.losses }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-green-700">胜率</span>
                  <span class="font-medium text-green-900">{{ gameStats.llm.winRate }}%</span>
                </div>
              </div>
            </div>

            <!-- PVP游戏统计 -->
            <div class="bg-purple-50 rounded-lg p-4">
              <h3 class="text-lg font-semibold text-purple-900 mb-3">玩家对战</h3>
              <div class="space-y-2">
                <div class="flex justify-between">
                  <span class="text-purple-700">总场次</span>
                  <span class="font-medium text-purple-900">{{ gameStats.pvp.totalGames }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-purple-700">胜利</span>
                  <span class="font-medium text-green-600">{{ gameStats.pvp.wins }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-purple-700">失败</span>
                  <span class="font-medium text-red-600">{{ gameStats.pvp.losses }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-purple-700">胜率</span>
                  <span class="font-medium text-purple-900">{{ gameStats.pvp.winRate }}%</span>
                </div>
              </div>
            </div>
          </div>
          <div v-else class="text-center py-8 text-gray-500">
            暂无游戏统计数据
          </div>
        </div>
      </div>

      <!-- 游戏历史 -->
      <div class="bg-white shadow rounded-lg mb-8">
        <div class="px-6 py-4 border-b border-gray-200 flex justify-between items-center">
          <h2 class="text-lg font-medium text-gray-900">最近游戏</h2>
          <button
            @click="loadGameHistory"
            class="text-indigo-600 hover:text-indigo-500 text-sm font-medium"
          >
            刷新
          </button>
        </div>
        <div class="p-6">
          <div v-if="historyLoading" class="flex justify-center py-8">
            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
          </div>
          <div v-else-if="gameHistory.length > 0" class="space-y-4">
            <div
              v-for="game in gameHistory"
              :key="game.id"
              class="border rounded-lg p-4 hover:bg-gray-50"
            >
              <div class="flex justify-between items-start">
                <div class="flex-1">
                  <div class="flex items-center space-x-2">
                    <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                          :class="getGameTypeClass(game.gameType)">
                      {{ getGameTypeName(game.gameType) }}
                    </span>
                    <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                          :class="getResultClass(game.result)">
                      {{ getResultText(game.result) }}
                    </span>
                  </div>
                  <p class="text-sm text-gray-600 mt-1">
                    游戏时长: {{ formatDuration(game.duration) }}
                  </p>
                </div>
                <div class="text-right">
                  <p class="text-sm text-gray-500">{{ formatDate(game.createdAt) }}</p>
                </div>
              </div>
            </div>
          </div>
          <div v-else class="text-center py-8 text-gray-500">
            暂无游戏历史记录
          </div>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="flex space-x-4">
        <button
          @click="showChangePassword = true"
          class="bg-yellow-600 text-white px-6 py-2 rounded-md hover:bg-yellow-700 focus:outline-none focus:ring-2 focus:ring-yellow-500"
        >
          修改密码
        </button>
        <button
          @click="handleLogout"
          class="bg-red-600 text-white px-6 py-2 rounded-md hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-500"
        >
          退出登录
        </button>
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

// 获取结果文本
const getResultText = (result: string) => {
  const resultMap: Record<string, string> = {
    'win': '胜利',
    'loss': '失败',
    'draw': '平局'
  }
  return resultMap[result] || result
}

// 获取结果样式
const getResultClass = (result: string) => {
  const classMap: Record<string, string> = {
    'win': 'bg-green-100 text-green-800',
    'loss': 'bg-red-100 text-red-800',
    'draw': 'bg-yellow-100 text-yellow-800'
  }
  return classMap[result] || 'bg-gray-100 text-gray-800'
}

// 加载游戏统计
const loadGameStats = async () => {
  statsLoading.value = true
  try {
    gameStats.value = await statisticsApi.getUserStats()
  } catch (error) {
    console.error('Failed to load game stats:', error)
  } finally {
    statsLoading.value = false
  }
}

// 加载游戏历史
const loadGameHistory = async () => {
  historyLoading.value = true
  try {
    const response = await statisticsApi.getUserGameHistory(1, 10)
    gameHistory.value = response.games
  } catch (error) {
    console.error('Failed to load game history:', error)
  } finally {
    historyLoading.value = false
  }
}

// 处理资料更新
const handleProfileUpdated = () => {
  showEditProfile.value = false
  // 刷新用户信息
  userStore.fetchUserProfile().catch(console.error)
}

// 处理退出登录
const handleLogout = async () => {
  if (confirm('确定要退出登录吗？')) {
    try {
      await userStore.logout()
      router.push('/login')
    } catch (error) {
      console.error('Logout failed:', error)
    }
  }
}

// 组件挂载时加载数据
onMounted(() => {
  loadGameStats()
  loadGameHistory()
})
</script>