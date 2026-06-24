package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/nearestdev/completionist/internal/models"
)

type ClientMessage struct {
	Client  *Client
	Message *WSMessage
}

type MessagingRepository interface {
	SendRoomMessage(msg *models.RoomMessage) (*models.RoomMessage, error)
	JoinRoom(roomID, userID int64) error
	LeaveRoom(roomID, userID int64) error
	UpdateRoomState(roomID, userID int64, req *models.UpdateRoomStateRequest) (*models.RoomState, error)
}

type Hub struct {
	clients       map[int64]*Client
	rooms         map[int64]map[int64]bool
	register      chan *Client
	unregister    chan *Client
	handleMessage chan *ClientMessage
	mu            sync.RWMutex
	repo          MessagingRepository
}

func NewHub(repo MessagingRepository) *Hub {
	return &Hub{
		clients:       make(map[int64]*Client),
		rooms:         make(map[int64]map[int64]bool),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		handleMessage: make(chan *ClientMessage),
		repo:          repo,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.userID] = client
			h.mu.Unlock()
			log.Printf("Client registered: userID=%d, username=%s", client.userID, client.username)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.userID]; ok {
				delete(h.clients, client.userID)
				close(client.send)
				for roomID := range client.rooms {
					if members, exists := h.rooms[roomID]; exists {
						delete(members, client.userID)
						if len(members) == 0 {
							delete(h.rooms, roomID)
						}
					}
				}
				log.Printf("Client unregistered: userID=%d, username=%s", client.userID, client.username)
			}
			h.mu.Unlock()

		case msg := <-h.handleMessage:
			h.processMessage(msg)
		}
	}
}

func (h *Hub) processMessage(msg *ClientMessage) {
	log.Printf("Processing message type: %s from userID=%d", msg.Message.Type, msg.Client.userID)

	switch msg.Message.Type {
	case MessageTypeRoomChat:
		var payload RoomChatPayload
		if err := json.Unmarshal(msg.Message.Payload, &payload); err != nil {
			log.Printf("Error unmarshaling room chat payload: %v", err)
			return
		}

		payload.UserID = msg.Client.userID
		payload.Username = msg.Client.username

		savedMsg, err := h.repo.SendRoomMessage(&models.RoomMessage{
			RoomID:  payload.RoomID,
			UserID:  payload.UserID,
			Content: payload.Content,
		})
		if err != nil {
			log.Printf("Error saving room message: %v", err)
			return
		}

		payload.ID = savedMsg.ID
		payload.CreatedAt = savedMsg.CreatedAt.Format(time.RFC3339)

		responsePayload, _ := json.Marshal(payload)
		responseMsg := &WSMessage{
			Type:    MessageTypeRoomChat,
			Payload: responsePayload,
		}

		h.SendToRoom(payload.RoomID, responseMsg)
	}
}

func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

func (h *Hub) UnregisterClient(client *Client) {
	h.unregister <- client
}

func (h *Hub) SendToUser(userID int64, msg *WSMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if client, ok := h.clients[userID]; ok {
		client.SendMessage(msg)
	}
}

func (h *Hub) SendToRoom(roomID int64, msg *WSMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if members, ok := h.rooms[roomID]; ok {
		for userID := range members {
			if client, exists := h.clients[userID]; exists {
				client.SendMessage(msg)
			}
		}
	}
}

func (h *Hub) SendToRoomExcept(roomID, exceptUserID int64, msg *WSMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if members, ok := h.rooms[roomID]; ok {
		for userID := range members {
			if userID != exceptUserID {
				if client, exists := h.clients[userID]; exists {
					client.SendMessage(msg)
				}
			}
		}
	}
}

func (h *Hub) AddClientToRoom(roomID, userID int64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.rooms[roomID]; !ok {
		h.rooms[roomID] = make(map[int64]bool)
	}
	h.rooms[roomID][userID] = true
	if client, ok := h.clients[userID]; ok {
		client.JoinRoom(roomID)
	}
}

func (h *Hub) RemoveClientFromRoom(roomID, userID int64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if members, ok := h.rooms[roomID]; ok {
		delete(members, userID)
		if len(members) == 0 {
			delete(h.rooms, roomID)
		}
	}
	if client, ok := h.clients[userID]; ok {
		client.LeaveRoom(roomID)
	}
}

func (h *Hub) GetClient(userID int64) *Client {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.clients[userID]
}

func (h *Hub) IsUserOnline(userID int64) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.clients[userID]
	return ok
}

func (h *Hub) GetRoomOnlineUsers(roomID int64) []int64 {
	h.mu.RLock()
	defer h.mu.RUnlock()
	users := []int64{}
	if members, ok := h.rooms[roomID]; ok {
		for userID := range members {
			users = append(users, userID)
		}
	}
	return users
}
