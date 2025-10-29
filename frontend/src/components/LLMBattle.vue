<template>
  <div class="llm-battle-container">
    <div class="game-header">
      <button @click="backToHome" class="back-button">
        ← 返回首页
      </button>
      
      <div class="game-info">
        <h1>与AI对战</h1>
        <div class="game-status">
          <span v-if="isGameActive" class="status playing">
            🧠 AI对战进行中
          </span>
          <span v-else class="status waiting">
            🎯 准备开始
          </span>
        </div>
      </div>
    </div>

    <!-- 模型选择区域 -->
    <div v-if="!isGameActive" class="model-selection">
      <div class="model-selector-card">
        <h2>选择AI模型</h2>
        <div v-if="isLoading" class="loading-state">
          <div class="spinner"></div>
          <span>加载模型列表...</span>
        </div>
        <div v-else-if="availableModels.length === 0" class="empty-state">
          <p>暂无可用模型</p>
          <button @click="loadAvailableModels" class="btn btn-primary">重新加载</button>
        </div>
        <div v-else class="models-grid">
          <div
            v-for="model in availableModels"
            :key="model.id"
            class="model-card"
            :class="{ 'selected': selectedModel?.id === model.id }"
            @click="selectModel(model)"
          >
            <div class="model-icon">🤖</div>
            <h3>{{ model.name }}</h3>
            <p class="model-provider">{{ model.provider }}</p>
            <p v-if="model.description" class="model-description">{{ model.description }}</p>
          </div>
        </div>
        
        <div v-if="selectedModel" class="model-config">
          <h3>模型配置</h3>
          <div class="config-item">
            <label>温度 (Temperature):</label>
            <input
              type="range"
              min="0"
              max="1"
              step="0.1"
              v-model.number="modelConfig.temperature"
              class="config-slider"
            />
            <span>{{ modelConfig.temperature }}</span>
          </div>
          <div class="config-item">
            <label>最大令牌数:</label>
            <input
              type="number"
              min="100"
              max="2000"
              v-model.number="modelConfig.maxTokens"
              class="config-input"
            />
          </div>
        </div>

        <button
          @click="startNewGame"
          :disabled="!selectedModel || isLoading"
          class="btn btn-primary start-game-btn"
        >
          <span v-if="isLoading">开始游戏中...</span>
          <span v-else>开始游戏</span>
        </button>
      </div>
    </div>

    <!-- 游戏区域 -->
    <div v-if="isGameActive" class="game-content">
      <div class="board-section">
        <Board
          :board="convertedBoard"
          :current-player="convertedCurrentPlayer"
          :last-move="lastMove"
          :can-move="canMakeMove"
          @move="handlePlayerMove"
        />
      </div>

      <div class="control-section">
        <ControlPanel
          :game-status="convertedGameStatus"
          :current-player="convertedCurrentPlayer"
          :last-move="lastMove"
          :move-count="moveCount"
          :is-processing="isLoading"
          :is-ai-thinking="isAiThinking"
          :statistics="statistics"
          :show-undo-button="false"
          :show-hint-button="false"
          :show-mode-selector="false"
          @restart="handleRestart"
        />
        
        <!-- LLM特有的信息面板 -->
        <div class="llm-info-panel">
          <div class="model-info">
            <h3>当前模型</h3>
            <div class="current-model">
              <span class="model-name">{{ selectedModel?.name }}</span>
              <span class="model-provider">{{ selectedModel?.provider }}</span>
            </div>
          </div>
          
          <div v-if="isAiThinking" class="ai-thinking">
            <div class="thinking-indicator">
              <div class="thinking-dots">
                <span></span>
                <span></span>
                <span></span>
              </div>
              <span>AI正在思考最佳落子...</span>
            </div>
          </div>
          
          <div class="game-actions">
            <button @click="handleEndGame" class="btn btn-secondary">
              结束游戏
            </button>
            <button @click="handleSwitchModel" class="btn btn-info">
              切换模型
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 错误提示 -->
    <div v-if="error" class="error-toast" @click="clearError">
      <div class="error-content">
        <span class="error-icon">⚠️</span>
        <span class="error-text">{{ error }}</span>
        <button class="error-close">&times;</button>
      </div>
    </div>

    <!-- 游戏结果弹窗 -->
    <div v-if="showGameResult" class="game-result-overlay" @click="closeGameResult">
      <div class="game-result-modal" @click.stop>
        <div class="result-header">
          <h2>游戏结束</h2>
          <button @click="closeGameResult" class="close-btn">&times;</button>
        </div>
        
        <div class="result-content">
          <div class="result-icon">
            <span v-if="gameStatus === 'human_win'">🎉</span>
            <span v-else-if="gameStatus === 'ai_win'">🤖</span>
            <span v-else>🤝</span>
          </div>
          
          <div class="result-message">
            <h3 v-if="gameStatus === 'human_win'">恭喜获胜！</h3>
            <h3 v-else-if="gameStatus === 'ai_win'">AI获胜</h3>
            <h3 v-else>平局</h3>
            
            <p v-if="gameStatus === 'human_win'">您成功击败了AI模型！</p>
            <p v-else-if="gameStatus === 'ai_win'">AI模型展现了强大的棋力</p>
            <p v-else>双方棋力相当，难分胜负</p>
          </div>
          
          <div class="result-stats">
            <div class="stat-item">
              <span class="label">总步数:</span>
              <span class="value">{{ moveCount }}</span>
            </div>
            <div class="stat-item">
              <span class="label">使用模型:</span>
              <span class="value">{{ selectedModel?.name }}</span>
            </div>
          </div>
          
          <div class="result-actions">
            <button @click="handleRestart" class="btn btn-primary">再来一局</button>
            <button @click="handleSwitchModel" class="btn btn-secondary">切换模型</button>
            <button @click="backToHome" class="btn btn-info">返回首页</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'
