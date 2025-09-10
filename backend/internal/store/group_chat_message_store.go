package store

import (
	"database/sql"
	"fmt"

	"github.com/tajjjjr/social-network/backend/internal/models"
)

type groupChatMessageStore struct {
	db *sql.DB
}

func NewGroupChatMessageStore(db *sql.DB) GroupChatMessageStore {
	return &groupChatMessageStore{db: db}
}

func (s *groupChatMessageStore) CreateGroupChatMessage(message *models.GroupChatMessage) (*models.GroupChatMessage, error) {
	stmt, err := s.db.Prepare("INSERT INTO group_chat_messages (group_id, sender_id, content) VALUES (?, ?, ?)")
	if err != nil {
		return nil, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(message.GroupID, message.SenderID, message.Content)
	if err != nil {
		return nil, fmt.Errorf("error executing statement: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert ID: %w", err)
	}

	message.ID = fmt.Sprintf("%d", id)
	return message, nil
}

func (s *groupChatMessageStore) GetGroupChatMessages(groupPublicID string, limit, offset int) ([]*models.GroupChatMessage, error) {
	// First get the internal group ID from public_id
	var groupID int64
	err := s.db.QueryRow("SELECT id FROM Groups WHERE public_id = ? AND type = 'group'", groupPublicID).Scan(&groupID)
	if err != nil {
		return nil, fmt.Errorf("error finding group: %w", err)
	}

	rows, err := s.db.Query("SELECT id, group_id, sender_id, content, created_at FROM group_chat_messages WHERE group_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?", groupID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying group chat messages: %w", err)
	}
	defer rows.Close()

	var messages []*models.GroupChatMessage
	for rows.Next() {
		var message models.GroupChatMessage
		err := rows.Scan(&message.ID, &message.GroupID, &message.SenderID, &message.Content, &message.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning group chat message: %w", err)
		}
		messages = append(messages, &message)
	}

	return messages, nil
}
