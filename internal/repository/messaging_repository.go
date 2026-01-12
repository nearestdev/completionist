package repository

import (
	"database/sql"
	"fmt"

	"github.com/GATEOPENERZ/completionist-api/internal/models"
	"github.com/jmoiron/sqlx"
)

type MessagingRepository struct {
	DB *sqlx.DB
}

func NewMessagingRepository(db *sqlx.DB) *MessagingRepository {
	return &MessagingRepository{DB: db}
}

func (r *MessagingRepository) SendDirectMessage(msg *models.DirectMessage) (*models.DirectMessage, error) {
	var result models.DirectMessage
	err := r.DB.QueryRowx(`
		INSERT INTO direct_messages (sender_id, receiver_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, sender_id, receiver_id, content, is_read, created_at
	`, msg.SenderID, msg.ReceiverID, msg.Content).StructScan(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *MessagingRepository) GetConversation(userID, otherUserID int64, limit int, beforeID int64) ([]models.DirectMessageResponse, error) {
	messages := []models.DirectMessageResponse{}
	
	query := `
		SELECT 
			dm.id,
			dm.sender_id,
			s.username AS sender_username,
			dm.receiver_id,
			rec.username AS receiver_username,
			dm.content,
			dm.is_read,
			dm.created_at
		FROM direct_messages dm
		JOIN users s ON dm.sender_id = s.id
		JOIN users rec ON dm.receiver_id = rec.id
		WHERE ((dm.sender_id = $1 AND dm.receiver_id = $2)
		   OR (dm.sender_id = $2 AND dm.receiver_id = $1))
	`
	args := []interface{}{userID, otherUserID}
	argIdx := 3

	if beforeID > 0 {
		query += fmt.Sprintf(" AND dm.id < $%d", argIdx)
		args = append(args, beforeID)
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY dm.created_at DESC LIMIT $%d", argIdx)
	args = append(args, limit)

	err := r.DB.Select(&messages, query, args...)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *MessagingRepository) GetConversations(userID int64) ([]models.ConversationSummary, error) {
	conversations := []models.ConversationSummary{}
	err := r.DB.Select(&conversations, `
		WITH latest_messages AS (
			SELECT DISTINCT ON (
				CASE 
					WHEN sender_id = $1 THEN receiver_id
					ELSE sender_id
				END
			)
				CASE 
					WHEN sender_id = $1 THEN receiver_id
					ELSE sender_id
				END AS other_user_id,
				content AS last_message,
				created_at AS last_message_time,
				sender_id = $1 AS is_sent_by_current_user
			FROM direct_messages
			WHERE sender_id = $1 OR receiver_id = $1
			ORDER BY 
				CASE 
					WHEN sender_id = $1 THEN receiver_id
					ELSE sender_id
				END,
				created_at DESC
		),
		unread_counts AS (
			SELECT sender_id, COUNT(*) as unread_count
			FROM direct_messages
			WHERE receiver_id = $1 AND is_read = false
			GROUP BY sender_id
		)
		SELECT 
			lm.other_user_id,
			u.username AS other_username,
			lm.last_message,
			lm.last_message_time,
			lm.is_sent_by_current_user,
			COALESCE(uc.unread_count, 0) AS unread_count
		FROM latest_messages lm
		JOIN users u ON lm.other_user_id = u.id
		LEFT JOIN unread_counts uc ON lm.other_user_id = uc.sender_id
		ORDER BY lm.last_message_time DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	return conversations, nil
}

func (r *MessagingRepository) MarkAsRead(receiverID, senderID int64) error {
	_, err := r.DB.Exec(`
		UPDATE direct_messages
		SET is_read = true
		WHERE receiver_id = $1 AND sender_id = $2 AND is_read = false
	`, receiverID, senderID)
	return err
}

func (r *MessagingRepository) CreateRoom(room *models.Room) (*models.Room, error) {
	var result models.Room
	err := r.DB.QueryRowx(`
		INSERT INTO rooms (name, description, room_type, creator_id, is_public, max_members)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, description, room_type, creator_id, is_public, max_members, created_at, updated_at
	`, room.Name, room.Description, room.RoomType, room.CreatorID, room.IsPublic, room.MaxMembers).StructScan(&result)
	if err != nil {
		return nil, err
	}
	_, err = r.DB.Exec(`
		INSERT INTO room_members (room_id, user_id, is_moderator)
		VALUES ($1, $2, true)
	`, result.ID, result.CreatorID)
	if err != nil {
		return nil, err
	}
	if result.RoomType == models.RoomTypeMusic || result.RoomType == models.RoomTypeWatch {
		_, err = r.DB.Exec(`
			INSERT INTO room_state (room_id)
			VALUES ($1)
		`, result.ID)
		if err != nil {
			return nil, err
		}
	}
	return &result, nil
}

func (r *MessagingRepository) GetRoom(roomID int64) (*models.Room, error) {
	var room models.Room
	err := r.DB.Get(&room, `
		SELECT id, name, description, room_type, creator_id, is_public, max_members, created_at, updated_at
		FROM rooms
		WHERE id = $1
	`, roomID)
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *MessagingRepository) ListRooms(roomType *models.RoomType, limit int) ([]models.RoomResponse, error) {
	query := `
		SELECT 
			r.id,
			r.name,
			r.description,
			r.room_type,
			r.creator_id,
			u.username AS creator_name,
			r.is_public,
			r.max_members,
			COUNT(rm.id) AS member_count,
			r.created_at,
			r.updated_at
		FROM rooms r
		JOIN users u ON r.creator_id = u.id
		LEFT JOIN room_members rm ON r.id = rm.room_id
		WHERE r.is_public = true
	`
	args := []interface{}{}
	if roomType != nil {
		query += " AND r.room_type = $1"
		args = append(args, *roomType)
		query += " GROUP BY r.id, u.username ORDER BY r.created_at DESC LIMIT $2"
		args = append(args, limit)
	} else {
		query += " GROUP BY r.id, u.username ORDER BY r.created_at DESC LIMIT $1"
		args = append(args, limit)
	}

	rooms := []models.RoomResponse{}
	err := r.DB.Select(&rooms, query, args...)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func (r *MessagingRepository) GetRoomWithDetails(roomID, userID int64) (*models.RoomResponse, error) {
	var room models.RoomResponse
	err := r.DB.Get(&room, `
		SELECT 
			r.id,
			r.name,
			r.description,
			r.room_type,
			r.creator_id,
			u.username AS creator_name,
			r.is_public,
			r.max_members,
			COUNT(rm.id) AS member_count,
			EXISTS(SELECT 1 FROM room_members WHERE room_id = r.id AND user_id = $2) AS is_member,
			EXISTS(SELECT 1 FROM room_members WHERE room_id = r.id AND user_id = $2 AND is_moderator = true) AS is_moderator,
			r.created_at,
			r.updated_at
		FROM rooms r
		JOIN users u ON r.creator_id = u.id
		LEFT JOIN room_members rm ON r.id = rm.room_id
		WHERE r.id = $1
		GROUP BY r.id, u.username
	`, roomID, userID)
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *MessagingRepository) JoinRoom(roomID, userID int64) error {
	var maxMembers, currentCount int
	err := r.DB.Get(&maxMembers, "SELECT max_members FROM rooms WHERE id = $1", roomID)
	if err != nil {
		return err
	}
	err = r.DB.Get(&currentCount, "SELECT COUNT(*) FROM room_members WHERE room_id = $1", roomID)
	if err != nil {
		return err
	}
	if currentCount >= maxMembers {
		return fmt.Errorf("room is full")
	}
	_, err = r.DB.Exec(`
		INSERT INTO room_members (room_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT (room_id, user_id) DO NOTHING
	`, roomID, userID)
	return err
}

func (r *MessagingRepository) LeaveRoom(roomID, userID int64) error {
	_, err := r.DB.Exec(`
		DELETE FROM room_members
		WHERE room_id = $1 AND user_id = $2
	`, roomID, userID)
	return err
}

func (r *MessagingRepository) DeleteRoom(roomID int64) error {
	_, err := r.DB.Exec("DELETE FROM rooms WHERE id = $1", roomID)
	return err
}

func (r *MessagingRepository) IsRoomMember(roomID, userID int64) (bool, error) {
	var exists bool
	err := r.DB.Get(&exists, `
		SELECT EXISTS(SELECT 1 FROM room_members WHERE room_id = $1 AND user_id = $2)
	`, roomID, userID)
	return exists, err
}

func (r *MessagingRepository) IsRoomModerator(roomID, userID int64) (bool, error) {
	var isMod bool
	err := r.DB.Get(&isMod, `
		SELECT is_moderator FROM room_members
		WHERE room_id = $1 AND user_id = $2
	`, roomID, userID)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return isMod, err
}

func (r *MessagingRepository) GetRoomMembers(roomID int64) ([]models.RoomMemberResponse, error) {
	members := []models.RoomMemberResponse{}
	err := r.DB.Select(&members, `
		SELECT 
			rm.user_id,
			u.username,
			rm.is_moderator,
			rm.joined_at
		FROM room_members rm
		JOIN users u ON rm.user_id = u.id
		WHERE rm.room_id = $1
		ORDER BY rm.joined_at ASC
	`, roomID)
	return members, err
}

func (r *MessagingRepository) SendRoomMessage(msg *models.RoomMessage) (*models.RoomMessage, error) {
	var result models.RoomMessage
	err := r.DB.QueryRowx(`
		INSERT INTO room_messages (room_id, user_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, room_id, user_id, content, created_at
	`, msg.RoomID, msg.UserID, msg.Content).StructScan(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *MessagingRepository) GetUserRooms(userID int64) ([]int64, error) {
	var roomIDs []int64
	err := r.DB.Select(&roomIDs, "SELECT room_id FROM room_members WHERE user_id = $1", userID)
	return roomIDs, err
}

func (r *MessagingRepository) GetRoomMessages(roomID int64, limit int) ([]models.RoomMessageResponse, error) {
	messages := []models.RoomMessageResponse{}
	err := r.DB.Select(&messages, `
		SELECT 
			rm.id,
			rm.room_id,
			rm.user_id,
			u.username,
			rm.content,
			rm.created_at
		FROM room_messages rm
		JOIN users u ON rm.user_id = u.id
		WHERE rm.room_id = $1
		ORDER BY rm.created_at DESC
		LIMIT $2
	`, roomID, limit)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *MessagingRepository) GetRoomState(roomID int64) (*models.RoomState, error) {
	var state models.RoomState
	err := r.DB.Get(&state, `
		SELECT room_id, current_media_url, current_media_title, current_position_ms, is_playing, updated_by, updated_at
		FROM room_state
		WHERE room_id = $1
	`, roomID)
	if err != nil {
		return nil, err
	}
	return &state, nil
}

func (r *MessagingRepository) UpdateRoomState(roomID, userID int64, req *models.UpdateRoomStateRequest) (*models.RoomState, error) {
	query := "UPDATE room_state SET updated_by = $1, updated_at = NOW()"
	args := []interface{}{userID}
	argIndex := 2

	if req.CurrentMediaURL != nil {
		query += fmt.Sprintf(", current_media_url = $%d", argIndex)
		args = append(args, *req.CurrentMediaURL)
		argIndex++
	}
	if req.CurrentMediaTitle != nil {
		query += fmt.Sprintf(", current_media_title = $%d", argIndex)
		args = append(args, *req.CurrentMediaTitle)
		argIndex++
	}
	if req.CurrentPositionMs != nil {
		query += fmt.Sprintf(", current_position_ms = $%d", argIndex)
		args = append(args, *req.CurrentPositionMs)
		argIndex++
	}
	if req.IsPlaying != nil {
		query += fmt.Sprintf(", is_playing = $%d", argIndex)
		args = append(args, *req.IsPlaying)
		argIndex++
	}

	query += fmt.Sprintf(" WHERE room_id = $%d RETURNING room_id, current_media_url, current_media_title, current_position_ms, is_playing, updated_by, updated_at", argIndex)
	args = append(args, roomID)

	var state models.RoomState
	err := r.DB.QueryRowx(query, args...).StructScan(&state)
	if err != nil {
		return nil, err
	}
	return &state, nil
}
