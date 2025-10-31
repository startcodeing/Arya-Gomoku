import { createRouter, createWebHistory } from 'vue-router'
import Home from '../components/Home.vue'
import AIGame from '../components/AIGame.vue'
import RoomList from '../components/RoomList.vue'
import RoomLobby from '../components/RoomLobby.vue'
import PVPGame from '../components/PVPGame.vue'
import GameResult from '../components/GameResult.vue'
import InvitePage from '../components/InvitePage.vue'
import LLMBattle from '../components/LLMBattle.vue'

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
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router