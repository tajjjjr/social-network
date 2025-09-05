package store

import (
	"database/sql"
	"time"

	"github.com/tajjjjr/social-network/backend/internal/models"
)

type GroupPostStore interface {
	CreateGroupPost(post *models.GroupPost) (*models.GroupPost, error)
	GetGroupPostByID(postID int64) (*models.GroupPost, error)
	GetGroupPosts(groupID int64, userID int64, limit, offset int) ([]*models.GroupPost, error)
	UpdateGroupPost(postID, userID int64, content string, imageData []byte, imageMimeType string) (*models.GroupPost, error)
	DeleteGroupPost(postID, userID int64) error
	CreateGroupPostComment(comment *models.GroupPostComment) (*models.GroupPostComment, error)
	GetGroupPostComments(postID int64, userID int64) ([]*models.GroupPostComment, error)
	UpdateGroupPostComment(commentID, userID int64, content string, imageData []byte, imageMimeType string) (*models.GroupPostComment, error)
	DeleteGroupPostComment(commentID, userID int64) error
	CanUserDeleteGroupContent(groupID, userID int64) (bool, error)
}

type groupPostStore struct {
	db *sql.DB
}

func NewGroupPostStore(db *sql.DB) GroupPostStore {
	return &groupPostStore{db: db}
}

