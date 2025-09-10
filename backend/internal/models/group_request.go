package models

import "time"

type GroupRequest struct {
	ID        string     `json:"id"`
	GroupID   string     `json:"group_id"`
	UserID    int64      `json:"user_id"`
	Status    string     `json:"status"` // e.g., "pending", "approved", "rejected"
	CreatedAt time.Time  `json:"created_at"`
}
