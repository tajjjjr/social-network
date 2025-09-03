package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tajjjjr/social-network/backend/internal/models"
)

// MockGroupService is a mock implementation of the GroupService for testing.
type MockGroupService struct {
	CreateGroupFunc      func(group *models.Group) (*models.Group, error)
	GetGroupByIDFunc     func(groupID int64) (*models.Group, error)
	SearchPublicGroupsFunc func(query string) ([]*models.Group, error)
	GetAllPublicGroupsFunc func() ([]*models.Group, error)
	GetUserGroupsFunc    func(userID int64) ([]*models.Group, error)
}

func (m *MockGroupService) CreateGroup(group *models.Group) (*models.Group, error) {
	if m.CreateGroupFunc != nil {
		return m.CreateGroupFunc(group)
	}
	return group, nil
}

func (m *MockGroupService) GetGroupByID(groupID int64) (*models.Group, error) {
	if m.GetGroupByIDFunc != nil {
		return m.GetGroupByIDFunc(groupID)
	}
	return nil, errors.New("GetGroupByID not implemented")
}

func (m *MockGroupService) SearchPublicGroups(query string) ([]*models.Group, error) {
	if m.SearchPublicGroupsFunc != nil {
		return m.SearchPublicGroupsFunc(query)
	}
	return nil, errors.New("SearchPublicGroups not implemented")
}

func (m *MockGroupService) GetAllPublicGroups() ([]*models.Group, error) {
	if m.GetAllPublicGroupsFunc != nil {
		return m.GetAllPublicGroupsFunc()
	}
	return nil, errors.New("GetAllPublicGroups not implemented")
}

func (m *MockGroupService) GetUserGroups(userID int64) ([]*models.Group, error) {
	if m.GetUserGroupsFunc != nil {
		return m.GetUserGroupsFunc(userID)
	}
	return nil, errors.New("GetUserGroups not implemented")
}

// MockGroupRequestService is a mock implementation of the GroupRequestService for testing.
type MockGroupRequestService struct {
	SendJoinRequestFunc    func(groupID, userID int64) (*models.GroupRequest, error)
	ApproveJoinRequestFunc func(requestID int64, approverID int64) error
	RejectJoinRequestFunc  func(requestID int64, rejecterID int64) error
}

func (m *MockGroupRequestService) SendJoinRequest(groupID, userID int64) (*models.GroupRequest, error) {
	if m.SendJoinRequestFunc != nil {
		return m.SendJoinRequestFunc(groupID, userID)
	}
	return nil, errors.New("SendJoinRequest not implemented")
}

func (m *MockGroupRequestService) ApproveJoinRequest(requestID int64, approverID int64) error {
	if m.ApproveJoinRequestFunc != nil {
		return m.ApproveJoinRequestFunc(requestID, approverID)
	}
	return errors.New("ApproveJoinRequest not implemented")
}

func (m *MockGroupRequestService) RejectJoinRequest(requestID int64, rejecterID int64) error {
	if m.RejectJoinRequestFunc != nil {
		return m.RejectJoinRequestFunc(requestID, rejecterID)
	}
	return errors.New("RejectJoinRequest not implemented")
}

// MockGroupChatMessageService is a mock implementation of the GroupChatMessageService for testing.
type MockGroupChatMessageService struct {
	SendGroupChatMessageFunc func(groupID, senderID int64, content string) (*models.GroupChatMessage, error)
	GetGroupChatMessagesFunc func(groupID int64, userID int64, limit, offset int) ([]*models.GroupChatMessage, error)
}

func (m *MockGroupChatMessageService) SendGroupChatMessage(groupID, senderID int64, content string) (*models.GroupChatMessage, error) {
	if m.SendGroupChatMessageFunc != nil {
		return m.SendGroupChatMessageFunc(groupID, senderID, content)
	}
	return nil, errors.New("SendGroupChatMessage not implemented")
}

