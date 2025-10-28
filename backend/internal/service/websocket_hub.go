package service

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"gomoku-backend/internal/model"
)

// Client represents a WebSocket client
type Client struct {
	ID     string
	Conn   *websocket.Conn
	Send   chan []byte
	RoomID string
	Player *model.PVPPlayer
	Hub    *Hub
}

// Hub maintains the set of active clients and broadcasts messages to the clients
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Inbound messages from the clients
	broadcast chan []byte

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Room-specific clients
	rooms map[string]map[*Client]bool

	// Mutex for thread safety
	mutex sync.RWMutex

	// Game service reference
	gameService *GameService
}

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from any origin in development
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// NewHub creates a new WebSocket hub
func NewHub(gameService *GameService) *Hub {
	return &Hub{
		clients:     make(map[*Client]bool),
		broadcast:   make(chan []byte),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		rooms:       make(map[string]map[*Client]bool),
		gameService: gameService,
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastMessage(message)
		}
	}
}

// handlePingMessage handles ping messages for keeping connection alive
func (c *Client) handlePingMessage(wsMessage *model.WSMessage) {
	// Respond with pong message
	pongMessage := model.WSMessage{
		Type: "pong",
		Data: map[string]interface{}{
			"timestamp": wsMessage.Data,
		},
	}
	
	// Send pong response to the client
	c.Hub.sendToClient(c, pongMessage)
}

// registerClient registers a new client
func (h *Hub) registerClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.clients[client] = true

	// Add client to room
	if client.RoomID != "" {
		if h.rooms[client.RoomID] == nil {
			h.rooms[client.RoomID] = make(map[*Client]bool)
		}
		h.rooms[client.RoomID][client] = true
		
		// Get current room state
		room := h.gameService.GetRoom(client.RoomID)
		if room != nil {
			log.Printf("客户端注册成功，广播房间状态: roomID=%s, playerID=%s, 房间玩家数=%d", 
				client.RoomID, client.ID, len(room.Players))
			
			// Broadcast updated room state to ALL clients in the room
			updateData := model.RoomUpdateData{
				Room:   room,
				Player: client.Player,
			}
			
			message := model.WSMessage{
				Type: "room_updated",
				Data: updateData,
			}
			
			// Broadcast to all clients in the room (including the new one)
			log.Printf("registerClient: 准备调用broadcastToRoomInternal - 房间ID: %s", client.RoomID)
			h.broadcastToRoomInternal(client.RoomID, message)
			log.Printf("registerClient: broadcastToRoomInternal调用完成 - 房间ID: %s", client.RoomID)
			
			// Also send current game state if game is in progress
			if room.Status == "playing" && room.Game != nil {
				gameMessage := model.WSMessage{
					Type: "game_update",
					Data: model.GameUpdateData{
						Game: room.Game,
					},
				}
				h.broadcastToRoomInternal(client.RoomID, gameMessage)
			}
			
			// Send player joined notification to other clients
			joinedMessage := model.WSMessage{
				Type: "player_joined",
				Data: map[string]interface{}{
					"player": client.Player,
					"room":   room,
				},
			}
			
			// Send to all clients in room except the newly joined one
			for roomClient := range h.rooms[client.RoomID] {
				if roomClient.ID != client.ID {
					h.sendToClient(roomClient, joinedMessage)
				}
			}
		}
	}

	log.Printf("Client %s registered for room %s, total clients in room: %d", 
		client.ID, client.RoomID, len(h.rooms[client.RoomID]))
}

