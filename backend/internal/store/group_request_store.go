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
	stmt, err := s.db.Prepare("INSERT INTO group_requests (group_id, user_id, status) VALUES (?, ?, ?)")
	if err != nil {
		return nil, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(request.GroupID, request.UserID, request.Status)
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

func (s *groupRequestStore) AddUserToGroup(groupID, userID int64) error {
	stmt, err := s.db.Prepare("INSERT OR IGNORE INTO Group_Members (group_id, user_id, is_accepted) VALUES (?, ?, 1)")
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(groupID, userID)
	if err != nil {
		return fmt.Errorf("error adding user to group: %w", err)
	}
	return nil
}
