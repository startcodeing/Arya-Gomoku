import type { WebSocketMessage, ConnectionStatus } from '../types/pvp'

export class WebSocketService {
  private ws: WebSocket | null = null
  private reconnectAttempts = 0
  private maxReconnectAttempts = 5
  private reconnectDelay = 1000
  private heartbeatInterval: number | null = null
  private connectionTimeout: number | null = null
  private isManualClose = false

  // Event handlers
  public onOpen: (() => void) | null = null
  public onClose: (() => void) | null = null
  public onError: ((error: Event) => void) | null = null
  public onMessage: ((message: WebSocketMessage) => void) | null = null
  public onReconnecting: (() => void) | null = null
  public onReconnected: (() => void) | null = null

  // WebSocket URL - adjust according to your backend configuration
  private getWebSocketUrl(): string {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const host = import.meta.env.VITE_WS_HOST || window.location.host
    return `${protocol}//${host}/api/ws`
  }

  // Connect to WebSocket server
  async connect(roomId: string, playerId: string): Promise<void> {
    return new Promise((resolve, reject) => {
      try {
        this.isManualClose = false
        const url = `${this.getWebSocketUrl()}?roomId=${roomId}&playerId=${playerId}`
        
        this.ws = new WebSocket(url)
        
        // Set connection timeout
        this.connectionTimeout = window.setTimeout(() => {
          if (this.ws && this.ws.readyState === WebSocket.CONNECTING) {
            this.ws.close()
            reject(new Error('WebSocket连接超时'))
          }
        }, 10000)

        this.ws.onopen = (event) => {
          console.log('WebSocket连接已建立')
          
          if (this.connectionTimeout) {
            clearTimeout(this.connectionTimeout)
            this.connectionTimeout = null
          }
          
          this.reconnectAttempts = 0
          this.startHeartbeat()
          
          if (this.onOpen) {
            this.onOpen()
          }
          
          if (this.reconnectAttempts > 0 && this.onReconnected) {
            this.onReconnected()
          }
          
          resolve()
        }

        this.ws.onclose = (event) => {
          console.log('WebSocket连接已关闭', event.code, event.reason)
          
          this.stopHeartbeat()
          
          if (this.connectionTimeout) {
            clearTimeout(this.connectionTimeout)
            this.connectionTimeout = null
          }
          
          if (this.onClose) {
            this.onClose()
          }
          
          // Auto-reconnect if not manually closed
          if (!this.isManualClose && this.reconnectAttempts < this.maxReconnectAttempts) {
            this.attemptReconnect(roomId, playerId)
          }
        }

        this.ws.onerror = (event) => {
          console.error('WebSocket错误:', event)
          
          if (this.onError) {
            this.onError(event)
          }
          
          if (this.ws?.readyState === WebSocket.CONNECTING) {
            reject(new Error('WebSocket连接失败'))
          }
        }

        this.ws.onmessage = (event) => {
          try {
            const message: WebSocketMessage = JSON.parse(event.data)
            
            if (this.onMessage) {
              this.onMessage(message)
            }
          } catch (error) {
            console.error('解析WebSocket消息失败:', error, event.data)
          }
        }

      } catch (error) {
        reject(error)
      }
    })
  }

  // Disconnect from WebSocket server
  disconnect(): void {
    this.isManualClose = true
    this.stopHeartbeat()
    
    if (this.connectionTimeout) {
      clearTimeout(this.connectionTimeout)
      this.connectionTimeout = null
    }
    
    if (this.ws) {
      this.ws.close(1000, 'Manual disconnect')
      this.ws = null
    }
  }

  // Send message to server
  send(message: WebSocketMessage): boolean {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      console.warn('WebSocket未连接，无法发送消息')
      return false
    }

    try {
      this.ws.send(JSON.stringify(message))
      return true
    } catch (error) {
      console.error('发送WebSocket消息失败:', error)
      return false
    }
  }

  // Get current connection status
  getConnectionStatus(): ConnectionStatus {
    if (!this.ws) {
      return 'disconnected'
    }

    switch (this.ws.readyState) {
      case WebSocket.CONNECTING:
        return this.reconnectAttempts > 0 ? 'reconnecting' : 'connecting'
      case WebSocket.OPEN:
        return 'connected'
      case WebSocket.CLOSING:
      case WebSocket.CLOSED:
        return 'disconnected'
      default:
        return 'error'
    }
  }

  // Check if WebSocket is connected
  isConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN
  }

  // Attempt to reconnect
  private async attemptReconnect(roomId: string, playerId: string): Promise<void> {
    if (this.isManualClose) {
      return
    }

    this.reconnectAttempts++
    console.log(`尝试重连 (${this.reconnectAttempts}/${this.maxReconnectAttempts})`)
    
    if (this.onReconnecting) {
      this.onReconnecting()
    }

    // Exponential backoff
    const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1)
    
    setTimeout(async () => {
      try {
        await this.connect(roomId, playerId)
      } catch (error) {
        console.error('重连失败:', error)
        
        if (this.reconnectAttempts >= this.maxReconnectAttempts) {
          console.error('达到最大重连次数，停止重连')
          if (this.onError) {
            this.onError(new Event('Max reconnect attempts reached'))
          }
        }
      }
    }, delay)
  }

  // Start heartbeat to keep connection alive
  private startHeartbeat(): void {
    this.stopHeartbeat()
    
    this.heartbeatInterval = window.setInterval(() => {
      if (this.isConnected()) {
        this.send({
          type: 'ping',
          data: {},
          timestamp: new Date().toISOString()
        })
      }
    }, 30000) // Send ping every 30 seconds
  }

  // Stop heartbeat
  private stopHeartbeat(): void {
    if (this.heartbeatInterval) {
      clearInterval(this.heartbeatInterval)
      this.heartbeatInterval = null
    }
  }

  // Set maximum reconnect attempts
  setMaxReconnectAttempts(attempts: number): void {
    this.maxReconnectAttempts = Math.max(0, attempts)
  }

  // Set reconnect delay
  setReconnectDelay(delay: number): void {
    this.reconnectDelay = Math.max(1000, delay)
  }

  // Reset reconnect attempts counter
  resetReconnectAttempts(): void {
    this.reconnectAttempts = 0
  }
}

// Singleton instance for global use
let globalWebSocketService: WebSocketService | null = null

export const getGlobalWebSocketService = (): WebSocketService => {
  if (!globalWebSocketService) {
    globalWebSocketService = new WebSocketService()
  }
  return globalWebSocketService
}

export const destroyGlobalWebSocketService = (): void => {
  if (globalWebSocketService) {
    globalWebSocketService.disconnect()
    globalWebSocketService = null
  }
}

export default WebSocketService