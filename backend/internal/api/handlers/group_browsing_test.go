package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tajjjjr/social-network/backend/internal/models"
)

func setUserIDInContext(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

type mockGroupService struct {
	groups    []*models.Group
	userGroups []*models.Group
}

func (m *mockGroupService) CreateGroup(group *models.Group) (*models.Group, error) {
	return group, nil
}

func (m *mockGroupService) GetGroupByID(groupID int64) (*models.Group, error) {
	return &models.Group{ID: groupID}, nil
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

func TestGetAllPublicGroups(t *testing.T) {
	mockGroups := []*models.Group{
		{ID: 1, Title: "Test Group 1", Description: "Description 1", Privacy: "public"},
		{ID: 2, Title: "Test Group 2", Description: "Description 2", Privacy: "public"},
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
		{ID: 1, Title: "My Group 1", Description: "My Description 1", Privacy: "public"},
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