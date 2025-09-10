package service

import (
	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/store"
)

type GroupEventServiceInterface interface {
	GetGroupEvents(groupPublicID string) ([]*models.GroupEvent, error)
	CreateGroupEvent(event *models.GroupEvent) (*models.GroupEvent, error)
}

type GroupEventService struct {
	groupEventStore store.GroupEventStore
}

func NewGroupEventService(groupEventStore store.GroupEventStore) GroupEventServiceInterface {
	return &GroupEventService{groupEventStore: groupEventStore}
}

func (s *GroupEventService) GetGroupEvents(groupPublicID string) ([]*models.GroupEvent, error) {
	return s.groupEventStore.GetGroupEvents(groupPublicID)
}

func (s *GroupEventService) CreateGroupEvent(event *models.GroupEvent) (*models.GroupEvent, error) {
	return s.groupEventStore.CreateGroupEvent(event)
}
