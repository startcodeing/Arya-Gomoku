<template>
  <div class="bg-white rounded-lg shadow-md p-6">
    <div class="flex items-center justify-between mb-4">
      <h3 class="text-lg font-semibold text-gray-800 flex items-center">
        <BarChart3 class="h-5 w-5 mr-2 text-emerald-500" />
        AI分析
      </h3>
      <div class="flex items-center space-x-2">
        <button
          @click="toggleAutoAnalysis"
          class="px-3 py-1 rounded-full text-xs font-medium transition-colors duration-200"
          :class="autoAnalysis ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-600'"
        >
          {{ autoAnalysis ? '自动分析' : '手动分析' }}
        </button>
        <button
          v-if="hasAnalysisData"
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

    <div v-if="!hasAnalysisData" class="text-center py-8 text-gray-500">
      <TrendingUp class="h-12 w-12 mx-auto mb-3 text-gray-300" />
      <p>暂无分析数据</p>
      <p class="text-sm mt-1">AI落子后将显示分析结果</p>
    </div>

    <div v-else class="space-y-4">
      <!-- 最新分析概览 -->
      <div class="bg-gradient-to-r from-emerald-50 to-blue-50 rounded-lg p-4 border border-emerald-200">
        <div class="flex items-center justify-between mb-3">
          <div class="flex items-center">
            <Brain class="h-5 w-5 mr-2 text-emerald-600" />
            <span class="font-medium text-emerald-800">最新AI分析</span>
          </div>
          <div class="text-sm text-emerald-600">
            第{{ lastAnalysis?.moveNumber }}步
          </div>
        </div>
        
        <div class="grid grid-cols-2 gap-4">
          <div class="bg-white rounded-lg p-3 shadow-sm">
            <div class="flex items-center justify-between mb-1">
              <span class="text-sm text-gray-600">置信度</span>
              <div
                class="w-3 h-3 rounded-full"
                :class="getConfidenceColor(lastAnalysis?.confidence || 0)"
              ></div>
            </div>
            <div class="text-lg font-bold text-emerald-600">
              {{ formatConfidence(lastAnalysis?.confidence || 0) }}
            </div>
          </div>
          
          <div class="bg-white rounded-lg p-3 shadow-sm">
            <div class="flex items-center justify-between mb-1">
              <span class="text-sm text-gray-600">思考时间</span>
              <Clock class="h-3 w-3 text-gray-400" />
            </div>
            <div class="text-lg font-bold text-emerald-600">
              {{ formatThinkingTime(lastAnalysis?.thinkingTime || 0) }}
            </div>
          </div>
        </div>
        
        <div v-if="lastAnalysis?.reasoning" class="mt-3 p-3 bg-white rounded-lg shadow-sm">
          <div class="text-sm text-gray-600 mb-1">AI推理过程</div>
          <div class="text-sm text-gray-800">{{ lastAnalysis.reasoning }}</div>
        </div>
      </div>

      <!-- 局面评估 -->
      <div class="bg-gray-50 rounded-lg p-4">
        <div class="flex items-center mb-3">
          <Target class="h-5 w-5 mr-2 text-blue-500" />
          <span class="font-medium text-gray-700">局面评估</span>
        </div>
        
        <div class="space-y-3">
          <!-- 优势评估条 -->
          <div class="space-y-2">
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">局面优势</span>
              <span class="font-medium" :class="getAdvantageColor()">
                {{ getAdvantageText() }}
              </span>
            </div>
            <div class="w-full bg-gray-200 rounded-full h-2">
              <div
                class="h-2 rounded-full transition-all duration-500"
                :class="getAdvantageBarColor()"
                :style="{ width: getAdvantagePercentage() + '%' }"
              ></div>
            </div>
          </div>
          
          <!-- 威胁等级 -->
          <div class="grid grid-cols-2 gap-3">
            <div class="text-center p-2 bg-white rounded-lg">
              <div class="text-sm text-gray-600 mb-1">攻击威胁</div>
              <div class="text-lg font-bold" :class="getThreatColor(currentAnalysis.attackThreat)">
                {{ getThreatLevel(currentAnalysis.attackThreat) }}
              </div>
            </div>
            <div class="text-center p-2 bg-white rounded-lg">
              <div class="text-sm text-gray-600 mb-1">防守压力</div>
              <div class="text-lg font-bold" :class="getThreatColor(currentAnalysis.defenseThreat)">
                {{ getThreatLevel(currentAnalysis.defenseThreat) }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 展开的详细分析 -->
      <div v-if="isExpanded" class="space-y-4 border-t pt-4">
        <!-- 候选位置分析 -->
        <div class="bg-purple-50 rounded-lg p-4">
          <div class="flex items-center mb-3">
            <MapPin class="h-5 w-5 mr-2 text-purple-500" />
            <span class="font-medium text-purple-700">候选位置分析</span>
          </div>
          
          <div class="space-y-2">
            <div
              v-for="(candidate, index) in getCandidatePositions()"
              :key="index"
              class="flex items-center justify-between p-2 bg-white rounded-lg"
            >
              <div class="flex items-center space-x-3">
                <div class="w-6 h-6 bg-purple-100 rounded-full flex items-center justify-center text-xs font-medium text-purple-700">
                  {{ index + 1 }}
                </div>
                <span class="text-sm font-medium">{{ formatPosition(candidate.x, candidate.y) }}</span>
                <div
                  class="w-3 h-3 rounded-full"
                  :class="getScoreColor(candidate.score)"
                ></div>
              </div>
              <div class="text-sm font-medium text-purple-600">
                {{ candidate.score.toFixed(1) }}
              </div>
            </div>
          </div>
        </div>

        <!-- 策略分析 -->
        <div class="bg-yellow-50 rounded-lg p-4">
          <div class="flex items-center mb-3">
            <Lightbulb class="h-5 w-5 mr-2 text-yellow-500" />
            <span class="font-medium text-yellow-700">策略分析</span>
          </div>
          
          <div class="space-y-2 text-sm">
            <div class="flex justify-between">
              <span class="text-gray-600">主要策略:</span>
              <span class="font-medium text-yellow-700">{{ getMainStrategy() }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-600">风险评估:</span>
              <span class="font-medium" :class="getRiskColor()">{{ getRiskLevel() }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-600">建议行动:</span>
              <span class="font-medium text-yellow-700">{{ getRecommendedAction() }}</span>
            </div>
          </div>
        </div>

        <!-- 历史趋势 -->
        <div class="bg-indigo-50 rounded-lg p-4">
          <div class="flex items-center mb-3">
            <TrendingUp class="h-5 w-5 mr-2 text-indigo-500" />
            <span class="font-medium text-indigo-700">分析趋势</span>
          </div>
          
          <div class="space-y-2">
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">平均置信度</span>
              <span class="font-medium text-indigo-600">{{ getAverageConfidence() }}</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">分析准确率</span>
              <span class="font-medium text-indigo-600">{{ getAnalysisAccuracy() }}</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">策略一致性</span>
              <span class="font-medium text-indigo-600">{{ getStrategyConsistency() }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="flex space-x-2 pt-2">
        <button
          @click="refreshAnalysis"
          :disabled="isAnalyzing"
          class="flex-1 px-3 py-2 bg-emerald-100 hover:bg-emerald-200 text-emerald-700 rounded-lg transition-colors duration-200 disabled:opacity-50 disabled:cursor-not-allowed text-sm"
        >
          <RefreshCw class="h-4 w-4 inline mr-1" :class="{ 'animate-spin': isAnalyzing }" />
          {{ isAnalyzing ? '分析中...' : '刷新分析' }}
        </button>
        <button
          @click="exportAnalysis"
          class="flex-1 px-3 py-2 bg-blue-100 hover:bg-blue-200 text-blue-700 rounded-lg transition-colors duration-200 text-sm"
        >
          <Download class="h-4 w-4 inline mr-1" />
          导出分析
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import {
  BarChart3,
  ChevronDown,
  TrendingUp,
  Brain,
  Clock,
  Target,
  MapPin,
  Lightbulb,
  RefreshCw,
  Download
} from 'lucide-vue-next'
import { useLLMGameStore } from '../../stores/llmGame'
import { storeToRefs } from 'pinia'
import { llmGameUtils } from '../../services/llmApi'

// 使用store
const llmGameStore = useLLMGameStore()
const {
  moveHistory,
  currentGame
} = storeToRefs(llmGameStore)

// 组件状态
const isExpanded = ref(false)
const autoAnalysis = ref(true)
const isAnalyzing = ref(false)

// 模拟分析数据
const currentAnalysis = ref({
  advantage: 0.2, // -1 到 1，负数表示AI优势，正数表示玩家优势
  attackThreat: 0.6, // 0 到 1
  defenseThreat: 0.4, // 0 到 1
  mainStrategy: 'balanced',
  riskLevel: 'medium'
})

// 计算属性
const hasAnalysisData = computed(() => {
  return moveHistory.value.length > 0
})

const lastAnalysis = computed(() => {
  const aiMoves = moveHistory.value.filter(move => move.player === 2)
  if (aiMoves.length === 0) return null
  
  const lastAIMove = aiMoves[aiMoves.length - 1]
  return {
    moveNumber: moveHistory.value.length,
    confidence: lastAIMove.confidence || 0.8,
    thinkingTime: lastAIMove.thinkingTime || 2000,
    reasoning: lastAIMove.reasoning || '选择此位置可以形成有效的攻防平衡，同时限制对手的发展空间。'
  }
})

// 切换展开状态
function toggleExpanded() {
  isExpanded.value = !isExpanded.value
}

// 切换自动分析
function toggleAutoAnalysis() {
  autoAnalysis.value = !autoAnalysis.value
}

// 格式化置信度
function formatConfidence(confidence: number): string {
  return llmGameUtils.formatConfidence(confidence)
}

// 获取置信度颜色
function getConfidenceColor(confidence: number): string {
  return llmGameUtils.getConfidenceColor(confidence)
}

// 格式化思考时间
function formatThinkingTime(time: number): string {
  return `${(time / 1000).toFixed(1)}s`
}

// 获取优势文本
function getAdvantageText(): string {
  const advantage = currentAnalysis.value.advantage
  if (advantage > 0.3) return '玩家大优'
  if (advantage > 0.1) return '玩家小优'
  if (advantage > -0.1) return '势均力敌'
  if (advantage > -0.3) return 'AI小优'
  return 'AI大优'
}

// 获取优势颜色
function getAdvantageColor(): string {
  const advantage = currentAnalysis.value.advantage
  if (advantage > 0.1) return 'text-blue-600'
  if (advantage > -0.1) return 'text-gray-600'
  return 'text-purple-600'
}

// 获取优势进度条颜色
function getAdvantageBarColor(): string {
  const advantage = currentAnalysis.value.advantage
  if (advantage > 0.1) return 'bg-blue-500'
  if (advantage > -0.1) return 'bg-gray-400'
  return 'bg-purple-500'
}

// 获取优势百分比
function getAdvantagePercentage(): number {
  const advantage = currentAnalysis.value.advantage
  return Math.abs(advantage) * 50 + 50
}

// 获取威胁等级
function getThreatLevel(threat: number): string {
  if (threat > 0.8) return '极高'
  if (threat > 0.6) return '高'
  if (threat > 0.4) return '中'
  if (threat > 0.2) return '低'
  return '极低'
}

// 获取威胁颜色
function getThreatColor(threat: number): string {
  if (threat > 0.8) return 'text-red-600'
  if (threat > 0.6) return 'text-orange-600'
  if (threat > 0.4) return 'text-yellow-600'
  if (threat > 0.2) return 'text-green-600'
  return 'text-gray-600'
}

// 获取候选位置
function getCandidatePositions() {
  // 模拟候选位置数据
  return [
    { x: 7, y: 7, score: 8.5 },
    { x: 6, y: 8, score: 7.2 },
    { x: 8, y: 6, score: 6.9 },
    { x: 7, y: 6, score: 6.1 },
    { x: 9, y: 7, score: 5.8 }
  ].slice(0, 3)
}

// 格式化位置
function formatPosition(x: number, y: number): string {
  const letters = 'ABCDEFGHIJKLMNO'
  return `${letters[x]}${y + 1}`
}

// 获取分数颜色
function getScoreColor(score: number): string {
  if (score > 8) return 'bg-green-500'
  if (score > 6) return 'bg-yellow-500'
  return 'bg-red-500'
}

// 获取主要策略
function getMainStrategy(): string {
  const strategies = {
    aggressive: '积极进攻',
    defensive: '稳健防守',
    balanced: '攻防平衡',
    opportunistic: '伺机而动'
  }
  return strategies[currentAnalysis.value.mainStrategy as keyof typeof strategies] || '未知'
}

// 获取风险等级
function getRiskLevel(): string {
  const levels = {
    low: '低风险',
    medium: '中等风险',
    high: '高风险'
  }
  return levels[currentAnalysis.value.riskLevel as keyof typeof levels] || '未知'
}

// 获取风险颜色
function getRiskColor(): string {
  const risk = currentAnalysis.value.riskLevel
  if (risk === 'high') return 'text-red-600'
  if (risk === 'medium') return 'text-yellow-600'
  return 'text-green-600'
}

// 获取建议行动
function getRecommendedAction(): string {
  const advantage = currentAnalysis.value.advantage
  const attackThreat = currentAnalysis.value.attackThreat
  
  if (attackThreat > 0.7) return '优先防守'
  if (advantage > 0.2) return '扩大优势'
  if (advantage < -0.2) return '寻找机会'
  return '稳步推进'
}

// 获取平均置信度
function getAverageConfidence(): string {
  const aiMoves = moveHistory.value.filter(move => move.player === 2 && move.confidence)
  if (aiMoves.length === 0) return '-'
  
  const total = aiMoves.reduce((sum, move) => sum + (move.confidence || 0), 0)
  const avg = total / aiMoves.length
  return formatConfidence(avg)
}

// 获取分析准确率
function getAnalysisAccuracy(): string {
  // 模拟准确率计算
  return '87.5%'
}

// 获取策略一致性
function getStrategyConsistency(): string {
  // 模拟一致性计算
  return '良好'
}

// 刷新分析
async function refreshAnalysis() {
  if (isAnalyzing.value) return
  
  isAnalyzing.value = true
  
  try {
    // 模拟分析过程
    await new Promise(resolve => setTimeout(resolve, 2000))
    
    // 更新分析数据
    currentAnalysis.value = {
      advantage: (Math.random() - 0.5) * 0.8,
      attackThreat: Math.random() * 0.8 + 0.1,
      defenseThreat: Math.random() * 0.8 + 0.1,
      mainStrategy: ['aggressive', 'defensive', 'balanced', 'opportunistic'][Math.floor(Math.random() * 4)],
      riskLevel: ['low', 'medium', 'high'][Math.floor(Math.random() * 3)]
    }
  } catch (error) {
    console.error('分析失败:', error)
  } finally {
    isAnalyzing.value = false
  }
}

// 导出分析
function exportAnalysis() {
  const data = {
    gameId: currentGame.value?.id || 'unknown',
    timestamp: new Date().toISOString(),
    analysis: {
      lastMove: lastAnalysis.value,
      currentAnalysis: currentAnalysis.value,
      candidates: getCandidatePositions(),
      statistics: {
        averageConfidence: getAverageConfidence(),
        analysisAccuracy: getAnalysisAccuracy(),
        strategyConsistency: getStrategyConsistency()
      }
    }
  }
  
  const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `gomoku-analysis-${new Date().toISOString().slice(0, 10)}.json`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
}
</script>