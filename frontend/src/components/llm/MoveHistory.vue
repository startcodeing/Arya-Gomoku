<template>
  <div class="bg-white rounded-lg shadow-md p-6">
    <div class="flex items-center justify-between mb-4">
      <h3 class="text-lg font-semibold text-gray-800 flex items-center">
        <History class="h-5 w-5 mr-2 text-indigo-500" />
        移动历史
      </h3>
      <div class="flex items-center space-x-2">
        <span class="text-sm text-gray-500">{{ moveHistory.length }} 步</span>
        <button
          v-if="moveHistory.length > 0"
          @click="toggleExpanded"
          class="p-1 hover:bg-gray-100 rounded transition-colors duration-200"
        >
          <ChevronDown
            class="h-4 w-4 text-gray-500 transition-transform duration-200"
            :class="{ 'rotate-180': isExpanded }"
          />
        </button>
      </div>
    </div>

    <div v-if="moveHistory.length === 0" class="text-center py-8 text-gray-500">
      <Clock class="h-12 w-12 mx-auto mb-3 text-gray-300" />
      <p>暂无移动记录</p>
      <p class="text-sm mt-1">开始游戏后将显示移动历史</p>
    </div>

    <div v-else class="space-y-4">
      <!-- 最近几步预览 -->
      <div class="space-y-2">
        <div
          v-for="(move, index) in getRecentMoves()"
          :key="index"
          class="flex items-center justify-between p-3 rounded-lg transition-all duration-200"
          :class="getMoveItemStyle(move, index)"
        >
          <div class="flex items-center space-x-3">
            <!-- 步数 -->
            <div
              class="w-8 h-8 rounded-full flex items-center justify-center text-sm font-bold"
              :class="getMoveNumberStyle(move)"
            >
              {{ moveHistory.length - index }}
            </div>
            
            <!-- 玩家信息 -->
            <div class="flex items-center space-x-2">
              <div
                class="w-4 h-4 rounded-full"
                :class="{
                  'bg-black': move.player === 1,
                  'bg-white border border-gray-300': move.player === 2
                }"
              ></div>
              <span class="font-medium" :class="getPlayerTextColor(move)">
                {{ move.player === 1 ? '你' : 'AI' }}
              </span>
            </div>
            
            <!-- 位置 -->
            <div class="text-sm text-gray-600">
              {{ formatPosition(move.x, move.y) }}
            </div>
          </div>
          
          <div class="flex items-center space-x-2">
            <!-- 思考时间 -->
            <div v-if="move.thinkingTime" class="text-xs text-gray-500 flex items-center">
              <Brain class="h-3 w-3 mr-1" />
              {{ formatThinkingTime(move.thinkingTime) }}
            </div>
            
            <!-- 置信度 -->
            <div v-if="move.confidence" class="text-xs flex items-center">
              <div
                class="w-2 h-2 rounded-full mr-1"
                :class="getConfidenceColor(move.confidence)"
              ></div>
              {{ formatConfidence(move.confidence) }}
            </div>
            
            <!-- 操作按钮 -->
            <button
              @click="highlightMove(move)"
              class="p-1 hover:bg-gray-100 rounded transition-colors duration-200"
              :title="'定位到第' + (moveHistory.length - index) + '步'"
            >
              <MapPin class="h-3 w-3 text-gray-400" />
            </button>
          </div>
        </div>
      </div>

      <!-- 展开的完整历史 -->
      <div v-if="isExpanded && moveHistory.length > 5" class="border-t pt-4">
        <div class="max-h-64 overflow-y-auto space-y-2">
          <div
            v-for="(move, index) in getOlderMoves()"
            :key="'old-' + index"
            class="flex items-center justify-between p-2 rounded-lg bg-gray-50 hover:bg-gray-100 transition-colors duration-200"
          >
            <div class="flex items-center space-x-3">
              <div class="w-6 h-6 rounded-full bg-gray-200 flex items-center justify-center text-xs font-medium text-gray-600">
                {{ moveHistory.length - index - 5 }}
              </div>
              <div class="flex items-center space-x-2">
                <div
                  class="w-3 h-3 rounded-full"
                  :class="{
                    'bg-black': move.player === 1,
                    'bg-white border border-gray-300': move.player === 2
                  }"
                ></div>
                <span class="text-sm text-gray-600">
                  {{ move.player === 1 ? '你' : 'AI' }}
                </span>
              </div>
              <div class="text-sm text-gray-500">
                {{ formatPosition(move.x, move.y) }}
              </div>
            </div>
            
            <div class="flex items-center space-x-2">
              <div v-if="move.thinkingTime" class="text-xs text-gray-400">
                {{ formatThinkingTime(move.thinkingTime) }}
              </div>
              <button
                @click="highlightMove(move)"
                class="p-1 hover:bg-gray-200 rounded transition-colors duration-200"
              >
                <MapPin class="h-3 w-3 text-gray-400" />
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- 统计信息 -->
      <div class="border-t pt-4">
        <div class="grid grid-cols-2 gap-4 text-sm">
          <div class="bg-blue-50 rounded-lg p-3">
            <div class="flex items-center justify-between">
              <span class="text-blue-600 font-medium">你的平均用时</span>
              <span class="text-blue-700 font-bold">{{ getPlayerAverageTime() }}</span>
            </div>
          </div>
          <div class="bg-purple-50 rounded-lg p-3">
            <div class="flex items-center justify-between">
              <span class="text-purple-600 font-medium">AI平均用时</span>
              <span class="text-purple-700 font-bold">{{ getAIAverageTime() }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="flex space-x-2 pt-2">
        <button
          @click="exportHistory"
          class="flex-1 px-3 py-2 bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-lg transition-colors duration-200 text-sm"
        >
          <Download class="h-4 w-4 inline mr-1" />
          导出历史
        </button>
        <button
          @click="clearHistory"
          :disabled="!canClearHistory"
          class="flex-1 px-3 py-2 bg-red-100 hover:bg-red-200 text-red-700 rounded-lg transition-colors duration-200 disabled:opacity-50 disabled:cursor-not-allowed text-sm"
        >
          <Trash2 class="h-4 w-4 inline mr-1" />
          清空历史
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import {
  History,
  ChevronDown,
  Clock,
  Brain,
  MapPin,
  Download,
  Trash2
} from 'lucide-vue-next'
import { useLLMGameStore } from '../../stores/llmGame'
import { storeToRefs } from 'pinia'
import { llmGameUtils } from '../../services/llmApi'
import type { LLMMove } from '../../types/game'

// 使用store
const llmGameStore = useLLMGameStore()
const {
  moveHistory,
  isGameFinished
} = storeToRefs(llmGameStore)

// 组件状态
const isExpanded = ref(false)
const highlightedMove = ref<LLMMove | null>(null)

// 切换展开状态
function toggleExpanded() {
  isExpanded.value = !isExpanded.value
}

// 获取最近的移动（最多5步）
function getRecentMoves(): LLMMove[] {
  return moveHistory.value.slice(-5).reverse()
}

// 获取更早的移动
function getOlderMoves(): LLMMove[] {
  if (moveHistory.value.length <= 5) return []
  return moveHistory.value.slice(0, -5).reverse()
}

// 获取移动项样式
function getMoveItemStyle(move: LLMMove, index: number): string {
  const isLatest = index === 0
  const isHighlighted = highlightedMove.value && 
    highlightedMove.value.x === move.x && 
    highlightedMove.value.y === move.y &&
    highlightedMove.value.timestamp === move.timestamp
  
  if (isHighlighted) {
    return 'bg-yellow-100 border border-yellow-300'
  } else if (isLatest) {
    return 'bg-blue-50 border border-blue-200'
  } else {
    return 'bg-gray-50 hover:bg-gray-100'
  }
}

// 获取移动编号样式
function getMoveNumberStyle(move: LLMMove): string {
  if (move.player === 1) {
    return 'bg-blue-100 text-blue-700'
  } else {
    return 'bg-purple-100 text-purple-700'
  }
}

// 获取玩家文本颜色
function getPlayerTextColor(move: LLMMove): string {
  return move.player === 1 ? 'text-blue-700' : 'text-purple-700'
}

// 格式化位置
function formatPosition(x: number, y: number): string {
  const letters = 'ABCDEFGHIJKLMNO'
  return `${letters[x]}${y + 1}`
}

// 格式化思考时间
function formatThinkingTime(time: number): string {
  return `${(time / 1000).toFixed(1)}s`
}

// 格式化置信度
function formatConfidence(confidence: number): string {
  return llmGameUtils.formatConfidence(confidence)
}

// 获取置信度颜色
function getConfidenceColor(confidence: number): string {
  return llmGameUtils.getConfidenceColor(confidence)
}

// 高亮移动
function highlightMove(move: LLMMove) {
  highlightedMove.value = move
  
  // 3秒后取消高亮
  setTimeout(() => {
    if (highlightedMove.value === move) {
      highlightedMove.value = null
    }
  }, 3000)
  
  // 这里可以触发棋盘上的高亮效果
  // 可以通过事件或store来通知棋盘组件
}

// 获取玩家平均用时
function getPlayerAverageTime(): string {
  const playerMoves = moveHistory.value.filter(move => move.player === 1 && move.thinkingTime)
  if (playerMoves.length === 0) return '-'
  
  const totalTime = playerMoves.reduce((sum, move) => sum + (move.thinkingTime || 0), 0)
  const avgTime = totalTime / playerMoves.length
  return `${(avgTime / 1000).toFixed(1)}s`
}

// 获取AI平均用时
function getAIAverageTime(): string {
  const aiMoves = moveHistory.value.filter(move => move.player === 2 && move.thinkingTime)
  if (aiMoves.length === 0) return '-'
  
  const totalTime = aiMoves.reduce((sum, move) => sum + (move.thinkingTime || 0), 0)
  const avgTime = totalTime / aiMoves.length
  return `${(avgTime / 1000).toFixed(1)}s`
}

// 是否可以清空历史
const canClearHistory = computed(() => {
  return moveHistory.value.length > 0 && isGameFinished.value
})

// 导出历史
function exportHistory() {
  if (moveHistory.value.length === 0) return
  
  const data = {
    gameId: 'current-game',
    timestamp: new Date().toISOString(),
    moves: moveHistory.value.map((move, index) => ({
      step: index + 1,
      player: move.player === 1 ? 'Player' : 'AI',
      position: formatPosition(move.x, move.y),
      coordinates: { x: move.x, y: move.y },
      timestamp: move.timestamp,
      thinkingTime: move.thinkingTime,
      confidence: move.confidence
    }))
  }
  
  const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `gomoku-history-${new Date().toISOString().slice(0, 10)}.json`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
}

// 清空历史
function clearHistory() {
  if (!canClearHistory.value) return
  
  if (confirm('确定要清空移动历史吗？此操作不可撤销。')) {
    // 这里可以调用store的方法来清空历史
    // 目前只是清空本地状态
    highlightedMove.value = null
  }
}
</script>