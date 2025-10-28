<template>
  <div class="bg-white rounded-lg shadow-md p-6">
    <div class="flex items-center justify-between mb-4">
      <h3 class="text-lg font-semibold text-gray-800 flex items-center">
        <Grid3x3 class="h-5 w-5 mr-2 text-purple-500" />
        游戏棋盘
      </h3>
      <div class="text-sm text-gray-600">
        {{ gameStats.totalMoves }} 步 | {{ gameStats.gameDuration }}
      </div>
    </div>

    <!-- 棋盘容器 -->
    <div class="relative bg-yellow-100 rounded-lg p-4 overflow-hidden">
      <!-- 棋盘网格 -->
      <div 
        class="relative bg-yellow-200 rounded"
        :style="{ width: boardSize + 'px', height: boardSize + 'px' }"
      >
        <!-- 网格线 -->
        <svg
          :width="boardSize"
          :height="boardSize"
          class="absolute inset-0"
          style="pointer-events: none;"
        >
          <!-- 垂直线 -->
          <line
            v-for="i in 15"
            :key="'v' + i"
            :x1="cellSize * (i - 0.5)"
            :y1="cellSize * 0.5"
            :x2="cellSize * (i - 0.5)"
            :y2="cellSize * 14.5"
            stroke="#8B4513"
            stroke-width="1"
          />
          <!-- 水平线 -->
          <line
            v-for="i in 15"
            :key="'h' + i"
            :x1="cellSize * 0.5"
            :y1="cellSize * (i - 0.5)"
            :x2="cellSize * 14.5"
            :y2="cellSize * (i - 0.5)"
            stroke="#8B4513"
            stroke-width="1"
          />
          <!-- 星位点 -->
          <circle
            v-for="point in starPoints"
            :key="`star-${point.x}-${point.y}`"
            :cx="cellSize * (point.x + 0.5)"
            :cy="cellSize * (point.y + 0.5)"
            r="3"
            fill="#8B4513"
          />
        </svg>

        <!-- 棋子 -->
        <div
          v-for="(row, y) in board"
          :key="'row' + y"
          class="absolute"
          :style="{ top: cellSize * y + 'px' }"
        >
          <div
            v-for="(cell, x) in row"
            :key="'cell' + x + '-' + y"
            class="absolute flex items-center justify-center cursor-pointer transition-all duration-200"
            :style="{
              left: cellSize * x + 'px',
              width: cellSize + 'px',
              height: cellSize + 'px'
            }"
            :class="{
              'hover:bg-blue-200 hover:bg-opacity-30': canMakeMove(x, y),
              'cursor-not-allowed': !canMakeMove(x, y)
            }"
            @click="handleCellClick(x, y)"
          >
            <!-- 棋子 -->
            <div
              v-if="cell !== 0"
              class="rounded-full border-2 shadow-lg transition-all duration-300"
              :class="{
                'bg-black border-gray-800': cell === 1,
                'bg-white border-gray-300': cell === 2,
                'ring-4 ring-blue-400 ring-opacity-50': isLastMove(x, y),
                'ring-4 ring-red-400 ring-opacity-50': isAILastMove(x, y)
              }"
              :style="{
                width: (cellSize * 0.8) + 'px',
                height: (cellSize * 0.8) + 'px'
              }"
            >
              <!-- 棋子内部标记 -->
              <div
                v-if="isLastMove(x, y) || isAILastMove(x, y)"
                class="w-full h-full rounded-full flex items-center justify-center"
              >
                <div
                  class="w-2 h-2 rounded-full"
                  :class="{
                    'bg-white': cell === 1,
                    'bg-black': cell === 2
                  }"
                ></div>
              </div>
            </div>

            <!-- 悬停预览 -->
            <div
              v-else-if="canMakeMove(x, y)"
              class="rounded-full bg-black bg-opacity-20 opacity-0 hover:opacity-60 transition-opacity duration-200"
              :style="{
                width: (cellSize * 0.6) + 'px',
                height: (cellSize * 0.6) + 'px'
              }"
            ></div>
          </div>
        </div>

        <!-- AI思考指示器 -->
        <div
          v-if="isThinking"
          class="absolute inset-0 bg-purple-500 bg-opacity-10 flex items-center justify-center rounded"
        >
          <div class="bg-white rounded-lg shadow-lg p-4 flex items-center">
            <Brain class="h-6 w-6 text-purple-500 animate-pulse mr-3" />
            <span class="text-purple-700 font-medium">AI正在思考...</span>
          </div>
        </div>

        <!-- 游戏结束遮罩 -->
        <div
          v-if="isGameFinished"
          class="absolute inset-0 bg-black bg-opacity-20 flex items-center justify-center rounded"
        >
          <div class="bg-white rounded-lg shadow-lg p-6 text-center">
            <div class="text-2xl font-bold mb-2" :class="getResultColor()">
              {{ gameResultMessage }}
            </div>
            <div class="text-gray-600 mb-4">
              游戏时长: {{ gameStats.gameDuration }}
            </div>
            <button
              @click="handleRestartGame"
              class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors duration-200"
            >
              再来一局
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 棋盘说明 -->
    <div class="mt-4 flex items-center justify-between text-sm text-gray-600">
      <div class="flex items-center space-x-4">
        <div class="flex items-center">
          <div class="w-4 h-4 bg-black rounded-full mr-2"></div>
          <span>你 (黑子)</span>
        </div>
        <div class="flex items-center">
          <div class="w-4 h-4 bg-white border border-gray-300 rounded-full mr-2"></div>
          <span>AI (白子)</span>
        </div>
      </div>
      <div class="flex items-center">
        <div class="w-3 h-3 bg-blue-400 rounded-full mr-2"></div>
        <span>最后一步</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { Grid3x3, Brain } from 'lucide-vue-next'
