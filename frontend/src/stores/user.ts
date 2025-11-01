import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '../services/authApi'

export interface User {
  id: string
  username: string
  email: string
  nickname: string
  role: string
  avatar?: string
  createdAt: string
  updatedAt: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface RegisterRequest {
  username: string
  email: string
  nickname: string
  password: string
}

export interface AuthResponse {
  user: User
  token: string
  refreshToken: string
}

export const useUserStore = defineStore('user', () => {
  // 状态
  const user = ref<User | null>(null)
  const token = ref<string | null>(null)
  const refreshToken = ref<string | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  let tokenRefreshTimer: NodeJS.Timeout | null = null

  // 计算属性
  const isAuthenticated = computed(() => !!user.value && !!token.value)
  const userDisplayName = computed(() => user.value?.nickname || user.value?.username || '')

  // 检查令牌是否过期
  const isTokenExpired = (token: string): boolean => {
    try {
      const payload = JSON.parse(atob(token.split('.')[1]))
      const currentTime = Date.now() / 1000
      return payload.exp < currentTime
    } catch (error) {
      return true
    }
  }

  // 设置令牌刷新定时器
  const setupTokenRefreshTimer = (token: string) => {
    if (tokenRefreshTimer) {
      clearTimeout(tokenRefreshTimer)
    }

    try {
      const payload = JSON.parse(atob(token.split('.')[1]))
      const expirationTime = payload.exp * 1000
      const currentTime = Date.now()
      const refreshTime = expirationTime - currentTime - 5 * 60 * 1000 // 提前5分钟刷新

      if (refreshTime > 0) {
        tokenRefreshTimer = setTimeout(async () => {
          try {
            await refreshAuthToken()
          } catch (error) {
            console.error('Auto refresh token failed:', error)
            clearAuth()
          }
        }, refreshTime)
      }
    } catch (error) {
      console.error('Failed to setup token refresh timer:', error)
    }
  }

  // 清除令牌刷新定时器
  const clearTokenRefreshTimer = () => {
    if (tokenRefreshTimer) {
      clearTimeout(tokenRefreshTimer)
      tokenRefreshTimer = null
    }
  }

  // 初始化：从 localStorage 恢复状态
  const initializeAuth = async () => {
    const savedToken = localStorage.getItem('auth_token')
    const savedRefreshToken = localStorage.getItem('refresh_token')
    const savedUser = localStorage.getItem('user_info')

    if (savedToken && savedUser) {
      // 检查令牌是否过期
      if (isTokenExpired(savedToken)) {
        // 尝试刷新令牌
        if (savedRefreshToken) {
          try {
            await refreshAuthToken()
            return
          } catch (error) {
            console.error('Failed to refresh token:', error)
            clearAuth()
            return
          }
        } else {
          clearAuth()
          return
        }
      }

      token.value = savedToken
      refreshToken.value = savedRefreshToken
      try {
        user.value = JSON.parse(savedUser)
        // 设置自动刷新定时器
        setupTokenRefreshTimer(savedToken)
      } catch (error) {
        console.error('Failed to parse saved user info:', error)
        clearAuth()
      }
    }
  }

  // 保存认证信息到 localStorage
  const saveAuth = (authData: AuthResponse) => {
    user.value = authData.user
    token.value = authData.token
    refreshToken.value = authData.refreshToken

    localStorage.setItem('auth_token', authData.token)
    localStorage.setItem('refresh_token', authData.refreshToken)
    localStorage.setItem('user_info', JSON.stringify(authData.user))

    // 设置自动刷新定时器
    setupTokenRefreshTimer(authData.token)
  }

  // 清除认证信息
  const clearAuth = () => {
    user.value = null
    token.value = null
    refreshToken.value = null

    localStorage.removeItem('auth_token')
    localStorage.removeItem('refresh_token')
    localStorage.removeItem('user_info')

    // 清除定时器
    clearTokenRefreshTimer()
  }

  // 登录
  const login = async (credentials: LoginRequest) => {
    loading.value = true
    try {
      const response = await authApi.login(credentials)
      saveAuth(response)
      return response
    } catch (error) {
      clearAuth()
      throw error
    } finally {
      loading.value = false
    }
  }

  // 注册
  const register = async (userData: RegisterRequest) => {
    loading.value = true
    try {
      const response = await authApi.register(userData)
      saveAuth(response)
      return response
    } catch (error) {
      clearAuth()
      throw error
    } finally {
      loading.value = false
    }
  }

  // 登出
  const logout = async () => {
    loading.value = true
    try {
      if (token.value) {
        await authApi.logout()
      }
    } catch (error) {
      console.error('Logout error:', error)
    } finally {
      clearAuth()
      loading.value = false
    }
  }

  // 刷新令牌
  const refreshAuthToken = async () => {
    if (!refreshToken.value) {
      throw new Error('No refresh token available')
    }

    try {
      const response = await authApi.refreshToken(refreshToken.value)
      saveAuth(response)
      return response
    } catch (error) {
      clearAuth()
      throw error
    }
  }

  // 获取用户信息
  const fetchUserProfile = async () => {
    if (!token.value) {
      throw new Error('No authentication token')
    }

    loading.value = true
    try {
      const userProfile = await authApi.getProfile()
      user.value = userProfile
      localStorage.setItem('user_info', JSON.stringify(userProfile))
      return userProfile
    } catch (error) {
      throw error
    } finally {
      loading.value = false
    }
  }

  // 更新用户信息
  const updateProfile = async (updates: Partial<User>) => {
    if (!token.value) {
      throw new Error('No authentication token')
    }

    loading.value = true
    try {
      const updatedUser = await authApi.updateProfile(updates)
      user.value = updatedUser
      localStorage.setItem('user_info', JSON.stringify(updatedUser))
      return updatedUser
    } catch (error) {
      throw error
    } finally {
      loading.value = false
    }
  }

  // 修改密码
  const changePassword = async (passwordData: { currentPassword: string; newPassword: string }) => {
    if (!token.value) {
      throw new Error('No authentication token')
    }

    loading.value = true
    try {
      await authApi.changePassword(passwordData.currentPassword, passwordData.newPassword)
    } catch (error) {
      throw error
    } finally {
      loading.value = false
    }
  }

  return {
    // 状态
    user,
    token,
    refreshToken,
    loading,
    
    // 计算属性
    isAuthenticated,
    userDisplayName,
    
    // 方法
    initializeAuth,
    login,
    register,
    logout,
    refreshAuthToken,
    fetchUserProfile,
    updateProfile,
    changePassword,
    clearAuth
  }
})