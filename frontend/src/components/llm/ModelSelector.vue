<template>
  <div class="bg-white rounded-lg shadow-md p-6">
    <h3 class="text-lg font-semibold text-gray-800 mb-4 flex items-center">
      <Brain class="h-5 w-5 mr-2 text-blue-500" />
      选择AI模型
    </h3>

    <!-- 模型列表 -->
    <div class="space-y-3">
      <div
        v-for="model in availableModels"
        :key="model.name"
        class="border rounded-lg p-3 cursor-pointer transition-all duration-200"
        :class="{
          'border-blue-500 bg-blue-50': selectedModel === model.name,
          'border-gray-200 hover:border-gray-300': selectedModel !== model.name,
          'opacity-50 cursor-not-allowed': model.status !== 'available'
        }"
        @click="handleModelSelect(model)"
      >
        <div class="flex items-center justify-between">
          <div class="flex items-center">
            <div
              class="w-3 h-3 rounded-full mr-3"
              :class="getStatusColor(model.status)"
            ></div>
            <div>
              <div class="font-medium text-gray-800">{{ model.displayName }}</div>
              <div class="text-sm text-gray-500">{{ model.provider }}</div>
            </div>
          </div>
          <div class="text-right">
            <div
              class="text-xs px-2 py-1 rounded-full"
              :class="getStatusBadgeClass(model.status)"
            >
              {{ getStatusText(model.status) }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 配置按钮 -->
    <div class="mt-4 pt-4 border-t border-gray-200">
      <button
        @click="showConfigModal = true"
        class="w-full flex items-center justify-center px-4 py-2 bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-lg transition-colors duration-200"
        :disabled="isLoading"
      >
        <Settings class="h-4 w-4 mr-2" />
        配置模型
      </button>
    </div>

    <!-- 配置模态框 -->
    <ModelConfigModal
      v-if="showConfigModal"
      :model-name="selectedModel"
      @close="showConfigModal = false"
      @updated="handleConfigUpdated"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Brain, Settings } from 'lucide-vue-next'
import { useLLMGameStore } from '../../stores/llmGame'
import { storeToRefs } from 'pinia'
import { llmGameUtils } from '../../services/llmApi'
import type { LLMModel } from '../../types/game'
import ModelConfigModal from './ModelConfigModal.vue'

// 使用store
const llmGameStore = useLLMGameStore()
const { availableModels, selectedModel, isLoading, isGameActive } = storeToRefs(llmGameStore)
const { selectModel } = llmGameStore

// 本地状态
const showConfigModal = ref(false)

// 处理模型选择
function handleModelSelect(model: LLMModel) {
  if (model.status !== 'available') {
    return
  }
  
  if (isGameActive.value) {
    const shouldSwitch = confirm('切换模型将结束当前游戏，确定要切换吗？')
    if (!shouldSwitch) {
      return
    }
    llmGameStore.endGame()
  }
  
  selectModel(model.name)
}

// 获取状态颜色
function getStatusColor(status: string): string {
  switch (status) {
    case 'available':
      return 'bg-green-500'
    case 'not_configured':
      return 'bg-yellow-500'
    case 'unavailable':
      return 'bg-red-500'
    default:
      return 'bg-gray-500'
  }
}

// 获取状态徽章样式
function getStatusBadgeClass(status: string): string {
  switch (status) {
    case 'available':
      return 'bg-green-100 text-green-800'
    case 'not_configured':
      return 'bg-yellow-100 text-yellow-800'
    case 'unavailable':
      return 'bg-red-100 text-red-800'
    default:
      return 'bg-gray-100 text-gray-800'
  }
}

// 获取状态文本
function getStatusText(status: string): string {
  return llmGameUtils.getModelStatusMessage(status)
}

// 处理配置更新
function handleConfigUpdated() {
  showConfigModal.value = false
  // 重新加载模型列表
  llmGameStore.loadAvailableModels()
}
</script>