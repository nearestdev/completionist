package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/nearestdev/completionist/internal/httpx"
	"github.com/nearestdev/completionist/internal/middleware"
	"github.com/nearestdev/completionist/internal/models"
	"github.com/nearestdev/completionist/internal/websocket"
)

func (h *Handler) SendDirectMessage(w http.ResponseWriter, r *http.Request) {
	senderID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req models.SendDirectMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Content == "" {
		httpx.JSONError(w, http.StatusBadRequest, "Content cannot be empty")
		return
	}

	if senderID == req.ReceiverID {
		httpx.JSONError(w, http.StatusBadRequest, "Cannot send message to yourself")
		return
	}

	receiver, err := h.UserRepo.FindByID(req.ReceiverID)
	if err != nil {
		httpx.JSONError(w, http.StatusNotFound, "Receiver not found")
		return
	}

	sender, err := h.UserRepo.FindByID(senderID)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to fetch sender details")
		return
	}

	dm := &models.DirectMessage{
		SenderID:   senderID,
		ReceiverID: req.ReceiverID,
		Content:    req.Content,
		ReplyToID:  req.ReplyToID,
	}

	result, err := h.MessagingRepo.SendDirectMessage(dm)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to send message")
		return
	}

	payload := websocket.DMPayload{
		ID:               result.ID,
		SenderID:         result.SenderID,
		SenderUsername:   sender.Username,
		ReceiverID:       result.ReceiverID,
		ReceiverUsername: receiver.Username,
		Content:          result.Content,
		IsRead:           result.IsRead,
		ReplyToID:        result.ReplyToID,
		IsPinned:         result.IsPinned,
		CreatedAt:        result.CreatedAt.Format(time.RFC3339),
	}

	payloadBytes, _ := json.Marshal(payload)
	wsMsg := &websocket.WSMessage{
		Type:    websocket.MessageTypeDM,
		Payload: payloadBytes,
	}

	h.WSHub.SendToUser(req.ReceiverID, wsMsg)
	h.WSHub.SendToUser(senderID, wsMsg)

	httpx.JSON(w, http.StatusCreated, result)
}

func (h *Handler) GetConversation(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	otherUserIDStr := chi.URLParam(r, "userId")
	otherUserID, err := strconv.ParseInt(otherUserIDStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 && v <= 100 {
			limit = v
		}
	}

	var beforeID int64
	if c := r.URL.Query().Get("cursor"); c != "" {
		if v, err := strconv.ParseInt(c, 10, 64); err == nil {
			beforeID = v
		}
	}

	messages, err := h.MessagingRepo.GetConversation(userID, otherUserID, limit, beforeID)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to get conversation")
		return
	}

	_ = h.MessagingRepo.MarkAsRead(userID, otherUserID)

	httpx.JSON(w, http.StatusOK, messages)
}

func (h *Handler) GetConversations(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	conversations, err := h.MessagingRepo.GetConversations(userID)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to get conversations")
		return
	}

	httpx.JSON(w, http.StatusOK, conversations)
}

func (h *Handler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	creatorID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req models.CreateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" {
		httpx.JSONError(w, http.StatusBadRequest, "Room name is required")
		return
	}

	if req.MaxMembers <= 0 {
		req.MaxMembers = 50
	}

	room := &models.Room{
		Name:        req.Name,
		Description: req.Description,
		RoomType:    req.RoomType,
		CreatorID:   creatorID,
		IsPublic:    req.IsPublic,
		MaxMembers:  req.MaxMembers,
	}

	result, err := h.MessagingRepo.CreateRoom(room)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to create room")
		return
	}

	h.WSHub.AddClientToRoom(result.ID, creatorID)

	httpx.JSON(w, http.StatusCreated, result)
}

func (h *Handler) GetRooms(w http.ResponseWriter, r *http.Request) {
	roomTypeStr := r.URL.Query().Get("type")
	var roomType *models.RoomType
	if roomTypeStr != "" {
		rt := models.RoomType(roomTypeStr)
		roomType = &rt
	}

	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 && v <= 100 {
			limit = v
		}
	}

	rooms, err := h.MessagingRepo.ListRooms(roomType, limit)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to get rooms")
		return
	}

	httpx.JSON(w, http.StatusOK, rooms)
}

