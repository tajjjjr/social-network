package service

import (
	"github.com/tajjjjr/social-network/backend/internal/models"
)

// AuthServiceInterface defines the interface for the auth service.
type AuthServiceInterface interface {
	AuthenticateUser(email, password string) (*models.User, string, error)
	DeleteSession(sessionID string) (int, error)
	GetUserIDBySession(sessionID string) (int, error)
	CreateUser(user *models.User) (*models.User, error)
	ValidateEmail(email string) (bool, error)
	UserExists(email string) (bool, error)
	UserNewEditEmailExist(email string, userid int64) (bool, error)
	EditUserProfile(user *models.User, userid int64) error
}

// PostServiceInterface defines the interface for the post service.
type PostServiceInterface interface {
	CreatePost(post *models.Post, imageData []byte, imageMimeType string) (int64, error)
	CreatePostWithViewers(post *models.Post, imageData []byte, imageMimeType string, viewerIDs []int64) (int64, error)
	GetPostByID(id int64) (*models.Post, error)
	GetPosts(userID int64) ([]*models.Post, error)
	GetPostsPaginated(userID int64, limit, offset int) ([]*models.Post, error)
	GetPostsCount(userID int64) (int, error)
	UpdatePost(postID, userID int64, content string, imageData []byte, imageMimeType string) (*models.Post, error)
	CreateComment(comment *models.Comment, imageData []byte, imageMimeType string) (int64, error)
	GetCommentsByPostID(postID, userID int64) ([]*models.Comment, error)
	DeletePost(postID, userID int64) error
	SearchUsers(query string, currentUserID int64) ([]*models.User, error)
	UpdateComment(commentID, userID int64, content string, imageData []byte, imageMimeType string) (*models.Comment, error)
	DeleteComment(commentID, userID int64) error
	GetCommentByID(commentID int64) (*models.Comment, error)
}

type FollowServiceInterface interface {
	IsAccountPublic(followeeID int64) (bool, error)
	CreateFollowForPublicAccount(followerid, followeeid int64) (int64, error)
	CreateFollowForPrivateAccount(followrid, followeeid int64) (int64, error)
	GetUserInfo(userID int64) (string, string, error)
	AddtoNotification(follower_id int64, message string) error
}

type UnfollowServiceInterface interface {
	GetFollowConnectionID(followerID, followeeID int64) (int64, error)
	DeleteFollowConnection(followConnectionID int64) error
}

type FollowRequestServiceInterface interface {
	AcceptedFollowConnection(followConnectionID int64) error
	RejectedFollowConnection(followConnectionID int64) error
	CancelFollowRequest(followConnectionID int64) error
	RetrieveUserName(userID int64) (string, string, error)
	GetRequestInfo(requestID int64) (int64, int64, error)
	GetRequestIDByUsers(followerID, followeeID int64) (int64, error)
	AddtoNotification(follower_id int64, message string) error
	GetPendingFollowRequest(userid int64) (models.FollowRequestUserResponse, error)
}

type ProfileServiceInterface interface {
	GetUserOwnProfile(userid int64) (models.ProfileDetails, error)
	GetUserProfile(userid, LoggedInUser int64) (models.ProfileDetails, error)
	GetUserPosts(userid int64) ([]models.Post, error)
	GetFollowersList(userid int64) (models.FollowListResponse, error)
	GetFolloweesList(userid int64) (models.FollowListResponse, error)
	GetUserPhotos(userId int64) ([]models.Photo, error)
}

type GroupService interface {
	CreateGroup(group *models.Group) (*models.Group, error)
	GetGroupByID(groupID int64) (*models.Group, error)
	SearchPublicGroups(query string) ([]*models.Group, error)
	GetAllPublicGroups() ([]*models.Group, error)
	GetUserGroups(userID int64) ([]*models.Group, error)
}

type GroupRequestService interface {
	SendJoinRequest(groupID, userID int64) (*models.GroupRequest, error)
	ApproveJoinRequest(requestID int64, approverID int64) error
	RejectJoinRequest(requestID int64, rejecterID int64) error
}

type GroupChatMessageService interface {
	SendGroupChatMessage(groupID, senderID int64, content string) (*models.GroupChatMessage, error)
	GetGroupChatMessages(groupID int64, userID int64, limit, offset int) ([]*models.GroupChatMessage, error)
}
