package models

import "time"

type GroupChatMessage struct {
	ID        int64     `json:"-" db:"id"`
	PublicID  string    `json:"id" db:"public_id"`
	GroupID   string    `json:"group_id" db:"group_id"`
	SenderID  int64     `json:"sender_id" db:"user_id"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