func (h *Handler) GetRoom(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	roomIDStr := chi.URLParam(r, "roomId")
	roomID, err := strconv.ParseInt(roomIDStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid room ID")
		return
	}

	room, err := h.MessagingRepo.GetRoomWithDetails(roomID, userID)
	if err != nil {
		httpx.JSONError(w, http.StatusNotFound, "Room not found")
		return
	}

	httpx.JSON(w, http.StatusOK, room)
}

func (h *Handler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	roomIDStr := chi.URLParam(r, "roomId")
	roomID, err := strconv.ParseInt(roomIDStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid room ID")
		return
	}

	if err := h.MessagingRepo.JoinRoom(roomID, userID); err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.WSHub.AddClientToRoom(roomID, userID)

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) LeaveRoom(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	roomIDStr := chi.URLParam(r, "roomId")
	roomID, err := strconv.ParseInt(roomIDStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid room ID")
		return
	}

	if err := h.MessagingRepo.LeaveRoom(roomID, userID); err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to leave room")
		return
	}

	h.WSHub.RemoveClientFromRoom(roomID, userID)

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetRoomMessages(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	roomIDStr := chi.URLParam(r, "roomId")
	roomID, err := strconv.ParseInt(roomIDStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid room ID")
		return
	}

	isMember, err := h.MessagingRepo.IsRoomMember(roomID, userID)
	if err != nil || !isMember { // TODO: Allow public rooms to be viewed without joining? For now, enforcing membership
		httpx.JSONError(w, http.StatusForbidden, "Not a member of this room")
		return
	}

	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 && v <= 100 {
			limit = v
		}
	}

	messages, err := h.MessagingRepo.GetRoomMessages(roomID, limit)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to get messages")
		return
	}

	httpx.JSON(w, http.StatusOK, messages)
}

func (h *Handler) UpdateRoomState(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	roomIDStr := chi.URLParam(r, "roomId")
	roomID, err := strconv.ParseInt(roomIDStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid room ID")
		return
	}

	isMember, err := h.MessagingRepo.IsRoomMember(roomID, userID)
	if err != nil || !isMember {
		httpx.JSONError(w, http.StatusForbidden, "Not a member of this room")
		return
	}

	var req models.UpdateRoomStateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	state, err := h.MessagingRepo.UpdateRoomState(roomID, userID, &req)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to update room state")
		return
	}

	httpx.JSON(w, http.StatusOK, state)
}

func (h *Handler) GetRoomState(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	roomIDStr := chi.URLParam(r, "roomId")
	roomID, err := strconv.ParseInt(roomIDStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid room ID")
		return
	}

	isMember, err := h.MessagingRepo.IsRoomMember(roomID, userID)
	if err != nil || !isMember {
		httpx.JSONError(w, http.StatusForbidden, "Not a member of this room")
		return
	}

	state, err := h.MessagingRepo.GetRoomState(roomID)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to get room state")
		return
	}

	httpx.JSON(w, http.StatusOK, state)
}

func (h *Handler) GetRoomMembers(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	roomIDStr := chi.URLParam(r, "roomId")
	roomID, err := strconv.ParseInt(roomIDStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid room ID")
		return
	}

	isMember, err := h.MessagingRepo.IsRoomMember(roomID, userID)
	if err != nil || !isMember {
		httpx.JSONError(w, http.StatusForbidden, "Not a member of this room")
		return
	}

	members, err := h.MessagingRepo.GetRoomMembers(roomID)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to get room members")
		return
	}

	httpx.JSON(w, http.StatusOK, members)
}

func (h *Handler) UploadFile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	_ = userID

	r.ParseMultipartForm(10 << 20)

	file, header, err := r.FormFile("file")
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Failed to read file")
		return
	}
	defer file.Close()

	url, err := h.FileService.Upload(r.Context(), file, header.Filename, header.Header.Get("Content-Type"))
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to upload file")
		return
	}

	httpx.JSON(w, http.StatusOK, map[string]string{"url": url})
}

