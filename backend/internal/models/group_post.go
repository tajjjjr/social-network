package models

import "time"

type GroupPost struct {
	ID           int64     `json:"-"`
	PublicID     string    `json:"id"`
	GroupID      string    `json:"group_id"`
	UserID       int64     `json:"user_id"`
	Content      string    `json:"content"`
	Image        string    `json:"image,omitempty"`
	LikeCount    int       `json:"like_count"`
	DislikeCount int       `json:"dislike_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Additional fields for API responses
	Author       *User  `json:"author,omitempty"`
	UserReaction string `json:"user_reaction,omitempty"`
	IsEdited     bool   `json:"is_edited,omitempty"`
}

type GroupPostComment struct {
	ID              int64     `json:"-"`
	PublicID        string    `json:"id"`
	GroupPostID     string    `json:"group_post_id"`
	UserID          int64     `json:"user_id"`
	ParentCommentID *string   `json:"parent_comment_id,omitempty"`
	Content         string    `json:"content"`
	Image           string    `json:"image,omitempty"`
	LikeCount       int       `json:"like_count"`
	DislikeCount    int       `json:"dislike_count"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	// Additional fields for API responses
	Author       *User  `json:"author,omitempty"`
	UserReaction string `json:"user_reaction,omitempty"`
}

type GroupPermission struct {
	ID             int64     `json:"id"`
	GroupID        int64     `json:"group_id"`
	UserID         int64     `json:"user_id"`
	PermissionType string    `json:"permission_type"` // 'admin', 'moderator'
	GrantedBy      int64     `json:"granted_by"`
	GrantedAt      time.Time `json:"granted_at"`
}
