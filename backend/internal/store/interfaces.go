package store

import "github.com/tajjjjr/social-network/backend/internal/models"

// PostStoreInterface defines the interface for post-related database operations.
type PostStoreInterface interface {
	CreatePost(post *models.Post) (int64, error)
	CreateComment(comment *models.Comment) (int64, error)
	GetPostByID(id int64) (*models.Post, error)
	GetPosts(userID int64) ([]*models.Post, error)
	GetPostsPaginated(userID int64, limit, offset int) ([]*models.Post, error)
	GetPostsCount(userID int64) (int, error)
	UpdatePost(postID int64, content, imagePath string) (*models.Post, error)
	GetCommentsByPostID(postID, userID int64) ([]*models.Comment, error)
	DeletePost(postID int64) error
	AddPostViewers(postID int64, viewerIDs []int64) error
	SearchUsers(query string, currentUserID int64) ([]*models.User, error)

	UpdateComment(commentID int64, content, imagePath string) (*models.Comment, error)
	DeleteComment(commentID int64) error
	GetCommentByID(commentID int64) (*models.Comment, error)
}

type GroupStore interface {
	CreateGroup(group *models.Group) (*models.Group, error)
	GetGroupByID(groupID int64) (*models.Group, error)
	SearchPublicGroups(query string) ([]*models.Group, error)
	GetAllPublicGroups() ([]*models.Group, error)
	GetUserGroups(userID int64) ([]*models.Group, error)
}

type GroupRequestStore interface {
	CreateGroupRequest(request *models.GroupRequest) (*models.GroupRequest, error)
	GetGroupRequestByID(requestID int64) (*models.GroupRequest, error)
	UpdateGroupRequestStatus(requestID int64, status string) error
	AddUserToGroup(groupID, userID int64) error
}

type GroupChatMessageStore interface {
	CreateGroupChatMessage(message *models.GroupChatMessage) (*models.GroupChatMessage, error)
	GetGroupChatMessages(groupID int64, limit, offset int) ([]*models.GroupChatMessage, error)
}

type GroupMemberStoreInterface interface {
	IsGroupMember(groupID, userID int64) (bool, error)
	AddGroupMember(groupID, userID int64, role string) (*models.GroupMember, error)
	RemoveGroupMember(groupID, userID int64) error
}
