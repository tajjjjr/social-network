package store

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/tajjjjr/social-network/backend/internal/models"
)

type GroupMemberStore interface {
	GetGroupMembers(publicID string) ([]*models.User, error)
	IsGroupMember(publicID string, userID int64) (bool, error)
	AddGroupMember(publicID string, userID int64, role string) (*models.GroupMember, error)
	RemoveGroupMember(publicID string, userID int64) error
}

type groupMemberStore struct {
	db *sql.DB
}

func NewGroupMemberStore(db *sql.DB) GroupMemberStore {
	return &groupMemberStore{db: db}
}

func (s *groupMemberStore) getGroupIDByPublicID(publicID string) (string, error) {
	var id string
	err := s.db.QueryRow("SELECT id FROM Groups WHERE public_id = ? AND type = 'group'", publicID).Scan(&id)
	return id, err
}

func (s *groupMemberStore) GetGroupMembers(publicID string) ([]*models.User, error) {
	groupID, err := s.getGroupIDByPublicID(publicID)
	if err != nil {
		return nil, err
	}
	rows, err := s.db.Query(`
		SELECT u.id, u.firstname, u.lastname, u.nickname, u.avatar, u.email
		FROM Users u
		JOIN Groups gm ON u.id = gm.user_id
		WHERE gm.group_id = ? AND gm.type = 'member' AND gm.status = 'active'
		ORDER BY u.firstname, u.lastname
	`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*models.User
	for rows.Next() {
		var member models.User
		err := rows.Scan(&member.ID, &member.FirstName, &member.LastName, &member.Nickname, &member.Avatar, &member.Email)
		if err != nil {
			return nil, err
		}
		members = append(members, &member)
	}

	return members, nil
}

func (s *groupMemberStore) IsGroupMember(publicID string, userID int64) (bool, error) {
	groupID, err := s.getGroupIDByPublicID(publicID)
	if err != nil {
		return false, err
	}
	var count int
	err = s.db.QueryRow(`
		SELECT COUNT(*) FROM Groups 
		WHERE group_id = ? AND user_id = ? AND type = 'member' AND status = 'active'
	`, groupID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *groupMemberStore) AddGroupMember(publicID string, userID int64, role string) (*models.GroupMember, error) {
	groupID, err := s.getGroupIDByPublicID(publicID)
	if err != nil {
		return nil, err
	}
	memberID := uuid.New().String()
	_, err = s.db.Exec("INSERT INTO Groups (id, type, group_id, user_id, role, status) VALUES (?, 'member', ?, ?, ?, 'active')", memberID, groupID, userID, role)
	if err != nil {
		return nil, err
	}

	return &models.GroupMember{
		ID:      1, // Dummy ID since we use string IDs internally
		GroupID: publicID,
		UserID:  userID,
		Role:    role,
	}, nil
}

func (s *groupMemberStore) RemoveGroupMember(publicID string, userID int64) error {
	groupID, err := s.getGroupIDByPublicID(publicID)
	if err != nil {
		return err
	}
	_, err = s.db.Exec("DELETE FROM Groups WHERE type = 'member' AND group_id = ? AND user_id = ?", groupID, userID)
	return err
}