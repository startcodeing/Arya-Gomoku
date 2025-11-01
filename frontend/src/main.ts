import { createApp } from 'vue'
import { createPinia } from 'pinia'
import './style.css'
import App from './App.vue'
import router from './router'
import { useUserStore } from './stores/user'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)

// 初始化用户状态
const userStore = useUserStore()
userStore.initializeAuth().then(() => {
  app.mount('#app')
}).catch((error) => {
  console.error('Failed to initialize auth:', error)
  app.mount('#app')
})