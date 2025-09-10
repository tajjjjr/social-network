package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/pkg/utils"
)

// MockGroupService is a mock implementation of the GroupService for testing.
type MockGroupService struct {
	CreateGroupFunc        func(group *models.Group) (*models.Group, error)
	GetGroupByIDFunc       func(groupID string) (*models.Group, error)
	SearchPublicGroupsFunc func(query string) ([]*models.Group, error)
	GetAllPublicGroupsFunc func() ([]*models.Group, error)
	GetUserGroupsFunc      func(userID int64) ([]*models.Group, error)
	JoinGroupFunc          func(groupID string, userID int64) error
	LeaveGroupFunc         func(groupID string, userID int64) error
	IsGroupMemberFunc      func(groupID string, userID int64) (bool, error)
}

func (m *MockGroupService) CreateGroup(group *models.Group) (*models.Group, error) {
	if m.CreateGroupFunc != nil {
		return m.CreateGroupFunc(group)
	}
	return group, nil
}

func (m *MockGroupService) GetGroupByID(groupID string) (*models.Group, error) {
	if m.GetGroupByIDFunc != nil {
		return m.GetGroupByIDFunc(groupID)
	}
	return &models.Group{PublicID: groupID}, nil
}

func (m *MockGroupService) SearchPublicGroups(query string) ([]*models.Group, error) {
	if m.SearchPublicGroupsFunc != nil {
		return m.SearchPublicGroupsFunc(query)
	}
	return []*models.Group{}, nil
}

func (m *MockGroupService) GetAllPublicGroups() ([]*models.Group, error) {
	if m.GetAllPublicGroupsFunc != nil {
		return m.GetAllPublicGroupsFunc()
	}
	return []*models.Group{}, nil
}

func (m *MockGroupService) GetUserGroups(userID int64) ([]*models.Group, error) {
	if m.GetUserGroupsFunc != nil {
		return m.GetUserGroupsFunc(userID)
	}
	return []*models.Group{}, nil
}

func (m *MockGroupService) JoinGroup(groupID string, userID int64) error {
	if m.JoinGroupFunc != nil {
		return m.JoinGroupFunc(groupID, userID)
	}
	return nil
}

func (m *MockGroupService) LeaveGroup(groupID string, userID int64) error {
	if m.LeaveGroupFunc != nil {
		return m.LeaveGroupFunc(groupID, userID)
	}
	return nil
}

func (m *MockGroupService) IsGroupMember(groupID string, userID int64) (bool, error) {
	if m.IsGroupMemberFunc != nil {
		return m.IsGroupMemberFunc(groupID, userID)
	}
	return true, nil
}

// MockGroupRequestService is a mock implementation of the GroupRequestService for testing.
type MockGroupRequestService struct {
	SendJoinRequestFunc    func(groupID string, userID int64) (*models.GroupRequest, error)
	ApproveJoinRequestFunc func(requestID int64, approverID int64) error
	RejectJoinRequestFunc  func(requestID int64, rejecterID int64) error
}

func (m *MockGroupRequestService) SendJoinRequest(groupID string, userID int64) (*models.GroupRequest, error) {
	if m.SendJoinRequestFunc != nil {
		return m.SendJoinRequestFunc(groupID, userID)
	}
	return &models.GroupRequest{ID: "req-1", GroupID: groupID, UserID: userID, Status: "pending"}, nil
}

func (m *MockGroupRequestService) ApproveJoinRequest(requestID int64, approverID int64) error {
	if m.ApproveJoinRequestFunc != nil {
		return m.ApproveJoinRequestFunc(requestID, approverID)
	}
	return nil
}

func (m *MockGroupRequestService) RejectJoinRequest(requestID int64, rejecterID int64) error {
	if m.RejectJoinRequestFunc != nil {
		return m.RejectJoinRequestFunc(requestID, rejecterID)
	}
	return nil
}

// MockGroupChatMessageService is a mock implementation of the GroupChatMessageService for testing.
type MockGroupChatMessageService struct {
	SendGroupChatMessageFunc func(groupID string, senderID int64, content string) (*models.GroupChatMessage, error)
	GetGroupChatMessagesFunc func(groupID string, userID int64, limit, offset int) ([]*models.GroupChatMessage, error)
}

func (m *MockGroupChatMessageService) SendGroupChatMessage(groupID string, senderID int64, content string) (*models.GroupChatMessage, error) {
	if m.SendGroupChatMessageFunc != nil {
		return m.SendGroupChatMessageFunc(groupID, senderID, content)
	}
	return &models.GroupChatMessage{ID: "msg-1", GroupID: "1", SenderID: senderID, Content: content}, nil
}

func (m *MockGroupChatMessageService) GetGroupChatMessages(groupID string, userID int64, limit, offset int) ([]*models.GroupChatMessage, error) {
	if m.GetGroupChatMessagesFunc != nil {
		return m.GetGroupChatMessagesFunc(groupID, userID, limit, offset)
	}
	return []*models.GroupChatMessage{}, nil
}

