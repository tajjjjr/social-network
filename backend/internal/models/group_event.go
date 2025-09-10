package models

import "time"

type GroupEvent struct {
	ID          int64     `json:"-"`
	PublicID    string    `json:"public_id"`
	GroupID     int64     `json:"-"`
	GroupPubID  string    `json:"group_public_id"`
	CreatorID   int64     `json:"creator_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EventDate   time.Time `json:"event_date"`
	CreatedAt   time.Time `json:"created_at"`
}