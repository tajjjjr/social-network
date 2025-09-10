package models

import "time"

type GroupMember struct {
	ID        int64     `json:"id"`
	GroupID   string     `json:"group_id"`
	UserID    int64     `json:"user_id"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}
