package store

import (
	"database/sql"
	"github.com/tajjjjr/social-network/backend/internal/models"
)

type GroupMemberStore interface {
	GetGroupMembers(groupID int64) ([]*models.User, error)
	IsGroupMember(groupID, userID int64) (bool, error)
	AddGroupMember(groupID, userID int64, role string) (*models.GroupMember, error)
	RemoveGroupMember(groupID, userID int64) error
}

type groupMemberStore struct {
	db *sql.DB
}

func NewGroupMemberStore(db *sql.DB) GroupMemberStore {
	return &groupMemberStore{db: db}
}

func (s *groupMemberStore) GetGroupMembers(groupID int64) ([]*models.User, error) {
	rows, err := s.db.Query(`
		SELECT u.id, u.firstname, u.lastname, u.nickname, u.avatar, u.email
		FROM Users u
		JOIN Group_Members gm ON u.id = gm.user_id
		WHERE gm.group_id = ? AND gm.is_accepted = 1
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

func (s *groupMemberStore) IsGroupMember(groupID, userID int64) (bool, error) {
	var count int
	err := s.db.QueryRow(`
		SELECT COUNT(*) FROM Group_Members 
		WHERE group_id = ? AND user_id = ? AND is_accepted = 1
	`, groupID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *groupMemberStore) AddGroupMember(groupID, userID int64, role string) (*models.GroupMember, error) {
	stmt, err := s.db.Prepare("INSERT INTO Group_Members (group_id, user_id, role, is_accepted) VALUES (?, ?, ?, 1)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(groupID, userID, role)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &models.GroupMember{
		ID:      id,
		GroupID: groupID,
		UserID:  userID,
		Role:    role,
	}, nil
}

func (s *groupMemberStore) RemoveGroupMember(groupID, userID int64) error {
	_, err := s.db.Exec("DELETE FROM Group_Members WHERE group_id = ? AND user_id = ?", groupID, userID)
	return err
}