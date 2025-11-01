<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 导航栏 -->
    <nav class="bg-white shadow-sm border-b">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center space-x-4">
            <router-link to="/user" class="text-gray-500 hover:text-gray-700">
              ← 返回用户中心
            </router-link>
            <h1 class="text-xl font-semibold text-gray-900">数据统计</h1>
          </div>
        </div>
      </div>
    </nav>

    <div class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <!-- 总体统计 -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <div class="bg-white overflow-hidden shadow rounded-lg">
          <div class="p-5">
            <div class="flex items-center">
              <div class="flex-shrink-0">
                <div class="w-8 h-8 bg-blue-500 rounded-md flex items-center justify-center">
                  <span class="text-white text-sm">🎮</span>
                </div>
              </div>
              <div class="ml-5 w-0 flex-1">
                <dl>
                  <dt class="text-sm font-medium text-gray-500 truncate">总游戏数</dt>
                  <dd class="text-lg font-medium text-gray-900">{{ stats.totalGames }}</dd>
                </dl>
              </div>
            </div>
          </div>
        </div>

        <div class="bg-white overflow-hidden shadow rounded-lg">
          <div class="p-5">
            <div class="flex items-center">
              <div class="flex-shrink-0">
                <div class="w-8 h-8 bg-green-500 rounded-md flex items-center justify-center">
                  <span class="text-white text-sm">🏆</span>
                </div>
              </div>
              <div class="ml-5 w-0 flex-1">
                <dl>
                  <dt class="text-sm font-medium text-gray-500 truncate">胜率</dt>
                  <dd class="text-lg font-medium text-gray-900">{{ (stats.winRate * 100).toFixed(1) }}%</dd>
                </dl>
              </div>
            </div>
          </div>
        </div>

        <div class="bg-white overflow-hidden shadow rounded-lg">
          <div class="p-5">
            <div class="flex items-center">
              <div class="flex-shrink-0">
                <div class="w-8 h-8 bg-yellow-500 rounded-md flex items-center justify-center">
                  <span class="text-white text-sm">⭐</span>
                </div>
              </div>
              <div class="ml-5 w-0 flex-1">
                <dl>
                  <dt class="text-sm font-medium text-gray-500 truncate">获胜场次</dt>
                  <dd class="text-lg font-medium text-gray-900">{{ stats.winGames }}</dd>
                </dl>
              </div>
            </div>
          </div>
        </div>

        <div class="bg-white overflow-hidden shadow rounded-lg">
          <div class="p-5">
            <div class="flex items-center">
              <div class="flex-shrink-0">
                <div class="w-8 h-8 bg-purple-500 rounded-md flex items-center justify-center">
                  <span class="text-white text-sm">⏱️</span>
                </div>
              </div>
              <div class="ml-5 w-0 flex-1">
                <dl>
                  <dt class="text-sm font-medium text-gray-500 truncate">平均时长</dt>
                  <dd class="text-lg font-medium text-gray-900">{{ formatDuration(stats.avgDuration) }}</dd>
                </dl>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 游戏类型统计 -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
        <div class="bg-white shadow rounded-lg">
          <div class="px-6 py-4 border-b border-gray-200">
            <h3 class="text-lg font-medium text-gray-900">游戏类型分布</h3>
          </div>
          <div class="p-6">
            <div class="space-y-4">
              <div class="flex items-center justify-between">
                <div class="flex items-center">
                  <div class="w-4 h-4 bg-blue-500 rounded mr-3"></div>
                  <span class="text-sm text-gray-700">AI对战</span>
                </div>
                <div class="text-right">
                  <div class="text-sm font-medium text-gray-900">{{ stats.aiGames }}</div>
                  <div class="text-xs text-gray-500">{{ getPercentage(stats.aiGames, stats.totalGames) }}%</div>
                </div>
              </div>
              
              <div class="flex items-center justify-between">
                <div class="flex items-center">
                  <div class="w-4 h-4 bg-purple-500 rounded mr-3"></div>
                  <span class="text-sm text-gray-700">LLM对战</span>
                </div>
                <div class="text-right">
                  <div class="text-sm font-medium text-gray-900">{{ stats.llmGames }}</div>
                  <div class="text-xs text-gray-500">{{ getPercentage(stats.llmGames, stats.totalGames) }}%</div>
                </div>
              </div>
              
              <div class="flex items-center justify-between">
                <div class="flex items-center">
                  <div class="w-4 h-4 bg-green-500 rounded mr-3"></div>
                  <span class="text-sm text-gray-700">双人对战</span>
                </div>
                <div class="text-right">
                  <div class="text-sm font-medium text-gray-900">{{ stats.pvpGames }}</div>
                  <div class="text-xs text-gray-500">{{ getPercentage(stats.pvpGames, stats.totalGames) }}%</div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="bg-white shadow rounded-lg">
          <div class="px-6 py-4 border-b border-gray-200">
            <h3 class="text-lg font-medium text-gray-900">胜负统计</h3>
          </div>
          <div class="p-6">
            <div class="space-y-4">
              <div class="flex items-center justify-between">
                <div class="flex items-center">
                  <div class="w-4 h-4 bg-green-500 rounded mr-3"></div>
                  <span class="text-sm text-gray-700">胜利</span>
                </div>
                <div class="text-right">
                  <div class="text-sm font-medium text-gray-900">{{ stats.winGames }}</div>
                  <div class="text-xs text-gray-500">{{ getPercentage(stats.winGames, stats.totalGames) }}%</div>
                </div>
              </div>
              
              <div class="flex items-center justify-between">
                <div class="flex items-center">
                  <div class="w-4 h-4 bg-red-500 rounded mr-3"></div>
                  <span class="text-sm text-gray-700">失败</span>
                </div>
                <div class="text-right">
                  <div class="text-sm font-medium text-gray-900">{{ stats.loseGames }}</div>
                  <div class="text-xs text-gray-500">{{ getPercentage(stats.loseGames, stats.totalGames) }}%</div>
                </div>
              </div>
              
              <div class="flex items-center justify-between">
                <div class="flex items-center">
                  <div class="w-4 h-4 bg-yellow-500 rounded mr-3"></div>
                  <span class="text-sm text-gray-700">平局</span>
                </div>
                <div class="text-right">
                  <div class="text-sm font-medium text-gray-900">{{ stats.drawGames }}</div>
                  <div class="text-xs text-gray-500">{{ getPercentage(stats.drawGames, stats.totalGames) }}%</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 排行榜 -->
      <div class="bg-white shadow rounded-lg">
        <div class="px-6 py-4 border-b border-gray-200">
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-medium text-gray-900">排行榜</h3>
            <select
              v-model="selectedRankingType"
              @change="fetchRankings"
              class="border border-gray-300 rounded-md px-3 py-1 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
            >
              <option value="">全部</option>
              <option value="ai">AI对战</option>
              <option value="llm">LLM对战</option>
              <option value="pvp">双人对战</option>
            </select>
          </div>
        </div>

        <div v-if="loadingRankings" class="p-6 text-center">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600 mx-auto"></div>
          <p class="mt-2 text-gray-500">加载中...</p>
        </div>

        <div v-else-if="rankings.length === 0" class="p-6 text-center text-gray-500">
          暂无排行榜数据
        </div>

        <div v-else class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  排名
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  用户
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  胜场
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  胜率
                </th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr
                v-for="player in rankings"
                :key="player.userId"
                :class="{ 'bg-indigo-50': player.userId === userStore.user?.id }"
              >
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="flex items-center">
                    <span
                      v-if="player.rank <= 3"
                      :class="getRankClass(player.rank)"
                      class="w-6 h-6 rounded-full flex items-center justify-center text-white text-xs font-bold mr-2"
                    >
                      {{ player.rank }}
                    </span>
                    <span v-else class="text-sm text-gray-900">{{ player.rank }}</span>
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="flex items-center">
                    <div class="w-8 h-8 bg-gray-300 rounded-full flex items-center justify-center mr-3">
                      <span class="text-white text-xs font-medium">
                        {{ (player.nickname || player.username).charAt(0).toUpperCase() }}
                      </span>
                    </div>
                    <div>
                      <div class="text-sm font-medium text-gray-900">
                        {{ player.nickname || player.username }}
                        <span v-if="player.userId === userStore.user?.id" class="text-indigo-600">(我)</span>
                      </div>
                      <div class="text-sm text-gray-500">@{{ player.username }}</div>
                    </div>
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {{ player.winGames }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {{ (player.winRate * 100).toFixed(1) }}%
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useUserStore } from '../../stores/user'
import { gameApi, type UserGameStats, type PlayerRanking } from '../../services/gameApi'