func (h *Handler) ReactToDirectMessage(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	msgIDStr := chi.URLParam(r, "messageId")
	msgID, err := strconv.ParseInt(msgIDStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid message ID")
		return
	}

	var req struct {
		Reaction string `json:"reaction"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.MessagingRepo.ReactToDirectMessage(userID, msgID, req.Reaction); err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to save reaction")
		return
	}

	user, _ := h.UserRepo.FindByID(userID)

	payload := websocket.ReactionPayload{
		MessageID:       msgID,
		UserID:          userID,
		Username:        user.Username,
		Reaction:        req.Reaction,
		IsDirectMessage: true,
	}

	payloadBytes, _ := json.Marshal(payload)
	wsMsg := &websocket.WSMessage{
		Type:    websocket.MessageTypeReaction,
		Payload: payloadBytes,
	}

	var dm struct {
		SenderID   int64 `db:"sender_id"`
		ReceiverID int64 `db:"receiver_id"`
	}
	h.MessagingRepo.DB.Get(&dm, "SELECT sender_id, receiver_id FROM direct_messages WHERE id = $1", msgID)

	h.WSHub.SendToUser(dm.SenderID, wsMsg)
	h.WSHub.SendToUser(dm.ReceiverID, wsMsg)

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ReactToRoomMessage(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	msgIDStr := chi.URLParam(r, "messageId")
	msgID, err := strconv.ParseInt(msgIDStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid message ID")
		return
	}

	var req struct {
		Reaction string `json:"reaction"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.MessagingRepo.ReactToRoomMessage(userID, msgID, req.Reaction); err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to save reaction")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) PinDirectMessage(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	msgIDStr := chi.URLParam(r, "messageId")
	msgID, err := strconv.ParseInt(msgIDStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid message ID")
		return
	}

	var dm struct {
		IsPinned   bool  `db:"is_pinned"`
		SenderID   int64 `db:"sender_id"`
		ReceiverID int64 `db:"receiver_id"`
	}
	err = h.MessagingRepo.DB.Get(&dm, "SELECT is_pinned, sender_id, receiver_id FROM direct_messages WHERE id = $1", msgID)
	if err != nil {
		httpx.JSONError(w, http.StatusNotFound, "Message not found")
		return
	}

	newPinState := !dm.IsPinned
	if newPinState {
		err = h.MessagingRepo.PinDirectMessage(msgID)
	} else {
		err = h.MessagingRepo.UnpinDirectMessage(msgID)
	}
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to toggle pin")
		return
	}

	payload := websocket.PinPayload{
		MessageID:       msgID,
		IsPinned:        newPinState,
		IsDirectMessage: true,
	}

	payloadBytes, _ := json.Marshal(payload)
	wsMsg := &websocket.WSMessage{
		Type:    websocket.MessageTypePin,
		Payload: payloadBytes,
	}

	h.WSHub.SendToUser(dm.SenderID, wsMsg)
	h.WSHub.SendToUser(dm.ReceiverID, wsMsg)

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) PinRoomMessage(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	roomIDStr := chi.URLParam(r, "roomId")
	roomID, err := strconv.ParseInt(roomIDStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid room ID")
		return
	}

	msgIDStr := chi.URLParam(r, "messageId")
	msgID, err := strconv.ParseInt(msgIDStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid message ID")
		return
	}

	isMod, err := h.MessagingRepo.IsRoomModerator(roomID, userID)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to check permissions")
		return
	}
	if !isMod {
		httpx.JSONError(w, http.StatusForbidden, "Only moderators can pin messages")
		return
	}

	if err := h.MessagingRepo.PinRoomMessage(msgID); err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to pin message")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetPinnedRoomMessages(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	roomIDStr := chi.URLParam(r, "roomId")
	roomID, err := strconv.ParseInt(roomIDStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid room ID")
		return
	}

	isMember, err := h.MessagingRepo.IsRoomMember(roomID, userID)
	if err != nil || !isMember {
		httpx.JSONError(w, http.StatusForbidden, "Not a member of this room")
		return
	}

	messages, err := h.MessagingRepo.GetPinnedRoomMessages(roomID)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to get pinned messages")
		return
	}

	httpx.JSON(w, http.StatusOK, messages)
}
