package models

import "time"

type DirectMessage struct {
	ID         int64     `db:"id" json:"id"`
	SenderID   int64     `db:"sender_id" json:"senderId"`
	ReceiverID int64     `db:"receiver_id" json:"receiverId"`
	Content    string    `db:"content" json:"content"`
	IsRead     bool      `db:"is_read" json:"isRead"`
	ReplyToID  *int64    `db:"reply_to_id" json:"replyToId,omitempty"`
	IsPinned   bool      `db:"is_pinned" json:"isPinned"`
	CreatedAt  time.Time `db:"created_at" json:"createdAt"`
}

type SendDirectMessageRequest struct {
	ReceiverID int64  `json:"receiverId"`
	Content    string `json:"content"`
	ReplyToID  *int64 `json:"replyToId,omitempty"`
}

type DirectMessageResponse struct {
	ID               int64                   `db:"id" json:"id"`
	SenderID         int64                   `db:"sender_id" json:"senderId"`
	SenderUsername   string                  `db:"sender_username" json:"senderUsername"`
	ReceiverID       int64                   `db:"receiver_id" json:"receiverId"`
	ReceiverUsername string                  `db:"receiver_username" json:"receiverUsername"`
	Content          string                  `db:"content" json:"content"`
	IsRead           bool                    `db:"is_read" json:"isRead"`
	ReplyToID        *int64                  `db:"reply_to_id" json:"replyToId,omitempty"`
	IsPinned         bool                    `db:"is_pinned" json:"isPinned"`
	Reactions        []DirectMessageReaction `json:"reactions"`
	Attachments      []Attachment            `json:"attachments"`
	CreatedAt        time.Time               `db:"created_at" json:"createdAt"`
}

type DirectMessageReaction struct {
	ID        int64     `db:"id" json:"id"`
	MessageID int64     `db:"message_id" json:"messageId"`
	UserID    int64     `db:"user_id" json:"userId"`
	Username  string    `db:"username" json:"username"`
	Reaction  string    `db:"reaction" json:"reaction"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

type ConversationSummary struct {
	OtherUserID         int64     `db:"other_user_id" json:"otherUserId"`
	OtherUsername       string    `db:"other_username" json:"otherUsername"`
	LastMessage         string    `db:"last_message" json:"lastMessage"`
	LastMessageTime     time.Time `db:"last_message_time" json:"lastMessageTime"`
	UnreadCount         int64     `db:"unread_count" json:"unreadCount"`
	IsSentByCurrentUser bool      `db:"is_sent_by_current_user" json:"isSentByCurrentUser"`
}

type RoomType string

const (
	RoomTypeChat  RoomType = "chat"
	RoomTypeMusic RoomType = "music"
	RoomTypeWatch RoomType = "watch"
)

type Room struct {
	ID          int64     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description *string   `db:"description" json:"description,omitempty"`
	RoomType    RoomType  `db:"room_type" json:"roomType"`
	CreatorID   int64     `db:"creator_id" json:"creatorId"`
	IsPublic    bool      `db:"is_public" json:"isPublic"`
	MaxMembers  int       `db:"max_members" json:"maxMembers"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
}

type CreateRoomRequest struct {
	Name        string   `json:"name"`
	Description *string  `json:"description,omitempty"`
	RoomType    RoomType `json:"roomType"`
	IsPublic    bool     `json:"isPublic"`
	MaxMembers  int      `json:"maxMembers"`
}

type RoomResponse struct {
	ID          int64     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description *string   `db:"description" json:"description,omitempty"`
	RoomType    RoomType  `db:"room_type" json:"roomType"`
	CreatorID   int64     `db:"creator_id" json:"creatorId"`
	CreatorName string    `db:"creator_name" json:"creatorName"`
	IsPublic    bool      `db:"is_public" json:"isPublic"`
	MaxMembers  int       `db:"max_members" json:"maxMembers"`
	MemberCount int       `db:"member_count" json:"memberCount"`
	IsMember    bool      `db:"is_member" json:"isMember"`
	IsModerator bool      `db:"is_moderator" json:"isModerator"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
}

type RoomMember struct {
	ID          int64     `db:"id" json:"id"`
	RoomID      int64     `db:"room_id" json:"roomId"`
	UserID      int64     `db:"user_id" json:"userId"`
	IsModerator bool      `db:"is_moderator" json:"isModerator"`
	JoinedAt    time.Time `db:"joined_at" json:"joinedAt"`
}

type RoomMemberResponse struct {
	UserID      int64     `db:"user_id" json:"userId"`
	Username    string    `db:"username" json:"username"`
	IsModerator bool      `db:"is_moderator" json:"isModerator"`
	JoinedAt    time.Time `db:"joined_at" json:"joinedAt"`
}

type RoomMessage struct {
	ID        int64     `db:"id" json:"id"`
	RoomID    int64     `db:"room_id" json:"roomId"`
	UserID    int64     `db:"user_id" json:"userId"`
	Content   string    `db:"content" json:"content"`
	ReplyToID *int64    `db:"reply_to_id" json:"replyToId,omitempty"`
	IsPinned  bool      `db:"is_pinned" json:"isPinned"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

type SendRoomMessageRequest struct {
	Content   string `json:"content"`
	ReplyToID *int64 `json:"replyToId,omitempty"`
}

type RoomMessageResponse struct {
	ID          int64                 `db:"id" json:"id"`
	RoomID      int64                 `db:"room_id" json:"roomId"`
	UserID      int64                 `db:"user_id" json:"userId"`
	Username    string                `db:"username" json:"username"`
	Content     string                `db:"content" json:"content"`
	ReplyToID   *int64                `db:"reply_to_id" json:"replyToId,omitempty"`
	IsPinned    bool                  `db:"is_pinned" json:"isPinned"`
	Reactions   []RoomMessageReaction `json:"reactions"`
	Attachments []Attachment          `json:"attachments"`
	CreatedAt   time.Time             `db:"created_at" json:"createdAt"`
}

type RoomMessageReaction struct {
	ID        int64     `db:"id" json:"id"`
	MessageID int64     `db:"message_id" json:"messageId"`
	UserID    int64     `db:"user_id" json:"userId"`
	Username  string    `db:"username" json:"username"`
	Reaction  string    `db:"reaction" json:"reaction"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

type RoomState struct {
	RoomID            int64     `db:"room_id" json:"roomId"`
	CurrentMediaURL   *string   `db:"current_media_url" json:"currentMediaUrl,omitempty"`
	CurrentMediaTitle *string   `db:"current_media_title" json:"currentMediaTitle,omitempty"`
	CurrentPositionMs int64     `db:"current_position_ms" json:"currentPositionMs"`
	IsPlaying         bool      `db:"is_playing" json:"isPlaying"`
	UpdatedBy         *int64    `db:"updated_by" json:"updatedBy,omitempty"`
	UpdatedAt         time.Time `db:"updated_at" json:"updatedAt"`
}

type UpdateRoomStateRequest struct {
	CurrentMediaURL   *string `json:"currentMediaUrl,omitempty"`
	CurrentMediaTitle *string `json:"currentMediaTitle,omitempty"`
	CurrentPositionMs *int64  `json:"currentPositionMs,omitempty"`
	IsPlaying         *bool   `json:"isPlaying,omitempty"`
}