import Board from './Board.vue'
import ControlPanel from './ControlPanel.vue'
import { useLLMGameStore } from '../stores/llmGame'
import { GameStatus, Player } from '../types/game'

const router = useRouter()
const llmGameStore = useLLMGameStore()

// 从store获取状态
const {
  gameState,
  currentSession,
  availableModels,
  selectedModel,
  modelConfig,
  isLoading,
  isAiThinking,
  error,
  moveHistory,
  statistics,
  board,
  gameStatus,
  currentPlayer,
  lastMove,
  moveCount,
  canMakeMove,
  isGameActive
} = storeToRefs(llmGameStore)

// 从store获取方法
const {
  loadAvailableModels,
  startNewGame,
  makePlayerMove,
  restartGame,
  selectModel,
  updateModelConfig,
  clearError,
  initialize
} = llmGameStore

// 本地状态
const showGameResult = ref(false)

// 计算属性
const GameStatus = computed(() => ({
  PLAYING: 'playing',
  HUMAN_WIN: 'human_win',
  AI_WIN: 'ai_win',
  DRAW: 'draw'
}))

// 计算属性
const GameStatusEnum = computed(() => ({
  PLAYING: 'playing',
  HUMAN_WIN: 'human_win',
  AI_WIN: 'ai_win',
  DRAW: 'draw'
}))

// 转换LLM游戏状态为通用游戏状态
const convertedGameStatus = computed(() => {
  // gameStatus.value 已经是 GameStatus 枚举，直接返回
  return gameStatus.value || GameStatus.PLAYING
})

// 转换当前玩家为Player枚举
const convertedCurrentPlayer = computed(() => {
  switch (currentPlayer.value) {
    case 1:
      return Player.HUMAN
    case 2:
      return Player.AI
    default:
      return Player.HUMAN
  }
})

// 转换棋盘数据
const convertedBoard = computed(() => {
  return board.value.map(row => 
    row.map(cell => {
      switch (cell) {
        case 1:
          return Player.HUMAN
        case 2:
          return Player.AI
        default:
          return Player.NONE
      }
    })
  )
})

