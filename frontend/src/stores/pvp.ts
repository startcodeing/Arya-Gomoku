import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  type Player,
  type Room,
  type CreateRoomRequest,
  type CreateRoomResponse,
  type JoinRoomRequest,
  type JoinRoomResponse,
  type RoomListResponse,
  type ChatMessage,
  type ConnectionStatus,
  type PVPGame
} from '../types/pvp'
import { pvpApi } from '../services/pvpApi'
import { getGlobalWebSocketService } from '../services/websocket'

export const usePvpStore = defineStore('pvp', () => {
  const rooms = ref<Room[]>([])
  const currentRoom = ref<Room | null>(null)
  const currentPlayer = ref<Player | null>(null)
  const currentGame = ref<PVPGame | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)
  const chatMessages = ref<ChatMessage[]>([])
  const connectionStatus = ref<ConnectionStatus>('disconnected')

  async function fetchRooms() {
    isLoading.value = true
    error.value = null
    try {
      console.log('开始获取房间列表...')
      const resp: RoomListResponse = await pvpApi.getActiveRooms()
      // 处理后端返回null的情况
      rooms.value = resp.rooms || []
      console.log('房间列表获取成功:', rooms.value)
    } catch (e: any) {
      console.error('获取房间列表失败:', e)
      error.value = e?.message ?? '获取房间列表失败'
    } finally {
      isLoading.value = false
      console.log('房间列表加载完成，isLoading:', isLoading.value)
    }
  }

  async function createRoom(payload: CreateRoomRequest) {
    isLoading.value = true
    error.value = null
    try {
      const resp: CreateRoomResponse = await pvpApi.createRoom(payload)
      const room = resp.room
      currentRoom.value = room
      // Backend only returns room on create; pick the creator as first player
      const player = resp.player ?? room.players[0]
      currentPlayer.value = player
      const ws = getGlobalWebSocketService()
      
      // 设置WebSocket事件监听器
      setupWebSocketEventHandlers(ws)
      
      // 不阻塞创建流程：WebSocket连接失败也不影响导航
      ws.connect(room.id, player.id).catch((err) => {
        console.warn('WebSocket连接失败，将在房间内重试', err)
      })
    } catch (e: any) {
      error.value = e?.message ?? '创建房间失败'
      throw e
    } finally {
      isLoading.value = false
    }
  }

  async function joinRoom(roomId: string, payload: JoinRoomRequest) {
    isLoading.value = true
    error.value = null
    try {
      const resp: JoinRoomResponse = await pvpApi.joinRoom(roomId, payload)
      currentRoom.value = resp.room
      currentPlayer.value = resp.player
      const ws = getGlobalWebSocketService()
      
      // 设置WebSocket事件监听器
      setupWebSocketEventHandlers(ws)
      
      // 不阻塞加入流程：WebSocket连接失败也不影响导航
      ws.connect(resp.room.id, resp.player.id).catch((err) => {
        console.warn('WebSocket连接失败，将在房间内重试', err)
      })
    } catch (e: any) {
      error.value = e?.message ?? '加入房间失败'
      throw e
    } finally {
      isLoading.value = false
    }
  }

  // 计算属性
  const isConnected = computed(() => connectionStatus.value === 'connected')
  const canStartGame = computed(() => {
    if (!currentRoom.value || currentRoom.value.status !== 'waiting') return false
    if (currentRoom.value.players.length < 2) return false
    return currentRoom.value.players.every(p => p.isReady)
  })
  const isMyTurn = computed(() => {
    if (!currentGame.value || !currentPlayer.value) return false
    return currentGame.value.currentPlayer === currentPlayer.value.id
  })

  // 获取房间信息
  async function getRoom(roomId: string) {
    isLoading.value = true
    error.value = null
    try {
      const resp = await pvpApi.getRoom(roomId)
      currentRoom.value = resp.room
      return resp.room
    } catch (e: any) {
      const errorMessage = e.response?.data?.message || e.message || '获取房间信息失败'
      error.value = errorMessage
      console.error('获取房间信息失败:', e)
      throw new Error(errorMessage)
    } finally {
      isLoading.value = false
    }
  }

  // 离开房间
  async function leaveRoom() {
    if (!currentRoom.value) return
    
    isLoading.value = true
    error.value = null
    try {
      // 断开WebSocket连接
      const ws = getGlobalWebSocketService()
      ws.disconnect()
      
      // 清理状态
      currentRoom.value = null
      currentPlayer.value = null
      currentGame.value = null
      chatMessages.value = []
      connectionStatus.value = 'disconnected'
    } catch (e: any) {
      error.value = e?.message ?? '离开房间失败'
      throw e
    } finally {
      isLoading.value = false
    }
  }

  // 切换准备状态
  function toggleReady() {
    if (!currentPlayer.value) return
    
    const newReadyState = !currentPlayer.value.isReady
    console.log('发送准备状态变更:', { ready: newReadyState })
    
    // 通过WebSocket发送准备状态变更
    const ws = getGlobalWebSocketService()
    const success = ws.send({
      type: 'ready',
      data: {
        ready: newReadyState
      }
    })
    
    console.log('WebSocket消息发送结果:', success)
  }

  // 开始游戏
  async function startGame() {
    if (!currentRoom.value) return
    
    isLoading.value = true
    error.value = null
    try {
      await pvpApi.startGame(currentRoom.value.id)
    } catch (e: any) {
      error.value = e?.message ?? '开始游戏失败'
      throw e
    } finally {
      isLoading.value = false
    }
  }

  // 发送聊天消息
  async function sendChatMessage(message: string) {
    if (!currentRoom.value || !currentPlayer.value) {
      throw new Error('未加入房间或玩家信息不存在')
    }

    const ws = getGlobalWebSocketService()
    if (!ws.isConnected()) {
      throw new Error('WebSocket未连接')
    }

    ws.send({
      type: 'chat_message',
      data: {
        roomId: currentRoom.value.id,
        playerId: currentPlayer.value.id,
        message: message
      }
    })
  }

  // 进行移动
  async function makeMove(x: number, y: number) {
    if (!currentRoom.value || !currentPlayer.value) {
      throw new Error('未加入房间或玩家信息不存在')
    }

    if (!currentGame.value) {
      throw new Error('游戏未开始')
    }

    if (currentGame.value.currentPlayer !== currentPlayer.value.id) {
      throw new Error('不是你的回合')
    }

    try {
      isLoading.value = true
      error.value = null

      const response = await pvpApi.makeMove(currentRoom.value.id, {
        x,
        y,
        playerId: currentPlayer.value.id
      })

      // 更新游戏状态
      if (response.room) {
        currentRoom.value = response.room
        currentGame.value = response.room.game
      }

    } catch (err: any) {
      error.value = err.message || '移动失败'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  // 清除错误
  function clearError() {
    error.value = null
  }

  // WebSocket事件处理
  function handleWebSocketMessage(message: any) {
    console.log('收到WebSocket消息:', message.type, message.data)
    switch (message.type) {
      case 'room_updated':
        console.log('处理房间更新消息:', message.data)
        if (message.data.room) {
          currentRoom.value = message.data.room
          console.log('房间状态已更新:', currentRoom.value)
        }
        break
      case 'game_start':
        if (message.data.game) {
          currentGame.value = message.data.game
          if (currentRoom.value) {
            currentRoom.value.status = 'playing'
          }
        }
        break
      case 'game_update':
        if (message.data.game) {
          currentGame.value = message.data.game
        }
        break
      case 'chat_message':
        if (message.data) {
          chatMessages.value.push({
            id: Date.now().toString(),
            playerId: message.data.playerId,
            playerName: message.data.playerName,
            message: message.data.message,
            timestamp: message.data.timestamp
          })
        }
        break
      case 'player_joined':
      case 'player_left':
        if (message.data.room) {
          currentRoom.value = message.data.room
        }
        break
      case 'error':
        error.value = message.data.message || '发生未知错误'
        break
    }
  }

  // 设置WebSocket连接状态
  function setConnectionStatus(status: ConnectionStatus) {
    connectionStatus.value = status
  }

  // 设置WebSocket事件监听器
  function setupWebSocketEventHandlers(ws: any) {
    ws.onOpen = () => {
      setConnectionStatus('connected')
      console.log('WebSocket连接已建立')
    }

    ws.onClose = () => {
      setConnectionStatus('disconnected')
      console.log('WebSocket连接已断开')
    }

    ws.onError = (error: Event) => {
      setConnectionStatus('error')
      console.error('WebSocket连接错误:', error)
    }

    ws.onMessage = (message: any) => {
      handleWebSocketMessage(message)
    }

    ws.onReconnecting = () => {
      setConnectionStatus('reconnecting')
      console.log('WebSocket重连中...')
    }

    ws.onReconnected = () => {
      setConnectionStatus('connected')
      console.log('WebSocket重连成功')
    }
  }

  return {
    // 状态
    rooms,
    currentRoom,
    currentPlayer,
    currentGame,
    isLoading,
    error,
    chatMessages,
    connectionStatus,
    
    // 计算属性
    isConnected,
    canStartGame,
    isMyTurn,
    
    // 方法
    fetchRooms,
    createRoom,
    joinRoom,
    getRoom,
    leaveRoom,
    toggleReady,
    startGame,
    makeMove,
    sendChatMessage,
    clearError,
    handleWebSocketMessage,
    setConnectionStatus,
    setupWebSocketEventHandlers,
  }
})