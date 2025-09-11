package service

import (
	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/store"
)

type GroupMemberServiceInterface interface {
	GetGroupMembers(groupID string) ([]*models.User, error)
	AddGroupMember(groupID string, userID int64, role string) (*models.GroupMember, error)
}

type GroupMemberService struct {
	groupMemberStore store.GroupMemberStore
}

func NewGroupMemberService(groupMemberStore store.GroupMemberStore) GroupMemberServiceInterface {
	return &GroupMemberService{groupMemberStore: groupMemberStore}
}

func (s *GroupMemberService) GetGroupMembers(groupID string) ([]*models.User, error) {
	return s.groupMemberStore.GetGroupMembers(groupID)
}

func (s *GroupMemberService) AddGroupMember(groupID string, userID int64, role string) (*models.GroupMember, error) {
	return s.groupMemberStore.AddGroupMember(groupID, userID, role)
}
