package models

import "time"

type Group struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatorID   int64     `json:"creator_id"`
	Privacy     string    `json:"privacy"` // e.g., "public", "private"
	Avatar      string    `json:"avatar,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}
