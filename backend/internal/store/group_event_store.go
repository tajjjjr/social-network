package store

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/tajjjjr/social-network/backend/internal/models"
	"time"
)

type GroupEventStore interface {
	GetGroupEvents(groupPublicID string) ([]*models.GroupEvent, error)
	CreateGroupEvent(event *models.GroupEvent) (*models.GroupEvent, error)
}

type groupEventStore struct {
	db *sql.DB
}

func NewGroupEventStore(db *sql.DB) GroupEventStore {
	return &groupEventStore{db: db}
}

func (s *groupEventStore) getGroupIDByPublicID(publicID string) (int64, error) {
	var id int64
	err := s.db.QueryRow("SELECT id FROM Groups WHERE public_id = ?", publicID).Scan(&id)
	return id, err
}

func (s *groupEventStore) CreateGroupEvent(event *models.GroupEvent) (*models.GroupEvent, error) {
	event.PublicID = uuid.New().String()
	groupID, err := s.getGroupIDByPublicID(event.GroupPubID)
	if err != nil {
		return nil, err
	}
	event.GroupID = groupID

	result, err := s.db.Exec(`
		INSERT INTO Groups (public_id, type, group_id, user_id, title, content, data)
		VALUES (?, 'event', ?, ?, ?, ?, ?)`,
		event.PublicID, groupID, event.CreatorID, event.Title, event.Description, event.EventDate.Format("2006-01-02T15:04:05Z07:00"))
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	event.ID = id
	return event, nil
}

func (s *groupEventStore) GetGroupEvents(groupPublicID string) ([]*models.GroupEvent, error) {
	rows, err := s.db.Query(`
		SELECT e.id, e.public_id, g.public_id, e.user_id, e.title, e.content, e.data, e.created_at
		FROM Groups e
		JOIN Groups g ON e.group_id = g.id
		WHERE e.type = 'event' AND g.public_id = ?
		ORDER BY e.data ASC
	`, groupPublicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*models.GroupEvent
	for rows.Next() {
		var event models.GroupEvent
		var eventDateStr string
		var idStr string
		err := rows.Scan(&idStr, &event.PublicID, &event.GroupPubID, &event.CreatorID, &event.Title,
			&event.Description, &eventDateStr, &event.CreatedAt)
		if err != nil {
			return nil, err
		}
		// Parse the event date from the data field
		if eventDate, parseErr := time.Parse("2006-01-02T15:04:05Z07:00", eventDateStr); parseErr == nil {
			event.EventDate = eventDate
		}
		events = append(events, &event)
	}

	return events, nil
}