func (m *MockGroupChatMessageService) GetGroupChatMessages(groupID int64, userID int64, limit, offset int) ([]*models.GroupChatMessage, error) {
	if m.GetGroupChatMessagesFunc != nil {
		return m.GetGroupChatMessagesFunc(groupID, userID, limit, offset)
	}
	return nil, errors.New("GetGroupChatMessages not implemented")
}

func TestCreateGroup(t *testing.T) {
	// Test case 1: Successful group creation
	t.Run("Successful group creation", func(t *testing.T) {
		mockGroupService := &MockGroupService{
			CreateGroupFunc: func(group *models.Group) (*models.Group, error) {
				group.ID = 1
				return group, nil
			},
		}
		mockGroupRequestService := &MockGroupRequestService{}
		mockGroupChatMessageService := &MockGroupChatMessageService{}

		h := NewGroupHandler(mockGroupService, mockGroupRequestService, mockGroupChatMessageService)

		group := &models.Group{
			Title:       "Test Group",
			Description: "This is a test group.",
			CreatorID:   1,
			Privacy:     "public",
		}

		body, _ := json.Marshal(group)
		req, err := http.NewRequest("POST", "/groups", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		// Add user ID to context (as int64 for CreateGroup)
		ctx := context.WithValue(req.Context(), userIDKey, int64(1))
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		h.CreateGroup(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusCreated)
		}

		var createdGroup models.Group
		if err := json.NewDecoder(rr.Body).Decode(&createdGroup); err != nil {
			t.Fatal(err)
		}

		if createdGroup.ID != 1 {
			t.Errorf("handler returned unexpected group ID: got %v want %v",
				createdGroup.ID, 1)
		}
	})

	// Test case 2: Invalid request body
	t.Run("Invalid request body", func(t *testing.T) {
		mockGroupService := &MockGroupService{}
		mockGroupRequestService := &MockGroupRequestService{}
		mockGroupChatMessageService := &MockGroupChatMessageService{}

		h := NewGroupHandler(mockGroupService, mockGroupRequestService, mockGroupChatMessageService)

		req, err := http.NewRequest("POST", "/groups", bytes.NewReader([]byte("invalid json")))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		h.CreateGroup(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})

	// Test case 3: User ID not found in context
	t.Run("User ID not found in context", func(t *testing.T) {
		mockGroupService := &MockGroupService{}
		mockGroupRequestService := &MockGroupRequestService{}
		mockGroupChatMessageService := &MockGroupChatMessageService{}

		h := NewGroupHandler(mockGroupService, mockGroupRequestService, mockGroupChatMessageService)

		group := &models.Group{
			Title:       "Test Group",
			Description: "This is a test group.",
			Privacy:     "public",
		}

		body, _ := json.Marshal(group)
		req, err := http.NewRequest("POST", "/groups", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		h.CreateGroup(rr, req)

		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusUnauthorized)
		}
	})
}

