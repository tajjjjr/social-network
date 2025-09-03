package service

import (
	"fmt"

	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/store"
)

type groupRequestService struct {
	groupRequestStore store.GroupRequestStore
	groupService      GroupService
}

func NewGroupRequestService(groupRequestStore store.GroupRequestStore, groupService GroupService) GroupRequestService {
	return &groupRequestService{groupRequestStore: groupRequestStore, groupService: groupService}
}

func (s *groupRequestService) SendJoinRequest(groupID, userID int64) (*models.GroupRequest, error) {
	group, err := s.groupService.GetGroupByID(int64(groupID))
	if err != nil {
		return nil, fmt.Errorf("group not found: %w", err)
	}

	if group.Privacy == "private" {
		return nil, fmt.Errorf("cannot send join request to a private group")
	}

	// For public groups, automatically approve and add user as member
	request := &models.GroupRequest{
		GroupID: int64(groupID),
		UserID:  int64(userID),
		Status:  "approved",
	}

	createdRequest, err := s.groupRequestStore.CreateGroupRequest(request)
	if err != nil {
		return nil, fmt.Errorf("failed to create group request: %w", err)
	}

	// Add user to Group_Members table
	err = s.groupRequestStore.AddUserToGroup(groupID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to add user to group: %w", err)
	}

	return createdRequest, nil
}

func (s *groupRequestService) ApproveJoinRequest(requestID int64, approverID int64) error {
	request, err := s.groupRequestStore.GetGroupRequestByID(int64(requestID))
	if err != nil {
		return fmt.Errorf("failed to get group request: %w", err)
	}

	if request.Status != "pending" {
		return fmt.Errorf("request is not pending")
	}

	err = s.groupRequestStore.UpdateGroupRequestStatus(requestID, "approved")
	if err != nil {
		return fmt.Errorf("failed to approve group request: %w", err)
	}

	return nil
}

func (s *groupRequestService) RejectJoinRequest(requestID int64, rejecterID int64) error {
	request, err := s.groupRequestStore.GetGroupRequestByID(int64(requestID))
	if err != nil {
		return fmt.Errorf("failed to get group request: %w", err)
	}

	if request.Status != "pending" {
		return fmt.Errorf("request is not pending")
	}

	err = s.groupRequestStore.UpdateGroupRequestStatus(requestID, "rejected")
	if err != nil {
		return fmt.Errorf("failed to reject group request: %w", err)
	}

	return nil
}
