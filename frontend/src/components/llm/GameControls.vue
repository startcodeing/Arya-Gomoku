<template>
  <div class="bg-white rounded-lg shadow-md p-6">
    <h3 class="text-lg font-semibold text-gray-800 mb-4 flex items-center">
      <Gamepad2 class="h-5 w-5 mr-2 text-green-500" />
      游戏控制
    </h3>

    <div class="space-y-3">
      <!-- 开始新游戏 -->
      <button
        @click="handleStartGame"
        :disabled="!canStartGame || isLoading"
        class="w-full flex items-center justify-center px-4 py-3 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed text-white rounded-lg transition-colors duration-200 font-medium"
      >
        <Play class="h-4 w-4 mr-2" />
        <span v-if="isLoading">启动中...</span>
        <span v-else-if="isGameActive">重新开始</span>
        <span v-else>开始游戏</span>
      </button>

      <!-- 结束游戏 -->
      <button
        v-if="isGameActive"
        @click="handleEndGame"
        :disabled="isLoading"
        class="w-full flex items-center justify-center px-4 py-3 bg-red-600 hover:bg-red-700 disabled:bg-gray-400 disabled:cursor-not-allowed text-white rounded-lg transition-colors duration-200"
      >
        <Square class="h-4 w-4 mr-2" />
        结束游戏
      </button>

      <!-- 重新开始 -->
      <button
        v-if="isGameActive && !isGameFinished"
        @click="handleRestartGame"
        :disabled="isLoading || isThinking"
        class="w-full flex items-center justify-center px-4 py-3 bg-yellow-600 hover:bg-yellow-700 disabled:bg-gray-400 disabled:cursor-not-allowed text-white rounded-lg transition-colors duration-200"
      >
        <RotateCcw class="h-4 w-4 mr-2" />
        重新开始
      </button>
    </div>

    <!-- 游戏提示 -->
    <div class="mt-4 pt-4 border-t border-gray-200">
      <div class="text-sm text-gray-600">
        <div v-if="!selectedModelInfo" class="flex items-center text-yellow-600">
          <AlertTriangle class="h-4 w-4 mr-2" />
          请先选择一个可用的AI模型
        </div>
        <div v-else-if="selectedModelInfo.status !== 'available'" class="flex items-center text-red-600">
          <AlertTriangle class="h-4 w-4 mr-2" />
          所选模型未配置，请先配置API密钥
        </div>
        <div v-else-if="!isGameActive" class="flex items-center text-blue-600">
          <Info class="h-4 w-4 mr-2" />
          点击"开始游戏"开始与 {{ selectedModelInfo.display_name }} 对战
        </div>
        <div v-else-if="isThinking" class="flex items-center text-purple-600">
          <Brain class="h-4 w-4 mr-2 animate-pulse" />
          AI正在思考中...
        </div>
        <div v-else-if="isGameFinished" class="flex items-center text-green-600">
          <CheckCircle class="h-4 w-4 mr-2" />
          游戏已结束
        </div>
        <div v-else class="flex items-center text-green-600">
          <Play class="h-4 w-4 mr-2" />
          轮到你下棋
        </div>
      </div>
    </div>

    <!-- 快捷操作 -->
    <div class="mt-4 pt-4 border-t border-gray-200">
      <div class="text-sm font-medium text-gray-700 mb-2">快捷操作</div>
      <div class="grid grid-cols-2 gap-2">
        <button
          @click="showGameHistory"
          :disabled="!currentGame"
          class="flex items-center justify-center px-3 py-2 bg-gray-100 hover:bg-gray-200 disabled:bg-gray-50 disabled:cursor-not-allowed text-gray-700 rounded text-xs transition-colors duration-200"
        >
          <History class="h-3 w-3 mr-1" />
          历史
        </button>
        <button
          @click="showGameStats"
          :disabled="!currentGame"
          class="flex items-center justify-center px-3 py-2 bg-gray-100 hover:bg-gray-200 disabled:bg-gray-50 disabled:cursor-not-allowed text-gray-700 rounded text-xs transition-colors duration-200"
        >
          <BarChart3 class="h-3 w-3 mr-1" />
          统计
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import {
  Gamepad2,
  Play,
  Square,
  RotateCcw,
  AlertTriangle,
  Info,
  Brain,
  CheckCircle,
  History,
  BarChart3
} from 'lucide-vue-next'
import { useLLMGameStore } from '../../stores/llmGame'
import { storeToRefs } from 'pinia'

// 使用store
const llmGameStore = useLLMGameStore()
const {
  selectedModel,
  selectedModelInfo,
  isGameActive,
  isLoading,
  isThinking,
  isGameFinished,
  currentGame
} = storeToRefs(llmGameStore)

const {
  startNewGame,
  endGame,
  restartGame
} = llmGameStore

// 计算属性
const canStartGame = computed(() => {
  return selectedModelInfo.value && selectedModelInfo.value.status === 'available'
})

// 处理开始游戏
async function handleStartGame() {
  if (!canStartGame.value) {
    return
  }

  if (isGameActive.value) {
    const shouldRestart = confirm('这将开始一个新游戏，确定要继续吗？')
    if (!shouldRestart) {
      return
    }
  }

  try {
    await startNewGame()
  } catch (error) {
    console.error('Failed to start game:', error)
  }
}

// 处理结束游戏
async function handleEndGame() {
  const shouldEnd = confirm('确定要结束当前游戏吗？')
  if (shouldEnd) {
    await endGame()
  }
}

// 处理重新开始游戏
async function handleRestartGame() {
  const shouldRestart = confirm('确定要重新开始游戏吗？')
  if (shouldRestart) {
    try {
      await restartGame()
    } catch (error) {
      console.error('Failed to restart game:', error)
    }
  }
}

// 显示游戏历史
function showGameHistory() {
  // TODO: 实现游戏历史显示
  alert('游戏历史功能即将推出')
}

// 显示游戏统计
function showGameStats() {
  // TODO: 实现游戏统计显示
  alert('游戏统计功能即将推出')
}
</script>