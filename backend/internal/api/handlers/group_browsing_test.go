package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/pkg/utils"
)

func setUserIDInContext(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, utils.User_id, userID)
}

type mockGroupService struct {
	groups     []*models.Group
	userGroups []*models.Group
}

func (m *mockGroupService) CreateGroup(group *models.Group) (*models.Group, error) {
	return group, nil
}

func (m *mockGroupService) GetGroupByID(groupID string) (*models.Group, error) {
	return &models.Group{PublicID: groupID}, nil
}

func (m *mockGroupService) SearchPublicGroups(query string) ([]*models.Group, error) {
	return m.groups, nil
}

func (m *mockGroupService) GetAllPublicGroups() ([]*models.Group, error) {
	return m.groups, nil
}

func (m *mockGroupService) GetUserGroups(userID int64) ([]*models.Group, error) {
	return m.userGroups, nil
}

// Add missing method to satisfy service.GroupService interface
func (m *mockGroupService) IsGroupMember(groupID string, userID int64) (bool, error) {
	// For testing, return true if userGroups contains the groupID
	for _, g := range m.userGroups {
		if g.PublicID == groupID {
			return true, nil
		}
	}
	return false, nil
}

// Implement JoinGroup to satisfy service.GroupService interface
func (m *mockGroupService) JoinGroup(groupID string, userID int64) error {
	// For testing, just append the group to userGroups if not already present
	for _, g := range m.userGroups {
		if g.PublicID == groupID {
			return nil // already joined
		}
	}
	m.userGroups = append(m.userGroups, &models.Group{PublicID: groupID})
	return nil
}

// Implement LeaveGroup to satisfy service.GroupService interface
func (m *mockGroupService) LeaveGroup(groupID string, userID int64) error {
	// For testing, remove the group from userGroups if present
	for i, g := range m.userGroups {
		if g.PublicID == groupID {
			m.userGroups = append(m.userGroups[:i], m.userGroups[i+1:]...)
			return nil
		}
	}
	return nil // group not found, nothing to do
}

func TestGetAllPublicGroups(t *testing.T) {
	mockGroups := []*models.Group{
		{PublicID: "1", Title: "Test Group 1", Description: "Description 1", Privacy: "public"},
		{PublicID: "2", Title: "Test Group 2", Description: "Description 2", Privacy: "public"},
	}

	mockService := &mockGroupService{groups: mockGroups}
	handler := &GroupHandler{groupService: mockService}

	req := httptest.NewRequest("GET", "/groups", nil)
	w := httptest.NewRecorder()

	handler.GetAllPublicGroups(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var groups []*models.Group
	if err := json.NewDecoder(w.Body).Decode(&groups); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(groups) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(groups))
	}
}

func TestGetUserGroups(t *testing.T) {
	mockUserGroups := []*models.Group{
		{PublicID: "1", Title: "My Group 1", Description: "My Description 1", Privacy: "public"},
	}

	mockService := &mockGroupService{userGroups: mockUserGroups}
	handler := &GroupHandler{groupService: mockService}

	req := httptest.NewRequest("GET", "/groups/my-groups", nil)
	req = req.WithContext(setUserIDInContext(req.Context(), int64(123)))
	w := httptest.NewRecorder()

	handler.GetUserGroups(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var groups []*models.Group
	if err := json.NewDecoder(w.Body).Decode(&groups); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(groups) != 1 {
		t.Errorf("Expected 1 group, got %d", len(groups))
	}
}
