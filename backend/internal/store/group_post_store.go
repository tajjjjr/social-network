package store

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/tajjjjr/social-network/backend/internal/models"
)

type GroupPostStore interface {
	CreateGroupPost(post *models.GroupPost) (*models.GroupPost, error)
	GetGroupPostByID(postPublicID string) (*models.GroupPost, error)
	GetGroupPosts(groupID string, userID int64, limit, offset int) ([]*models.GroupPost, error)
	UpdateGroupPost(postPublicID string, userID int64, content string, imageData []byte, imageMimeType string) (*models.GroupPost, error)
	DeleteGroupPost(postPublicID string, userID int64) error
	CreateGroupPostComment(comment *models.GroupPostComment) (*models.GroupPostComment, error)
	GetGroupPostComments(postPublicID string, userID int64) ([]*models.GroupPostComment, error)
	UpdateGroupPostComment(commentPublicID string, userID int64, content string, imageData []byte, imageMimeType string) (*models.GroupPostComment, error)
	DeleteGroupPostComment(commentPublicID string, userID int64) error
	CanUserDeleteGroupContent(groupID string, userID int64) (bool, error)
}

type groupPostStore struct {
	db *sql.DB
}

func NewGroupPostStore(db *sql.DB) GroupPostStore {
	return &groupPostStore{db: db}
}

func (s *groupPostStore) CreateGroupPost(post *models.GroupPost) (*models.GroupPost, error) {
	post.PublicID = uuid.New().String()
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	
	_, err := s.db.Exec(`
		INSERT INTO Groups (id, public_id, type, group_id, user_id, content, image, created_at, updated_at)
		VALUES (?, ?, 'post', ?, ?, ?, ?, ?, ?)
	`, post.PublicID, post.PublicID, post.GroupID, post.UserID, post.Content, post.Image, post.CreatedAt, post.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *groupPostStore) GetGroupPostByID(postPublicID string) (*models.GroupPost, error) {
	var post models.GroupPost
	err := s.db.QueryRow(`
		SELECT public_id, group_id, user_id, content, COALESCE(image, ''), created_at, updated_at
		FROM Groups WHERE public_id = ? AND type = 'post'
	`, postPublicID).Scan(&post.PublicID, &post.GroupID, &post.UserID, &post.Content, &post.Image, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (s *groupPostStore) GetGroupPosts(groupID string, userID int64, limit, offset int) ([]*models.GroupPost, error) {
	rows, err := s.db.Query(`
		SELECT public_id, group_id, user_id, content, COALESCE(image, ''), created_at, updated_at
		FROM Groups WHERE group_id = ? AND type = 'post'
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, groupID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var posts []*models.GroupPost
	for rows.Next() {
		var post models.GroupPost
		err := rows.Scan(&post.PublicID, &post.GroupID, &post.UserID, &post.Content, &post.Image, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	return posts, nil
}

func (s *groupPostStore) UpdateGroupPost(postPublicID string, userID int64, content string, imageData []byte, imageMimeType string) (*models.GroupPost, error) {
	_, err := s.db.Exec(`
		UPDATE Groups SET content = ?, updated_at = CURRENT_TIMESTAMP
		WHERE public_id = ? AND user_id = ? AND type = 'post'
	`, content, postPublicID, userID)
	if err != nil {
		return nil, err
	}
	return s.GetGroupPostByID(postPublicID)
}

func (s *groupPostStore) DeleteGroupPost(postPublicID string, userID int64) error {
	_, err := s.db.Exec(`
		DELETE FROM Groups 
		WHERE public_id = ? AND user_id = ? AND type = 'post'
	`, postPublicID, userID)
	return err
}

func (s *groupPostStore) CreateGroupPostComment(comment *models.GroupPostComment) (*models.GroupPostComment, error) {
	comment.PublicID = uuid.New().String()
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()
	
	_, err := s.db.Exec(`
		INSERT INTO Groups (id, public_id, type, group_id, user_id, content, image, created_at, updated_at)
		VALUES (?, ?, 'comment', ?, ?, ?, ?, ?, ?)
	`, comment.PublicID, comment.PublicID, comment.GroupPostID, comment.UserID, comment.Content, comment.Image, comment.CreatedAt, comment.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (s *groupPostStore) GetGroupPostComments(postPublicID string, userID int64) ([]*models.GroupPostComment, error) {
	rows, err := s.db.Query(`
		SELECT public_id, group_id, user_id, content, COALESCE(image, ''), created_at, updated_at
		FROM Groups WHERE group_id = ? AND type = 'comment'
		ORDER BY created_at ASC
	`, postPublicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var comments []*models.GroupPostComment
	for rows.Next() {
		var comment models.GroupPostComment
		err := rows.Scan(&comment.PublicID, &comment.GroupPostID, &comment.UserID, &comment.Content, &comment.Image, &comment.CreatedAt, &comment.UpdatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	return comments, nil
}

func (s *groupPostStore) UpdateGroupPostComment(commentPublicID string, userID int64, content string, imageData []byte, imageMimeType string) (*models.GroupPostComment, error) {
	_, err := s.db.Exec(`
		UPDATE Groups SET content = ?, updated_at = CURRENT_TIMESTAMP
		WHERE public_id = ? AND user_id = ? AND type = 'comment'
	`, content, commentPublicID, userID)
	if err != nil {
		return nil, err
	}
	
	var comment models.GroupPostComment
	err = s.db.QueryRow(`
		SELECT public_id, group_id, user_id, content, COALESCE(image, ''), created_at, updated_at
		FROM Groups WHERE public_id = ? AND type = 'comment'
	`, commentPublicID).Scan(&comment.PublicID, &comment.GroupPostID, &comment.UserID, &comment.Content, &comment.Image, &comment.CreatedAt, &comment.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (s *groupPostStore) DeleteGroupPostComment(commentPublicID string, userID int64) error {
	_, err := s.db.Exec(`
		DELETE FROM Groups 
		WHERE public_id = ? AND user_id = ? AND type = 'comment'
	`, commentPublicID, userID)
	return err
}

func (s *groupPostStore) CanUserDeleteGroupContent(groupID string, userID int64) (bool, error) {
	// Check if user is group creator or admin
	var count int
	err := s.db.QueryRow(`
		SELECT COUNT(*) FROM (
			SELECT user_id FROM Groups WHERE id = ? AND user_id = ? AND type = 'group'
			UNION
			SELECT user_id FROM Groups WHERE group_id = ? AND user_id = ? AND type = 'member' AND role = 'admin' AND status = 'active'
		)
	`, groupID, userID, groupID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *groupPostStore) CanUserDeletePost(postID, userID int64) (bool, error) {
	// Check if user owns the post
	var ownerID int64
	var groupID string
	err := s.db.QueryRow("SELECT user_id, group_id FROM Groups WHERE id = ? AND type = 'post'", postID).Scan(&ownerID, &groupID)
	if err != nil {
		return false, err
	}
	if ownerID == userID {
		return true, nil
	}

	// Check if user is group admin
	return s.CanUserDeleteGroupContent(groupID, userID)
}

func (s *groupPostStore) CanUserDeleteComment(commentID, userID int64) (bool, error) {
	// Check if user owns the comment
	var ownerID int64
	var groupID string
	err := s.db.QueryRow(`
		SELECT c.user_id, c.group_id 
		FROM Groups c 
		WHERE c.id = ? AND c.type = 'comment'`, commentID).Scan(&ownerID, &groupID)
	if err != nil {
		return false, err
	}
	if ownerID == userID {
		return true, nil
	}

	// Check if user is group admin
	return s.CanUserDeleteGroupContent(groupID, userID)
}