// unregisterClient unregisters a client
func (h *Hub) unregisterClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, ok := h.clients[client]; ok {
		log.Printf("客户端断开连接: roomID=%s, playerID=%s", client.RoomID, client.ID)
		
		delete(h.clients, client)
		close(client.Send)

		// Remove from room
		if client.RoomID != "" && h.rooms[client.RoomID] != nil {
			delete(h.rooms[client.RoomID], client)
			remainingClients := len(h.rooms[client.RoomID])
			
			if remainingClients == 0 {
				delete(h.rooms, client.RoomID)
				log.Printf("房间 %s 已清空，删除房间", client.RoomID)
			} else {
				log.Printf("房间 %s 剩余客户端数: %d", client.RoomID, remainingClients)
			}
		}

		// Handle player leaving
		if client.Player != nil {
			h.gameService.HandlePlayerDisconnect(client.RoomID, client.Player.ID)
			
			// Get updated room state after player disconnect
			room := h.gameService.GetRoom(client.RoomID)
			if room != nil {
				log.Printf("广播玩家离开事件: roomID=%s, playerID=%s, 房间剩余玩家数=%d", 
					client.RoomID, client.ID, len(room.Players))
				
				// Broadcast player left notification
				leftMessage := model.WSMessage{
					Type: "player_left",
					Data: map[string]interface{}{
						"player": client.Player,
						"room":   room,
					},
				}
				h.broadcastToRoomInternal(client.RoomID, leftMessage)
				
				// Also broadcast updated room state
				updateData := model.RoomUpdateData{
					Room:   room,
					Player: client.Player,
				}
				roomUpdateMessage := model.WSMessage{
					Type: "room_updated",
					Data: updateData,
				}
				h.broadcastToRoomInternal(client.RoomID, roomUpdateMessage)
			}
		}

		log.Printf("Client %s unregistered from room %s", client.ID, client.RoomID)
	}
}

// broadcastMessage broadcasts a message to all clients
func (h *Hub) broadcastMessage(message []byte) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	for client := range h.clients {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(h.clients, client)
		}
	}
}

// BroadcastToRoom broadcasts a message to all clients in a specific room
func (h *Hub) BroadcastToRoom(roomID string, message interface{}) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	h.broadcastToRoomInternal(roomID, message)
}

// broadcastToRoomInternal broadcasts a message to all clients in a specific room
// This method assumes the caller already holds the appropriate lock
func (h *Hub) broadcastToRoomInternal(roomID string, message interface{}) {
	log.Printf("开始广播消息到房间 %s", roomID)
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	if roomClients, exists := h.rooms[roomID]; exists {
		log.Printf("房间 %s 存在，客户端数量: %d", roomID, len(roomClients))
		clientCount := 0
		for client := range roomClients {
			select {
			case client.Send <- messageBytes:
				clientCount++
				log.Printf("消息已发送给客户端 %s", client.ID)
			default:
				log.Printf("客户端 %s 发送失败，清理连接", client.ID)
				close(client.Send)
				delete(h.clients, client)
				delete(roomClients, client)
			}
		}
		log.Printf("消息广播完成，成功发送给 %d 个客户端", clientCount)
	} else {
		log.Printf("房间 %s 不存在", roomID)
	}
}

// sendToClient sends a message to a specific client
func (h *Hub) sendToClient(client *Client, message interface{}) {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message for client %s: %v", client.ID, err)
		return
	}

	select {
	case client.Send <- messageBytes:
		log.Printf("消息已发送给客户端 %s", client.ID)
	default:
		log.Printf("客户端 %s 发送失败，清理连接", client.ID)
		close(client.Send)
		h.mutex.Lock()
		delete(h.clients, client)
		if roomClients, exists := h.rooms[client.RoomID]; exists {
			delete(roomClients, client)
		}
		h.mutex.Unlock()
	}
}

// GetRoomClients returns the number of clients in a room
func (h *Hub) GetRoomClients(roomID string) int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if roomClients, exists := h.rooms[roomID]; exists {
		return len(roomClients)
	}
	return 0
}

