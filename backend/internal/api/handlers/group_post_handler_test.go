package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/pkg/utils"
)

// MockGroupPostService is a mock implementation of the GroupPostService for testing.
type MockGroupPostService struct {
	CreateGroupPostFunc        func(post *models.GroupPost, imageData []byte, imageMimeType string) (*models.GroupPost, error)
	GetGroupPostsFunc          func(publicID string, userID int64, limit, offset int) ([]*models.GroupPost, error)
	UpdateGroupPostFunc        func(postPublicID string, userID int64, content string, imageData []byte, imageMimeType string) (*models.GroupPost, error)
	DeleteGroupPostFunc        func(postPublicID string, userID int64) error
	CreateGroupPostCommentFunc func(comment *models.GroupPostComment, imageData []byte, imageMimeType string) (*models.GroupPostComment, error)
	GetGroupPostCommentsFunc   func(postPublicID string, userID int64) ([]*models.GroupPostComment, error)
	IsGroupMemberFunc          func(publicID string, userID int64) (bool, error)
}

func (m *MockGroupPostService) CreateGroupPost(post *models.GroupPost, imageData []byte, imageMimeType string) (*models.GroupPost, error) {
	if m.CreateGroupPostFunc != nil {
		return m.CreateGroupPostFunc(post, imageData, imageMimeType)
	}
	post.PublicID = "post-1"
	return post, nil
}

func (m *MockGroupPostService) GetGroupPostByID(postPublicID string) (*models.GroupPost, error) {
	return &models.GroupPost{PublicID: postPublicID}, nil
}

func (m *MockGroupPostService) GetGroupPosts(publicID string, userID int64, limit, offset int) ([]*models.GroupPost, error) {
	if m.GetGroupPostsFunc != nil {
		return m.GetGroupPostsFunc(publicID, userID, limit, offset)
	}
	return []*models.GroupPost{}, nil
}

func (m *MockGroupPostService) UpdateGroupPost(postPublicID string, userID int64, content string, imageData []byte, imageMimeType string) (*models.GroupPost, error) {
	if m.UpdateGroupPostFunc != nil {
		return m.UpdateGroupPostFunc(postPublicID, userID, content, imageData, imageMimeType)
	}
	return &models.GroupPost{PublicID: postPublicID, Content: content}, nil
}

func (m *MockGroupPostService) DeleteGroupPost(postPublicID string, userID int64) error {
	if m.DeleteGroupPostFunc != nil {
		return m.DeleteGroupPostFunc(postPublicID, userID)
	}
	return nil
}

func (m *MockGroupPostService) CreateGroupPostComment(comment *models.GroupPostComment, imageData []byte, imageMimeType string) (*models.GroupPostComment, error) {
	if m.CreateGroupPostCommentFunc != nil {
		return m.CreateGroupPostCommentFunc(comment, imageData, imageMimeType)
	}
	comment.PublicID = "comment-1"
	return comment, nil
}

func (m *MockGroupPostService) GetGroupPostComments(postPublicID string, userID int64) ([]*models.GroupPostComment, error) {
	if m.GetGroupPostCommentsFunc != nil {
		return m.GetGroupPostCommentsFunc(postPublicID, userID)
	}
	return []*models.GroupPostComment{}, nil
}

func (m *MockGroupPostService) UpdateGroupPostComment(commentPublicID string, userID int64, content string, imageData []byte, imageMimeType string) (*models.GroupPostComment, error) {
	return &models.GroupPostComment{PublicID: commentPublicID, Content: content}, nil
}

func (m *MockGroupPostService) DeleteGroupPostComment(commentPublicID string, userID int64) error {
	return nil
}

func (m *MockGroupPostService) IsGroupMember(publicID string, userID int64) (bool, error) {
	if m.IsGroupMemberFunc != nil {
		return m.IsGroupMemberFunc(publicID, userID)
	}
	return true, nil
}

func TestCreateGroupPost(t *testing.T) {
	t.Run("Successful post creation", func(t *testing.T) {
		mockService := &MockGroupPostService{
			CreateGroupPostFunc: func(post *models.GroupPost, imageData []byte, imageMimeType string) (*models.GroupPost, error) {
				post.PublicID = "post-1"
				return post, nil
			},
			IsGroupMemberFunc: func(publicID string, userID int64) (bool, error) {
				return true, nil
			},
		}

		handler := NewGroupPostHandler(mockService)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.WriteField("content", "Test post content")
		writer.Close()

		req, err := http.NewRequest("POST", "/groups/1/posts", body)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.SetPathValue("groupID", "1")

		ctx := context.WithValue(req.Context(), utils.User_id, int64(1))
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		handler.CreateGroupPost(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}

		var createdPost models.GroupPost
		if err := json.NewDecoder(rr.Body).Decode(&createdPost); err != nil {
			t.Fatal(err)
		}

		if createdPost.PublicID != "post-1" {
			t.Errorf("Expected post PublicID 'post-1', got %s", createdPost.PublicID)
		}
	})

	t.Run("User not a group member", func(t *testing.T) {
		mockService := &MockGroupPostService{
			IsGroupMemberFunc: func(publicID string, userID int64) (bool, error) {
				return false, nil
			},
		}

		handler := NewGroupPostHandler(mockService)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.WriteField("content", "Test post content")
		writer.Close()

		req, err := http.NewRequest("POST", "/groups/1/posts", body)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.SetPathValue("groupID", "1")

		ctx := context.WithValue(req.Context(), utils.User_id, int64(1))
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		handler.CreateGroupPost(rr, req)

		if status := rr.Code; status != http.StatusForbidden {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusForbidden)
		}
	})

	t.Run("User ID not found in context", func(t *testing.T) {
		mockService := &MockGroupPostService{}
		handler := NewGroupPostHandler(mockService)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.WriteField("content", "Test post content")
		writer.Close()

		req, err := http.NewRequest("POST", "/groups/1/posts", body)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.SetPathValue("groupID", "1")

		rr := httptest.NewRecorder()
		handler.CreateGroupPost(rr, req)

		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
		}
	})
}

func TestGetGroupPosts(t *testing.T) {
	t.Run("Successful get posts", func(t *testing.T) {
		mockService := &MockGroupPostService{
			GetGroupPostsFunc: func(publicID string, userID int64, limit, offset int) ([]*models.GroupPost, error) {
				return []*models.GroupPost{
					{PublicID: "post-1", Content: "Post 1"},
					{PublicID: "post-2", Content: "Post 2"},
				}, nil
			},
			IsGroupMemberFunc: func(publicID string, userID int64) (bool, error) {
				return true, nil
			},
		}

		handler := NewGroupPostHandler(mockService)

		req, err := http.NewRequest("GET", "/groups/1/posts", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.SetPathValue("groupID", "1")

		ctx := context.WithValue(req.Context(), utils.User_id, int64(1))
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		handler.GetGroupPosts(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var posts []*models.GroupPost
		if err := json.NewDecoder(rr.Body).Decode(&posts); err != nil {
			t.Fatal(err)
		}

		if len(posts) != 2 {
			t.Errorf("Expected 2 posts, got %d", len(posts))
		}
	})

	t.Run("Service error", func(t *testing.T) {
		mockService := &MockGroupPostService{
			IsGroupMemberFunc: func(publicID string, userID int64) (bool, error) {
				return false, errors.New("service error")
			},
		}

		handler := NewGroupPostHandler(mockService)

		req, err := http.NewRequest("GET", "/groups/1/posts", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.SetPathValue("groupID", "1")

		ctx := context.WithValue(req.Context(), utils.User_id, int64(1))
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		handler.GetGroupPosts(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
	})
}