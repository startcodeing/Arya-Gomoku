<template>
  <div class="bg-white rounded-lg shadow-md p-6">
    <div class="flex items-center justify-between mb-4">
      <h3 class="text-lg font-semibold text-gray-800 flex items-center">
        <Activity class="h-5 w-5 mr-2 text-green-500" />
        游戏状态
      </h3>
      <div class="flex items-center space-x-2">
        <div
          class="w-3 h-3 rounded-full"
          :class="getStatusColor()"
        ></div>
        <span class="text-sm font-medium" :class="getStatusTextColor()">
          {{ getStatusText() }}
        </span>
      </div>
    </div>

    <!-- 当前回合 -->
    <div class="space-y-4">
      <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
        <div class="flex items-center">
          <User class="h-5 w-5 mr-2 text-blue-500" />
          <span class="font-medium text-gray-700">当前回合</span>
        </div>
        <div class="flex items-center">
          <div
            class="w-4 h-4 rounded-full mr-2"
            :class="{
              'bg-black': currentPlayer === 1,
              'bg-white border border-gray-300': currentPlayer === 2
            }"
          ></div>
          <span class="font-semibold" :class="getCurrentPlayerColor()">
            {{ getCurrentPlayerText() }}
          </span>
        </div>
      </div>

      <!-- 游戏进度 -->
      <div class="space-y-3">
        <div class="flex items-center justify-between">
          <span class="text-sm font-medium text-gray-600">游戏进度</span>
          <span class="text-sm text-gray-500">{{ gameStats.totalMoves }}/225 步</span>
        </div>
        <div class="w-full bg-gray-200 rounded-full h-2">
          <div
            class="bg-gradient-to-r from-blue-500 to-purple-500 h-2 rounded-full transition-all duration-300"
            :style="{ width: getProgressPercentage() + '%' }"
          ></div>
        </div>
      </div>

      <!-- 游戏统计 -->
      <div class="grid grid-cols-2 gap-4">
        <div class="text-center p-3 bg-blue-50 rounded-lg">
          <div class="text-2xl font-bold text-blue-600">{{ gameStats.playerMoves }}</div>
          <div class="text-sm text-blue-500">你的步数</div>
        </div>
        <div class="text-center p-3 bg-purple-50 rounded-lg">
          <div class="text-2xl font-bold text-purple-600">{{ gameStats.aiMoves }}</div>
          <div class="text-sm text-purple-500">AI步数</div>
        </div>
      </div>

      <!-- 游戏时长 -->
      <div class="flex items-center justify-between p-3 bg-yellow-50 rounded-lg">
        <div class="flex items-center">
          <Clock class="h-5 w-5 mr-2 text-yellow-500" />
          <span class="font-medium text-gray-700">游戏时长</span>
        </div>
        <span class="font-semibold text-yellow-600">{{ gameStats.gameDuration }}</span>
      </div>

      <!-- AI状态 -->
      <div v-if="currentGame" class="space-y-2">
        <div class="flex items-center justify-between p-3 bg-purple-50 rounded-lg">
          <div class="flex items-center">
            <Brain class="h-5 w-5 mr-2 text-purple-500" />
            <span class="font-medium text-gray-700">AI模型</span>
          </div>
          <span class="font-semibold text-purple-600">{{ selectedModel?.name || 'Unknown' }}</span>
        </div>

        <div v-if="isThinking" class="p-3 bg-orange-50 rounded-lg border border-orange-200">
          <div class="flex items-center">
            <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-orange-500 mr-2"></div>
            <span class="text-orange-700 font-medium">AI正在思考中...</span>
          </div>
          <div class="mt-2 text-sm text-orange-600">
            分析棋局，计算最佳落子位置
          </div>
        </div>
      </div>

      <!-- 游戏结果 -->
      <div v-if="isGameFinished" class="p-4 rounded-lg border-2" :class="getResultStyle()">
        <div class="flex items-center justify-center mb-2">
          <component
            :is="getResultIcon()"
            class="h-8 w-8 mr-3"
            :class="getResultIconColor()"
          />
          <span class="text-xl font-bold" :class="getResultTextColor()">
            {{ gameResultMessage }}
          </span>
        </div>
        <div class="text-center text-sm" :class="getResultSubTextColor()">
          {{ getResultSubText() }}
        </div>
      </div>

      <!-- 快捷操作 -->
      <div class="flex space-x-2 pt-2">
        <button
          v-if="!isGameFinished"
          @click="handleSurrender"
          :disabled="!currentGame || isThinking"
          class="flex-1 px-3 py-2 bg-red-100 hover:bg-red-200 text-red-700 rounded-lg transition-colors duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <Flag class="h-4 w-4 inline mr-1" />
          认输
        </button>
        <button
          v-if="isGameFinished"
          @click="handleRestart"
          class="flex-1 px-3 py-2 bg-green-100 hover:bg-green-200 text-green-700 rounded-lg transition-colors duration-200"
        >
          <RotateCcw class="h-4 w-4 inline mr-1" />
          再来一局
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import {
  Activity,
  User,
  Clock,
  Brain,
  Flag,
  RotateCcw,
  Trophy,
  X,
  Minus
} from 'lucide-vue-next'
import { useLLMGameStore } from '../../stores/llmGame'
import { storeToRefs } from 'pinia'
import { LLMGameStatus } from '../../types/game'

