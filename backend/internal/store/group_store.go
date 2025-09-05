package store

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tajjjjr/social-network/backend/internal/models"
)

type groupStore struct {
	db *sql.DB
}

func NewGroupStore(db *sql.DB) GroupStore {
	return &groupStore{db: db}
}

func (s *groupStore) CreateGroup(group *models.Group) (*models.Group, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			fmt.Printf("Failed to rollback transaction: %v\n", err)
		}
	}()

	stmt, err := tx.Prepare("INSERT INTO Groups (creator_id, title, description, privacy, avatar) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(group.CreatorID, group.Title, group.Description, group.Privacy, group.Avatar)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	group.ID = id

	// Add creator as member
	memberStmt, err := tx.Prepare("INSERT INTO Group_Members (group_id, user_id, is_accepted) VALUES (?, ?, 1)")
	if err != nil {
		return nil, err
	}
	defer memberStmt.Close()

	_, err = memberStmt.Exec(group.ID, group.CreatorID)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return group, nil
}

func (s *groupStore) GetGroupByID(groupID int64) (*models.Group, error) {
	var group models.Group
	err := s.db.QueryRow("SELECT id, creator_id, title, description, privacy, avatar, created_at FROM Groups WHERE id = ?", groupID).Scan(
		&group.ID,
		&group.CreatorID,
		&group.Title,
		&group.Description,
		&group.Privacy,
		&group.Avatar,
		&group.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("group not found")
		}
		return nil, err
	}
	return &group, nil
}

func (s *groupStore) SearchPublicGroups(query string) ([]*models.Group, error) {
	rows, err := s.db.Query("SELECT id, creator_id, title, description, privacy, avatar, created_at FROM Groups WHERE title LIKE ? AND privacy = 'public'", "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*models.Group
	for rows.Next() {
		var group models.Group
		err := rows.Scan(
			&group.ID,
			&group.CreatorID,
			&group.Title,
			&group.Description,
			&group.Privacy,
			&group.Avatar,
			&group.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		groups = append(groups, &group)
	}

	return groups, nil
}

func (s *groupStore) GetAllPublicGroups() ([]*models.Group, error) {
	rows, err := s.db.Query("SELECT id, creator_id, title, description, privacy, avatar, created_at FROM Groups WHERE privacy = 'public' ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*models.Group
	for rows.Next() {
		var group models.Group
		err := rows.Scan(
			&group.ID,
			&group.CreatorID,
			&group.Title,
			&group.Description,
			&group.Privacy,
			&group.Avatar,
			&group.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		groups = append(groups, &group)
	}

	return groups, nil
}

func (s *groupStore) GetUserGroups(userID int64) ([]*models.Group, error) {
	rows, err := s.db.Query(`
		SELECT g.id, g.creator_id, g.title, g.description, g.privacy, g.avatar, g.created_at 
		FROM Groups g 
		JOIN Group_Members gm ON g.id = gm.group_id 
		WHERE gm.user_id = ? AND gm.is_accepted = 1
		ORDER BY g.created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*models.Group
	for rows.Next() {
		var group models.Group
		err := rows.Scan(
			&group.ID,
			&group.CreatorID,
			&group.Title,
			&group.Description,
			&group.Privacy,
			&group.Avatar,
			&group.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		groups = append(groups, &group)
	}

	return groups, nil
}
