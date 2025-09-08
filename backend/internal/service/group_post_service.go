package service

import (
	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/store"
)

type GroupPostServiceInterface interface {
	CreateGroupPost(post *models.GroupPost, imageData []byte, imageMimeType string) (int64, error)
	GetGroupPostByID(postID int64) (*models.GroupPost, error)
	GetGroupPosts(groupID string, userID int64, limit, offset int) ([]*models.GroupPost, error) // migrated groupID to string
	UpdateGroupPost(postID, userID int64, content string, imageData []byte, imageMimeType string) (*models.GroupPost, error)
	DeleteGroupPost(postID, userID int64) error
	CreateGroupPostComment(comment *models.GroupPostComment, imageData []byte, imageMimeType string) (int64, error)
	GetGroupPostComments(postID int64, userID int64) ([]*models.GroupPostComment, error)
	UpdateGroupPostComment(commentID, userID int64, content string, imageData []byte, imageMimeType string) (*models.GroupPostComment, error)
	DeleteGroupPostComment(commentID, userID int64) error
	IsGroupMember(groupID string, userID int64) (bool, error) // migrated groupID to string
}

type GroupPostService struct {
	groupPostStore   store.GroupPostStore
	groupMemberStore store.GroupMemberStore
}

func NewGroupPostService(groupPostStore store.GroupPostStore, groupMemberStore store.GroupMemberStore) GroupPostServiceInterface {
	return &GroupPostService{
		groupPostStore:   groupPostStore,
		groupMemberStore: groupMemberStore,
	}
}

func (s *GroupPostService) CreateGroupPost(post *models.GroupPost, imageData []byte, imageMimeType string) (int64, error) {
	// TODO: Handle image upload similar to regular posts
	if imageData != nil {
		// Save image and set post.Image path
	}
	
	createdPost, err := s.groupPostStore.CreateGroupPost(post)
	if err != nil {
		return 0, err
	}
	return createdPost.ID, nil
}

func (s *GroupPostService) GetGroupPostByID(postID int64) (*models.GroupPost, error) {
	return s.groupPostStore.GetGroupPostByID(postID)
}

func (s *GroupPostService) GetGroupPosts(groupID string, userID int64, limit, offset int) ([]*models.GroupPost, error) {
	return s.groupPostStore.GetGroupPosts(groupID, userID, limit, offset)
}

func (s *GroupPostService) UpdateGroupPost(postID, userID int64, content string, imageData []byte, imageMimeType string) (*models.GroupPost, error) {
	return s.groupPostStore.UpdateGroupPost(postID, userID, content, imageData, imageMimeType)
}

func (s *GroupPostService) DeleteGroupPost(postID, userID int64) error {
	return s.groupPostStore.DeleteGroupPost(postID, userID)
}

func (s *GroupPostService) CreateGroupPostComment(comment *models.GroupPostComment, imageData []byte, imageMimeType string) (int64, error) {
	// TODO: Handle image upload similar to regular comments
	if imageData != nil {
		// Save image and set comment.Image path
	}
	
	createdComment, err := s.groupPostStore.CreateGroupPostComment(comment)
	if err != nil {
		return 0, err
	}
	return createdComment.ID, nil
}

func (s *GroupPostService) GetGroupPostComments(postID int64, userID int64) ([]*models.GroupPostComment, error) {
	return s.groupPostStore.GetGroupPostComments(postID, userID)
}

func (s *GroupPostService) UpdateGroupPostComment(commentID, userID int64, content string, imageData []byte, imageMimeType string) (*models.GroupPostComment, error) {
	return s.groupPostStore.UpdateGroupPostComment(commentID, userID, content, imageData, imageMimeType)
}

func (s *GroupPostService) DeleteGroupPostComment(commentID, userID int64) error {
	return s.groupPostStore.DeleteGroupPostComment(commentID, userID)
}

func (s *GroupPostService) IsGroupMember(groupID string, userID int64) (bool, error) {
	return s.groupMemberStore.IsGroupMember(groupID, userID)
}