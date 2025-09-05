package service

import (
	"github.com/tajjjjr/social-network/backend/internal/store"
)

type GroupEventServiceInterface interface {
	GetGroupEvents(groupID int64) ([]*store.GroupEvent, error)
}

type GroupEventService struct {
	groupEventStore store.GroupEventStore
}

func NewGroupEventService(groupEventStore store.GroupEventStore) GroupEventServiceInterface {
	return &GroupEventService{groupEventStore: groupEventStore}
}

func (s *GroupEventService) GetGroupEvents(groupID int64) ([]*store.GroupEvent, error) {
	return s.groupEventStore.GetGroupEvents(groupID)
}