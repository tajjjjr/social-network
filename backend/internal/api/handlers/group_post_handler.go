package handlers

import (
	"net/http"
	"strconv"

	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/service"
	"github.com/tajjjjr/social-network/backend/pkg/utils"
)

type GroupPostHandler struct {
	groupPostService service.GroupPostServiceInterface
}

func NewGroupPostHandler(groupPostService service.GroupPostServiceInterface) *GroupPostHandler {
	return &GroupPostHandler{groupPostService: groupPostService}
}

func (h *GroupPostHandler) CreateGroupPost(w http.ResponseWriter, r *http.Request) {
	groupID := r.PathValue("groupID")
	if groupID == "" {
		utils.RespondJSON(w, http.StatusBadRequest, utils.Response{Message: "Invalid group ID"})
		return
	}

	userID, ok := r.Context().Value(utils.User_id).(int64)
	if !ok {
		utils.RespondJSON(w, http.StatusUnauthorized, utils.Response{Message: "Unauthorized"})
		return
	}

	// Check if user is a member of the group
	isMember, err := h.groupPostService.IsGroupMember(groupID, userID)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Response{Message: "Failed to check group membership"})
		return
	}
	if !isMember {
		utils.RespondJSON(w, http.StatusForbidden, utils.Response{Message: "You must be a group member to post"})
		return
	}

	err = r.ParseMultipartForm(20 << 20)
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, utils.Response{Message: "Unable to parse form"})
		return
	}

	post := &models.GroupPost{
		GroupID: groupID,
		UserID:  userID,
		Content: r.FormValue("content"),
	}

	imageData, imageMimeType, status, err := handleImageUpload(r)
	if err != nil {
		utils.RespondJSON(w, status, utils.Response{Message: err.Error()})
		return
	}

	id, err := h.groupPostService.CreateGroupPost(post, imageData, imageMimeType)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Response{Message: err.Error()})
		return
	}

	post.ID = id
	utils.RespondJSON(w, http.StatusCreated, post)
}

func (h *GroupPostHandler) GetGroupPosts(w http.ResponseWriter, r *http.Request) {
	groupID := r.PathValue("groupID")
	if groupID == "" {
		utils.RespondJSON(w, http.StatusBadRequest, utils.Response{Message: "Invalid group ID"})
		return
	}

	userID, ok := r.Context().Value(utils.User_id).(int64)
	if !ok {
		utils.RespondJSON(w, http.StatusUnauthorized, utils.Response{Message: "Unauthorized"})
		return
	}

	// Check if user is a member of the group
	isMember, err := h.groupPostService.IsGroupMember(groupID, userID)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Response{Message: "Failed to check group membership"})
		return
	}
	if !isMember {
		utils.RespondJSON(w, http.StatusForbidden, utils.Response{Message: "You must be a group member to view posts"})
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 15
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 {
			limit = l
		}
	}

	offsetStr := r.URL.Query().Get("offset")
	offset := 0
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	posts, err := h.groupPostService.GetGroupPosts(groupID, userID, limit, offset)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Response{Message: "Internal server error"})
		return
	}

	utils.RespondJSON(w, http.StatusOK, posts)
}

func (h *GroupPostHandler) UpdateGroupPost(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.PathValue("postID")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, utils.Response{Message: "Invalid post ID"})
		return
	}

	userID, ok := r.Context().Value(utils.User_id).(int64)
	if !ok {
		utils.RespondJSON(w, http.StatusUnauthorized, utils.Response{Message: "Unauthorized"})
		return
	}

	err = r.ParseMultipartForm(20 << 20)
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, utils.Response{Message: "Unable to parse form"})
		return
	}

	content := r.FormValue("content")
	if content == "" {
		utils.RespondJSON(w, http.StatusBadRequest, utils.Response{Message: "Content is required"})
		return
	}

	imageData, imageMimeType, status, err := handleImageUpload(r)
	if err != nil {
		utils.RespondJSON(w, status, utils.Response{Message: err.Error()})
		return
	}

	updatedPost, err := h.groupPostService.UpdateGroupPost(postID, userID, content, imageData, imageMimeType)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Response{Message: err.Error()})
		return
	}

	utils.RespondJSON(w, http.StatusOK, updatedPost)
}

func (h *GroupPostHandler) DeleteGroupPost(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.PathValue("postID")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, utils.Response{Message: "Invalid post ID"})
		return
	}

	userID, ok := r.Context().Value(utils.User_id).(int64)
	if !ok {
		utils.RespondJSON(w, http.StatusUnauthorized, utils.Response{Message: "Unauthorized"})
		return
	}

	err = h.groupPostService.DeleteGroupPost(postID, userID)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Response{Message: err.Error()})
		return
	}

	utils.RespondJSON(w, http.StatusNoContent, utils.Response{Message: "Post deleted successfully"})
}

func (h *GroupPostHandler) CreateGroupPostComment(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.PathValue("postID")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, utils.Response{Message: "Invalid post ID"})
		return
	}

	userID, ok := r.Context().Value(utils.User_id).(int64)
	if !ok {
		utils.RespondJSON(w, http.StatusUnauthorized, utils.Response{Message: "Unauthorized"})
		return
	}

	err = r.ParseMultipartForm(20 << 20)
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, utils.Response{Message: "Unable to parse form"})
		return
	}

	comment := &models.GroupPostComment{
		GroupPostID: postID,
		UserID:      userID,
		Content:     r.FormValue("content"),
	}

	imageData, imageMimeType, status, err := handleImageUpload(r)
	if err != nil {
		utils.RespondJSON(w, status, utils.Response{Message: err.Error()})
		return
	}

	id, err := h.groupPostService.CreateGroupPostComment(comment, imageData, imageMimeType)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Response{Message: err.Error()})
		return
	}

	comment.ID = id
	utils.RespondJSON(w, http.StatusCreated, comment)
}

func (h *GroupPostHandler) GetGroupPostComments(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.PathValue("postID")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, utils.Response{Message: "Invalid post ID"})
		return
	}

	userID, ok := r.Context().Value(utils.User_id).(int64)
	if !ok {
		utils.RespondJSON(w, http.StatusUnauthorized, utils.Response{Message: "Unauthorized"})
		return
	}

	comments, err := h.groupPostService.GetGroupPostComments(postID, userID)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Response{Message: "Internal server error"})
		return
	}

	utils.RespondJSON(w, http.StatusOK, comments)
}