func TestSendJoinRequest(t *testing.T) {
	// Test case 1: Successful join request
	t.Run("Successful join request", func(t *testing.T) {
		mockGroupService := &MockGroupService{
			GetGroupByIDFunc: func(groupID int64) (*models.Group, error) {
				return &models.Group{ID: groupID, Privacy: "public"}, nil
			},
		}
		mockGroupRequestService := &MockGroupRequestService{
			SendJoinRequestFunc: func(groupID, userID int64) (*models.GroupRequest, error) {
				return &models.GroupRequest{ID: 1, GroupID: groupID, UserID: userID, Status: "pending"}, nil
			},
		}
		mockGroupChatMessageService := &MockGroupChatMessageService{}

		h := NewGroupHandler(mockGroupService, mockGroupRequestService, mockGroupChatMessageService)

		req, err := http.NewRequest("POST", "/groups/1/join-request", nil)
		if err != nil {
			t.Fatal(err)
		}
		// Set path variables
		req.SetPathValue("groupID", "1")

		// Add user ID to context (as int64 for SendJoinRequest)
		ctx := context.WithValue(req.Context(), userIDKey, int64(101))
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		h.SendJoinRequest(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusCreated)
		}

		var response struct {
			Request *models.GroupRequest `json:"request"`
			Message string               `json:"message"`
		}
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		if response.Request.ID != 1 {
			t.Errorf("Expected request ID 1, got %v", response.Request.ID)
		}
	})

	// Test case 2: Invalid group ID
	t.Run("Invalid group ID", func(t *testing.T) {
		mockGroupService := &MockGroupService{}
		mockGroupRequestService := &MockGroupRequestService{}
		mockGroupChatMessageService := &MockGroupChatMessageService{}

		h := NewGroupHandler(mockGroupService, mockGroupRequestService, mockGroupChatMessageService)

		req, err := http.NewRequest("POST", "/groups/invalid/join-request", nil)
		if err != nil {
			t.Fatal(err)
		}
		// Set path variables
		req.SetPathValue("groupID", "invalid")

		// Add user ID to context (as int64 for SendJoinRequest)
		ctx := context.WithValue(req.Context(), userIDKey, int64(101))
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		h.SendJoinRequest(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})

	// Test case 3: User ID not found in context
	t.Run("User ID not found in context", func(t *testing.T) {
		mockGroupService := &MockGroupService{}
		mockGroupRequestService := &MockGroupRequestService{}
		mockGroupChatMessageService := &MockGroupChatMessageService{}

		h := NewGroupHandler(mockGroupService, mockGroupRequestService, mockGroupChatMessageService)

		req, err := http.NewRequest("POST", "/groups/1/join-request", nil)
		if err != nil {
			t.Fatal(err)
		}
		// Set path variables
		req.SetPathValue("groupID", "1")

		rr := httptest.NewRecorder()
		h.SendJoinRequest(rr, req)

		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusUnauthorized)
		}
	})

	// Test case 4: Service error
	t.Run("Service error", func(t *testing.T) {
		mockGroupService := &MockGroupService{
			GetGroupByIDFunc: func(groupID int64) (*models.Group, error) {
				return &models.Group{ID: groupID, Privacy: "public"}, nil
			},
		}
		mockGroupRequestService := &MockGroupRequestService{
			SendJoinRequestFunc: func(groupID, userID int64) (*models.GroupRequest, error) {
				return nil, errors.New("service error")
			},
		}
		mockGroupChatMessageService := &MockGroupChatMessageService{}

		h := NewGroupHandler(mockGroupService, mockGroupRequestService, mockGroupChatMessageService)

		req, err := http.NewRequest("POST", "/groups/1/join-request", nil)
		if err != nil {
			t.Fatal(err)
		}
		// Set path variables
		req.SetPathValue("groupID", "1")

		// Add user ID to context (as int64 for SendJoinRequest)
		ctx := context.WithValue(req.Context(), userIDKey, int64(101))
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		h.SendJoinRequest(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}
	})
}

