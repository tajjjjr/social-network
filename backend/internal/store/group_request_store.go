package store

import (
	"database/sql"
	"fmt"

	"github.com/tajjjjr/social-network/backend/internal/models"
)

type groupRequestStore struct {
	db *sql.DB
}

func NewGroupRequestStore(db *sql.DB) GroupRequestStore {
	return &groupRequestStore{db: db}
}

func (s *groupRequestStore) CreateGroupRequest(request *models.GroupRequest) (*models.GroupRequest, error) {
	// Get internal group ID from public_id
	var groupID int64
	err := s.db.QueryRow("SELECT id FROM Groups WHERE public_id = ? AND type = 'group'", request.GroupID).Scan(&groupID)
	if err != nil {
		return nil, fmt.Errorf("error finding group: %w", err)
	}

	result, err := s.db.Exec(`
		INSERT INTO Groups (type, group_id, user_id, status) 
		VALUES ('request', ?, ?, ?)`,
		groupID, request.UserID, request.Status)
	if err != nil {
		return nil, fmt.Errorf("error executing statement: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert ID: %w", err)
	}

	request.ID = id
	return request, nil
}

func (s *groupRequestStore) GetGroupRequestByID(requestID int64) (*models.GroupRequest, error) {
	var request models.GroupRequest
	err := s.db.QueryRow("SELECT id, group_id, user_id, status, created_at FROM group_requests WHERE id = ?", requestID).Scan(
		&request.ID,
		&request.GroupID,
		&request.UserID,
		&request.Status,
		&request.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("group request not found")
		}
		return nil, fmt.Errorf("error getting group request by ID: %w", err)
	}
	return &request, nil
}

func (s *groupRequestStore) UpdateGroupRequestStatus(requestID int64, status string) error {
	stmt, err := s.db.Prepare("UPDATE group_requests SET status = ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(status, requestID)
	if err != nil {
		return fmt.Errorf("error updating group request status: %w", err)
	}
	return nil
}

func (s *groupRequestStore) AddUserToGroupWithRole(groupPublicID string, userID int64, role string) error {
	// Get internal group ID from public_id
	var groupID int64
	err := s.db.QueryRow("SELECT id FROM Groups WHERE public_id = ? AND type = 'group'", groupPublicID).Scan(&groupID)
	if err != nil {
		return fmt.Errorf("error finding group: %w", err)
	}

	// Define clear roles: admin (creator) and member (default)
	if role != "admin" {
		role = "member"
	}

	// Clear SQL statement for joining group with proper role assignment
	stmt, err := s.db.Prepare(`
		INSERT OR REPLACE INTO Group_Members (group_id, user_id, role, is_accepted) 
		VALUES (?, ?, ?, 1)
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(groupID, userID, role)
	if err != nil {
		return fmt.Errorf("error adding user to group: %w", err)
	}
	return nil
}

func (s *groupRequestStore) AddUserToGroup(groupPublicID string, userID int64) error {
	return s.AddUserToGroupWithRole(groupPublicID, userID, "member")
}

func (s *groupRequestStore) IsUserMember(groupPublicID string, userID int64) (bool, error) {
	// Get internal group ID from public_id
	var groupID int64
	err := s.db.QueryRow("SELECT id FROM Groups WHERE public_id = ? AND type = 'group'", groupPublicID).Scan(&groupID)
	if err != nil {
		return false, fmt.Errorf("error finding group: %w", err)
	}

	var count int
	err = s.db.QueryRow(`
		SELECT COUNT(*) FROM (
			SELECT creator_id as user_id FROM Groups WHERE id = ? AND creator_id = ?
			UNION
			SELECT user_id FROM Group_Members WHERE group_id = ? AND user_id = ? AND is_accepted = 1
		)
	`, groupID, userID, groupID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *groupRequestStore) HasPendingRequest(groupPublicID string, userID int64) (bool, error) {
	// Get internal group ID from public_id
	var groupID int64
	err := s.db.QueryRow("SELECT id FROM Groups WHERE public_id = ? AND type = 'group'", groupPublicID).Scan(&groupID)
	if err != nil {
		return false, fmt.Errorf("error finding group: %w", err)
	}

	var count int
	err = s.db.QueryRow("SELECT COUNT(*) FROM group_requests WHERE group_id = ? AND user_id = ? AND status = 'pending'", groupID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
