import axios from 'axios'
import type { User, LoginRequest, RegisterRequest, AuthResponse } from '../stores/user'

// 创建 axios 实例
const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器：添加认证令牌
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('auth_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器：处理认证错误
api.interceptors.response.use(
  (response) => {
    return response
  },
  async (error) => {
    const originalRequest = error.config

    // 如果是 401 错误且不是刷新令牌请求，尝试刷新令牌
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true

      try {
        const refreshToken = localStorage.getItem('refresh_token')
        if (refreshToken) {
          const response = await authApi.refreshToken(refreshToken)
          localStorage.setItem('auth_token', response.token)
          localStorage.setItem('refresh_token', response.refreshToken)
          
          // 重新发送原始请求
          originalRequest.headers.Authorization = `Bearer ${response.token}`
          return api(originalRequest)
        }
      } catch (refreshError) {
        // 刷新令牌失败，清除本地存储并跳转到登录页
        localStorage.removeItem('auth_token')
        localStorage.removeItem('refresh_token')
        localStorage.removeItem('user_info')
        window.location.href = '/login'
        return Promise.reject(refreshError)
      }
    }

    return Promise.reject(error)
  }
)

export const authApi = {
  // 登录
  async login(credentials: LoginRequest): Promise<AuthResponse> {
    try {
      const response = await api.post('/auth/login', {
        username: credentials.username,
        password: credentials.password,
        remember_me: false
      })
      
      if (!response.data.success) {
        throw new Error(response.data.message || '登录失败')
      }
      
      return {
        user: response.data.user,
        token: response.data.access_token,
        refreshToken: response.data.refresh_token
      }
    } catch (error: any) {
      throw new Error(error.response?.data?.message || '登录失败')
    }
  },

  // 注册
  async register(userData: RegisterRequest): Promise<AuthResponse> {
    try {
      const response = await api.post('/auth/register', userData)
      
      if (!response.data.success) {
        throw new Error(response.data.message || '注册失败')
      }
      
      // 注册成功后自动登录
      return await this.login({
        username: userData.username,
        password: userData.password
      })
    } catch (error: any) {
      throw new Error(error.response?.data?.message || '注册失败')
    }
  },

  // 登出
  async logout(): Promise<void> {
    try {
      await api.post('/auth/logout')
    } catch (error: any) {
      console.error('Logout error:', error)
      // 即使登出请求失败，也要清除本地存储
    }
  },

  // 刷新令牌
  async refreshToken(refreshToken: string): Promise<AuthResponse> {
    try {
      const response = await api.post('/auth/refresh', { refresh_token: refreshToken })
      
      if (!response.data.success) {
        throw new Error(response.data.message || '令牌刷新失败')
      }
      
      return {
        user: response.data.user || null,
        token: response.data.access_token,
        refreshToken: response.data.refresh_token
      }
    } catch (error: any) {
      throw new Error(error.response?.data?.message || '令牌刷新失败')
    }
  },

  // 获取用户信息
  async getProfile(): Promise<User> {
    try {
      const response = await api.get('/auth/profile')
      
      if (!response.data.success) {
        throw new Error(response.data.message || '获取用户信息失败')
      }
      
      return response.data.user
    } catch (error: any) {
      throw new Error(error.response?.data?.message || '获取用户信息失败')
    }
  },

  // 更新用户信息
  async updateProfile(updates: Partial<User>): Promise<User> {
    try {
      const response = await api.put('/auth/profile', {
        nickname: updates.nickname,
        avatar_url: updates.avatar
      })
      
      if (!response.data.success) {
        throw new Error(response.data.message || '更新用户信息失败')
      }
      
      // 更新后重新获取用户信息
      return await this.getProfile()
    } catch (error: any) {
      throw new Error(error.response?.data?.message || '更新用户信息失败')
    }
  },

  // 修改密码
  async changePassword(currentPassword: string, newPassword: string): Promise<void> {
    try {
      const response = await api.post('/auth/change-password', {
        old_password: currentPassword,
        new_password: newPassword
      })
      
      if (!response.data.success) {
        throw new Error(response.data.message || '修改密码失败')
      }
    } catch (error: any) {
      throw new Error(error.response?.data?.message || '修改密码失败')
    }
  },

  // 验证邮箱
  async verifyEmail(token: string): Promise<void> {
    try {
      await api.post('/auth/verify-email', { token })
    } catch (error: any) {
      throw new Error(error.response?.data?.message || '邮箱验证失败')
    }
  },

  // 重置密码请求
  async requestPasswordReset(email: string): Promise<void> {
    try {
      await api.post('/auth/reset-password', { email })
    } catch (error: any) {
      throw new Error(error.response?.data?.message || '密码重置请求失败')
    }
  },

  // 确认密码重置
  async confirmPasswordReset(token: string, newPassword: string): Promise<void> {
    try {
      await api.post('/auth/reset-password/confirm', {
        token,
        newPassword
      })
    } catch (error: any) {
      throw new Error(error.response?.data?.message || '密码重置失败')
    }
  }
}

export default api