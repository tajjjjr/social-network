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
}

func NewGroupHandler(groupService service.GroupService, groupRequestService service.GroupRequestService, groupChatMessageService service.GroupChatMessageService) *GroupHandler {
	return &GroupHandler{groupService: groupService, groupRequestService: groupRequestService, groupChatMessageService: groupChatMessageService}
}

func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var group models.Group
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	creatorID, ok := r.Context().Value(utils.User_id).(int64)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}
	group.CreatorID = creatorID

	if group.Privacy == "" {
		group.Privacy = "public"
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
	groupIDStr := r.PathValue("groupID")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value(utils.User_id).(int64)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

		request, err := h.groupRequestService.SendJoinRequest(int64(groupID), int64(userID))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to send join request: %v", err), http.StatusInternalServerError)
		return
	}

	resp := struct {
		Request *models.GroupRequest `json:"request"`
		Message string               `json:"message"`
	}{
		Request: request,
		Message: "Join request approved",
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
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
	groupIDStr := r.PathValue("groupID")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

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

	message, err := h.groupChatMessageService.SendGroupChatMessage(int64(groupID), int64(senderID), requestBody.Content)
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
	groupIDStr := r.PathValue("groupID")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

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

	messages, err := h.groupChatMessageService.GetGroupChatMessages(int64(groupID), int64(userID), limit, offset)
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
		http.Error(w, fmt.Sprintf("Failed to get user groups: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(groups); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}