func TestApproveJoinRequest(t *testing.T) {
	// Test case 1: Successful approval
	t.Run("Successful approval", func(t *testing.T) {
		mockGroupService := &MockGroupService{}
		mockGroupRequestService := &MockGroupRequestService{
			ApproveJoinRequestFunc: func(requestID int64, approverID int64) error {
				return nil
			},
		}
		mockGroupChatMessageService := &MockGroupChatMessageService{}

		h := NewGroupHandler(mockGroupService, mockGroupRequestService, mockGroupChatMessageService)

		req, err := http.NewRequest("PUT", "/groups/1/join-request/1/approve", nil)
		if err != nil {
			t.Fatal(err)
		}
		// Set path variables
		req.SetPathValue("groupID", "1")
		req.SetPathValue("requestID", "1")

		// Add user ID to context (as int for ApproveJoinRequest)
		ctx := context.WithValue(req.Context(), userIDKey, 101)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		h.ApproveJoinRequest(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		var response map[string]string
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		if response["message"] != "Join request approved" {
			t.Errorf("Expected message 'Join request approved', got %s", response["message"])
		}
	})

	// Test case 2: Invalid request ID
	t.Run("Invalid request ID", func(t *testing.T) {
		mockGroupService := &MockGroupService{}
		mockGroupRequestService := &MockGroupRequestService{}
		mockGroupChatMessageService := &MockGroupChatMessageService{}

		h := NewGroupHandler(mockGroupService, mockGroupRequestService, mockGroupChatMessageService)

		req, err := http.NewRequest("PUT", "/groups/1/join-request/invalid/approve", nil)
		if err != nil {
			t.Fatal(err)
		}
		// Set path variables
		req.SetPathValue("groupID", "1")
		req.SetPathValue("requestID", "invalid")

		// Add user ID to context (as int for ApproveJoinRequest)
		ctx := context.WithValue(req.Context(), userIDKey, 101)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		h.ApproveJoinRequest(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})

	// Test case 3: Approver ID not found in context
	t.Run("Approver ID not found in context", func(t *testing.T) {
		mockGroupService := &MockGroupService{}
		mockGroupRequestService := &MockGroupRequestService{}
		mockGroupChatMessageService := &MockGroupChatMessageService{}

		h := NewGroupHandler(mockGroupService, mockGroupRequestService, mockGroupChatMessageService)

		req, err := http.NewRequest("PUT", "/groups/1/join-request/1/approve", nil)
		if err != nil {
			t.Fatal(err)
		}
		// Set path variables
		req.SetPathValue("groupID", "1")
		req.SetPathValue("requestID", "1")

		rr := httptest.NewRecorder()
		h.ApproveJoinRequest(rr, req)

		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusUnauthorized)
		}
	})

	// Test case 4: Service error
	t.Run("Service error", func(t *testing.T) {
		mockGroupService := &MockGroupService{}
		mockGroupRequestService := &MockGroupRequestService{
			ApproveJoinRequestFunc: func(requestID int64, approverID int64) error {
				return errors.New("service error")
			},
		}
		mockGroupChatMessageService := &MockGroupChatMessageService{}

		h := NewGroupHandler(mockGroupService, mockGroupRequestService, mockGroupChatMessageService)

		req, err := http.NewRequest("PUT", "/groups/1/join-request/1/approve", nil)
		if err != nil {
			t.Fatal(err)
		}
		// Set path variables
		req.SetPathValue("groupID", "1")
		req.SetPathValue("requestID", "1")

		// Add user ID to context (as int for ApproveJoinRequest)
		ctx := context.WithValue(req.Context(), userIDKey, 101)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		h.ApproveJoinRequest(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}
	})
}