// ServeWS handles websocket requests from the peer
func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request, roomID, playerID string) {
	log.Printf("WebSocket连接请求: roomID=%s, playerID=%s", roomID, playerID)
	
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	// Get player information
	room := h.gameService.GetRoom(roomID)
	if room == nil {
		log.Printf("WebSocket连接失败: 房间不存在 roomID=%s", roomID)
		// Send error message before closing
		errorMsg := model.WSMessage{
			Type: "error",
			Data: map[string]interface{}{
				"message": "房间不存在或已关闭",
				"code":    "ROOM_NOT_FOUND",
			},
		}
		if msgBytes, err := json.Marshal(errorMsg); err == nil {
			conn.WriteMessage(websocket.TextMessage, msgBytes)
		}
		conn.Close()
		return
	}

	player := room.GetPlayer(playerID)
	if player == nil {
		log.Printf("WebSocket连接失败: 玩家不在房间中 roomID=%s, playerID=%s", roomID, playerID)
		// Send error message before closing
		errorMsg := model.WSMessage{
			Type: "error",
			Data: map[string]interface{}{
				"message": "玩家不在此房间中",
				"code":    "PLAYER_NOT_IN_ROOM",
			},
		}
		if msgBytes, err := json.Marshal(errorMsg); err == nil {
			conn.WriteMessage(websocket.TextMessage, msgBytes)
		}
		conn.Close()
		return
	}

	log.Printf("WebSocket连接成功: roomID=%s, playerID=%s, playerName=%s", roomID, playerID, player.Name)

	client := &Client{
		ID:     playerID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		RoomID: roomID,
		Player: player,
		Hub:    h,
	}

	// Register the client
	h.register <- client

	// Start goroutines for reading and writing
	go client.writePump()
	go client.readPump()
}

// readPump pumps messages from the websocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		log.Printf("readPump结束 - 玩家: %s, 房间: %s", c.Player.Name, c.RoomID)
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	log.Printf("readPump开始 - 玩家: %s, 房间: %s", c.Player.Name, c.RoomID)
	c.Conn.SetReadLimit(512)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		log.Printf("等待读取消息 - 玩家: %s", c.Player.Name)
		_, messageBytes, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("读取消息错误 - 玩家: %s, 错误: %v", c.Player.Name, err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		log.Printf("收到原始消息 - 玩家: %s, 长度: %d, 内容: %s", c.Player.Name, len(messageBytes), string(messageBytes))

		var wsMessage model.WSMessage
		if err := json.Unmarshal(messageBytes, &wsMessage); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		c.handleMessage(&wsMessage)
	}
}

// writePump pumps messages from the hub to the websocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			log.Printf("writePump: 准备发送消息给客户端 %s, 消息长度: %d", c.ID, len(message))
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				log.Printf("writePump: 发送通道已关闭，客户端 %s", c.ID)
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("writePump: NextWriter错误，客户端 %s: %v", c.ID, err)
				return
			}
			w.Write(message)
			log.Printf("writePump: 消息内容已写入，客户端 %s: %s", c.ID, string(message))

			// Add queued messages to the current message
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				log.Printf("writePump: Close错误，客户端 %s: %v", c.ID, err)
				return
			}
			log.Printf("writePump: 消息发送完成，客户端 %s", c.ID)

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage handles incoming WebSocket messages
func (c *Client) handleMessage(wsMessage *model.WSMessage) {
    log.Printf("收到WebSocket消息 - 类型: %s, 玩家: %s, 房间: %s", wsMessage.Type, c.Player.Name, c.RoomID)
    
    switch wsMessage.Type {
    case "join":
        c.handleJoinMessage(wsMessage)
    case "leave":
        c.handleLeaveMessage(wsMessage)
    case "move":
        c.handleMoveMessage(wsMessage)
    case "chat":
        c.handleChatMessage(wsMessage)
    case "ready":
        c.handleReadyMessage(wsMessage)
    case "ping":
        c.handlePingMessage(wsMessage)
    case "draw_offer":
        c.handleDrawOfferMessage(wsMessage)
    case "draw_response":
        c.handleDrawResponseMessage(wsMessage)
    case "resign":
        c.handleResignMessage(wsMessage)
    default:
        log.Printf("Unknown message type: %s", wsMessage.Type)
    }
}

// handleJoinMessage handles player join messages
func (c *Client) handleJoinMessage(wsMessage *model.WSMessage) {
	// Player is already joined when WebSocket connection is established
	// Send room update to all clients
	room := c.Hub.gameService.GetRoom(c.RoomID)
	if room != nil {
		updateData := model.RoomUpdateData{
			Room:   room,
			Player: c.Player,
		}
		c.Hub.BroadcastToRoom(c.RoomID, model.WSMessage{
			Type: "player_joined",
			Data: updateData,
		})
	}
}

