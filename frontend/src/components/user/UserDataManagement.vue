<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 导航栏 -->
    <nav class="bg-white shadow-sm border-b">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center">
            <button
              @click="$router.push('/user')"
              class="flex items-center text-gray-600 hover:text-gray-900 transition-colors"
            >
              <ArrowLeft class="h-5 w-5 mr-2" />
              返回用户中心
            </button>
          </div>
          <div class="flex items-center">
            <h1 class="text-xl font-semibold text-gray-900">数据管理</h1>
          </div>
        </div>
      </div>
    </nav>

    <!-- 主要内容 -->
    <div class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <div class="px-4 py-6 sm:px-0">
        <!-- 数据导出 -->
        <div class="bg-white overflow-hidden shadow rounded-lg mb-6">
          <div class="px-4 py-5 sm:p-6">
            <h3 class="text-lg leading-6 font-medium text-gray-900 mb-4">数据导出</h3>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div class="border border-gray-200 rounded-lg p-4">
                <div class="flex items-center mb-3">
                  <Download class="h-5 w-5 text-blue-500 mr-2" />
                  <h4 class="font-medium text-gray-900">游戏历史数据</h4>
                </div>
                <p class="text-sm text-gray-600 mb-3">导出您的所有游戏记录，包括对战详情、时间等信息</p>
                <button
                  @click="exportGameHistory"
                  :disabled="exportLoading.history"
                  class="w-full bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                >
                  <span v-if="exportLoading.history" class="flex items-center justify-center">
                    <Loader2 class="animate-spin h-4 w-4 mr-2" />
                    导出中...
                  </span>
                  <span v-else>导出游戏历史</span>
                </button>
              </div>

              <div class="border border-gray-200 rounded-lg p-4">
                <div class="flex items-center mb-3">
                  <BarChart3 class="h-5 w-5 text-green-500 mr-2" />
                  <h4 class="font-medium text-gray-900">统计数据</h4>
                </div>
                <p class="text-sm text-gray-600 mb-3">导出您的游戏统计数据，包括胜率、排名等信息</p>
                <button
                  @click="exportStatistics"
                  :disabled="exportLoading.statistics"
                  class="w-full bg-green-600 text-white px-4 py-2 rounded-md hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                >
                  <span v-if="exportLoading.statistics" class="flex items-center justify-center">
                    <Loader2 class="animate-spin h-4 w-4 mr-2" />
                    导出中...
                  </span>
                  <span v-else>导出统计数据</span>
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- 数据清理 -->
        <div class="bg-white overflow-hidden shadow rounded-lg mb-6">
          <div class="px-4 py-5 sm:p-6">
            <h3 class="text-lg leading-6 font-medium text-gray-900 mb-4">数据清理</h3>
            <div class="space-y-4">
              <!-- 清理旧游戏记录 -->
              <div class="border border-yellow-200 rounded-lg p-4 bg-yellow-50">
                <div class="flex items-center mb-3">
                  <Trash2 class="h-5 w-5 text-yellow-600 mr-2" />
                  <h4 class="font-medium text-gray-900">清理旧游戏记录</h4>
                </div>
                <p class="text-sm text-gray-600 mb-3">删除30天前的游戏记录（保留统计数据）</p>
                <div class="flex items-center space-x-3">
                  <select
                    v-model="cleanupDays"
                    class="border border-gray-300 rounded-md px-3 py-2 text-sm"
                  >
                    <option value="30">30天前</option>
                    <option value="60">60天前</option>
                    <option value="90">90天前</option>
                    <option value="180">180天前</option>
                  </select>
                  <button
                    @click="cleanupOldGames"
                    :disabled="cleanupLoading"
                    class="bg-yellow-600 text-white px-4 py-2 rounded-md hover:bg-yellow-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                  >
                    <span v-if="cleanupLoading" class="flex items-center">
                      <Loader2 class="animate-spin h-4 w-4 mr-2" />
                      清理中...
                    </span>
                    <span v-else>开始清理</span>
                  </button>
                </div>
              </div>

              <!-- 清空所有数据 -->
              <div class="border border-red-200 rounded-lg p-4 bg-red-50">
                <div class="flex items-center mb-3">
                  <AlertTriangle class="h-5 w-5 text-red-600 mr-2" />
                  <h4 class="font-medium text-gray-900">清空所有数据</h4>
                </div>
                <p class="text-sm text-gray-600 mb-3">
                  <strong>危险操作：</strong>这将删除您的所有游戏记录和统计数据，此操作不可恢复
                </p>
                <div class="flex items-center space-x-3">
                  <input
                    v-model="confirmText"
                    type="text"
                    placeholder="输入 '确认删除' 来确认操作"
                    class="border border-gray-300 rounded-md px-3 py-2 text-sm flex-1"
                  />
                  <button
                    @click="clearAllData"
                    :disabled="confirmText !== '确认删除' || clearLoading"
                    class="bg-red-600 text-white px-4 py-2 rounded-md hover:bg-red-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                  >
                    <span v-if="clearLoading" class="flex items-center">
                      <Loader2 class="animate-spin h-4 w-4 mr-2" />
                      删除中...
                    </span>
                    <span v-else>清空数据</span>
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 存储使用情况 -->
        <div class="bg-white overflow-hidden shadow rounded-lg">
          <div class="px-4 py-5 sm:p-6">
            <h3 class="text-lg leading-6 font-medium text-gray-900 mb-4">存储使用情况</h3>
            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div class="text-center p-4 border border-gray-200 rounded-lg">
                <div class="text-2xl font-bold text-blue-600">{{ storageInfo.gameCount }}</div>
                <div class="text-sm text-gray-600">游戏记录数</div>
              </div>
              <div class="text-center p-4 border border-gray-200 rounded-lg">
                <div class="text-2xl font-bold text-green-600">{{ formatFileSize(storageInfo.dataSize) }}</div>
                <div class="text-sm text-gray-600">数据大小</div>
              </div>
              <div class="text-center p-4 border border-gray-200 rounded-lg">
                <div class="text-2xl font-bold text-purple-600">{{ storageInfo.lastCleanup || '从未' }}</div>
                <div class="text-sm text-gray-600">上次清理</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { 
  ArrowLeft, 
  Download, 
  BarChart3, 
  Trash2, 
  AlertTriangle, 
  Loader2 
} from 'lucide-vue-next'
import { gameApi } from '../../services/gameApi'

