package service

import (
	"errors"
	"testing"

	"github.com/tajjjjr/social-network/backend/internal/models"
)

type MockGroupPostStore struct {
	CreateGroupPostFunc        func(post *models.GroupPost) (*models.GroupPost, error)
	GetGroupPostsFunc          func(groupID string, userID int64, limit, offset int) ([]*models.GroupPost, error) // migrated groupID to string
	UpdateGroupPostFunc        func(postID, userID int64, content string, imageData []byte, imageMimeType string) (*models.GroupPost, error)
	DeleteGroupPostFunc        func(postID, userID int64) error
	CreateGroupPostCommentFunc func(comment *models.GroupPostComment) (*models.GroupPostComment, error)
	GetGroupPostCommentsFunc   func(postID int64, userID int64) ([]*models.GroupPostComment, error)
}

func (m *MockGroupPostStore) CreateGroupPost(post *models.GroupPost) (*models.GroupPost, error) {
	if m.CreateGroupPostFunc != nil {
		return m.CreateGroupPostFunc(post)
	}
	post.ID = 1
	return post, nil
}

func (m *MockGroupPostStore) GetGroupPostByID(postID int64) (*models.GroupPost, error) {
	return &models.GroupPost{ID: postID}, nil
}

func (m *MockGroupPostStore) GetGroupPosts(groupID string, userID int64, limit, offset int) ([]*models.GroupPost, error) {
	if m.GetGroupPostsFunc != nil {
		return m.GetGroupPostsFunc(groupID, userID, limit, offset)
	}
	return []*models.GroupPost{}, nil
}

func (m *MockGroupPostStore) UpdateGroupPost(postID, userID int64, content string, imageData []byte, imageMimeType string) (*models.GroupPost, error) {
	if m.UpdateGroupPostFunc != nil {
		return m.UpdateGroupPostFunc(postID, userID, content, imageData, imageMimeType)
	}
	return &models.GroupPost{ID: postID, Content: content}, nil
}

func (m *MockGroupPostStore) DeleteGroupPost(postID, userID int64) error {
	if m.DeleteGroupPostFunc != nil {
		return m.DeleteGroupPostFunc(postID, userID)
	}
	return nil
}

func (m *MockGroupPostStore) CreateGroupPostComment(comment *models.GroupPostComment) (*models.GroupPostComment, error) {
	if m.CreateGroupPostCommentFunc != nil {
		return m.CreateGroupPostCommentFunc(comment)
	}
	comment.ID = 1
	return comment, nil
}

func (m *MockGroupPostStore) GetGroupPostComments(postID int64, userID int64) ([]*models.GroupPostComment, error) {
	if m.GetGroupPostCommentsFunc != nil {
		return m.GetGroupPostCommentsFunc(postID, userID)
	}
	return []*models.GroupPostComment{}, nil
}

func (m *MockGroupPostStore) UpdateGroupPostComment(commentID, userID int64, content string, imageData []byte, imageMimeType string) (*models.GroupPostComment, error) {
	return &models.GroupPostComment{ID: commentID, Content: content}, nil
}

func (m *MockGroupPostStore) DeleteGroupPostComment(commentID, userID int64) error {
	return nil
}

func (m *MockGroupPostStore) CanUserDeleteGroupContent(groupID string, userID int64) (bool, error) {
	return true, nil
}

func TestCreateGroupPost(t *testing.T) {
	t.Run("Successful post creation", func(t *testing.T) {
		mockStore := &MockGroupPostStore{
			CreateGroupPostFunc: func(post *models.GroupPost) (*models.GroupPost, error) {
				post.ID = 1
				return post, nil
			},
		}

		service := NewGroupPostService(mockStore, nil)

		post := &models.GroupPost{
			GroupID: "1", // migrated to string
			UserID:  1,
			Content: "Test content",
		}

		id, err := service.CreateGroupPost(post, nil, "")
		if err != nil {
			t.Fatalf("CreateGroupPost failed: %v", err)
		}

		if id != 1 {
			t.Errorf("Expected post ID 1, got %d", id)
		}
	})

	t.Run("Store error", func(t *testing.T) {
		mockStore := &MockGroupPostStore{
			CreateGroupPostFunc: func(post *models.GroupPost) (*models.GroupPost, error) {
				return nil, errors.New("store error")
			},
		}

		service := NewGroupPostService(mockStore, nil)

		post := &models.GroupPost{
			GroupID: "1", // migrated to string
			UserID:  1,
			Content: "Test content",
		}

		_, err := service.CreateGroupPost(post, nil, "")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

func TestGetGroupPosts(t *testing.T) {
	t.Run("Successful get posts", func(t *testing.T) {
		mockStore := &MockGroupPostStore{
			GetGroupPostsFunc: func(groupID string, userID int64, limit, offset int) ([]*models.GroupPost, error) {
				return []*models.GroupPost{
					{ID: 1, GroupID: groupID, UserID: userID, Content: "Post 1"},
					{ID: 2, GroupID: groupID, UserID: userID, Content: "Post 2"},
				}, nil
			},
		}

		service := NewGroupPostService(mockStore, nil)

		posts, err := service.GetGroupPosts("1", 1, 10, 0)
		if err != nil {
			t.Fatalf("GetGroupPosts failed: %v", err)
		}

		if len(posts) != 2 {
			t.Errorf("Expected 2 posts, got %d", len(posts))
		}
	})
}

func TestCreateGroupPostComment(t *testing.T) {
	t.Run("Successful comment creation", func(t *testing.T) {
		mockStore := &MockGroupPostStore{
			CreateGroupPostCommentFunc: func(comment *models.GroupPostComment) (*models.GroupPostComment, error) {
				comment.ID = 1
				return comment, nil
			},
		}

		service := NewGroupPostService(mockStore, nil)

		comment := &models.GroupPostComment{
			GroupPostID: "1", // migrated to string
			UserID:      1,
			Content:     "Test comment",
		}

		id, err := service.CreateGroupPostComment(comment, nil, "")
		if err != nil {
			t.Fatalf("CreateGroupPostComment failed: %v", err)
		}

		if id != 1 {
			t.Errorf("Expected comment ID 1, got %d", id)
		}
	})
}