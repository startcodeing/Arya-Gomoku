<template>
  <div
    v-if="isThinking"
    class="bg-gradient-to-r from-purple-50 to-blue-50 rounded-lg shadow-md p-6 border border-purple-200"
  >
    <div class="flex items-center justify-between mb-4">
      <h3 class="text-lg font-semibold text-purple-800 flex items-center">
        <Brain class="h-5 w-5 mr-2 text-purple-600 animate-pulse" />
        AI思考中
      </h3>
      <div class="flex items-center space-x-2">
        <div class="flex space-x-1">
          <div
            v-for="i in 3"
            :key="i"
            class="w-2 h-2 bg-purple-500 rounded-full animate-bounce"
            :style="{ animationDelay: `${(i - 1) * 0.2}s` }"
          ></div>
        </div>
        <span class="text-sm text-purple-600 font-medium">{{ thinkingTime }}s</span>
      </div>
    </div>

    <!-- 思考进度 -->
    <div class="space-y-4">
      <!-- 当前分析阶段 -->
      <div class="bg-white rounded-lg p-4 shadow-sm">
        <div class="flex items-center justify-between mb-2">
          <span class="text-sm font-medium text-gray-700">当前阶段</span>
          <span class="text-xs text-purple-600 bg-purple-100 px-2 py-1 rounded-full">
            {{ getCurrentPhase() }}
          </span>
        </div>
        <div class="text-sm text-gray-600">{{ getCurrentPhaseDescription() }}</div>
      </div>

      <!-- 思考进度条 -->
      <div class="space-y-2">
        <div class="flex justify-between text-sm">
          <span class="text-gray-600">分析进度</span>
          <span class="text-purple-600 font-medium">{{ getProgressPercentage() }}%</span>
        </div>
        <div class="w-full bg-gray-200 rounded-full h-2">
          <div
            class="bg-gradient-to-r from-purple-500 to-blue-500 h-2 rounded-full transition-all duration-500 relative overflow-hidden"
            :style="{ width: getProgressPercentage() + '%' }"
          >
            <!-- 进度条动画效果 -->
            <div class="absolute inset-0 bg-white opacity-30 animate-pulse"></div>
          </div>
        </div>
      </div>

      <!-- 分析详情 -->
      <div class="grid grid-cols-2 gap-4">
        <div class="bg-blue-50 rounded-lg p-3">
          <div class="flex items-center mb-1">
            <Target class="h-4 w-4 mr-1 text-blue-500" />
            <span class="text-sm font-medium text-blue-700">分析深度</span>
          </div>
          <div class="text-lg font-bold text-blue-600">{{ getAnalysisDepth() }}</div>
          <div class="text-xs text-blue-500">层</div>
        </div>
        
        <div class="bg-green-50 rounded-lg p-3">
          <div class="flex items-center mb-1">
            <Zap class="h-4 w-4 mr-1 text-green-500" />
            <span class="text-sm font-medium text-green-700">候选位置</span>
          </div>
          <div class="text-lg font-bold text-green-600">{{ getCandidatePositions() }}</div>
          <div class="text-xs text-green-500">个</div>
        </div>
      </div>

      <!-- AI思考提示 -->
      <div class="bg-yellow-50 rounded-lg p-3 border border-yellow-200">
        <div class="flex items-start">
          <Lightbulb class="h-4 w-4 mr-2 text-yellow-500 mt-0.5 flex-shrink-0" />
          <div class="text-sm text-yellow-700">
            <div class="font-medium mb-1">AI正在分析</div>
            <div class="text-xs">{{ getThinkingTip() }}</div>
          </div>
        </div>
      </div>

      <!-- 模型信息 -->
      <div v-if="selectedModel" class="bg-purple-50 rounded-lg p-3 border border-purple-200">
        <div class="flex items-center justify-between mb-2">
          <div class="flex items-center">
            <Cpu class="h-4 w-4 mr-1 text-purple-500" />
            <span class="text-sm font-medium text-purple-700">{{ selectedModel.name }}</span>
          </div>
          <div class="text-xs text-purple-600 bg-purple-100 px-2 py-1 rounded-full">
            {{ selectedModel.provider }}
          </div>
        </div>
        <div class="text-xs text-purple-600">
          {{ getModelThinkingDescription() }}
        </div>
      </div>

      <!-- 取消按钮 -->
      <div class="flex justify-center pt-2">
        <button
          @click="handleCancelThinking"
          :disabled="!canCancel"
          class="px-4 py-2 bg-red-100 hover:bg-red-200 text-red-700 rounded-lg transition-colors duration-200 disabled:opacity-50 disabled:cursor-not-allowed text-sm"
        >
          <X class="h-4 w-4 inline mr-1" />
          取消思考
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import {
  Brain,
  Target,
  Zap,
  Lightbulb,
  Cpu,
  X
} from 'lucide-vue-next'
import { useLLMGameStore } from '../../stores/llmGame'
import { storeToRefs } from 'pinia'

