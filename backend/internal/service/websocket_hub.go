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
	}

	log.Printf("Client %s registered for room %s", client.ID, client.RoomID)
}

// unregisterClient unregisters a client
func (h *Hub) unregisterClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.Send)

		// Remove from room
		if client.RoomID != "" && h.rooms[client.RoomID] != nil {
			delete(h.rooms[client.RoomID], client)
			if len(h.rooms[client.RoomID]) == 0 {
				delete(h.rooms, client.RoomID)
			}
		}

		// Handle player leaving
		if client.Player != nil {
			h.gameService.HandlePlayerDisconnect(client.RoomID, client.Player.ID)
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
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	// Get player information
	room := h.gameService.GetRoom(roomID)
	if room == nil {
		conn.Close()
		return
	}

	player := room.GetPlayer(playerID)
	if player == nil {
		conn.Close()
		return
	}

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
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, messageBytes, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

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
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current message
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

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
		c.Hub.BroadcastToRoom(c.RoomID, model.WSMessage{
			Type: "room_updated",
			Data: updateData,
		})
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