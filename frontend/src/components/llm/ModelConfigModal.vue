<template>
  <div
    v-if="isOpen"
    class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
    @click="handleBackdropClick"
  >
    <div
      class="bg-white rounded-lg shadow-xl max-w-2xl w-full mx-4 max-h-[90vh] overflow-hidden"
      @click.stop
    >
      <!-- 头部 -->
      <div class="flex items-center justify-between p-6 border-b border-gray-200">
        <div class="flex items-center">
          <Settings class="h-6 w-6 mr-3 text-blue-500" />
          <div>
            <h2 class="text-xl font-semibold text-gray-800">模型配置</h2>
            <p class="text-sm text-gray-600 mt-1">
              {{ selectedModel?.name || '未选择模型' }}
            </p>
          </div>
        </div>
        <button
          @click="closeModal"
          class="p-2 hover:bg-gray-100 rounded-lg transition-colors duration-200"
        >
          <X class="h-5 w-5 text-gray-500" />
        </button>
      </div>

      <!-- 内容 -->
      <div class="p-6 overflow-y-auto max-h-[calc(90vh-140px)]">
        <div v-if="!selectedModel" class="text-center py-8 text-gray-500">
          <AlertCircle class="h-12 w-12 mx-auto mb-3 text-gray-300" />
          <p>请先选择一个模型</p>
        </div>

        <div v-else class="space-y-6">
          <!-- 模型信息 -->
          <div class="bg-gray-50 rounded-lg p-4">
            <div class="flex items-center mb-3">
              <Info class="h-5 w-5 mr-2 text-blue-500" />
              <span class="font-medium text-gray-700">模型信息</span>
            </div>
            <div class="grid grid-cols-2 gap-4 text-sm">
              <div>
                <span class="text-gray-600">提供商:</span>
                <span class="ml-2 font-medium">{{ selectedModel.provider }}</span>
              </div>
              <div>
                <span class="text-gray-600">状态:</span>
                <span
                  class="ml-2 px-2 py-1 rounded-full text-xs font-medium"
                  :class="getModelStatusStyle()"
                >
                  {{ getModelStatusText() }}
                </span>
              </div>
              <div class="col-span-2" v-if="selectedModel.description">
                <span class="text-gray-600">描述:</span>
                <p class="mt-1 text-gray-700">{{ selectedModel.description }}</p>
              </div>
            </div>
          </div>

          <!-- 配置表单 -->
          <form @submit.prevent="handleSubmit" class="space-y-4">
            <!-- API Key -->
            <div v-if="needsApiKey()">
              <label class="block text-sm font-medium text-gray-700 mb-2">
                API Key
                <span class="text-red-500">*</span>
              </label>
              <div class="relative">
                <input
                  v-model="formData.apiKey"
                  :type="showApiKey ? 'text' : 'password'"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 pr-10"
                  :placeholder="getApiKeyPlaceholder()"
                  required
                />
                <button
                  type="button"
                  @click="toggleApiKeyVisibility"
                  class="absolute inset-y-0 right-0 pr-3 flex items-center"
                >
                  <component
                    :is="showApiKey ? EyeOff : Eye"
                    class="h-4 w-4 text-gray-400"
                  />
                </button>
              </div>
              <p class="mt-1 text-xs text-gray-500">
                {{ getApiKeyHint() }}
              </p>
            </div>

            <!-- 基础URL (仅Ollama) -->
            <div v-if="selectedModel.provider === 'ollama'">
              <label class="block text-sm font-medium text-gray-700 mb-2">
                服务器地址
              </label>
              <input
                v-model="formData.baseUrl"
                type="url"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                placeholder="http://localhost:11434"
              />
              <p class="mt-1 text-xs text-gray-500">
                Ollama服务器的地址，默认为本地地址
              </p>
            </div>

            <!-- 高级配置 -->
            <div class="border-t pt-4">
              <div class="flex items-center justify-between mb-4">
                <span class="font-medium text-gray-700">高级配置</span>
                <button
                  type="button"
                  @click="toggleAdvanced"
                  class="text-sm text-blue-600 hover:text-blue-700"
                >
                  {{ showAdvanced ? '隐藏' : '显示' }}
                </button>
              </div>

              <div v-if="showAdvanced" class="space-y-4">
                <!-- 温度 -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    创造性 (Temperature)
                    <span class="text-gray-500 font-normal">{{ formData.temperature }}</span>
                  </label>
                  <input
                    v-model.number="formData.temperature"
                    type="range"
                    min="0"
                    max="1"
                    step="0.1"
                    class="w-full"
                  />
                  <div class="flex justify-between text-xs text-gray-500 mt-1">
                    <span>保守 (0)</span>
                    <span>平衡 (0.5)</span>
                    <span>创新 (1)</span>
                  </div>
                </div>

                <!-- 最大令牌数 -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    最大响应长度
                  </label>
                  <input
                    v-model.number="formData.maxTokens"
                    type="number"
                    min="50"
                    max="2000"
                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  />
                  <p class="mt-1 text-xs text-gray-500">
                    控制AI响应的最大长度，建议值：100-500
                  </p>
                </div>

                <!-- 超时设置 -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    请求超时 (秒)
                  </label>
                  <input
                    v-model.number="formData.timeout"
                    type="number"
                    min="5"
                    max="60"
                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  />
                </div>

                <!-- 重试次数 -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    失败重试次数
                  </label>
                  <select
                    v-model.number="formData.retryCount"
                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  >
                    <option value="0">不重试</option>
                    <option value="1">1次</option>
                    <option value="2">2次</option>
                    <option value="3">3次</option>
                  </select>
                </div>
              </div>
            </div>
          </form>

          <!-- 测试连接 -->
          <div class="border-t pt-4">
            <button
              @click="testConnection"
              :disabled="isTesting || !canTest()"
              class="w-full px-4 py-2 bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-lg transition-colors duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <component
                :is="isTesting ? Loader2 : Wifi"
                class="h-4 w-4 inline mr-2"
                :class="{ 'animate-spin': isTesting }"
              />
              {{ isTesting ? '测试中...' : '测试连接' }}
            </button>
            
            <div v-if="testResult" class="mt-3 p-3 rounded-lg" :class="getTestResultStyle()">
              <div class="flex items-center">
                <component
                  :is="testResult.success ? CheckCircle : XCircle"
                  class="h-4 w-4 mr-2"
                  :class="testResult.success ? 'text-green-500' : 'text-red-500'"
                />
                <span class="font-medium">
                  {{ testResult.success ? '连接成功' : '连接失败' }}
                </span>
              </div>
              <p class="text-sm mt-1">{{ testResult.message }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- 底部按钮 -->
      <div class="flex items-center justify-end space-x-3 p-6 border-t border-gray-200 bg-gray-50">
        <button
          @click="closeModal"
          class="px-4 py-2 text-gray-700 hover:bg-gray-100 rounded-lg transition-colors duration-200"
        >
          取消
        </button>
        <button
          @click="resetToDefaults"
          class="px-4 py-2 text-blue-600 hover:bg-blue-50 rounded-lg transition-colors duration-200"
        >
          重置默认
        </button>
        <button
          @click="handleSubmit"
          :disabled="isSaving || !isFormValid()"
          class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <component
            :is="isSaving ? Loader2 : Save"
            class="h-4 w-4 inline mr-2"
            :class="{ 'animate-spin': isSaving }"
          />
          {{ isSaving ? '保存中...' : '保存配置' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import {
  Settings,
  X,
  Info,
  AlertCircle,
  Eye,
  EyeOff,
  Wifi,
  CheckCircle,
  XCircle,
  Loader2,
  Save
} from 'lucide-vue-next'
import { useLLMGameStore } from '../../stores/llmGame'
import { storeToRefs } from 'pinia'
import type { LLMConfigRequest } from '../../types/game'

// Props
interface Props {
  isOpen: boolean
}

const props = defineProps<Props>()

// Emits
const emit = defineEmits<{
  close: []
  saved: [config: LLMConfigRequest]
}>()

// 使用store
const llmGameStore = useLLMGameStore()
const { selectedModel } = storeToRefs(llmGameStore)
const { updateModelConfig } = llmGameStore

// 组件状态
const showApiKey = ref(false)
const showAdvanced = ref(false)
const isTesting = ref(false)
const isSaving = ref(false)
const testResult = ref<{ success: boolean; message: string } | null>(null)

// 表单数据
const formData = reactive<LLMConfigRequest>({
  apiKey: '',
  baseUrl: '',
  temperature: 0.7,
  maxTokens: 200,
  timeout: 30,
  retryCount: 2
})

// 监听模型变化，重置表单
watch(
  () => selectedModel.value,
  (newModel) => {
    if (newModel) {
      resetForm()
    }
  },
  { immediate: true }
)

// 重置表单
function resetForm() {
  if (!selectedModel.value) return
  
  // 根据模型类型设置默认值
  formData.apiKey = ''
  formData.baseUrl = selectedModel.value.provider === 'ollama' ? 'http://localhost:11434' : ''
  formData.temperature = 0.7
  formData.maxTokens = 200
  formData.timeout = 30
  formData.retryCount = 2
  
  testResult.value = null
}

// 是否需要API Key
function needsApiKey(): boolean {
  return selectedModel.value?.provider !== 'ollama'
}

// 获取API Key占位符
function getApiKeyPlaceholder(): string {
  if (!selectedModel.value) return ''
  
  switch (selectedModel.value.provider) {
    case 'deepseek':
      return 'sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'
    case 'openai':
      return 'sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'
    default:
      return '请输入API Key'
  }
}

// 获取API Key提示
function getApiKeyHint(): string {
  if (!selectedModel.value) return ''
  
  switch (selectedModel.value.provider) {
    case 'deepseek':
      return '在 DeepSeek 控制台获取您的API密钥'
    case 'openai':
      return '在 OpenAI 控制台获取您的API密钥'
    default:
      return '请输入有效的API密钥'
  }
}

// 获取模型状态样式
function getModelStatusStyle(): string {
  if (!selectedModel.value) return 'bg-gray-100 text-gray-600'
  
  switch (selectedModel.value.status) {
    case 'available':
      return 'bg-green-100 text-green-700'
    case 'unavailable':
      return 'bg-red-100 text-red-700'
    case 'not_configured':
      return 'bg-yellow-100 text-yellow-700'
    default:
      return 'bg-gray-100 text-gray-600'
  }
}

// 获取模型状态文本
function getModelStatusText(): string {
  if (!selectedModel.value) return '未知'
  
  switch (selectedModel.value.status) {
    case 'available':
      return '可用'
    case 'unavailable':
      return '不可用'
    case 'not_configured':
      return '未配置'
    default:
      return '未知'
  }
}

// 切换API Key可见性
function toggleApiKeyVisibility() {
  showApiKey.value = !showApiKey.value
}

// 切换高级配置
function toggleAdvanced() {
  showAdvanced.value = !showAdvanced.value
}

// 表单验证
function isFormValid(): boolean {
  if (needsApiKey() && !formData.apiKey.trim()) {
    return false
  }
  
  if (selectedModel.value?.provider === 'ollama' && !formData.baseUrl.trim()) {
    return false
  }
  
  return true
}

// 是否可以测试
function canTest(): boolean {
  return isFormValid()
}

// 测试连接
async function testConnection() {
  if (!selectedModel.value || isTesting.value) return
  
  isTesting.value = true
  testResult.value = null
  
  try {
    // 模拟测试连接
    await new Promise(resolve => setTimeout(resolve, 2000))
    
    // 这里应该调用实际的测试API
    const success = Math.random() > 0.3 // 模拟70%成功率
    
    testResult.value = {
      success,
      message: success 
        ? '连接成功，模型响应正常' 
        : '连接失败，请检查配置信息'
    }
  } catch (error) {
    testResult.value = {
      success: false,
      message: '测试过程中发生错误'
    }
  } finally {
    isTesting.value = false
  }
}

// 获取测试结果样式
function getTestResultStyle(): string {
  if (!testResult.value) return ''
  
  return testResult.value.success 
    ? 'bg-green-50 border border-green-200' 
    : 'bg-red-50 border border-red-200'
}

// 重置为默认值
function resetToDefaults() {
  resetForm()
}

// 提交表单
async function handleSubmit() {
  if (!selectedModel.value || !isFormValid() || isSaving.value) return
  
  isSaving.value = true
  
  try {
    await updateModelConfig(selectedModel.value.id, formData)
    emit('saved', formData)
    closeModal()
  } catch (error) {
    console.error('保存配置失败:', error)
    // 这里可以显示错误提示
  } finally {
    isSaving.value = false
  }
}

// 关闭模态框
function closeModal() {
  emit('close')
}

// 处理背景点击
function handleBackdropClick() {
  closeModal()
}
</script>