func (s *groupPostStore) CreateGroupPost(post *models.GroupPost) (*models.GroupPost, error) {
	stmt, err := s.db.Prepare("INSERT INTO Group_Posts (group_id, user_id, content, image) VALUES (?, ?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(post.GroupID, post.UserID, post.Content, post.Image)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	post.ID = id
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	return post, nil
}

func (s *groupPostStore) GetGroupPostByID(postID int64) (*models.GroupPost, error) {
	var post models.GroupPost
	err := s.db.QueryRow(`
		SELECT id, group_id, user_id, content, image, like_count, dislike_count, created_at, updated_at 
		FROM Group_Posts WHERE id = ?`, postID).Scan(
		&post.ID, &post.GroupID, &post.UserID, &post.Content, &post.Image,
		&post.LikeCount, &post.DislikeCount, &post.CreatedAt, &post.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (s *groupPostStore) GetGroupPosts(groupID int64, userID int64, limit, offset int) ([]*models.GroupPost, error) {
	rows, err := s.db.Query(`
		SELECT gp.id, gp.group_id, gp.user_id, gp.content, gp.image, gp.like_count, gp.dislike_count, 
		       gp.created_at, gp.updated_at
		FROM Group_Posts gp
		WHERE gp.group_id = ?
		ORDER BY gp.created_at DESC
		LIMIT ? OFFSET ?`, groupID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.GroupPost
	for rows.Next() {
		var post models.GroupPost
		err := rows.Scan(&post.ID, &post.GroupID, &post.UserID, &post.Content, &post.Image,
			&post.LikeCount, &post.DislikeCount, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	return posts, nil
}

func (s *groupPostStore) UpdateGroupPost(postID, userID int64, content string, imageData []byte, imageMimeType string) (*models.GroupPost, error) {
	// Check if user owns the post
	var ownerID int64
	err := s.db.QueryRow("SELECT user_id FROM Group_Posts WHERE id = ?", postID).Scan(&ownerID)
	if err != nil {
		return nil, err
	}
	if ownerID != userID {
		return nil, sql.ErrNoRows
	}

	stmt, err := s.db.Prepare("UPDATE Group_Posts SET content = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(content, postID)
	if err != nil {
		return nil, err
	}

	return s.GetGroupPostByID(postID)
}

func (s *groupPostStore) DeleteGroupPost(postID, userID int64) error {
	// Check if user owns the post or is group admin
	canDelete, err := s.canUserDeletePost(postID, userID)
	if err != nil {
		return err
	}
	if !canDelete {
		return sql.ErrNoRows
	}

	_, err = s.db.Exec("DELETE FROM Group_Posts WHERE id = ?", postID)
	return err
}

func (s *groupPostStore) CreateGroupPostComment(comment *models.GroupPostComment) (*models.GroupPostComment, error) {
	stmt, err := s.db.Prepare("INSERT INTO Group_Post_Comments (group_post_id, user_id, parent_comment_id, content, image) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(comment.GroupPostID, comment.UserID, comment.ParentCommentID, comment.Content, comment.Image)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	comment.ID = id
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()
	return comment, nil
}

func (s *groupPostStore) GetGroupPostComments(postID int64, userID int64) ([]*models.GroupPostComment, error) {
	rows, err := s.db.Query(`
		SELECT id, group_post_id, user_id, parent_comment_id, content, image, like_count, dislike_count, created_at, updated_at
		FROM Group_Post_Comments 
		WHERE group_post_id = ?
		ORDER BY created_at ASC`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.GroupPostComment
	for rows.Next() {
		var comment models.GroupPostComment
		err := rows.Scan(&comment.ID, &comment.GroupPostID, &comment.UserID, &comment.ParentCommentID,
			&comment.Content, &comment.Image, &comment.LikeCount, &comment.DislikeCount,
			&comment.CreatedAt, &comment.UpdatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	return comments, nil
}

func (s *groupPostStore) UpdateGroupPostComment(commentID, userID int64, content string, imageData []byte, imageMimeType string) (*models.GroupPostComment, error) {
	// Check if user owns the comment
	var ownerID int64
	err := s.db.QueryRow("SELECT user_id FROM Group_Post_Comments WHERE id = ?", commentID).Scan(&ownerID)
	if err != nil {
		return nil, err
	}
	if ownerID != userID {
		return nil, sql.ErrNoRows
	}

	stmt, err := s.db.Prepare("UPDATE Group_Post_Comments SET content = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(content, commentID)
	if err != nil {
		return nil, err
	}

	var comment models.GroupPostComment
	err = s.db.QueryRow(`
		SELECT id, group_post_id, user_id, parent_comment_id, content, image, like_count, dislike_count, created_at, updated_at
		FROM Group_Post_Comments WHERE id = ?`, commentID).Scan(
		&comment.ID, &comment.GroupPostID, &comment.UserID, &comment.ParentCommentID,
		&comment.Content, &comment.Image, &comment.LikeCount, &comment.DislikeCount,
		&comment.CreatedAt, &comment.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (s *groupPostStore) DeleteGroupPostComment(commentID, userID int64) error {
	// Check if user owns the comment or is group admin
	canDelete, err := s.canUserDeleteComment(commentID, userID)
	if err != nil {
		return err
	}
	if !canDelete {
		return sql.ErrNoRows
	}

	_, err = s.db.Exec("DELETE FROM Group_Post_Comments WHERE id = ?", commentID)
	return err
}

func (s *groupPostStore) CanUserDeleteGroupContent(groupID, userID int64) (bool, error) {
	// Check if user is group creator
	var creatorID int64
	err := s.db.QueryRow("SELECT creator_id FROM Groups WHERE id = ?", groupID).Scan(&creatorID)
	if err != nil {
		return false, err
	}
	if creatorID == userID {
		return true, nil
	}

	// Check if user has admin permissions
	var count int
	err = s.db.QueryRow("SELECT COUNT(*) FROM Group_Permissions WHERE group_id = ? AND user_id = ? AND permission_type = 'admin'", groupID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *groupPostStore) canUserDeletePost(postID, userID int64) (bool, error) {
	// Check if user owns the post
	var ownerID, groupID int64
	err := s.db.QueryRow("SELECT user_id, group_id FROM Group_Posts WHERE id = ?", postID).Scan(&ownerID, &groupID)
	if err != nil {
		return false, err
	}
	if ownerID == userID {
		return true, nil
	}

	// Check if user is group admin
	return s.CanUserDeleteGroupContent(groupID, userID)
}

func (s *groupPostStore) canUserDeleteComment(commentID, userID int64) (bool, error) {
	// Check if user owns the comment
	var ownerID int64
	var groupID int64
	err := s.db.QueryRow(`
		SELECT gpc.user_id, gp.group_id 
		FROM Group_Post_Comments gpc 
		JOIN Group_Posts gp ON gpc.group_post_id = gp.id 
		WHERE gpc.id = ?`, commentID).Scan(&ownerID, &groupID)
	if err != nil {
		return false, err
	}
	if ownerID == userID {
		return true, nil
	}

	// Check if user is group admin
	return s.CanUserDeleteGroupContent(groupID, userID)
}