<template>
  <div class="bg-white rounded-lg shadow-md p-6">
    <div class="flex items-center justify-between mb-4">
      <h3 class="text-lg font-semibold text-gray-800 flex items-center">
        <Info class="h-5 w-5 mr-2 text-blue-500" />
        游戏信息
      </h3>
      <button
        v-if="currentGame"
        @click="toggleExpanded"
        class="p-1 hover:bg-gray-100 rounded transition-colors duration-200"
      >
        <ChevronDown
          class="h-4 w-4 text-gray-500 transition-transform duration-200"
          :class="{ 'rotate-180': isExpanded }"
        />
      </button>
    </div>

    <div v-if="!currentGame" class="text-center py-8 text-gray-500">
      <GamepadIcon class="h-12 w-12 mx-auto mb-3 text-gray-300" />
      <p>暂无游戏信息</p>
      <p class="text-sm mt-1">开始一局新游戏来查看详细信息</p>
    </div>

    <div v-else class="space-y-4">
      <!-- 基本信息 -->
      <div class="grid grid-cols-2 gap-4">
        <div class="bg-blue-50 rounded-lg p-3">
          <div class="flex items-center mb-1">
            <Hash class="h-4 w-4 mr-1 text-blue-500" />
            <span class="text-sm font-medium text-blue-700">游戏ID</span>
          </div>
          <div class="text-xs text-blue-600 font-mono">{{ currentGame.id.slice(0, 8) }}...</div>
        </div>
        
        <div class="bg-green-50 rounded-lg p-3">
          <div class="flex items-center mb-1">
            <Calendar class="h-4 w-4 mr-1 text-green-500" />
            <span class="text-sm font-medium text-green-700">开始时间</span>
          </div>
          <div class="text-xs text-green-600">{{ formatTime(currentGame.createdAt) }}</div>
        </div>
      </div>

      <!-- 模型信息 -->
      <div class="bg-purple-50 rounded-lg p-4">
        <div class="flex items-center justify-between mb-3">
          <div class="flex items-center">
            <Brain class="h-5 w-5 mr-2 text-purple-500" />
            <span class="font-medium text-purple-700">AI模型信息</span>
          </div>
          <div
            class="px-2 py-1 rounded-full text-xs font-medium"
            :class="getModelStatusStyle()"
          >
            {{ getModelStatusText() }}
          </div>
        </div>
        
        <div class="space-y-2">
          <div class="flex justify-between items-center">
            <span class="text-sm text-gray-600">模型名称:</span>
            <span class="text-sm font-medium text-purple-600">{{ selectedModel?.name || 'Unknown' }}</span>
          </div>
          <div class="flex justify-between items-center">
            <span class="text-sm text-gray-600">提供商:</span>
            <span class="text-sm font-medium text-purple-600">{{ selectedModel?.provider || 'Unknown' }}</span>
          </div>
          <div v-if="selectedModel?.description" class="text-xs text-gray-500 mt-2">
            {{ selectedModel.description }}
          </div>
        </div>
      </div>

      <!-- 展开的详细信息 -->
      <div v-if="isExpanded" class="space-y-4 border-t pt-4">
        <!-- 游戏配置 -->
        <div class="bg-gray-50 rounded-lg p-4">
          <div class="flex items-center mb-3">
            <Settings class="h-5 w-5 mr-2 text-gray-500" />
            <span class="font-medium text-gray-700">游戏配置</span>
          </div>
          <div class="grid grid-cols-2 gap-3 text-sm">
            <div class="flex justify-between">
              <span class="text-gray-600">棋盘大小:</span>
              <span class="font-medium">15×15</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-600">获胜条件:</span>
              <span class="font-medium">五子连珠</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-600">先手:</span>
              <span class="font-medium">玩家 (黑子)</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-600">后手:</span>
              <span class="font-medium">AI (白子)</span>
            </div>
          </div>
        </div>

        <!-- 性能统计 -->
        <div class="bg-yellow-50 rounded-lg p-4">
          <div class="flex items-center mb-3">
            <Zap class="h-5 w-5 mr-2 text-yellow-500" />
            <span class="font-medium text-yellow-700">性能统计</span>
          </div>
          <div class="grid grid-cols-2 gap-3 text-sm">
            <div class="flex justify-between">
              <span class="text-gray-600">平均响应时间:</span>
              <span class="font-medium">{{ getAverageResponseTime() }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-600">最快响应:</span>
              <span class="font-medium">{{ getFastestResponse() }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-600">最慢响应:</span>
              <span class="font-medium">{{ getSlowestResponse() }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-600">总请求数:</span>
              <span class="font-medium">{{ gameStats.aiMoves }}</span>
            </div>
          </div>
        </div>

        <!-- 游戏规则提示 -->
        <div class="bg-blue-50 rounded-lg p-4">
          <div class="flex items-center mb-3">
            <BookOpen class="h-5 w-5 mr-2 text-blue-500" />
            <span class="font-medium text-blue-700">游戏规则</span>
          </div>
          <ul class="text-sm text-blue-600 space-y-1">
            <li>• 玩家执黑子先行，AI执白子后行</li>
            <li>• 在棋盘上连成五子者获胜</li>
            <li>• 可以横向、纵向或斜向连成五子</li>
            <li>• 棋盘填满且无人获胜则为平局</li>
          </ul>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="flex space-x-2 pt-2">
        <button
          @click="copyGameInfo"
          class="flex-1 px-3 py-2 bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-lg transition-colors duration-200 text-sm"
        >
          <Copy class="h-4 w-4 inline mr-1" />
          复制信息
        </button>
        <button
          v-if="currentGame.status === 'finished'"
          @click="shareGame"
          class="flex-1 px-3 py-2 bg-blue-100 hover:bg-blue-200 text-blue-700 rounded-lg transition-colors duration-200 text-sm"
        >
          <Share2 class="h-4 w-4 inline mr-1" />
          分享游戏
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import {
  Info,
  ChevronDown,
  Hash,
  Calendar,
  Brain,
  Settings,
  Zap,
  BookOpen,
  Copy,
  Share2,
  Gamepad2 as GamepadIcon
} from 'lucide-vue-next'
import { useLLMGameStore } from '../../stores/llmGame'
import { storeToRefs } from 'pinia'
import { llmGameUtils } from '../../services/llmApi'

// 使用store
const llmGameStore = useLLMGameStore()
const {
  currentGame,
  selectedModel,
  gameStats
} = storeToRefs(llmGameStore)

// 组件状态
const isExpanded = ref(false)

// 切换展开状态
function toggleExpanded() {
  isExpanded.value = !isExpanded.value
}

// 格式化时间
function formatTime(timestamp: string): string {
  return llmGameUtils.formatTimestamp(timestamp)
}

// 获取模型状态样式
function getModelStatusStyle(): string {
  if (!selectedModel.value) return 'bg-gray-100 text-gray-600'
  
  switch (selectedModel.value.status) {
    case 'available':
      return 'bg-green-100 text-green-700'
    case 'unavailable':
      return 'bg-red-100 text-red-700'
    case 'not_configured':
      return 'bg-yellow-100 text-yellow-700'
    default:
      return 'bg-gray-100 text-gray-600'
  }
}

// 获取模型状态文本
function getModelStatusText(): string {
  if (!selectedModel.value) return '未知'
  
  switch (selectedModel.value.status) {
    case 'available':
      return '可用'
    case 'unavailable':
      return '不可用'
    case 'not_configured':
      return '未配置'
    default:
      return '未知'
  }
}

// 获取平均响应时间
function getAverageResponseTime(): string {
  if (!currentGame.value || gameStats.value.aiMoves === 0) return '-'
  
  const aiMoves = currentGame.value.moves.filter(move => move.player === 2)
  if (aiMoves.length === 0) return '-'
  
  const totalTime = aiMoves.reduce((sum, move) => {
    return sum + (move.thinkingTime || 0)
  }, 0)
  
  const avgTime = totalTime / aiMoves.length
  return `${(avgTime / 1000).toFixed(1)}s`
}

// 获取最快响应时间
function getFastestResponse(): string {
  if (!currentGame.value) return '-'
  
  const aiMoves = currentGame.value.moves.filter(move => move.player === 2 && move.thinkingTime)
  if (aiMoves.length === 0) return '-'
  
  const fastest = Math.min(...aiMoves.map(move => move.thinkingTime || 0))
  return `${(fastest / 1000).toFixed(1)}s`
}

// 获取最慢响应时间
function getSlowestResponse(): string {
  if (!currentGame.value) return '-'
  
  const aiMoves = currentGame.value.moves.filter(move => move.player === 2 && move.thinkingTime)
  if (aiMoves.length === 0) return '-'
  
  const slowest = Math.max(...aiMoves.map(move => move.thinkingTime || 0))
  return `${(slowest / 1000).toFixed(1)}s`
}

// 复制游戏信息
async function copyGameInfo() {
  if (!currentGame.value) return
  
  const info = `
游戏ID: ${currentGame.value.id}
开始时间: ${formatTime(currentGame.value.createdAt)}
AI模型: ${selectedModel.value?.name || 'Unknown'}
游戏状态: ${currentGame.value.status}
总步数: ${gameStats.value.totalMoves}
游戏时长: ${gameStats.value.gameDuration}
  `.trim()
  
  try {
    await navigator.clipboard.writeText(info)
    // 这里可以添加成功提示
    console.log('游戏信息已复制到剪贴板')
  } catch (error) {
    console.error('复制失败:', error)
  }
}

// 分享游戏
function shareGame() {
  if (!currentGame.value) return
  
  const shareData = {
    title: '五子棋 AI 对战',
    text: `我刚刚与 ${selectedModel.value?.name || 'AI'} 进行了一局五子棋对战！`,
    url: window.location.href
  }
  
  if (navigator.share) {
    navigator.share(shareData).catch(console.error)
  } else {
    // 降级处理：复制链接
    copyGameInfo()
  }
}
</script>