// MockGroupMemberService is a mock implementation of the GroupMemberService for testing.
type MockGroupMemberService struct {
	GetGroupMembersFunc func(groupID string) ([]*models.User, error)
	AddGroupMemberFunc  func(groupID string, userID int64, role string) (*models.GroupMember, error)
}

func (m *MockGroupMemberService) GetGroupMembers(groupID string) ([]*models.User, error) {
	if m.GetGroupMembersFunc != nil {
		return m.GetGroupMembersFunc(groupID)
	}
	return []*models.User{}, nil
}

func (m *MockGroupMemberService) AddGroupMember(groupID string, userID int64, role string) (*models.GroupMember, error) {
	if m.AddGroupMemberFunc != nil {
		return m.AddGroupMemberFunc(groupID, userID, role)
	}
	return &models.GroupMember{ID: 1, GroupID: groupID, UserID: userID, Role: role}, nil
}

// MockGroupEventService is a mock implementation of the GroupEventService for testing.
type MockGroupEventService struct {
	CreateGroupEventFunc func(event *models.GroupEvent) (*models.GroupEvent, error)
	GetGroupEventsFunc   func(groupID string) ([]*models.GroupEvent, error)
}

func (m *MockGroupEventService) CreateGroupEvent(event *models.GroupEvent) (*models.GroupEvent, error) {
	if m.CreateGroupEventFunc != nil {
		return m.CreateGroupEventFunc(event)
	}
	return event, nil
}

func (m *MockGroupEventService) GetGroupEvents(groupID string) ([]*models.GroupEvent, error) {
	if m.GetGroupEventsFunc != nil {
		return m.GetGroupEventsFunc(groupID)
	}
	return []*models.GroupEvent{}, nil
}

func TestCreateGroup(t *testing.T) {
	t.Run("Successful group creation", func(t *testing.T) {
		mockGroupService := &MockGroupService{
			CreateGroupFunc: func(group *models.Group) (*models.Group, error) {
				group.PublicID = "group-1"
				return group, nil
			},
		}

		h := NewGroupHandler(mockGroupService, &MockGroupRequestService{}, &MockGroupChatMessageService{}, &MockGroupMemberService{}, &MockGroupEventService{}, nil)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.WriteField("title", "Test Group")
		writer.WriteField("description", "This is a test group.")
		writer.WriteField("privacy", "public")
		writer.Close()

		req, err := http.NewRequest("POST", "/groups", body)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())

		ctx := context.WithValue(req.Context(), utils.User_id, int64(1))
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		h.CreateGroup(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}

		var createdGroup models.Group
		if err := json.NewDecoder(rr.Body).Decode(&createdGroup); err != nil {
			t.Fatal(err)
		}

		if createdGroup.PublicID != "group-1" {
			t.Errorf("handler returned unexpected group ID: got %v want %v", createdGroup.PublicID, "group-1")
		}
	})

	t.Run("User ID not found in context", func(t *testing.T) {
		h := NewGroupHandler(&MockGroupService{}, &MockGroupRequestService{}, &MockGroupChatMessageService{}, &MockGroupMemberService{}, &MockGroupEventService{}, nil)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.WriteField("title", "Test Group")
		writer.Close()

		req, err := http.NewRequest("POST", "/groups", body)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())

		rr := httptest.NewRecorder()
		h.CreateGroup(rr, req)

		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
		}
	})
}

func TestSendJoinRequest(t *testing.T) {
	t.Run("Successful join request", func(t *testing.T) {
		mockGroupRequestService := &MockGroupRequestService{
			SendJoinRequestFunc: func(groupID string, userID int64) (*models.GroupRequest, error) {
				return &models.GroupRequest{ID: "req-1", GroupID: groupID, UserID: userID, Status: "pending"}, nil
			},
		}

		h := NewGroupHandler(&MockGroupService{}, mockGroupRequestService, &MockGroupChatMessageService{}, &MockGroupMemberService{}, &MockGroupEventService{}, nil)

		req, err := http.NewRequest("POST", "/groups/1/join-request", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.SetPathValue("groupID", "1")

		ctx := context.WithValue(req.Context(), utils.User_id, int64(101))
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		h.SendJoinRequest(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}
	})

	t.Run("User ID not found in context", func(t *testing.T) {
		h := NewGroupHandler(&MockGroupService{}, &MockGroupRequestService{}, &MockGroupChatMessageService{}, &MockGroupMemberService{}, &MockGroupEventService{}, nil)

		req, err := http.NewRequest("POST", "/groups/1/join-request", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.SetPathValue("groupID", "1")

		rr := httptest.NewRecorder()
		h.SendJoinRequest(rr, req)

		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
		}
	})
}