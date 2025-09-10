package models

import "time"

// GroupRecord represents a unified record in the Groups table
type GroupRecord struct {
	ID        string    `json:"id" db:"id"`
	Type      string    `json:"type" db:"type"`
	GroupID   *string   `json:"group_id" db:"group_id"`
	UserID    *int64    `json:"user_id" db:"user_id"`
	Title     *string   `json:"title" db:"title"`
	Content   *string   `json:"content" db:"content"`
	Role      string    `json:"role" db:"role"`
	Status    string    `json:"status" db:"status"`
	Privacy   string    `json:"privacy" db:"privacy"`
	Image     *string   `json:"image" db:"image"`
	Data      *string   `json:"data" db:"data"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Group struct {
	ID          string    `json:"id"`
	PublicID    string    `json:"public_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatorID   int64     `json:"creator_id"`
	Privacy     string    `json:"privacy"`
	Avatar      string    `json:"avatar,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// GroupView represents a group for API responses
type GroupView struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Privacy     string    `json:"privacy"`
	Image       string    `json:"image"`
	CreatorID   int64     `json:"creator_id"`
	MemberCount int       `json:"member_count"`
	CreatedAt   time.Time `json:"created_at"`
}

// MemberView represents a group member for API responses
type MemberView struct {
	UserID    int64  `json:"user_id"`
	Role      string `json:"role"`
	Status    string `json:"status"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Avatar    string `json:"avatar"`
}

// PostView represents a group post for API responses
type PostView struct {
	ID        string    `json:"id"`
	GroupID   string    `json:"group_id"`
	UserID    int64     `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Image     string    `json:"image"`
	Author    *User     `json:"author"`
	CreatedAt time.Time `json:"created_at"`
}


