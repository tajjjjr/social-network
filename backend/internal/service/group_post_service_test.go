package service

import (
	"errors"
	"testing"

	"github.com/tajjjjr/social-network/backend/internal/models"
)

type MockGroupPostStore struct {
	CreateGroupPostFunc        func(post *models.GroupPost) (*models.GroupPost, error)
	GetGroupPostsFunc          func(groupID string, userID int64, limit, offset int) ([]*models.GroupPost, error)
	UpdateGroupPostFunc        func(postPublicID string, userID int64, content string, imageData []byte, imageMimeType string) (*models.GroupPost, error)
	DeleteGroupPostFunc        func(postPublicID string, userID int64) error
	CreateGroupPostCommentFunc func(comment *models.GroupPostComment) (*models.GroupPostComment, error)
	GetGroupPostCommentsFunc   func(postPublicID string, userID int64) ([]*models.GroupPostComment, error)
	IsGroupMemberFunc          func(groupID string, userID int64) (bool, error)
}

func (m *MockGroupPostStore) CreateGroupPost(post *models.GroupPost) (*models.GroupPost, error) {
	if m.CreateGroupPostFunc != nil {
		return m.CreateGroupPostFunc(post)
	}
	post.ID = 1
	return post, nil
}

func (m *MockGroupPostStore) GetGroupPostByID(postPublicID string) (*models.GroupPost, error) {
	return &models.GroupPost{PublicID: postPublicID}, nil
}

func (m *MockGroupPostStore) GetGroupPosts(groupID string, userID int64, limit, offset int) ([]*models.GroupPost, error) {
	if m.GetGroupPostsFunc != nil {
		return m.GetGroupPostsFunc(groupID, userID, limit, offset)
	}
	return []*models.GroupPost{}, nil
}

func (m *MockGroupPostStore) UpdateGroupPost(postPublicID string, userID int64, content string, imageData []byte, imageMimeType string) (*models.GroupPost, error) {
	if m.UpdateGroupPostFunc != nil {
		return m.UpdateGroupPostFunc(postPublicID, userID, content, imageData, imageMimeType)
	}
	return &models.GroupPost{PublicID: postPublicID, Content: content}, nil
}

func (m *MockGroupPostStore) DeleteGroupPost(postPublicID string, userID int64) error {
	if m.DeleteGroupPostFunc != nil {
		return m.DeleteGroupPostFunc(postPublicID, userID)
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

func (m *MockGroupPostStore) GetGroupPostComments(postPublicID string, userID int64) ([]*models.GroupPostComment, error) {
	if m.GetGroupPostCommentsFunc != nil {
		return m.GetGroupPostCommentsFunc(postPublicID, userID)
	}
	return []*models.GroupPostComment{}, nil
}

func (m *MockGroupPostStore) UpdateGroupPostComment(commentPublicID string, userID int64, content string, imageData []byte, imageMimeType string) (*models.GroupPostComment, error) {
	return &models.GroupPostComment{PublicID: commentPublicID, Content: content}, nil
}

func (m *MockGroupPostStore) DeleteGroupPostComment(commentPublicID string, userID int64) error {
	return nil
}

func (m *MockGroupPostStore) CanUserDeleteGroupContent(groupID string, userID int64) (bool, error) {
	return true, nil
}

func (m *MockGroupPostStore) IsGroupMember(groupID string, userID int64) (bool, error) {
	if m.IsGroupMemberFunc != nil {
		return m.IsGroupMemberFunc(groupID, userID)
	}
	return true, nil
}

type MockGroupStore struct {
	GetGroupByIDFunc func(publicID string) (*models.Group, error)
}

func (m *MockGroupStore) CreateGroup(group *models.Group) (*models.Group, error) {
	return group, nil
}

func (m *MockGroupStore) GetGroupByID(publicID string) (*models.Group, error) {
	if m.GetGroupByIDFunc != nil {
		return m.GetGroupByIDFunc(publicID)
	}
	return &models.Group{ID: "1", PublicID: publicID}, nil
}

func (m *MockGroupStore) SearchPublicGroups(query string) ([]*models.Group, error) {
	return []*models.Group{}, nil
}

func (m *MockGroupStore) GetAllPublicGroups() ([]*models.Group, error) {
	return []*models.Group{}, nil
}

func (m *MockGroupStore) JoinGroup(publicID string, userID int64) error {
	return nil
}

func (m *MockGroupStore) LeaveGroup(publicID string, userID int64) error {
	return nil
}

func (m *MockGroupStore) IsGroupMember(publicID string, userID int64) (bool, error) {
	return true, nil
}

func (m *MockGroupStore) GetUserGroups(userID int64) ([]*models.Group, error) {
	return []*models.Group{}, nil
}

func TestCreateGroupPost(t *testing.T) {
	t.Run("Successful post creation", func(t *testing.T) {
		mockStore := &MockGroupPostStore{
			CreateGroupPostFunc: func(post *models.GroupPost) (*models.GroupPost, error) {
				post.ID = 1
				return post, nil
			},
		}
		mockGroupStore := &MockGroupStore{}
		service := NewGroupPostService(mockStore, nil, mockGroupStore)

		post := &models.GroupPost{
			GroupID: "1",
			UserID:  1,
			Content: "Test content",
		}

		createdPost, err := service.CreateGroupPost(post, nil, "")
		if err != nil {
			t.Fatalf("CreateGroupPost failed: %v", err)
		}

		if createdPost.ID != 1 {
			t.Errorf("Expected post ID 1, got %d", createdPost.ID)
		}
	})

	t.Run("Store error", func(t *testing.T) {
		mockStore := &MockGroupPostStore{
			CreateGroupPostFunc: func(post *models.GroupPost) (*models.GroupPost, error) {
				return nil, errors.New("store error")
			},
		}
		mockGroupStore := &MockGroupStore{}
		service := NewGroupPostService(mockStore, nil, mockGroupStore)

		post := &models.GroupPost{
			GroupID: "1",
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
					{ID: 1, GroupID: "1", UserID: userID, Content: "Post 1"},
					{ID: 2, GroupID: "1", UserID: userID, Content: "Post 2"},
				}, nil
			},
		}
		mockGroupStore := &MockGroupStore{
			GetGroupByIDFunc: func(publicID string) (*models.Group, error) {
				return &models.Group{ID: "1", PublicID: publicID}, nil
			},
		}
		service := NewGroupPostService(mockStore, nil, mockGroupStore)

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
		mockGroupStore := &MockGroupStore{}
		service := NewGroupPostService(mockStore, nil, mockGroupStore)

		comment := &models.GroupPostComment{
			GroupPostID: "1",
			UserID:      1,
			Content:     "Test comment",
		}

		createdComment, err := service.CreateGroupPostComment(comment, nil, "")
		if err != nil {
			t.Fatalf("CreateGroupPostComment failed: %v", err)
		}

		if createdComment.ID != 1 {
			t.Errorf("Expected comment ID 1, got %d", createdComment.ID)
		}
	})
}