// 使用store
const llmGameStore = useLLMGameStore()
const {
  currentGame,
  selectedModel,
  isThinking,
  isGameFinished,
  gameResultMessage,
  gameStats
} = storeToRefs(llmGameStore)

const {
  endGame,
  restartGame
} = llmGameStore

// 当前玩家
const currentPlayer = computed(() => {
  if (!currentGame.value) return 1
  return currentGame.value.moves.length % 2 === 0 ? 1 : 2
})

// 获取状态颜色
function getStatusColor(): string {
  if (!currentGame.value) return 'bg-gray-400'
  
  switch (currentGame.value.status) {
    case LLMGameStatus.WAITING:
      return 'bg-yellow-400'
    case LLMGameStatus.PLAYING:
      return 'bg-green-400'
    case LLMGameStatus.FINISHED:
      return 'bg-blue-400'
    case LLMGameStatus.ERROR:
      return 'bg-red-400'
    default:
      return 'bg-gray-400'
  }
}

// 获取状态文本颜色
function getStatusTextColor(): string {
  if (!currentGame.value) return 'text-gray-600'
  
  switch (currentGame.value.status) {
    case LLMGameStatus.WAITING:
      return 'text-yellow-600'
    case LLMGameStatus.PLAYING:
      return 'text-green-600'
    case LLMGameStatus.FINISHED:
      return 'text-blue-600'
    case LLMGameStatus.ERROR:
      return 'text-red-600'
    default:
      return 'text-gray-600'
  }
}

// 获取状态文本
function getStatusText(): string {
  if (!currentGame.value) return '未开始'
  
  switch (currentGame.value.status) {
    case LLMGameStatus.WAITING:
      return '等待中'
    case LLMGameStatus.PLAYING:
      return isThinking.value ? 'AI思考中' : '游戏进行中'
    case LLMGameStatus.FINISHED:
      return '游戏结束'
    case LLMGameStatus.ERROR:
      return '游戏错误'
    default:
      return '未知状态'
  }
}

// 获取当前玩家文本
function getCurrentPlayerText(): string {
  if (isGameFinished.value) return '游戏结束'
  if (isThinking.value) return 'AI思考中'
  return currentPlayer.value === 1 ? '你的回合' : 'AI回合'
}

// 获取当前玩家颜色
function getCurrentPlayerColor(): string {
  if (isGameFinished.value) return 'text-gray-600'
  if (isThinking.value) return 'text-purple-600'
  return currentPlayer.value === 1 ? 'text-blue-600' : 'text-purple-600'
}

// 获取进度百分比
function getProgressPercentage(): number {
  return Math.min((gameStats.value.totalMoves / 225) * 100, 100)
}

// 获取结果样式
function getResultStyle(): string {
  if (gameResultMessage.value.includes('你获胜')) {
    return 'bg-green-50 border-green-200'
  } else if (gameResultMessage.value.includes('AI获胜')) {
    return 'bg-red-50 border-red-200'
  } else {
    return 'bg-yellow-50 border-yellow-200'
  }
}

// 获取结果图标
function getResultIcon() {
  if (gameResultMessage.value.includes('你获胜')) {
    return Trophy
  } else if (gameResultMessage.value.includes('AI获胜')) {
    return X
  } else {
    return Minus
  }
}

// 获取结果图标颜色
function getResultIconColor(): string {
  if (gameResultMessage.value.includes('你获胜')) {
    return 'text-green-500'
  } else if (gameResultMessage.value.includes('AI获胜')) {
    return 'text-red-500'
  } else {
    return 'text-yellow-500'
  }
}

// 获取结果文本颜色
function getResultTextColor(): string {
  if (gameResultMessage.value.includes('你获胜')) {
    return 'text-green-700'
  } else if (gameResultMessage.value.includes('AI获胜')) {
    return 'text-red-700'
  } else {
    return 'text-yellow-700'
  }
}

// 获取结果子文本颜色
function getResultSubTextColor(): string {
  if (gameResultMessage.value.includes('你获胜')) {
    return 'text-green-600'
  } else if (gameResultMessage.value.includes('AI获胜')) {
    return 'text-red-600'
  } else {
    return 'text-yellow-600'
  }
}

// 获取结果子文本
function getResultSubText(): string {
  if (gameResultMessage.value.includes('你获胜')) {
    return '恭喜！你战胜了AI'
  } else if (gameResultMessage.value.includes('AI获胜')) {
    return 'AI获得了胜利，再接再厉！'
  } else {
    return '平局，势均力敌！'
  }
}

// 处理认输
async function handleSurrender() {
  if (!currentGame.value) return
  
  try {
    await endGame()
  } catch (error) {
    console.error('Failed to surrender:', error)
  }
}

// 处理重新开始
async function handleRestart() {
  try {
    await restartGame()
  } catch (error) {
    console.error('Failed to restart game:', error)
  }
}
</script>