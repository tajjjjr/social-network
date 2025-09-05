package store

import (
	"database/sql"
	"time"
)

type GroupEvent struct {
	ID          int64     `json:"id"`
	GroupID     int64     `json:"group_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EventTime   time.Time `json:"event_time"`
	CreatedBy   int64     `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

type GroupEventStore interface {
	GetGroupEvents(groupID int64) ([]*GroupEvent, error)
}

type groupEventStore struct {
	db *sql.DB
}

func NewGroupEventStore(db *sql.DB) GroupEventStore {
	return &groupEventStore{db: db}
}

func (s *groupEventStore) GetGroupEvents(groupID int64) ([]*GroupEvent, error) {
	rows, err := s.db.Query(`
		SELECT id, group_id, title, description, event_time, created_by, created_at
		FROM Group_Events
		WHERE group_id = ?
		ORDER BY event_time ASC
	`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*GroupEvent
	for rows.Next() {
		var event GroupEvent
		err := rows.Scan(&event.ID, &event.GroupID, &event.Title, &event.Description, 
			&event.EventTime, &event.CreatedBy, &event.CreatedAt)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	return events, nil
}