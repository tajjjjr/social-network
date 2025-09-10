package utils

import "context"

type UserID string
type GroupID string

type Response struct {
	Message string `json:"message,omitempty"`
}

type PaginationMeta struct {
	CurrentPage int  `json:"currentPage"`
	TotalPages  int  `json:"totalPages"`
	TotalPosts  int  `json:"totalPosts"`
	HasMore     bool `json:"hasMore"`
	Limit       int  `json:"limit"`
}

type PostsResponse struct {
	Posts      interface{}    `json:"posts"`
	Pagination PaginationMeta `json:"pagination"`
}

var User_id UserID = "userID"
var Group_id GroupID = "groupID"

// SetUserContext is a helper function for tests to set user context
func SetUserContext(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, User_id, userID)
}

func SetGroupContext(ctx context.Context, groupID string) context.Context {
	return context.WithValue(ctx, Group_id, groupID)
}
