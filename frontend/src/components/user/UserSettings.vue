<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 导航栏 -->
    <nav class="bg-white shadow-sm border-b">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center space-x-4">
            <router-link to="/user" class="text-gray-500 hover:text-gray-700">
              ← 返回用户中心
            </router-link>
            <h1 class="text-xl font-semibold text-gray-900">个人设置</h1>
          </div>
        </div>
      </div>
    </nav>

    <div class="max-w-4xl mx-auto py-6 sm:px-6 lg:px-8">
      <!-- 个人信息 -->
      <div class="bg-white shadow rounded-lg mb-6">
        <div class="px-6 py-4 border-b border-gray-200">
          <h3 class="text-lg font-medium text-gray-900">个人信息</h3>
        </div>
        <form @submit.prevent="updateProfile" class="p-6 space-y-6">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <label for="username" class="block text-sm font-medium text-gray-700">
                用户名
              </label>
              <input
                id="username"
                v-model="profileForm.username"
                type="text"
                disabled
                class="mt-1 block w-full border border-gray-300 rounded-md px-3 py-2 bg-gray-50 text-gray-500 cursor-not-allowed"
              />
              <p class="mt-1 text-xs text-gray-500">用户名不可修改</p>
            </div>

            <div>
              <label for="email" class="block text-sm font-medium text-gray-700">
                邮箱
              </label>
              <input
                id="email"
                v-model="profileForm.email"
                type="email"
                required
                class="mt-1 block w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-500"
              />
            </div>

            <div>
              <label for="nickname" class="block text-sm font-medium text-gray-700">
                昵称
              </label>
              <input
                id="nickname"
                v-model="profileForm.nickname"
                type="text"
                class="mt-1 block w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-500"
              />
            </div>

            <div>
              <label for="role" class="block text-sm font-medium text-gray-700">
                角色
              </label>
              <input
                id="role"
                :value="getRoleName(userStore.user?.role || '')"
                type="text"
                disabled
                class="mt-1 block w-full border border-gray-300 rounded-md px-3 py-2 bg-gray-50 text-gray-500 cursor-not-allowed"
              />
            </div>
          </div>

          <div class="flex justify-end">
            <button
              type="submit"
              :disabled="profileLoading"
              class="bg-indigo-600 text-white px-4 py-2 rounded-md hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span v-if="profileLoading" class="flex items-center">
                <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                保存中...
              </span>
              <span v-else>保存更改</span>
            </button>
          </div>

          <div v-if="profileMessage" :class="profileMessageClass" class="text-sm text-center p-3 rounded-md">
            {{ profileMessage }}
          </div>
        </form>
      </div>

      <!-- 修改密码 -->
      <div class="bg-white shadow rounded-lg mb-6">
        <div class="px-6 py-4 border-b border-gray-200">
          <h3 class="text-lg font-medium text-gray-900">修改密码</h3>
        </div>
        <form @submit.prevent="changePassword" class="p-6 space-y-6">
          <div>
            <label for="currentPassword" class="block text-sm font-medium text-gray-700">
              当前密码
            </label>
            <input
              id="currentPassword"
              v-model="passwordForm.currentPassword"
              type="password"
              required
              class="mt-1 block w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-500"
            />
          </div>

          <div>
            <label for="newPassword" class="block text-sm font-medium text-gray-700">
              新密码
            </label>
            <input
              id="newPassword"
              v-model="passwordForm.newPassword"
              type="password"
              required
              class="mt-1 block w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-500"
            />
            <p class="mt-1 text-xs text-gray-500">密码长度至少6个字符</p>
          </div>

          <div>
            <label for="confirmPassword" class="block text-sm font-medium text-gray-700">
              确认新密码
            </label>
            <input
              id="confirmPassword"
              v-model="passwordForm.confirmPassword"
              type="password"
              required
              class="mt-1 block w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-500"
            />
          </div>

          <div class="flex justify-end">
            <button
              type="submit"
              :disabled="passwordLoading"
              class="bg-red-600 text-white px-4 py-2 rounded-md hover:bg-red-700 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span v-if="passwordLoading" class="flex items-center">
                <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                修改中...
              </span>
              <span v-else>修改密码</span>
            </button>
          </div>

          <div v-if="passwordMessage" :class="passwordMessageClass" class="text-sm text-center p-3 rounded-md">
            {{ passwordMessage }}
          </div>
        </form>
      </div>

      <!-- 游戏偏好设置 -->
      <div class="bg-white shadow rounded-lg mb-6">
        <div class="px-6 py-4 border-b border-gray-200">
          <h3 class="text-lg font-medium text-gray-900">游戏偏好设置</h3>
        </div>
        <div class="p-6 space-y-6">
          <div class="flex items-center justify-between">
            <div>
              <h4 class="text-sm font-medium text-gray-900">音效</h4>
              <p class="text-sm text-gray-500">开启游戏音效</p>
            </div>
            <label class="relative inline-flex items-center cursor-pointer">
              <input
                v-model="preferences.soundEnabled"
                type="checkbox"
                class="sr-only peer"
              />
              <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-indigo-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-indigo-600"></div>
            </label>
          </div>

          <div class="flex items-center justify-between">
            <div>
              <h4 class="text-sm font-medium text-gray-900">自动保存游戏</h4>
              <p class="text-sm text-gray-500">自动保存游戏进度</p>
            </div>
            <label class="relative inline-flex items-center cursor-pointer">
              <input
                v-model="preferences.autoSave"
                type="checkbox"
                class="sr-only peer"
              />
              <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-indigo-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-indigo-600"></div>
            </label>
          </div>

          <div>
            <label for="defaultBoardSize" class="block text-sm font-medium text-gray-700">
              默认棋盘大小
            </label>
            <select
              id="defaultBoardSize"
              v-model="preferences.defaultBoardSize"
              class="mt-1 block w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-500"
            >
              <option value="15">15x15</option>
              <option value="19">19x19</option>
            </select>
          </div>

          <div class="flex justify-end">
            <button
              @click="savePreferences"
              :disabled="preferencesLoading"
              class="bg-indigo-600 text-white px-4 py-2 rounded-md hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span v-if="preferencesLoading" class="flex items-center">
                <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                保存中...
              </span>
              <span v-else>保存偏好设置</span>
            </button>
          </div>

          <div v-if="preferencesMessage" :class="preferencesMessageClass" class="text-sm text-center p-3 rounded-md">
            {{ preferencesMessage }}
          </div>
        </div>
      </div>

      <!-- 账户安全 -->
      <div class="bg-white shadow rounded-lg">
        <div class="px-6 py-4 border-b border-gray-200">
          <h3 class="text-lg font-medium text-gray-900">账户安全</h3>
        </div>
        <div class="p-6 space-y-4">
          <div class="flex items-center justify-between">
            <div>
              <h4 class="text-sm font-medium text-gray-900">账户状态</h4>
              <p class="text-sm text-gray-500">
                账户创建时间：{{ formatDate(userStore.user?.createdAt || '') }}
              </p>
            </div>
            <span class="px-2 py-1 text-xs font-medium bg-green-100 text-green-800 rounded-full">
              正常
            </span>
          </div>

          <div class="border-t border-gray-200 pt-4">
            <button
              @click="logout"
              class="w-full bg-red-600 text-white px-4 py-2 rounded-md hover:bg-red-700"
            >
              退出登录
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../../stores/user'

