<template>
  <div class="register-container">
    <!-- 背景装饰 -->
    <div class="background-decoration">
      <div class="gradient-overlay"></div>
      <div class="floating-shape shape-1"></div>
      <div class="floating-shape shape-2"></div>
    </div>
    
    <!-- 浮动装饰元素 -->
    <div class="floating-elements">
      <div class="floating-dot dot-1"></div>
      <div class="floating-dot dot-2"></div>
      <div class="floating-dot dot-3"></div>
    </div>
    
    <div class="register-wrapper">
      <!-- 注册卡片 -->
      <div class="register-card">
        <!-- 卡片头部 -->
        <div class="card-header">
          <div class="icon-container">
            <div class="icon-wrapper">
              <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z"></path>
              </svg>
            </div>
          </div>
          <h2 class="title">创建账户</h2>
          <p class="subtitle">加入五子棋世界，开始你的对弈之旅</p>
        </div>

        <!-- 表单内容 -->
        <div class="form-container">
          <form @submit.prevent="handleRegister" class="register-form">
            <!-- 用户名输入 -->
            <div class="form-group">
              <label for="username" class="form-label">
                用户名 <span class="required">*</span>
              </label>
              <div class="input-wrapper">
                <div class="input-icon">
                  <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
                  </svg>
                </div>
                <input
                  id="username"
                  v-model="form.username"
                  @blur="validateUsername"
                  @input="clearError('username')"
                  type="text"
                  :class="['form-input', errors.username ? 'error' : '']"
                  placeholder="请输入用户名"
                  required
                />
              </div>
              <p v-if="errors.username" class="error-message">{{ errors.username }}</p>
              <p class="help-text">3-20个字符，只能包含字母、数字和下划线</p>
            </div>

            <!-- 邮箱输入 -->
            <div class="form-group">
              <label for="email" class="form-label">
                邮箱地址 <span class="required">*</span>
              </label>
              <div class="input-wrapper">
                <div class="input-icon">
                  <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 12a4 4 0 10-8 0 4 4 0 008 0zm0 0v1.5a2.5 2.5 0 005 0V12a9 9 0 10-9 9m4.5-1.206a8.959 8.959 0 01-4.5 1.207"></path>
                  </svg>
                </div>
                <input
                  id="email"
                  v-model="form.email"
                  @blur="validateEmail"
                  @input="clearError('email')"
                  type="email"
                  :class="['form-input', errors.email ? 'error' : '']"
                  placeholder="请输入邮箱地址"
                  required
                />
              </div>
              <p v-if="errors.email" class="error-message">{{ errors.email }}</p>
            </div>

            <!-- 昵称输入 -->
            <div class="form-group">
              <label for="nickname" class="form-label">
                昵称 <span class="optional">(可选)</span>
              </label>
              <div class="input-wrapper">
                <div class="input-icon">
                  <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z"></path>
                  </svg>
                </div>
                <input
                  id="nickname"
                  v-model="form.nickname"
                  @blur="validateNickname"
                  @input="clearError('nickname')"
                  type="text"
                  :class="['form-input', errors.nickname ? 'error' : '']"
                  placeholder="请输入昵称（不填写将使用用户名）"
                />
              </div>
              <p v-if="errors.nickname" class="error-message">{{ errors.nickname }}</p>
              <p class="help-text">1-20个字符，用于显示</p>
            </div>

            <!-- 密码输入 -->
            <div class="form-group">
              <label for="password" class="form-label">
                密码 <span class="required">*</span>
              </label>
              <div class="input-wrapper">
                <div class="input-icon">
                  <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"></path>
                  </svg>
                </div>
                <input
                  id="password"
                  v-model="form.password"
                  @blur="validatePassword"
                  @input="clearError('password')"
                  :type="showPassword ? 'text' : 'password'"
                  :class="['form-input password-input', errors.password ? 'error' : '']"
                  placeholder="请输入密码"
                  required
                />
                <button
                  type="button"
                  @click="showPassword = !showPassword"
                  class="password-toggle"
                >
                  <svg v-if="showPassword" class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L3 3m6.878 6.878L21 21"></path>
                  </svg>
                  <svg v-else class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
                  </svg>
                </button>
              </div>
              <p v-if="errors.password" class="error-message">{{ errors.password }}</p>
              
              <!-- 密码强度指示器 -->
              <div v-if="form.password" class="password-strength">
                <div class="strength-header">
                  <span class="strength-label">密码强度</span>
                  <span :class="['strength-text', passwordStrength.textColor]">{{ passwordStrength.text }}</span>
                </div>
                <div class="strength-bar">
                  <div 
                    :class="['strength-fill', passwordStrength.color]" 
                    :style="{ width: passwordStrength.width }"
                  ></div>
                </div>
                <div class="strength-checks">
                  <div class="check-item">
                    <svg :class="['check-icon', passwordChecks.length ? 'valid' : 'invalid']" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"></path>
                    </svg>
                    <span :class="passwordChecks.length ? 'valid' : 'invalid'">至少8个字符</span>
                  </div>
                  <div class="check-item">
                    <svg :class="['check-icon', passwordChecks.hasLetter ? 'valid' : 'invalid']" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"></path>
                    </svg>
                    <span :class="passwordChecks.hasLetter ? 'valid' : 'invalid'">包含字母</span>
                  </div>
                  <div class="check-item">
                    <svg :class="['check-icon', passwordChecks.hasNumber ? 'valid' : 'invalid']" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"></path>
                    </svg>
                    <span :class="passwordChecks.hasNumber ? 'valid' : 'invalid'">包含数字</span>
                  </div>
                  <div class="check-item">
                    <svg :class="['check-icon', passwordChecks.hasSpecial ? 'valid' : 'invalid']" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"></path>
                    </svg>
                    <span :class="passwordChecks.hasSpecial ? 'valid' : 'invalid'">包含特殊字符</span>
                  </div>
                </div>
              </div>
            </div>

            <!-- 错误提示 -->
            <div v-if="error" class="alert alert-error">
              <div class="alert-icon">
                <svg class="icon" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"></path>
                </svg>
              </div>
              <div class="alert-content">
                <p class="alert-message">{{ error }}</p>
              </div>
            </div>

            <!-- 成功提示 -->
            <div v-if="success" class="alert alert-success">
              <div class="alert-icon">
                <svg class="icon" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"></path>
                </svg>
              </div>
              <div class="alert-content">
                <p class="alert-message">{{ success }}</p>
              </div>
            </div>

            <!-- 注册按钮 -->
            <button
              type="submit"
              :disabled="!isFormValid || isLoading"
              class="register-button"
            >
              <span v-if="isLoading" class="loading-spinner"></span>
              <span class="button-content">
                <svg v-if="!isLoading" class="button-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z"></path>
                </svg>
                {{ isLoading ? '注册中...' : '创建账户' }}
              </span>
            </button>
          </form>

          <!-- 底部链接 -->
          <div class="login-link">
            <p>
              已有账户？
              <router-link to="/login" class="link">立即登录</router-link>
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

