package models

import "time"

type GroupMember struct {
	ID        int64     `json:"-" db:"id"`
	PublicID  string    `json:"id" db:"public_id"`
	GroupID   string    `json:"group_id" db:"group_id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Role      string    `json:"role" db:"role"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
