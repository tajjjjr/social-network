package models

import "time"

type GroupRequest struct {
	ID        int64     `json:"-" db:"id"`
	PublicID  string    `json:"id" db:"public_id"`
	GroupID   string    `json:"group_id" db:"group_id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Status    string    `json:"status" db:"status"` // e.g., "pending", "approved", "rejected"
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