const form = ref({
  username: '',
  email: '',
  nickname: '',
  password: ''
})

const errors = ref<Record<string, string>>({})
const isLoading = ref(false)
const error = ref('')
const success = ref('')
const showPassword = ref(false)

// 密码强度检查
const passwordChecks = computed(() => ({
  length: form.value.password.length >= 8,
  hasLetter: /[a-zA-Z]/.test(form.value.password),
  hasNumber: /\d/.test(form.value.password),
  hasSpecial: /[!@#$%^&*(),.?":{}|<>]/.test(form.value.password)
}))

// 密码强度计算
const passwordStrength = computed(() => {
  const checks = Object.values(passwordChecks.value)
  const score = checks.filter(Boolean).length
  
  if (score === 0) return { width: '0%', color: 'weak', text: '', textColor: 'weak' }
  if (score === 1) return { width: '25%', color: 'weak', text: '弱', textColor: 'weak' }
  if (score === 2) return { width: '50%', color: 'fair', text: '一般', textColor: 'fair' }
  if (score === 3) return { width: '75%', color: 'good', text: '良好', textColor: 'good' }
  return { width: '100%', color: 'strong', text: '强', textColor: 'strong' }
})

// 表单验证
const isFormValid = computed(() => {
  return form.value.username && 
         form.value.email && 
         form.value.password && 
         !Object.values(errors.value).some(error => error)
})

// 验证函数
const validateUsername = () => {
  const username = form.value.username.trim()
  if (!username) {
    errors.value.username = '用户名不能为空'
  } else if (username.length < 3 || username.length > 20) {
    errors.value.username = '用户名长度必须在3-20个字符之间'
  } else if (!/^[a-zA-Z0-9_]+$/.test(username)) {
    errors.value.username = '用户名只能包含字母、数字和下划线'
  } else {
    delete errors.value.username
  }
}

const validateEmail = () => {
  const email = form.value.email.trim()
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!email) {
    errors.value.email = '邮箱不能为空'
  } else if (!emailRegex.test(email)) {
    errors.value.email = '请输入有效的邮箱地址'
  } else {
    delete errors.value.email
  }
}