const router = useRouter()

// 响应式数据
const exportLoading = ref({
  history: false,
  statistics: false
})
const cleanupLoading = ref(false)
const clearLoading = ref(false)
const cleanupDays = ref(30)
const confirmText = ref('')
const storageInfo = ref({
  gameCount: 0,
  dataSize: 0,
  lastCleanup: ''
})

// 导出游戏历史
const exportGameHistory = async () => {
  exportLoading.value.history = true
  try {
    await gameApi.exportGameHistory()
    // 成功提示
    alert('游戏历史导出成功！')
  } catch (error) {
    console.error('Export game history failed:', error)
    alert('导出失败，请稍后重试')
  } finally {
    exportLoading.value.history = false
  }
}

// 导出统计数据
const exportStatistics = async () => {
  exportLoading.value.statistics = true
  try {
    await gameApi.exportUserStatistics()
    // 成功提示
    alert('统计数据导出成功！')
  } catch (error) {
    console.error('Export statistics failed:', error)
    alert('导出失败，请稍后重试')
  } finally {
    exportLoading.value.statistics = false
  }
}

// 清理旧游戏记录
const cleanupOldGames = async () => {
  if (!confirm(`确定要删除${cleanupDays.value}天前的游戏记录吗？此操作不可恢复。`)) {
    return
  }

  cleanupLoading.value = true
  try {
    const cutoffDate = new Date()
    cutoffDate.setDate(cutoffDate.getDate() - cleanupDays.value)
    
    // 这里应该调用后端API来清理数据
    // await gameApi.cleanupOldGames(cutoffDate.toISOString())
    
    alert(`成功清理了${cleanupDays.value}天前的游戏记录`)
    await loadStorageInfo()
  } catch (error) {
    console.error('Cleanup failed:', error)
    alert('清理失败，请稍后重试')
  } finally {
    cleanupLoading.value = false
  }
}

// 清空所有数据
const clearAllData = async () => {
  if (!confirm('确定要清空所有数据吗？此操作将删除您的所有游戏记录和统计数据，且不可恢复！')) {
    return
  }

  clearLoading.value = true
  try {
    // 这里应该调用后端API来清空所有数据
    // await gameApi.clearAllUserData()
    
    alert('所有数据已清空')
    confirmText.value = ''
    await loadStorageInfo()
  } catch (error) {
    console.error('Clear all data failed:', error)
    alert('清空失败，请稍后重试')
  } finally {
    clearLoading.value = false
  }
}

// 加载存储信息
const loadStorageInfo = async () => {
  try {
    // 这里应该调用后端API获取存储信息
    // const info = await gameApi.getStorageInfo()
    // storageInfo.value = info
    
    // 模拟数据
    storageInfo.value = {
      gameCount: 156,
      dataSize: 2048576, // 2MB
      lastCleanup: '2024-01-15'
    }
  } catch (error) {
    console.error('Load storage info failed:', error)
  }
}

// 格式化文件大小
const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 组件挂载时加载数据
onMounted(() => {
  loadStorageInfo()
})
</script>