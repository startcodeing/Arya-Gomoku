import { createRouter, createWebHistory } from 'vue-router'
import Home from '../components/Home.vue'
import RoomList from '../components/RoomList.vue'
import RoomLobby from '../components/RoomLobby.vue'
import PVPGame from '../components/PVPGame.vue'
import GameResult from '../components/GameResult.vue'
import InvitePage from '../components/InvitePage.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
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
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router