const validateNickname = () => {
  const nickname = form.value.nickname.trim()
  if (nickname && (nickname.length < 1 || nickname.length > 20)) {
    errors.value.nickname = '昵称长度必须在1-20个字符之间'
  } else {
    delete errors.value.nickname
  }
}

const validatePassword = () => {
  const password = form.value.password
  if (!password) {
    errors.value.password = '密码不能为空'
  } else if (password.length < 8 || password.length > 50) {
    errors.value.password = '密码长度必须在8-50个字符之间'
  } else {
    delete errors.value.password
  }
}

const clearError = (field: string) => {
  if (errors.value[field]) {
    delete errors.value[field]
  }
  error.value = ''
}

// 注册处理
const handleRegister = async () => {
  // 验证所有字段
  validateUsername()
  validateEmail()
  validateNickname()
  validatePassword()
  
  if (Object.keys(errors.value).length > 0) {
    return
  }
  
  isLoading.value = true
  error.value = ''
  success.value = ''
  
  try {
    await userStore.register({
      username: form.value.username.trim(),
      email: form.value.email.trim(),
      nickname: form.value.nickname.trim() || form.value.username.trim(),
      password: form.value.password
    })
    
    success.value = '注册成功！正在跳转到登录页面...'
    
    setTimeout(() => {
      router.push('/login')
    }, 2000)
  } catch (err: any) {
    error.value = err.message || '注册失败，请重试'
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
/* 主容器 */
.register-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #1e293b 0%, #7c3aed 50%, #1e293b 100%);
  position: relative;
  overflow: hidden;
  padding: 1rem;
}

/* 背景装饰 */
.background-decoration {
  position: absolute;
  inset: 0;
  opacity: 0.2;
}

.gradient-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, rgba(16, 185, 129, 0.1) 0%, transparent 50%, rgba(59, 130, 246, 0.1) 100%);
}

.floating-shape {
  position: absolute;
  border-radius: 50%;
  filter: blur(60px);
}

.shape-1 {
  top: 25%;
  left: 25%;
  width: 24rem;
  height: 24rem;
  background: rgba(16, 185, 129, 0.05);
}

.shape-2 {
  bottom: 25%;
  right: 25%;
  width: 24rem;
  height: 24rem;
  background: rgba(59, 130, 246, 0.05);
}

/* 浮动装饰元素 */
.floating-elements {
  position: absolute;
  inset: 0;
  pointer-events: none;
}

.floating-dot {
  position: absolute;
  border-radius: 50%;
  filter: blur(20px);
  animation: pulse 3s ease-in-out infinite;
}

.dot-1 {
  top: 2.5rem;
  right: 2.5rem;
  width: 5rem;
  height: 5rem;
  background: rgba(16, 185, 129, 0.2);
}

.dot-2 {
  top: 8rem;
  left: 5rem;
  width: 4rem;
  height: 4rem;
  background: rgba(34, 197, 94, 0.2);
  animation-delay: 1s;
}