const userStore = useUserStore()

const stats = ref<UserGameStats>({
  userId: '',
  totalGames: 0,
  winGames: 0,
  loseGames: 0,
  drawGames: 0,
  winRate: 0,
  aiGames: 0,
  llmGames: 0,
  pvpGames: 0,
  avgDuration: 0,
  lastGameTime: ''
})

const rankings = ref<PlayerRanking[]>([])
const selectedRankingType = ref('')
const loadingRankings = ref(false)

const fetchUserStats = async () => {
  if (!userStore.user?.id) return

  try {
    const userStats = await gameApi.getUserStats(userStore.user.id)
    stats.value = userStats
  } catch (error) {
    console.error('Failed to fetch user stats:', error)
  }
}

const fetchRankings = async () => {
  loadingRankings.value = true
  try {
    const data = await gameApi.getTopPlayers(selectedRankingType.value, 20)
    rankings.value = data
  } catch (error) {
    console.error('Failed to fetch rankings:', error)
  } finally {
    loadingRankings.value = false
  }
}

const getPercentage = (value: number, total: number): string => {
  if (total === 0) return '0'
  return ((value / total) * 100).toFixed(1)
}

const getRankClass = (rank: number): string => {
  const classes: Record<number, string> = {
    1: 'bg-yellow-500',
    2: 'bg-gray-400',
    3: 'bg-yellow-600'
  }
  return classes[rank] || 'bg-gray-300'
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

onMounted(() => {
  fetchUserStats()
  fetchRankings()
})
</script>