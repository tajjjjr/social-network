package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/service"
	"github.com/tajjjjr/social-network/backend/pkg/utils"
)

type GroupHandler struct {
	groupService            service.GroupService
	groupRequestService     service.GroupRequestService
	groupChatMessageService service.GroupChatMessageService
	groupMemberService      service.GroupMemberServiceInterface
	groupEventService       service.GroupEventServiceInterface
	groupPostService        service.GroupPostServiceInterface
}

func NewGroupHandler(groupService service.GroupService, groupRequestService service.GroupRequestService, groupChatMessageService service.GroupChatMessageService, groupMemberService service.GroupMemberServiceInterface, groupEventService service.GroupEventServiceInterface, groupPostService service.GroupPostServiceInterface) *GroupHandler {
	return &GroupHandler{groupService: groupService, groupRequestService: groupRequestService, groupChatMessageService: groupChatMessageService, groupMemberService: groupMemberService, groupEventService: groupEventService, groupPostService: groupPostService}
}

func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	creatorID, ok := r.Context().Value(utils.User_id).(int64)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	// Parse multipart form data
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	group := models.Group{
		CreatorID:   creatorID,
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Privacy:     r.FormValue("privacy"),
	}
	if group.Privacy == "" {
		group.Privacy = "public"
	}

	// Handle avatar upload if present
	file, header, err := r.FormFile("avatar")
	if err == nil {
		defer file.Close()
		avatarPath, err := UploadAvatarImage(file, header)
		if err != nil {
			http.Error(w, "Failed to upload avatar: "+err.Error(), http.StatusInternalServerError)
			return
		}
		group.Avatar = avatarPath
	}

	newGroup, err := h.groupService.CreateGroup(&group)
	if err != nil {
		http.Error(w, "Failed to create group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newGroup); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *GroupHandler) SendJoinRequest(w http.ResponseWriter, r *http.Request) {
	publicID := r.PathValue("groupID")

	userID, ok := r.Context().Value(utils.User_id).(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	_, err := h.groupRequestService.SendJoinRequest(publicID, userID)
	if err != nil {
		if err.Error() == "cannot send join request to a private group" || err.Error() == "user is already the group creator" || err.Error() == "user is already a member of this group" {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Successfully joined group"}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *GroupHandler) ApproveJoinRequest(w http.ResponseWriter, r *http.Request) {
	requestIDStr := r.PathValue("requestID")
	requestID, err := strconv.Atoi(requestIDStr)
	if err != nil {
		http.Error(w, "Invalid request ID", http.StatusBadRequest)
		return
	}

	approverID, ok := r.Context().Value(utils.User_id).(int64)
	if !ok {
		http.Error(w, "Approver ID not found in context", http.StatusUnauthorized)
		return
	}

	err = h.groupRequestService.ApproveJoinRequest(int64(requestID), int64(approverID))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to approve join request: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Join request approved"}); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *GroupHandler) RejectJoinRequest(w http.ResponseWriter, r *http.Request) {
	requestIDStr := r.PathValue("requestID")
	requestID, err := strconv.Atoi(requestIDStr)
	if err != nil {
		http.Error(w, "Invalid request ID", http.StatusBadRequest)
		return
	}

	rejecterID, ok := r.Context().Value(utils.User_id).(int64)
	if !ok {
		http.Error(w, "Rejecter ID not found in context", http.StatusUnauthorized)
		return
	}

	err = h.groupRequestService.RejectJoinRequest(int64(requestID), int64(rejecterID))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to reject join request: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Join request rejected"}); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *GroupHandler) SendGroupChatMessage(w http.ResponseWriter, r *http.Request) {
	publicID := r.PathValue("groupID")

	senderID, ok := r.Context().Value(utils.User_id).(int64)
	if !ok {
		http.Error(w, "Sender ID not found in context", http.StatusUnauthorized)
		return
	}

	var requestBody struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	message, err := h.groupChatMessageService.SendGroupChatMessage(publicID, senderID, requestBody.Content)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to send message: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(message); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *GroupHandler) GetGroupChatMessages(w http.ResponseWriter, r *http.Request) {
	publicID := r.PathValue("groupID")

	userID, ok := r.Context().Value(utils.User_id).(int64)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10 // Default limit
	}

	offsetStr := r.URL.Query().Get("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0 // Default offset
	}

	messages, err := h.groupChatMessageService.GetGroupChatMessages(publicID, userID, limit, offset)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get messages: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *GroupHandler) SearchPublicGroups(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")

	groups, err := h.groupService.SearchPublicGroups(query)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to search public groups: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(groups); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *GroupHandler) GetAllPublicGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := h.groupService.GetAllPublicGroups()
	if err != nil {
		// fmt.Printf("Error getting public groups: %v\n", err)
		http.Error(w, fmt.Sprintf("Failed to get public groups: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(groups); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *GroupHandler) GetUserGroups(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(utils.User_id).(int64)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	groups, err := h.groupService.GetUserGroups(userID)
	if err != nil {
		// fmt.Printf("Error getting user groups for user %d: %v\n", userID, err)
		http.Error(w, fmt.Sprintf("Failed to get user groups: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(groups); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *GroupHandler) GetGroupByID(w http.ResponseWriter, r *http.Request) {
	groupID := r.PathValue("groupID")
	group, err := h.groupService.GetGroupByID(groupID)
	if err != nil {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(group); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *GroupHandler) GetGroupEvents(w http.ResponseWriter, r *http.Request) {
	publicID := r.PathValue("groupID")

	events, err := h.groupEventService.GetGroupEvents(publicID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get group events: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(events); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *GroupHandler) CreateGroupEvent(w http.ResponseWriter, r *http.Request) {
	groupPublicID := r.PathValue("groupID")

	creatorID, ok := r.Context().Value(utils.User_id).(int64)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	var event models.GroupEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	event.GroupPubID = groupPublicID
	event.CreatorID = creatorID

	newEvent, err := h.groupEventService.CreateGroupEvent(&event)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create group event: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newEvent); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *GroupHandler) GetGroupMembers(w http.ResponseWriter, r *http.Request) {
	publicID := r.PathValue("groupID")

	members, err := h.groupMemberService.GetGroupMembers(publicID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get group members: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(members); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *GroupHandler) GetGroupStats(w http.ResponseWriter, r *http.Request) {
	publicID := r.PathValue("groupID")

	userID, ok := r.Context().Value(utils.User_id).(int64)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	// Get member count
	members, err := h.groupMemberService.GetGroupMembers(publicID)
	if err != nil {
		members = []*models.User{}
	}

	// Get posts count by fetching all posts
	posts, err := h.groupPostService.GetGroupPosts(publicID, userID, 1000, 0)
	if err != nil {
		posts = []*models.GroupPost{}
	}

	// Get events count
	events, err := h.groupEventService.GetGroupEvents(publicID)
	if err != nil {
		events = []*models.GroupEvent{}
	}

	stats := map[string]interface{}{
		"members": len(members),
		"posts":   len(posts),
		"events":  len(events),
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *GroupHandler) LeaveGroup(w http.ResponseWriter, r *http.Request) {
	groupID := r.PathValue("groupID")
	userID, ok := r.Context().Value(utils.User_id).(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := h.groupService.LeaveGroup(groupID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Left group successfully"}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