.dot-3 {
  bottom: 5rem;
  right: 8rem;
  width: 6rem;
  height: 6rem;
  background: rgba(20, 184, 166, 0.2);
  animation-delay: 2s;
}

@keyframes pulse {
  0%, 100% { opacity: 0.4; transform: scale(1); }
  50% { opacity: 0.8; transform: scale(1.1); }
}

/* 注册包装器 */
.register-wrapper {
  position: relative;
  z-index: 10;
  max-width: 28rem;
  width: 100%;
}

/* 注册卡片 */
.register-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
  border-radius: 1.5rem;
  border: 1px solid rgba(255, 255, 255, 0.2);
  overflow: hidden;
}

/* 卡片头部 */
.card-header {
  padding: 2rem 2rem 1.5rem;
  background: linear-gradient(135deg, rgba(16, 185, 129, 0.1) 0%, rgba(59, 130, 246, 0.1) 100%);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.icon-container {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 1rem;
}

.icon-wrapper {
  width: 3rem;
  height: 3rem;
  background: linear-gradient(135deg, #10b981 0%, #3b82f6 100%);
  border-radius: 1rem;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 10px 25px -5px rgba(16, 185, 129, 0.4);
}

.icon-wrapper .icon {
  width: 1.5rem;
  height: 1.5rem;
  color: white;
}

.title {
  font-size: 1.875rem;
  font-weight: 700;
  text-align: center;
  background: linear-gradient(135deg, #059669 0%, #2563eb 100%);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  margin-bottom: 0.5rem;
}

.subtitle {
  text-align: center;
  color: #6b7280;
  font-size: 0.875rem;
}

/* 表单容器 */
.form-container {
  padding: 2rem;
}

.register-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

/* 表单组 */
.form-group {
  display: flex;
  flex-direction: column;
}

.form-label {
  display: block;
  font-size: 0.875rem;
  font-weight: 500;
  color: #374151;
  margin-bottom: 0.5rem;
}

.required {
  color: #ef4444;
}

.optional {
  color: #9ca3af;
}

/* 输入框包装器 */
.input-wrapper {
  position: relative;
}

.input-icon {
  position: absolute;
  top: 50%;
  left: 0.75rem;
  transform: translateY(-50%);
  display: flex;
  align-items: center;
  pointer-events: none;
}

.input-icon .icon {
  width: 1.25rem;
  height: 1.25rem;
  color: #9ca3af;
}

/* 表单输入框 */
.form-input {
  width: 100%;
  padding: 0.75rem 1rem 0.75rem 2.5rem;
  color: #111827;
  background: white;
  border: 1px solid #d1d5db;
  border-radius: 0.75rem;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
  font-size: 1rem;
  transition: all 0.2s ease-in-out;
}

.form-input::placeholder {
  color: #9ca3af;
}

.form-input:focus {
  outline: none;
  border-color: #10b981;
  box-shadow: 0 0 0 3px rgba(16, 185, 129, 0.1);
}

.form-input.error {
  border-color: #ef4444;
}

.form-input.error:focus {
  border-color: #ef4444;
  box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.1);
}

.password-input {
  padding-right: 3rem;
}

/* 密码切换按钮 */
.password-toggle {
  position: absolute;
  top: 50%;
  right: 0.75rem;
  transform: translateY(-50%);
  display: flex;
  align-items: center;
  color: #9ca3af;
  background: none;
  border: none;
  cursor: pointer;
  transition: color 0.2s ease-in-out;
}

.password-toggle:hover {
  color: #6b7280;
}

.password-toggle .icon {
  width: 1.25rem;
  height: 1.25rem;
}

/* 错误消息 */
.error-message {
  margin-top: 0.25rem;
  font-size: 0.875rem;
  color: #dc2626;
}

/* 帮助文本 */
.help-text {
  margin-top: 0.25rem;
  font-size: 0.75rem;
  color: #6b7280;
}

/* 密码强度指示器 */
.password-strength {
  margin-top: 0.75rem;
}

.strength-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 0.5rem;
}

.strength-label {
  font-size: 0.75rem;
  color: #6b7280;
}

.strength-text {
  font-size: 0.75rem;
  font-weight: 500;
}

.strength-text.weak { color: #ef4444; }
.strength-text.fair { color: #f59e0b; }
.strength-text.good { color: #3b82f6; }
.strength-text.strong { color: #10b981; }

.strength-bar {
  width: 100%;
  height: 0.5rem;
  background: #e5e7eb;
  border-radius: 9999px;
  overflow: hidden;
}

.strength-fill {
  height: 100%;
  transition: all 0.3s ease-in-out;
  border-radius: 9999px;
}

.strength-fill.weak { background: #ef4444; }
.strength-fill.fair { background: #f59e0b; }
.strength-fill.good { background: #3b82f6; }
.strength-fill.strong { background: #10b981; }

.strength-checks {
  margin-top: 0.5rem;
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.check-item {
  display: flex;
  align-items: center;
  font-size: 0.75rem;
}

.check-icon {
  width: 0.75rem;
  height: 0.75rem;
  margin-right: 0.25rem;
}

.check-icon.valid { color: #10b981; }
.check-icon.invalid { color: #9ca3af; }

.check-item span.valid { color: #059669; }
.check-item span.invalid { color: #6b7280; }

/* 警告框 */
.alert {
  padding: 1rem;
  border-radius: 0.5rem;
  display: flex;
  align-items: flex-start;
}

.alert-error {
  background: #fef2f2;
  border-left: 4px solid #ef4444;
}

.alert-success {
  background: #f0fdf4;
  border-left: 4px solid #22c55e;
}

.alert-icon {
  flex-shrink: 0;
  margin-right: 0.75rem;
}

.alert-icon .icon {
  width: 1.25rem;
  height: 1.25rem;
}

.alert-error .alert-icon .icon {
  color: #ef4444;
}

.alert-success .alert-icon .icon {
  color: #22c55e;
}

.alert-content {
  flex: 1;
}

.alert-message {
  font-size: 0.875rem;
  margin: 0;
}

.alert-error .alert-message {
  color: #b91c1c;
}

.alert-success .alert-message {
  color: #15803d;
}

/* 注册按钮 */
.register-button {
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 0.75rem 1.5rem;
  background: linear-gradient(135deg, #10b981 0%, #3b82f6 100%);
  color: white;
  font-weight: 500;
  border-radius: 0.75rem;
  box-shadow: 0 10px 25px -5px rgba(16, 185, 129, 0.4);
  border: none;
  cursor: pointer;
  transition: all 0.2s ease-in-out;
  font-size: 1rem;
}

.register-button:hover:not(:disabled) {
  box-shadow: 0 20px 40px -10px rgba(16, 185, 129, 0.4);
  background: linear-gradient(135deg, #059669 0%, #2563eb 100%);
  transform: translateY(-1px);
}

.register-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}

.button-content {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.button-icon {
  width: 1.25rem;
  height: 1.25rem;
}

/* 加载动画 */
.loading-spinner {
  width: 1.25rem;
  height: 1.25rem;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top: 2px solid white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-right: 0.75rem;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* 登录链接 */
.login-link {
  margin-top: 1.5rem;
  text-align: center;
}

.login-link p {
  font-size: 0.875rem;
  color: #6b7280;
  margin: 0;
}

.link {
  font-weight: 500;
  color: #10b981;
  text-decoration: none;
  transition: color 0.2s ease-in-out;
}

.link:hover {
  color: #059669;
}

/* 响应式设计 */
@media (max-width: 640px) {
  .register-container {
    padding: 0.5rem;
  }
  
  .register-wrapper {
    max-width: 100%;
  }
  
  .card-header {
    padding: 1.5rem 1.5rem 1rem;
  }
  
  .form-container {
    padding: 1.5rem;
  }
  
  .title {
    font-size: 1.5rem;
  }
  
  .floating-dot {
    display: none;
  }
}
</style>