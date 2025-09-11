package service

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/store"
	"github.com/tajjjjr/social-network/backend/pkg/utils"
)

type GroupPostServiceInterface interface {
	CreateGroupPost(post *models.GroupPost, imageData []byte, imageMimeType string) (*models.GroupPost, error)
	GetGroupPostByID(postPublicID string) (*models.GroupPost, error)
	GetGroupPosts(publicID string, userID int64, limit, offset int) ([]*models.GroupPost, error)
	UpdateGroupPost(postPublicID string, userID int64, content string, imageData []byte, imageMimeType string) (*models.GroupPost, error)
	DeleteGroupPost(postPublicID string, userID int64) error
	CreateGroupPostComment(comment *models.GroupPostComment, imageData []byte, imageMimeType string) (*models.GroupPostComment, error)
	GetGroupPostComments(postPublicID string, userID int64) ([]*models.GroupPostComment, error)
	UpdateGroupPostComment(commentPublicID string, userID int64, content string, imageData []byte, imageMimeType string) (*models.GroupPostComment, error)
	DeleteGroupPostComment(commentPublicID string, userID int64) error
	IsGroupMember(publicID string, userID int64) (bool, error)
}

type GroupPostService struct {
	groupPostStore   store.GroupPostStore
	groupMemberStore store.GroupMemberStore
	groupStore       store.GroupStore
}

func NewGroupPostService(groupPostStore store.GroupPostStore, groupMemberStore store.GroupMemberStore, groupStore store.GroupStore) GroupPostServiceInterface {
	return &GroupPostService{
		groupPostStore:   groupPostStore,
		groupMemberStore: groupMemberStore,
		groupStore:       groupStore,
	}
}

func (s *GroupPostService) CreateGroupPost(post *models.GroupPost, imageData []byte, imageMimeType string) (*models.GroupPost, error) {
	if len(imageData) > 0 {
		imagePath, err := s.saveImage(imageData, "posts")
		if err != nil {
			return nil, err
		}
		post.Image = imagePath
	}
	
	return s.groupPostStore.CreateGroupPost(post)
}

func (s *GroupPostService) GetGroupPostByID(postPublicID string) (*models.GroupPost, error) {
	return s.groupPostStore.GetGroupPostByID(postPublicID)
}

func (s *GroupPostService) GetGroupPosts(publicID string, userID int64, limit, offset int) ([]*models.GroupPost, error) {
	group, err := s.groupStore.GetGroupByID(publicID)
	if err != nil {
		return nil, err
	}
	return s.groupPostStore.GetGroupPosts(group.PublicID, userID, limit, offset)
}

func (s *GroupPostService) UpdateGroupPost(postPublicID string, userID int64, content string, imageData []byte, imageMimeType string) (*models.GroupPost, error) {
	// Get existing post to preserve image if no new image provided
	existingPost, err := s.groupPostStore.GetGroupPostByID(postPublicID)
	if err != nil {
		return nil, fmt.Errorf("post not found")
	}
	
	// Handle image update if provided
	var imagePath string
	if len(imageData) > 0 {
		savedImagePath, err := s.saveImage(imageData, "posts")
		if err != nil {
			return nil, err
		}
		imagePath = savedImagePath
	} else {
		// Keep existing image if no new image provided
		imagePath = existingPost.Image
	}
	
	return s.groupPostStore.UpdateGroupPost(postPublicID, userID, content, []byte(imagePath), "")
}

func (s *GroupPostService) DeleteGroupPost(postPublicID string, userID int64) error {
	return s.groupPostStore.DeleteGroupPost(postPublicID, userID)
}

func (s *GroupPostService) CreateGroupPostComment(comment *models.GroupPostComment, imageData []byte, imageMimeType string) (*models.GroupPostComment, error) {
	if len(imageData) > 0 {
		imagePath, err := s.saveImage(imageData, "comments")
		if err != nil {
			return nil, err
		}
		comment.Image = imagePath
	}
	
	return s.groupPostStore.CreateGroupPostComment(comment)
}

func (s *GroupPostService) GetGroupPostComments(postPublicID string, userID int64) ([]*models.GroupPostComment, error) {
	return s.groupPostStore.GetGroupPostComments(postPublicID, userID)
}

func (s *GroupPostService) UpdateGroupPostComment(commentPublicID string, userID int64, content string, imageData []byte, imageMimeType string) (*models.GroupPostComment, error) {
	// Handle image update if provided
	var imagePath string
	if len(imageData) > 0 {
		savedImagePath, err := s.saveImage(imageData, "comments")
		if err != nil {
			return nil, err
		}
		imagePath = savedImagePath
	}
	
	return s.groupPostStore.UpdateGroupPostComment(commentPublicID, userID, content, []byte(imagePath), "")
}

func (s *GroupPostService) DeleteGroupPostComment(commentPublicID string, userID int64) error {
	return s.groupPostStore.DeleteGroupPostComment(commentPublicID, userID)
}

func (s *GroupPostService) IsGroupMember(publicID string, userID int64) (bool, error) {
	return s.groupMemberStore.IsGroupMember(publicID, userID)
}

// saveImage handles the logic for validating, naming, and saving an uploaded image.
func (s *GroupPostService) saveImage(imageData []byte, subDir string) (string, error) {
	detectedFormat, err := utils.DetectImageFormat(bytes.NewReader(imageData))
	if err != nil {
		return "", fmt.Errorf("image signature check failed: %w", err)
	}
	extension, ok := groupFormatToExtension(detectedFormat)
	if !ok {
		return "", fmt.Errorf("unsupported image format: %s", detectedFormat)
	}
	imageFileName := fmt.Sprintf("%s%s", uuid.New().String(), extension)
	saveDir := filepath.Join("attachments", subDir)
	imagePath := filepath.Join(subDir, imageFileName)
	imageWritePath := filepath.Join(saveDir, imageFileName)
	if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}
	if err := os.WriteFile(imageWritePath, imageData, 0644); err != nil {
		return "", fmt.Errorf("failed to save image: %w", err)
	}
	return imagePath, nil
}

// groupFormatToExtension maps an ImageFormat to a file extension for group posts.
func groupFormatToExtension(format utils.ImageFormat) (string, bool) {
	switch format {
	case utils.JPEG:
		return ".jpg", true
	case utils.PNG:
		return ".png", true
	case utils.GIF:
		return ".gif", true
	case utils.WebP:
		return ".webp", true
	case utils.BMP:
		return ".bmp", true
	case utils.TIFF:
		return ".tiff", true
	default:
		return "", false
	}
}