<template>
  <div class="control-panel">
    <div class="panel-header">
      <h2>五子棋人机对战</h2>
    </div>

    <div class="game-info">
      <div class="status-section">
        <div class="status-item">
          <span class="label">游戏状态:</span>
          <span class="value" :class="getStatusClass(gameStatus)">
            {{ formatGameStatus(gameStatus) }}
          </span>
        </div>
        
        <div class="status-item" v-if="gameStatus === 'playing'">
          <span class="label">当前回合:</span>
          <span class="value" :class="getCurrentPlayerClass(currentPlayer)">
            {{ getCurrentPlayerText(currentPlayer) }}
          </span>
        </div>

        <div class="status-item" v-if="lastMove">
          <span class="label">最后一步:</span>
          <span class="value">
            {{ `(${lastMove.x + 1}, ${lastMove.y + 1})` }}
          </span>
        </div>
      </div>

      <div class="move-counter">
        <span class="label">步数:</span>
        <span class="value">{{ moveCount }}</span>
      </div>
    </div>

    <div class="controls">
      <button 
        class="btn btn-primary"
        @click="handleRestart"
        :disabled="isProcessing"
      >
        <span v-if="isProcessing">重新开始中...</span>
        <span v-else>重新开始</span>
      </button>

      <button 
        class="btn btn-secondary"
        @click="handleUndo"
        :disabled="!canUndo || isProcessing"
        v-if="showUndoButton"
      >
        悔棋
      </button>

      <button 
        class="btn btn-info"
        @click="handleHint"
        :disabled="!canGetHint || isProcessing"
        v-if="showHintButton"
      >
        提示
      </button>
    </div>

    <div class="game-mode" v-if="showModeSelector">
      <div class="mode-section">
        <span class="label">游戏模式:</span>
        <select 
          v-model="selectedMode" 
          @change="handleModeChange"
          :disabled="isProcessing"
          class="mode-select"
        >
          <option value="ai">人机对战</option>
          <option value="pvp" disabled>在线对战 (即将推出)</option>
        </select>
      </div>
    </div>

    <div class="statistics" v-if="showStatistics">
      <h3>游戏统计</h3>
      <div class="stats-grid">
        <div class="stat-item">
          <span class="stat-label">玩家胜利:</span>
          <span class="stat-value">{{ statistics.humanWins }}</span>
        </div>
        <div class="stat-item">
          <span class="stat-label">AI胜利:</span>
          <span class="stat-value">{{ statistics.aiWins }}</span>
        </div>
        <div class="stat-item">
          <span class="stat-label">平局:</span>
          <span class="stat-value">{{ statistics.draws }}</span>
        </div>
        <div class="stat-item">
          <span class="stat-label">总局数:</span>
          <span class="stat-value">{{ statistics.totalGames }}</span>
        </div>
      </div>
    </div>

    <div class="ai-thinking" v-if="isAiThinking">
      <div class="thinking-indicator">
        <div class="spinner"></div>
        <span>AI思考中...</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Player, GameStatus, type Position } from '../types/game'
import { formatGameStatus, getCurrentPlayerText } from '../utils/gameLogic'

// Props
interface Props {
  gameStatus: GameStatus
  currentPlayer: Player
  lastMove?: Position
  moveCount?: number
  isProcessing?: boolean
  isAiThinking?: boolean
  canUndo?: boolean
  canGetHint?: boolean
  showUndoButton?: boolean
  showHintButton?: boolean
  showModeSelector?: boolean
  showStatistics?: boolean
  statistics?: {
    humanWins: number
    aiWins: number
    draws: number
    totalGames: number
  }
}

const props = withDefaults(defineProps<Props>(), {
  moveCount: 0,
  isProcessing: false,
  isAiThinking: false,
  canUndo: false,
  canGetHint: false,
  showUndoButton: false,
  showHintButton: false,
  showModeSelector: true,
  showStatistics: true,
  statistics: () => ({
    humanWins: 0,
    aiWins: 0,
    draws: 0,
    totalGames: 0
  })
})

