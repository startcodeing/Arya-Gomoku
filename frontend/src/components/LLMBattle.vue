<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 p-4">
    <div class="max-w-7xl mx-auto">
      <!-- 页面标题 -->
      <div class="text-center mb-8">
        <h1 class="text-4xl font-bold text-gray-800 mb-2">与AI对战</h1>
        <p class="text-gray-600">挑战强大的LLM模型，体验智能五子棋对战</p>
      </div>

      <!-- 错误提示 -->
      <div v-if="error" class="mb-6 bg-red-50 border border-red-200 rounded-lg p-4">
        <div class="flex items-center">
          <AlertCircle class="h-5 w-5 text-red-500 mr-2" />
          <span class="text-red-700">{{ error }}</span>
          <button @click="clearError" class="ml-auto text-red-500 hover:text-red-700">
            <X class="h-4 w-4" />
          </button>
        </div>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-4 gap-6">
        <!-- 左侧控制面板 -->
        <div class="lg:col-span-1 space-y-6">
          <!-- 模型选择器 -->
          <ModelSelector />
          
          <!-- 游戏控制 -->
          <GameControls />
          
          <!-- 游戏信息 -->
          <GameInfo />
          
          <!-- AI思考状态 -->
          <AIThinkingStatus v-if="isThinking" />
        </div>

        <!-- 中间棋盘区域 -->
        <div class="lg:col-span-2">
          <GameBoard />
        </div>

        <!-- 右侧信息面板 -->
        <div class="lg:col-span-1 space-y-6">
          <!-- 游戏状态 -->
          <GameStatus />
          
          <!-- 移动历史 -->
          <MoveHistory />
          
          <!-- AI分析 -->
          <AIAnalysis v-if="lastAIMove" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { AlertCircle, X } from 'lucide-vue-next'
import { useLLMGameStore } from '../stores/llmGame'
import { storeToRefs } from 'pinia'

// 导入子组件
import ModelSelector from './llm/ModelSelector.vue'
import GameControls from './llm/GameControls.vue'
import GameInfo from './llm/GameInfo.vue'
import AIThinkingStatus from './llm/AIThinkingStatus.vue'
import GameBoard from './llm/GameBoard.vue'
import GameStatus from './llm/GameStatus.vue'
import MoveHistory from './llm/MoveHistory.vue'
import AIAnalysis from './llm/AIAnalysis.vue'

// 使用store
const llmGameStore = useLLMGameStore()
const { error, isThinking, lastAIMove } = storeToRefs(llmGameStore)
const { clearError, initialize, endGame } = llmGameStore

// 组件挂载时初始化
onMounted(async () => {
  await initialize()
})

// 组件卸载时清理
onUnmounted(() => {
  // 如果有活跃游戏，询问是否结束
  if (llmGameStore.isGameActive) {
    const shouldEnd = confirm('离开页面将结束当前游戏，确定要离开吗？')
    if (shouldEnd) {
      endGame()
    }
  }
})
</script>

<style scoped>
/* 自定义样式 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>