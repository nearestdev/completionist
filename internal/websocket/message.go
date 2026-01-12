package websocket

import "encoding/json"

type MessageType string

const (
	MessageTypeDM         MessageType = "dm"
	MessageTypeRoomChat   MessageType = "room_chat"
	MessageTypeRoomJoin   MessageType = "room_join"
	MessageTypeRoomLeave  MessageType = "room_leave"
	MessageTypeRoomState  MessageType = "room_state"
	MessageTypeTyping     MessageType = "typing"
	MessageTypePing       MessageType = "ping"
	MessageTypePong       MessageType = "pong"
	MessageTypeError      MessageType = "error"
	MessageTypeRoomUsers  MessageType = "room_users"
)

type WSMessage struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type DMPayload struct {
	ID               int64  `json:"id"`
	SenderID         int64  `json:"senderId"`
	SenderUsername   string `json:"senderUsername"`
	ReceiverID       int64  `json:"receiverId"`
	ReceiverUsername string `json:"receiverUsername"`
	Content          string `json:"content"`
	IsRead           bool   `json:"isRead"`
	CreatedAt        string `json:"createdAt"`
}

type RoomChatPayload struct {
	ID        int64  `json:"id"`
	RoomID    int64  `json:"roomId"`
	UserID    int64  `json:"userId"`
	Username  string `json:"username"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
}

type RoomJoinPayload struct {
	RoomID   int64  `json:"roomId"`
	UserID   int64  `json:"userId"`
	Username string `json:"username"`
}

type RoomLeavePayload struct {
	RoomID   int64  `json:"roomId"`
	UserID   int64  `json:"userId"`
	Username string `json:"username"`
}

type RoomUsersPayload struct {
	RoomID int64   `json:"roomId"`
	Users  []int64 `json:"users"`
}

type RoomStatePayload struct {
	RoomID            int64   `json:"roomId"`
	CurrentMediaURL   *string `json:"currentMediaUrl,omitempty"`
	CurrentMediaTitle *string `json:"currentMediaTitle,omitempty"`
	CurrentPositionMs int64   `json:"currentPositionMs"`
	IsPlaying         bool    `json:"isPlaying"`
	UpdatedBy         int64   `json:"updatedBy"`
	UpdatedByUsername string  `json:"updatedByUsername"`
	UpdatedAt         string  `json:"updatedAt"`
}

type TypingPayload struct {
	RoomID   *int64 `json:"roomId,omitempty"`
	UserID   int64  `json:"userId"`
	Username string `json:"username"`
}

type ErrorPayload struct {
	Message string `json:"message"`
}