// Emits
interface Emits {
  (e: 'restart'): void
  (e: 'undo'): void
  (e: 'hint'): void
  (e: 'modeChange', mode: string): void
}

const emit = defineEmits<Emits>()

// 响应式数据
const selectedMode = ref('ai')

// 计算属性
const winRate = computed(() => {
  const total = props.statistics.totalGames
  if (total === 0) return 0
  return Math.round((props.statistics.humanWins / total) * 100)
})

// 方法
function handleRestart() {
  emit('restart')
}

function handleUndo() {
  emit('undo')
}

function handleHint() {
  emit('hint')
}

function handleModeChange() {
  emit('modeChange', selectedMode.value)
}

function getStatusClass(status: GameStatus): string {
  switch (status) {
    case GameStatus.HUMAN_WIN:
      return 'status-win'
    case GameStatus.AI_WIN:
      return 'status-lose'
    case GameStatus.DRAW:
      return 'status-draw'
    case GameStatus.PLAYING:
      return 'status-playing'
    default:
      return ''
  }
}

function getCurrentPlayerClass(player: Player): string {
  switch (player) {
    case Player.HUMAN:
      return 'player-human'
    case Player.AI:
      return 'player-ai'
    default:
      return ''
  }
}
</script>

<style scoped>
.control-panel {
  background: #ffffff;
  border-radius: 10px;
  padding: 20px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  min-width: 280px;
  max-width: 350px;
}

.panel-header h2 {
  margin: 0 0 20px 0;
  color: #333;
  text-align: center;
  font-size: 1.5em;
  font-weight: 600;
}

.game-info {
  margin-bottom: 20px;
}

.status-section {
  margin-bottom: 15px;
}

.status-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  padding: 5px 0;
}

.label {
  font-weight: 500;
  color: #666;
}

.value {
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 4px;
}

.status-playing {
  color: #2196F3;
  background: #E3F2FD;
}

.status-win {
  color: #4CAF50;
  background: #E8F5E8;
}

.status-lose {
  color: #F44336;
  background: #FFEBEE;
}

.status-draw {
  color: #FF9800;
  background: #FFF3E0;
}

.player-human {
  color: #333;
  background: #F5F5F5;
}

.player-ai {
  color: #2196F3;
  background: #E3F2FD;
}

.move-counter {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px;
  background: #F8F9FA;
  border-radius: 5px;
}

.controls {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-bottom: 20px;
}

.btn {
  padding: 10px 15px;
  border: none;
  border-radius: 5px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  text-align: center;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-primary {
  background: #2196F3;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: #1976D2;
}

.btn-secondary {
  background: #6C757D;
  color: white;
}

.btn-secondary:hover:not(:disabled) {
  background: #5A6268;
}

.btn-info {
  background: #17A2B8;
  color: white;
}

.btn-info:hover:not(:disabled) {
  background: #138496;
}

.game-mode {
  margin-bottom: 20px;
}

.mode-section {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.mode-select {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  background: white;
}

.statistics h3 {
  margin: 0 0 15px 0;
  color: #333;
  font-size: 1.1em;
  text-align: center;
}

.stats-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 10px;
  background: #F8F9FA;
  border-radius: 5px;
}

.stat-label {
  font-size: 12px;
  color: #666;
  margin-bottom: 4px;
}

.stat-value {
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.ai-thinking {
  margin-top: 15px;
  padding: 15px;
  background: #E3F2FD;
  border-radius: 5px;
  text-align: center;
}

.thinking-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  color: #2196F3;
  font-weight: 500;
}

.spinner {
  width: 20px;
  height: 20px;
  border: 2px solid #E3F2FD;
  border-top: 2px solid #2196F3;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .control-panel {
    min-width: 250px;
    max-width: 300px;
    padding: 15px;
  }
  
  .panel-header h2 {
    font-size: 1.3em;
  }
  
  .stats-grid {
    grid-template-columns: 1fr;
  }
}
</style>