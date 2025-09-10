package store

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/tajjjjr/social-network/backend/internal/models"
)

type groupStore struct {
	db *sql.DB
}

func NewGroupStore(db *sql.DB) GroupStore {
	return &groupStore{db: db}
}

func (s *groupStore) CreateGroup(group *models.Group) (*models.Group, error) {
	group.PublicID = uuid.New().String()
	result, err := s.db.Exec(`
		INSERT INTO Groups (public_id, type, user_id, title, content, role, privacy, image) 
		VALUES (?, 'group', ?, ?, ?, 'admin', ?, ?)`,
		group.PublicID, group.CreatorID, group.Title, group.Description, group.Privacy, group.Avatar)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	group.ID = id
	return group, nil
}

func (s *groupStore) GetGroupByID(publicID string) (*models.Group, error) {
	var group models.Group
	var avatar sql.NullString
	err := s.db.QueryRow(`
		SELECT id, public_id, user_id, title, content, privacy, COALESCE(image, ''), created_at 
		FROM Groups WHERE type = 'group' AND public_id = ?`, publicID).Scan(
		&group.ID, &group.PublicID, &group.CreatorID, &group.Title, &group.Description,
		&group.Privacy, &avatar, &group.CreatedAt)
	if err != nil {
		return nil, err
	}
	group.Avatar = avatar.String
	return &group, err
}

func (s *groupStore) SearchPublicGroups(query string) ([]*models.Group, error) {
	rows, err := s.db.Query(`
		SELECT id, COALESCE(public_id, id), COALESCE(user_id, 0), COALESCE(title, ''), COALESCE(content, ''), COALESCE(privacy, 'public'), COALESCE(image, ''), COALESCE(created_at, CURRENT_TIMESTAMP)
		FROM Groups 
		WHERE type = 'group' AND title LIKE ? AND COALESCE(privacy, 'public') = 'public'
		ORDER BY COALESCE(created_at, CURRENT_TIMESTAMP) DESC`,
		"%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*models.Group
	for rows.Next() {
		var group models.Group
		err := rows.Scan(&group.ID, &group.PublicID, &group.CreatorID, &group.Title, &group.Description,
			&group.Privacy, &group.Avatar, &group.CreatedAt)
		if err != nil {
			return nil, err
		}
		groups = append(groups, &group)
	}
	return groups, nil
}

func (s *groupStore) GetAllPublicGroups() ([]*models.Group, error) {
	rows, err := s.db.Query(`
		SELECT id, public_id, COALESCE(user_id, 0), COALESCE(title, ''), COALESCE(content, ''), COALESCE(privacy, 'public'), COALESCE(image, ''), COALESCE(created_at, CURRENT_TIMESTAMP)
		FROM Groups 
		WHERE type = 'group' AND COALESCE(privacy, 'public') = 'public'
		ORDER BY COALESCE(created_at, CURRENT_TIMESTAMP) DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*models.Group
	for rows.Next() {
		var group models.Group
		err := rows.Scan(&group.ID, &group.PublicID, &group.CreatorID, &group.Title, &group.Description,
			&group.Privacy, &group.Avatar, &group.CreatedAt)
		if err != nil {
			return nil, err
		}
		groups = append(groups, &group)
	}
	return groups, nil
}

func (s *groupStore) getGroupIDByPublicID(publicID string) (int64, error) {
	var id int64
	err := s.db.QueryRow("SELECT id FROM Groups WHERE public_id = ? AND type = 'group'", publicID).Scan(&id)
	return id, err
}

func (s *groupStore) JoinGroup(publicID string, userID int64) error {
	groupID, err := s.getGroupIDByPublicID(publicID)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`
		INSERT INTO Groups (type, group_id, user_id, role, status) 
		VALUES ('member', ?, ?, 'member', 'active')`,
		groupID, userID)
	return err
}

func (s *groupStore) LeaveGroup(publicID string, userID int64) error {
	groupID, err := s.getGroupIDByPublicID(publicID)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`DELETE FROM Groups WHERE type = 'member' AND group_id = ? AND user_id = ?`, groupID, userID)
	return err
}

func (s *groupStore) IsGroupMember(publicID string, userID int64) (bool, error) {
	groupID, err := s.getGroupIDByPublicID(publicID)
	if err != nil {
		return false, err
	}
	var count int
	err = s.db.QueryRow(`
		SELECT COUNT(*) FROM Groups 
		WHERE (type = 'group' AND id = ? AND user_id = ?) 
		OR (type = 'member' AND group_id = ? AND user_id = ? AND status = 'active')`,
		groupID, userID, groupID, userID).Scan(&count)
	return count > 0, err
}

func (s *groupStore) GetUserGroups(userID int64) ([]*models.Group, error) {
	rows, err := s.db.Query(`
		SELECT DISTINCT id, COALESCE(public_id, id), COALESCE(user_id, 0), COALESCE(title, ''), COALESCE(content, ''), COALESCE(privacy, 'public'), COALESCE(image, ''), COALESCE(created_at, CURRENT_TIMESTAMP)
		FROM Groups 
		WHERE (type = 'group' AND user_id = ?)
		OR (type = 'group' AND id IN (
			SELECT group_id FROM Groups WHERE type = 'member' AND user_id = ? AND status = 'active'
		))
		ORDER BY COALESCE(created_at, CURRENT_TIMESTAMP) DESC`, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*models.Group
	for rows.Next() {
		var group models.Group
		err := rows.Scan(&group.ID, &group.PublicID, &group.CreatorID, &group.Title, &group.Description,
			&group.Privacy, &group.Avatar, &group.CreatedAt)
		if err != nil {
			return nil, err
		}
		groups = append(groups, &group)
	}
	return groups, nil
}