import { useLLMGameStore } from '../../stores/llmGame'
import { storeToRefs } from 'pinia'

// 使用store
const llmGameStore = useLLMGameStore()
const {
  board,
  isThinking,
  isGameFinished,
  gameResultMessage,
  gameStats,
  moveHistory,
  lastAIMove
} = storeToRefs(llmGameStore)

const {
  makeMove,
  canMakeMove,
  getLastMove,
  restartGame
} = llmGameStore

// 棋盘尺寸配置
const cellSize = ref(30)
const boardSize = computed(() => cellSize.value * 15)

// 星位点坐标
const starPoints = [
  { x: 3, y: 3 }, { x: 11, y: 3 }, { x: 7, y: 7 },
  { x: 3, y: 11 }, { x: 11, y: 11 }
]

// 处理棋盘点击
async function handleCellClick(x: number, y: number) {
  if (!canMakeMove(x, y)) {
    return
  }

  try {
    await makeMove(x, y)
  } catch (error) {
    console.error('Failed to make move:', error)
  }
}

// 检查是否是最后一步移动
function isLastMove(x: number, y: number): boolean {
  const lastMove = getLastMove()
  return lastMove ? lastMove.x === x && lastMove.y === y && lastMove.player === 1 : false
}

// 检查是否是AI的最后一步移动
function isAILastMove(x: number, y: number): boolean {
  return lastAIMove.value ? lastAIMove.value.x === x && lastAIMove.value.y === y : false
}

// 获取游戏结果颜色
function getResultColor(): string {
  if (gameResultMessage.value.includes('你获胜')) {
    return 'text-green-600'
  } else if (gameResultMessage.value.includes('AI获胜')) {
    return 'text-red-600'
  } else {
    return 'text-yellow-600'
  }
}

// 处理重新开始游戏
async function handleRestartGame() {
  try {
    await restartGame()
  } catch (error) {
    console.error('Failed to restart game:', error)
  }
}
</script>

<style scoped>
/* 棋盘动画效果 */
.piece-enter-active {
  transition: all 0.3s ease;
}

.piece-enter-from {
  transform: scale(0);
  opacity: 0;
}

.piece-enter-to {
  transform: scale(1);
  opacity: 1;
}
</style>