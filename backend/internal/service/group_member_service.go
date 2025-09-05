package service

import (
	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/store"
)

type GroupMemberServiceInterface interface {
	GetGroupMembers(groupID int64) ([]*models.User, error)
}

type GroupMemberService struct {
	groupMemberStore store.GroupMemberStore
}

func NewGroupMemberService(groupMemberStore store.GroupMemberStore) GroupMemberServiceInterface {
	return &GroupMemberService{groupMemberStore: groupMemberStore}
}

func (s *GroupMemberService) GetGroupMembers(groupID int64) ([]*models.User, error) {
	return s.groupMemberStore.GetGroupMembers(groupID)
}