const router = useRouter()
const userStore = useUserStore()

const profileLoading = ref(false)
const passwordLoading = ref(false)
const preferencesLoading = ref(false)

const profileMessage = ref('')
const passwordMessage = ref('')
const preferencesMessage = ref('')

const profileMessageClass = ref('')
const passwordMessageClass = ref('')
const preferencesMessageClass = ref('')

const profileForm = reactive({
  username: '',
  email: '',
  nickname: ''
})

const passwordForm = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const preferences = reactive({
  soundEnabled: true,
  autoSave: true,
  defaultBoardSize: '15'
})

const updateProfile = async () => {
  profileLoading.value = true
  profileMessage.value = ''

  try {
    await userStore.updateProfile({
      email: profileForm.email,
      nickname: profileForm.nickname
    })
    
    profileMessage.value = '个人信息更新成功'
    profileMessageClass.value = 'bg-green-100 text-green-800'
  } catch (error: any) {
    profileMessage.value = error.message || '更新失败'
    profileMessageClass.value = 'bg-red-100 text-red-800'
  } finally {
    profileLoading.value = false
  }
}

const changePassword = async () => {
  if (passwordForm.newPassword !== passwordForm.confirmPassword) {
    passwordMessage.value = '两次输入的密码不一致'
    passwordMessageClass.value = 'bg-red-100 text-red-800'
    return
  }

  if (passwordForm.newPassword.length < 6) {
    passwordMessage.value = '密码长度至少6个字符'
    passwordMessageClass.value = 'bg-red-100 text-red-800'
    return
  }

  passwordLoading.value = true
  passwordMessage.value = ''

  try {
    await userStore.changePassword({
      currentPassword: passwordForm.currentPassword,
      newPassword: passwordForm.newPassword
    })
    
    passwordMessage.value = '密码修改成功'
    passwordMessageClass.value = 'bg-green-100 text-green-800'
    
    // 清空表单
    passwordForm.currentPassword = ''
    passwordForm.newPassword = ''
    passwordForm.confirmPassword = ''
  } catch (error: any) {
    passwordMessage.value = error.message || '密码修改失败'
    passwordMessageClass.value = 'bg-red-100 text-red-800'
  } finally {
    passwordLoading.value = false
  }
}

const savePreferences = async () => {
  preferencesLoading.value = true
  preferencesMessage.value = ''

  try {
    // 保存到本地存储
    localStorage.setItem('gamePreferences', JSON.stringify(preferences))
    
    preferencesMessage.value = '偏好设置保存成功'
    preferencesMessageClass.value = 'bg-green-100 text-green-800'
  } catch (error) {
    preferencesMessage.value = '保存失败'
    preferencesMessageClass.value = 'bg-red-100 text-red-800'
  } finally {
    preferencesLoading.value = false
  }
}

const logout = async () => {
  if (!confirm('确定要退出登录吗？')) return

  try {
    await userStore.logout()
    router.push('/login')
  } catch (error) {
    console.error('Logout failed:', error)
  }
}

const getRoleName = (role: string): string => {
  const roleNames: Record<string, string> = {
    admin: '管理员',
    user: '普通用户'
  }
  return roleNames[role] || role
}

const formatDate = (dateString: string): string => {
  if (!dateString) return ''
  return new Date(dateString).toLocaleDateString('zh-CN')
}

const loadPreferences = () => {
  try {
    const saved = localStorage.getItem('gamePreferences')
    if (saved) {
      const savedPrefs = JSON.parse(saved)
      Object.assign(preferences, savedPrefs)
    }
  } catch (error) {
    console.error('Failed to load preferences:', error)
  }
}

onMounted(() => {
  if (userStore.user) {
    profileForm.username = userStore.user.username
    profileForm.email = userStore.user.email
    profileForm.nickname = userStore.user.nickname || ''
  }
  
  loadPreferences()
})
</script>