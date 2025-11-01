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
            <h1 class="text-xl font-semibold text-gray-900">游戏历史</h1>
          </div>
        </div>
      </div>
    </nav>

    <div class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <!-- 筛选器 -->
      <div class="bg-white shadow rounded-lg mb-6">
        <div class="px-6 py-4">
          <div class="flex flex-wrap gap-4 items-center">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">游戏类型</label>
              <select
                v-model="filters.gameType"
                @change="fetchGameHistory"
                class="border border-gray-300 rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
              >
                <option value="">全部</option>
                <option value="ai">AI对战</option>
                <option value="llm">LLM对战</option>
                <option value="pvp">双人对战</option>
              </select>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">游戏状态</label>
              <select
                v-model="filters.status"
                @change="fetchGameHistory"
                class="border border-gray-300 rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
              >
                <option value="">全部</option>
                <option value="completed">已完成</option>
                <option value="abandoned">已放弃</option>
              </select>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">搜索</label>
              <input
                v-model="searchKeyword"
                @input="debounceSearch"
                type="text"
                placeholder="搜索游戏..."
                class="border border-gray-300 rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
              />
            </div>

            <div class="flex-1"></div>
            
            <button
              @click="exportData"
              :disabled="loading"
              class="bg-green-600 text-white px-4 py-2 rounded-md text-sm hover:bg-green-700 disabled:opacity-50"
            >
              导出数据
            </button>
          </div>
        </div>
      </div>

      <!-- 游戏列表 -->
      <div class="bg-white shadow rounded-lg">
        <div class="px-6 py-4 border-b border-gray-200">
          <h2 class="text-lg font-medium text-gray-900">游戏记录</h2>
        </div>

        <div v-if="loading" class="p-6 text-center">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600 mx-auto"></div>
          <p class="mt-2 text-gray-500">加载中...</p>
        </div>

        <div v-else-if="gameHistory.games.length === 0" class="p-6 text-center text-gray-500">
          暂无游戏记录
        </div>

        <div v-else class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  游戏类型
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  状态
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  结果
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  时长
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  开始时间
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  操作
                </th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-for="game in gameHistory.games" :key="game.id" class="hover:bg-gray-50">
                <td class="px-6 py-4 whitespace-nowrap">
                  <span :class="getGameTypeClass(game.type)" class="px-2 py-1 text-xs font-medium rounded-full">
                    {{ getGameTypeName(game.type) }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span :class="getStatusClass(game.status)" class="px-2 py-1 text-xs font-medium rounded-full">
                    {{ getStatusName(game.status) }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span v-if="game.winner" :class="getResultClass(game.winner)" class="px-2 py-1 text-xs font-medium rounded-full">
                    {{ getResultName(game.winner) }}
                  </span>
                  <span v-else class="text-gray-400">-</span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {{ formatDuration(game.duration) }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {{ formatDate(game.createdAt) }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                  <button
                    @click="viewGameDetail(game)"
                    class="text-indigo-600 hover:text-indigo-900 mr-3"
                  >
                    查看详情
                  </button>
                  <button
                    @click="deleteGame(game)"
                    class="text-red-600 hover:text-red-900"
                  >
                    删除
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- 分页 -->
        <div v-if="gameHistory.total > 0" class="px-6 py-4 border-t border-gray-200">
          <div class="flex items-center justify-between">
            <div class="text-sm text-gray-700">
              显示 {{ (currentPage - 1) * pageSize + 1 }} 到 {{ Math.min(currentPage * pageSize, gameHistory.total) }} 条，
              共 {{ gameHistory.total }} 条记录
            </div>
            <div class="flex space-x-2">
              <button
                @click="changePage(currentPage - 1)"
                :disabled="currentPage <= 1"
                class="px-3 py-1 border border-gray-300 rounded-md text-sm disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
              >
                上一页
              </button>
              <button
                @click="changePage(currentPage + 1)"
                :disabled="currentPage >= totalPages"
                class="px-3 py-1 border border-gray-300 rounded-md text-sm disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
              >
                下一页
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useUserStore } from '../../stores/user'
import { gameApi, type GameHistory, type GameSummary } from '../../services/gameApi'

const userStore = useUserStore()

const loading = ref(false)
const searchKeyword = ref('')
const currentPage = ref(1)
const pageSize = ref(20)

const filters = reactive({
  gameType: '',
  status: ''
})

const gameHistory = ref<GameHistory>({
  games: [],
  total: 0,
  page: 1,
  pageSize: 20
})

const totalPages = computed(() => Math.ceil(gameHistory.value.total / pageSize.value))

let searchTimeout: NodeJS.Timeout

const debounceSearch = () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    if (searchKeyword.value.trim()) {
      searchGames()
    } else {
      fetchGameHistory()
    }
  }, 500)
}