func TestRejectJoinRequest(t *testing.T) {
	// Test case 1: Successful rejection
	t.Run("Successful rejection", func(t *testing.T) {
		mockGroupService := &MockGroupService{}
		mockGroupRequestService := &MockGroupRequestService{
			RejectJoinRequestFunc: func(requestID int64, rejecterID int64) error {
				return nil
			},
		}
		mockGroupChatMessageService := &MockGroupChatMessageService{}

		h := NewGroupHandler(mockGroupService, mockGroupRequestService, mockGroupChatMessageService)

		req, err := http.NewRequest("PUT", "/groups/1/join-request/1/reject", nil)
		if err != nil {
			t.Fatal(err)
		}
		// Set path variables
		req.SetPathValue("groupID", "1")
		req.SetPathValue("requestID", "1")

		// Add user ID to context (as int for RejectJoinRequest)
		ctx := context.WithValue(req.Context(), userIDKey, 101)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		h.RejectJoinRequest(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		var response map[string]string
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		if response["message"] != "Join request rejected" {
			t.Errorf("Expected message 'Join request rejected', got %s", response["message"])
		}
	})

	// Test case 2: Invalid request ID
	t.Run("Invalid request ID", func(t *testing.T) {
		mockGroupService := &MockGroupService{}
		mockGroupRequestService := &MockGroupRequestService{}
		mockGroupChatMessageService := &MockGroupChatMessageService{}

		h := NewGroupHandler(mockGroupService, mockGroupRequestService, mockGroupChatMessageService)

		req, err := http.NewRequest("PUT", "/groups/1/join-request/invalid/reject", nil)
		if err != nil {
			t.Fatal(err)
		}
		// Set path variables
		req.SetPathValue("groupID", "1")
		req.SetPathValue("requestID", "invalid")

		// Add user ID to context (as int for RejectJoinRequest)
		ctx := context.WithValue(req.Context(), userIDKey, 101)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		h.RejectJoinRequest(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})

	// Test case 3: Rejecter ID not found in context
	t.Run("Rejecter ID not found in context", func(t *testing.T) {
		mockGroupService := &MockGroupService{}
		mockGroupRequestService := &MockGroupRequestService{}
		mockGroupChatMessageService := &MockGroupChatMessageService{}

		h := NewGroupHandler(mockGroupService, mockGroupRequestService, mockGroupChatMessageService)

		req, err := http.NewRequest("PUT", "/groups/1/join-request/1/reject", nil)
		if err != nil {
			t.Fatal(err)
		}
		// Set path variables
		req.SetPathValue("groupID", "1")
		req.SetPathValue("requestID", "1")

		rr := httptest.NewRecorder()
		h.RejectJoinRequest(rr, req)

		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusUnauthorized)
		}
	})

	// Test case 4: Service error
	t.Run("Service error", func(t *testing.T) {
		mockGroupService := &MockGroupService{}
		mockGroupRequestService := &MockGroupRequestService{
			RejectJoinRequestFunc: func(requestID int64, rejecterID int64) error {
				return errors.New("service error")
			},
		}
		mockGroupChatMessageService := &MockGroupChatMessageService{}

		h := NewGroupHandler(mockGroupService, mockGroupRequestService, mockGroupChatMessageService)

		req, err := http.NewRequest("PUT", "/groups/1/join-request/1/reject", nil)
		if err != nil {
			t.Fatal(err)
		}
		// Set path variables
		req.SetPathValue("groupID", "1")
		req.SetPathValue("requestID", "1")

		// Add user ID to context (as int for RejectJoinRequest)
		ctx := context.WithValue(req.Context(), userIDKey, 101)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		h.RejectJoinRequest(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}
	})
}