// handleLeaveMessage handles player leave messages
func (c *Client) handleLeaveMessage(wsMessage *model.WSMessage) {
	c.Hub.gameService.LeaveRoom(c.RoomID, c.Player.ID)
	
	updateData := model.RoomUpdateData{
		Player: c.Player,
	}
	c.Hub.BroadcastToRoom(c.RoomID, model.WSMessage{
		Type: "player_left",
		Data: updateData,
	})
}

// handleMoveMessage handles game move messages
func (c *Client) handleMoveMessage(wsMessage *model.WSMessage) {
	data, ok := wsMessage.Data.(map[string]interface{})
	if !ok {
		return
	}

	x, xOk := data["x"].(float64)
	y, yOk := data["y"].(float64)
	if !xOk || !yOk {
		return
	}

	room, move, err := c.Hub.gameService.MakeMove(c.RoomID, c.Player.ID, int(x), int(y))
	if err != nil {
		log.Printf("Error making move: %v", err)
		return
	}

	if move != nil && room != nil && room.Game != nil {
		updateData := model.GameUpdateData{
			Game:     room.Game,
			LastMove: move,
		}
		c.Hub.BroadcastToRoom(c.RoomID, model.WSMessage{
			Type: "move_made",
			Data: updateData,
		})

		// Check if game is finished
		if room.Game.Status == "finished" {
			c.Hub.BroadcastToRoom(c.RoomID, model.WSMessage{
				Type: "game_ended",
				Data: updateData,
			})
		}
	}
}

// handleChatMessage handles chat messages
func (c *Client) handleChatMessage(wsMessage *model.WSMessage) {
	data, ok := wsMessage.Data.(map[string]interface{})
	if !ok {
		return
	}

	message, msgOk := data["message"].(string)
	if !msgOk {
		return
	}

	chatData := model.ChatMessageData{
		PlayerID:   c.Player.ID,
		PlayerName: c.Player.Name,
		Message:    message,
		Timestamp:  time.Now(),
	}

	c.Hub.BroadcastToRoom(c.RoomID, model.WSMessage{
		Type: "chat_message",
		Data: chatData,
	})
}

// handleReadyMessage handles player ready status changes
func (c *Client) handleReadyMessage(wsMessage *model.WSMessage) {
	log.Printf("收到准备状态消息 - 玩家: %s, 房间: %s", c.Player.Name, c.RoomID)
	
	data, ok := wsMessage.Data.(map[string]interface{})
	if !ok {
		log.Printf("无法解析消息数据")
		return
	}

	ready, readyOk := data["ready"].(bool)
	if !readyOk {
		log.Printf("无法解析ready字段")
		return
	}

	log.Printf("设置玩家 %s 准备状态为: %v", c.Player.Name, ready)
	err := c.Hub.gameService.SetPlayerReady(c.RoomID, c.Player.ID, ready)
	if err != nil {
		log.Printf("Error setting player ready: %v", err)
		return
	}
	log.Printf("玩家准备状态设置成功")

	room := c.Hub.gameService.GetRoom(c.RoomID)
	if room != nil {
		log.Printf("获取房间信息成功，准备广播房间更新")
		updateData := model.RoomUpdateData{
			Room:   room,
			Player: c.Player,
		}
		log.Printf("广播房间更新消息到房间 %s", c.RoomID)
		log.Printf("调用BroadcastToRoom前 - 房间ID: %s", c.RoomID)
		c.Hub.BroadcastToRoom(c.RoomID, model.WSMessage{
			Type: "room_updated",
			Data: updateData,
		})
		log.Printf("调用BroadcastToRoom后 - 房间ID: %s", c.RoomID)
		log.Printf("房间更新消息广播完成")

		// Check if game can start
		if room.CanStartGame() && room.Status == "waiting" {
			err := c.Hub.gameService.StartGame(c.RoomID)
			if err != nil {
				log.Printf("Error starting game: %v", err)
				return
			}

			// Get updated room with game
			room = c.Hub.gameService.GetRoom(c.RoomID)
			if room != nil && room.Game != nil {
				gameData := model.GameUpdateData{
					Game: room.Game,
				}
				c.Hub.BroadcastToRoom(c.RoomID, model.WSMessage{
					Type: "game_started",
					Data: gameData,
				})
			}
		}
	}
}

