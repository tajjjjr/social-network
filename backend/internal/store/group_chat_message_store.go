package store

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/tajjjjr/social-network/backend/internal/models"
)

type groupChatMessageStore struct {
	db *sql.DB
}

func NewGroupChatMessageStore(db *sql.DB) GroupChatMessageStore {
	return &groupChatMessageStore{db: db}
}

func (s *groupChatMessageStore) CreateGroupChatMessage(message *models.GroupChatMessage) (*models.GroupChatMessage, error) {
	// Get internal group ID from public_id
	var groupID int64
	err := s.db.QueryRow("SELECT id FROM Groups WHERE public_id = ? AND type = 'group'", message.GroupID).Scan(&groupID)
	if err != nil {
		return nil, fmt.Errorf("error finding group: %w", err)
	}

	result, err := s.db.Exec(`
		INSERT INTO Groups (type, group_id, user_id, content) 
		VALUES ('message', ?, ?, ?)`,
		groupID, message.SenderID, message.Content)
	if err != nil {
		return nil, fmt.Errorf("error executing statement: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert ID: %w", err)
	}

	message.ID = id
	return message, nil
}

func (s *groupChatMessageStore) GetGroupChatMessages(groupPublicID string, limit, offset int) ([]*models.GroupChatMessage, error) {
	// Get the internal group ID from public_id
	var groupID int64
	err := s.db.QueryRow("SELECT id FROM Groups WHERE public_id = ? AND type = 'group'", groupPublicID).Scan(&groupID)
	if err != nil {
		return nil, fmt.Errorf("error finding group: %w", err)
	}

	rows, err := s.db.Query(`
		SELECT id, group_id, user_id, content, created_at 
		FROM Groups 
		WHERE type = 'message' AND group_id = ? 
		ORDER BY created_at DESC LIMIT ? OFFSET ?`, groupID, limit, offset)
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
		// Set the public group ID for the response
		message.GroupID = groupPublicID
		messages = append(messages, &message)
	}

	return messages, nil
}
