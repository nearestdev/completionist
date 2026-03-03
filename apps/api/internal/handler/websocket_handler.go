package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/GATEOPENERZ/completionist-api/internal/middleware"
	"github.com/GATEOPENERZ/completionist-api/internal/websocket"
	ws "github.com/gorilla/websocket"
)

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.UserRepo.FindByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := websocket.NewClient(h.WSHub, conn, userID, user.Username)
	h.WSHub.RegisterClient(client)

	if roomIDs, err := h.MessagingRepo.GetUserRooms(userID); err == nil {
		for _, roomID := range roomIDs {
			h.WSHub.AddClientToRoom(roomID, userID)

			onlineUsers := h.WSHub.GetRoomOnlineUsers(roomID)
			payloadBytes, _ := json.Marshal(websocket.RoomUsersPayload{
				RoomID: roomID,
				Users:  onlineUsers,
			})
			h.WSHub.SendToRoom(roomID, &websocket.WSMessage{
				Type:    websocket.MessageTypeRoomUsers,
				Payload: payloadBytes,
			})
		}
	}

	go client.WritePump()
	go client.ReadPump()

	log.Printf("WebSocket connection established for user: %s (ID: %d)", user.Username, userID)
}