// 使用store
const llmGameStore = useLLMGameStore()
const {
  isThinking,
  selectedModel
} = storeToRefs(llmGameStore)

// 组件状态
const thinkingTime = ref(0)
const thinkingPhase = ref(0)
const canCancel = ref(true)

// 定时器
let thinkingTimer: NodeJS.Timeout | null = null
let phaseTimer: NodeJS.Timeout | null = null

// 思考阶段
const phases = [
  { name: '初始化', description: '准备分析棋局' },
  { name: '局面评估', description: '评估当前棋局形势' },
  { name: '候选生成', description: '生成可能的落子位置' },
  { name: '深度搜索', description: '深度分析各种可能性' },
  { name: '策略优化', description: '优化落子策略' },
  { name: '最终决策', description: '确定最佳落子位置' }
]

// 思考提示
const thinkingTips = [
  '正在分析棋盘上的威胁和机会...',
  '计算各种落子位置的得分...',
  '评估攻击和防守的平衡...',
  '寻找最优的连子机会...',
  '分析对手可能的反击...',
  '优化长期战略布局...'
]

// 监听思考状态变化
const startThinking = () => {
  if (!isThinking.value) return
  
  thinkingTime.value = 0
  thinkingPhase.value = 0
  canCancel.value = true
  
  // 开始计时
  thinkingTimer = setInterval(() => {
    thinkingTime.value += 0.1
  }, 100)
  
  // 阶段切换
  phaseTimer = setInterval(() => {
    if (thinkingPhase.value < phases.length - 1) {
      thinkingPhase.value++
    }
  }, 2000)
  
  // 5秒后不允许取消
  setTimeout(() => {
    canCancel.value = false
  }, 5000)
}

const stopThinking = () => {
  if (thinkingTimer) {
    clearInterval(thinkingTimer)
    thinkingTimer = null
  }
  if (phaseTimer) {
    clearInterval(phaseTimer)
    phaseTimer = null
  }
  thinkingTime.value = 0
  thinkingPhase.value = 0
  canCancel.value = true
}

// 监听思考状态
const unwatch = computed(() => isThinking.value)
const watchThinking = () => {
  if (isThinking.value) {
    startThinking()
  } else {
    stopThinking()
  }
}

// 获取当前阶段
function getCurrentPhase(): string {
  return phases[thinkingPhase.value]?.name || '思考中'
}

// 获取当前阶段描述
function getCurrentPhaseDescription(): string {
  return phases[thinkingPhase.value]?.description || '正在处理...'
}

// 获取进度百分比
function getProgressPercentage(): number {
  const baseProgress = (thinkingPhase.value / (phases.length - 1)) * 80
  const timeProgress = Math.min((thinkingTime.value / 10) * 20, 20)
  return Math.min(baseProgress + timeProgress, 95)
}

// 获取分析深度
function getAnalysisDepth(): number {
  return Math.min(Math.floor(thinkingTime.value / 2) + 3, 8)
}

// 获取候选位置数
function getCandidatePositions(): number {
  return Math.max(15 - Math.floor(thinkingTime.value), 3)
}

// 获取思考提示
function getThinkingTip(): string {
  const index = Math.floor(thinkingTime.value / 3) % thinkingTips.length
  return thinkingTips[index]
}

// 获取模型思考描述
function getModelThinkingDescription(): string {
  if (!selectedModel.value) return ''
  
  switch (selectedModel.value.provider) {
    case 'deepseek':
      return '使用深度学习算法分析最优落子策略'
    case 'openai':
      return '运用GPT模型进行棋局推理和决策'
    case 'ollama':
      return '本地模型正在进行棋局分析'
    default:
      return '正在运行AI算法分析棋局'
  }
}

// 处理取消思考
function handleCancelThinking() {
  if (!canCancel.value) return
  
  // 这里可以调用store的方法来取消AI思考
  // 目前只是停止本地计时器
  stopThinking()
}

// 生命周期
onMounted(() => {
  // 监听思考状态变化
  const stopWatcher = computed(() => isThinking.value)
  const unwatch = stopWatcher.effect(() => {
    watchThinking()
  })
  
  // 清理函数
  onUnmounted(() => {
    unwatch()
    stopThinking()
  })
})

onUnmounted(() => {
  stopThinking()
})
</script>

<style scoped>
@keyframes bounce {
  0%, 80%, 100% {
    transform: translateY(0);
  }
  40% {
    transform: translateY(-10px);
  }
}

.animate-bounce {
  animation: bounce 1.4s infinite;
}
</style>