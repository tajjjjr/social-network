package models

import (
	"testing"
	"time"
)

func TestGroupModel(t *testing.T) {
	group := Group{
		ID:          int64(1),
		PublicID:    "group-1",
		Title:       "Test Group",
		Description: "This is a test group.",
		CreatorID:   101,
		Privacy:     "public",
		CreatedAt:   time.Now(),
	}

	if group.ID != int64(1) {
		t.Errorf("Expected ID 1, got %d", group.ID)
	}
	if group.PublicID != "group-1" {
		t.Errorf("Expected PublicID 'group-1', got %s", group.PublicID)
	}
	if group.Title != "Test Group" {
		t.Errorf("Expected Title 'Test Group', got %s", group.Title)
	}
	if group.Description != "This is a test group." {
		t.Errorf("Expected Description 'This is a test group.', got %s", group.Description)
	}
	if group.CreatorID != 101 {
		t.Errorf("Expected CreatorID 101, got %d", group.CreatorID)
	}
	if group.Privacy != "public" {
		t.Errorf("Expected Privacy 'public', got %s", group.Privacy)
	}
}