// handleDrawOfferMessage handles a player's draw offer
func (c *Client) handleDrawOfferMessage(wsMessage *model.WSMessage) {
    log.Printf("收到求和请求 - 玩家: %s, 房间: %s", c.Player.Name, c.RoomID)

    // 仅在游戏进行中允许求和
    room := c.Hub.gameService.GetRoom(c.RoomID)
    if room == nil || room.Game == nil || room.Game.Status != "playing" {
        log.Printf("求和请求被拒绝：游戏未进行中")
        // 返回错误给发起者
        c.Hub.sendToClient(c, model.WSMessage{
            Type: "error",
            Data: map[string]interface{}{
                "message": "当前无法求和",
                "code":    "DRAW_NOT_ALLOWED",
            },
        })
        return
    }

    // 广播求和请求给房间内其他玩家
    c.Hub.BroadcastToRoom(c.RoomID, model.WSMessage{
        Type: "draw_offer",
        Data: map[string]interface{}{
            "fromPlayerId":   c.Player.ID,
            "fromPlayerName": c.Player.Name,
            "timestamp":      time.Now(),
        },
    })
}

// handleDrawResponseMessage handles draw accept/reject from the other player
func (c *Client) handleDrawResponseMessage(wsMessage *model.WSMessage) {
    log.Printf("收到求和回应 - 玩家: %s, 房间: %s", c.Player.Name, c.RoomID)

    data, ok := wsMessage.Data.(map[string]interface{})
    if !ok {
        return
    }

    accept, ok := data["accept"].(bool)
    if !ok {
        return
    }

    room := c.Hub.gameService.GetRoom(c.RoomID)
    if room == nil || room.Game == nil || room.Game.Status != "playing" {
        log.Printf("求和回应忽略：游戏未进行中")
        return
    }

    if accept {
        // 标记为平局结束
        room.Game.Status = "finished"
        room.Game.Winner = "" // 无胜者表示平局
        now := time.Now()
        room.Game.EndedAt = &now

        // 广播游戏结束（平局）
        updateData := model.GameUpdateData{
            Game: room.Game,
        }
        c.Hub.BroadcastToRoom(c.RoomID, model.WSMessage{
            Type: "game_ended",
            Data: updateData,
        })

        // 额外广播求和已接受事件，便于前端提示
        c.Hub.BroadcastToRoom(c.RoomID, model.WSMessage{
            Type: "draw_accepted",
            Data: map[string]interface{}{
                "byPlayerId":   c.Player.ID,
                "byPlayerName": c.Player.Name,
                "timestamp":    time.Now(),
            },
        })
    } else {
        // 广播求和被拒绝
        c.Hub.BroadcastToRoom(c.RoomID, model.WSMessage{
            Type: "draw_rejected",
            Data: map[string]interface{}{
                "byPlayerId":   c.Player.ID,
                "byPlayerName": c.Player.Name,
                "timestamp":    time.Now(),
            },
        })
    }
}

// handleResignMessage handles resign action from a player
func (c *Client) handleResignMessage(wsMessage *model.WSMessage) {
    log.Printf("收到认输请求 - 玩家: %s, 房间: %s", c.Player.Name, c.RoomID)

    room := c.Hub.gameService.GetRoom(c.RoomID)
    if room == nil || room.Game == nil || room.Game.Status != "playing" {
        log.Printf("认输请求忽略：游戏未进行中")
        return
    }

    // 设置另一位玩家为胜者
    var winnerID string
    for _, p := range room.Players {
        if p.ID != c.Player.ID {
            winnerID = p.ID
            break
        }
    }
    room.Game.Status = "finished"
    room.Game.Winner = winnerID
    now := time.Now()
    room.Game.EndedAt = &now

    updateData := model.GameUpdateData{Game: room.Game}
    c.Hub.BroadcastToRoom(c.RoomID, model.WSMessage{
        Type: "game_ended",
        Data: updateData,
    })

    // 额外广播认输事件，便于前端提示
    c.Hub.BroadcastToRoom(c.RoomID, model.WSMessage{
        Type: "resigned",
        Data: map[string]interface{}{
            "playerId":   c.Player.ID,
            "playerName": c.Player.Name,
            "timestamp":  time.Now(),
        },
    })
}