// 监听游戏状态变化
watch(gameStatus, (newStatus) => {
  if (newStatus !== 'playing') {
    setTimeout(() => {
      showGameResult.value = true
    }, 1000) // 延迟显示结果，让用户看到最后一步
  }
})

// 方法
function backToHome() {
  router.push('/')
}

async function handlePlayerMove(x: number, y: number) {
  await makePlayerMove(x, y)
}

function handleRestart() {
  showGameResult.value = false
  restartGame()
  if (selectedModel.value) {
    startNewGame()
  }
}

function handleEndGame() {
  const confirmed = confirm('确定要结束当前游戏吗？')
  if (confirmed) {
    restartGame()
  }
}

function handleSwitchModel() {
  showGameResult.value = false
  restartGame()
}

function closeGameResult() {
  showGameResult.value = false
}

// 生命周期
onMounted(async () => {
  await initialize()
})
</script>

<style scoped>
.llm-battle-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.game-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 30px;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  border-radius: 15px;
  padding: 20px;
}

.back-button {
  background: rgba(255, 255, 255, 0.2);
  border: none;
  color: white;
  padding: 10px 20px;
  border-radius: 8px;
  cursor: pointer;
  font-size: 16px;
  transition: all 0.3s ease;
}

.back-button:hover {
  background: rgba(255, 255, 255, 0.3);
  transform: translateX(-2px);
}

.game-info h1 {
  color: white;
  font-size: 2.5rem;
  margin: 0;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.3);
}

.game-status {
  margin-top: 10px;
}

.status {
  padding: 8px 16px;
  border-radius: 20px;
  font-weight: 500;
  font-size: 14px;
}

.status.playing {
  background: rgba(34, 197, 94, 0.2);
  color: #16a34a;
  border: 1px solid rgba(34, 197, 94, 0.3);
}

.status.waiting {
  background: rgba(251, 191, 36, 0.2);
  color: #d97706;
  border: 1px solid rgba(251, 191, 36, 0.3);
}

.model-selection {
  max-width: 800px;
  margin: 0 auto;
}

.model-selector-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  padding: 30px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
}

.model-selector-card h2 {
  text-align: center;
  color: #374151;
  margin-bottom: 30px;
  font-size: 1.8rem;
}

.loading-state, .empty-state {
  text-align: center;
  padding: 40px;
  color: #6b7280;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #f3f4f6;
  border-top: 4px solid #3b82f6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 20px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.models-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.model-card {
  background: white;
  border: 2px solid #e5e7eb;
  border-radius: 15px;
  padding: 20px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
}

.model-card:hover {
  border-color: #3b82f6;
  transform: translateY(-2px);
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.1);
}

.model-card.selected {
  border-color: #3b82f6;
  background: #eff6ff;
}

.model-icon {
  font-size: 2rem;
  margin-bottom: 10px;
}

.model-card h3 {
  color: #374151;
  margin: 10px 0 5px;
  font-size: 1.2rem;
}

.model-provider {
  color: #6b7280;
  font-size: 0.9rem;
  margin-bottom: 10px;
}

.model-description {
  color: #9ca3af;
  font-size: 0.8rem;
  line-height: 1.4;
}

.model-config {
  background: #f9fafb;
  border-radius: 10px;
  padding: 20px;
  margin-bottom: 20px;
}

.model-config h3 {
  color: #374151;
  margin-bottom: 15px;
}

.config-item {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 15px;
}

.config-item label {
  min-width: 120px;
  color: #6b7280;
  font-size: 0.9rem;
}

.config-slider {
  flex: 1;
}

.config-input {
  width: 100px;
  padding: 5px 10px;
  border: 1px solid #d1d5db;
  border-radius: 5px;
}

.start-game-btn {
  width: 100%;
  padding: 15px;
  font-size: 1.1rem;
  font-weight: 600;
}

.game-content {
  display: grid;
  grid-template-columns: 1fr 400px;
  gap: 30px;
  max-width: 1400px;
  margin: 0 auto;
}

.board-section {
  display: flex;
  justify-content: center;
  align-items: center;
}

