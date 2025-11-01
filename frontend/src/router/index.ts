import { createRouter, createWebHistory } from 'vue-router'
import Home from '../components/Home.vue'
import AIGame from '../components/AIGame.vue'
import RoomList from '../components/RoomList.vue'
import RoomLobby from '../components/RoomLobby.vue'
import PVPGame from '../components/PVPGame.vue'
import GameResult from '../components/GameResult.vue'
import InvitePage from '../components/InvitePage.vue'
import LLMBattle from '../components/LLMBattle.vue'
import Login from '../components/auth/Login.vue'
import Register from '../components/auth/Register.vue'
import UserProfile from '../components/user/UserProfile.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/ai-game',
    name: 'AIGame',
    component: AIGame
  },
  {
    path: '/pvp',
    name: 'PVP',
    component: RoomList
  },
  {
    path: '/room/:id',
    name: 'RoomLobby',
    component: RoomLobby,
    props: true
  },
  {
    path: '/game/:id',
    name: 'PVPGame',
    component: PVPGame,
    props: true
  },
  {
    path: '/result/:id',
    name: 'GameResult',
    component: GameResult,
    props: true
  },
  {
    path: '/invite/:id',
    name: 'InvitePage',
    component: InvitePage,
    props: true
  },
  {
    path: '/llm-battle',
    name: 'LLMBattle',
    component: LLMBattle
  },
  {
    path: '/login',
    name: 'Login',
    component: Login
  },
  {
    path: '/register',
    name: 'Register',
    component: Register
  },
  {
    path: '/profile',
    name: 'UserProfile',
    component: UserProfile
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach(async (to, from, next) => {
  try {
    // 动态导入用户store以避免循环依赖
    const { useUserStore } = await import('../stores/user')
    const userStore = useUserStore()
    
    // 公开路由，不需要认证
    const publicRoutes = ['Home', 'Login', 'Register']
    
    // 如果是公开路由，直接通过
    if (publicRoutes.includes(to.name as string)) {
      // 如果已登录用户访问登录/注册页，重定向到首页
      if ((to.name === 'Login' || to.name === 'Register') && userStore.isAuthenticated) {
        next({ name: 'Home' })
        return
      }
      next()
      return
    }

    // 对于受保护的路由，检查认证状态
    if (!userStore.isAuthenticated) {
      // 如果有token但用户信息为空，尝试初始化
      const token = localStorage.getItem('auth_token')
      if (token) {
        try {
          await userStore.initializeAuth()
          // 初始化成功，检查是否已认证
          if (userStore.isAuthenticated) {
            next()
            return
          }
        } catch (error) {
          console.error('认证初始化失败:', error)
          // 清除无效的认证信息
          localStorage.removeItem('auth_token')
          localStorage.removeItem('refresh_token')
          localStorage.removeItem('user_info')
        }
      }
      
      // 未登录，跳转到登录页
      next({ 
        name: 'Login', 
        query: { redirect: to.fullPath } 
      })
      return
    }

    // 已认证，允许访问
    next()
  } catch (error) {
    console.error('路由守卫错误:', error)
    // 发生错误时，重定向到首页
    next({ name: 'Home' })
  }
})

export default router