const fetchGameHistory = async () => {
  if (!userStore.user?.id) return

  loading.value = true
  try {
    const data = await gameApi.getUserGameHistory(
      userStore.user.id,
      currentPage.value,
      pageSize.value,
      filters.gameType || undefined
    )
    gameHistory.value = data
  } catch (error) {
    console.error('Failed to fetch game history:', error)
  } finally {
    loading.value = false
  }
}

const searchGames = async () => {
  if (!searchKeyword.value.trim()) return

  loading.value = true
  try {
    const data = await gameApi.searchGames(
      searchKeyword.value,
      currentPage.value,
      pageSize.value
    )
    gameHistory.value = data
  } catch (error) {
    console.error('Failed to search games:', error)
  } finally {
    loading.value = false
  }
}

const changePage = (page: number) => {
  currentPage.value = page
  if (searchKeyword.value.trim()) {
    searchGames()
  } else {
    fetchGameHistory()
  }
}

const exportData = async () => {
  if (!userStore.user?.id) return

  try {
    const blob = await gameApi.exportUserData(userStore.user.id, 'json')
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `game-data-${new Date().toISOString().split('T')[0]}.json`
    document.body.appendChild(a)
    a.click()
    window.URL.revokeObjectURL(url)
    document.body.removeChild(a)
  } catch (error) {
    console.error('Failed to export data:', error)
  }
}

const viewGameDetail = (game: GameSummary) => {
  // TODO: 实现游戏详情查看
  console.log('View game detail:', game)
}

const deleteGame = async (game: GameSummary) => {
  if (!confirm('确定要删除这条游戏记录吗？')) return

  try {
    await gameApi.deleteGame(game.id, game.type as 'ai' | 'llm' | 'pvp')
    fetchGameHistory()
  } catch (error) {
    console.error('Failed to delete game:', error)
  }
}

const getGameTypeName = (type: string) => {
  const names: Record<string, string> = {
    ai: 'AI对战',
    llm: 'LLM对战',
    pvp: '双人对战'
  }
  return names[type] || type
}

const getGameTypeClass = (type: string) => {
  const classes: Record<string, string> = {
    ai: 'bg-blue-100 text-blue-800',
    llm: 'bg-purple-100 text-purple-800',
    pvp: 'bg-green-100 text-green-800'
  }
  return classes[type] || 'bg-gray-100 text-gray-800'
}

const getStatusName = (status: string) => {
  const names: Record<string, string> = {
    completed: '已完成',
    abandoned: '已放弃',
    in_progress: '进行中'
  }
  return names[status] || status
}

const getStatusClass = (status: string) => {
  const classes: Record<string, string> = {
    completed: 'bg-green-100 text-green-800',
    abandoned: 'bg-red-100 text-red-800',
    in_progress: 'bg-yellow-100 text-yellow-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

const getResultName = (winner: string) => {
  if (winner === userStore.user?.id) return '胜利'
  if (winner === 'draw') return '平局'
  return '失败'
}

const getResultClass = (winner: string) => {
  if (winner === userStore.user?.id) return 'bg-green-100 text-green-800'
  if (winner === 'draw') return 'bg-yellow-100 text-yellow-800'
  return 'bg-red-100 text-red-800'
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
  return new Date(dateString).toLocaleString('zh-CN')
}

onMounted(() => {
  fetchGameHistory()
})
</script>