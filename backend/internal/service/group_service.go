package service

import (
	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/store"
)

type groupService struct {
	groupStore store.GroupStore
}

func NewGroupService(groupStore store.GroupStore) GroupService {
	return &groupService{groupStore: groupStore}
}

func (s *groupService) CreateGroup(group *models.Group) (*models.Group, error) {
	return s.groupStore.CreateGroup(group)
}

func (s *groupService) GetGroupByID(groupID int64) (*models.Group, error) {
	return s.groupStore.GetGroupByID(groupID)
}

func (s *groupService) SearchPublicGroups(query string) ([]*models.Group, error) {
	return s.groupStore.SearchPublicGroups(query)
}

func (s *groupService) GetAllPublicGroups() ([]*models.Group, error) {
	return s.groupStore.GetAllPublicGroups()
}

func (s *groupService) GetUserGroups(userID int64) ([]*models.Group, error) {
	return s.groupStore.GetUserGroups(userID)
}
