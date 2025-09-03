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
	stmt, err := s.db.Prepare("INSERT INTO groups (creator_id, title, description, privacy) VALUES (?, ?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(group.CreatorID, group.Title, group.Description, group.Privacy)
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

func (s *groupStore) GetGroupByID(groupID int64) (*models.Group, error) {
	var group models.Group
	err := s.db.QueryRow("SELECT id, creator_id, title, description, privacy, created_at FROM groups WHERE id = ?", groupID).Scan(
		&group.ID,
		&group.CreatorID,
		&group.Title,
		&group.Description,
		&group.Privacy,
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
	rows, err := s.db.Query("SELECT id, creator_id, title, description, privacy, created_at FROM groups WHERE privacy = 'public' AND title LIKE ?", "%"+query+"%")
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
	rows, err := s.db.Query("SELECT id, creator_id, title, description, privacy, created_at FROM groups WHERE privacy = 'public' ORDER BY created_at DESC")
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
		SELECT g.id, g.creator_id, g.title, g.description, g.privacy, g.created_at 
		FROM groups g 
		JOIN group_members gm ON g.id = gm.group_id 
		WHERE gm.user_id = ? 
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
			&group.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		groups = append(groups, &group)
	}

	return groups, nil
}