func TestSendGroupChatMessage(t *testing.T) {
	// Test case 1: Successful message sending
	t.Run("Successful message sending", func(t *testing.T) {
		mockGroupService := &MockGroupService{
			GetGroupByIDFunc: func(groupID int64) (*models.Group, error) {
				return &models.Group{ID: groupID, Privacy: "private"}, nil
			},
		}
		mockGroupRequestService := &MockGroupRequestService{}
		mockGroupChatMessageService := &MockGroupChatMessageService{
			SendGroupChatMessageFunc: func(groupID, senderID int64, content string) (*models.GroupChatMessage, error) {
				return &models.GroupChatMessage{ID: 1, GroupID: groupID, SenderID: senderID, Content: content}, nil
			},
		}

		h := NewGroupHandler(mockGroupService, mockGroupRequestService, mockGroupChatMessageService)

		body, _ := json.Marshal(map[string]string{"content": "Hello Group!"})
		req, err := http.NewRequest("POST", "/groups/1/chat", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		// Set path variables
		req.SetPathValue("groupID", "1")

		// Add user ID to context (as int for SendGroupChatMessage)
		ctx := context.WithValue(req.Context(), userIDKey, 101)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		h.SendGroupChatMessage(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusCreated)
		}

		var response models.GroupChatMessage
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		if response.Content != "Hello Group!" {
			t.Errorf("Expected message content 'Hello Group!', got %s", response.Content)
		}
	})

	// Test case 2: Invalid group ID
	t.Run("Invalid group ID", func(t *testing.T) {
		mockGroupService := &MockGroupService{}
		mockGroupRequestService := &MockGroupRequestService{}
		mockGroupChatMessageService := &MockGroupChatMessageService{}

		h := NewGroupHandler(mockGroupService, mockGroupRequestService, mockGroupChatMessageService)

		body, _ := json.Marshal(map[string]string{"content": "Hello Group!"})
		req, err := http.NewRequest("POST", "/groups/invalid/chat", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		// Set path variables
		req.SetPathValue("groupID", "invalid")

		// Add user ID to context (as int for SendGroupChatMessage)
		ctx := context.WithValue(req.Context(), userIDKey, 101)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		h.SendGroupChatMessage(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})

	// Test case 3: Invalid request body
	t.Run("Invalid request body", func(t *testing.T) {
		mockGroupService := &MockGroupService{}
		mockGroupRequestService := &MockGroupRequestService{}
		mockGroupChatMessageService := &MockGroupChatMessageService{}

		h := NewGroupHandler(mockGroupService, mockGroupRequestService, mockGroupChatMessageService)

		req, err := http.NewRequest("POST", "/groups/1/chat", bytes.NewReader([]byte("invalid json")))
		if err != nil {
			t.Fatal(err)
		}
		// Set path variables
		req.SetPathValue("groupID", "1")

		// Add user ID to context (as int for SendGroupChatMessage)
		ctx := context.WithValue(req.Context(), userIDKey, 101)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		h.SendGroupChatMessage(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})

	// Test case 4: Sender ID not found in context
	t.Run("Sender ID not found in context", func(t *testing.T) {
		mockGroupService := &MockGroupService{}
		mockGroupRequestService := &MockGroupRequestService{}
		mockGroupChatMessageService := &MockGroupChatMessageService{}

		h := NewGroupHandler(mockGroupService, mockGroupRequestService, mockGroupChatMessageService)

		body, _ := json.Marshal(map[string]string{"content": "Hello Group!"})
		req, err := http.NewRequest("POST", "/groups/1/chat", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		// Set path variables
		req.SetPathValue("groupID", "1")

		rr := httptest.NewRecorder()
		h.SendGroupChatMessage(rr, req)

		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusUnauthorized)
		}
	})

	// Test case 5: Service error
	t.Run("Service error", func(t *testing.T) {
		mockGroupService := &MockGroupService{
			GetGroupByIDFunc: func(groupID int64) (*models.Group, error) {
				return &models.Group{ID: groupID, Privacy: "private"}, nil
			},
		}
		mockGroupRequestService := &MockGroupRequestService{}
		mockGroupChatMessageService := &MockGroupChatMessageService{
			SendGroupChatMessageFunc: func(groupID, senderID int64, content string) (*models.GroupChatMessage, error) {
				return nil, errors.New("service error")
			},
		}

		h := NewGroupHandler(mockGroupService, mockGroupRequestService, mockGroupChatMessageService)

		body, _ := json.Marshal(map[string]string{"content": "Hello Group!"})
		req, err := http.NewRequest("POST", "/groups/1/chat", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		// Set path variables
		req.SetPathValue("groupID", "1")

		// Add user ID to context (as int for SendGroupChatMessage)
		ctx := context.WithValue(req.Context(), userIDKey, 101)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		h.SendGroupChatMessage(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}
	})
}

func TestGetGroupChatMessages(t *testing.T) {
	// Test case 1: Successful message retrieval
	t.Run("Successful message retrieval", func(t *testing.T) {
		mockGroupService := &MockGroupService{
			GetGroupByIDFunc: func(groupID int64) (*models.Group, error) {
				return &models.Group{ID: groupID, Privacy: "private"}, nil
			},
		}
		mockGroupRequestService := &MockGroupRequestService{}
		mockGroupChatMessageService := &MockGroupChatMessageService{
			GetGroupChatMessagesFunc: func(groupID int64, userID int64, limit, offset int) ([]*models.GroupChatMessage, error) {
				return []*models.GroupChatMessage{
					{ID: 1, GroupID: groupID, SenderID: 101, Content: "Message 1"},
					{ID: 2, GroupID: groupID, SenderID: 102, Content: "Message 2"},
				}, nil
			},
		}

		h := NewGroupHandler(mockGroupService, mockGroupRequestService, mockGroupChatMessageService)

		req, err := http.NewRequest("GET", "/groups/1/chat?limit=10&offset=0", nil)
		if err != nil {
			t.Fatal(err)
		}
		// Set path variables
		req.SetPathValue("groupID", "1")

		// Add user ID to context (as int for GetGroupChatMessages)
		ctx := context.WithValue(req.Context(), userIDKey, 101)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		h.GetGroupChatMessages(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		var response []models.GroupChatMessage
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		if len(response) != 2 {
			t.Errorf("Expected 2 messages, got %d", len(response))
		}
		if response[0].Content != "Message 1" {
			t.Errorf("Expected first message content 'Message 1', got %s", response[0].Content)
		}
	})

	// Test case 2: Invalid group ID
	t.Run("Invalid group ID", func(t *testing.T) {
		mockGroupService := &MockGroupService{}
		mockGroupRequestService := &MockGroupRequestService{}
		mockGroupChatMessageService := &MockGroupChatMessageService{}

		h := NewGroupHandler(mockGroupService, mockGroupRequestService, mockGroupChatMessageService)

		req, err := http.NewRequest("GET", "/groups/invalid/chat", nil)
		if err != nil {
			t.Fatal(err)
		}
		// Set path variables
		req.SetPathValue("groupID", "invalid")

		// Add user ID to context (as int for GetGroupChatMessages)
		ctx := context.WithValue(req.Context(), userIDKey, 101)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		h.GetGroupChatMessages(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})

	// Test case 3: User ID not found in context
	t.Run("User ID not found in context", func(t *testing.T) {
		mockGroupService := &MockGroupService{}
		mockGroupRequestService := &MockGroupRequestService{}
		mockGroupChatMessageService := &MockGroupChatMessageService{}

		h := NewGroupHandler(mockGroupService, mockGroupRequestService, mockGroupChatMessageService)

		req, err := http.NewRequest("GET", "/groups/1/chat", nil)
		if err != nil {
			t.Fatal(err)
		}
		// Set path variables
		req.SetPathValue("groupID", "1")

		rr := httptest.NewRecorder()
		h.GetGroupChatMessages(rr, req)

		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusUnauthorized)
		}
	})

	// Test case 4: Service error
	t.Run("Service error", func(t *testing.T) {
		mockGroupService := &MockGroupService{
			GetGroupByIDFunc: func(groupID int64) (*models.Group, error) {
				return &models.Group{ID: groupID, Privacy: "private"}, nil
			},
		}
		mockGroupRequestService := &MockGroupRequestService{}
		mockGroupChatMessageService := &MockGroupChatMessageService{
			GetGroupChatMessagesFunc: func(groupID int64, userID int64, limit, offset int) ([]*models.GroupChatMessage, error) {
				return nil, errors.New("service error")
			},
		}

		h := NewGroupHandler(mockGroupService, mockGroupRequestService, mockGroupChatMessageService)

		req, err := http.NewRequest("GET", "/groups/1/chat", nil)
		if err != nil {
			t.Fatal(err)
		}
		// Set path variables
		req.SetPathValue("groupID", "1")

		// Add user ID to context (as int for GetGroupChatMessages)
		ctx := context.WithValue(req.Context(), userIDKey, 101)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		h.GetGroupChatMessages(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}
	})

	// Test case 5: Default limit and offset
	t.Run("Default limit and offset", func(t *testing.T) {
		mockGroupService := &MockGroupService{
			GetGroupByIDFunc: func(groupID int64) (*models.Group, error) {
				return &models.Group{ID: groupID, Privacy: "private"}, nil
			},
		}
		mockGroupRequestService := &MockGroupRequestService{}
		mockGroupChatMessageService := &MockGroupChatMessageService{
			GetGroupChatMessagesFunc: func(groupID int64, userID int64, limit, offset int) ([]*models.GroupChatMessage, error) {
				// Check that default values are used
				if limit != 10 || offset != 0 {
					return nil, errors.New("unexpected limit or offset")
				}
				return []*models.GroupChatMessage{
					{ID: 1, GroupID: groupID, SenderID: 101, Content: "Message 1"},
				}, nil
			},
		}

		h := NewGroupHandler(mockGroupService, mockGroupRequestService, mockGroupChatMessageService)

		req, err := http.NewRequest("GET", "/groups/1/chat", nil)
		if err != nil {
			t.Fatal(err)
		}
		// Set path variables
		req.SetPathValue("groupID", "1")

		// Add user ID to context (as int for GetGroupChatMessages)
		ctx := context.WithValue(req.Context(), userIDKey, 101)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		h.GetGroupChatMessages(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	})
}

func TestSearchPublicGroups(t *testing.T) {
	// Test case 1: Successful search
	t.Run("Successful search", func(t *testing.T) {
		mockGroupService := &MockGroupService{
			SearchPublicGroupsFunc: func(query string) ([]*models.Group, error) {
				return []*models.Group{
					{ID: 1, Title: "Public Group 1", Description: "Desc 1", Privacy: "public"},
					{ID: 2, Title: "Another Public Group", Description: "Desc 2", Privacy: "public"},
				}, nil
			},
		}
		mockGroupRequestService := &MockGroupRequestService{}
		mockGroupChatMessageService := &MockGroupChatMessageService{}

		h := NewGroupHandler(mockGroupService, mockGroupRequestService, mockGroupChatMessageService)

		req, err := http.NewRequest("GET", "/groups/search?query=public", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		h.SearchPublicGroups(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		var groups []models.Group
		if err := json.NewDecoder(rr.Body).Decode(&groups); err != nil {
			t.Fatal(err)
		}

		if len(groups) != 2 {
			t.Errorf("Expected 2 groups, got %d", len(groups))
		}
		if groups[0].Title != "Public Group 1" {
			t.Errorf("Expected first group title 'Public Group 1', got %s", groups[0].Title)
		}
	})

	// Test case 2: No results found
	t.Run("No results found", func(t *testing.T) {
		mockGroupService := &MockGroupService{
			SearchPublicGroupsFunc: func(query string) ([]*models.Group, error) {
				return []*models.Group{}, nil
			},
		}
		mockGroupRequestService := &MockGroupRequestService{}
		mockGroupChatMessageService := &MockGroupChatMessageService{}

		h := NewGroupHandler(mockGroupService, mockGroupRequestService, mockGroupChatMessageService)

		req, err := http.NewRequest("GET", "/groups/search?query=nonexistent", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		h.SearchPublicGroups(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		var groups []models.Group
		if err := json.NewDecoder(rr.Body).Decode(&groups); err != nil {
			t.Fatal(err)
		}

		if len(groups) != 0 {
			t.Errorf("Expected 0 groups, got %d", len(groups))
		}
	})

	// Test case 3: Service error
	t.Run("Service error", func(t *testing.T) {
		mockGroupService := &MockGroupService{
			SearchPublicGroupsFunc: func(query string) ([]*models.Group, error) {
				return nil, errors.New("service error")
			},
		}
		mockGroupRequestService := &MockGroupRequestService{}
		mockGroupChatMessageService := &MockGroupChatMessageService{}

		h := NewGroupHandler(mockGroupService, mockGroupRequestService, mockGroupChatMessageService)

		req, err := http.NewRequest("GET", "/groups/search?query=error", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		h.SearchPublicGroups(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}
	})
}