.control-section {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.llm-info-panel {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 15px;
  padding: 20px;
}

.model-info h3 {
  color: #374151;
  margin-bottom: 10px;
  font-size: 1.1rem;
}

.current-model {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.model-name {
  font-weight: 600;
  color: #1f2937;
}

.model-provider {
  font-size: 0.9rem;
  color: #6b7280;
}

.ai-thinking {
  margin: 20px 0;
  padding: 15px;
  background: #eff6ff;
  border-radius: 10px;
  border-left: 4px solid #3b82f6;
}

.thinking-indicator {
  display: flex;
  align-items: center;
  gap: 10px;
  color: #1e40af;
}

.thinking-dots {
  display: flex;
  gap: 4px;
}

.thinking-dots span {
  width: 8px;
  height: 8px;
  background: #3b82f6;
  border-radius: 50%;
  animation: thinking 1.4s ease-in-out infinite both;
}

.thinking-dots span:nth-child(1) { animation-delay: -0.32s; }
.thinking-dots span:nth-child(2) { animation-delay: -0.16s; }

@keyframes thinking {
  0%, 80%, 100% {
    transform: scale(0);
  }
  40% {
    transform: scale(1);
  }
}

.game-actions {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-top: 20px;
}

.btn {
  padding: 10px 20px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.3s ease;
  text-align: center;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-primary {
  background: #3b82f6;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: #2563eb;
  transform: translateY(-1px);
}

.btn-secondary {
  background: #6b7280;
  color: white;
}

.btn-secondary:hover:not(:disabled) {
  background: #4b5563;
}

.btn-info {
  background: #06b6d4;
  color: white;
}

.btn-info:hover:not(:disabled) {
  background: #0891b2;
}

.error-toast {
  position: fixed;
  top: 20px;
  right: 20px;
  background: #fef2f2;
  border: 1px solid #fecaca;
  border-radius: 10px;
  padding: 15px;
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.1);
  cursor: pointer;
  z-index: 1000;
  animation: slideIn 0.3s ease;
}

@keyframes slideIn {
  from {
    transform: translateX(100%);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}

.error-content {
  display: flex;
  align-items: center;
  gap: 10px;
  color: #dc2626;
}

.error-close {
  background: none;
  border: none;
  color: #dc2626;
  cursor: pointer;
  font-size: 1.2rem;
  margin-left: auto;
}

.game-result-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.game-result-modal {
  background: white;
  border-radius: 20px;
  padding: 30px;
  max-width: 500px;
  width: 90%;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.2);
  animation: slideUp 0.3s ease;
}

@keyframes slideUp {
  from {
    transform: translateY(50px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

.result-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.result-header h2 {
  color: #374151;
  margin: 0;
}

.close-btn {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #6b7280;
}

.result-content {
  text-align: center;
}

.result-icon {
  font-size: 4rem;
  margin-bottom: 20px;
}

.result-message h3 {
  color: #374151;
  margin-bottom: 10px;
  font-size: 1.5rem;
}

.result-message p {
  color: #6b7280;
  margin-bottom: 20px;
}

.result-stats {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 15px;
  margin-bottom: 30px;
  padding: 20px;
  background: #f9fafb;
  border-radius: 10px;
}

.stat-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stat-item .label {
  color: #6b7280;
  font-size: 0.9rem;
}

.stat-item .value {
  color: #374151;
  font-weight: 600;
}

.result-actions {
  display: flex;
  gap: 10px;
  justify-content: center;
}

@media (max-width: 1024px) {
  .game-content {
    grid-template-columns: 1fr;
    gap: 20px;
  }
  
  .control-section {
    order: -1;
  }
}

@media (max-width: 768px) {
  .llm-battle-container {
    padding: 10px;
  }
  
  .game-header {
    flex-direction: column;
    gap: 15px;
    text-align: center;
  }
  
  .game-info h1 {
    font-size: 2rem;
  }
  
  .models-grid {
    grid-template-columns: 1fr;
  }
  
  .result-actions {
    flex-direction: column;
  }
}
</style>