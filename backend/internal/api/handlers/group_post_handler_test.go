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

type MockGroupPostService struct {
	CreateGroupPostFunc        func(post *models.GroupPost, imageData []byte, imageMimeType string) (int64, error)
	GetGroupPostsFunc          func(groupID int64, userID int64, limit, offset int) ([]*models.GroupPost, error)
	UpdateGroupPostFunc        func(postID, userID int64, content string, imageData []byte, imageMimeType string) (*models.GroupPost, error)
	DeleteGroupPostFunc        func(postID, userID int64) error
	CreateGroupPostCommentFunc func(comment *models.GroupPostComment, imageData []byte, imageMimeType string) (int64, error)
	GetGroupPostCommentsFunc   func(postID int64, userID int64) ([]*models.GroupPostComment, error)
}

func (m *MockGroupPostService) CreateGroupPost(post *models.GroupPost, imageData []byte, imageMimeType string) (int64, error) {
	if m.CreateGroupPostFunc != nil {
		return m.CreateGroupPostFunc(post, imageData, imageMimeType)
	}
	return 1, nil
}

func (m *MockGroupPostService) GetGroupPostByID(postID int64) (*models.GroupPost, error) {
	return nil, errors.New("not implemented")
}

func (m *MockGroupPostService) GetGroupPosts(groupID int64, userID int64, limit, offset int) ([]*models.GroupPost, error) {
	if m.GetGroupPostsFunc != nil {
		return m.GetGroupPostsFunc(groupID, userID, limit, offset)
	}
	return []*models.GroupPost{}, nil
}

func (m *MockGroupPostService) UpdateGroupPost(postID, userID int64, content string, imageData []byte, imageMimeType string) (*models.GroupPost, error) {
	if m.UpdateGroupPostFunc != nil {
		return m.UpdateGroupPostFunc(postID, userID, content, imageData, imageMimeType)
	}
	return &models.GroupPost{ID: postID, Content: content}, nil
}

func (m *MockGroupPostService) DeleteGroupPost(postID, userID int64) error {
	if m.DeleteGroupPostFunc != nil {
		return m.DeleteGroupPostFunc(postID, userID)
	}
	return nil
}

func (m *MockGroupPostService) CreateGroupPostComment(comment *models.GroupPostComment, imageData []byte, imageMimeType string) (int64, error) {
	if m.CreateGroupPostCommentFunc != nil {
		return m.CreateGroupPostCommentFunc(comment, imageData, imageMimeType)
	}
	return 1, nil
}

func (m *MockGroupPostService) GetGroupPostComments(postID int64, userID int64) ([]*models.GroupPostComment, error) {
	if m.GetGroupPostCommentsFunc != nil {
		return m.GetGroupPostCommentsFunc(postID, userID)
	}
	return []*models.GroupPostComment{}, nil
}

func (m *MockGroupPostService) UpdateGroupPostComment(commentID, userID int64, content string, imageData []byte, imageMimeType string) (*models.GroupPostComment, error) {
	return nil, errors.New("not implemented")
}

func (m *MockGroupPostService) DeleteGroupPostComment(commentID, userID int64) error {
	return errors.New("not implemented")
}

func TestCreateGroupPost(t *testing.T) {
	t.Run("Successful group post creation", func(t *testing.T) {
		mockService := &MockGroupPostService{
			CreateGroupPostFunc: func(post *models.GroupPost, imageData []byte, imageMimeType string) (int64, error) {
				return 1, nil
			},
		}

		handler := NewGroupPostHandler(mockService)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		if err := writer.WriteField("content", "Test group post content"); err != nil {
			t.Fatal(err)
		}
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

		var response models.GroupPost
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		if response.ID != 1 {
			t.Errorf("Expected post ID 1, got %d", response.ID)
		}
	})

	t.Run("Invalid group ID", func(t *testing.T) {
		mockService := &MockGroupPostService{}
		handler := NewGroupPostHandler(mockService)

		req, err := http.NewRequest("POST", "/groups/invalid/posts", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.SetPathValue("groupID", "invalid")

		rr := httptest.NewRecorder()
		handler.CreateGroupPost(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockService := &MockGroupPostService{}
		handler := NewGroupPostHandler(mockService)

		req, err := http.NewRequest("POST", "/groups/1/posts", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.SetPathValue("groupID", "1")

		rr := httptest.NewRecorder()
		handler.CreateGroupPost(rr, req)

		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
		}
	})
}

func TestGetGroupPosts(t *testing.T) {
	t.Run("Successful get group posts", func(t *testing.T) {
		mockService := &MockGroupPostService{
			GetGroupPostsFunc: func(groupID int64, userID int64, limit, offset int) ([]*models.GroupPost, error) {
				return []*models.GroupPost{
					{ID: 1, GroupID: groupID, UserID: userID, Content: "Test post 1"},
					{ID: 2, GroupID: groupID, UserID: userID, Content: "Test post 2"},
				}, nil
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
}

func TestCreateGroupPostComment(t *testing.T) {
	t.Run("Successful comment creation", func(t *testing.T) {
		mockService := &MockGroupPostService{
			CreateGroupPostCommentFunc: func(comment *models.GroupPostComment, imageData []byte, imageMimeType string) (int64, error) {
				return 1, nil
			},
		}

		handler := NewGroupPostHandler(mockService)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		if err := writer.WriteField("content", "Test comment"); err != nil {
			t.Fatal(err)
		}
		writer.Close()

		req, err := http.NewRequest("POST", "/groups/1/posts/1/comments", body)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.SetPathValue("postID", "1")

		ctx := context.WithValue(req.Context(), utils.User_id, int64(1))
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		handler.CreateGroupPostComment(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}

		var response models.GroupPostComment
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		if response.ID != 1 {
			t.Errorf("Expected comment ID 1, got %d", response.ID)
		}
	})
}