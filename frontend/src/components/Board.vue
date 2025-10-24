<template>
  <div class="board-container">
    <div class="board" ref="boardRef">
      <!-- 棋盘网格线 -->
      <svg class="board-grid" :width="boardSize" :height="boardSize">
        <!-- 垂直线 -->
        <line
          v-for="i in 15"
          :key="`v-${i}`"
          :x1="(i - 1) * cellSize + cellSize / 2"
          :y1="cellSize / 2"
          :x2="(i - 1) * cellSize + cellSize / 2"
          :y2="boardSize - cellSize / 2"
          stroke="#8B4513"
          stroke-width="1"
        />
        <!-- 水平线 -->
        <line
          v-for="i in 15"
          :key="`h-${i}`"
          :x1="cellSize / 2"
          :y1="(i - 1) * cellSize + cellSize / 2"
          :x2="boardSize - cellSize / 2"
          :y2="(i - 1) * cellSize + cellSize / 2"
          stroke="#8B4513"
          stroke-width="1"
        />
        <!-- 天元和星位 -->
        <circle
          v-for="point in starPoints"
          :key="`star-${point.x}-${point.y}`"
          :cx="point.x * cellSize + cellSize / 2"
          :cy="point.y * cellSize + cellSize / 2"
          r="3"
          fill="#8B4513"
        />
      </svg>

      <!-- 棋子 -->
      <div
        v-for="(row, y) in board"
        :key="`row-${y}`"
        class="board-row"
      >
        <div
          v-for="(cell, x) in row"
          :key="`cell-${x}-${y}`"
          class="board-cell"
          :style="{
            width: cellSize + 'px',
            height: cellSize + 'px'
          }"
          @click="handleCellClick(x, y)"
        >
          <!-- 棋子 -->
          <div
            v-if="cell !== 0"
            class="piece"
            :class="{
              'piece-human': cell === 1,
              'piece-ai': cell === 2,
              'piece-last': lastMove && lastMove.x === x && lastMove.y === y
            }"
          >
            <!-- 最后一步标记 -->
            <div
              v-if="lastMove && lastMove.x === x && lastMove.y === y"
              class="last-move-marker"
            ></div>
          </div>
          
          <!-- 悬停预览 -->
          <div
            v-if="cell === 0 && canMove && currentPlayer === 1"
            class="piece-preview"
            :class="{ 'piece-preview-visible': hoveredCell?.x === x && hoveredCell?.y === y }"
          ></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { Player, type Position } from '../types/game'

// Props
interface Props {
  board: number[][]
  currentPlayer: Player
  lastMove?: Position
  canMove?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  canMove: true
})

// Emits
interface Emits {
  (e: 'move', x: number, y: number): void
}

const emit = defineEmits<Emits>()

// 响应式数据
const boardRef = ref<HTMLElement>()
const hoveredCell = ref<Position | null>(null)
const cellSize = ref(30)
const boardSize = computed(() => cellSize.value * 15)

// 星位坐标（天元和四个角的星位）
const starPoints = [
  { x: 7, y: 7 },   // 天元
  { x: 3, y: 3 },   // 左上
  { x: 11, y: 3 },  // 右上
  { x: 3, y: 11 },  // 左下
  { x: 11, y: 11 }  // 右下
]

// 处理点击事件
function handleCellClick(x: number, y: number) {
  if (!props.canMove || props.currentPlayer !== Player.HUMAN) return
  if (props.board[y][x] !== Player.NONE) return
  
  emit('move', x, y)
}

// 处理鼠标移动
function handleMouseMove(event: MouseEvent) {
  if (!boardRef.value || !props.canMove || props.currentPlayer !== Player.HUMAN) {
    hoveredCell.value = null
    return
  }

  const rect = boardRef.value.getBoundingClientRect()
  const x = Math.floor((event.clientX - rect.left) / cellSize.value)
  const y = Math.floor((event.clientY - rect.top) / cellSize.value)

  if (x >= 0 && x < 15 && y >= 0 && y < 15 && props.board[y][x] === Player.NONE) {
    hoveredCell.value = { x, y }
  } else {
    hoveredCell.value = null
  }
}

// 处理鼠标离开
function handleMouseLeave() {
  hoveredCell.value = null
}

// 响应式调整棋盘大小
function updateBoardSize() {
  if (!boardRef.value) return
  
  const container = boardRef.value.parentElement
  if (!container) return
  
  const containerWidth = container.clientWidth
  const containerHeight = container.clientHeight
  const maxSize = Math.min(containerWidth, containerHeight) - 40
  
  cellSize.value = Math.max(20, Math.min(40, Math.floor(maxSize / 15)))
}

// 生命周期
onMounted(() => {
  updateBoardSize()
  window.addEventListener('resize', updateBoardSize)
  
  if (boardRef.value) {
    boardRef.value.addEventListener('mousemove', handleMouseMove)
    boardRef.value.addEventListener('mouseleave', handleMouseLeave)
  }
})

onUnmounted(() => {
  window.removeEventListener('resize', updateBoardSize)
  
  if (boardRef.value) {
    boardRef.value.removeEventListener('mousemove', handleMouseMove)
    boardRef.value.removeEventListener('mouseleave', handleMouseLeave)
  }
})
</script>

<style scoped>
.board-container {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 20px;
  background: linear-gradient(135deg, #DEB887 0%, #D2B48C 100%);
  border-radius: 10px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.board {
  position: relative;
  background: #F5DEB3;
  border: 2px solid #8B4513;
  border-radius: 5px;
  cursor: pointer;
  user-select: none;
}

.board-grid {
  position: absolute;
  top: 0;
  left: 0;
  pointer-events: none;
}

.board-row {
  display: flex;
}

.board-cell {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

.piece {
  width: 85%;
  height: 85%;
  border-radius: 50%;
  position: relative;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
  transition: all 0.2s ease;
}

.piece-human {
  background: radial-gradient(circle at 30% 30%, #ffffff, #000000);
  border: 1px solid #333;
}

.piece-ai {
  background: radial-gradient(circle at 30% 30%, #ffffff, #ffffff);
  border: 1px solid #ccc;
}

.piece-last {
  box-shadow: 0 0 0 2px #ff4444, 0 2px 4px rgba(0, 0, 0, 0.3);
}

.last-move-marker {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 6px;
  height: 6px;
  background: #ff4444;
  border-radius: 50%;
}

.piece-preview {
  width: 85%;
  height: 85%;
  border-radius: 50%;
  background: radial-gradient(circle at 30% 30%, #ffffff, #000000);
  opacity: 0;
  transition: opacity 0.2s ease;
  border: 1px solid #333;
}

.piece-preview-visible {
  opacity: 0.5;
}

.board:hover .piece-preview-visible {
  opacity: 0.7;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .board-container {
    padding: 10px;
  }
}

@media (max-width: 480px) {
  .board-container {
    padding: 5px;
  }
}
</style>